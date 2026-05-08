package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mattermost/mattermost-plugin-agents/external/pluginmcp"
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
	mcpRegister = func(server *pluginmcp.Server) error {
		return registerErr
	}

	api := &plugintest.API{}
	api.On("LogWarn", "MCP registration unavailable; continuing plugin activation", "err", registerErr.Error()).Once()

	plugin := &Plugin{}
	plugin.API = api

	err := plugin.OnActivate()
	require.NoError(t, err)
	require.NotNil(t, plugin.currentMCPServer())
	api.AssertExpectations(t)
}

func TestEnsureMCPServerInitializesOnceConcurrently(t *testing.T) {
	originalMCPNewServer := mcpNewServer
	t.Cleanup(func() {
		mcpNewServer = originalMCPNewServer
	})

	var newServerCalls atomic.Int32
	mcpNewServer = func(api pluginmcp.PluginAPI, config pluginmcp.Config) *pluginmcp.Server {
		newServerCalls.Add(1)
		time.Sleep(10 * time.Millisecond)
		return originalMCPNewServer(api, config)
	}

	plugin := &Plugin{}
	plugin.API = &plugintest.API{}

	const callers = 16
	start := make(chan struct{})
	errs := make(chan error, callers)

	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- plugin.ensureMCPServer()
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	require.NotNil(t, plugin.currentMCPServer())
	require.Equal(t, int32(1), newServerCalls.Load())
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
	mcpUnregister = func(server *pluginmcp.Server) error {
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
	mcpUnregister = func(server *pluginmcp.Server) error {
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
