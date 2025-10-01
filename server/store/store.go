package store

import "github.com/itstar-tech/mattermost-plugin-demo/server/models"

type Store interface {
	GetSession(sessionId string) (*models.Session, error)
	CreateSession(userId string) (*models.Session, error)
	CloseSession(sessionId string) (*models.Session, error)
	GetSessionByUserId(userId string) (*models.Session, error)
	GetSessionsUnclosed() ([]models.Session, error)
	GetWhatsappChannels() ([]models.WhatsappChannel, error)
	CreateWhatsappChannel(channelId string) (*models.WhatsappChannel, error)
	Close() error
}
