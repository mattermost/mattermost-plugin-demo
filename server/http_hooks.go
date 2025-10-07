package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
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
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) initializeAPI() {
	router := mux.NewRouter()

	router.HandleFunc("/status", p.handleStatus)
	router.HandleFunc("/hello", p.handleHello)
	router.HandleFunc("/dynamic_arg_test_url", p.handleDynamicArgTest)
	router.HandleFunc("/check_auth_header", p.handleCheckAuthHeader)

	webhook := router.PathPrefix("/webhook").Subrouter()
	webhook.Use(p.withDelay)
	webhook.HandleFunc("/outgoing", p.handleOutgoingWebhook).Methods(http.MethodPost)

	interativeRouter := router.PathPrefix("/interactive").Subrouter()
	interativeRouter.Use(p.withDelay)
	interativeRouter.HandleFunc("/button/1", p.handleInteractiveAction)

	dialogRouter := router.PathPrefix("/dialog").Subrouter()
	dialogRouter.Use(p.withDelay)
	dialogRouter.HandleFunc("/1", p.handleDialog1)
	dialogRouter.HandleFunc("/2", p.handleDialog2)
	dialogRouter.HandleFunc("/3", p.handleDialog3)
	dialogRouter.HandleFunc("/error", p.handleDialogWithError)

	ephemeralRouter := router.PathPrefix("/ephemeral").Subrouter()
	ephemeralRouter.Use(p.withDelay)
	ephemeralRouter.HandleFunc("/update", p.handleEphemeralUpdate)
	ephemeralRouter.HandleFunc("/delete", p.handleEphemeralDelete)

	p.router = router
}

func (p *Plugin) withDelay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		delay := p.getConfiguration().IntegrationRequestDelay
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		}

		next.ServeHTTP(w, r)
	})
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
		p.API.LogError("Failed to write status", "err", err.Error())
	}
}

func (p *Plugin) handleHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World!")); err != nil {
		p.API.LogError("Failed to write hello world", "err", err.Error())
	}
}

// The Authorization header should be an empty string if the request is by an
// authenticated user.
func (p *Plugin) handleCheckAuthHeader(w http.ResponseWriter, r *http.Request) {
	isAuthenticatedUser := r.Header.Get("Mattermost-User-ID") != ""
	authHeader := r.Header.Get(model.HeaderAuth)

	responseMessage := ""

	if isAuthenticatedUser {
		responseMessage += "You are an authenticated user. The Authorization header should be an empty string.\n"
	}

	responseMessage += fmt.Sprintf("Authorization header: %s", authHeader)

	if _, err := w.Write([]byte(responseMessage)); err != nil {
		p.API.LogError("Failed to write checkAuthHeader message", "err", err.Error())
	}
}

func (p *Plugin) handleOutgoingWebhook(w http.ResponseWriter, r *http.Request) {
	var request model.OutgoingWebhookPayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode OutgoingWebhookPayload", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	s, err := PrettyJSON(request)
	if err != nil {
		p.API.LogError("Failed to Marshal payload back to JSON", "err", err.Error())
		return
	}

	text := "```\n" + s + "\n```"
	resp := model.OutgoingWebhookResponse{
		Text: &text,
	}

	p.writeJSON(w, resp)
}

