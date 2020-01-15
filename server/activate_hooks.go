package main

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
)

const minimumServerVersion = "5.12.0"

func (p *Plugin) checkServerVersion() error {
	serverVersion, err := semver.Parse(p.API.GetServerVersion())
	if err != nil {
		return errors.Wrap(err, "failed to parse server version")
	}

	r := semver.MustParseRange(">=" + minimumServerVersion)
	if !r(serverVersion) {
		return fmt.Errorf("this plugin requires Mattermost v%s or later", minimumServerVersion)
	}

	return nil
}

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	if err := p.checkServerVersion(); err != nil {
		return err
	}

	if ok, err := p.checkRequiredServerConfiguration(); err != nil {
		return errors.Wrap(err, "could not check required server configuration")
	} else if !ok {
		return errors.New("server configuration is not compatible")
	}

	configuration := p.getConfiguration()

	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	for _, team := range teams {
		_, ok := configuration.demoChannelIds[team.Id]
		if !ok {
			p.API.LogWarn("No demo channel id for team", "team", team.Id)
			continue
		}

		msg := fmt.Sprintf("OnActivate: %s", manifest.Id)
		if err := p.postPluginMessage(team.Id, msg); err != nil {
			return errors.Wrap(err, "failed to post OnActivate message")
		}
	}

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {
	configuration := p.getConfiguration()

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	for _, team := range teams {
		_, ok := configuration.demoChannelIds[team.Id]
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

func (p *Plugin) checkRequiredServerConfiguration() (bool, error) {
	path, err := p.API.GetBundlePath()
	if err != nil {
		return false, errors.Wrap(err, "failed to get bundle path")
	}

	manifest, _, err := model.FindManifest(path)
	if err != nil {
		return false, errors.Wrap(err, "failed to find manifest file")
	}

	return p.Helpers.CheckRequiredServerConfiguration(manifest.RequiredConfig)
}
