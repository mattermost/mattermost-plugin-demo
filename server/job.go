package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) BackgroundJob() {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	for _, channelId := range configuration.demoChannelIds {
		_, err := p.API.CreatePost(&model.Post{
			UserId:    p.botId,
			ChannelId: channelId,
			Message:   "Background job executed",
		})
		if err != nil {
			p.API.LogError(
				"failed to post BackgroundJob message",
				"channel_id", channelId,
				"error", err.Error(),
			)
		}
	}
}
