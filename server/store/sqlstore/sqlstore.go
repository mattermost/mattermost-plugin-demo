package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/itstar-tech/mattermost-plugin-demo/server/models"
	"github.com/itstar-tech/mattermost-plugin-demo/server/store"
	_ "github.com/lib/pq"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type SQLStore struct {
	db     *sql.DB
	api    plugin.API
	driver string
}

func New(api plugin.API) (store.Store, error) {
	config := api.GetUnsanitizedConfig()

	dataSource := *config.SqlSettings.DataSource
	driverName := *config.SqlSettings.DriverName

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLStore{
		db:     db,
		api:    api,
		driver: driverName,
	}, nil
}

func (s *SQLStore) GetMessage(messageId string) (*models.Message, error) {
	query := `
		SELECT id, channel_id, user_id, content, created_at, updated_at
		FROM messages
		WHERE id = $1
	`

	var message models.Message
	err := s.db.QueryRow(query, messageId).Scan(
		&message.ID,
		&message.ChannelID,
		&message.UserID,
		&message.Content,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message with id %s not found", messageId)
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

func (s *SQLStore) CreateMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (id, channel_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(query,
		message.ID,
		message.ChannelID,
		message.UserID,
		message.Content,
		message.CreatedAt,
		message.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (s *SQLStore) UpdateMessage(message *models.Message) error {
	query := `
		UPDATE messages
		SET content = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := s.db.Exec(query, message.Content, message.UpdatedAt, message.ID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message with id %s not found", message.ID)
	}

	return nil
}

func (s *SQLStore) DeleteMessage(messageId string) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := s.db.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message with id %s not found", messageId)
	}

	return nil
}

func (s *SQLStore) GetMessagesByChannel(channelId string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, channel_id, user_id, content, created_at, updated_at
		FROM messages
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, channelId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages by channel: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.ChannelID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return messages, nil
}

func (s *SQLStore) GetMessagesByUser(userId string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, channel_id, user_id, content, created_at, updated_at
		FROM messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages by user: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.ChannelID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return messages, nil
}

func (s *SQLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
