package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// ReactionHasBeenAdded is invoked after the reaction has been committed to the database.
//
// Note that this method will be called for reactions added by plugins, including the plugin that
// added the reaction.
//
// This demo implementation logs a message to the demo channel whenever a reaction is added to a post.
func (p *Plugin) ReactionHasBeenAdded(c *plugin.Context, reaction *model.Reaction) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(reaction.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", reaction.UserId,
			"error", err.Error(),
		)
		return
	}

	post, err := p.API.GetPost(reaction.PostId)
	if err != nil {
		p.API.LogError(
			"Failed to query post",
			"post_id", reaction.PostId,
			"error", err.Error(),
		)
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError(
			"Failed to query channel",
			"channel_id", post.ChannelId,
			"error", err.Error(),
		)
		return
	}

	postURL := fmt.Sprintf("%s/_redirect/pl/%s", *p.API.GetConfig().ServiceSettings.SiteURL, reaction.PostId)
	msg := fmt.Sprintf("ReactionHasBeenAdded: @%s, :%s:, [<jump to convo>](%s)", user.Username, reaction.EmojiName, postURL)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"Failed to post ReactionHasBeenAdded message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}

// ReactionHasBeenRemoved is invoked after the removal of the reaction has been committed to the database.
//
// Note that this method will be called for reactions removed by plugins, including the plugin that
// removed the reaction.
//
// This demo implementation logs a message to the demo channel whenever reaction is removed from a post.
func (p *Plugin) ReactionHasBeenRemoved(c *plugin.Context, reaction *model.Reaction) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(reaction.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", reaction.UserId,
			"error", err.Error(),
		)
		return
	}

	post, err := p.API.GetPost(reaction.PostId)
	if err != nil {
		p.API.LogError(
			"Failed to query post",
			"post_id", reaction.PostId,
			"error", err.Error(),
		)
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError(
			"Failed to query channel",
			"channel_id", post.ChannelId,
			"error", err.Error(),
		)
		return
	}

	postURL := fmt.Sprintf("%s/_redirect/pl/%s", *p.API.GetConfig().ServiceSettings.SiteURL, reaction.PostId)
	msg := fmt.Sprintf("ReactionHasBeenRemoved: @%s, :%s:, [<jump to convo>](%s)", user.Username, reaction.EmojiName, postURL)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"Failed to post ReactionHasBeenRemoved message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}
