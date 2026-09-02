package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func postMmBlocks(t *testing.T, path string, body any) *http.Response {
	return postMmBlocksWithPlugin(t, &Plugin{}, path, body)
}

func postMmBlocksWithPlugin(t *testing.T, plugin *Plugin, path string, body any) *http.Response {
	t.Helper()

	plugin.initializeAPI()

	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, path, &buf)
	r.Header.Set("Content-Type", "application/json")
	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	require.NotNil(t, result)
	return result
}

func postMmBlocksWithOpenDialogAPI(t *testing.T, path string, body any) *http.Response {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v4/actions/dialogs/open", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	}))
	t.Cleanup(srv.Close)

	siteURL := srv.URL
	api := &plugintest.API{}
	api.On("GetConfig").Return(&model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	})
	t.Cleanup(func() { api.AssertExpectations(t) })

	plugin := &Plugin{}
	plugin.API = api
	return postMmBlocksWithPlugin(t, plugin, path, body)
}

func decodeMmBlocksResponse(t *testing.T, result *http.Response) mmBlocksActionResponse {
	t.Helper()
	defer result.Body.Close()

	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, "application/json", result.Header.Get("Content-Type"))

	var resp mmBlocksActionResponse
	require.NoError(t, json.NewDecoder(result.Body).Decode(&resp))
	return resp
}

func TestMmBlocksIntegration(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration", map[string]any{"user_name": "alice"})
	resp := decodeMmBlocksResponse(t, result)

	assert.Equal(t, "Demo mm_blocks integration OK (user: alice).", resp.EphemeralText)
	assert.True(t, resp.SkipSlackParsing)
}

func TestMmBlocksIntegrationUpdate(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration_update", map[string]any{})
	resp := decodeMmBlocksResponse(t, result)

	require.NotNil(t, resp.Update)
	assert.Equal(t, "E2E mm_blocks post updated (message field).", resp.Update.Message)
	assert.Equal(t, mmBlocksUpdatedMarker, resp.Update.Props["mm_blocks"].([]any)[0].(map[string]any)["text"])
}

func TestMmBlocksEchoQuery(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration_echo_query?zeta=z&alpha=a", map[string]any{})
	resp := decodeMmBlocksResponse(t, result)

	assert.Equal(t, "Demo mm_blocks query OK (alpha=a&zeta=z)", resp.EphemeralText)
}

func TestMmBlocksEchoContext(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration_echo_context", map[string]any{
		"context": map[string]any{"test_marker": "hello"},
	})
	resp := decodeMmBlocksResponse(t, result)

	assert.Equal(t, "Demo mm_blocks context OK (test_marker: hello).", resp.EphemeralText)
}

func TestMmBlocksStaticSelect(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration_static_select", map[string]any{
		"context": map[string]any{"selected_option": "opt_beta"},
	})
	resp := decodeMmBlocksResponse(t, result)

	assert.Equal(t, "Demo mm_blocks static_select OK (selected_option: opt_beta).", resp.EphemeralText)
}

func TestMmBlocksEchoFormValues(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_integration_echo_form_values", map[string]any{
		"context": map[string]any{
			"form_values": map[string]any{
				"title":  "Demo ticket",
				"labels": []any{"a", "b"},
			},
		},
	})
	resp := decodeMmBlocksResponse(t, result)

	assert.Equal(t, "Demo mm_blocks form_values OK (labels=a,b&title=Demo ticket)", resp.EphemeralText)
}

func TestMmBlocksLookup(t *testing.T) {
	t.Run("unfiltered", func(t *testing.T) {
		result := postMmBlocks(t, "/mm_blocks_integration_lookup", map[string]any{})
		resp := decodeMmBlocksResponse(t, result)
		require.Len(t, resp.Items, 4)
		assert.Equal(t, "Alpha", resp.Items[0].Text)
	})

	t.Run("query param", func(t *testing.T) {
		result := postMmBlocks(t, "/mm_blocks_integration_lookup?query=alpha", map[string]any{})
		resp := decodeMmBlocksResponse(t, result)
		require.Len(t, resp.Items, 1)
		assert.Equal(t, "opt_alpha", resp.Items[0].Value)
	})
}

