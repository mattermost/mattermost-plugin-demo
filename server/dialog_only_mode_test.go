package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newHookTestPlugin returns a plugin whose API fails the test on any call that has no expectation.
func newHookTestPlugin(api *plugintest.API, cfg *configuration) *Plugin {
	p := &Plugin{}
	p.API = api
	p.botID = "bot-id"
	p.setConfiguration(cfg)
	return p
}

func TestDialogOnlyModeSkipsEventHooks(t *testing.T) {
	cfg := &configuration{
		DialogOnlyMode: true,
		Username:       "demo_plugin",
		demoChannelIDs: map[string]string{"team-id": "channel-id"},
	}
	post := &model.Post{UserId: "user-id", ChannelId: "channel-id", Message: "@demo_plugin hello"}
	member := &model.ChannelMember{UserId: "user-id", ChannelId: "channel-id"}
	teamMember := &model.TeamMember{UserId: "user-id", TeamId: "team-id"}
	user := &model.User{Id: "user-id", Username: "someone"}

	hooks := map[string]func(p *Plugin){
		"UserHasBeenCreated":     func(p *Plugin) { p.UserHasBeenCreated(nil, user) },
		"UserHasBeenDeactivated": func(p *Plugin) { p.UserHasBeenDeactivated(nil, user) },
		"UserHasLoggedIn":        func(p *Plugin) { p.UserHasLoggedIn(nil, user) },
		"ChannelHasBeenCreated":  func(p *Plugin) { p.ChannelHasBeenCreated(nil, &model.Channel{TeamId: "team-id"}) },
		"UserHasJoinedChannel":   func(p *Plugin) { p.UserHasJoinedChannel(nil, member, nil) },
		"UserHasLeftChannel":     func(p *Plugin) { p.UserHasLeftChannel(nil, member, nil) },
		"UserHasJoinedTeam":      func(p *Plugin) { p.UserHasJoinedTeam(nil, teamMember, nil) },
		"UserHasLeftTeam":        func(p *Plugin) { p.UserHasLeftTeam(nil, teamMember, nil) },
		"MessageHasBeenPosted":   func(p *Plugin) { p.MessageHasBeenPosted(nil, post) },
		"MessageHasBeenUpdated":  func(p *Plugin) { p.MessageHasBeenUpdated(nil, post, post) },
		"MessageHasBeenDeleted":  func(p *Plugin) { p.MessageHasBeenDeleted(nil, post) },
		"ReactionHasBeenAdded":   func(p *Plugin) { p.ReactionHasBeenAdded(nil, &model.Reaction{UserId: "user-id"}) },
		"ReactionHasBeenRemoved": func(p *Plugin) { p.ReactionHasBeenRemoved(nil, &model.Reaction{UserId: "user-id"}) },
		"BackgroundJob":          func(p *Plugin) { p.BackgroundJob() },
	}

	for name, call := range hooks {
		t.Run(name, func(t *testing.T) {
			api := &plugintest.API{}
			call(newHookTestPlugin(api, cfg.Clone()))
			api.AssertExpectations(t)
		})
	}

	t.Run("MessageWillBePosted lets the post through unchanged", func(t *testing.T) {
		api := &plugintest.API{}
		got, reason := newHookTestPlugin(api, cfg.Clone()).MessageWillBePosted(nil, post)
		assert.Equal(t, post, got)
		assert.Empty(t, reason)
		api.AssertExpectations(t)
	})

	t.Run("FileWillBeUploaded allows the upload without querying teams", func(t *testing.T) {
		api := &plugintest.API{}
		info, reason := newHookTestPlugin(api, cfg.Clone()).FileWillBeUploaded(nil, &model.FileInfo{Name: "a.txt"}, *bytes.NewReader([]byte("a")), &bytes.Buffer{})
		assert.Nil(t, info)
		assert.Empty(t, reason)
		api.AssertExpectations(t)
	})

	t.Run("ConfigurationWillBeSaved saves without querying teams", func(t *testing.T) {
		api := &plugintest.API{}
		newCfg := &model.Config{}
		newCfg.PluginSettings.Plugins = map[string]map[string]any{manifest.Id: {"dialogonlymode": true}}
		got, err := newHookTestPlugin(api, cfg.Clone()).ConfigurationWillBeSaved(newCfg)
		assert.Nil(t, got)
		assert.NoError(t, err)
		api.AssertExpectations(t)
	})
}

func TestPostPluginMessage(t *testing.T) {
	t.Run("posts to the team's demo channel", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("CreatePost", mock.MatchedBy(func(post *model.Post) bool {
			return post.ChannelId == "channel-id" && post.UserId == "bot-id"
		})).Return(&model.Post{}, nil).Once()

		p := newHookTestPlugin(api, &configuration{demoChannelIDs: map[string]string{"team-id": "channel-id"}})
		assert.Nil(t, p.postPluginMessage("team-id", "hello"))
		api.AssertExpectations(t)
	})

	t.Run("skips a team without a demo channel", func(t *testing.T) {
		api := &plugintest.API{}
		p := newHookTestPlugin(api, &configuration{demoChannelIDs: map[string]string{"team-id": "channel-id"}})
		assert.Nil(t, p.postPluginMessage("other-team-id", "hello"))
		api.AssertExpectations(t)
	})
}

