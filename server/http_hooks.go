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
	case "/ephemeral/update":
		p.handleEphemeralUpdate(w, r)
	case "/ephemeral/delete":
		p.handleEphemeralDelete(w, r)
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
	w.Write(responseJSON)
}

func (p *Plugin) handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (p *Plugin) handleEphemeralUpdate(w http.ResponseWriter, r *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(r.Body)
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL

	count := request.Context["count"].(float64) + 1

	post := &model.Post{
		Id:        request.PostId,
		ChannelId: request.ChannelId,
		Message: "updated ephemeral action",
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"count": count,
								},
								URL: fmt.Sprintf("%s/plugins/%s/ephemeral/update", URL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: fmt.Sprintf("Update %d", int(count)),
						},
						{
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{},
								URL:     fmt.Sprintf("%s/plugins/%s/ephemeral/delete", URL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Delete",
						},
					},
				},
			},
		},
	}
	p.API.UpdateEphemeralPost(request.UserId, post)

	resp := &model.PostActionIntegrationResponse{}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.ToJson())

}

func (p *Plugin) handleEphemeralDelete(w http.ResponseWriter, r *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(r.Body)

	post := &model.Post{
		Id: request.PostId,
	}
	p.API.DeleteEphemeralPost(request.UserId, post)

	resp := &model.PostActionIntegrationResponse{}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp.ToJson())
}
