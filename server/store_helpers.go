package main

import (
	"fmt"
	"time"

	"github.com/itstar-tech/mattermost-plugin-demo/server/models"
	"github.com/mattermost/mattermost/server/public/model"
)

// Store helper methods for the plugin

// SaveMessageToStore saves a message to the database store
func (p *Plugin) SaveMessageToStore(post *model.Post) error {
	if p.store == nil {
		return fmt.Errorf("store not initialized")
	}

	message := &models.Message{
		ID:        post.Id,
		ChannelID: post.ChannelId,
		UserID:    post.UserId,
		Content:   post.Message,
		CreatedAt: time.Unix(post.CreateAt/1000, 0),
		UpdatedAt: time.Unix(post.UpdateAt/1000, 0),
	}

	return p.store.CreateMessage(message)
}

// GetMessageFromStore retrieves a message from the database store
func (p *Plugin) GetMessageFromStore(messageId string) (*models.Message, error) {
	if p.store == nil {
		return nil, fmt.Errorf("store not initialized")
	}

	return p.store.GetMessage(messageId)
}

// GetChannelMessagesFromStore retrieves messages for a channel from the database store
func (p *Plugin) GetChannelMessagesFromStore(channelId string, limit, offset int) ([]*models.Message, error) {
	if p.store == nil {
		return nil, fmt.Errorf("store not initialized")
	}

	return p.store.GetMessagesByChannel(channelId, limit, offset)
}

// GetUserMessagesFromStore retrieves messages for a user from the database store
func (p *Plugin) GetUserMessagesFromStore(userId string, limit, offset int) ([]*models.Message, error) {
	if p.store == nil {
		return nil, fmt.Errorf("store not initialized")
	}

	return p.store.GetMessagesByUser(userId, limit, offset)
}

// UpdateMessageInStore updates a message in the database store
func (p *Plugin) UpdateMessageInStore(message *models.Message) error {
	if p.store == nil {
		return fmt.Errorf("store not initialized")
	}

	message.UpdatedAt = time.Now()
	return p.store.UpdateMessage(message)
}

// DeleteMessageFromStore deletes a message from the database store
func (p *Plugin) DeleteMessageFromStore(messageId string) error {
	if p.store == nil {
		return fmt.Errorf("store not initialized")
	}

	return p.store.DeleteMessage(messageId)
}
