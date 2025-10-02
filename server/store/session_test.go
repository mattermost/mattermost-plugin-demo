// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

func TestSessionCRUD(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Create and Get Session", func(t *testing.T) {
		// Clean table before test
		TruncateTable(t, store, "session")

		session := &model.Session{
			ID:     "test-session-1",
			UserID: "user1",
		}

		// Create session
		err := store.CreateSession(session)
		require.NoError(t, err)

		// Verify session was created
		count := GetTableRowCount(t, store, "session")
		assert.Equal(t, 1, count)

		// Get session by ID
		retrieved, err := store.GetSessionByID(session.ID)
		require.NoError(t, err)
		assert.Equal(t, session.ID, retrieved.ID)
		assert.Equal(t, session.UserID, retrieved.UserID)
		assert.Greater(t, retrieved.CreateAt, int64(0))
		assert.Nil(t, retrieved.ClosedAt)
	})

	t.Run("Update Session", func(t *testing.T) {
		// Clean table
		TruncateTable(t, store, "session")

		session := &model.Session{
			ID:     "test-session-2",
			UserID: "user2",
		}

		// Create session
		err := store.CreateSession(session)
		require.NoError(t, err)

		// Update session with closed time
		closedAt := int64(1234567890)
		session.ClosedAt = &closedAt
		session.UserID = "user2-updated"

		err = store.UpdateSession(session)
		require.NoError(t, err)

		// Verify update
		retrieved, err := store.GetSessionByID(session.ID)
		require.NoError(t, err)
		assert.Equal(t, "user2-updated", retrieved.UserID)
		assert.NotNil(t, retrieved.ClosedAt)
		assert.Equal(t, closedAt, *retrieved.ClosedAt)
	})

	t.Run("Delete Session", func(t *testing.T) {
		// Clean table
		TruncateTable(t, store, "session")

		session := &model.Session{
			ID:     "test-session-3",
			UserID: "user3",
		}

		// Create session
		err := store.CreateSession(session)
		require.NoError(t, err)

		// Delete session
		err = store.DeleteSession(session.ID)
		require.NoError(t, err)

		// Verify deletion
		count := GetTableRowCount(t, store, "session")
		assert.Equal(t, 0, count)

		// Try to get deleted session
		_, err = store.GetSessionByID(session.ID)
		assert.Error(t, err)
	})

	t.Run("Get All Sessions", func(t *testing.T) {
		// Clean table
		TruncateTable(t, store, "session")

		// Create multiple sessions
		sessions := []*model.Session{
			{ID: "session-1", UserID: "user1"},
			{ID: "session-2", UserID: "user2"},
			{ID: "session-3", UserID: "user3"},
		}

		for _, session := range sessions {
			err := store.CreateSession(session)
			require.NoError(t, err)
		}

		// Get all sessions
		retrieved, err := store.GetSessions()
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)
	})

	t.Run("Session Not Found", func(t *testing.T) {
		TruncateTable(t, store, "session")

		_, err := store.GetSessionByID("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete Non-Existent Session", func(t *testing.T) {
		TruncateTable(t, store, "session")

		err := store.DeleteSession("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestSessionValidation(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Create Session Without ID", func(t *testing.T) {
		TruncateTable(t, store, "session")
		session := &model.Session{
			UserID: "user1",
		}

		// Should automatically generate ID
		err := store.CreateSession(session)
		require.NoError(t, err)
		assert.NotEmpty(t, session.ID)

		// Clean up
		_ = store.DeleteSession(session.ID)
	})

	t.Run("Create Session Without UserID", func(t *testing.T) {
		TruncateTable(t, store, "session")
		session := &model.Session{
			ID: "test-session",
		}

		// Should fail validation
		err := store.CreateSession(session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID")
	})

	t.Run("Update Invalid Session", func(t *testing.T) {
		TruncateTable(t, store, "session")
		session := &model.Session{
			ID: "test-session",
			// Missing UserID
		}

		err := store.UpdateSession(session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})
}

func TestSessionConcurrency(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	// Clean table
	TruncateTable(t, store, "session")

	t.Run("Concurrent Session Creation", func(t *testing.T) {
		const numGoroutines = 10

		done := make(chan bool, numGoroutines)

		// Create sessions concurrently
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				session := &model.Session{
					UserID: "concurrent-user",
				}
				err := store.CreateSession(session)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify all sessions were created
		count := GetTableRowCount(t, store, "session")
		assert.Equal(t, numGoroutines, count)
	})
}

func TestSessionWithBinaryParams(t *testing.T) {
	store, cleanup := SetupTestStoreWithBinaryParams(t, true)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	// Verify binary params are enabled
	assert.True(t, store.isBinaryParams, "Binary params should be enabled")

	t.Run("CRUD Operations With Binary Params", func(t *testing.T) {
		TruncateTable(t, store, "session")
		session := &model.Session{
			ID:     "binary-test-session",
			UserID: "binary-user",
		}

		// Create
		err := store.CreateSession(session)
		require.NoError(t, err)

		// Read
		retrieved, err := store.GetSessionByID(session.ID)
		require.NoError(t, err)
		assert.Equal(t, session.ID, retrieved.ID)
		assert.Equal(t, session.UserID, retrieved.UserID)

		// Update
		closedAt := int64(9999999999)
		retrieved.ClosedAt = &closedAt
		err = store.UpdateSession(retrieved)
		require.NoError(t, err)

		// Verify update
		updated, err := store.GetSessionByID(session.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.ClosedAt)
		assert.Equal(t, closedAt, *updated.ClosedAt)

		// Delete
		err = store.DeleteSession(session.ID)
		require.NoError(t, err)
	})
}

func TestGetSchemaName(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Get Schema Name", func(t *testing.T) {
		TruncateTable(t, store, "session")
		schema, err := store.GetSchemaName()
		require.NoError(t, err)
		assert.NotEmpty(t, schema)
		// PostgreSQL default schema is 'public'
		assert.Equal(t, "public", schema)
	})
}

func TestMigrations(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Verify Session Table Created", func(t *testing.T) {
		TruncateTable(t, store, "session")
		exists := TableExists(t, store, "session")
		assert.True(t, exists, "Session table should exist after setup")
	})

	t.Run("Verify Session Table Columns", func(t *testing.T) {
		TruncateTable(t, store, "session")
		columns := GetTableColumns(t, store, "session")
		assert.Contains(t, columns, "id")
		assert.Contains(t, columns, "user_id")
		assert.Contains(t, columns, "create_at")
		assert.Contains(t, columns, "closed_at")
	})
}

func TestStoreShutdown(t *testing.T) {
	store, cleanup := SetupTestStore(t)
	defer cleanup()

	if store == nil {
		t.Skip("PostgreSQL not available")
	}

	t.Run("Shutdown Store", func(t *testing.T) {
		TruncateTable(t, store, "session")
		err := store.Shutdown()
		assert.NoError(t, err)
	})
}
