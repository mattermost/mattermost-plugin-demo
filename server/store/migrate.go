// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path"
	"text/template"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
	sqlUtils "github.com/mattermost/mattermost/server/public/utils/sql"
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	"github.com/mattermost/morph/drivers/mysql"
	"github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/sources/embedded"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

const (
	migrationDBPPingRetries = 5

	migrationAssetsDir = "migrations"

	migrationLockKey = "plugin-demo-lock-key"
)

//go:embed migrations/*.sql
var Assets embed.FS

func (s *SQLStore) Migrate(migrationTimeoutSeconds int) error {
	var driver drivers.Driver
	var err error

	s.pluginAPI.LogDebug("Obtaining migration connection")
	db, err := s.getMigrationConnection()
	if err != nil {
		return err
	}

	defer func() {
		s.pluginAPI.LogDebug("Closing migration connection")
		if dbErr := db.Close(); dbErr != nil {
			s.pluginAPI.LogError("Failed to close migration connection", "error", dbErr.Error())
		}
	}()

	switch s.dbType {
	case model.DBTypePostgres:
		driver, err = postgres.WithInstance(db)
	case model.DBTypeMySQL:
		driver, err = mysql.WithInstance(db)
	default:
		err = fmt.Errorf("unknown DB type encountered, dbtype: %s", s.dbType)
		s.pluginAPI.LogError("Unknown DB type encountered", "error", err.Error())
		return err
	}

	if err != nil {
		s.pluginAPI.LogError("Failed to create database driver instance", "error", err.Error())
		return err
	}

	migrationAssets, err := s.generateMigrationAssets()
	if err != nil {
		return err
	}

	src, err := embedded.WithInstance(migrationAssets)
	if err != nil {
		s.pluginAPI.LogError("Failed to generate migration sources from migration assets", "error", err.Error())
		return err
	}

	engineOptions := []morph.EngineOption{
		morph.WithLock(migrationLockKey),
		morph.SetMigrationTableName(fmt.Sprintf("%sschema_migrations", s.tablePrefix)),
		morph.SetStatementTimeoutInSeconds(migrationTimeoutSeconds),
	}

	s.pluginAPI.LogDebug("Creating migration engine")

	engine, err := morph.New(context.Background(), driver, src, engineOptions...)
	if err != nil {
		s.pluginAPI.LogError("Failed to create database migration engine", "error", err.Error())
		return err
	}

	defer func() {
		s.pluginAPI.LogDebug("Closing database migration engine")
		if err := engine.Close(); err != nil {
			s.pluginAPI.LogError("Failed to clone database emigration engine", "error", err.Error())
		}
	}()

	return s.runMigrations(engine, driver)
}

func (s *SQLStore) getMigrationConnection() (*sql.DB, error) {
	connectionString := s.connectionString

	if s.dbType == model.DBTypeMySQL {
		var err error
		connectionString, err = sqlUtils.ResetReadTimeout(connectionString)
		if err != nil {
			s.pluginAPI.LogError("failed to reset read timeout on MySQL connection string", "error", err.Error())
			return nil, err
		}

		connectionString, err = sqlUtils.AppendMultipleStatementsFlag(connectionString)
		if err != nil {
			s.pluginAPI.LogError("failed to append multi statement flag on MySQL connection string", "error", err.Error())
			return nil, err
		}
	}

	sqlSettings := s.pluginAPI.GetUnsanitizedConfig().SqlSettings

	logger, err := mlog.NewLogger()
	if err != nil {
		s.pluginAPI.LogError("failed to crete new mLog logger instance", "error", err.Error())
		return nil, err
	}

	return sqlUtils.SetupConnection(logger, "master", connectionString, &sqlSettings, migrationDBPPingRetries)
}

func (s *SQLStore) generateMigrationAssets() (*embedded.AssetSource, error) {
	assetList, err := Assets.ReadDir(migrationAssetsDir)
	if err != nil {
		s.pluginAPI.LogError("Failed to read migration asset dir", "error", err.Error())
		return nil, err
	}

	assetNamesForDriver := make([]string, len(assetList))
	for i, asset := range assetList {
		assetNamesForDriver[i] = asset.Name()
	}

	templateParams := map[string]interface{}{
		"prefix":   s.tablePrefix,
		"postgres": s.dbType == model.DBTypePostgres,
		"mysql":    s.dbType == model.DBTypeMySQL,
	}

	migrationAssets := &embedded.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			asset, err := Assets.ReadFile(path.Join(migrationAssetsDir, name))
			if err != nil {
				s.pluginAPI.LogError("Failed to read migration file", "fileName", name, "error", err.Error())
				return nil, err
			}

			tmpl, err := template.New("sql").Funcs(s.GetTemplateHelperFuncs()).Parse(string(asset))
			if err != nil {
				s.pluginAPI.LogError("Failed to parse migration template", "fileName", name, "error", err.Error())
				return nil, err
			}

			buffer := bytes.NewBufferString("")
			if err := tmpl.Execute(buffer, templateParams); err != nil {
				s.pluginAPI.LogError("Failed to execute migration template", "fileName", name, "error", err.Error())
				return nil, err
			}

			s.pluginAPI.LogDebug("Generated migration SQL", "migrationName", name, "sql", buffer.String())

			return buffer.Bytes(), nil
		},
	}

	return migrationAssets, nil
}

func (s *SQLStore) GetTemplateHelperFuncs() template.FuncMap {
	// these are all referenced from Focalboard.
	// See source for more such utility functions here -
	// https://github.com/mattermost/focalboard/blob/7a31925d8a7469a0568c795fc175237207e3d0c8/server/services/store/sqlstore/migrate.go#L306
	return template.FuncMap{
		"addColumnIfNeeded":     s.genAddColumnIfNeeded,
		"dropColumnIfNeeded":    s.genDropColumnIfNeeded,
		"addConstraintIfNeeded": s.genAddConstraintIfNeeded,
		"createIndexIfNeeded":   s.genCreateIndexIfNeeded,
		"dropIndexIfNeeded":     s.genDropIndexIfNeeded,
		"renameColumnIfNeeded":  s.genRenameColumnIfNeeded,
	}
}

func (s *SQLStore) runMigrations(engine *morph.Morph, driver drivers.Driver) error {
	appliedMigrations, err := driver.AppliedMigrations()
	if err != nil {
		s.pluginAPI.LogError("Failed to get currently applied migrations", "error", err.Error())
		return err
	}

	s.pluginAPI.LogDebug("Applying all remaining migrations...", "current_version", len(appliedMigrations))

	if err := engine.ApplyAll(); err != nil {
		s.pluginAPI.LogError("Failed to apply migrations", "current_version", len(appliedMigrations), "error", err.Error())
		return err
	}

	return nil
}
