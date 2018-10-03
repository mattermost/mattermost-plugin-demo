package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOnActivate(t *testing.T) {
	teamID := model.NewId()

	api1 := &plugintest.API{}
	api1.On("GetServerVersion").Return("5.4.0")
	api1.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamID}}, nil)
	api1.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
	api1.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
	defer api1.AssertExpectations(t)

	api2 := &plugintest.API{}
	api2.On("GetServerVersion").Return("5.5.0")
	api2.On("GetTeams").Return([]*model.Team{&model.Team{Id: teamID}}, nil)
	api2.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(&model.Post{}, nil)
	api2.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
	defer api2.AssertExpectations(t)

	api3 := &plugintest.API{}
	api3.On("GetServerVersion").Return("5.3.0")
	defer api3.AssertExpectations(t)

	api4 := &plugintest.API{}
	api4.On("GetServerVersion").Return("5.4.0")
	api4.On("GetTeams").Return(nil, &model.AppError{})
	api4.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	defer api4.AssertExpectations(t)

	for name, test := range map[string]struct {
		API         *plugintest.API
		ShouldError bool
	}{
		"All fine, current server version the same as required": {
			API:         api1,
			ShouldError: false,
		},
		"All fine, current server version newer as required": {
			API:         api2,
			ShouldError: false,
		},
		"Current server version is to low": {
			API:         api3,
			ShouldError: true,
		},
		"GetTeams fails": {
			API:         api4,
			ShouldError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			p := Plugin{}
			p.SetAPI(test.API)
			err := p.OnActivate()

			if test.ShouldError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})
	}

}
