package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

const (
	commandTriggerCrash             = "crash"
	commandTriggerHooks             = "demo_plugin"
	commandTriggerDialog            = "dialog"
	commandTriggerEphemeral         = "ephemeral"
	commandTriggerEphemeralOverride = "ephemeral_override"
	commandTriggerInteractive       = "interactive"
	commandTriggerMentions          = "show_mentions"
	commandTriggerListFiles         = "list_files"
	commandTriggerAutocompleteTest  = "autocomplete_test"

	dialogElementNameNumber = "somenumber"
	dialogElementNameEmail  = "someemail"

	dialogStateSome                = "somestate"
	dialogStateRelativeCallbackURL = "relativecallbackstate"
	dialogIntroductionText         = "**Some** _introductory_ paragraph in Markdown formatted text with [link](https://mattermost.com)"

	commandDialogHelp = "###### Interactive Dialog Slash Command Help\n" +
		"- `/dialog` - Open an Interactive Dialog. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog no-elements` - Open an Interactive Dialog with no elements. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog relative-callback-url` - Open an Interactive Dialog with relative callback URL. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog introduction-text` - Open an Interactive Dialog with optional introduction text. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog multi-select` - Open an Interactive Dialog with multi-select fields. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog error` - Open an Interactive Dialog which always returns an general error.\n" +
		"- `/dialog error-no-elements` - Open an Interactive Dialog with no elements which always returns an general error.\n" +
		"- `/dialog help` - Show this help text"
)

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{

		Trigger:          commandTriggerHooks,
		AutoComplete:     true,
		AutoCompleteHint: "(true|false)",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		AutocompleteData: getCommandHooksAutocompleteData(),
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerHooks)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerCrash,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Crashes Demo Plugin",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerCrash)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerEphemeral,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates an ephemeral post capabilities.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerEphemeral)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerEphemeralOverride,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates an ephemeral post overridden in the webapp.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerEphemeralOverride)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerDialog,
		AutoComplete:     true,
		AutoCompleteDesc: "Open an Interactive Dialog.",
		DisplayName:      "Demo Plugin Command",
		AutocompleteData: getCommandDialogAutocompleteData(),
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerDialog)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerInteractive,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates  interactive message buttons.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerInteractive)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerMentions,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates access to mentions in the message.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerMentions)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerListFiles,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates the file search plugin api.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerInteractive)
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerAutocompleteTest,
		AutoComplete:     true,
		AutocompleteData: getAutocompleteTestAutocompleteData(),
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

func getCommandHooksAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData(commandTriggerHooks, "", "Enables or disables the demo plugin hooks.")
	command.AddStaticListArgument("", true, []model.AutocompleteListItem{
		{
			Item:     "true",
			HelpText: "Enable demo plugin hooks",
		}, {
			Item:     "false",
			HelpText: "Disable demo plugin hooks",
		},
	})
	return command
}

func getCommandDialogAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData(commandTriggerDialog, "", "Open an Interactive Dialog.")

	noElements := model.NewAutocompleteData("no-elements", "", "Open an Interactive Dialog with no elements.")
	command.AddCommand(noElements)

	relativeCallbackURL := model.NewAutocompleteData("relative-callback-url", "", "Open an Interactive Dialog with a relative callback url.")
	command.AddCommand(relativeCallbackURL)

	introText := model.NewAutocompleteData("introduction-text", "", "Open an Interactive Dialog with an introduction text.")
	command.AddCommand(introText)

	error := model.NewAutocompleteData("error", "", "Open an Interactive Dialog with error.")
	command.AddCommand(error)

	errorNoElements := model.NewAutocompleteData("error-no-elements", "", "Open an Interactive Dialog with error no elements.")
	command.AddCommand(errorNoElements)

	multiSelect := model.NewAutocompleteData("multi-select", "", "Open an Interactive Dialog with multi-select fields.")
	command.AddCommand(multiSelect)

	help := model.NewAutocompleteData("help", "", "")
	command.AddCommand(help)

	return command
}

func getAutocompleteTestAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData(commandTriggerAutocompleteTest, "", "Test an autocomplete.")

	dynamicArg := model.NewAutocompleteData("dynamic-arg", "", "Test a dynamic argument")
	dynamicArg.AddDynamicListArgument("Some dynamic argument", "dynamic_arg_test_url", true)
	command.AddCommand(dynamicArg)

	namedArg := model.NewAutocompleteData("named-arg", "", "Test a named argument")
	namedArg.AddNamedTextArgument("name", "Input named argument with pattern p([a-z]+)ch", "", "p([a-z]+)ch", true)
	command.AddCommand(namedArg)

	optionalArg := model.NewAutocompleteData("optional-arg", "", "Test an optional argument")
	optionalArg.AddNamedTextArgument("name1", "Optional named argument", "", "", false)
	optionalArg.AddNamedTextArgument("name2", "Optional named argument with pattern p([a-z]+)ch", "", "p([a-z]+)ch", false)
	command.AddCommand(optionalArg)

	return command
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
//
// This demo implementation responds to a /demo_plugin command, allowing the user to enable
// or disable the demo plugin's hooks functionality (but leave the command and webapp enabled).
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	delay := p.getConfiguration().IntegrationRequestDelay
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case commandTriggerCrash:
		return p.executeCommandCrash(), nil
	case commandTriggerHooks:
		return p.executeCommandHooks(args), nil
	case commandTriggerEphemeral:
		return p.executeCommandEphemeral(args), nil
	case commandTriggerEphemeralOverride:
		return p.executeCommandEphemeralOverride(args), nil
	case commandTriggerDialog:
		return p.executeCommandDialog(args), nil
	case commandTriggerListFiles:
		return p.executeCommandListFiles(args), nil
	case commandTriggerInteractive:
		return p.executeCommandInteractive(args), nil
	case commandTriggerMentions:
		return p.executeCommandMentions(args), nil
	case commandTriggerAutocompleteTest:
		return p.executeAutocompleteTest(args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s", args.Command),
		}, nil
	}
}

func (p *Plugin) executeCommandCrash() *model.CommandResponse {
	go p.crash()
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Crashing plugin",
	}
}

func (p *Plugin) executeCommandHooks(args *model.CommandArgs) *model.CommandResponse {
	configuration := p.getConfiguration()

	if strings.HasSuffix(args.Command, "true") {
		if !configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         "The demo plugin hooks are already enabled.",
			}
		}

		p.setEnabled(true)
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Enabled demo plugin hooks.",
		}
	}

	if strings.HasSuffix(args.Command, "false") {
		if configuration.disabled {
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         "The demo plugin hooks are already disabled.",
			}
		}

		p.setEnabled(false)
		p.emitStatusChange()

		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Disabled demo plugin hooks.",
		}
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Unknown command action: %s", args.Command),
	}
}

