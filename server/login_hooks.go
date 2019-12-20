package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// UserWillLogIn before the login of the user is returned. Returning a non empty string will reject the login event.
// If you don't need to reject the login event, see UserHasLoggedIn
//
// This demo implementation rejects login attempts by the demo user.
func (p *Plugin) UserWillLogIn(c *plugin.Context, user *model.User) string {
	configuration := p.getConfiguration()

	if user.Username == configuration.Username {
		return "the demo user is not allowed to login"
	}

	return ""
}

// UserHasLoggedIn is invoked after a user has logged in.
//
// This demo implementation logs a message to the demo channel whenever a user logs in.
func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
	configuration := p.getConfiguration()

	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams UserHasLoggedIn",
			"error", err.Error(),
		)
		return
	}

	for _, team := range teams {
		msg := fmt.Sprintf("User @%s has logged in", user.Username)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			p.API.LogError(
				"failed to post UserHasLoggedIn message",
				"channel_id", configuration.demoChannelIds[team.Id],
				"error", err.Error(),
			)
		}
	}
}
