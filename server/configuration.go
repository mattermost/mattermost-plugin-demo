package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
type configuration struct {
	// The user to use as part of the demo plugin, created automatically if it does not exist.
	Username string

	// The channel to use as part of the demo plugin, created for each team automatically if it does not exist.
	ChannelName string

	// LastName is the last name of the demo user.
	LastName string

	// TextStyle controls the text style of the messages posted by the demo user.
	TextStyle string

	// RandomSecret is a generated key that, when mentioned in a message by a user, will trigger the demo user to post the 'SecretMessage'.
	RandomSecret string

	// SecretMessage is the message posted to the demo channel when the 'RandomSecret' is pasted somewhere in the team.
	SecretMessage string

	// EnableMentionUser controls whether the 'MentionUser' is prepended to all demo messages or not.
	EnableMentionUser bool

	// MentionUser is the user that is prepended to demo messages when enabled.
	MentionUser string

	// SecretNumber is an integer that, when mentioned in a message by a user, will trigger the demo user to post a message.
	SecretNumber int

	// disabled tracks whether or not the plugin has been disabled after activation. It always starts enabled.
	disabled bool

	// demoUserId is the id of the user specified above.
	demoUserId string

	// demoChannelIds maps team ids to the channels created for each using the channel name above.
	demoChannelIds map[string]string
}

// Clone deep copies the configuration. Your implementation may only require a shallow copy if
// your configuration has no reference types.
func (c *configuration) Clone() *configuration {
	// Deep copy demoChannelIds, a reference type.
	demoChannelIds := make(map[string]string)
	for key, value := range c.demoChannelIds {
		demoChannelIds[key] = value
	}

	return &configuration{
		Username:          c.Username,
		ChannelName:       c.ChannelName,
		LastName:          c.LastName,
		TextStyle:         c.TextStyle,
		RandomSecret:      c.RandomSecret,
		SecretMessage:     c.SecretMessage,
		EnableMentionUser: c.EnableMentionUser,
		MentionUser:       c.MentionUser,
		disabled:          c.disabled,
		demoUserId:        c.demoUserId,
		demoChannelIds:    demoChannelIds,
	}
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) diffConfiguration(newConfiguration *configuration) {
	oldConfiguration := p.getConfiguration()
	configurationDiff := make(map[string]interface{})

	if newConfiguration.Username != oldConfiguration.Username {
		configurationDiff["username"] = newConfiguration.Username
	}
	if newConfiguration.ChannelName != oldConfiguration.ChannelName {
		configurationDiff["channel_name"] = newConfiguration.ChannelName
	}
	if newConfiguration.LastName != oldConfiguration.LastName {
		configurationDiff["lastname"] = newConfiguration.LastName
	}
	if newConfiguration.TextStyle != oldConfiguration.TextStyle {
		configurationDiff["text_style"] = newConfiguration.ChannelName
	}
	if newConfiguration.RandomSecret != oldConfiguration.RandomSecret {
		configurationDiff["random_secret"] = "<HIDDEN>"
	}
	if newConfiguration.SecretMessage != oldConfiguration.SecretMessage {
		configurationDiff["secret_message"] = newConfiguration.SecretMessage
	}
	if newConfiguration.EnableMentionUser != oldConfiguration.EnableMentionUser {
		configurationDiff["enable_mention_user"] = newConfiguration.EnableMentionUser
	}
	if newConfiguration.MentionUser != oldConfiguration.MentionUser {
		configurationDiff["mention_user"] = newConfiguration.MentionUser
	}
	if newConfiguration.SecretNumber != oldConfiguration.SecretNumber {
		configurationDiff["secret_number"] = newConfiguration.SecretNumber
	}

	if len(configurationDiff) == 0 {
		return
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogWarn("failed to query teams OnConfigChange", "err", err)
		return
	}

	for _, team := range teams {
		demoChannelId, ok := newConfiguration.demoChannelIds[team.Id]
		if !ok {
			p.API.LogWarn("No demo channel id for team", "team", team.Id)
			continue
		}

		newConfigurationData, jsonErr := json.Marshal(newConfiguration)
		if jsonErr != nil {
			p.API.LogWarn("failed to marshal new configuration", "err", err)
			return
		}

		fileInfo, err := p.API.UploadFile(newConfigurationData, demoChannelId, "configuration.json")
		if err != nil {
			p.API.LogWarn("failed to attach new configuration", "err", err)
			return
		}

		if _, err := p.API.CreatePost(&model.Post{
			UserId:    p.botId,
			ChannelId: demoChannelId,
			Message:   "OnConfigChange: loading new configuration",
			Type:      "custom_demo_plugin",
			Props:     configurationDiff,
			FileIds:   model.StringArray{fileInfo.Id},
		}); err != nil {
			p.API.LogWarn("failed to post OnConfigChange message", "err", err)
			return
		}
	}
}

