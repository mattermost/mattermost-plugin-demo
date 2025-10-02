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
		"closed_at",
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
		var closedAt sql.NullInt64

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CreateAt,
			&closedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "SessionsFromRows failed to scan session row")
		}

		if closedAt.Valid {
			session.ClosedAt = &closedAt.Int64
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (s *SQLStore) CreateSession(session *model.Session) error {
	session.SetDefaults()

	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "CreateSession: invalid session")
	}

	_, err := s.getQueryBuilder().
		Insert(s.tablePrefix+"session").
		Columns(s.sessionColumns()...).
		Values(
			session.ID,
			session.UserID,
			session.CreateAt,
			session.ClosedAt,
		).
		Exec()

	if err != nil {
		return errors.Wrap(err, "CreateSession: failed to insert session into database")
	}

	return nil
}

func (s *SQLStore) GetSessionByID(id string) (*model.Session, error) {
	row := s.getQueryBuilder().
		Select(s.sessionColumns()...).
		From(s.tablePrefix+"session").
		Where("id = ?", id).
		QueryRow()

	var session model.Session
	var closedAt sql.NullInt64

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CreateAt,
		&closedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "GetSessionByID: failed to scan session")
	}

	if closedAt.Valid {
		session.ClosedAt = &closedAt.Int64
	}

	return &session, nil
}

func (s *SQLStore) UpdateSession(session *model.Session) error {
	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "UpdateSession: invalid session")
	}

	_, err := s.getQueryBuilder().
		Update(s.tablePrefix+"session").
		Set("user_id", session.UserID).
		Set("create_at", session.CreateAt).
		Set("closed_at", session.ClosedAt).
		Where("id = ?", session.ID).
		Exec()

	if err != nil {
		return errors.Wrap(err, "UpdateSession: failed to update session in database")
	}

	return nil
}

func (s *SQLStore) DeleteSession(id string) error {
	result, err := s.getQueryBuilder().
		Delete(s.tablePrefix+"session").
		Where("id = ?", id).
		Exec()

	if err != nil {
		return errors.Wrap(err, "DeleteSession: failed to delete session from database")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "DeleteSession: failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("DeleteSession: session not found")
	}

	return nil
}

func (s *SQLStore) GetActiveSessionByUserID(userID string) (*model.Session, error) {
	row := s.getQueryBuilder().
		Select(s.sessionColumns()...).
		From(s.tablePrefix+"session").
		Where("user_id = ? AND closed_at IS NULL", userID).
		OrderBy("create_at DESC").
		Limit(1).
		QueryRow()

	var session model.Session
	var closedAt sql.NullInt64

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CreateAt,
		&closedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "GetActiveSessionByUserID: failed to scan session")
	}

	if closedAt.Valid {
		session.ClosedAt = &closedAt.Int64
	}

	return &session, nil
}