func TestMmBlocksDialogReturn(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		result := postMmBlocks(t, "/mm_blocks_dialog_return", map[string]any{
			"context": map[string]any{"marker": "abc"},
		})
		resp := decodeMmBlocksResponse(t, result)
		assert.Equal(t, "dialog", resp.Type)
		require.NotNil(t, resp.BlockDialog)
		assert.Equal(t, "Demo Blocks (return)", resp.BlockDialog.Title)
		assert.Contains(t, resp.BlockDialog.Actions, mmBlocksActionSubmit)
		assert.Equal(t, pluginPath("/mm_blocks_dialog_submit"), resp.BlockDialog.Actions[mmBlocksActionSubmit].URL)
	})

	t.Run("scenario simple", func(t *testing.T) {
		result := postMmBlocks(t, "/mm_blocks_dialog_return", map[string]any{
			"context": map[string]any{"scenario": "simple"},
		})
		resp := decodeMmBlocksResponse(t, result)
		require.NotNil(t, resp.BlockDialog)
		assert.Equal(t, "Demo Simple Dialog", resp.BlockDialog.Title)
	})
}

func TestMmBlocksDialogOpenRequiresTriggerID(t *testing.T) {
	result := postMmBlocks(t, "/mm_blocks_dialog_open", map[string]any{})
	resp := decodeMmBlocksResponse(t, result)
	assert.Equal(t, "ok", resp.Type)
	assert.Equal(t, "mm_blocks_dialog_open requires trigger_id", resp.Error)
}

func TestMmBlocksDialogOpen(t *testing.T) {
	result := postMmBlocksWithOpenDialogAPI(t, "/mm_blocks_dialog_open", map[string]any{
		"trigger_id": "trig",
		"context":    map[string]any{"marker": "abc"},
	})
	resp := decodeMmBlocksResponse(t, result)
	assert.Equal(t, "ok", resp.Type)
	assert.Empty(t, resp.Error)
	assert.Nil(t, resp.BlockDialog)
}

func TestMmBlocksDialogSubmitAndCancel(t *testing.T) {
	submit := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_submit", map[string]any{
		"context": map[string]any{
			"step": "2",
			"form_values": map[string]any{
				"notes": "hi",
			},
		},
	}))
	assert.Equal(t, "Demo mm_blocks dialog submit OK step=2 (notes=hi)", submit.EphemeralText)

	cancel := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_cancel", map[string]any{}))
	assert.Equal(t, "Demo mm_blocks dialog cancelled (reason=cancel)", cancel.EphemeralText)
}

func TestMmBlocksDialogRefreshAndMultistep(t *testing.T) {
	refresh := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_refresh", map[string]any{
		"context": map[string]any{"form_values": map[string]any{"title": "My ticket"}},
	}))
	assert.Equal(t, "refresh", refresh.Type)
	require.NotNil(t, refresh.BlockDialog)
	assert.Equal(t, "Step 2", refresh.BlockDialog.Title)

	step1 := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_multistep", map[string]any{
		"context": map[string]any{"step": "1"},
	}))
	assert.Equal(t, "Step 2 - Work Info", step1.BlockDialog.Title)

	step2 := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_multistep", map[string]any{
		"context": map[string]any{"step": "2"},
	}))
	assert.Equal(t, "Step 3 - Final Details", step2.BlockDialog.Title)
}

func TestMmBlocksDialogErrorsAndGoto(t *testing.T) {
	errorsResp := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_errors", map[string]any{}))
	assert.Equal(t, "Title looks wrong", errorsResp.Errors["title"])

	errorResp := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_error", map[string]any{}))
	assert.Equal(t, "Demo mm_blocks dialog top-level error", errorResp.Error)

	gotoResp := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_goto", map[string]any{}))
	assert.Equal(t, "/", gotoResp.GotoLocation)
}

