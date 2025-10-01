// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"fmt"
	"strings"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (s *SQLStore) genAddColumnIfNeeded(tableName, columnName, dataType, constraint string) (string, error) {
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normalizedTableName := s.normalizeTablename(tableName)

	switch s.dbType {
	case model.DBTypeMySQL:
		vars := map[string]string{
			"schema":          s.schemaName,
			"table_name":      tableName,
			"norm_table_name": normalizedTableName,
			"column_name":     columnName,
			"data_type":       dataType,
			"constraint":      constraint,
		}

		x := replaceVars(`
			SET @stmt = (SELECT IF(
				(
				  SELECT COUNT(column_name) FROM INFORMATION_SCHEMA.COLUMNS
				  WHERE table_name = '[[table_name]]'
				  AND table_schema = '[[schema]]'
				  AND column_name = '[[column_name]]'
				) > 0,
				'SELECT 1;',
				'ALTER TABLE [[norm_table_name]] ADD COLUMN [[column_name]] [[data_type]] [[constraint]];'
			));
			PREPARE addColumnIfNeeded FROM @stmt;
			EXECUTE addColumnIfNeeded;
			DEALLOCATE PREPARE addColumnIfNeeded;
		`, vars)

		s.pluginAPI.LogInfo("************************************************************")
		s.pluginAPI.LogInfo(x)
		s.pluginAPI.LogInfo("************************************************************")

		return x, nil

	case model.DBTypePostgres:
		return fmt.Sprintf("\nALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s %s;\n", normalizedTableName, columnName, dataType, constraint), nil

	default:
		return "", ErrUnsupportedDatabaseType
	}
}

func addPrefixIfNeeded(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}

func (s *SQLStore) normalizeTablename(tableName string) string {
	if s.schemaName != "" && !strings.HasPrefix(tableName, s.schemaName+".") {
		schemaName := s.schemaName
		if s.dbType == model.DBTypeMySQL {
			schemaName = "`" + schemaName + "`"
		}
		tableName = schemaName + "." + tableName
	}
	return tableName
}

func (s *SQLStore) genDropColumnIfNeeded(tableName, columnName string) (string, error) {
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normTableName := s.normalizeTablename(tableName)

	switch s.dbType {
	case model.DBTypeMySQL:
		vars := map[string]string{
			"schema":          s.schemaName,
			"table_name":      tableName,
			"norm_table_name": normTableName,
			"column_name":     columnName,
		}
		return replaceVars(`
			SET @stmt = (SELECT IF(
				(
				  SELECT COUNT(column_name) FROM INFORMATION_SCHEMA.COLUMNS
				  WHERE table_name = '[[table_name]]'
				  AND table_schema = '[[schema]]'
				  AND column_name = '[[column_name]]'
				) > 0,
				'ALTER TABLE [[norm_table_name]] DROP COLUMN [[column_name]];',
				'SELECT 1;'
			));
			PREPARE dropColumnIfNeeded FROM @stmt;
			EXECUTE dropColumnIfNeeded;
			DEALLOCATE PREPARE dropColumnIfNeeded;
		`, vars), nil
	case model.DBTypePostgres:
		return fmt.Sprintf("\nALTER TABLE %s DROP COLUMN IF EXISTS %s;\n", normTableName, columnName), nil
	default:
		return "", ErrUnsupportedDatabaseType
	}
}

func (s *SQLStore) genAddConstraintIfNeeded(tableName, constraintName, constraintType, constraintDefinition string) (string, error) {
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normTableName := s.normalizeTablename(tableName)

	var query string

	vars := map[string]string{
		"schema":                s.schemaName,
		"constraint_name":       constraintName,
		"constraint_type":       constraintType,
		"table_name":            tableName,
		"constraint_definition": constraintDefinition,
		"norm_table_name":       normTableName,
	}

	switch s.dbType {
	case model.DBTypeMySQL:
		query = replaceVars(`
			SET @stmt = (SELECT IF(
				(
				SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
				WHERE constraint_schema = '[[schema]]'
				AND constraint_name = '[[constraint_name]]'
				AND constraint_type = '[[constraint_type]]'
				AND table_name = '[[table_name]]'
				) > 0,
				'SELECT 1;',
				'ALTER TABLE [[norm_table_name]] ADD CONSTRAINT [[constraint_name]] [[constraint_definition]];'
			));
			PREPARE addConstraintIfNeeded FROM @stmt;
			EXECUTE addConstraintIfNeeded;
			DEALLOCATE PREPARE addConstraintIfNeeded;
		`, vars)
	case model.DBTypePostgres:
		query = replaceVars(`
		DO
		$$
		BEGIN
		IF NOT EXISTS (
			SELECT * FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
				WHERE constraint_schema = '[[schema]]'
				AND constraint_name = '[[constraint_name]]'
				AND constraint_type = '[[constraint_type]]'
				AND table_name = '[[table_name]]'
		) THEN
			ALTER TABLE [[norm_table_name]] ADD CONSTRAINT [[constraint_name]] [[constraint_definition]];
		END IF;
		END;
		$$
		LANGUAGE plpgsql;
		`, vars)
	}

	return query, nil
}

