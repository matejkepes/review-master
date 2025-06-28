#!/bin/bash

# Configuration
SSH_KEY="/Users/matejkepes/Downloads/google-reviews-mysql-key-pair-london-1.pem"
SERVER="ubuntu@ec2-35-178-184-3.eu-west-2.compute.amazonaws.com"
REMOTE_SQL_DIR="~/google_reviews/sql"
REMOTE_BACKUP_DIR="~/google_reviews/db-backup"
LOCAL_SQL_DIR="./sql"
DB_NAME="google_reviews"
DATE_FORMAT=$(date +"%y%m%d")
DRY_RUN=false

# Function to display help
show_help() {
    echo "Database Migration Script"
    echo "-------------------------"
    echo "This script automates the process of deploying database migrations to the server."
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help          Display this help message"
    echo "  -k, --key FILE      SSH key file (default: $SSH_KEY)"
    echo "  -s, --server HOST   Server hostname (default: $SERVER)"
    echo "  -d, --db NAME       Database name (default: $DB_NAME)"
    echo "  -l, --local DIR     Local SQL directory (default: $LOCAL_SQL_DIR)"
    echo "  --dry-run           Check what would be copied without making changes"
    echo
    echo "Example:"
    echo "  $0 --key ~/.ssh/my-key.pem --server user@example.com"
    echo
    echo "Process:"
    echo "  1. Checks existing migrations on the server"
    echo "  2. Copies new migration files"
    echo "  3. Creates database backups (schema-only and full)"
    echo "  4. Runs new migrations"
    echo
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -k|--key)
                SSH_KEY="$2"
                shift 2
                ;;
            -s|--server)
                SERVER="$2"
                shift 2
                ;;
            -d|--db)
                DB_NAME="$2"
                shift 2
                ;;
            -l|--local)
                LOCAL_SQL_DIR="$2"
                shift 2
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            *)
                echo "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Function to check if a migration exists on the server
check_existing_migrations() {
    echo "Checking existing migrations on server..."
    ssh -i $SSH_KEY $SERVER "ls -1 $REMOTE_SQL_DIR" > /tmp/remote_migrations.txt
    
    # Find the highest numbered migration file on the server
    highest_remote=$(ssh -i $SSH_KEY $SERVER "ls -1 $REMOTE_SQL_DIR" | grep -E '^[0-9]+_' | sort -n | tail -1 | cut -d'_' -f1)
    echo "Highest migration on server: $highest_remote"
    
    # Find all local migrations with higher numbers
    for local_file in $(ls -1 $LOCAL_SQL_DIR | grep -E '^[0-9]+_' | sort -n); do
        local_num=$(echo $local_file | cut -d'_' -f1)
        if [[ $local_num -gt $highest_remote ]]; then
            echo "Found new migration to copy: $local_file"
            new_migrations+=("$local_file")
        fi
    done
    
    if [ ${#new_migrations[@]} -eq 0 ]; then
        echo "No new migrations to apply. Exiting."
        exit 0
    fi
}

# Function to copy new migrations to server
copy_migrations() {
    echo "Copying new migrations to server..."
    for migration in "${new_migrations[@]}"; do
        echo "Copying $migration..."
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY RUN] Would copy $migration to server"
        else
            scp -i $SSH_KEY "$LOCAL_SQL_DIR/$migration" "$SERVER:$REMOTE_SQL_DIR/"
            if [ $? -ne 0 ]; then
                echo "Error copying migration file $migration. Exiting."
                exit 1
            fi
        fi
    done
}

# Function to create database backups
create_backups() {
    echo "Creating database backups..."
    
    # Prompt for MySQL password
    read -sp "Enter MySQL root password: " MYSQL_PASSWORD
    echo
    
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] Would create schema-only backup: ${DB_NAME}_nodata_$DATE_FORMAT.sql"
        echo "[DRY RUN] Would create full backup: ${DB_NAME}_$DATE_FORMAT.sql"
        return
    fi
    
    # Create schema-only backup
    echo "Creating schema-only backup..."
    ssh -i $SSH_KEY $SERVER "mysqldump $DB_NAME -uroot -p$MYSQL_PASSWORD --no-data > $REMOTE_BACKUP_DIR/${DB_NAME}_nodata_$DATE_FORMAT.sql"
    if [ $? -ne 0 ]; then
        echo "Error creating schema-only backup. Exiting."
        exit 1
    fi
    
    # Create full backup
    echo "Creating full backup..."
    ssh -i $SSH_KEY $SERVER "mysqldump $DB_NAME -uroot -p$MYSQL_PASSWORD > $REMOTE_BACKUP_DIR/${DB_NAME}_$DATE_FORMAT.sql"
    if [ $? -ne 0 ]; then
        echo "Error creating full backup. Exiting."
        exit 1
    fi
}

# Function to run migrations
run_migrations() {
    echo "Running migrations..."
    for migration in "${new_migrations[@]}"; do
        echo "Running migration $migration..."
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY RUN] Would run migration: $migration"
        else
            ssh -i $SSH_KEY $SERVER "mysql $DB_NAME -uroot -p$MYSQL_PASSWORD < $REMOTE_SQL_DIR/$migration"
            if [ $? -ne 0 ]; then
                echo "Error running migration $migration. Exiting."
                exit 1
            fi
            echo "Migration $migration completed successfully."
        fi
    done
}

# Main execution
main() {
    echo "Starting database migration process..."
    new_migrations=()
    
    check_existing_migrations
    copy_migrations
    create_backups
    run_migrations
    
    echo "All migrations completed successfully!"
}

# Parse arguments
parse_args "$@"

# Execute main function
main 