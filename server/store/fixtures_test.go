// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixtureLoader(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Load Session Fixtures", func(t *testing.T) {
		TruncateTable(t, store, "session")

		// Create fixture loader
		loader := NewFixtureLoader(store, "testdata/fixtures")

		// Load sessions from fixtures file
		err := loader.InsertSessions("sessions.json")
		require.NoError(t, err, "Should load session fixtures")

		// Verify sessions were loaded
		count := GetTableRowCount(t, store, "session")
		assert.Equal(t, 5, count, "Should have loaded 5 sessions from fixtures")

		// Verify we can query the fixtures
		sessions, err := store.GetSessions()
		require.NoError(t, err)
		assert.Len(t, sessions, 5)

		// Verify specific fixture data
		session, err := store.GetSessionByID("session-fixture-001")
		require.NoError(t, err)
		assert.Equal(t, "session-fixture-001", session.ID)
		assert.Equal(t, "user-001", session.UserID)
		assert.Nil(t, session.ClosedAt)

		// Verify closed session
		closedSession, err := store.GetSessionByID("session-fixture-003")
		require.NoError(t, err)
		assert.NotNil(t, closedSession.ClosedAt)
		assert.Equal(t, int64(1704326400000), *closedSession.ClosedAt)
	})

	t.Run("Load All Fixtures", func(t *testing.T) {
		TruncateTable(t, store, "session")

		loader := NewFixtureLoader(store, "testdata/fixtures")

		err := loader.LoadAllFixtures()
		require.NoError(t, err, "Should load all fixtures")

		count := GetTableRowCount(t, store, "session")
		assert.Greater(t, count, 0, "Should have loaded fixtures")
	})

	t.Run("Fixture File Not Found", func(t *testing.T) {
		loader := NewFixtureLoader(store, "testdata/fixtures")

		err := loader.InsertSessions("nonexistent.json")
		assert.Error(t, err, "Should fail when fixture file doesn't exist")
		assert.Contains(t, err.Error(), "failed to read fixture file")
	})
}
