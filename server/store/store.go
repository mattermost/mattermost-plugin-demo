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
}