func TestEnsureDemoChannelsSkipsTeamThatFails(t *testing.T) {
	api := &plugintest.API{}
	notFound := model.NewAppError("GetChannelByNameForTeamName", "app.channel.get_by_name.missing.app_error", nil, "", http.StatusNotFound)
	exists := model.NewAppError("CreateChannel", "store.sql_channel.save_channel.exists.app_error", nil, "", http.StatusBadRequest)

	api.On("GetTeams").Return([]*model.Team{{Id: "team-1", Name: "one"}, {Id: "team-2", Name: "two"}}, nil).Once()
	api.On("GetChannelByNameForTeamName", "one", "demo_plugin", false).Return((*model.Channel)(nil), notFound).Once()
	api.On("GetChannelByNameForTeamName", "two", "demo_plugin", false).Return((*model.Channel)(nil), notFound).Once()
	api.On("CreateChannel", mock.MatchedBy(func(c *model.Channel) bool { return c.TeamId == "team-1" })).Return((*model.Channel)(nil), exists).Once()
	api.On("CreateChannel", mock.MatchedBy(func(c *model.Channel) bool { return c.TeamId == "team-2" })).Return(&model.Channel{Id: "channel-2"}, nil).Once()
	api.On("LogWarn", "Failed to create demo channel", "team_id", "team-1", "error", mock.Anything).Once()

	p := newHookTestPlugin(api, &configuration{})
	ids, err := p.ensureDemoChannels(&configuration{ChannelName: "demo_plugin"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"team-2": "channel-2"}, ids)
	api.AssertExpectations(t)
}

func TestConfigurationWillBeSavedWithoutPluginSettingsLeavesSaveAlone(t *testing.T) {
	api := &plugintest.API{}
	p := newHookTestPlugin(api, &configuration{DialogOnlyMode: true, Username: "demo_plugin"})

	newCfg := &model.Config{}
	newCfg.PluginSettings.Plugins = map[string]map[string]any{}
	got, err := p.ConfigurationWillBeSaved(newCfg)
	require.NoError(t, err)
	assert.Nil(t, got)
	api.AssertExpectations(t)
}

func TestConfigurationWillBeSavedLeavesLiveConfigurationUnchanged(t *testing.T) {
	api := &plugintest.API{}
	p := newHookTestPlugin(api, &configuration{DialogOnlyMode: true, Username: "demo_plugin"})
	live := p.getConfiguration()

	newCfg := &model.Config{}
	newCfg.PluginSettings.Plugins = map[string]map[string]any{manifest.Id: {"username": "someone_else"}}
	_, err := p.ConfigurationWillBeSaved(newCfg)
	require.NoError(t, err)

	assert.Equal(t, "demo_plugin", live.Username)
	api.AssertExpectations(t)
}

func TestConfigurationWillBeSavedTurningOnDialogOnlyModePostsNothing(t *testing.T) {
	api := &plugintest.API{}
	p := newHookTestPlugin(api, &configuration{Username: "demo_plugin", demoChannelIDs: map[string]string{"team-id": "channel-id"}})

	newCfg := &model.Config{}
	newCfg.PluginSettings.Plugins = map[string]map[string]any{manifest.Id: {"dialogonlymode": true, "username": "demo_plugin"}}
	got, err := p.ConfigurationWillBeSaved(newCfg)
	assert.Nil(t, got)
	assert.NoError(t, err)
	api.AssertExpectations(t)
}

func TestDiffConfigurationLeavesSecretsOutOfAttachment(t *testing.T) {
	api := &plugintest.API{}
	var attachment []byte
	api.On("GetTeams").Return([]*model.Team{{Id: "team-id"}}, nil).Once()
	api.On("UploadFile", mock.MatchedBy(func(data []byte) bool {
		attachment = data
		return true
	}), "channel-id", "configuration.json").Return(&model.FileInfo{Id: "file-id"}, nil).Once()
	api.On("CreatePost", mock.MatchedBy(func(post *model.Post) bool {
		return post.ChannelId == "channel-id" && post.FileIds[0] == "file-id"
	})).Return(&model.Post{}, nil).Once()

	p := newHookTestPlugin(api, &configuration{Username: "demo_plugin"})
	p.diffConfiguration(&configuration{
		Username:       "renamed",
		ServiceAPIKey:  "sk-secret-value",
		RandomSecret:   "random-secret-value",
		demoChannelIDs: map[string]string{"team-id": "channel-id"},
	})

	assert.Contains(t, string(attachment), `"Username":"renamed"`)
	assert.NotContains(t, string(attachment), "sk-secret-value")
	assert.NotContains(t, string(attachment), "random-secret-value")
	api.AssertExpectations(t)
}
