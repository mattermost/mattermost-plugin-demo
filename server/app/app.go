// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	"github.com/itstar-tech/mattermost-plugin-demo/server/store"
)

type WhatsappApp struct {
	api       plugin.API
	store     store.Store
	getConfig func() *model.Config
	apiClient *pluginapi.Client
	botID     string
}

func New(
	api plugin.API,
	store store.Store,
	getConfigFunc func() *model.Config,
	driver plugin.Driver,
) (*WhatsappApp, error) {
	app := &WhatsappApp{
		api:       api,
		store:     store,
		getConfig: getConfigFunc,
		apiClient: pluginapi.NewClient(api, driver),
	}

	return app, nil
}
