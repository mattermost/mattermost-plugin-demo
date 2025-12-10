package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"
)

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	if p.client == nil {
		p.client = pluginapi.NewClient(p.API, p.Driver)
	}

	if err := p.checkRequiredServerConfiguration(); err != nil {
		return errors.Wrap(err, "server configuration is not compatible")
	}

	p.initializeAPI()

	configuration := p.getConfiguration()

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	// Skip team messages and background job in minimal mode
	if configuration.DialogOnlyMode {
		p.API.LogInfo("Demo plugin activated in minimal mode (dialog command only)")
		return nil
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	for _, team := range teams {
		_, ok := configuration.demoChannelIDs[team.Id]
		if !ok {
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

	// Skip cleanup in minimal mode
	if configuration.DialogOnlyMode {
		return nil
	}

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