func TestMmBlocksDialogChild(t *testing.T) {
	missing := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_child", map[string]any{}))
	assert.Equal(t, "mm_blocks_dialog_child requires trigger_id", missing.Error)
	assert.True(t, missing.KeepDialogOpen)

	child := decodeMmBlocksResponse(t, postMmBlocksWithOpenDialogAPI(t, "/mm_blocks_dialog_child", map[string]any{
		"trigger_id": "trig",
		"context":    map[string]any{"source": "Details"},
	}))
	assert.Equal(t, "ok", child.Type)
	assert.Empty(t, child.Error)
	assert.True(t, child.KeepDialogOpen)
	assert.Nil(t, child.BlockDialog)
}

func TestMmBlocksDialogFieldRefresh(t *testing.T) {
	opts := mmBlocksDialogOptions{Title: "Custom Refresh", Marker: "mark-1"}
	initial := getMmBlocksFieldRefreshDialog("web", "Acme", opts)
	require.Contains(t, initial.Actions, mmBlocksActionFieldRefresh)
	assert.Equal(t, "Custom Refresh", initial.Actions[mmBlocksActionFieldRefresh].Context["title"])
	assert.Equal(t, "mark-1", initial.Actions[mmBlocksActionFieldRefresh].Context["marker"])

	resp := decodeMmBlocksResponse(t, postMmBlocks(t, "/mm_blocks_dialog_field_refresh", map[string]any{
		"context": map[string]any{
			"title":       "Custom Refresh",
			"marker":      "mark-1",
			"form_values": map[string]any{"project_type": "web", "project_name": "Acme"},
		},
	}))
	assert.Equal(t, "refresh", resp.Type)
	require.NotNil(t, resp.BlockDialog)
	assert.Equal(t, "Custom Refresh", resp.BlockDialog.Title)
	require.Contains(t, resp.BlockDialog.Actions, mmBlocksActionFieldRefresh)
	assert.Equal(t, "Custom Refresh", resp.BlockDialog.Actions[mmBlocksActionFieldRefresh].Context["title"])
	assert.Equal(t, "mark-1", resp.BlockDialog.Actions[mmBlocksActionFieldRefresh].Context["marker"])

	foundFramework := false
	foundMarker := false
	for _, block := range resp.BlockDialog.Blocks {
		b, ok := block.(map[string]any)
		if !ok {
			continue
		}
		if b["name"] == "framework" {
			foundFramework = true
		}
		if text, _ := b["text"].(string); strings.Contains(text, "mark-1") {
			foundMarker = true
		}
	}
	assert.True(t, foundFramework)
	assert.True(t, foundMarker)
}

func TestMmBlocksUnknownMethodIsNotFoundOrNotAllowed(t *testing.T) {
	plugin := &Plugin{}
	plugin.initializeAPI()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/mm_blocks_integration", nil)
	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	defer result.Body.Close()
	_, _ = io.Copy(io.Discard, result.Body)
	assert.Contains(t, []int{http.StatusNotFound, http.StatusMethodNotAllowed}, result.StatusCode)
}

func TestPostOpenBlockDialog(t *testing.T) {
	var got struct {
		TriggerID   string          `json:"trigger_id"`
		BlockDialog *mmBlocksDialog `json:"block_dialog"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v4/actions/dialogs/open", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	}))
	t.Cleanup(srv.Close)

	dialog := getMmBlocksDialog(mmBlocksDialogOptions{Title: "Demo Blocks (open)", Marker: "abc"})
	require.NoError(t, postOpenBlockDialog(srv.Client(), srv.URL, "trig-1", dialog))
	assert.Equal(t, "trig-1", got.TriggerID)
	require.NotNil(t, got.BlockDialog)
	assert.Equal(t, "Demo Blocks (open)", got.BlockDialog.Title)
	assert.Contains(t, got.BlockDialog.Actions, mmBlocksActionSubmit)
}
