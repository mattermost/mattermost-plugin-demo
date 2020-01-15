package main

import (
	"testing"

	"github.com/blang/semver"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnActivate(t *testing.T) {
	teamId := model.NewId()
	channelId := model.NewId()
	demoChannelIds := map[string]string{
		teamId: channelId,
	}

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
		"check server config fails, could not read manifest": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("GetBundlePath").Return("", nil)

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
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(false, nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fullfiled, but RegisterCommand fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(&model.AppError{})
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fullfiled, but GetTeams fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return(nil, &model.AppError{})
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fullfiled, but CreatePost fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil, &model.AppError{})
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)

				return helpers
			},
			ShouldError: true,
		},
		"minimum supported version fullfiled": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetServerVersion").Return(minimumServerVersion)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)

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
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("GetBundlePath").Return("../", nil)

				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("CheckRequiredServerConfiguration", mock.AnythingOfType("*model.Config")).Return(true, nil)

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
				demoChannelIds: demoChannelIds,
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
	teamId := model.NewId()
	channelId := model.NewId()
	demoChannelIds := map[string]string{
		teamId: channelId,
	}

	for name, test := range map[string]struct {
		SetupAPI    func() *plugintest.API
		ShouldError bool
	}{
		"all fine": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
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
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
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
				demoChannelIds: demoChannelIds,
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
