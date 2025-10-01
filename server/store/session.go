package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
	_ "github.com/lib/pq"
	mm_model "github.com/mattermost/mattermost/server/public/model"
)

func (S SQLStore) GetSession(sessionId string) (*model.WhatsappSession, error) {
	query := "SELECT id, user_id, created_at, closed_at FROM whatsapp_plugin_session WHERE id = $1"
	var session model.WhatsappSession
	err := S.db.QueryRow(query, sessionId).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ClosedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session with id %s not found", sessionId)
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return &session, nil
}

func (S SQLStore) CreateSession(userId string) (*model.WhatsappSession, error) {
	newId := mm_model.NewId()
	session := &model.WhatsappSession{
		ID:        newId,
		UserID:    userId,
		CreatedAt: time.Now(),
	}
	query := "INSERT INTO whatsapp_plugin_session (id, user_id, created_at) VALUES ($1, $2, $3)"
	_, err := S.db.Exec(query, session.ID, session.UserID, session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return session, nil
}

func (S SQLStore) CloseSession(sessionId string) (*model.WhatsappSession, error) {
	session, err := S.GetSession(sessionId)
	if err != nil {
		return nil, err
	}
	query := "UPDATE whatsapp_plugin_session SET closed_at = $1 WHERE id = $2"
	_, err = S.db.Exec(query, time.Now(), sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to close session: %w", err)
	}
	session.ClosedAt = time.Now()
	return session, nil
}

func (S SQLStore) GetSessionByUserId(userId string) (*model.WhatsappSession, error) {
	query := "SELECT id, user_id, created_at, closed_at FROM whatsapp_plugin_session WHERE user_id = $1 AND closed_at IS NULL ORDER BY created_at DESC LIMIT 1"
	var session model.WhatsappSession
	err := S.db.QueryRow(query, userId).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ClosedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session with user id %s not found", userId)
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return &session, nil
}

func (S SQLStore) GetSessionsUnclosed() ([]model.WhatsappSession, error) {
	query := "SELECT id, user_id, created_at, closed_at FROM whatsapp_plugin_session WHERE closed_at IS NULL"
	var sessions []model.WhatsappSession
	rows, err := S.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var session model.WhatsappSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.CreatedAt,
			&session.ClosedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