// OnConfigurationChange is invoked when configuration changes may have been made.
//
// This demo implementation ensures the configured demo user and channel are created for use
// by the plugin.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)
	var err error

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	configuration.demoUserId, err = p.ensureDemoUser(configuration)
	if err != nil {
		return errors.Wrap(err, "failed to ensure demo user")
	}

	botId, ensureBotError := p.Helpers.EnsureBot(&model.Bot{
		Username:    "demoplugin",
		DisplayName: "Demo Plugin Bot",
		Description: "A bot account created by the demo plugin.",
	}, plugin.IconImagePath("/assets/github.svg"))
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure demo bot.")
	}

	p.botId = botId

	configuration.demoChannelIds, err = p.ensureDemoChannels(configuration)
	if err != nil {
		return errors.Wrap(err, "failed to ensure demo channels")
	}

	p.diffConfiguration(configuration)

	p.setConfiguration(configuration)

	return nil
}

func (p *Plugin) ensureDemoUser(configuration *configuration) (string, error) {
	var err *model.AppError

	// Check for the configured user. Ignore any error, since it's hard to distinguish runtime
	// errors from a user simply not existing.
	user, _ := p.API.GetUserByUsername(configuration.Username)

	// Ensure the configured user exists.
	if user == nil {
		user, err = p.API.CreateUser(&model.User{
			Username:  configuration.Username,
			Password:  "password",
			Email:     fmt.Sprintf("%s@example.com", configuration.Username),
			Nickname:  "Demo Day",
			FirstName: "Demo",
			LastName:  configuration.LastName,
			Position:  "Bot",
		})

		if err != nil {
			return "", err
		}
	}

	if user.LastName != configuration.LastName {
		user.LastName = configuration.LastName
		user, err = p.API.UpdateUser(user)

		if err != nil {
			return "", err
		}
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return "", err
	}

	for _, team := range teams {
		// Ignore any error.
		p.API.CreateTeamMember(team.Id, configuration.demoUserId)
	}

	return user.Id, nil
}

func (p *Plugin) ensureDemoChannels(configuration *configuration) (map[string]string, error) {
	teams, err := p.API.GetTeams()
	if err != nil {
		return nil, err
	}

	demoChannelIds := make(map[string]string)
	for _, team := range teams {
		// Check for the configured channel. Ignore any error, since it's hard to
		// distinguish runtime errors from a channel simply not existing.
		channel, _ := p.API.GetChannelByNameForTeamName(team.Name, configuration.ChannelName, false)

		// Ensure the configured channel exists.
		if channel == nil {
			channel, err = p.API.CreateChannel(&model.Channel{
				TeamId:      team.Id,
				Type:        model.CHANNEL_OPEN,
				DisplayName: "Demo Plugin",
				Name:        configuration.ChannelName,
				Header:      "The channel used by the demo plugin.",
				Purpose:     "This channel was created by a plugin for testing.",
			})

			if err != nil {
				return nil, err
			}
		}

		// Save the ids for later use.
		demoChannelIds[team.Id] = channel.Id
	}

	return demoChannelIds, nil
}

// setEnabled wraps setConfiguration to configure if the plugin is enabled.
func (p *Plugin) setEnabled(enabled bool) {
	var configuration = p.getConfiguration().Clone()
	configuration.disabled = !enabled

	p.setConfiguration(configuration)
}
