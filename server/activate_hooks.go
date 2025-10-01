package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/itstar-tech/mattermost-plugin-demo/server/store"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"
)

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account

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

func (p *Plugin) OnActivate() error {
	if p.client == nil {
		p.client = pluginapi.NewClient(p.API, p.Driver)
	}

	if err := p.checkRequiredServerConfiguration(); err != nil {
		return errors.Wrap(err, "server configuration is not compatible")
	}

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	sqlStore, err := p.initStore()
	if err != nil {
		return err
	}

	p.store = sqlStore

	p.initializeAPI()

	configuration := p.getConfiguration()

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	for _, team := range teams {
		_, ok := configuration.demoChannelIDs[team.Id]
		if !ok {
			p.API.LogWarn("No demo channel id for team", "team", team.Id)
			continue
		}

		msg := fmt.Sprintf("OnActivate: %s", manifest.Id)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			return errors.Wrap(err, "failed to post OnActivate message")
		}
	}

	job, cronErr := cluster.Schedule(
		p.API,
		"BackgroundJob",
		cluster.MakeWaitForRoundedInterval(15*time.Minute),
		p.BackgroundJob,
	)
	if cronErr != nil {
		return errors.Wrap(cronErr, "failed to schedule background job")
	}
	p.backgroundJob = job

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {
	configuration := p.getConfiguration()

	// TODO: Close database store
	//if p.store != nil {
	//	if sqlStore, ok := p.store.(*store.SQLStore); ok {
	//		if err := sqlStore.Close(); err != nil {
	//			p.API.LogError("Failed to close database store", "err", err)
	//		}
	//	}
	//}

	if p.backgroundJob != nil {
		if err := p.backgroundJob.Close(); err != nil {
			p.API.LogError("Failed to close background job", "err", err)
		}
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	for _, team := range teams {
		_, ok := configuration.demoChannelIDs[team.Id]
		if !ok {
			p.API.LogWarn("No demo channel id for team", "team", team.Id)
			continue
		}

		msg := fmt.Sprintf("OnDeactivate: %s", manifest.Id)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			return errors.Wrap(err, "failed to post OnDeactivate message")
		}
	}

	return nil
}

func (p *Plugin) checkRequiredServerConfiguration() error {
	config := p.client.Configuration.GetConfig()
	if config.ServiceSettings.EnableGifPicker == nil || !*config.ServiceSettings.EnableGifPicker {
		return errors.New("ServiceSettings.EnableGifPicker must be enabled")
	}

	if config.FileSettings.EnablePublicLink == nil || !*config.FileSettings.EnablePublicLink {
		return errors.New("FileSettings.EnablePublicLink must be enabled")
	}

	return nil
}
