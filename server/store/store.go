package store

import (
	"text/template"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

type Store interface {
	Shutdown() error
	Migrate(migrationTimeoutSeconds int) error
	GetTemplateHelperFuncs() template.FuncMap
	GetSchemaName() (string, error)
	GetSessions() ([]*model.Session, error)
	CreateSession(session *model.Session) error
	GetSessionByID(id string) (*model.Session, error)
	UpdateSession(session *model.Session) error
	DeleteSession(id string) error
}
