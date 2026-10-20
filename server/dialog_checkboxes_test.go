package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The checkbox dialogs are consumed by the mobile e2e suite, so they must pass
// the same server-side validation the real /dialogs/open endpoint applies.
// Without this, a bad default or a label_position on the wrong element type only
// surfaces as an opaque "Failed to open Interactive Dialog" at test time.
func TestCheckboxDialogsAreValid(t *testing.T) {
	t.Run("checkbox_group", func(t *testing.T) {
		d := getDialogCheckboxGroup()
		require.NoError(t, d.IsValid())
	})

	t.Run("checkbox_matrix", func(t *testing.T) {
		d := getDialogCheckboxMatrix()
		require.NoError(t, d.IsValid())
	})
}

// Field names and option values are part of the e2e contract — the spec builds
// testIDs like AppFormElement.<name>.checkbox.<value>.button from them.
func TestCheckboxDialogContract(t *testing.T) {
	group := getDialogCheckboxGroup()

	byName := map[string]int{}
	for i, e := range group.Elements {
		byName[e.Name] = i
	}
	for _, name := range []string{"services", "regions", "notify_before", "notify_after", "priority"} {
		require.Contains(t, byName, name, "e2e spec depends on this field name")
	}

	services := group.Elements[byName["services"]]
	assert.False(t, services.Optional, "services must stay required for the validation test")
	assert.Empty(t, services.Default, "services must have no default so submit is blocked")

	regions := group.Elements[byName["regions"]]
	assert.Equal(t, "us,apac", regions.Default, "defaults test asserts these are pre-checked")

	assert.Equal(t, "before", group.Elements[byName["notify_before"]].LabelPosition)
	assert.Equal(t, "after", group.Elements[byName["notify_after"]].LabelPosition)

	priority := group.Elements[byName["priority"]]
	assert.Equal(t, "radio", priority.Type)
	assert.True(t, priority.Optional, "Clear selection only renders for optional radios")

	matrix := getDialogCheckboxMatrix()
	matrixByName := map[string]int{}
	for i, e := range matrix.Elements {
		matrixByName[e.Name] = i
	}
	require.Contains(t, matrixByName, "permissions")
	require.Contains(t, matrixByName, "environment_owner")

	permissions := matrix.Elements[matrixByName["permissions"]]
	assert.Equal(t, "multiple", permissions.MatrixConfig.RowSelection)
	assert.Equal(t, "posts:view,edit", permissions.Default)
	assert.Greater(t, len(permissions.MatrixConfig.Columns), 4, "needs enough columns to force horizontal scroll")

	owner := matrix.Elements[matrixByName["environment_owner"]]
	assert.Equal(t, "single", owner.MatrixConfig.RowSelection)
}

func TestFormatCheckboxValue(t *testing.T) {
	for name, tc := range map[string]struct {
		value    any
		expected string
	}{
		"checkbox_group sorted":      {[]any{"web", "api"}, "api,web"},
		"checkbox_group single":      {[]any{"api"}, "api"},
		"empty array":                {[]any{}, "(none)"},
		"nil":                        {nil, "(none)"},
		"empty string":               {"", "(none)"},
		"radio string":               {"p2", "p2"},
		"matrix single row":          {[]any{"posts:edit,view"}, "posts:edit,view"},
		"matrix sorts columns":       {[]any{"posts:view,edit"}, "posts:edit,view"},
		"matrix sorts rows":          {[]any{"files:delete", "posts:view"}, "files:delete;posts:view"},
		"matrix trims column spaces": {[]any{"posts: view , edit "}, "posts:edit,view"},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, formatCheckboxValue(tc.value))
		})
	}
}
