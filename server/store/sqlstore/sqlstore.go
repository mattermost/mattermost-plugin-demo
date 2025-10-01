package sqlstore

import (
	"database/sql"
	"errors"
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

func (S SQLStore) GetSession(sessionId string) (*models.Session, error) {
	query := "SELECT id, user_id, created_at, closed_at FROM whatsapp_plugin_session WHERE id = $1"
	var session models.Session
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

func (S SQLStore) CreateSession(userId string) (*models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLStore) CloseSession(sessionId string) (*models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLStore) GetSessionByUserId(userId string) (*models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLStore) GetSessionsUnclosed() ([]models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLStore) GetWhatsappChannels() ([]models.WhatsappChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLStore) CreateWhatsappChannel(channelId string) (*models.WhatsappChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (S *SQLStore) Close() error {
	if S.db != nil {
		return S.db.Close()
	}
	return nil
}

func New(api plugin.API) (store.Store, error) {
	config := api.GetUnsanitizedConfig()

	dataSource := "postgres://mmuser:mmuser_password@postgres:5432/mattermost?connect_timeout=30&sslmode=disable"

	driverName := *config.SqlSettings.DriverName

	api.LogInfo("Connecting to DB",
		"driver", driverName,
		"dataSource", dataSource,
	)

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		api.LogInfo("Erorr connecting to DB",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		api.LogInfo("Erorr connecting to DB",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLStore{
		db:     db,
		api:    api,
		driver: driverName,
	}, nil
}
