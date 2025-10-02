package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

// Helper method for the demo plugin. Posts a message to the "demo" channel
// for the team specified. If the teamID specified is empty, the method
// will post the message to the "demo" channel for each team.
func (p *Plugin) postPluginMessage(teamID, msg string) *model.AppError {
	configuration := p.getConfiguration()

	if configuration.Disabled {
		return nil
	}

	if configuration.EnableMentionUser {
		msg = fmt.Sprintf("tag @%s | %s", configuration.MentionUser, msg)
	}
	msg = fmt.Sprintf("%s%s%s", configuration.TextStyle, msg, configuration.TextStyle)

	return nil
}
