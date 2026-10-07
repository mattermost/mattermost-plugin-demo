package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

type sortie struct {
	ID           string
	Aircraft     string
	Origin       string
	Status       string // "airborne", "scheduled", "landed", "complete"
	DepartureUTC string
	CargoLbs     int
	Waypoints    []cargoWaypoint
}

type cargoWaypoint struct {
	ID       string
	Location string
	ETA      string
	Tasks    []cargoTask
}

type cargoTask struct {
	ID        string
	Name      string
	Completed bool
}

var sorties = []sortie{
	{
		ID: "S-101", Aircraft: "C-17A", Origin: "Ramstein AB",
		Status: "airborne", DepartureUTC: "06:30Z", CargoLbs: 42000,
		Waypoints: []cargoWaypoint{
			{ID: "W-1", Location: "FOB Alpha", ETA: "08:45Z", Tasks: []cargoTask{
				{ID: "T-101-1", Name: "Offload medical supplies", Completed: true},
				{ID: "T-101-2", Name: "Verify pallet count", Completed: false},
				{ID: "T-101-3", Name: "Refuel aircraft", Completed: false},
			}},
			{ID: "W-2", Location: "FOB Bravo", ETA: "10:15Z", Tasks: []cargoTask{
				{ID: "T-101-4", Name: "Offload equip containers", Completed: false},
				{ID: "T-101-5", Name: "Pick up return cargo", Completed: false},
			}},
		},
	},
	{
		ID: "S-102", Aircraft: "KC-135R", Origin: "Incirlik AB",
		Status: "scheduled", DepartureUTC: "14:00Z", CargoLbs: 18500,
		Waypoints: []cargoWaypoint{
			{ID: "W-3", Location: "Camp Delta", ETA: "16:30Z", Tasks: []cargoTask{
				{ID: "T-102-1", Name: "Deliver fuel bladders", Completed: false},
				{ID: "T-102-2", Name: "Brief ground crew", Completed: false},
			}},
		},
	},
	{
		ID: "S-103", Aircraft: "C-130J", Origin: "Aviano AB",
		Status: "landed", DepartureUTC: "03:15Z", CargoLbs: 28750,
		Waypoints: []cargoWaypoint{
			{ID: "W-4", Location: "FOB Charlie", ETA: "05:00Z", Tasks: []cargoTask{
				{ID: "T-103-1", Name: "Offload ammunition", Completed: true},
				{ID: "T-103-2", Name: "Manifest verification", Completed: true},
				{ID: "T-103-3", Name: "Coordinate with QM", Completed: false},
			}},
		},
	},
}

func sortieStatusIcon(status string) string {
	switch status {
	case "airborne":
		return "✈️" // ✈️
	case "scheduled":
		return "\U0001f550" // 🕐
	case "landed":
		return "\U0001f7e2" // 🟢
	case "complete":
		return "✅" // ✅
	default:
		return "⬜"
	}
}

func findSortie(id string) (*sortie, int) {
	for i := range sorties {
		if sorties[i].ID == id {
			return &sorties[i], i
		}
	}
	return nil, -1
}

func findCargoWaypoint(s *sortie, id string) (*cargoWaypoint, int) {
	for i := range s.Waypoints {
		if s.Waypoints[i].ID == id {
			return &s.Waypoints[i], i
		}
	}
	return nil, -1
}

func taskProgress(wp cargoWaypoint) (done, total int) {
	total = len(wp.Tasks)
	for _, t := range wp.Tasks {
		if t.Completed {
			done++
		}
	}
	return
}

func formatCargoWeight(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	cut := len(s) - 3
	return s[:cut] + "," + s[cut:]
}

// Dialog A: Cargo Operations Board
func getDialogCargoBoard() model.Dialog {
	counts := map[string]int{}
	for _, s := range sorties {
		counts[s.Status]++
	}
	parts := []string{}
	for _, status := range []string{"airborne", "scheduled", "landed", "complete"} {
		if c := counts[status]; c > 0 {
			parts = append(parts, fmt.Sprintf("%s %d %s", sortieStatusIcon(status), c, status))
		}
	}
	intro := fmt.Sprintf("**%d Active Sorties** — %s\n\n", len(sorties), strings.Join(parts, " · "))
	intro += "| Sortie | Aircraft | Origin | Departure | Cargo | Status |\n|---|---|---|---|---|---|\n"
	for _, s := range sorties {
		intro += fmt.Sprintf("| **%s** | %s | %s | %s | %s lbs | %s %s |\n",
			s.ID, s.Aircraft, s.Origin, s.DepartureUTC,
			formatCargoWeight(s.CargoLbs),
			sortieStatusIcon(s.Status), s.Status,
		)
	}
	intro += "\nSelect a sortie to view its waypoint schedule."

	elements := make([]model.DialogElement, 0, len(sorties))
	for _, s := range sorties {
		elements = append(elements, model.DialogElement{
			DisplayName: fmt.Sprintf("%s %s", sortieStatusIcon(s.Status), s.ID),
			Name:        "btn_" + s.ID,
			Type:        "action_button",
			ActionButton: &model.DialogActionButton{
				URL:     fmt.Sprintf("/plugins/%s/dialog/cargo/sortie", manifest.Id),
				Context: map[string]string{"sortie_id": s.ID},
			},
		})
	}

	return model.Dialog{
		CallbackId:       "cargo_board",
		Title:            "Cargo Operations",
		IntroductionText: intro,
		SubmitLabel:      "Close",
		NotifyOnCancel:   false,
		State:            "cargo_board",
		Elements:         elements,
	}
}

