// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (api *Handlers) handleUpdateSession(w http.ResponseWriter, r *http.Request) {
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

	session, err := api.app.GetSessionByID(sessionID)
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
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}

	if err := api.RequireSystemAdmin(w, r); err != nil {
		return
	}

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

	jsonResponse(w, http.StatusCreated, session)
}

func (api *Handlers) handleGetSession(w http.ResponseWriter, r *http.Request) {
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

	session, err := api.app.GetSessionByID(sessionID)
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, session)
}

func (api *Handlers) handleListSessions(w http.ResponseWriter, r *http.Request) {
	if err := api.RequireAuthentication(w, r); err != nil {
		return
	}

	if err := api.RequireSystemAdmin(w, r); err != nil {
		return
	}

	sessions, err := api.app.GetSessions()
	if err != nil {
		http.Error(w, "Failed to list sessions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, sessions)
}

func (api *Handlers) handleDeleteSession(w http.ResponseWriter, r *http.Request) {
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

	if err := api.app.DeleteSession(sessionID); err != nil {
		http.Error(w, "Failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ReturnStatusOK(w)
}
