package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-plugin-agents/public/mcphelper"
)

const mcpBasePath = "/mcp"

var (
	mcpRegister = func(server *mcphelper.Server) error {
		return server.Register()
	}
	mcpUnregister = func(server *mcphelper.Server) error {
		return server.Unregister()
	}
)

func (p *Plugin) ensureMCPServer() error {
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

	p.mcpServer = mcphelper.NewServer(p.API, mcphelper.PluginMCPServer{
		PluginID: manifest.Id,
		Name:     serverName + " MCP",
		Path:     mcpBasePath,
		Version:  manifest.Version,
	})

	p.registerMCPTools()
	return nil
}

func (p *Plugin) registerMCPServerBestEffort() {
	if p.mcpServer == nil {
		p.API.LogWarn("MCP registration unavailable; continuing plugin activation", "reason", "server not initialized")
		return
	}

	if err := mcpRegister(p.mcpServer); err != nil {
		p.API.LogWarn("MCP registration unavailable; continuing plugin activation", "err", err.Error())
	}
}

func (p *Plugin) unregisterMCPServerBestEffort() {
	if p.mcpServer == nil {
		return
	}

	if err := mcpUnregister(p.mcpServer); err != nil {
		p.API.LogWarn("MCP unregister failed; continuing plugin shutdown", "err", err.Error())
	}
}

func (p *Plugin) serveMCPIfMatch(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != mcpBasePath && !strings.HasPrefix(r.URL.Path, mcpBasePath+"/") {
		return false
	}

	if p.mcpServer == nil {
		http.NotFound(w, r)
		return true
	}

	p.mcpServer.ServeHTTP(w, r)
	return true
}
