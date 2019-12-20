package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// UserHasBeenCreated is invoked when a new user is created.
//
// This demo implementation logs a message to the demo channel in the team whenever a new user
// is created.
func (p *Plugin) UserHasBeenCreated(c *plugin.Context, user *model.User) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams UserHasBeenCreated",
			"error", err.Error(),
		)
		return
	}

	for _, team := range teams {
		msg := fmt.Sprintf("User_ID @%s has been created in", user.Id)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			p.API.LogError(
				"failed to post UserHasBeenCreated message",
				"channel_id", configuration.demoChannelIds[team.Id],
				"error", err.Error(),
			)
		}
	}
}
