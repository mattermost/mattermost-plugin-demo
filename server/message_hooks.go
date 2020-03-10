package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// MessageWillBePosted is invoked when a message is posted by a user before it is committed to the
// database. If you also want to act on edited posts, see MessageWillBeUpdated. Return values
// should be the modified post or nil if rejected and an explanation for the user.
//
// If you don't need to modify or reject posts, use MessageHasBeenPosted instead.
//
// Note that this method will be called for posts created by plugins, including the plugin that created the post.
//
// This demo implementation rejects posts in the demo channel, as well as posts that @-mention
// the demo plugin user.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return post, ""
	}

	// Always allow posts by the demo plugin user and demo plugin bot.
	if post.UserId == p.botId || post.UserId == configuration.demoUserId {
		return post, ""
	}

	// Reject posts by other users in the demo channels, effectively making it read-only.
	for _, channelId := range configuration.demoChannelIds {
		if channelId == post.ChannelId {
			p.API.SendEphemeralPost(post.UserId, &model.Post{
				UserId:    configuration.demoUserId,
				ChannelId: channelId,
				Message:   "Posting is not allowed in this channel.",
			})

			return nil, "disallowing post in demo channel"
		}
	}

	// Reject posts mentioning the demo plugin user.
	if strings.Contains(post.Message, fmt.Sprintf("@%s", configuration.Username)) {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			UserId:    configuration.demoUserId,
			ChannelId: post.ChannelId,
			Message:   "Shh! You must not talk about the demo plugin user.",
		})

		return nil, plugin.DismissPostError
	}

	// Otherwise, allow the post through.
	return post, ""
}

// MessageWillBeUpdated is invoked when a message is updated by a user before it is committed to
// the database. If you also want to act on new posts, see MessageWillBePosted. Return values
// should be the modified post or nil if rejected and an explanation for the user. On rejection,
// the post will be kept in its previous state.
//
// If you don't need to modify or rejected updated posts, use MessageHasBeenUpdated instead.
//
// Note that this method will be called for posts updated by plugins, including the plugin that
// updated the post.
//
// This demo implementation rejects posts that @-mention the demo plugin user.
func (p *Plugin) MessageWillBeUpdated(c *plugin.Context, newPost, oldPost *model.Post) (*model.Post, string) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return newPost, ""
	}

	// Reject posts mentioning the demo plugin user.
	if strings.Contains(newPost.Message, fmt.Sprintf("@%s", configuration.Username)) {
		p.API.SendEphemeralPost(newPost.UserId, &model.Post{
			UserId:    configuration.demoUserId,
			ChannelId: newPost.ChannelId,
			Message:   "You must not talk about the demo plugin user.",
		})

		return nil, "disallowing mention of demo plugin user"
	}

	// Otherwise, allow the post through.
	return newPost, ""
}

// MessageHasBeenPosted is invoked after the message has been committed to the database. If you
// need to modify or reject the post, see MessageWillBePosted Note that this method will be called
// for posts created by plugins, including the plugin that created the post.
//
// This demo implementation logs a message to the demo channel whenever a message is posted,
// unless by the demo plugin user itself.
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	// Ignore posts by the demo plugin user and demo plugin bot.
	if post.UserId == p.botId || post.UserId == configuration.demoUserId {
		return
	}

	user, err := p.API.GetUser(post.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", post.UserId)
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", post.ChannelId)
		return
	}

	msg := fmt.Sprintf("MessageHasBeenPosted: @%s, ~%s", user.Username, channel.Name)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"failed to post MessageHasBeenPosted message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}

	// Check if the Random Secret was posted
	if strings.Contains(post.Message, configuration.RandomSecret) {
		msg = fmt.Sprintf("The random secret %q has been entered by @%s!\n%s",
			configuration.RandomSecret, user.Username, configuration.SecretMessage,
		)
		if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
			p.API.LogError(
				"failed to post random secret message",
				"channel_id", channel.Id,
				"user_id", user.Id,
				"error", err.Error(),
			)
		}
	}

	if strings.Contains(post.Message, strconv.Itoa(configuration.SecretNumber)) {
		msg = fmt.Sprintf("The random number %d has been entered by @%s!",
			configuration.SecretNumber, user.Username)
		if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
			p.API.LogError(
				"failed to post random secret message",
				"channel_id", channel.Id,
				"user_id", user.Id,
				"error", err.Error(),
			)
		}
	}
}

// MessageHasBeenUpdated is invoked after a message is updated and has been updated in the
// database. If you need to modify or reject the post, see MessageWillBeUpdated Note that this
// method will be called for posts created by plugins, including the plugin that created the post.
//
// This demo implementation logs a message to the demo channel whenever a message is updated,
// unless by the demo plugin user itself.
func (p *Plugin) MessageHasBeenUpdated(c *plugin.Context, newPost, oldPost *model.Post) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	// Ignore updates by the demo plugin user.
	if newPost.UserId == configuration.demoUserId {
		return
	}

	user, err := p.API.GetUser(newPost.UserId)
	if err != nil {
		p.API.LogError("failed to query user", "user_id", newPost.UserId)
		return
	}

	channel, err := p.API.GetChannel(newPost.ChannelId)
	if err != nil {
		p.API.LogError("failed to query channel", "channel_id", newPost.ChannelId)
		return
	}

	msg := fmt.Sprintf("MessageHasBeenUpdated: @%s, ~%s", user.Username, channel.Name)
	if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
		p.API.LogError(
			"failed to post MessageHasBeenUpdated message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}

	// Check if the Random Secret was posted
	if strings.Contains(newPost.Message, configuration.RandomSecret) {
		msg = fmt.Sprintf("The random secret %q has been entered by @%s!\n%s",
			configuration.RandomSecret, user.Username, configuration.SecretMessage,
		)
		if err := p.postPluginMessage(channel.TeamId, msg); err != nil {
			p.API.LogError(
				"failed to post random secret message",
				"channel_id", channel.Id,
				"user_id", user.Id,
				"error", err.Error(),
			)
		}
	}
}
