package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

const (
	commandMmBlocksHelp = "###### mm_blocks Slash Command Help\n" +
		"- `/mm_blocks` - Show this help text\n" +
		"- `/mm_blocks demo` - Post a message with mm_blocks buttons, a select, and dialog openers.\n" +
		"- `/mm_blocks integration` - Post a button that returns an ephemeral integration response.\n" +
		"- `/mm_blocks update` - Post a button that updates the post in place.\n" +
		"- `/mm_blocks query` - Post a button whose action URL includes query parameters.\n" +
		"- `/mm_blocks context` - Post a button that echoes static action context.\n" +
		"- `/mm_blocks select` - Post a static select that echoes the selected option.\n" +
		"- `/mm_blocks form` - Post form fields that echo `context.form_values`.\n" +
		"- `/mm_blocks lookup` - Post a dynamic select whose options come from this plugin.\n" +
		"- `/mm_blocks dialog [scenario]` - Post a button that opens a blocks dialog (type:dialog).\n" +
		"- `/mm_blocks dialog-open` - Post a button that opens a blocks dialog via the open path.\n" +
		"- `/mm_blocks example [name]` - Post a fixture mm_blocks example. Omit the name to list available examples.\n" +
		"- `/mm_blocks help` - Show this help text\n\n" +
		"Dialog scenarios: `simple`, `full`, `boolean`, `users_channels`, `multiselect`, `multiselect_defaults`, `dynamic`, `empty_required`, `file_upload`, `field_refresh`, `multistep_1`, `action_parent`, `datetime_basic`, `datetime_mindate`, `datetime_interval`, `datetime_relative`, `datetime_timezone`, `datetime_manual`"
)

func mmBlocksHelpText() string {
	return commandMmBlocksHelp + "\n\nExamples: `" + strings.Join(mmBlocksExampleNames(), "`, `") + "`"
}

func getCommandMmBlocksAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData(commandTriggerMmBlocks, "", "Post mm_blocks interactive messages and dialogs.")

	for _, item := range []struct {
		name, help string
	}{
		{"demo", "Post a message with mm_blocks buttons, a select, and dialog openers."},
		{"integration", "Post a button that returns an ephemeral integration response."},
		{"update", "Post a button that updates the post in place."},
		{"query", "Post a button whose action URL includes query parameters."},
		{"context", "Post a button that echoes static action context."},
		{"select", "Post a static select that echoes the selected option."},
		{"form", "Post form fields that echo context.form_values."},
		{"lookup", "Post a dynamic select whose options come from this plugin."},
		{"dialog-open", "Post a button that opens a blocks dialog via the open path."},
		{"help", "Show mm_blocks command help."},
	} {
		command.AddCommand(model.NewAutocompleteData(item.name, "", item.help))
	}

	dialog := model.NewAutocompleteData("dialog", "[scenario]", "Post a button that opens a blocks dialog.")
	for _, scenario := range mmBlocksDialogScenarios() {
		dialog.AddCommand(model.NewAutocompleteData(scenario, "", "Open the "+scenario+" blocks dialog."))
	}
	command.AddCommand(dialog)

	example := model.NewAutocompleteData("example", "[name]", "Post a fixture mm_blocks example.")
	for _, name := range mmBlocksExampleNames() {
		example.AddCommand(model.NewAutocompleteData(name, "", mmBlocksExamples[name].Message))
	}
	command.AddCommand(example)

	return command
}

