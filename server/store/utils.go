// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"strings"

	"github.com/mattermost/squirrel"
	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

var (
	ErrUnsupportedDatabaseType = errors.New("database type is unsupported")
)

// replaceVars replaces instances of variable placeholders with the
// values provided via a map.  Variable placeholders are of the form
// `[[var_name]]`.
func replaceVars(s string, vars map[string]string) string {
	for key, val := range vars {
		placeholder := "[[" + key + "]]"
		val = strings.ReplaceAll(val, "'", "\\'")
		s = strings.ReplaceAll(s, placeholder, val)
	}
	return s
}

func (s *SQLStore) GetSchemaName() (string, error) {
	var query squirrel.SelectBuilder

	switch s.dbType {
	case model.DBTypeMySQL:
		query = s.getQueryBuilder().Select("DATABASE()")
	case model.DBTypePostgres:
		query = s.getQueryBuilder().Select("current_schema()")
	default:
		return "", ErrUnsupportedDatabaseType
	}

	scanner := query.QueryRow()

	var result string
	err := scanner.Scan(&result)
	if err != nil && !model.IsErrNotFound(err) {
		return "", err
	}

	return result, nil
}
