#!/bin/bash
# EngLog Database Setup Script
# "Data is the new oil" - Clive Humby 🛢️

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default environment variables
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-englog}
DB_USER=${DB_USER:-englog}
DB_PASSWORD=${DB_PASSWORD:-englog_dev_password}

echo -e "${BLUE}🔧 EngLog Database Setup${NC}"
echo "========================="

# Check if PostgreSQL is running
echo -e "${YELLOW}📊 Checking PostgreSQL connection...${NC}"
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres >/dev/null 2>&1; then
    echo -e "${RED}❌ PostgreSQL is not running or not accessible${NC}"
    echo "Please ensure PostgreSQL is running and accessible with the provided credentials"
    echo "Host: $DB_HOST:$DB_PORT"
    exit 1
fi
echo -e "${GREEN}✅ PostgreSQL is running${NC}"

# Create database if it doesn't exist
echo -e "${YELLOW}🏗️  Creating database if it doesn't exist...${NC}"
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || echo "Database already exists"

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo -e "${YELLOW}📦 Installing goose...${NC}"
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Database connection string
DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run migrations
echo -e "${YELLOW}🚀 Running database migrations...${NC}"
goose -dir migrations postgres "$DB_URL" up

# Check migration status
echo -e "${YELLOW}📋 Migration status:${NC}"
goose -dir migrations postgres "$DB_URL" status

# Generate SQLC code
echo -e "${YELLOW}⚙️  Generating SQLC code...${NC}"
if command -v sqlc &> /dev/null; then
    sqlc generate
    echo -e "${GREEN}✅ SQLC code generated${NC}"
else
    echo -e "${YELLOW}⚠️  SQLC not found, installing...${NC}"
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    sqlc generate
    echo -e "${GREEN}✅ SQLC code generated${NC}"
fi

# Run a test query to verify everything works
echo -e "${YELLOW}🧪 Testing database connection...${NC}"
if psql "$DB_URL" -c "SELECT 'Database setup successful!' as status;" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Database setup completed successfully!${NC}"
    echo ""
    echo -e "${BLUE}📊 Database Information:${NC}"
    echo "  Host: $DB_HOST:$DB_PORT"
    echo "  Database: $DB_NAME"
    echo "  User: $DB_USER"
    echo ""
    echo -e "${GREEN}🎉 You can now start the EngLog services!${NC}"
else
    echo -e "${RED}❌ Database test failed${NC}"
    exit 1
fi
