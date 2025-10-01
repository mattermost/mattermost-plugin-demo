// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"database/sql"
	"net/http"

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

const (
	DBTypeMySQL    = "mysql"
	DBTypePostgres = "postgres"
)

func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}

	// check if this is a sql.ErrNotFound
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	// check if this is a plugin API error
	if errors.Is(err, pluginapi.ErrNotFound) {
		return true
	}

	// check if this is a Mattermost AppError with a Not Found status
	var appErr *mmModel.AppError
	if errors.As(err, &appErr) {
		if appErr.StatusCode == http.StatusNotFound {
			return true
		}
	}

	return false
}
