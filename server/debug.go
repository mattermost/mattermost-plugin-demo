// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"

	pluginModel "github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

var (
	createSessionCommand = "createsession"
	closeSessionCommand  = "closesession"
	getSessionCommand    = "getsession"
	listSessionsCommand  = "listsessions"
)

func (p *Plugin) registerDebugCommands() error {
	commands := []string{createSessionCommand, closeSessionCommand, getSessionCommand, listSessionsCommand}
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
	case "/" + listSessionsCommand:
		return p.executeListSessionsCommand(ctx, args)
	}
	return nil, nil
}

func (p *Plugin) executeListSessionsCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	sessions, err := p.store.GetSessions()
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
	session := &pluginModel.Session{
		UserID: args.UserId,
	}

	err := p.store.CreateSession(session)
	if err != nil {
		return &model.CommandResponse{Text: "Failed to create session: " + err.Error()}, nil
	}

	return &model.CommandResponse{Text: "Session created successfully with ID: " + session.ID}, nil
}
