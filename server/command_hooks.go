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
	commandTriggerToast             = "toast"

	dialogElementNameNumber   = "somenumber"
	dialogElementNameEmail    = "someemail"
	dialogElementNameDate     = "somedate"
	dialogElementNameDatetime = "somedatetime"

	dialogStateSome                = "somestate"
	dialogStateRelativeCallbackURL = "relativecallbackstate"
	dialogIntroductionText         = "**Some** _introductory_ paragraph in Markdown formatted text with [link](https://mattermost.com)"

	commandDialogHelp = "###### Interactive Dialog Slash Command Help\n" +
		"- `/dialog` - Open an Interactive Dialog. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog basic` - Open a simple Interactive Dialog with one optional text field for basic testing.\n" +
		"- `/dialog boolean` - Open an Interactive Dialog with boolean fields for testing toggle functionality.\n" +
		"- `/dialog textfields` - Open an Interactive Dialog with various text field types for testing input validation.\n" +
		"- `/dialog selectfields` - Open an Interactive Dialog with select, radio, user, and channel selectors.\n" +
		"- `/dialog no-elements` - Open an Interactive Dialog with no elements. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog relative-callback-url` - Open an Interactive Dialog with relative callback URL. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog introduction-text` - Open an Interactive Dialog with optional introduction text. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog dynamic-select` - Open an Interactive Dialog with dynamic select fields. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog date` - Open an Interactive Dialog with date and datetime fields for testing.\n" +
		"- `/dialog multi-select` - Open an Interactive Dialog with multi-select fields. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog error` - Open an Interactive Dialog which always returns an general error.\n" +
		"- `/dialog error-no-elements` - Open an Interactive Dialog with no elements which always returns an general error.\n" +
		"- `/dialog field-refresh` - Open an Interactive Dialog with field refresh functionality.\n" +
		"- `/dialog multistep` - Open a multi-step Interactive Dialog demonstrating form refresh on submit.\n" +
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

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerToast,
		AutoComplete:     true,
		AutoCompleteDesc: "Demonstrates the toast notification API.",
		AutocompleteData: getCommandToastAutocompleteData(),
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerToast)
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

	date := model.NewAutocompleteData("date", "", "Open an Interactive Dialog with date and datetime fields.")
	command.AddCommand(date)

	error := model.NewAutocompleteData("error", "", "Open an Interactive Dialog with error.")
	command.AddCommand(error)

	errorNoElements := model.NewAutocompleteData("error-no-elements", "", "Open an Interactive Dialog with error no elements.")
	command.AddCommand(errorNoElements)

	dynamicSelect := model.NewAutocompleteData("dynamic-select", "", "Open an Interactive Dialog with dynamic select fields.")
	command.AddCommand(dynamicSelect)

	fieldRefresh := model.NewAutocompleteData("field-refresh", "", "Open an Interactive Dialog with field refresh functionality.")
	command.AddCommand(fieldRefresh)

	multistep := model.NewAutocompleteData("multistep", "", "Open a multi-step Interactive Dialog with form refresh on submit.")
	command.AddCommand(multistep)

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

func getCommandToastAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData(commandTriggerToast, "[--all-sessions] [position] [message]", "Send a toast notification.")

	// Add --all-sessions flag
	allSessions := model.NewAutocompleteData("--all-sessions", "[position] [message]", "Send toast to all sessions")

	// Add position options to both main command and --all-sessions
	for _, parent := range []*model.AutocompleteData{command, allSessions} {
		topLeft := model.NewAutocompleteData("top-left", "[message]", "Show toast at top-left")
		topLeft.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(topLeft)

		topCenter := model.NewAutocompleteData("top-center", "[message]", "Show toast at top-center")
		topCenter.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(topCenter)

		topRight := model.NewAutocompleteData("top-right", "[message]", "Show toast at top-right")
		topRight.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(topRight)

		bottomLeft := model.NewAutocompleteData("bottom-left", "[message]", "Show toast at bottom-left")
		bottomLeft.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(bottomLeft)

		bottomCenter := model.NewAutocompleteData("bottom-center", "[message]", "Show toast at bottom-center")
		bottomCenter.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(bottomCenter)

		bottomRight := model.NewAutocompleteData("bottom-right", "[message]", "Show toast at bottom-right (default)")
		bottomRight.AddTextArgument("Message to display", "[message]", "")
		parent.AddCommand(bottomRight)
	}

	command.AddCommand(allSessions)

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
	case commandTriggerToast:
		return p.executeCommandToast(c, args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown command: %s. Use `/dialog help` for available commands.", args.Command),
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
	case "basic":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/3", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogBasic(),
		}
	case "boolean":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/3", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogBoolean(),
		}
	case "textfields":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/3", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogTextFields(),
		}
	case "selectfields":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/3", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogSelectFields(),
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
	case "dynamic-select":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/1", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithDynamicSelectElements(),
		}
	case "date":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/date", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithDateElements(),
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
			Dialog:    getDialogBasic(),
		}
	case "error-no-elements":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("/plugins/%s/dialog/error", manifest.Id),
			Dialog:    getDialogWithoutElements(dialogStateSome),
		}
	case "field-refresh":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/field-refresh", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogWithFieldRefresh(""), // Start with no project type selected
		}
	case "multistep":
		dialogRequest = model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:       fmt.Sprintf("%s/plugins/%s/dialog/multistep", *serverConfig.ServiceSettings.SiteURL, manifest.Id),
			Dialog:    getDialogStep1(),
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

func (p *Plugin) executeCommandToast(c *plugin.Context, args *model.CommandArgs) *model.CommandResponse {
	fields := strings.Fields(args.Command)

	// Default values
	position := "bottom-right"
	message := "This is a demo toast notification!"
	connectionID := ""
	allSessions := false

	// Check if --all-sessions flag is present in the first position
	startIndex := 1
	if len(fields) >= 2 && fields[1] == "--all-sessions" {
		allSessions = true
		startIndex = 2
	}

	// If --all-sessions is NOT set, use the session ID
	if !allSessions {
		var found bool
		connectionID, found = p.GetConnectionIDForSession(c.SessionId)
		if !found {
			p.API.LogWarn("Failed to get connection ID for session", "session_id", c.SessionId)
		}
	}

	// Parse command arguments: /toast [--all-sessions] [position] [message]
	if len(fields) >= startIndex+1 {
		position = fields[startIndex]
	}
	if len(fields) >= startIndex+2 {
		// Join all remaining fields as the message
		message = strings.Join(fields[startIndex+1:], " ")
	}

	// Send the toast message using the plugin API
	options := model.SendToastMessageOptions{
		Position: position,
	}

	if err := p.client.Frontend.SendToastMessage(args.UserId, connectionID, message, options); err != nil {
		errorMessage := "Failed to send toast notification"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         errorMessage,
		}
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("Toast notification sent to position: %s", position),
	}
}
