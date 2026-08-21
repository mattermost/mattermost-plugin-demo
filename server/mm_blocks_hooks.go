package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

func (p *Plugin) handleMmBlocksIntegration(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	userName := request.UserName
	if userName == "" {
		userName = "unknown"
	}

	p.writeJSON(w, mmBlocksActionResponse{
		EphemeralText:    fmt.Sprintf("%s integration OK (user: %s).", mmBlocksResponsePrefix, userName),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationUpdate(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.decodePostActionRequest(w, r); !ok {
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Update: &model.Post{
			Message: "E2E mm_blocks post updated (message field).",
			Props: model.StringInterface{
				"mm_blocks": []any{
					map[string]any{
						"type": "text",
						"text": mmBlocksUpdatedMarker,
					},
				},
			},
		},
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationEchoQuery(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.decodePostActionRequest(w, r); !ok {
		return
	}

	query := r.URL.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]string, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, k+"="+strings.Join(query[k], ","))
	}

	p.writeJSON(w, mmBlocksActionResponse{
		EphemeralText:    fmt.Sprintf("%s query OK (%s)", mmBlocksResponsePrefix, strings.Join(entries, "&")),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationEchoContext(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	marker := contextValueAsString(request.Context, "test_marker")
	p.writeJSON(w, mmBlocksActionResponse{
		EphemeralText:    fmt.Sprintf("%s context OK (test_marker: %s).", mmBlocksResponsePrefix, marker),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationStaticSelect(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	label := contextValueAsString(request.Context, "selected_option")
	p.writeJSON(w, mmBlocksActionResponse{
		EphemeralText:    fmt.Sprintf("%s static_select OK (selected_option: %s).", mmBlocksResponsePrefix, label),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationEchoFormValues(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	summary := formatFormValuesSummary(getUpstreamFormValues(request.Context))
	p.writeJSON(w, mmBlocksActionResponse{
		EphemeralText:    fmt.Sprintf("%s form_values OK (%s)", mmBlocksResponsePrefix, summary),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksIntegrationLookup(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	searchText := contextString(request.Context, "query")
	if searchText == "" {
		searchText = r.URL.Query().Get("query")
	}

	p.writeJSON(w, model.LookupDialogResponse{
		Items: getMmBlocksLookupOptions(searchText),
	})
}

func (p *Plugin) handleMmBlocksDialogOpen(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	if request.TriggerId == "" {
		p.writeJSON(w, mmBlocksActionResponse{
			Type:  "ok",
			Error: "mm_blocks_dialog_open requires trigger_id",
		})
		return
	}

	dialog := getMmBlocksDialog(mmBlocksDialogOptions{
		Title:  "Demo Blocks (open)",
		Marker: contextString(request.Context, "marker"),
	})
	if err := p.openBlockDialog(request.TriggerId, dialog); err != nil {
		if p.API != nil {
			p.API.LogError("Failed to open mm_blocks dialog", "err", err.Error())
		}
		p.writeJSON(w, mmBlocksActionResponse{
			Type:  "ok",
			Error: err.Error(),
		})
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{Type: "ok"})
}

func (p *Plugin) handleMmBlocksDialogReturn(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	opts := mmBlocksDialogOptions{
		Marker: contextString(request.Context, "marker"),
	}
	scenario := contextString(request.Context, "scenario")
	if scenario == "" || scenario == "default" {
		opts.Title = "Demo Blocks (return)"
		p.writeJSON(w, mmBlocksActionResponse{
			Type:        "dialog",
			BlockDialog: getMmBlocksDialog(opts),
		})
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:        "dialog",
		BlockDialog: getMmBlocksDialogByScenario(scenario, opts),
	})
}

func (p *Plugin) handleMmBlocksDialogSubmit(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	summary := formatFormValuesSummary(getUpstreamFormValues(request.Context))
	step := ""
	if s := contextString(request.Context, "step"); s != "" {
		step = " step=" + s
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:             "ok",
		EphemeralText:    fmt.Sprintf("%s dialog submit OK%s (%s)", mmBlocksResponsePrefix, step, summary),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksDialogCancel(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	reason := contextString(request.Context, "reason")
	if reason == "" {
		reason = "cancel"
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:             "ok",
		EphemeralText:    fmt.Sprintf("%s dialog cancelled (reason=%s)", mmBlocksResponsePrefix, reason),
		SkipSlackParsing: true,
	})
}

func (p *Plugin) handleMmBlocksDialogRefresh(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	previousTitle := "Demo ticket"
	if title := contextString(getUpstreamFormValues(request.Context), "title"); title != "" {
		previousTitle = title
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:        "refresh",
		BlockDialog: getMmBlocksDialogStep2(previousTitle),
	})
}

func (p *Plugin) handleMmBlocksDialogFieldRefresh(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	formValues := getUpstreamFormValues(request.Context)
	p.writeJSON(w, mmBlocksActionResponse{
		Type:        "refresh",
		BlockDialog: getMmBlocksFieldRefreshDialog(contextString(formValues, "project_type"), contextString(formValues, "project_name")),
	})
}

func (p *Plugin) handleMmBlocksDialogMultistep(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	step := contextString(request.Context, "step")
	if step == "" {
		step = "1"
	}

	var blockDialog *mmBlocksDialog
	if step == "1" {
		blockDialog = getMmBlocksMultistep2Dialog()
	} else {
		blockDialog = getMmBlocksMultistep3Dialog()
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:        "refresh",
		BlockDialog: blockDialog,
	})
}

func (p *Plugin) handleMmBlocksDialogChild(w http.ResponseWriter, r *http.Request) {
	request, ok := p.decodePostActionRequest(w, r)
	if !ok {
		return
	}

	if request.TriggerId == "" {
		p.writeJSON(w, mmBlocksActionResponse{
			Type:           "ok",
			KeepDialogOpen: true,
			Error:          "mm_blocks_dialog_child requires trigger_id",
		})
		return
	}

	source := contextString(request.Context, "source")
	if source == "" {
		source = "Unknown"
	}

	if err := p.openBlockDialog(request.TriggerId, getMmBlocksChildContentDialog(source)); err != nil {
		if p.API != nil {
			p.API.LogError("Failed to open mm_blocks child dialog", "err", err.Error())
		}
		p.writeJSON(w, mmBlocksActionResponse{
			Type:           "ok",
			KeepDialogOpen: true,
			Error:          err.Error(),
		})
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{Type: "ok", KeepDialogOpen: true})
}

func (p *Plugin) handleMmBlocksDialogErrors(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.decodePostActionRequest(w, r); !ok {
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type: "ok",
		Errors: map[string]string{
			"title": "Title looks wrong",
			"email": "Email is invalid",
			"pick":  "Pick something else",
		},
	})
}

func (p *Plugin) handleMmBlocksDialogError(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.decodePostActionRequest(w, r); !ok {
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{
		Type:  "ok",
		Error: mmBlocksResponsePrefix + " dialog top-level error",
	})
}

func (p *Plugin) handleMmBlocksDialogGoto(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.decodePostActionRequest(w, r); !ok {
		return
	}

	p.writeJSON(w, mmBlocksActionResponse{
		GotoLocation: "/",
	})
}

func (p *Plugin) decodePostActionRequest(w http.ResponseWriter, r *http.Request) (*model.PostActionIntegrationRequest, bool) {
	defer r.Body.Close()

	var request model.PostActionIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil && !errors.Is(err, io.EOF) {
		if p.API != nil {
			p.API.LogError("Failed to decode PostActionIntegrationRequest", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}
	return &request, true
}

// openBlockDialog calls POST /api/v4/actions/dialogs/open so the server publishes
// the open_dialog websocket event (Path A). /mm_blocks_dialog_return remains the
// type:dialog response path (Path B).
func (p *Plugin) openBlockDialog(triggerID string, dialog *mmBlocksDialog) error {
	if p.API == nil {
		return nil
	}

	cfg := p.API.GetConfig()
	if cfg == nil || cfg.ServiceSettings.SiteURL == nil || *cfg.ServiceSettings.SiteURL == "" {
		return errors.New("ServiceSettings.SiteURL is not set")
	}

	return postOpenBlockDialog(&http.Client{Timeout: 15 * time.Second}, *cfg.ServiceSettings.SiteURL, triggerID, dialog)
}

func postOpenBlockDialog(client *http.Client, siteURL, triggerID string, dialog *mmBlocksDialog) error {
	body, err := json.Marshal(struct {
		TriggerID   string          `json:"trigger_id"`
		BlockDialog *mmBlocksDialog `json:"block_dialog"`
	}{
		TriggerID:   triggerID,
		BlockDialog: dialog,
	})
	if err != nil {
		return err
	}

	url := strings.TrimRight(siteURL, "/") + "/api/v4/actions/dialogs/open"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("dialogs/open returned %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
