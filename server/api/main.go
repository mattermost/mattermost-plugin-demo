// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/itstar-tech/mattermost-plugin-demo/server/app"
)

const (
	headerMattermostUserID = "Mattermost-User-ID"
)

type Handlers struct {
	app       *app.WhatsappApp
	pluginAPI plugin.API
	Router    *mux.Router
}

func New(app *app.WhatsappApp, pluginAPI plugin.API) *Handlers {
	api := &Handlers{
		app:       app,
		pluginAPI: pluginAPI,
	}

	api.initRoutes()
	return api
}

func (api *Handlers) initRoutes() {
	api.Router = mux.NewRouter()

	whatsapp := api.Router.PathPrefix("/whatsapp").Subrouter()

	// Endpoints de sesión Whatsapp
	whatsapp.HandleFunc("/session/{sessionID}", api.handleGetSession).Methods("GET")
	whatsapp.HandleFunc("/session", api.handleCreateSession).Methods("POST")
	whatsapp.HandleFunc("/session/{sessionID}/close", api.handleCloseSession).Methods("POST")
	whatsapp.HandleFunc("/session/by_user/{userID}", api.handleGetSessionByUserId).Methods("GET")
	whatsapp.HandleFunc("/session/unclosed", api.handleGetSessionsUnclosed).Methods("GET")

	// Endpoints de canales Whatsapp
	whatsapp.HandleFunc("/channels", api.handleGetChannels).Methods("GET")
	whatsapp.HandleFunc("/channels", api.handleCreateChannel).Methods("POST")
}

func (api *Handlers) handlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Pong")
}

func ReturnStatusOK(w http.ResponseWriter) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "OK"})
}

func jsonResponse(w http.ResponseWriter, code int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error marshaling data", http.StatusInternalServerError)
		return
	}

	setResponseHeader(w, "Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(bytes)
}

func setResponseHeader(w http.ResponseWriter, key string, value string) { //nolint:unparam
	w.Header().Set(key, value)
}
