package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
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
	case "/dialog/error":
		p.handleDialogWithError(w, r)
	case "/ephemeral/update":
		p.handleEphemeralUpdate(w, r)
	case "/ephemeral/delete":
		p.handleEphemeralDelete(w, r)
	case "/interactive/button/1":
		p.handleInteractiveAction(w, r)
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
		UserId:    p.botId,
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
			UserId:    p.botId,
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

	suffix := ""
	if request.State == dialogStateRelativeCallbackURL {
		suffix = "from relative callback URL"
	}

	if _, appErr = p.API.CreatePost(&model.Post{
		UserId:    p.botId,
		ChannelId: request.ChannelId,
		Message:   fmt.Sprintf("@%v confirmed an Interactive Dialog %v", user.Username, suffix),
	}); appErr != nil {
		p.API.LogError("failed to post handleDialog2 message", "err", appErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDialogWithError(w http.ResponseWriter, r *http.Request) {
	// Always return an error
	response := &model.SubmitDialogResponse{
		Error: "some error",
	}
	p.writeSubmitDialogResponse(w, response)
}

func (p *Plugin) handleEphemeralUpdate(w http.ResponseWriter, r *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(r.Body)

	if request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL
	count := request.Context["count"].(float64) + 1

	post := &model.Post{
		Id:        request.PostId,
		ChannelId: request.ChannelId,
		Message:   "updated ephemeral action",
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{{
				Actions: []*model.PostAction{{
					Integration: &model.PostActionIntegration{
						Context: model.StringInterface{
							"count": count,
						},
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/update", siteURL, manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: fmt.Sprintf("Update %d", int(count)),
				}, {
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/delete", siteURL, manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: "Delete",
				}},
			}},
		},
	}
	p.API.UpdateEphemeralPost(request.UserId, post)

	resp := &model.PostActionIntegrationResponse{}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.ToJson())

}

func (p *Plugin) handleEphemeralDelete(w http.ResponseWriter, r *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(r.Body)

	if request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.API.DeleteEphemeralPost(request.UserId, request.PostId)

	resp := &model.PostActionIntegrationResponse{}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.ToJson())
}

func (p *Plugin) handleInteractiveAction(w http.ResponseWriter, r *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(r.Body)
	if request == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("failed to get user for interactive action", "err", appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, postErr := p.API.GetPost(request.PostId)
	if postErr != nil {
		p.API.LogError("failed to get post for interactive action", "err", postErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rootId := post.RootId
	if rootId == "" {
		rootId = post.Id
	}

	requestJSON, jsonErr := json.MarshalIndent(request, "", "  ")
	if jsonErr != nil {
		p.API.LogError("failed to marshal json for interactive action", "err", jsonErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := "@%v clicked an interactive button.\n```json\n%v\n```"
	if _, appErr := p.API.CreatePost(&model.Post{
		UserId:    p.botId,
		ChannelId: request.ChannelId,
		RootId:    rootId,
		Message:   fmt.Sprintf(msg, user.Username, string(requestJSON)),
	}); appErr != nil {
		p.API.LogError("failed to post handleInteractiveAction message", "err", appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &model.PostActionIntegrationResponse{}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.ToJson())
}

func (p *Plugin) writeSubmitDialogResponse(w http.ResponseWriter, response *model.SubmitDialogResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response.ToJson()); err != nil {
		p.API.LogError("failed to write DialogResponse", "err", err.Error())
	}
}
