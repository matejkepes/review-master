#!/bin/bash

# Google Reviews Database Migration Script
# This script drops and recreates the database, then runs all SQL migration files in the correct order

# Default database connection parameters
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="google_reviews"
DB_USER="root"
DB_PASSWORD=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --host HOST         Database host (default: localhost)"
    echo "  -P, --port PORT         Database port (default: 3306)"
    echo "  -d, --database DATABASE Database name (default: google_reviews)"
    echo "  -u, --user USER         Database user (default: root)"
    echo "  --no-drop               Skip dropping the database (default: false)"
    echo "  --help                  Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 -h localhost -d google_reviews -u root"
    echo "  $0 --no-drop  # Skip dropping the database"
    echo ""
    echo "Note: You will be prompted for the MySQL password during execution."
}

# Parse command line arguments
DROP_DATABASE=true
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--host)
            DB_HOST="$2"
            shift 2
            ;;
        -P|--port)
            DB_PORT="$2"
            shift 2
            ;;
        -d|--database)
            DB_NAME="$2"
            shift 2
            ;;
        -u|--user)
            DB_USER="$2"
            shift 2
            ;;
        --no-drop)
            DROP_DATABASE=false
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check if mysql client is available
if ! command -v mysql &> /dev/null; then
    print_error "MySQL client is not installed or not in PATH"
    exit 1
fi

# Prompt for MySQL password
echo -n "Enter MySQL password for user '$DB_USER': "
read -s DB_PASSWORD
echo

# Build MySQL connection string
MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD"

# Test database connection (connect to mysql system database first)
print_info "Testing database connection..."
echo "SELECT 1;" | $MYSQL_CMD mysql &>/dev/null
if [[ $? -ne 0 ]]; then
    print_error "Failed to connect to database. Please check your connection parameters."
    exit 1
fi
print_info "Database connection successful"

# Drop and recreate database if requested
if [[ "$DROP_DATABASE" == true ]]; then
    print_warning "This will DROP the entire '$DB_NAME' database and recreate it!"
    echo -n "Are you sure you want to continue? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_info "Operation cancelled by user"
        exit 0
    fi
    
    print_info "Dropping database '$DB_NAME'..."
    echo "DROP DATABASE IF EXISTS \`$DB_NAME\`;" | $MYSQL_CMD mysql
    if [[ $? -ne 0 ]]; then
        print_error "Failed to drop database '$DB_NAME'"
        exit 1
    fi
    
    print_info "Creating database '$DB_NAME'..."
    echo "CREATE DATABASE \`$DB_NAME\` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;" | $MYSQL_CMD mysql
    if [[ $? -ne 0 ]]; then
        print_error "Failed to create database '$DB_NAME'"
        exit 1
    fi
    
    print_info "Database '$DB_NAME' has been recreated successfully"
fi

# Now connect to the target database
MYSQL_CMD="$MYSQL_CMD $DB_NAME"

# Change to the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_DIR="$SCRIPT_DIR/sql"

# Check if sql directory exists
if [[ ! -d "$SQL_DIR" ]]; then
    print_error "SQL directory not found: $SQL_DIR"
    exit 1
fi

print_info "Starting database migrations..."
print_info "SQL Directory: $SQL_DIR"
print_info "Database: $DB_NAME on $DB_HOST:$DB_PORT"

# Get all SQL files and sort them numerically
SQL_FILES=($(find "$SQL_DIR" -name "*.sql" -type f | sort -V))

if [[ ${#SQL_FILES[@]} -eq 0 ]]; then
    print_warning "No SQL files found in $SQL_DIR"
    exit 0
fi

print_info "Found ${#SQL_FILES[@]} migration files"

# Execute each SQL file
EXECUTED_COUNT=0
FAILED_COUNT=0

for sql_file in "${SQL_FILES[@]}"; do
    filename=$(basename "$sql_file")
    print_info "Executing: $filename"
    
    # Execute SQL file and capture both stdout and stderr
    if $MYSQL_CMD < "$sql_file" 2>/tmp/mysql_error.log; then
        print_info "✓ Successfully executed: $filename"
        ((EXECUTED_COUNT++))
    else
        print_error "✗ Failed to execute: $filename"
        if [[ -f /tmp/mysql_error.log ]]; then
            print_error "MySQL Error:"
            cat /tmp/mysql_error.log | sed 's/^/    /'
        fi
        ((FAILED_COUNT++))
        
        # Ask user if they want to continue
        echo -n "Do you want to continue with the remaining migrations? (y/n): "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_warning "Migration process stopped by user"
            break
        fi
    fi
done

# Cleanup temp file
rm -f /tmp/mysql_error.log

# Summary
echo ""
print_info "Migration Summary:"
print_info "  Total files found: ${#SQL_FILES[@]}"
print_info "  Successfully executed: $EXECUTED_COUNT"
if [[ $FAILED_COUNT -gt 0 ]]; then
    print_error "  Failed: $FAILED_COUNT"
else
    print_info "  Failed: $FAILED_COUNT"
fi

if [[ $FAILED_COUNT -eq 0 ]]; then
    print_info "All migrations completed successfully! 🎉"
    print_info "Database '$DB_NAME' is ready to use."
    exit 0
else
    print_warning "Some migrations failed. Please check the errors above."
    exit 1
fi 