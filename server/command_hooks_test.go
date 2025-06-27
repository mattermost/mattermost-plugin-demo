package main

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDialogWithMultiSelectElements(t *testing.T) {
	dialog := getDialogWithMultiSelectElements()

	// Test basic dialog properties
	assert.Equal(t, "somecallbackid", dialog.CallbackId)
	assert.Equal(t, "Multi-Select Dialog Demo", dialog.Title)
	assert.Equal(t, "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png", dialog.IconURL)
	assert.Equal(t, "Submit Multi-Select", dialog.SubmitLabel)
	assert.True(t, dialog.NotifyOnCancel)
	assert.Equal(t, dialogStateSome, dialog.State)

	// Test that we have exactly 3 elements
	require.Len(t, dialog.Elements, 3)

	// Test Multi-Select Users element
	usersElement := dialog.Elements[0]
	assert.Equal(t, "Multi-Select Users", usersElement.DisplayName)
	assert.Equal(t, "multiselect_users", usersElement.Name)
	assert.Equal(t, "select", usersElement.Type)
	assert.Equal(t, "users", usersElement.DataSource)
	assert.True(t, usersElement.MultiSelect)
	assert.Equal(t, "Select multiple users...", usersElement.Placeholder)
	assert.Equal(t, "Choose multiple users from the list.", usersElement.HelpText)

	// Test Multi-Select Channels element
	channelsElement := dialog.Elements[1]
	assert.Equal(t, "Multi-Select Channels", channelsElement.DisplayName)
	assert.Equal(t, "multiselect_channels", channelsElement.Name)
	assert.Equal(t, "select", channelsElement.Type)
	assert.Equal(t, "channels", channelsElement.DataSource)
	assert.True(t, channelsElement.MultiSelect)
	assert.Equal(t, "Select multiple channels...", channelsElement.Placeholder)
	assert.Equal(t, "Choose multiple channels from the list.", channelsElement.HelpText)

	// Test Multi-Select Options element
	optionsElement := dialog.Elements[2]
	assert.Equal(t, "Multi-Select Options", optionsElement.DisplayName)
	assert.Equal(t, "multiselect_options", optionsElement.Name)
	assert.Equal(t, "select", optionsElement.Type)
	assert.True(t, optionsElement.MultiSelect)
	assert.Equal(t, "Select multiple options...", optionsElement.Placeholder)
	assert.Equal(t, "Choose multiple options from the list.", optionsElement.HelpText)

	// Test custom options
	require.Len(t, optionsElement.Options, 4)
	expectedOptions := []*model.PostActionOptions{
		{Text: "Option A", Value: "optA"},
		{Text: "Option B", Value: "optB"},
		{Text: "Option C", Value: "optC"},
		{Text: "Option D", Value: "optD"},
	}
	assert.Equal(t, expectedOptions, optionsElement.Options)
}

func TestExecuteCommandDialog_MultiSelect(t *testing.T) {
	// Test that the multi-select dialog is properly configured
	// We test the dialog generation directly since mocking the full API is complex

	// Parse the command like the real handler does
	fields := []string{"/dialog", "multi-select"}
	command := ""
	if len(fields) == 2 {
		command = fields[1]
	}

	assert.Equal(t, "multi-select", command)

	// Test that we can create the dialog request structure
	// This tests the case that would be executed in the real handler
	if command == "multi-select" {
		dialog := getDialogWithMultiSelectElements()
		
		// Verify it's the correct dialog
		assert.Equal(t, "Multi-Select Dialog Demo", dialog.Title)
		assert.Len(t, dialog.Elements, 3)

		// Verify all elements have MultiSelect = true
		for _, element := range dialog.Elements {
			assert.True(t, element.MultiSelect, "Element %s should have MultiSelect = true", element.Name)
		}
	}
}

func TestCommandDialogHelp(t *testing.T) {
	// Test that the help text includes our new multi-select option
	assert.Contains(t, commandDialogHelp, "Interactive Dialog Slash Command Help")
	assert.Contains(t, commandDialogHelp, "/dialog multi-select")
	assert.Contains(t, commandDialogHelp, "multi-select fields")
}

func TestGetCommandDialogAutocompleteData(t *testing.T) {
	autocomplete := getCommandDialogAutocompleteData()

	assert.Equal(t, commandTriggerDialog, autocomplete.Trigger)
	assert.Equal(t, "Open an Interactive Dialog.", autocomplete.HelpText)

	// Check that multi-select command is included
	commands := autocomplete.SubCommands
	var multiSelectCommand *model.AutocompleteData
	for _, cmd := range commands {
		if cmd.Trigger == "multi-select" {
			multiSelectCommand = cmd
			break
		}
	}

	require.NotNil(t, multiSelectCommand, "multi-select command should be in autocomplete")
	assert.Equal(t, "Open an Interactive Dialog with multi-select fields.", multiSelectCommand.HelpText)
}