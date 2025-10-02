package store

import (
	"text/template"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

type Store interface {
	Shutdown() error
	Migrate(migrationTimeoutSeconds int) error
	GetSession(sessionId string) (*model.WhatsappSession, error)
	CreateSession(userId string) (*model.WhatsappSession, error)
	CloseSession(sessionId string) (*model.WhatsappSession, error)
	GetSessionByUserId(userId string) (*model.WhatsappSession, error)
	GetSessionsUnclosed() ([]model.WhatsappSession, error)
	GetWhatsappChannels() ([]model.WhatsappChannel, error)
	CreateWhatsappChannel(channelId string) (*model.WhatsappChannel, error)

	GetTemplateHelperFuncs() template.FuncMap
	GetSchemaName() (string, error)
}
