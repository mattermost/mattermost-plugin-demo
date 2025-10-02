package main

import (
	"database/sql"
	"net/http"
	"sync"

	root "github.com/itstar-tech/mattermost-plugin-demo"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"

	"github.com/itstar-tech/mattermost-plugin-demo/server/api"
	"github.com/itstar-tech/mattermost-plugin-demo/server/model"

	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/itstar-tech/mattermost-plugin-demo/server/app"
	"github.com/itstar-tech/mattermost-plugin-demo/server/store"
	"github.com/mattermost/mattermost/server/public/plugin"
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
	configuration *model.Config

	store       store.Store
	app         *app.WhatsappApp
	apiHandlers *api.Handlers

	whatsappBotID string

	jobs []*cluster.Job
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.apiHandlers.Router.ServeHTTP(w, r)
}

func (p *Plugin) OnActivate() error {
	buildMode := "Production"
	p.API.LogInfo("Starting up User Survey Plugin, build mode: " + buildMode)

	sqlStore, err := p.initStore()
	if err != nil {
		return err
	}

	_app, err := p.initApp(sqlStore, false)
	if err != nil {
		return err
	}

	_api := p.initAPI(_app)

	p.store = sqlStore
	p.app = _app
	p.apiHandlers = _api

	return nil
}

func (p *Plugin) OnDeactivate() error {
	err := p.store.Shutdown()
	if err != nil {
		p.API.LogError("failed to close database connection on plugin deactivation.", "error", err.Error())
		return err
	}

	return nil
}

func (p *Plugin) initStore() (store.Store, error) {
	storeParams, err := p.createStoreParams()
	if err != nil {
		return nil, err
	}

	return store.New(*storeParams)
}

func (p *Plugin) createStoreParams() (*store.Params, error) {
	mmConfig := p.API.GetUnsanitizedConfig()
	db, err := p.getMasterDB()
	if err != nil {
		return nil, err
	}

	return &store.Params{
		DBType:                  *mmConfig.SqlSettings.DriverName,
		ConnectionString:        *mmConfig.SqlSettings.DataSource,
		TablePrefix:             store.TablePrefix,
		SkipMigrations:          false,
		PluginAPI:               p.API,
		DB:                      db,
		Driver:                  p.Driver,
		MigrationTimeoutSeconds: *mmConfig.SqlSettings.MigrationsStatementTimeoutSeconds,
	}, nil
}

func (p *Plugin) getMasterDB() (*sql.DB, error) {
	client := pluginapi.NewClient(p.API, p.Driver)
	db, err := client.Store.GetMasterDB()
	if err != nil {
		p.API.LogError("failed to get master DB", "error", err.Error())
		return nil, err
	}

	return db, nil
}

func (p *Plugin) initApp(store store.Store, debugBuild bool) (*app.WhatsappApp, error) {
	getConfigFunc := func() *model.Config {
		return p.getConfiguration()
	}

	return app.New(p.API, store, getConfigFunc, p.Driver, debugBuild)
}

func (p *Plugin) initAPI(app *app.WhatsappApp) *api.Handlers {
	return api.New(app, p.API)
}
