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
	})

	for _, inc := range incidents {
		t.Run("triage/"+inc.ID, func(t *testing.T) {
			d := getDialogIncidentTriage(inc)
			require.NoError(t, d.IsValid())
		})

		t.Run("note/"+inc.ID, func(t *testing.T) {
			d := getDialogTimelineNote(inc.ID)
			require.NoError(t, d.IsValid())
		})
	}
}
