package main

import (
	"sync"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"

	root "github.com/mattermost/mattermost-plugin-demo"
)

var (
	manifest model.Manifest = root.Manifest
)

type Plugin struct {
	plugin.MattermostPlugin
	client *pluginapi.Client

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	router *mux.Router

	// BotId of the created bot account.
	botID string

	// backgroundJob is a job that executes periodically on only one plugin instance at a time
	backgroundJob *cluster.Job

	// Session tracking
	sessionToConn   map[string]string
	sessionToConnMu sync.RWMutex
}
