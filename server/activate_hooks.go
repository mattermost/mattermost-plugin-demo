package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
)

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
func (p *Plugin) OnActivate() error {
	configuration := p.getConfiguration()

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	for _, team := range teams {
		if _, err := p.API.CreatePost(&model.Post{
			UserId:    configuration.demoUserId,
			ChannelId: configuration.demoChannelIds[team.Id],
			Message: fmt.Sprintf(
				"OnActivate: %s", manifest.Id,
			),
			Type: "custom_demo_plugin",
			Props: map[string]interface{}{
				"username":     configuration.Username,
				"channel_name": configuration.ChannelName,
			},
		}); err != nil {
			return errors.Wrap(err, "failed to post OnActivate message")
		}

		if err := p.registerCommand(team.Id); err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	}

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {
	configuration := p.getConfiguration()

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	for _, team := range teams {
		if _, err := p.API.CreatePost(&model.Post{
			UserId:    configuration.demoUserId,
			ChannelId: configuration.demoChannelIds[team.Id],
			Message: fmt.Sprintf(
				"OnDeactivate: %s", manifest.Id,
			),
			Type: "custom_demo_plugin",
			Props: map[string]interface{}{
				"username":     configuration.Username,
				"channel_name": configuration.ChannelName,
			},
		}); err != nil {
			return errors.Wrap(err, "failed to post OnDeactivate message")
		}

		if err := p.registerCommand(team.Id); err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	}

	return nil
}
