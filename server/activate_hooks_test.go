package main

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost-plugin-agents/public/mcphelper"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/require"
)

func TestOnActivateContinuesWhenMCPRegistrationFails(t *testing.T) {
	originalOnActivateCore := onActivateCore
	originalMCPRegister := mcpRegister
	t.Cleanup(func() {
		onActivateCore = originalOnActivateCore
		mcpRegister = originalMCPRegister
	})

	onActivateCore = func(_ *Plugin) error {
		return nil
	}

	registerErr := errors.New("agents plugin unavailable")
	mcpRegister = func(server *mcphelper.Server) error {
		return registerErr
	}

	api := &plugintest.API{}
	api.On("LogWarn", "MCP registration unavailable; continuing plugin activation", "err", registerErr.Error()).Once()

	plugin := &Plugin{}
	plugin.API = api

	err := plugin.OnActivate()
	require.NoError(t, err)
	require.NotNil(t, plugin.mcpServer)
	api.AssertExpectations(t)
}

func TestOnDeactivateContinuesWhenMCPUnregisterFails(t *testing.T) {
	originalOnDeactivateCore := onDeactivateCore
	originalMCPUnregister := mcpUnregister
	t.Cleanup(func() {
		onDeactivateCore = originalOnDeactivateCore
		mcpUnregister = originalMCPUnregister
	})

	onDeactivateCore = func(_ *Plugin) error {
		return nil
	}

	unregisterErr := errors.New("agents plugin already stopped")
	mcpUnregister = func(server *mcphelper.Server) error {
		return unregisterErr
	}

	api := &plugintest.API{}
	api.On("LogWarn", "MCP unregister failed; continuing plugin shutdown", "err", unregisterErr.Error()).Once()

	plugin := &Plugin{}
	plugin.API = api
	require.NoError(t, plugin.ensureMCPServer())

	err := plugin.OnDeactivate()
	require.NoError(t, err)
	api.AssertExpectations(t)
}

func TestOnDeactivatePreservesCoreErrorWhenMCPUnregisterFails(t *testing.T) {
	originalOnDeactivateCore := onDeactivateCore
	originalMCPUnregister := mcpUnregister
	t.Cleanup(func() {
		onDeactivateCore = originalOnDeactivateCore
		mcpUnregister = originalMCPUnregister
	})

	coreErr := errors.New("core shutdown failed")
	onDeactivateCore = func(_ *Plugin) error {
		return coreErr
	}

	unregisterErr := errors.New("agents plugin already stopped")
	mcpUnregister = func(server *mcphelper.Server) error {
		return unregisterErr
	}

	api := &plugintest.API{}
	api.On("LogWarn", "MCP unregister failed; continuing plugin shutdown", "err", unregisterErr.Error()).Once()

	plugin := &Plugin{}
	plugin.API = api
	require.NoError(t, plugin.ensureMCPServer())

	err := plugin.OnDeactivate()
	require.ErrorIs(t, err, coreErr)
	api.AssertExpectations(t)
}
