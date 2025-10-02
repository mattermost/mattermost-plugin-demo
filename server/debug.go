// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

var (
	createSessionCommand = "createsession"
	closeSessionCommand  = "closesession"
	getSessionCommand    = "getsession"
	listSessionsCommand  = "listsessions"

	createChannelCommand = "createchannel"
	listChannelsCommand  = "listchannels"
)

func (p *Plugin) registerDebugCommands() error {
	commands := []string{createSessionCommand, closeSessionCommand, getSessionCommand, listSessionsCommand, createChannelCommand, listChannelsCommand}
	for _, cmd := range commands {
		err := p.API.RegisterCommand(&model.Command{
			Trigger:      cmd,
			AutoComplete: true,
		})
		if err != nil {
			p.API.LogError("registerDebugCommands: failed to register command", "cmd", cmd, "error", err.Error())
			return errors.Wrap(err, "registerDebugCommands: failed to register command "+cmd)
		}
	}
	return nil
}

func (p *Plugin) ExecuteCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) == 0 {
		return nil, nil
	}
	command := split[0]

	switch command {
	case "/" + createSessionCommand:
		return p.executeCreateSessionCommand(ctx, args)
	case "/" + closeSessionCommand:
		return p.executeCloseSessionCommand(ctx, args)
	case "/" + listSessionsCommand:
		return p.executeListSessionsCommand(ctx, args)
	case "/" + createChannelCommand:
		return p.executeCreateChannelCommand(ctx, args)
	case "/" + listChannelsCommand:
		return p.executeListChannelsCommand(ctx, args)
	}
	return nil, nil
}

func (p *Plugin) executeListSessionsCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	sessions, err := p.app.GetSessions()
	if err != nil {
		return &model.CommandResponse{Text: "Failed to list sessions: " + err.Error()}, nil
	}
	var ids []string
	for _, s := range sessions {
		ids = append(ids, s.ID)
	}
	return &model.CommandResponse{Text: "Unclosed session IDs: " + strings.Join(ids, ", ")}, nil
}

func (p *Plugin) executeCreateSessionCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	session, err := p.app.CreateSession(args.UserId)
	if err != nil {
		return &model.CommandResponse{Text: "Failed to create session: " + err.Error()}, nil
	}

	return &model.CommandResponse{Text: "Session created successfully with ID: " + session.ID}, nil
}

func (p *Plugin) executeCloseSessionCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return &model.CommandResponse{Text: "Usage: /closesession <session_id>"}, nil
	}

	sessionID := split[1]

	session, err := p.app.GetSessionByID(sessionID)
	if err != nil {
		return &model.CommandResponse{Text: "Failed to get session: " + err.Error()}, nil
	}

	if session.ClosedAt != nil {
		return &model.CommandResponse{Text: "Session is already closed"}, nil
	}

	now := model.GetMillis()
	session.ClosedAt = &now

	err = p.app.UpdateSession(session)
	if err != nil {
		return &model.CommandResponse{Text: "Failed to close session: " + err.Error()}, nil
	}

	return &model.CommandResponse{Text: "Session " + sessionID + " closed successfully"}, nil
}

func (p *Plugin) executeCreateChannelCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return &model.CommandResponse{Text: "Usage: /createchannel <channel_id> <phone_number> <phone_number_id>"}, nil
	}
	channelID := split[1]
	phoneNumber := split[2]
	phoneNumberId := split[3]
	channel, err := p.app.CreateChannel(channelID, phoneNumber, phoneNumberId)
	if err != nil {
		return &model.CommandResponse{Text: "Failed to create channel: " + err.Error()}, nil
	}
	return &model.CommandResponse{Text: "Channel created successfully with ID: " + channel.ID}, nil
}

func (p *Plugin) executeListChannelsCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	channels, err := p.app.GetChannels()
	if err != nil {
		return &model.CommandResponse{Text: "Failed to list channels: " + err.Error()}, nil
	}
	if len(channels) == 0 {
		return &model.CommandResponse{Text: "No channels found."}, nil
	}
	var ids []string
	for _, c := range channels {
		ids = append(ids, c.ID)
	}
	return &model.CommandResponse{Text: "Channel IDs: " + strings.Join(ids, ", ")}, nil
}
