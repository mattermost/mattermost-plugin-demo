package main

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnActivate(t *testing.T) {
	teamID := model.NewId()
	channelID := model.NewId()
	demoChannelIDs := map[string]string{
		teamID: channelID,
	}

	demoUserID := model.NewId()

	for name, test := range map[string]struct {
		SetupAPI     func(*plugintest.API) *plugintest.API
		SetupHelpers func(*plugintest.Helpers) *plugintest.Helpers
		ShouldError  bool
	}{
		"GetServerVersion not implemented, returns empty string": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return("")

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				return helpers
			},
			ShouldError: true,
		},
		"lesser minor version than minimumServerVersion": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				v := semver.MustParse(minimumServerVersion)
				if v.Minor == 0 {
					v.Major--
					v.Minor = 0
					v.Patch = 0
				} else {
					v.Minor--
					v.Patch = 0
				}
				api.On("GetServerVersion").Return(v.String())

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				return helpers
			},
			ShouldError: true,
		},
		"check server config fails, config is incompatible": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("LogError", "Server configuration is not compatible").Return()
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte("{}"), nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)
				api.On("CreateTeamMember", teamID, demoUserID).Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(false, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: false,
		},
		"minimum supported version fulfilled, but RegisterCommand fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(&model.AppError{})

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte("{}"), nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)
				api.On("CreateTeamMember", teamID, demoUserID).Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fulfilled, but GetTeams fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("GetTeams").Return(nil, &model.AppError{})

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte{}, nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fulfilled, but CreatePost fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil, &model.AppError{})

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte{}, nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)
				api.On("CreateTeamMember", teamID, demoUserID).Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fulfilled": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte("{}"), nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)
				api.On("CreateTeamMember", teamID, demoUserID).Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: false,
		},
		"greater minor version than minimumServerVersion": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				v := semver.MustParse(minimumServerVersion)
				require.Nil(t, v.IncrementMinor())
				api.On("GetServerVersion").Return(v.String())
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVGet", "cron_BackgroundJob").Return([]byte("{}"), nil).Maybe()
				api.On("KVSetWithOptions", "cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()
				api.On("KVSetWithOptions", "mutex_cron_BackgroundJob", mock.Anything, mock.AnythingOfType("model.PluginKVSetOptions")).Return(true, nil).Maybe()

				// OnConfigurationChange
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(&model.User{Id: demoUserID, Username: "demo_user"}, nil)
				api.On("CreateTeamMember", teamID, demoUserID).Return(&model.TeamMember{}, nil)
				api.On("GetChannelByNameForTeamName", "", "", false).Return(&model.Channel{}, nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(model.NewId(), nil)

				return helpers
			},
			ShouldError: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			api := test.SetupAPI(&plugintest.API{})
			helpers := test.SetupHelpers(&plugintest.Helpers{})
			defer api.AssertExpectations(t)

			p := Plugin{}
			p.setConfiguration(&configuration{
				demoChannelIDs: demoChannelIDs,
			})
			p.SetAPI(api)
			p.SetHelpers(helpers)
			err := p.OnActivate()

			if test.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestOnDeactivate(t *testing.T) {
	teamID := model.NewId()
	channelID := model.NewId()
	demoChannelIDs := map[string]string{
		teamID: channelID,
	}

	for name, test := range map[string]struct {
		SetupAPI    func() *plugintest.API
		ShouldError bool
	}{
		"all fine": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)

				return api
			},
			ShouldError: false,
		},
		"GetTeam fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetTeams").Return(nil, &model.AppError{})

				return api
			},
			ShouldError: true,
		},
		"CreatePost fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetTeams").Return([]*model.Team{{Id: teamID}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil, &model.AppError{})

				return api
			},
			ShouldError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			api := test.SetupAPI()
			defer api.AssertExpectations(t)

			p := Plugin{}
			p.setConfiguration(&configuration{
				demoChannelIDs: demoChannelIDs,
			})
			p.SetAPI(api)
			err := p.OnDeactivate()

			if test.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
