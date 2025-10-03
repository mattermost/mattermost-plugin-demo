package store

import (
	"text/template"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	MattermostModel "github.com/mattermost/mattermost/server/public/model"
)

type Store interface {
	Shutdown() error
	Migrate(migrationTimeoutSeconds int) error
	GetTemplateHelperFuncs() template.FuncMap
	GetSchemaName() (string, error)
	GetActiveUsers() ([]*MattermostModel.User, error)
	GetSessions() ([]*model.Session, error)
	CreateSession(session *model.Session) error
	GetSessionByUserId(id string) (*model.Session, error)
	UpdateSession(session *model.Session) error
	CloseSessionsFromUserId(id string) error
	GetChannels() ([]*model.Channel, error)
	CreateChannel(channel *model.Channel) error
	GetChannelByID(id string) (*model.Channel, error)
	UpdateChannel(channel *model.Channel) error
	DeleteChannel(id string) error
}
