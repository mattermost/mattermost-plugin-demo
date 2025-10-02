// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

// FixtureLoader handles loading test fixtures from JSON files
type FixtureLoader struct {
	store      *SQLStore
	fixtureDir string
}

// NewFixtureLoader creates a new fixture loader
func NewFixtureLoader(store *SQLStore, fixtureDir string) *FixtureLoader {
	return &FixtureLoader{
		store:      store,
		fixtureDir: fixtureDir,
	}
}

// SessionFixture represents a session in the fixture file
type SessionFixture struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	CreateAt int64  `json:"create_at"`
	ClosedAt *int64 `json:"closed_at"`
}

// LoadSessions loads session fixtures from a JSON file
func (fl *FixtureLoader) LoadSessions(filename string) ([]*model.Session, error) {
	filepath := filepath.Join(fl.fixtureDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filepath, err)
	}

	var fixtures []SessionFixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture file %s: %w", filepath, err)
	}

	sessions := make([]*model.Session, 0, len(fixtures))
	for _, fixture := range fixtures {
		session := &model.Session{
			ID:       fixture.ID,
			UserID:   fixture.UserID,
			CreateAt: fixture.CreateAt,
			ClosedAt: fixture.ClosedAt,
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// InsertSessions inserts sessions from fixtures into the database
func (fl *FixtureLoader) InsertSessions(filename string) error {
	sessions, err := fl.LoadSessions(filename)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := fl.store.CreateSession(session); err != nil {
			return fmt.Errorf("failed to insert session %s: %w", session.ID, err)
		}
	}

	return nil
}

// LoadAllFixtures loads all fixture files for the store
func (fl *FixtureLoader) LoadAllFixtures() error {
	// Load sessions
	if err := fl.InsertSessions("sessions.json"); err != nil {
		return fmt.Errorf("failed to load session fixtures: %w", err)
	}

	return nil
}
