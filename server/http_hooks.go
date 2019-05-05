package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
//
// The Mattermost-User-Id header will be present if (and only if) the request is by an
// authenticated user.
//
// This demo implementation sends back whether or not the plugin hooks are currently enabled. It
// is used by the web app to recover from a network reconnection and synchronize the state of the
// plugin's hooks.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/status":
		p.handleStatus(w, r)
	case "/hello":
		p.handleHello(w, r)
	case "/dialog/1":
		p.handleDialog1(w, r)
	case "/dialog/2":
		p.handleDialog2(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleStatus(w http.ResponseWriter, r *http.Request) {
	configuration := p.getConfiguration()

	var response = struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: !configuration.disabled,
	}

	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) handleHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World!")); err != nil {
		p.API.LogError("failed to write hello world", "err", err.Error())
	}
}

func (p *Plugin) handleDialog1(w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	request := model.SubmitDialogRequestFromJson(r.Body)
	if request == nil {
		p.API.LogError("failed to decode SubmitDialogRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.Cancelled {
		number, ok := request.Submission[dialogElementNameNumber].(float64)
		if !ok {
			p.API.LogError("request is missing field", "field", dialogElementNameNumber)
			w.WriteHeader(http.StatusOK)
			return
		}

		if number != 42 {
			response := &model.SubmitDialogResponse{
				Errors: map[string]string{
					dialogElementNameNumber: "This must be 42",
				},
			}
			p.writeSubmitDialogResponse(w, response)
			return
		}
	}

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("failed to get user for dialog", "err", appErr.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	msg := "@%v submitted an Interative Dialog"
	if request.Cancelled {
		msg = "@%v cancelled an Interative Dialog"
	}

	rootPost, appErr := p.API.CreatePost(&model.Post{
		UserId:    config.demoUserId,
		ChannelId: request.ChannelId,
		Message:   fmt.Sprintf(msg, user.Username),
	})
	if appErr != nil {
		p.API.LogError("failed to post handleDialog1 message", "err", appErr.Error())
		return
	}

	if !request.Cancelled {
		// Don't post the email address publicly
		request.Submission[dialogElementNameEmail] = "xxxxxxxxxxx"

		if _, appErr = p.API.CreatePost(&model.Post{
			UserId:    config.demoUserId,
			ChannelId: request.ChannelId,
			RootId:    rootPost.Id,
			Message:   "Data:",
			Type:      "custom_demo_plugin",
			Props:     request.Submission,
		}); appErr != nil {
			p.API.LogError("failed to post handleDialog1 message", "err", appErr.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDialog2(w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	request := model.SubmitDialogRequestFromJson(r.Body)
	if request == nil {
		p.API.LogError("failed to decode SubmitDialogRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("failed to get user for dialog", "err", appErr.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	if _, appErr = p.API.CreatePost(&model.Post{
		UserId:    config.demoUserId,
		ChannelId: request.ChannelId,
		Message:   fmt.Sprintf("@%v confirmed an Interactive Dialog", user.Username),
	}); appErr != nil {
		p.API.LogError("failed to post handleDialog1 message", "err", appErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) writeSubmitDialogResponse(w http.ResponseWriter, response *model.SubmitDialogResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response.ToJson()); err != nil {
		p.API.LogError("failed to write DialogResponse", "err", err.Error())
	}
}
