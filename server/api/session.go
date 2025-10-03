// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"

	MattermostModel "github.com/mattermost/mattermost/server/public/model"
)

const (
	WebsocketEventPreferenceUpdated = "whatsapp_preference_updated"
	WebSocketEventSessionUpdated    = "whatsapp_session_updated"
)

type SessionsResponse struct {
	ActiveUsers []*MattermostModel.User `json:"active_users"`
}

func (api *Handlers) handleUpdateSession(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}

	if err := api.RequireSystemAdmin(w, r); err != nil {
		return
	}

	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok || userID == "" {
		http.Error(w, "missing user ID in request", http.StatusBadRequest)
		return
	}

	session, err := api.app.GetSessionByUserId(userID)
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusNotFound)
		return
	}

	var updateData model.Session
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	session.UserID = updateData.UserID
	session.ClosedAt = updateData.ClosedAt

	if err := api.app.UpdateSession(session); err != nil {
		http.Error(w, "Failed to update session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, session)
}

func (api *Handlers) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID string `json:"userID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	session, err := api.app.CreateSession(requestData.UserID)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.PublishSessionUpdateEvent(requestData.UserID)
	if err != nil {
		return
	}

	jsonResponse(w, http.StatusCreated, session)
}

func (api *Handlers) handleGetSessionByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok || userID == "" {
		http.Error(w, "missing user ID in request", http.StatusBadRequest)
		return
	}

	sess, err := api.app.GetSessionByUserId(userID)
	if err != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		jsonResponse(w, http.StatusOK, nil)
		return
	}

	jsonResponse(w, http.StatusOK, sess)
}

func (api *Handlers) handleListActiveUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := api.app.GetActiveUsers()
	if err != nil {
		http.Error(w, "Failed to list active users: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp := SessionsResponse{ActiveUsers: users}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (api *Handlers) handleCloseSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok || userID == "" {
		http.Error(w, "missing user ID in request", http.StatusBadRequest)
		return
	}

	if err := api.app.CloseSessionsFromUserId(userID); err != nil {
		http.Error(w, "Failed to close session(s): "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := api.PublishSessionUpdateEvent(userID); err != nil {
		return
	}

	var closedAt interface{} = nil
	if sess, err := api.app.GetSessionByUserId(userID); err == nil {
		if sess.ClosedAt != nil {
			closedAt = *sess.ClosedAt
		}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"status":    "OK",
		"closed_at": closedAt,
	})
}

func (api *Handlers) PublishPreferenceUpdateEvent() error {
	activeUsers, err := api.app.GetActiveUsers()
	if err != nil {
		return errors.Wrap(err, "failed to get sessions from Mattermost API")
	}

	activeUsersJson := SessionsResponse{
		ActiveUsers: activeUsers,
	}

	jsonData, err := json.Marshal(activeUsersJson)
	if err != nil {
		return errors.Wrap(err, "failed to marshal active users to JSON")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return errors.Wrap(err, "failed to unmarshal active users to JSON")
	}

	api.pluginAPI.PublishWebSocketEvent(
		WebsocketEventPreferenceUpdated,
		payload,
		&MattermostModel.WebsocketBroadcast{},
	)
	return nil
}

func (api *Handlers) PublishSessionUpdateEvent(userID string) error {
	session, _ := api.app.GetSessionByUserId(userID)

	var payload map[string]interface{}

	if session == nil {
		payload = map[string]interface{}{
			"session": nil,
			"user_id": userID,
		}
	} else {
		jsonData, err := json.Marshal(session)
		if err != nil {
			return errors.Wrap(err, "failed to marshal session to JSON")
		}

		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return errors.Wrap(err, "failed to unmarshal session to JSON")
		}
	}

	api.pluginAPI.PublishWebSocketEvent(
		WebSocketEventSessionUpdated,
		payload,
		&MattermostModel.WebsocketBroadcast{UserId: userID},
	)
	return nil
}
