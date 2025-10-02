// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package store

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	demoModel "github.com/itstar-tech/mattermost-plugin-demo/server/model"
)

const (
	testTablePrefix = "test_demo_plugin_"
)

// TestConfig holds configuration for test database
type TestConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetTestConfig returns configuration for test database
// It can read from environment variables or use defaults
func GetTestConfig() TestConfig {
	return TestConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "mmuser"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "mmuser_password"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "mattermost_test"),
		SSLMode:  getEnvOrDefault("TEST_DB_SSLMODE", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetTestConnectionString builds a PostgreSQL connection string for tests
func (tc TestConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=10",
		tc.User,
		tc.Password,
		tc.Host,
		tc.Port,
		tc.DBName,
		tc.SSLMode,
	)
}

// GetTestConnectionStringWithBinaryParams builds a PostgreSQL connection string with binary_parameters=yes
func (tc TestConfig) GetTestConnectionStringWithBinaryParams() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=10&binary_parameters=yes",
		tc.User,
		tc.Password,
		tc.Host,
		tc.Port,
		tc.DBName,
		tc.SSLMode,
	)
}

// SetupTestDB creates a test database connection and ensures it's ready
func SetupTestDB(t *testing.T) (*sql.DB, TestConfig) {
	t.Helper()

	config := GetTestConfig()

	// First connect to postgres database to create test database if needed
	adminConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=%s&connect_timeout=10",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.SSLMode,
	)

	adminDB, err := sql.Open("postgres", adminConnString)
	if err != nil {
		t.Skipf("Could not connect to PostgreSQL (this is OK if not running locally): %v", err)
		return nil, config
	}
	defer adminDB.Close()

	// Test connection
	err = adminDB.Ping()
	if err != nil {
		t.Skipf("Could not ping PostgreSQL (this is OK if not running locally): %v", err)
		return nil, config
	}

	// Create test database if it doesn't exist
	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil {
		// Ignore error if database already exists
		t.Logf("Database may already exist: %v", err)
	}

	// Connect to test database
	db, err := sql.Open("postgres", config.GetConnectionString())
	require.NoError(t, err, "Should connect to test database")

	err = db.Ping()
	require.NoError(t, err, "Should ping test database")

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return db, config
}

// CleanupTestDB drops all test tables
func CleanupTestDB(t *testing.T, db *sql.DB, tablePrefix string) {
	t.Helper()

	if db == nil {
		return
	}

	// Check if the connection is still alive
	if err := db.Ping(); err != nil {
		// Connection is already closed, nothing to cleanup
		return
	}

	// Drop all tables with test prefix
	_, err := db.Exec(fmt.Sprintf(`
		DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE '%s%%') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`, tablePrefix))
	if err != nil {
		t.Logf("Warning: Could not cleanup test tables: %v", err)
	}

	// Note: We don't close the connection here anymore
	// It will be closed by store.Shutdown() in the cleanup function
}

// SetupTestStore creates a SQLStore instance for testing
func SetupTestStore(t *testing.T) (*SQLStore, func()) {
	return SetupTestStoreWithBinaryParams(t, false)
}

// SetupTestStoreWithBinaryParams creates a SQLStore instance for testing with optional binary params
func SetupTestStoreWithBinaryParams(t *testing.T, binaryParams bool) (*SQLStore, func()) {
	t.Helper()

	db, config := SetupTestDB(t)
	if db == nil {
		return nil, func() {}
	}

	connectionString := config.GetConnectionString()
	if binaryParams {
		connectionString = config.GetTestConnectionStringWithBinaryParams()
	}

	// Create store without migrations for unit tests
	// Integration tests with real migrations should be in separate test files
	store := &SQLStore{
		db:               db,
		dbType:           demoModel.DBTypePostgres,
		tablePrefix:      testTablePrefix,
		connectionString: connectionString,
		skipMigrations:   true,
	}

	var err error
	store.isBinaryParams, err = store.checkBinaryParams()
	require.NoError(t, err, "Should check binary params")

	store.schemaName, err = store.GetSchemaName()
	require.NoError(t, err, "Should get schema name")

	// Create tables manually for unit tests
	// This is faster than running migrations for each test
	err = createTestSchema(t, store)
	require.NoError(t, err, "Should create test schema")

	// Cleanup function
	cleanup := func() {
		// First cleanup tables (while connection is still open)
		CleanupTestDB(t, db, testTablePrefix)

		// Then shutdown the store (closes connection)
		if store != nil {
			_ = store.Shutdown()
		}
	}

	return store, cleanup
}

// createTestSchema creates the necessary tables for testing
func createTestSchema(t *testing.T, store *SQLStore) error {
	t.Helper()

	// Create session table
	_, err := store.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %ssession (
			id VARCHAR(26) PRIMARY KEY,
			user_id VARCHAR(26) NOT NULL,
			create_at BIGINT NOT NULL,
			closed_at BIGINT
		)
	`, store.tablePrefix))

	if err != nil {
		return fmt.Errorf("failed to create session table: %w", err)
	}

	// Create index
	_, err = store.db.Exec(fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_%ssession_user_id ON %ssession(user_id)
	`, store.tablePrefix, store.tablePrefix))

	if err != nil {
		return fmt.Errorf("failed to create session index: %w", err)
	}

	return nil
}

// TruncateTable truncates a specific table (useful for test isolation)
func TruncateTable(t *testing.T, store *SQLStore, tableName string) {
	t.Helper()

	fullTableName := store.tablePrefix + tableName
	_, err := store.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", fullTableName))
	require.NoError(t, err, "Should truncate table %s", fullTableName)
}

// GetTableRowCount returns the number of rows in a table
func GetTableRowCount(t *testing.T, store *SQLStore, tableName string) int {
	t.Helper()

	fullTableName := store.tablePrefix + tableName
	var count int
	err := store.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", fullTableName)).Scan(&count)
	require.NoError(t, err, "Should count rows in table %s", fullTableName)
	return count
}

// TableExists checks if a table exists in the database
func TableExists(t *testing.T, store *SQLStore, tableName string) bool {
	t.Helper()

	fullTableName := store.tablePrefix + tableName
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`
	err := store.db.QueryRow(query, fullTableName).Scan(&exists)
	require.NoError(t, err, "Should check if table exists")
	return exists
}

// GetTableColumns returns the column names for a table
func GetTableColumns(t *testing.T, store *SQLStore, tableName string) []string {
	t.Helper()

	fullTableName := store.tablePrefix + tableName
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := store.db.Query(query, fullTableName)
	require.NoError(t, err, "Should query table columns")
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		err := rows.Scan(&column)
		require.NoError(t, err, "Should scan column name")
		columns = append(columns, column)
	}

	return columns
}