func (s *SQLStore) genCreateIndexIfNeeded(tableName, columns string) (string, error) {
	indexName := getIndexName(tableName, columns)
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normTableName := s.normalizeTablename(tableName)

	switch s.dbType {
	case model.DBTypeMySQL:
		vars := map[string]string{
			"schema":          s.schemaName,
			"table_name":      tableName,
			"norm_table_name": normTableName,
			"index_name":      indexName,
			"columns":         columns,
		}
		return replaceVars(`
			SET @stmt = (SELECT IF(
				(
				  SELECT COUNT(index_name) FROM INFORMATION_SCHEMA.STATISTICS
				  WHERE table_name = '[[table_name]]'
				  AND table_schema = '[[schema]]'
				  AND index_name = '[[index_name]]'
				) > 0,
				'SELECT 1;',
				'CREATE INDEX [[index_name]] ON [[norm_table_name]] ([[columns]]);'
			));
			PREPARE createIndexIfNeeded FROM @stmt;
			EXECUTE createIndexIfNeeded;
			DEALLOCATE PREPARE createIndexIfNeeded;
		`, vars), nil
	case model.DBTypePostgres:
		return fmt.Sprintf("\nCREATE INDEX IF NOT EXISTS %s ON %s (%s);\n", indexName, normTableName, columns), nil
	default:
		return "", ErrUnsupportedDatabaseType
	}
}

func (s *SQLStore) genDropIndexIfNeeded(tableName, columns string) (string, error) {
	indexName := getIndexName(tableName, columns)
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normTableName := s.normalizeTablename(tableName)

	switch s.dbType {
	case model.DBTypeMySQL:
		vars := map[string]string{
			"schema":          s.schemaName,
			"table_name":      tableName,
			"norm_table_name": normTableName,
			"index_name":      indexName,
			"columns":         columns,
		}
		return replaceVars(`
			SET @stmt = (SELECT IF(
				(
				  SELECT COUNT(index_name) FROM INFORMATION_SCHEMA.STATISTICS
				  WHERE table_name = '[[table_name]]'
				  AND table_schema = '[[schema]]'
				  AND index_name = '[[index_name]]'
				) > 0,
				'DROP INDEX [[index_name]] ON [[norm_table_name]];',
				'SELECT 1;'
			));
			PREPARE dropIndexIfNeeded FROM @stmt;
			EXECUTE dropIndexIfNeeded;
			DEALLOCATE PREPARE dropIndexIfNeeded;
		`, vars), nil
	case model.DBTypePostgres:
		return fmt.Sprintf("\nDROP INDEX IF EXISTS %s;\n", indexName), nil
	default:
		return "", ErrUnsupportedDatabaseType
	}
}

func getIndexName(tableName string, columns string) string {
	var sb strings.Builder

	_, _ = sb.WriteString("idx_")
	_, _ = sb.WriteString(tableName)

	// allow developers to separate column names with spaces and/or commas
	columns = strings.ReplaceAll(columns, ",", " ")
	cols := strings.Split(columns, " ")

	for _, s := range cols {
		sub := strings.TrimSpace(s)
		if sub == "" {
			continue
		}

		_, _ = sb.WriteString("_")
		_, _ = sb.WriteString(s)
	}
	return sb.String()
}

func (s *SQLStore) genRenameColumnIfNeeded(tableName, oldColumnName, newColumnName, dataType string) (string, error) {
	tableName = addPrefixIfNeeded(tableName, s.tablePrefix)
	normTableName := s.normalizeTablename(tableName)

	vars := map[string]string{
		"schema":          s.schemaName,
		"table_name":      tableName,
		"norm_table_name": normTableName,
		"old_column_name": oldColumnName,
		"new_column_name": newColumnName,
		"data_type":       dataType,
	}

	switch s.dbType {
	case model.DBTypeMySQL:
		return replaceVars(`
			SET @stmt = (SELECT IF(
				(
				SELECT COUNT(column_name) FROM INFORMATION_SCHEMA.COLUMNS
				WHERE table_name = '[[table_name]]'
				AND table_schema = '[[schema]]'
				AND column_name = '[[new_column_name]]'
				) > 0,
				'SELECT 1;',
				'ALTER TABLE [[norm_table_name]] CHANGE [[old_column_name]] [[new_column_name]] [[data_type]];'
			));
			PREPARE renameColumnIfNeeded FROM @stmt;
			EXECUTE renameColumnIfNeeded;
			DEALLOCATE PREPARE renameColumnIfNeeded;
		`, vars), nil
	case model.DBTypePostgres:
		return replaceVars(`
			do $$
			begin
				if (SELECT COUNT(table_name) FROM INFORMATION_SCHEMA.COLUMNS
							WHERE table_name = '[[table_name]]'
							AND table_schema = '[[schema]]'
							AND column_name = '[[new_column_name]]'
				) = 0 then
					ALTER TABLE [[norm_table_name]] RENAME COLUMN [[old_column_name]] TO [[new_column_name]];
				end if;
			end$$;
		`, vars), nil
	default:
		return "", ErrUnsupportedDatabaseType
	}
}
