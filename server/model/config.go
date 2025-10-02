// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

type Config struct {
	// The user to use as part of the demo plugin, created automatically if it does not exist.
	Username string `json:"username"`

	// The channel to use as part of the demo plugin, created for each team automatically if it does not exist.
	ChannelName string `json:"channelname"`

	// LastName is the last name of the demo user.
	LastName string `json:"lastname"`

	// TextStyle controls the text style of the messages posted by the demo user.
	TextStyle string `json:"textstyle"`

	// RandomSecret is a generated key that, when mentioned in a message by a user, will trigger the demo user to post the 'SecretMessage'.
	RandomSecret string `json:"randomsecret"`

	// SecretMessage is the message posted to the demo channel when the 'RandomSecret' is pasted somewhere in the team.
	SecretMessage string `json:"secretmessage"`

	// EnableMentionUser controls whether the 'MentionUser' is prepended to all demo messages or not.
	EnableMentionUser bool `json:"enablementionuser"`

	// MentionUser is the user that is prepended to demo messages when enabled.
	MentionUser string `json:"mentionuser"`

	// SecretNumber is an integer that, when mentioned in a message by a user, will trigger the demo user to post a message.
	SecretNumber int `json:"secretnumber"`

	// A deplay in seconds that is applied to Slash Command responses, Post Actions responses and Interactive Dialog responses.
	// It's useful for testing.
	IntegrationRequestDelay int `json:"integrationrequestdelay"`

	// WhatsApp outgoing reaction webhook
	WebhookURL string `json:"webhookurl"`

	// disabled tracks whether the plugin has been disabled after activation. It always starts enabled.
	Disabled bool `json:"-"`

	// demoUserID is the id of the user specified above.
	DemoUserID string `json:"-"`

	// demoChannelIDs maps team ids to the channels created for each using the channel name above.
	DemoChannelIDs map[string]string `json:"-"`

	// WhatsApp monitored channels
	MonitoredChannels map[string]bool `json:"-"`

	// whatsAppAccessToken is the access token for the bot
	WhatsAppAccessToken string `json:"-"`

	// assistantAccessToken is the access token for the Assistant bot
	AssistantAccessToken string `json:"-"`

	EnabledUsers map[string]*mmModel.User `json:"-"`

	EnableAutoAssignment bool `json:"enableautoassignment"`
}

// Clone shallow copies the Configuration. Your implementation may require a deep copy if
// your Configuration has reference types.
func (c *Config) Clone() *Config {
	var clone = *c
	return &clone
}
