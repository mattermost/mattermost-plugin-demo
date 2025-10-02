// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"database/sql"
	"encoding/json"
	"net/url"

	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/squirrel"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

const (
	TablePrefix = "demo_plugin_"
)

type SQLStore struct {
	db               *sql.DB
	dbType           string
	tablePrefix      string
	connectionString string
	pluginAPI        plugin.API
	isBinaryParams   bool
	skipMigrations   bool
	schemaName       string
	apiClient        *pluginapi.Client
}

func New(params Params) (*SQLStore, error) {
	if err := params.IsValid(); err != nil {
		return nil, err
	}

	params.PluginAPI.LogInfo("initializing SQLStore...")
	store := &SQLStore{
		db:               params.DB,
		dbType:           params.DBType,
		tablePrefix:      params.TablePrefix,
		connectionString: params.ConnectionString,
		pluginAPI:        params.PluginAPI,
		skipMigrations:   params.SkipMigrations,
		apiClient:        pluginapi.NewClient(params.PluginAPI, params.Driver),
	}

	var err error

	store.isBinaryParams, err = store.checkBinaryParams()
	if err != nil {
		return nil, err
	}

	store.schemaName, err = store.GetSchemaName()
	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.New failed to get database schema name")
	}

	if !store.skipMigrations {
		if migrationErr := store.Migrate(params.MigrationTimeoutSeconds); migrationErr != nil {
			params.PluginAPI.LogError(`Table creation / migration failed`, "error", migrationErr.Error())
			return nil, migrationErr
		}
	}

	return store, nil
}

func (s *SQLStore) checkBinaryParams() (bool, error) {
	if s.dbType != model.DBTypePostgres {
		return false, nil
	}

	parsedURL, err := url.Parse(s.connectionString)
	if err != nil {
		s.pluginAPI.LogError("failed to parse database connection string URL", "error", err.Error())
		return false, err
	}

	return parsedURL.Query().Get("binary_parameters") == "yes", nil
}

func (s *SQLStore) Shutdown() error {
	return s.db.Close()
}

func (s *SQLStore) getQueryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(s.getQueryPlaceholder()).RunWith(s.db)
}

func (s *SQLStore) getMasterQueryBuilder() (*squirrel.StatementBuilderType, error) {
	masterDB, err := s.apiClient.Store.GetMasterDB()
	if err != nil {
		return nil, errors.Wrap(err, "getMasterQueryBuilder: failed to get master DB from plugin API client")
	}
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(s.getQueryPlaceholder()).RunWith(masterDB)
	return &queryBuilder, nil
}

func (s *SQLStore) getQueryPlaceholder() squirrel.PlaceholderFormat {
	if s.dbType == model.DBTypePostgres {
		return squirrel.Dollar
	}
	return squirrel.Question
}

func (s *SQLStore) MarshalJSONB(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// in Postgres, when using binary param for connection and sending data for a JSONB column,
	// if the code sends a byte array (As we normally do), Postgres assumes it to be already encoded and expected it
	// to be in the right encoding. Unfortunately, the default JSON encoding is not the same as Postgres' JSONB encoding.
	// Postgres' JSONB encoding expects the data to start with the JSON version, which currently is always `1` (or, 0x01 in hex).
	// So, we check here if we're using binary params (which is only set if we're using Postgres), and append the JSON version number
	// before the data byte array.
	if s.isBinaryParams {
		b = append([]byte{0x01}, b...)
	}

	return b, nil
}
