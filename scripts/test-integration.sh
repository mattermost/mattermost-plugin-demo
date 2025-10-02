#!/bin/bash

# Integration Test Script for Store
# This script runs integration tests against a dedicated PostgreSQL instance
# with real migrations and fixtures

set -e

echo "=== Starting Integration Tests ==="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
POSTGRES_SERVICE="postgres_test"
TEST_NETWORK="mattermost_plugin_testing"
TEST_DB="mattermost_test"
TEST_USER="mmuser"
TEST_PASSWORD="mmuser_password"

# Function to check if PostgreSQL is ready
wait_for_postgres() {
    echo "⏳ Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if docker compose exec -T $POSTGRES_SERVICE pg_isready -U $TEST_USER > /dev/null 2>&1; then
            echo -e "${GREEN}✓${NC} PostgreSQL is ready"
            return 0
        fi
        echo "  Attempt $i/30..."
        sleep 1
    done
    echo -e "${RED}✗${NC} PostgreSQL failed to start"
    return 1
}

# Function to cleanup
cleanup() {
    echo "🧹 Cleaning up..."
    docker compose down $POSTGRES_SERVICE
}

# Trap cleanup on exit
trap cleanup EXIT

echo "🚀 Starting PostgreSQL test instance..."
docker compose up -d $POSTGRES_SERVICE

if ! wait_for_postgres; then
    exit 1
fi

echo "🗄️  Dropping and recreating test database..."
docker compose exec -T $POSTGRES_SERVICE psql -U $TEST_USER -d postgres -c "DROP DATABASE IF EXISTS $TEST_DB;" || true
docker compose exec -T $POSTGRES_SERVICE psql -U $TEST_USER -d postgres -c "CREATE DATABASE $TEST_DB;"

echo "📦 Running migrations..."
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=$TEST_USER
export TEST_DB_PASSWORD=$TEST_PASSWORD
export TEST_DB_NAME=$TEST_DB
export TEST_DB_SSLMODE=disable
export RUN_INTEGRATION_TESTS=1

echo "🧪 Running integration tests..."
cd server/store
go test -v -tags=integration -count=1 ./...

echo -e "${GREEN}✓${NC} Integration tests completed successfully!"