func (p *Plugin) executeCommandMmBlocks(args *model.CommandArgs) *model.CommandResponse {
	fields := strings.Fields(args.Command)
	command := ""
	if len(fields) >= 2 {
		command = fields[1]
	}
	scenario := ""
	if len(fields) >= 3 {
		scenario = fields[2]
	}

	switch command {
	case "", "help":
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         mmBlocksHelpText(),
		}
	case "demo":
		return p.postMmBlocksDemo(args)
	case "integration":
		return p.postMmBlocksMessage(args, "Click to ping the demo plugin mm_blocks integration.",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_integration` and shows an ephemeral response."},
				map[string]any{"type": "button", "text": "Ping integration", "style": "primary", "action_id": "demo_mm_blocks_integration"},
			},
			map[string]any{
				"demo_mm_blocks_integration": mmBlocksPostAction("/mm_blocks_integration", map[string]any{}),
			},
		)
	case "update":
		return p.postMmBlocksMessage(args, "Click to update this post in place.",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_integration_update`."},
				map[string]any{"type": "button", "text": "Update post", "style": "primary", "action_id": "demo_mm_blocks_update"},
			},
			map[string]any{
				"demo_mm_blocks_update": mmBlocksPostAction("/mm_blocks_integration_update", map[string]any{}),
			},
		)
	case "query":
		return p.postMmBlocksMessage(args, "Click to echo query parameters from the action URL.",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_integration_echo_query?source=plugin&kind=button`."},
				map[string]any{"type": "button", "text": "Echo query", "style": "primary", "action_id": "demo_mm_blocks_query"},
			},
			map[string]any{
				"demo_mm_blocks_query": map[string]any{
					"type": "external",
					"url":  pluginPath("/mm_blocks_integration_echo_query") + "?source=plugin&kind=button",
				},
			},
		)
	case "context":
		return p.postMmBlocksMessage(args, "Click to echo static action context.",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_integration_echo_context` with `test_marker`."},
				map[string]any{"type": "button", "text": "Echo context", "style": "primary", "action_id": "demo_mm_blocks_context"},
			},
			map[string]any{
				"demo_mm_blocks_context": mmBlocksPostAction("/mm_blocks_integration_echo_context", map[string]any{"test_marker": "plugin-context"}),
			},
		)
	case "select":
		return p.postMmBlocksMessage(args, "Choose an option to echo `selected_option`.",
			[]any{
				map[string]any{"type": "text", "text": "Static select backed by `/mm_blocks_integration_static_select`."},
				map[string]any{
					"type":        "static_select",
					"placeholder": "Pick an option",
					"action_id":   "demo_mm_blocks_select",
					"options": []any{
						map[string]any{"text": "Alpha", "value": "opt_alpha"},
						map[string]any{"text": "Beta", "value": "opt_beta"},
						map[string]any{"text": "Gamma", "value": "opt_gamma"},
					},
				},
			},
			map[string]any{
				"demo_mm_blocks_select": mmBlocksPostAction("/mm_blocks_integration_static_select", map[string]any{}),
			},
		)
	case "form":
		return p.postMmBlocksMessage(args, "Fill the form; submit or change fields to echo `form_values`.",
			[]any{
				map[string]any{"type": "text", "text": "Form fields call `/mm_blocks_integration_echo_form_values`."},
				map[string]any{
					"type":        "text_input",
					"name":        "title",
					"label":       "Title",
					"placeholder": "Short title",
					"onChange":    "demo_mm_blocks_form",
				},
				map[string]any{
					"type":          "bool_input",
					"name":          "notify_email",
					"label":         "Notify",
					"placeholder":   "Send email",
					"initial_value": false,
					"onChange":      "demo_mm_blocks_form",
				},
				map[string]any{
					"type":        "select",
					"name":        "priority",
					"label":       "Priority",
					"placeholder": "Choose priority",
					"onChange":    "demo_mm_blocks_form",
					"options": []any{
						map[string]any{"text": "Low", "value": "low"},
						map[string]any{"text": "High", "value": "high"},
					},
				},
				map[string]any{
					"type":      "button",
					"text":      "Submit form",
					"style":     "primary",
					"subtype":   "submit",
					"action_id": "demo_mm_blocks_form",
				},
			},
			map[string]any{
				"demo_mm_blocks_form": mmBlocksPostAction("/mm_blocks_integration_echo_form_values", map[string]any{}),
			},
		)
	case "lookup":
		return p.postMmBlocksMessage(args, "Search the dynamic select, then submit.",
			[]any{
				map[string]any{"type": "text", "text": "Options come from `/mm_blocks_integration_lookup`."},
				map[string]any{
					"type":               "select",
					"name":               "pick",
					"label":              "Dynamic option",
					"placeholder":        "Type to search…",
					"data_source":        "dynamic",
					"data_source_action": "demo_mm_blocks_lookup",
				},
				map[string]any{
					"type":      "button",
					"text":      "Submit",
					"style":     "primary",
					"subtype":   "submit",
					"action_id": "demo_mm_blocks_form",
				},
			},
			map[string]any{
				"demo_mm_blocks_lookup": mmBlocksPostAction("/mm_blocks_integration_lookup", map[string]any{}),
				"demo_mm_blocks_form":   mmBlocksPostAction("/mm_blocks_integration_echo_form_values", map[string]any{}),
			},
		)
	case "dialog":
		if scenario == "" {
			scenario = "default"
		}
		buttonText := "Open blocks dialog"
		if scenario != "default" {
			buttonText = "Open " + scenario
		}
		return p.postMmBlocksMessage(args, "Click to open a blocks dialog ("+scenario+").",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_dialog_return` with `context.scenario`."},
				map[string]any{"type": "button", "text": buttonText, "style": "primary", "action_id": "demo_mm_blocks_dialog"},
			},
			map[string]any{
				"demo_mm_blocks_dialog": mmBlocksPostAction("/mm_blocks_dialog_return", map[string]any{"scenario": scenario, "marker": "slash-command"}),
			},
		)
	case "dialog-open":
		return p.postMmBlocksMessage(args, "Click to open a blocks dialog via the open path.",
			[]any{
				map[string]any{"type": "text", "text": "Calls `/mm_blocks_dialog_open`."},
				map[string]any{"type": "button", "text": "Open via dialogs/open", "style": "primary", "action_id": "demo_mm_blocks_dialog_open"},
			},
			map[string]any{
				"demo_mm_blocks_dialog_open": mmBlocksPostAction("/mm_blocks_dialog_open", map[string]any{"marker": "slash-command"}),
			},
		)
	case "example":
		return p.executeMmBlocksExample(args, scenario)
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s. Use `/mm_blocks help` for available commands.", command),
		}
	}
}

func (p *Plugin) postMmBlocksDemo(args *model.CommandArgs) *model.CommandResponse {
	return p.postMmBlocksMessage(args, "Demo mm_blocks interactive message.",
		[]any{
			map[string]any{"type": "text", "text": "Use the controls below to exercise plugin-hosted mm_blocks integrations (no webhook sidecar)."},
			map[string]any{"type": "divider"},
			map[string]any{
				"type": "container",
				"flow": "horizontal",
				"gap":  "medium",
				"content": []any{
					map[string]any{"type": "button", "text": "Ping integration", "style": "primary", "action_id": "demo_mm_blocks_integration"},
					map[string]any{"type": "button", "text": "Update post", "style": "default", "action_id": "demo_mm_blocks_update"},
					map[string]any{"type": "button", "text": "Open dialog", "style": "primary", "action_id": "demo_mm_blocks_dialog"},
					map[string]any{"type": "button", "text": "Open via dialogs/open", "style": "default", "action_id": "demo_mm_blocks_dialog_open"},
				},
			},
			map[string]any{
				"type":        "static_select",
				"placeholder": "Static select",
				"action_id":   "demo_mm_blocks_select",
				"options": []any{
					map[string]any{"text": "Alpha", "value": "opt_alpha"},
					map[string]any{"text": "Beta", "value": "opt_beta"},
				},
			},
		},
		map[string]any{
			"demo_mm_blocks_integration": mmBlocksPostAction("/mm_blocks_integration", map[string]any{}),
			"demo_mm_blocks_update":      mmBlocksPostAction("/mm_blocks_integration_update", map[string]any{}),
			"demo_mm_blocks_dialog":      mmBlocksPostAction("/mm_blocks_dialog_return", map[string]any{"marker": "slash-command"}),
			"demo_mm_blocks_dialog_open": mmBlocksPostAction("/mm_blocks_dialog_open", map[string]any{"marker": "slash-command"}),
			"demo_mm_blocks_select":      mmBlocksPostAction("/mm_blocks_integration_static_select", map[string]any{}),
		},
	)
}

func (p *Plugin) postMmBlocksMessage(args *model.CommandArgs, message string, blocks []any, actions map[string]any) *model.CommandResponse {
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   message,
		Props: model.StringInterface{
			"mm_blocks":         blocks,
			"mm_blocks_actions": actions,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		const errorMessage = "Failed to create mm_blocks post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}
