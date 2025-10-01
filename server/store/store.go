package store

import "github.com/itstar-tech/mattermost-plugin-demo/server/model"

type Store interface {
	GetSession(sessionId string) (*model.WhatsappSession, error)
	CreateSession(userId string) (*model.WhatsappSession, error)
	CloseSession(sessionId string) (*model.WhatsappSession, error)
	GetSessionByUserId(userId string) (*model.WhatsappSession, error)
	GetSessionsUnclosed() ([]model.WhatsappSession, error)
	GetWhatsappChannels() ([]model.WhatsappChannel, error)
	CreateWhatsappChannel(channelId string) (*model.WhatsappChannel, error)
	Close() error
}
