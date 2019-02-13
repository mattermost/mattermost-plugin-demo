package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) postPluginMessage(id, msg string) *model.AppError {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return nil
	}

	if configuration.EnableMentionUser {
		msg = fmt.Sprintf("tag @%s | %s", configuration.MentionUser, msg)
	}
	msg = fmt.Sprintf("%s%s%s", configuration.TextStyle, msg, configuration.TextStyle)

	_, err := p.API.CreatePost(&model.Post{
		UserId:    configuration.demoUserId,
		ChannelId: configuration.demoChannelIds[id],
		Message:   msg,
	})

	return err
}
