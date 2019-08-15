package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	commandTriggerCrash             = "crash"
	commandTriggerHooks             = "demo_plugin"
	commandTriggerDialog            = "dialog"
	commandTriggerEphemeral         = "ephemeral"
	commandTriggerEphemeralOverride = "ephemeral_override"
	commandTestCallback             = "test_callback"

	dialogElementNameNumber = "somenumber"
	dialogElementNameEmail  = "someemail"

	dialogStateSome                = "somestate"
	dialogStateRelativeCallbackURL = "relativecallbackstate"

	commandDialogHelp = "###### Interactive Dialog Slash Command Help\n" +
		"- `/dialog` - pen an Interactive Dialog. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog no-elements` - Open an Interactive Dialog with no elements. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog relative-callback-url` - Open an Interactive Dialog with relative callback URL. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog help` - Show this help text"
)

func (p *Plugin) registerCommands() error {
	if err := p.Helpers.RegisterCommand(&model.Command{
		Trigger:          commandTriggerHooks,
		AutoComplete:     true,
		AutoCompleteHint: "(true|false)",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
	}, func(c *plugin.Context, args *plugin.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandHooks(args), nil
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerHooks)
	}

	if err := p.Helpers.RegisterCommand(&model.Command{
		Trigger:          commandTriggerCrash,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Crashes Demo Plugin",
	}, func(c *plugin.Context, args *plugin.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandCrash(), nil
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerCrash)
	}

	if err := p.Helpers.RegisterCommand(&model.Command{
		Trigger:          commandTriggerEphemeral,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates an ephemeral post capabilities.",
	}, func(c *plugin.Context, args *plugin.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandEphemeral(args), nil
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerEphemeral)
	}

	if err := p.Helpers.RegisterCommand(&model.Command{
		Trigger:          commandTriggerEphemeralOverride,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates an ephemeral post overriden in the webapp.",
	}, func(c *plugin.Context, args *plugin.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandEphemeralOverride(args), nil
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerEphemeralOverride)
	}

	if err := p.Helpers.RegisterCommand(&model.Command{
		Trigger:          commandTriggerDialog,
		AutoComplete:     true,
		AutoCompleteDesc: "Open an Interactive Dialog.",
		DisplayName:      "Demo Plugin Command",
	}, func(c *plugin.Context, args *plugin.CommandArgs) (*model.CommandResponse, *model.AppError) {
		return p.executeCommandDialog(args), nil
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
	return p.Helpers.ExecuteCommand(c, args)
}

func (p *Plugin) executeCommandCrash() *model.CommandResponse {
	go p.crash()
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Crashing plugin",
	}
}

func (p *Plugin) executeCommandHooks(args *plugin.CommandArgs) *model.CommandResponse {
	configuration := p.getConfiguration()

	if strings.HasSuffix(args.Args[0], "true") {
		if !configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The demo plugin hooks are already enabled.",
			}
		}

		configuration.disabled = false
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Enabled demo plugin hooks.",
		}

	}

	if strings.HasSuffix(args.Args[0], "false") {
		if configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         "The demo plugin hooks are already disabled.",
			}
		}

		configuration.disabled = true
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Disabled demo plugin hooks.",
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Trigger),
	}
}

func (p *Plugin) executeCommandEphemeral(args *plugin.CommandArgs) *model.CommandResponse {
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL

	post := &model.Post{
		ChannelId: args.OriginalArgs.ChannelId,
		Message:   "test ephemeral actions",
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{{
				Actions: []*model.PostAction{{
					Integration: &model.PostActionIntegration{
						Context: model.StringInterface{
							"count": 0,
						},
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/update", siteURL, manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: "Update",
				}, {
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/delete", siteURL, manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: "Delete",
				}},
			}},
		},
	}

	_ = p.API.SendEphemeralPost(args.OriginalArgs.UserId, post)
	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandEphemeralOverride(args *plugin.CommandArgs) *model.CommandResponse {
	_ = p.API.SendEphemeralPost(args.OriginalArgs.UserId, &model.Post{
		ChannelId: args.OriginalArgs.ChannelId,
		Message:   "This is a demo of overriding an ephemeral post.",
		Props: model.StringInterface{
			"type": "custom_demo_plugin_ephemeral",
		},
	})
	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandDialog(args *plugin.CommandArgs) *model.CommandResponse {
	serverConfig := p.API.GetConfig()

	var dialogRequest model.OpenDialogRequest
	command := ""
	if len(args.Args) == 2 {
		command = args.Args[1]
	}

	switch command {
	case "help":
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         commandDialogHelp,
		}
	case "":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.OriginalArgs.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithSampleElements(),
		}
	case "no-elements":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.OriginalArgs.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/2", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateSome),
		}
	case "relative-callback-url":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.OriginalArgs.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog/2", manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateRelativeCallbackURL),
		}
	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + command),
		}
	}

	if err := p.API.OpenInteractiveDialog(dialogRequest); err != nil {
		errorMessage := "Failed to open Interactive Dialog"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}
	return &model.CommandResponse{}
}

func getDialogWithSampleElements() model.Dialog {
	return model.Dialog{
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
			Name:        dialogElementNameEmail,
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
		State:          dialogStateSome,
	}
}

func getDialogWithoutElements(state string) model.Dialog {
	return model.Dialog{
		CallbackId:     "somecallbackid",
		Title:          "Sample Confirmation Dialog",
		IconURL:        "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		Elements:       nil,
		SubmitLabel:    "Confirm",
		NotifyOnCancel: true,
		State:          state,
	}
}

func (p *Plugin) crash() {
	<-time.NewTimer(time.Second).C
	y := 0
	_ = 1 / y
}
