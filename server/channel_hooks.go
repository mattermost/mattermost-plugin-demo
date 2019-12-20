package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// ChannelHasBeenCreated is invoked after the channel has been committed to the database.
//
// This demo implementation logs a message to the demo channel whenever a channel is created.
func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	msg := fmt.Sprintf("ChannelHasBeenCreated: ~%s", channel.Name)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"failed to post ChannelHasBeenCreated message",
			"channel_id", channel.Id,
			"error", err.Error(),
		)
	}
}

// UserHasJoinedChannel is invoked after the membership has been committed to the database. If
// actor is not nil, the user was invited to the channel by the actor.
//
// This demo implementation logs a message to the demo channel whenever a user joins a channel.
func (p *Plugin) UserHasJoinedChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(channelMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", channelMember.UserId)
		return
	}

	channel, err := p.API.GetChannel(channelMember.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", channelMember.ChannelId)
		return
	}

	msg := fmt.Sprintf("UserHasJoinedChannel: @%s, ~%s", user.Username, channel.Name)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"failed to post UserHasJoinedChannel message",
			"user_id", channelMember.UserId,
			"error", err.Error(),
		)
	}
}

// UserHasLeftChannel is invoked after the membership has been removed from the database. If
// actor is not nil, the user was removed from the channel by the actor.
//
// This demo implementation logs a message to the demo channel whenever a user leaves a
// channel.
func (p *Plugin) UserHasLeftChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(channelMember.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", channelMember.UserId)
		return
	}

	channel, err := p.API.GetChannel(channelMember.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", channelMember.ChannelId)
		return
	}

	msg := fmt.Sprintf("UserHasLeftChannel: @%s, ~%s", user.Username, channel.Name)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"failed to post UserHasLeftChannel message",
			"user_id", channelMember.UserId,
			"error", err.Error(),
		)
	}
}
