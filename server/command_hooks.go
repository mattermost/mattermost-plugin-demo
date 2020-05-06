package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	svg "github.com/h2non/go-is-svg"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	commandTriggerCrash             = "crash"
	commandTriggerHooks             = "demo_plugin"
	commandTriggerDialog            = "dialog"
	commandTriggerEphemeral         = "ephemeral"
	commandTriggerEphemeralOverride = "ephemeral_override"
	commandTriggerInteractive       = "interactive"

	dialogElementNameNumber = "somenumber"
	dialogElementNameEmail  = "someemail"

	dialogStateSome                = "somestate"
	dialogStateRelativeCallbackURL = "relativecallbackstate"
	dialogIntroductionText         = "**Some** _introductory_ paragraph in Markdown formatted text with [link](https://mattermost.com)"

	commandDialogHelp = "###### Interactive Dialog Slash Command Help\n" +
		"- `/dialog` - pen an Interactive Dialog. Once submitted, user-entered input is posted back into a channel.\n" +
		"- `/dialog no-elements` - Open an Interactive Dialog with no elements. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog relative-callback-url` - Open an Interactive Dialog with relative callback URL. Once submitted, user's action is posted back into a channel.\n" +
		"- `/dialog introduction-text` - Open an Interactive Dialog with optional introduction text. Once submitted, user's action is posted back into a channel.\n" +
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
		AutoCompleteDesc: "Demonstrates an ephemeral post overriden in the webapp.",
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerEphemeralOverride)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerDialog,
		AutoComplete:     true,
		AutoCompleteDesc: "Open an Interactive Dialog.",
		DisplayName:      "Demo Plugin Command",
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

	jiraIconData, err := p.readFile("/assets/icon.svg")
	if err != nil {
		return errors.Wrap(err, "failed to read icon image")
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:              "jira",
		AutoComplete:         true,
		AutoCompleteHint:     "",
		AutoCompleteDesc:     "Demonstrates  interactive message buttons.",
		AutocompleteData:     createJiraAutocompleteData(),
		AutocompleteIconData: jiraIconData,
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerInteractive)
	}

	agendaIconData, err := p.readFile("/assets/github.svg")
	if err != nil {
		return errors.Wrap(err, "failed to read icon image")
	}
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:              "agenda",
		AutoComplete:         true,
		AutoCompleteHint:     "",
		AutoCompleteDesc:     "Demonstrates  interactive message buttons.",
		AutocompleteData:     createAgendaAutocompleteData(),
		AutocompleteIconData: agendaIconData,
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", commandTriggerInteractive)
	}

	return nil
}

func (p *Plugin) readFile(path string) (string, error) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return "", errors.Wrap(err, "failed to get bundle path")
	}

	imageBytes, err := ioutil.ReadFile(filepath.Join(bundlePath, path))
	if err != nil {
		return "", errors.Wrap(err, "failed to read image")
	}
	if !svg.Is(imageBytes) {
		return "", errors.Wrapf(err, "icon is not svg %s", "/assets/icon.svg")
	}
	iconData := fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(imageBytes))

	return iconData, nil
}

func getIcon(iconPath string) (string, error) {
	icon, err := ioutil.ReadFile(iconPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open icon at path %s", iconPath)
	}
	if !svg.Is(icon) {
		return "", errors.Wrapf(err, "icon is not svg %s", iconPath)
	}
	return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(icon)), nil
}

func createAgendaAutocompleteData() *model.AutocompleteData {
	agenda := model.NewAutocompleteData("agenda", "[command]", "Available commands: list, queue, setting, help")

	list := model.NewAutocompleteData("list", "", "Show a list of items queued for the next meeting")
	optionalListNextWeek := model.NewAutocompleteData("next-week", "(optional)", "If `next-week` is provided, it will list the agenda for the next calendar week.")
	list.AddCommand(optionalListNextWeek)
	agenda.AddCommand(list)

	queue := model.NewAutocompleteData("queue", "", "Queue `message` as a topic on the next meeting.")
	queue.AddStaticListArgument("If `next-week` is provided, it will queue for the meeting in the next calendar week.", []model.AutocompleteListItem{{
		HelpText: "If `next-week` is provided, it will queue for the meeting in the next calendar week.",
		Hint:     "(optional)",
		Item:     "next-week",
	}}, false)
	queue.AddTextArgument("Creates a post for user with the given message for the next meeting date.", "message", "")
	agenda.AddCommand(queue)

	setting := model.NewAutocompleteData("setting", "", "Update the setting.")
	schedule := model.NewAutocompleteData("schedule", "", "Update schedule.")
	schedule.AddTextArgument("Must be between 1-5", "weekday", "")
	setting.AddCommand(schedule)
	hashtag := model.NewAutocompleteData("hashtag", "", "Update hastag.")
	hashtag.AddTextArgument("input hashtag", "Default: Jan02", "")
	setting.AddCommand(hashtag)
	agenda.AddCommand(setting)

	help := model.NewAutocompleteData("help", "", "Mattermost Agenda plugin slash command help")
	agenda.AddCommand(help)
	return agenda
}

