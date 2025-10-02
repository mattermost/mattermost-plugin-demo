// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func (s *SQLStore) sessionColumns() []string {
	return []string{
		"id",
		"user_id",
		"create_at",
	}
}

func (s *SQLStore) GetSessions() ([]*model.Session, error) {

	rows, err := s.getQueryBuilder().
		Select(s.sessionColumns()...).
		From(s.tablePrefix + "session").
		Query()

	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetInProgressSurvey failed to fetch survey by status from database")
	}

	sessions, err := s.SessionsFromRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "GetSurveysByStatus: failed to map survey rows to surveys")
	}

	return sessions, nil
}

func (s *SQLStore) SessionsFromRows(rows *sql.Rows) ([]*model.Session, error) {
	sessions := []*model.Session{}

	for rows.Next() {
		var session model.Session

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CreateAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "SessionsFromRows failed to scan session row")
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}
