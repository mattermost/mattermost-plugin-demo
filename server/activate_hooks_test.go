package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOnActivate(t *testing.T) {
	teamId := model.NewId()
	channelId := model.NewId()
	demoChannelIds := map[string]string{
		teamId: channelId,
	}

	for name, test := range map[string]struct {
		SetupAPI    func() *plugintest.API
		ShouldError bool
	}{
		"GetServerVersion not implemented, returns empty string": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("")

				return api
			},
			ShouldError: true,
		},
		"below minimum supported version: 5.3.9": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.3.9")

				return api
			},
			ShouldError: true,
		},
		"minimum supported version: 5.4.0, but GetTeams fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.4.0")
				api.On("GetTeams").Return(nil, &model.AppError{})

				return api
			},
			ShouldError: true,
		},
		"minimum supported version: 5.4.0, but CreatePost fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.4.0")
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil, &model.AppError{})

				return api
			},
			ShouldError: true,
		},
		"minimum supported version: 5.4.0, but RegisterCommand fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.4.0")
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(&model.AppError{})

				return api
			},
			ShouldError: true,
		},
		"minimum supported version: 5.4.0": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.4.0")
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)

				return api
			},
			ShouldError: false,
		},
		"newer supported version: 5.5.0": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetServerVersion").Return("5.5.0")
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)

				return api
			},
			ShouldError: false,
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
				api.On("UnregisterCommand", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

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
		"RegisterCommand fails": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamId}}, nil)
				api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
				api.On("UnregisterCommand", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.AppError{})

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