func createJiraAutocompleteData() *model.AutocompleteData {
	jira := model.NewAutocompleteData("jira", "[command]", "Available commands: connect, assign, disconnect, create, transition, view, subscribe, settings, install cloud/server, uninstall cloud/server, help")

	connect := model.NewAutocompleteData("connect", "[url]", "Connect your Mattermost account to your Jira account")
	connect.AddTextArgument("connect to the server", "[url]", "")
	jira.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your Mattermost account from your Jira account")
	jira.AddCommand(disconnect)

	assign := model.NewAutocompleteData("assign", "[issue] [user]", "Change the assignee of a Jira issue")
	assign.AddDynamicListArgument("List of issues is downloading from your Jira account", "dynamic_issues", true)
	assign.AddDynamicListArgument("List of assignees is downloading from your Jira account", "dynamic_users", true)
	jira.AddCommand(assign)

	create := model.NewAutocompleteData("create", "[issue text]", "Create a new Issue")
	create.AddTextArgument("This text is optional, will be inserted into the description field", "[text]", "")
	jira.AddCommand(create)

	transition := model.NewAutocompleteData("transition", "[issue]", "Change the state of a Jira issue")
	transition.AddDynamicListArgument("List of issues is downloading from your Jira account", "dynamic_issues", true)
	transition.AddDynamicListArgument("List of states is downloading from your Jira account", "dynamic_states", true)
	jira.AddCommand(transition)

	subscribe := model.NewAutocompleteData("subscribe", "", "Configure the Jira notifications sent to this channel")
	jira.AddCommand(subscribe)

	view := model.NewAutocompleteData("view", "[issue]", "View the details of a specific Jira issue")
	view.AddDynamicListArgument("List of issues is downloading from your Jira account", "dynamic_issues", true)
	jira.AddCommand(view)

	settings := model.NewAutocompleteData("settings", "[notifications/...]", "Update your user settings")
	notifications := model.NewAutocompleteData("notifications", "[on/off]", "Turn notifications on or off")

	items := []model.AutocompleteListItem{
		{
			Item:     "on",
			HelpText: "Turn notifications on",
		},
		{
			Item:     "off",
			HelpText: "Turn notifications off",
		},
	}
	notifications.AddStaticListArgument("Turn notifications on or off", items, true)
	settings.AddCommand(notifications)
	jira.AddCommand(settings)

	timezone := model.NewAutocompleteData("timezone", "", "Update your timezone")
	timezone.AddNamedTextArgument("zone", "Set timezone", "[UTC+07:00]", "", true)
	jira.AddCommand(timezone)

	install := model.NewAutocompleteData("install", "[cloud/server]", "Connect Mattermost to a Jira instance")
	install.RoleID = model.SYSTEM_ADMIN_ROLE_ID
	cloud := model.NewAutocompleteData("cloud", "[URL]", "Connect to a Jira Cloud instance")
	urlPattern := "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"
	cloud.AddTextArgument("input URL of the Jira Cloud instance", "[URL]", urlPattern)
	install.AddCommand(cloud)
	server := model.NewAutocompleteData("server", "", "Connect to a Jira Server or Data Center instance")
	server.AddTextArgument("input URL of the Jira Server or Data Center instance", "[URL]", urlPattern)
	install.AddCommand(server)
	jira.AddCommand(install)

	uninstall := model.NewAutocompleteData("uninstall", "[cloud/server]", "Disconnect Mattermost from a Jira instance")
	uninstall.RoleID = model.SYSTEM_ADMIN_ROLE_ID
	cloud = model.NewAutocompleteData("cloud", "[URL]", "Disconnect from a Jira Cloud instance")
	cloud.AddTextArgument("input URL of the Jira Cloud instance", "[URL]", urlPattern)
	uninstall.AddCommand(cloud)
	server = model.NewAutocompleteData("server", "[URL]", "Disconnect from a Jira Server or Data Center instance")
	server.AddTextArgument("input URL of the Jira Server or Data Center instance", "[URL]", urlPattern)
	uninstall.AddCommand(server)
	jira.AddCommand(uninstall)

	return jira
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
	case commandTriggerInteractive:
		return p.executeCommandInteractive(args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}
}

func (p *Plugin) executeCommandCrash() *model.CommandResponse {
	go p.crash()
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Crashing plugin",
	}
}

func (p *Plugin) executeCommandHooks(args *model.CommandArgs) *model.CommandResponse {
	configuration := p.getConfiguration()

	if strings.HasSuffix(args.Command, "true") {
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

	if strings.HasSuffix(args.Command, "false") {
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
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
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
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
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

func (p *Plugin) executeCommandInteractive(args *model.CommandArgs) *model.CommandResponse {

	post := &model.Post{
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		UserId:    p.botId,
		Message:   "Test interactive button",
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{{
				Actions: []*model.PostAction{{
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("/plugins/%s/interactive/button/1", manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: "Interactive Button",
				}},
			}},
		},
	}

	_, err := p.API.CreatePost(post)
	if err != nil {
		errorMessage := "Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
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
