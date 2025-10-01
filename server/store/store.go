package store

import "github.com/itstar-tech/mattermost-plugin-demo/server/models"

type Store interface {
	GetMessage(messageId string) (*models.Message, error)
	CreateMessage(message *models.Message) error
	UpdateMessage(message *models.Message) error
	DeleteMessage(messageId string) error
	GetMessagesByChannel(channelId string, limit, offset int) ([]*models.Message, error)
	GetMessagesByUser(userId string, limit, offset int) ([]*models.Message, error)
}
