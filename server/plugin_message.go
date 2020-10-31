package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

// Helper method for the demo plugin. Posts a message to the "demo" channel
// for the team specified. If the teamID specified is empty, the method
// will post the message to the "demo" channel for each team.
func (p *Plugin) postPluginMessage(teamID, msg string) *model.AppError {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return nil
	}

	if configuration.EnableMentionUser {
		msg = fmt.Sprintf("tag @%s | %s", configuration.MentionUser, msg)
	}
	msg = fmt.Sprintf("%s%s%s", configuration.TextStyle, msg, configuration.TextStyle)

	if teamID != "" {
		_, err := p.API.CreatePost(&model.Post{
			UserId:    p.botID,
			ChannelId: configuration.demoChannelIDs[teamID],
			Message:   msg,
		})
		return err
	}

	for _, channelID := range configuration.demoChannelIDs {
		_, err := p.API.CreatePost(&model.Post{
			UserId:    p.botID,
			ChannelId: channelID,
			Message:   msg,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
