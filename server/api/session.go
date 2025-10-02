// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Handlers) handleGetSession(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	vars := mux.Vars(r)
	sessionID, ok := vars["sessionID"]
	if !ok {
		http.Error(w, "missing session ID in request", http.StatusBadRequest)
		return
	}
	session, err := api.app.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, session)
}

func (api *Handlers) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
		http.Error(w, "invalid request body or user_id missing", http.StatusBadRequest)
		return
	}
	session, err := api.app.CreateSession(req.UserID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, session)
}

func (api *Handlers) handleCloseSession(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	if err := api.RequireSystemAdmin(w, r); err != nil {
		return
	}
	vars := mux.Vars(r)
	sessionID, ok := vars["sessionID"]
	if !ok {
		http.Error(w, "missing session ID in request", http.StatusBadRequest)
		return
	}
	session, err := api.app.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	if session.ClosedAt.Unix() != 0 {
		http.Error(w, "Session already closed", http.StatusBadRequest)
		return
	}
	_, err = api.app.CloseSession(sessionID)
	if err != nil {
		http.Error(w, "Failed to close session", http.StatusInternalServerError)
		return
	}
	ReturnStatusOK(w)
}

func (api *Handlers) handleGetSessionByUserId(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		http.Error(w, "missing user ID in request", http.StatusBadRequest)
		return
	}
	session, err := api.app.GetSessionByUserId(userID)
	if err != nil {
		http.Error(w, "Failed to get session for user", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, session)
}

func (api *Handlers) handleGetSessionsUnclosed(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}
	sessions, err := api.app.GetSessionsUnclosed()
	if err != nil {
		http.Error(w, "Failed to get unclosed sessions", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, sessions)
}