func (p *Plugin) executeCommandEphemeral(args *model.CommandArgs) *model.CommandResponse {
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL

	post := &model.Post{
		ChannelId: args.ChannelId,
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
					Type: model.PostActionTypeButton,
					Name: "Update",
				}, {
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/delete", siteURL, manifest.Id),
					},
					Type: model.PostActionTypeButton,
					Name: "Delete",
				}},
			}},
		},
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandEphemeralOverride(args *model.CommandArgs) *model.CommandResponse {
	_ = p.API.SendEphemeralPost(args.UserId, &model.Post{
		ChannelId: args.ChannelId,
		Message:   "This is a demo of overriding an ephemeral post.",
		Props: model.StringInterface{
			"type": "custom_demo_plugin_ephemeral",
		},
	})
	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandDialog(args *model.CommandArgs) *model.CommandResponse {
	serverConfig := p.API.GetConfig()

	var dialogRequest model.OpenDialogRequest
	fields := strings.Fields(args.Command)
	command := ""
	if len(fields) == 2 {
		command = fields[1]
	}

	switch command {
	case "help":
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         commandDialogHelp,
		}
	case "":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithSampleElements(),
		}
	case "no-elements":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/2", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateSome),
		}
	case "relative-callback-url":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog/2", manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateRelativeCallbackURL),
		}
	case "introduction-text":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithIntroductionText(dialogIntroductionText),
		}
	case "multi-select":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithMultiSelectElements(),
		}
	case "error":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog/error", manifest.Id),
			Dialog:    getDialogWithSampleElements(),
		}
	case "error-no-elements":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog/error", manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateSome),
		}
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s", command),
		}
	}

	if err := p.API.OpenInteractiveDialog(dialogRequest); err != nil {
		errorMessage := "Failed to open Interactive Dialog"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
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
			DisplayName: "Password",
			Name:        "somepassword",
			Type:        "text",
			SubType:     "password",
			Placeholder: "Password",
			HelpText:    "This a password input in an interactive dialog triggered by a test integration.",
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
		}, {
			DisplayName: "Option Selector with default",
			Name:        "someoptionselector2",
			Type:        "select",
			Default:     "opt2",
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
		}, {
			DisplayName: "Boolean Selector",
			Name:        "someboolean",
			Type:        "bool",
			Placeholder: "Agree to the terms of service",
			HelpText:    "You must agree to the terms of service to proceed.",
		}, {
			DisplayName: "Boolean Selector",
			Name:        "someboolean_optional",
			Type:        "bool",
			Placeholder: "Sign up for monthly emails?",
			HelpText:    "It's up to you if you want to get monthly emails.",
			Optional:    true,
		}, {
			DisplayName: "Boolean Selector (default true)",
			Name:        "someboolean_default_true",
			Type:        "bool",
			Placeholder: "Enable secure login",
			HelpText:    "You must enable secure login to proceed.",
			Default:     "true",
		}, {
			DisplayName: "Boolean Selector (default true)",
			Name:        "someboolean_default_true_optional",
			Type:        "bool",
			Placeholder: "Enable painfully secure login",
			HelpText:    "You may optionally enable painfully secure login.",
			Default:     "true",
			Optional:    true,
		}, {
			DisplayName: "Boolean Selector (default false)",
			Name:        "someboolean_default_false",
			Type:        "bool",
			Placeholder: "Agree to the annoying terms of service",
			HelpText:    "You must also agree to the annoying terms of service to proceed.",
			Default:     "false",
		}, {
			DisplayName: "Boolean Selector (default false)",
			Name:        "someboolean_default_false_optional",
			Type:        "bool",
			Placeholder: "Throw-away account",
			HelpText:    "A throw-away account will be deleted after 24 hours.",
			Default:     "false",
			Optional:    true,
		}, {
			DisplayName: "Radio Option Selector",
			Name:        "someradiooptionselector",
			Type:        "radio",
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

func getDialogWithIntroductionText(introductionText string) model.Dialog {
	dialog := getDialogWithSampleElements()
	dialog.IntroductionText = introductionText
	return dialog
}

func getDialogWithMultiSelectElements() model.Dialog {
	return model.Dialog{
		CallbackId: "somecallbackid",
		Title:      "Multi-Select Dialog Demo",
		IconURL:    "http://www.mattermost.org/wp-content/uploads/2016/04/icon.png",
		Elements: []model.DialogElement{{
			DisplayName: "Multi-Select Users",
			Name:        "multiselect_users",
			Type:        "select",
			Placeholder: "Select multiple users...",
			HelpText:    "Choose multiple users from the list.",
			DataSource:  "users",
			MultiSelect: true,
		}, {
			DisplayName: "Multi-Select Channels",
			Name:        "multiselect_channels",
			Type:        "select",
			Placeholder: "Select multiple channels...",
			HelpText:    "Choose multiple channels from the list.",
			DataSource:  "channels",
			MultiSelect: true,
		}, {
			DisplayName: "Multi-Select Options",
			Name:        "multiselect_options",
			Type:        "select",
			Placeholder: "Select multiple options...",
			HelpText:    "Choose multiple options from the list.",
			MultiSelect: true,
			Options: []*model.PostActionOptions{{
				Text:  "Option A",
				Value: "optA",
			}, {
				Text:  "Option B",
				Value: "optB",
			}, {
				Text:  "Option C",
				Value: "optC",
			}, {
				Text:  "Option D",
				Value: "optD",
			}},
		}},
		SubmitLabel:    "Submit Multi-Select",
		NotifyOnCancel: true,
		State:          dialogStateSome,
	}
}

func (p *Plugin) executeCommandInteractive(args *model.CommandArgs) *model.CommandResponse {
	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   "Test interactive button",
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{{
				Actions: []*model.PostAction{{
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("/plugins/%s/interactive/button/1", manifest.Id),
					},
					Type: model.PostActionTypeButton,
					Name: "Interactive Button",
				}},
			}},
		},
	}

	_, err := p.API.CreatePost(post)
	if err != nil {
		const errorMessage = "Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}

func (p *Plugin) crash() {
	<-time.NewTimer(time.Second).C
	y := 0
	_ = 1 / y
}

func (p *Plugin) executeCommandMentions(args *model.CommandArgs) *model.CommandResponse {
	message := "The command `" + args.Command + "` contains the following different mentions.\n"
	message += "### Mentions to users in the team\n"
	if args.UserMentions == nil {
		message += "_There are no mentions to users in the team in your command_.\n"
	} else {
		message += "| User name | ID |\n"
		message += "|-----------|----|\n"
		for name, id := range args.UserMentions {
			message += fmt.Sprintf("|@%s|%s|\n", name, id)
		}
	}

	message += "\n### Mentions to public channels\n"
	if args.ChannelMentions == nil {
		message += "_There are no mentions to public channels in your command_.\n"
	} else {
		message += "| Channel name | ID |\n"
		message += "|--------------|----|\n"
		for name, id := range args.ChannelMentions {
			message += fmt.Sprintf("|~%s|%s|\n", name, id)
		}
	}

	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botID,
		Message:   message,
	}

	_, err := p.API.CreatePost(post)
	if err != nil {
		const errorMessage = "Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{}
}

func (p *Plugin) executeCommandListFiles(args *model.CommandArgs) *model.CommandResponse {
	fileInfos, err := p.API.GetFileInfos(0, 10, &model.GetFileInfosOptions{
		ChannelIds:     []string{args.ChannelId},
		SortDescending: true,
	})
	if err != nil {
		errorMessage := "Failed to get file list"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	team, err := p.API.GetTeam(args.TeamId)
	if err != nil {
		errorMessage := "Failed to get team name"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	permaLink := args.SiteURL + "/" + team.Name + "/pl/"
	attachments := make([]*model.SlackAttachment, 0, len(fileInfos))
	for _, f := range fileInfos {
		user, err := p.API.GetUser(f.CreatorId)
		if err != nil {
			errorMessage := "Failed to get username"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         errorMessage,
			}
		}
		fileLink, err := p.API.GetFileLink(f.Id)
		if err != nil {
			errorMessage := "Failed to get file public link"
			p.API.LogError(errorMessage, "err", err.Error())
			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         errorMessage,
			}
		}
		attachments = append(attachments,
			&model.SlackAttachment{
				Title:     f.Name,
				TitleLink: permaLink + f.PostId,
				Text:      fmt.Sprintf("uploaded by %s", user.Username),
				Fields: []*model.SlackAttachmentField{
					{
						Title: "Direct Download Link",
						Value: args.SiteURL + fileLink,
					}},
			},
		)
	}

	post := &model.Post{
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf("Last %d Files uploaded to this channel", len(fileInfos)),
		Props: model.StringInterface{
			"attachments": attachments,
		},
	}

	_ = p.API.SendEphemeralPost(args.UserId, post)
	return &model.CommandResponse{}
}

func (p *Plugin) executeAutocompleteTest(args *model.CommandArgs) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Executed command: %s", args.Command),
	}
}
