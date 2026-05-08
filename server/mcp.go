package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-plugin-agents/external/pluginmcp"
)

const mcpBasePath = "/mcp"

var (
	mcpNewServer = pluginmcp.NewServer
	mcpRegister  = func(server *pluginmcp.Server) error {
		return server.Register()
	}
	mcpUnregister = func(server *pluginmcp.Server) error {
		return server.Unregister()
	}
)

func (p *Plugin) ensureMCPServer() error {
	p.mcpServerLock.Lock()
	defer p.mcpServerLock.Unlock()

	if p.mcpServer != nil {
		return nil
	}

	if manifest.Id == "" {
		return errors.New("plugin manifest id is required for MCP server")
	}
	if manifest.Version == "" {
		return errors.New("plugin manifest version is required for MCP server")
	}

	serverName := strings.TrimSpace(manifest.Name)
	if serverName == "" {
		return errors.New("plugin manifest name is required for MCP server")
	}

	server := mcpNewServer(p.API, pluginmcp.Config{
		PluginID:       manifest.Id,
		Name:           serverName + " MCP",
		Path:           mcpBasePath,
		ExposeExternal: true,
		Version:        manifest.Version,
	})

	p.registerMCPTools(server)
	p.mcpServer = server
	return nil
}

func (p *Plugin) registerMCPServerBestEffort() {
	server := p.currentMCPServer()
	if server == nil {
		p.API.LogWarn("MCP registration unavailable; continuing plugin activation", "reason", "server not initialized")
		return
	}

	if err := mcpRegister(server); err != nil {
		p.API.LogWarn("MCP registration unavailable; continuing plugin activation", "err", err.Error())
	}
}

func (p *Plugin) unregisterMCPServerBestEffort() {
	server := p.currentMCPServer()
	if server == nil {
		return
	}

	if err := mcpUnregister(server); err != nil {
		p.API.LogWarn("MCP unregister failed; continuing plugin shutdown", "err", err.Error())
	}
}

func (p *Plugin) serveMCPIfMatch(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != mcpBasePath && !strings.HasPrefix(r.URL.Path, mcpBasePath+"/") {
		return false
	}

	server := p.currentMCPServer()
	if server == nil {
		http.NotFound(w, r)
		return true
	}

	server.ServeHTTP(w, r)
	return true
}

func (p *Plugin) currentMCPServer() *pluginmcp.Server {
	p.mcpServerLock.RLock()
	defer p.mcpServerLock.RUnlock()

	return p.mcpServer
}
