// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"database/sql"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

type Params struct {
	DBType                  string
	ConnectionString        string
	TablePrefix             string
	DB                      *sql.DB
	PluginAPI               plugin.API
	SkipMigrations          bool
	Driver                  plugin.Driver
	MigrationTimeoutSeconds int
}

func (p Params) IsValid() error {
	if p.ConnectionString == "" {
		return errors.New("SQLStore Params.IsValid: invalid param: ConnectionString cannot be empty")
	}

	return nil
}
