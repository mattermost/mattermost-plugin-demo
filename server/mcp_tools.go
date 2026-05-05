package main

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-plugin-agents/public/mcphelper"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EchoArgs struct {
	Message string `json:"message" jsonschema:"The string to echo back,minLength=1"`
}

type EchoOutput struct {
	Echoed string `json:"echoed" jsonschema:"The same string that was passed in"`
}

type AddTwoNumbersArgs struct {
	A int `json:"a" jsonschema:"First addend"`
	B int `json:"b" jsonschema:"Second addend"`
}

type AddTwoNumbersOutput struct {
	Sum int `json:"sum" jsonschema:"Sum of a and b"`
}

type GetUserDisplayNameArgs struct{}

type GetUserDisplayNameOutput struct {
	UserID      string `json:"user_id" jsonschema:"Mattermost user ID of the caller"`
	Username    string `json:"username" jsonschema:"Username of the caller"`
	DisplayName string `json:"display_name" jsonschema:"Full-name display name of the caller, falling back to username"`
}

func (p *Plugin) registerMCPTools() {
	mcphelper.AddTool(p.mcpServer, &mcp.Tool{
		Name:        "echo",
		Description: "Echo a string back to the caller. Useful for verifying the MCP round-trip.",
	}, p.echoHandler)

	mcphelper.AddTool(p.mcpServer, &mcp.Tool{
		Name:        "add_two_numbers",
		Description: "Return the sum of two integers. Exercises typed JSON-schema generation.",
	}, p.addTwoNumbersHandler)

	mcphelper.AddTool(p.mcpServer, &mcp.Tool{
		Name:        "get_user_display_name",
		Description: "Look up the calling user's display name. Exercises the X-Mattermost-UserID context propagation chain: server -> agents plugin -> PluginHTTP -> mcphelper.ServeHTTP -> tool handler.",
	}, p.getUserDisplayNameHandler)
}

func (p *Plugin) echoHandler(_ context.Context, _ *mcp.CallToolRequest, in EchoArgs) (*mcp.CallToolResult, EchoOutput, error) {
	return nil, EchoOutput{Echoed: in.Message}, nil
}

func (p *Plugin) addTwoNumbersHandler(_ context.Context, _ *mcp.CallToolRequest, in AddTwoNumbersArgs) (*mcp.CallToolResult, AddTwoNumbersOutput, error) {
	return nil, AddTwoNumbersOutput{Sum: in.A + in.B}, nil
}

func (p *Plugin) getUserDisplayNameHandler(ctx context.Context, _ *mcp.CallToolRequest, _ GetUserDisplayNameArgs) (*mcp.CallToolResult, GetUserDisplayNameOutput, error) {
	userID := mcphelper.GetUserID(ctx)
	if userID == "" {
		return nil, GetUserDisplayNameOutput{}, fmt.Errorf("no Mattermost user ID in tool context (did the request arrive via mcphelper.ServeHTTP?)")
	}

	user, err := p.client.User.Get(userID)
	if err != nil {
		return nil, GetUserDisplayNameOutput{}, fmt.Errorf("failed to fetch user %s: %w", userID, err)
	}

	return nil, GetUserDisplayNameOutput{
		UserID:      user.Id,
		Username:    user.Username,
		DisplayName: user.GetDisplayName(model.ShowFullName),
	}, nil
}
