package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-plugin-agents/public/bridgeclient"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEchoHandler(t *testing.T) {
	plugin := &Plugin{}

	callResult, out, err := plugin.echoHandler(context.Background(), nil, EchoArgs{Message: "hello, MCP demo!"})
	require.NoError(t, err)
	assert.Nil(t, callResult)
	assert.Equal(t, "hello, MCP demo!", out.Echoed)

	callResult, out, err = plugin.echoHandler(context.Background(), nil, EchoArgs{Message: "  keep spacing  "})
	require.NoError(t, err)
	assert.Nil(t, callResult)
	assert.Equal(t, "  keep spacing  ", out.Echoed)
}

func TestAddTwoNumbersHandler(t *testing.T) {
	plugin := &Plugin{}

	callResult, out, err := plugin.addTwoNumbersHandler(context.Background(), nil, AddTwoNumbersArgs{A: 2, B: 3})
	require.NoError(t, err)
	assert.Nil(t, callResult)
	assert.Equal(t, 5, out.Sum)

	callResult, out, err = plugin.addTwoNumbersHandler(context.Background(), nil, AddTwoNumbersArgs{A: -4, B: 7})
	require.NoError(t, err)
	assert.Nil(t, callResult)
	assert.Equal(t, 3, out.Sum)
}

func TestGetUserDisplayNameHandlerMissingUserContext(t *testing.T) {
	plugin := &Plugin{}

	callResult, out, err := plugin.getUserDisplayNameHandler(context.Background(), nil, GetUserDisplayNameArgs{})
	require.Error(t, err)
	assert.Nil(t, callResult)
	assert.Contains(t, err.Error(), "no Mattermost user ID in tool context")
	assert.Equal(t, GetUserDisplayNameOutput{}, out)
}

func TestGetUserDisplayNameHandlerUserLookupFailure(t *testing.T) {
	api := &plugintest.API{}
	userID := "user-id-123"
	apiErr := model.NewAppError("TestGetUserDisplayNameHandlerUserLookupFailure", "app.user.get", nil, "lookup failed", http.StatusInternalServerError)
	api.On("GetUser", userID).Return((*model.User)(nil), apiErr).Once()

	plugin := newToolTestPlugin(t, api)
	session := newToolTestSession(t, plugin, headerWithUserID(userID))
	toolName := lookupToolName(t, session, "__get_user_display_name")

	out, err := callGetUserDisplayNameTool(t, session, toolName)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch user "+userID)
	assert.Equal(t, GetUserDisplayNameOutput{}, out)
	api.AssertExpectations(t)
}

func TestGetUserDisplayNameHandlerSuccess(t *testing.T) {
	api := &plugintest.API{}
	userID := "user-id-123"
	user := &model.User{
		Id:        userID,
		Username:  "demo_user",
		FirstName: "Demo",
		LastName:  "User",
	}
	api.On("GetUser", userID).Return(user, (*model.AppError)(nil)).Once()

	plugin := newToolTestPlugin(t, api)
	session := newToolTestSession(t, plugin, headerWithUserID(userID))
	toolName := lookupToolName(t, session, "__get_user_display_name")

	out, err := callGetUserDisplayNameTool(t, session, toolName)
	require.NoError(t, err)
	assert.Equal(t, user.Id, out.UserID)
	assert.Equal(t, user.Username, out.Username)
	assert.Equal(t, user.GetDisplayName(model.ShowFullName), out.DisplayName)
	api.AssertExpectations(t)
}

func TestGetUserDisplayNameHandlerFallsBackToUsername(t *testing.T) {
	api := &plugintest.API{}
	userID := "user-id-456"
	user := &model.User{
		Id:       userID,
		Username: "demo_user",
	}
	api.On("GetUser", userID).Return(user, (*model.AppError)(nil)).Once()

	plugin := newToolTestPlugin(t, api)
	session := newToolTestSession(t, plugin, headerWithUserID(userID))
	toolName := lookupToolName(t, session, "__get_user_display_name")

	out, err := callGetUserDisplayNameTool(t, session, toolName)
	require.NoError(t, err)
	assert.Equal(t, user.GetDisplayName(model.ShowFullName), out.DisplayName)
	api.AssertExpectations(t)
}

func newToolTestPlugin(t *testing.T, api *plugintest.API) *Plugin {
	t.Helper()

	plugin := &Plugin{}
	plugin.API = api
	plugin.client = pluginapi.NewClient(api, nil)
	require.NoError(t, plugin.ensureMCPServer())
	return plugin
}

func newToolTestSession(t *testing.T, plugin *Plugin, extraHeaders http.Header) *mcp.ClientSession {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Mattermost-Plugin-ID", bridgeclient.AiPluginID)
		for key, values := range extraHeaders {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}
		plugin.ServeHTTP(nil, w, r)
	}))
	t.Cleanup(ts.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "demo-plugin-test-client", Version: "0.0.1"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: ts.URL + mcpBasePath}, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = session.Close()
	})
	return session
}

func lookupToolName(t *testing.T, session *mcp.ClientSession, suffix string) string {
	t.Helper()

	tools, err := session.ListTools(context.Background(), &mcp.ListToolsParams{})
	require.NoError(t, err)

	for _, tool := range tools.Tools {
		if strings.HasSuffix(tool.Name, suffix) {
			return tool.Name
		}
	}

	t.Fatalf("tool ending with %q not found", suffix)
	return ""
}

func callGetUserDisplayNameTool(t *testing.T, session *mcp.ClientSession, toolName string) (GetUserDisplayNameOutput, error) {
	t.Helper()

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      toolName,
		Arguments: map[string]any{},
	})
	if err != nil {
		return GetUserDisplayNameOutput{}, err
	}

	if result.IsError {
		if len(result.Content) == 0 {
			return GetUserDisplayNameOutput{}, errors.New("tool returned MCP error with empty content")
		}
		text, ok := result.Content[0].(*mcp.TextContent)
		if ok {
			return GetUserDisplayNameOutput{}, errors.New(text.Text)
		}
		return GetUserDisplayNameOutput{}, errors.New("tool returned MCP error content")
	}

	var out GetUserDisplayNameOutput
	payload, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return GetUserDisplayNameOutput{}, err
	}

	err = json.Unmarshal(payload, &out)
	return out, err
}

func headerWithUserID(userID string) http.Header {
	headers := http.Header{}
	headers.Set("X-Mattermost-UserID", userID)
	return headers
}
