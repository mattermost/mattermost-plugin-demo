package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
)

// KV keys are kept well under the 50-char limit (channel IDs are 26 chars).
func csSchemaKey(channelID string) string { return "cs_schema_" + channelID }
func csCustomKey(channelID string) string { return "cs_custom_" + channelID }

// customState is the typed payload for the fully-custom channel settings tab.
type customState struct {
	Note        string `json:"note"`
	PinGreeting bool   `json:"pinGreeting"`
}

// requireUser enforces that the request is from an authenticated user. It
// returns false if it already wrote an error response.
func (p *Plugin) requireUser(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Mattermost-User-ID") == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// requireChannelReader enforces auth and read access to the channel for read
// paths. It returns the authenticated user ID, or false if it already wrote an
// error response.
func (p *Plugin) requireChannelReader(w http.ResponseWriter, r *http.Request, channelID string) (string, bool) {
	if !p.requireUser(w, r) {
		return "", false
	}
	userID := r.Header.Get("Mattermost-User-ID")

	if !p.API.HasPermissionToChannel(userID, channelID, model.PermissionReadChannel) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return "", false
	}

	return userID, true
}

// requireChannelManager enforces auth and channel-properties permission for
// write paths. It returns the authenticated user ID, or false if it already
// wrote an error response.
func (p *Plugin) requireChannelManager(w http.ResponseWriter, r *http.Request, channelID string) (string, bool) {
	if !p.requireUser(w, r) {
		return "", false
	}
	userID := r.Header.Get("Mattermost-User-ID")

	channel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		p.API.LogError("Failed to get channel for channel settings", "channel_id", channelID, "err", appErr.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", false
	}

	permission := model.PermissionManagePublicChannelProperties
	if channel.Type == model.ChannelTypePrivate {
		permission = model.PermissionManagePrivateChannelProperties
	}

	if !p.API.HasPermissionToChannel(userID, channelID, permission) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return "", false
	}

	return userID, true
}

func (p *Plugin) handleGetChannelSettingsSchema(w http.ResponseWriter, r *http.Request) {
	channelID := mux.Vars(r)["channel_id"]
	if _, ok := p.requireChannelReader(w, r, channelID); !ok {
		return
	}

	values := map[string]string{}
	if err := p.client.KV.Get(csSchemaKey(channelID), &values); err != nil {
		p.API.LogError("Failed to get channel settings schema", "channel_id", channelID, "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if values == nil {
		values = map[string]string{}
	}

	p.writeJSON(w, values)
}

func (p *Plugin) handleSaveChannelSettingsSchema(w http.ResponseWriter, r *http.Request) {
	channelID := mux.Vars(r)["channel_id"]
	if _, ok := p.requireChannelManager(w, r, channelID); !ok {
		return
	}

	var values map[string]string
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if _, err := p.client.KV.Set(csSchemaKey(channelID), values); err != nil {
		p.API.LogError("Failed to save channel settings schema", "channel_id", channelID, "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	p.writeJSON(w, values)
}

func (p *Plugin) handleGetChannelSettingsCustom(w http.ResponseWriter, r *http.Request) {
	channelID := mux.Vars(r)["channel_id"]
	if _, ok := p.requireChannelReader(w, r, channelID); !ok {
		return
	}

	var state customState
	if err := p.client.KV.Get(csCustomKey(channelID), &state); err != nil {
		p.API.LogError("Failed to get channel settings custom state", "channel_id", channelID, "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	p.writeJSON(w, state)
}

func (p *Plugin) handleSaveChannelSettingsCustom(w http.ResponseWriter, r *http.Request) {
	channelID := mux.Vars(r)["channel_id"]
	if _, ok := p.requireChannelManager(w, r, channelID); !ok {
		return
	}

	var state customState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if _, err := p.client.KV.Set(csCustomKey(channelID), state); err != nil {
		p.API.LogError("Failed to save channel settings custom state", "channel_id", channelID, "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	p.writeJSON(w, state)
}
