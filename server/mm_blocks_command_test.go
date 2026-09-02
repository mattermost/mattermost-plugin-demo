package main

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMmBlocksExamples(t *testing.T) {
	names := mmBlocksExampleNames()
	require.NotEmpty(t, names)
	assert.Contains(t, names, "basic")
	assert.Contains(t, names, "form_all_inputs_submit")
	assert.Contains(t, names, "incident_response")

	for _, name := range names {
		example := mmBlocksExamples[name]
		assert.NotEmpty(t, example.Message, name)
		assert.NotNil(t, example.Props["mm_blocks"], name)
	}

	p := &Plugin{}
	list := p.executeMmBlocksExample(&model.CommandArgs{}, "")
	assert.Equal(t, model.CommandResponseTypeEphemeral, list.ResponseType)
	assert.Contains(t, list.Text, "`basic`")

	unknown := p.executeMmBlocksExample(&model.CommandArgs{}, "not-a-real-example")
	assert.Equal(t, model.CommandResponseTypeEphemeral, unknown.ResponseType)
	assert.Contains(t, unknown.Text, "Unknown mm_blocks example")

	help := p.executeCommandMmBlocks(&model.CommandArgs{Command: "/mm_blocks help"})
	assert.Contains(t, help.Text, "/mm_blocks example")
	assert.Contains(t, help.Text, "`form_text_inputs_submit`")
}

func TestEnsureMmBlocksExampleActions(t *testing.T) {
	t.Run("fills missing actions for interactive buttons", func(t *testing.T) {
		props := cloneExampleProps("basic")
		require.Nil(t, props["mm_blocks_actions"])

		ensureMmBlocksExampleActions(props)

		actions := props["mm_blocks_actions"].(map[string]any)
		assert.Contains(t, actions, "e2e_mm_blocks_primary")
		assert.Contains(t, actions, "e2e_mm_blocks_secondary")
	})

	t.Run("skips disabled buttons", func(t *testing.T) {
		props := cloneExampleProps("button_tooltip_disabled")
		ensureMmBlocksExampleActions(props)

		actions := props["mm_blocks_actions"].(map[string]any)
		assert.Contains(t, actions, "e2e_mm_blocks_run")
		assert.NotContains(t, actions, "e2e_mm_blocks_na")
	})

	t.Run("keeps fixture actions that already exist", func(t *testing.T) {
		props := cloneExampleProps("form_text_inputs_submit")
		original := props["mm_blocks_actions"].(map[string]any)["e2e_mm_blocks_form_submit"]

		ensureMmBlocksExampleActions(props)

		assert.Equal(t, original, props["mm_blocks_actions"].(map[string]any)["e2e_mm_blocks_form_submit"])
	})

	t.Run("fills nested static_select in collapsible", func(t *testing.T) {
		props := cloneExampleProps("collapsible_select")
		ensureMmBlocksExampleActions(props)

		actions := props["mm_blocks_actions"].(map[string]any)
		assert.Contains(t, actions, "e2e_mm_blocks_select")
	})

	t.Run("leaves non-interactive fixtures without actions", func(t *testing.T) {
		props := cloneExampleProps("image_size_presets")
		ensureMmBlocksExampleActions(props)
		assert.Nil(t, props["mm_blocks_actions"])
	})

	t.Run("incident response fills restart and resolve actions", func(t *testing.T) {
		props := cloneExampleProps("incident_response")
		ensureMmBlocksExampleActions(props)

		actions := props["mm_blocks_actions"].(map[string]any)
		assert.Contains(t, actions, "incident_resolve")
		assert.Contains(t, actions, "incident_restart_api_prod_2a")
		assert.Contains(t, actions, "incident_filter_env")
		assert.NotContains(t, actions, "incident_restart_worker_jobs_1a")
	})
}

func TestIncidentResponseStaysUnderClientBlockLimit(t *testing.T) {
	blocks, ok := mmBlocksExamples["incident_response"].Props["mm_blocks"].([]any)
	require.True(t, ok)
	assert.LessOrEqual(t, countMmBlocks(blocks), mmBlocksClientMaxTotal)
}

func cloneExampleProps(name string) model.StringInterface {
	example := mmBlocksExamples[name]
	props := model.StringInterface{}
	for k, v := range example.Props {
		props[k] = v
	}
	return props
}
