// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"net/http"
)

func (api *Handlers) handleGetChannels(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	channels, err := api.app.GetWhatsappChannels()
	if err != nil {
		http.Error(w, "Failed to get channels", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, channels)
}

func (api *Handlers) handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	var req struct {
		ChannelID string `json:"channel_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ChannelID == "" {
		http.Error(w, "invalid request body or channel_id missing", http.StatusBadRequest)
		return
	}
	channel, err := api.app.CreateWhatsappChannel(req.ChannelID)
	if err != nil {
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, channel)
}
