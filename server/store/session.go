// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	MattermostModel "github.com/mattermost/mattermost/server/public/model"
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
		Where("closed_at IS NULL").
		Query()

	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetSessions failed to fetch sessions from database")
	}

	sessions, err := s.SessionsFromRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetSessions: failed to map session rows to sessions")
	}

	return sessions, nil
}

func (s *SQLStore) GetActiveUsers() ([]*MattermostModel.User, error) {

	rows, err := s.getQueryBuilder().
		Select("user_id").
		From(s.tablePrefix + "session").
		Where("closed_at IS NULL").
		Query()

	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetSessions failed to fetch sessions from database")
	}

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, errors.Wrap(err, "SQLStore.GetActiveUsers failed to scan user ID from row")
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetActiveUsers rows iteration failed")
	}
	users, err := s.pluginAPI.GetUsersByIds(userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "SQLStore.GetActiveUsers failed to fetch users by IDs from plugin API")
	}

	return users, nil
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

func (s *SQLStore) GetSessionByUserId(id string) (*model.Session, error) {
	row := s.getQueryBuilder().
		Select(s.sessionColumns()...).
		From(s.tablePrefix+"session").
		Where("user_id = ?", id).
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

func (s *SQLStore) CloseSessionsFromUserId(id string) error {
	_, err := s.getQueryBuilder().
		Update(s.tablePrefix+"session").
		Set("closed_at", MattermostModel.GetMillis()).
		Where("user_id = ? AND closed_at IS NULL", id).
		Exec()

	if err != nil {
		return errors.Wrap(err, "CloseSessionsFromUserId: failed to close sessions in database")
	}

	return nil
}
