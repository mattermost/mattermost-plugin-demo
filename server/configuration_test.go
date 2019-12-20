package main

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfiguration(t *testing.T) {
	t.Run("null configuration", func(t *testing.T) {
		plugin := &Plugin{}
		assert.NotNil(t, plugin.getConfiguration())
	})

	t.Run("changing configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{disabled: false}
		plugin.setConfiguration(configuration1)
		assert.Equal(t, configuration1, plugin.getConfiguration())

		configuration2 := &configuration{disabled: true}
		plugin.setConfiguration(configuration2)
		assert.Equal(t, configuration2, plugin.getConfiguration())
		assert.NotEqual(t, configuration1, plugin.getConfiguration())
		assert.False(t, plugin.getConfiguration() == configuration1)
		assert.True(t, plugin.getConfiguration() == configuration2)
	})

	t.Run("setting same configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{}
		plugin.setConfiguration(configuration1)
		assert.Panics(t, func() {
			plugin.setConfiguration(configuration1)
		})
	})

	t.Run("clearing configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{disabled: true}
		plugin.setConfiguration(configuration1)
		assert.NotPanics(t, func() {
			plugin.setConfiguration(nil)
		})
		assert.NotNil(t, plugin.getConfiguration())
		assert.NotEqual(t, configuration1, plugin.getConfiguration())
	})
}

func TestOnConfigurationChange(t *testing.T) {
	user := &model.User{
		Id:       model.NewId(),
		Username: "demo_user",
	}
	teamId := model.NewId()
	channelId := model.NewId()
	demoChannelIds := map[string]string{
		teamId: channelId,
	}
	var apiConfiguration = new(configuration)

	for name, test := range map[string]struct {
		SetupAPI         func() *plugintest.API
		SetupHelpers     func() *plugintest.Helpers
		preConfiguration *configuration
		ShouldError      bool
	}{
		"same configuration": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", apiConfiguration).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
				api.On("CreateTeamMember", teamId, "").Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func() *plugintest.Helpers {
				helpers := &plugintest.Helpers{}
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)
				return helpers
			},
			preConfiguration: apiConfiguration,
			ShouldError:      false,
		},
		"different configuration": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", apiConfiguration).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
				api.On("CreateTeamMember", teamId, "").Return(&model.TeamMember{}, nil)
				channel := &model.Channel{
					Id: model.NewId(),
				}
				api.On("GetChannelByNameForTeamName", "", "", false).Return(channel, nil)
				api.On("UploadFile", mock.AnythingOfType("[]uint8"), channel.Id, "configuration.json").Return(&model.FileInfo{}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

				return api
			},
			SetupHelpers: func() *plugintest.Helpers {
				helpers := &plugintest.Helpers{}
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)
				return helpers
			},
			preConfiguration: &configuration{EnableMentionUser: true},
			ShouldError:      false,
		},
		"failure to ensure bot": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", apiConfiguration).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
				api.On("CreateTeamMember", teamId, "").Return(&model.TeamMember{}, nil)

				return api
			},
			SetupHelpers: func() *plugintest.Helpers {
				helpers := &plugintest.Helpers{}
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return("", errors.New("some error"))
				return helpers
			},
			preConfiguration: &configuration{EnableMentionUser: true},
			ShouldError:      true,
		},
		"bot icon exists": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", apiConfiguration).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
				api.On("CreateTeamMember", teamId, "").Return(&model.TeamMember{}, nil)
				channel := &model.Channel{
					Id: model.NewId(),
				}
				api.On("GetChannelByNameForTeamName", "", "", false).Return(channel, nil)
				api.On("UploadFile", mock.AnythingOfType("[]uint8"), channel.Id, "configuration.json").Return(&model.FileInfo{}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.AssertNotCalled(t, "SetBotIconImage", mock.Anything, mock.Anything)
				return api
			},
			SetupHelpers: func() *plugintest.Helpers {
				helpers := &plugintest.Helpers{}
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)
				return helpers
			},
			preConfiguration: &configuration{EnableMentionUser: true},
			ShouldError:      false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			api := test.SetupAPI()
			defer api.AssertExpectations(t)

			helpers := test.SetupHelpers()
			defer helpers.AssertExpectations(t)

			p := Plugin{}
			p.setConfiguration(&configuration{
				demoChannelIds: demoChannelIds,
			})
			p.SetAPI(api)
			p.SetHelpers(helpers)

			// The configuration set here allows us to test calling the
			// "OnConfigurationChange" hook from multiple states
			p.setConfiguration(test.preConfiguration)
			err := p.OnConfigurationChange()

			if test.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
