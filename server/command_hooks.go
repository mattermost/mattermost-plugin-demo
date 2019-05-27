package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const CommandTriggerPlugin = "demo_plugin"
const CommandTriggerEphemeral = "ephemeral"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTriggerPlugin,
		AutoComplete:     true,
		AutoCompleteHint: "(true|false)",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Demo Plugin Command",
		Description:      "A command used to enable or disable the demo plugin hooks.",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTriggerEphemeral,
		AutoComplete:     true,
		AutoCompleteHint: "",
		AutoCompleteDesc: "Demonstrates an ephemeral post capabilities.",
		DisplayName:      "Demo Plugin Ephemeral Command",
		Description:      "A command used to demonstrate ephemeral post capabilities.",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
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
	configuration := p.getConfiguration()

	if strings.HasPrefix(args.Command, "/"+CommandTriggerPlugin) {

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

		} else if strings.HasSuffix(args.Command, "false") {
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
	} else if strings.HasPrefix(args.Command, "/"+CommandTriggerEphemeral) {

		siteURL := *p.API.GetConfig().ServiceSettings.SiteURL

		post := &model.Post{
			ChannelId: args.ChannelId,
			Message:   "test ephemeral actions",
			Props: model.StringInterface{
				"attachments": []*model.SlackAttachment{{
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"count": 0,
								},
								URL: fmt.Sprintf("%s/plugins/%s/ephemeral/update", siteURL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Update",
						},
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("%s/plugins/%s/ephemeral/delete", siteURL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Delete",
						},
					},
				}},
			},
		}
		p.API.SendEphemeralPost(args.UserId, post)
		return &model.CommandResponse{}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
	}, nil
}
