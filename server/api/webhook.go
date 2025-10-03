package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/itstar-tech/mattermost-plugin-demo/server/model"
	MattermostModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/pkg/errors"
)

type WebhookEntry struct {
	PhoneNumber    string        `json:"phone_number"`
	PhoneNumberID  string        `json:"phone_number_id"`
	MessageDetails MessageParent `json:"message"`
}

type WebhookPayload struct {
	MessageDetails MessageParent `json:"message"`
}

type MessageParent struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
	Field            string    `json:"field"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WAID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      Text   `json:"text"`
	Image     Image  `json:"image"`
	Type      string `json:"type"`
}

type Text struct {
	Body string `json:"body"`
}

type Image struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	ID       string `json:"id"`
}

func (api *Handlers) handleWhatsAppWebhook(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload
	var botUserID = api.app.GetBotId()
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	phoneNumber := payload.MessageDetails.Metadata.DisplayPhoneNumber
	phoneNumberID := payload.MessageDetails.Metadata.PhoneNumberID

	webhook := WebhookEntry{
		PhoneNumber:    phoneNumber,
		PhoneNumberID:  phoneNumberID,
		MessageDetails: payload.MessageDetails,
	}

	err := api.processWhatsappWebhook(webhook, botUserID)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ReturnStatusOK(w)
}

func (m *WebhookEntry) makeName(phoneNumberId string, phoneNumber string) string {
	return "whatsapp_" + phoneNumberId + "_" + phoneNumber
}

func (api *Handlers) processWhatsappWebhook(webhook WebhookEntry, botUserID string) error {
	userName := webhook.MessageDetails.Contacts[0].Profile.Name
	phoneNumber := webhook.PhoneNumber

	channelName := webhook.makeName(webhook.PhoneNumberID, phoneNumber)

	teams, err := api.pluginAPI.GetTeams()
	if err != nil {
		return errors.Wrapf(err, "failed to get teams for webhook from %s", phoneNumber)
	}
	if len(teams) == 0 {
		return errors.New("no teams found on Mattermost server")
	}
	teamID := teams[0].Id

	channel, err := api.pluginAPI.GetChannelByName(teamID, channelName, false)

	if channel == nil || err != nil {
		newChannel := &MattermostModel.Channel{
			Name:        channelName,
			DisplayName: userName,
			Header:      "WhatsApp Chat with " + userName,
			TeamId:      teamID,
			Type:        MattermostModel.ChannelTypePrivate,
			Props: map[string]interface{}{
				"phone_number":    phoneNumber,
				"phone_number_id": webhook.PhoneNumberID,
				"channel_type":    "whatsapp",
			},
		}

		channel, err = api.pluginAPI.CreateChannel(newChannel)
		if err != nil {
			return errors.Wrapf(err, "failed to create channel '%s' for user %s", channelName, userName)
		}
	}

	for _, msg := range webhook.MessageDetails.Messages {
		err := api.processPost(&msg, botUserID, channel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (api *Handlers) processPost(message *Message, botUserId string, channel *MattermostModel.Channel) error {
	var newPost *MattermostModel.Post
	if message.Type == "image" {
		newPost = &MattermostModel.Post{
			UserId:    botUserId,
			ChannelId: channel.Id,
			Message:   message.Image.Caption,
		}
	} else if message.Type == "text" {
		newPost = &MattermostModel.Post{
			UserId:    botUserId,
			ChannelId: channel.Id,
			Message:   message.Text.Body,
		}
	}
	_, err := api.pluginAPI.CreatePost(newPost)
	if err != nil {
		return errors.Wrapf(err, "failed to create post for message (text: '%s') in channel %s",
			message.Text.Body, channel.Name)
	}
	return nil
}
