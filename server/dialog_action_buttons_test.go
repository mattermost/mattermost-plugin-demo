package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The incident dialogs must pass the same server-side validation the real
// /dialogs/open endpoint applies. Both Dialog.Title and DialogElement.DisplayName
// are capped at 24 *bytes*, and the emoji-prefixed labels here sit close to that
// cap — exceeding it only surfaces as an opaque "Failed to open Interactive
// Dialog" at demo/e2e time, with no hint as to which field is at fault.
func TestIncidentDialogsAreValid(t *testing.T) {
	t.Run("board", func(t *testing.T) {
		d := getDialogIncidentBoard()
		require.NoError(t, d.IsValid())

		// The board refreshes when a stacked triage dialog closes.
		require.True(t, d.RefreshOnReturn)
	})

	for _, inc := range incidents {
		t.Run("triage/"+inc.ID, func(t *testing.T) {
			d := getDialogIncidentTriage(inc, "")
			require.NoError(t, d.IsValid())

			// Triage refreshes when a stacked timeline note dialog closes.
			require.True(t, d.RefreshOnReturn)
		})

		t.Run("note/"+inc.ID, func(t *testing.T) {
			d := getDialogTimelineNote(inc.ID, nil)
			require.NoError(t, d.IsValid())
		})
	}
}

func TestTruncateLine(t *testing.T) {
	for name, tc := range map[string]struct {
		in       string
		maxRunes int
		want     string
	}{
		"short":              {"rettwert", 40, "rettwert"},
		"exact length":       {"abcde", 5, "abcde"},
		"cut with ellipsis":  {"abcdefgh", 5, "abcd…"},
		"newlines collapsed": {"line one\n\nline   two", 5, "line…"},
		"multibyte":          {"ééééééé", 5, "éééé…"},
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.want, truncateLine(tc.in, tc.maxRunes))
		})
	}
}

func TestIncidentTriageListsTimelineEntries(t *testing.T) {
	inc := incidents[0]

	incidentsMu.Lock()
	saved := incidentTimelines[inc.ID]
	incidentTimelines[inc.ID] = []timelineEntry{
		{Time: "22:40 UTC", Type: "update", Message: "Failover started"},
		{Time: "22:50 UTC", Type: "root_cause", Message: "Replica disk filled up because *binlog* retention was set far too high"},
	}
	incidentsMu.Unlock()
	t.Cleanup(func() {
		incidentsMu.Lock()
		incidentTimelines[inc.ID] = saved
		incidentsMu.Unlock()
	})

	intro := getDialogIncidentTriage(inc, "").IntroductionText

	require.Contains(t, intro, "- Timeline: 2 entries\n")
	require.Contains(t, intro, "\n  - 22:40 UTC, update: Failover started")

	// Long messages are truncated to one line, then Markdown-escaped.
	require.Contains(t, intro, "\n  - 22:50 UTC, root_cause: Replica disk filled up because \\*binlog\\*…")
}
