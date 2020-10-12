package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// UserHasJoinedTeam is invoked after the membership has been committed to the database. If
// actor is not nil, the user was added to the team by the actor.
//
// This demo implementation logs a message to the demo channel in the team whenever a user
// joins the team.
func (p *Plugin) UserHasJoinedTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(teamMember.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
		return
	}

	msg := fmt.Sprintf("UserHasJoinedTeam: @%s", user.Username)
	if err := p.postPluginMessage(teamMember.TeamId, msg); err != nil {
		p.API.LogError(
			"Failed to post UserHasJoinedTeam message",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
	}
}

// UserHasLeftTeam is invoked after the membership has been removed from the database. If actor
// is not nil, the user was removed from the team by the actor.
//
// This demo implementation logs a message to the demo channel in the team whenever a user
// leaves the team.
func (p *Plugin) UserHasLeftTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(teamMember.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
		return
	}

	msg := fmt.Sprintf("UserHasLeftTeam: @%s", user.Username)
	if err := p.postPluginMessage(teamMember.TeamId, msg); err != nil {
		p.API.LogError(
			"Failed to post UserHasLeftTeam message",
			"user_id", teamMember.UserId,
			"error", err.Error(),
		)
	}
}