func (p *Plugin) handleDialog1(w http.ResponseWriter, r *http.Request) {
	var request model.SubmitDialogRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode SubmitDialogRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !request.Cancelled {
		number, ok := request.Submission[dialogElementNameNumber].(float64)
		if !ok {
			p.API.LogError("Request is missing field", "field", dialogElementNameNumber)
			w.WriteHeader(http.StatusOK)
			return
		}

		if number != 42 {
			response := &model.SubmitDialogResponse{
				Errors: map[string]string{
					dialogElementNameNumber: "This must be 42",
				},
			}
			p.writeJSON(w, response)
			return
		}
	}

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("Failed to get user for dialog", "err", appErr.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	msg := "@%v submitted an Interative Dialog"
	if request.Cancelled {
		msg = "@%v canceled an Interative Dialog"
	}

	rootPost, appErr := p.API.CreatePost(&model.Post{
		UserId:    p.botID,
		ChannelId: request.ChannelId,
		Message:   fmt.Sprintf(msg, user.Username),
	})
	if appErr != nil {
		p.API.LogError("Failed to post handleDialog1 message", "err", appErr.Error())
		return
	}

	if !request.Cancelled {
		// Don't post the email address publicly
		request.Submission[dialogElementNameEmail] = "xxxxxxxxxxx"

		if _, appErr = p.API.CreatePost(&model.Post{
			UserId:    p.botID,
			ChannelId: request.ChannelId,
			RootId:    rootPost.Id,
			Message:   "Data:",
			Type:      "custom_demo_plugin",
			Props:     request.Submission,
		}); appErr != nil {
			p.API.LogError("Failed to post handleDialog1 message", "err", appErr.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDialog2(w http.ResponseWriter, r *http.Request) {
	var request model.SubmitDialogRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode SubmitDialogRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("Failed to get user for dialog", "err", appErr.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	suffix := ""
	if request.State == dialogStateRelativeCallbackURL {
		suffix = "from relative callback URL"
	}

	if _, appErr = p.API.CreatePost(&model.Post{
		UserId:    p.botID,
		ChannelId: request.ChannelId,
		Message:   fmt.Sprintf("@%v confirmed an Interactive Dialog %v", user.Username, suffix),
	}); appErr != nil {
		p.API.LogError("Failed to post handleDialog2 message", "err", appErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDialog3(w http.ResponseWriter, r *http.Request) {
	var request model.SubmitDialogRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode SubmitDialogRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var message string
	if request.Cancelled {
		message = "Dialog cancelled"
	} else {
		submission := request.Submission

		// Generic approach - format submission data as structured lines
		message = "Dialog Submitted:"
		for key, value := range submission {
			valueStr := fmt.Sprintf("%v", value)
			message += fmt.Sprintf("\n- %s: %s", key, valueStr)
		}
	}

	// Post the message to the channel
	if _, appErr := p.API.CreatePost(&model.Post{
		UserId:    p.botID,
		ChannelId: request.ChannelId,
		Message:   message,
	}); appErr != nil {
		p.API.LogError("Failed to post handleDialog3 message", "err", appErr.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Plugin) handleDialogWithError(w http.ResponseWriter, r *http.Request) {
	// Always return an error
	response := &model.SubmitDialogResponse{
		Error: "some error",
	}
	p.writeJSON(w, response)
}

func (p *Plugin) handleEphemeralUpdate(w http.ResponseWriter, r *http.Request) {
	var request model.PostActionIntegrationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode PostActionIntegrationRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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
					Type: model.PostActionTypeButton,
					Name: fmt.Sprintf("Update %d", int(count)),
				}, {
					Integration: &model.PostActionIntegration{
						URL: fmt.Sprintf("%s/plugins/%s/ephemeral/delete", siteURL, manifest.Id),
					},
					Type: model.PostActionTypeButton,
					Name: "Delete",
				}},
			}},
		},
	}
	p.API.UpdateEphemeralPost(request.UserId, post)

	resp := &model.PostActionIntegrationResponse{}
	p.writeJSON(w, resp)
}

func (p *Plugin) handleEphemeralDelete(w http.ResponseWriter, r *http.Request) {
	var request model.PostActionIntegrationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode PostActionIntegrationRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	p.API.DeleteEphemeralPost(request.UserId, request.PostId)

	resp := &model.PostActionIntegrationResponse{}
	p.writeJSON(w, resp)
}

func (p *Plugin) handleInteractiveAction(w http.ResponseWriter, r *http.Request) {
	var request model.PostActionIntegrationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.API.LogError("Failed to decode PostActionIntegrationRequest", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, appErr := p.API.GetUser(request.UserId)
	if appErr != nil {
		p.API.LogError("Failed to get user for interactive action", "err", appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, postErr := p.API.GetPost(request.PostId)
	if postErr != nil {
		p.API.LogError("Failed to get post for interactive action", "err", postErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rootID := post.RootId
	if rootID == "" {
		rootID = post.Id
	}

	requestJSON, jsonErr := json.MarshalIndent(request, "", "  ")
	if jsonErr != nil {
		p.API.LogError("Failed to marshal json for interactive action", "err", jsonErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := "@%v clicked an interactive button.\n```json\n%v\n```"
	if _, appErr := p.API.CreatePost(&model.Post{
		UserId:    p.botID,
		ChannelId: request.ChannelId,
		RootId:    rootID,
		Message:   fmt.Sprintf(msg, user.Username, string(requestJSON)),
	}); appErr != nil {
		p.API.LogError("Failed to post handleInteractiveAction message", "err", appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &model.PostActionIntegrationResponse{}
	p.writeJSON(w, resp)
}

func (p *Plugin) writeJSON(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		p.API.LogError("Failed to write JSON response", "err", err.Error())
	}
}

func (p *Plugin) handleDynamicArgTest(w http.ResponseWriter, r *http.Request) {
	queryArgs := []string{"user_input", "parsed", "root_id", "parent_id", "user_id", "site_url", "request_id", "session_id", "ip_address", "accept_language", "user_agent"}
	query := r.URL.Query()

	channelID := query.Get("channel_id")
	teamID := query.Get("team_id")
	userID := query.Get("user_id")
	rootID := query.Get("root_id")

	channel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		http.Error(w, fmt.Sprintf("Error getting channels: %s", appErr.Error()), http.StatusInternalServerError)
		return
	}
	team, appErr := p.API.GetTeam(teamID)
	if appErr != nil {
		http.Error(w, fmt.Sprintf("Error getting team: %s", appErr.Error()), http.StatusInternalServerError)
		return
	}
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, fmt.Sprintf("Error getting user: %s", appErr.Error()), http.StatusInternalServerError)
		return
	}

	argMap := map[string]string{}
	for _, arg := range queryArgs {
		argMap[arg] = query.Get(arg)
	}
	argMapString := ""
	for k, v := range argMap {
		argMapString = fmt.Sprintf("%s  * %s:%s\n", argMapString, k, v)
	}
	result := fmt.Sprintf("dynamic argument was triggered by **%v** from team **%v** in the **%v** channel, with these arguments\n\n%v", user.GetFullName(), team.DisplayName, channel.DisplayName, argMapString)
	post := &model.Post{
		ChannelId: channelID,
		RootId:    rootID,
		UserId:    p.botID,
		Message:   result,
	}

	_, appErr = p.API.CreatePost(post)
	if appErr != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %s", appErr.Error()), http.StatusInternalServerError)
		return
	}

	suggestions := []model.AutocompleteListItem{{
		Item:     "suggestion 1",
		HelpText: "help text 1",
		Hint:     "(hint)",
	}, {
		Item:     "suggestion 2",
		HelpText: "help text 2",
		Hint:     "(hint)",
	}}

	jsonBytes, err := json.Marshal(suggestions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting dynamic args: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(jsonBytes); err != nil {
		http.Error(w, fmt.Sprintf("Error getting dynamic args: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
