package main

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// UserHasBeenCreated is invoked when a new user is created.
//
// This demo implementation logs a message to the demo channel in the team whenever a new user
// is created.
func (p *Plugin) UserHasBeenCreated(c *plugin.Context, user *model.User) {
	configuration := p.getConfiguration()

	if configuration.Disabled {
		return
	}

}

// UserHasBeenDeactivated is invoked when a user is made inactive.
//
// This demo implementation logs a message to the demo channel in the team whenever a user
// is deactivated.
func (p *Plugin) UserHasBeenDeactivated(c *plugin.Context, user *model.User) {
	configuration := p.getConfiguration()

	if configuration.Disabled {
		return
	}

}