// Dialog B: Sortie Detail — waypoint progress table + select to navigate to Dialog C.
// Submits to /dialog/cargo/tasks; the handler returns Dialog C as a form replacement,
// keeping the stack at depth 2 so Dialog C's save can return a fresh Dialog B.
func getDialogSortieDetail(sortieID string) model.Dialog {
	s, _ := findSortie(sortieID)
	if s == nil {
		return model.Dialog{Title: "Error", IntroductionText: "Sortie not found."}
	}

	wayNames := make([]string, 0, len(s.Waypoints))
	for _, wp := range s.Waypoints {
		wayNames = append(wayNames, wp.Location)
	}
	intro := fmt.Sprintf("**%s** · %s · %s → %s\n",
		s.ID, s.Aircraft, s.Origin, strings.Join(wayNames, " → "))
	intro += fmt.Sprintf("Departure: **%s** · Status: %s **%s** · Cargo: **%s lbs**\n\n",
		s.DepartureUTC, sortieStatusIcon(s.Status), s.Status, formatCargoWeight(s.CargoLbs))

	intro += "| Waypoint | ETA | Progress |\n|---|---|---|\n"
	for _, wp := range s.Waypoints {
		done, total := taskProgress(wp)
		intro += fmt.Sprintf("| \U0001f4cd %s | %s | %d/%d ✅ |\n",
			wp.Location, wp.ETA, done, total)
	}

	options := make([]*model.PostActionOptions, 0, len(s.Waypoints))
	for _, wp := range s.Waypoints {
		done, total := taskProgress(wp)
		options = append(options, &model.PostActionOptions{
			Text:  fmt.Sprintf("\U0001f4cd %s — ETA %s (%d/%d ✅)", wp.Location, wp.ETA, done, total),
			Value: wp.ID,
		})
	}

	return model.Dialog{
		CallbackId:       "cargo_sortie_" + s.ID,
		Title:            s.ID + " — " + s.Aircraft,
		IntroductionText: intro,
		SubmitLabel:      "View Tasks",
		NotifyOnCancel:   true,
		State:            sortieID,
		Elements: []model.DialogElement{{
			DisplayName: "Waypoint",
			Name:        "waypoint_id",
			Type:        "select",
			Placeholder: "Select a waypoint...",
			Options:     options,
		}},
	}
}

// Dialog C: Waypoint Task Checklist — bool field per task
func getDialogLocationTasks(sortieID, waypointID string) model.Dialog {
	s, _ := findSortie(sortieID)
	if s == nil {
		return model.Dialog{Title: "Error", IntroductionText: "Sortie not found."}
	}
	wp, _ := findCargoWaypoint(s, waypointID)
	if wp == nil {
		return model.Dialog{Title: "Error", IntroductionText: "Waypoint not found."}
	}

	done, total := taskProgress(*wp)
	intro := fmt.Sprintf("**%s** · Sortie %s · ETA %s\n%d of %d tasks complete",
		wp.Location, s.ID, wp.ETA, done, total)

	elements := make([]model.DialogElement, 0, len(wp.Tasks))
	for _, t := range wp.Tasks {
		defaultVal := "false"
		if t.Completed {
			defaultVal = "true"
		}
		elements = append(elements, model.DialogElement{
			DisplayName: t.Name,
			Name:        t.ID,
			Type:        "bool",
			Default:     defaultVal,
			Optional:    true,
		})
	}

	// "FOB Charlie — Tasks" = 19 bytes, within 24-byte limit
	return model.Dialog{
		CallbackId:       "cargo_tasks",
		Title:            wp.Location + " — Tasks",
		IntroductionText: intro,
		SubmitLabel:      "Save",
		NotifyOnCancel:   true,
		State:            sortieID + "::" + waypointID,
		Elements:         elements,
	}
}
