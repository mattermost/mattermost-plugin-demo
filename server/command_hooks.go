package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	commandTriggerHooks  = "demo_plugin"
	commandTriggerDialog = "dialog"

	dialogElementNameNumber = "somenumber"
)

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerHooks,
		AutoComplete:     true,
		AutoCompleteHint: "(true|false)",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Demo Plugin Hooks Command",
		Description:      "A command used to enable or disable the demo plugin hooks.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerHooks)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerDialog,
		AutoComplete:     true,
		AutoCompleteDesc: "Open an Interactive Dialog.",
		DisplayName:      "Demo Plugin Command",
		Description:      "A command to open an Interactive Dialog.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerDialog)
	}

	return nil
}

func (p *Plugin) emitStatusChange() {
	configuration := p.getConfiguration()

	p.API.PublishWebSocketEvent("status_change", map[string]interface{}{
		"enabled": !configuration.disabled,
	}, &model.WebsocketBroadcast{})
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
//
// This demo implementation responds to a /demo_plugin command, allowing the user to enable
// or disable the demo plugin's hooks functionality (but leave the command and webapp enabled).
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case commandTriggerHooks:
		return p.executeCommandHooks(c, args)
	case commandTriggerDialog:
		return p.executeCommandDialog(c, args)
	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}
}

func (p *Plugin) executeCommandHooks(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	configuration := p.getConfiguration()

	if strings.HasSuffix(args.Command, "true") {
		if !configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The demo plugin hooks are already enabled.",
			}, nil
		}

		configuration.disabled = false
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Enabled demo plugin hooks.",
		}, nil

	}

	if strings.HasSuffix(args.Command, "false") {
		if configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The demo plugin hooks are already disabled.",
			}, nil
		}

		configuration.disabled = true
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Disabled demo plugin hooks.",
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
	}, nil
}

func (p *Plugin) executeCommandDialog(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	serverConfig := p.API.GetConfig()

	fields := strings.Fields(args.Command)
	var dialogRequest model.OpenDialogRequest

	if len(fields) == 2 && fields[1] == "no-elements" {
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/2", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog: model.Dialog{
				CallbackId:     "somecallbackid",
				Title:          "Sample Confirmation Dialog",
				IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
				Elements:       nil,
				SubmitLabel:    "Confirm",
				NotifyOnCancel: true,
				State:          "somestate",
			},
		}
	} else {
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog: model.Dialog{
				CallbackId: "somecallbackid",
				Title:      "Test Title",
				IconURL:    "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
				Elements: []model.DialogElement{{
					DisplayName: "Display Name",
					Name:        "realname",
					Type:        "text",
					Default:     "default text",
					Placeholder: "placeholder",
					HelpText:    "This a regular input in an interactive dialog triggered by a test integration.",
				}, {
					DisplayName: "Email",
					Name:        "someemail",
					Type:        "text",
					SubType:     "email",
					Placeholder: "placeholder@bladekick.com",
					HelpText:    "This a regular email input in an interactive dialog triggered by a test integration.",
				}, {
					DisplayName: "Number",
					Name:        dialogElementNameNumber,
					Type:        "text",
					SubType:     "number",
				}, {
					DisplayName: "Display Name Long Text Area",
					Name:        "realnametextarea",
					Type:        "textarea",
					Placeholder: "placeholder",
					Optional:    true,
					MinLength:   5,
					MaxLength:   100,
				}, {
					DisplayName: "User Selector",
					Name:        "someuserselector",
					Type:        "select",
					Placeholder: "Select a user...",
					HelpText:    "Choose a user from the list.",
					Optional:    true,
					MinLength:   5,
					MaxLength:   100,
					DataSource:  "users",
				}, {
					DisplayName: "Channel Selector",
					Name:        "somechannelselector",
					Type:        "select",
					Placeholder: "Select a channel...",
					HelpText:    "Choose a channel from the list.",
					Optional:    true,
					MinLength:   5,
					MaxLength:   100,
					DataSource:  "channels",
				}, {
					DisplayName: "Option Selector",
					Name:        "someoptionselector",
					Type:        "select",
					Placeholder: "Select an option...",
					HelpText:    "Choose an option from the list.",
					Options: []*model.PostActionOptions{{
						Text:  "Option1",
						Value: "opt1",
					}, {
						Text:  "Option2",
						Value: "opt2",
					}, {
						Text:  "Option3",
						Value: "opt3",
					}},
				}},
				SubmitLabel:    "Submit",
				NotifyOnCancel: true,
				State:          "somestate",
			},
		}
	}

	if err := p.API.OpenInteractiveDialog(dialogRequest); err != nil {
		errorMessage := "Failed to open Interactive Dialog"
		p.API.LogError(errorMessage, "err", err)
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}, nil
	}
	return &model.CommandResponse{}, nil
}
