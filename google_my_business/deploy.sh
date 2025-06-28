#!/bin/bash

# Deployment script for Google My Business Review Analysis

# Configuration
SSH_KEY="/Users/matejkepes/Downloads/google-reviews-fe-key-pair-london-1.pem"
REMOTE_USER="ubuntu"
REMOTE_HOST="ec2-3-10-143-220.eu-west-2.compute.amazonaws.com"
REMOTE_DIR="/home/ubuntu/Documents/code/golang/google_my_business"
EXECUTABLE="my_business"
BACKUP_EXECUTABLE="${EXECUTABLE}.bkp"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display help
function show_help() {
    echo -e "${BLUE}Google My Business Deployment Script${NC}"
    echo ""
    echo "Usage:"
    echo "  ./deploy.sh [OPTION]"
    echo ""
    echo "Options:"
    echo "  rollback    Rollback to previous version"
    echo "  help        Display this help message"
    echo "  -h, --help  Display this help message"
    echo ""
    echo "Default behavior (no arguments) is to deploy the application."
    echo "This script will build the Go executable locally and deploy it to the server."
    echo "Before deployment, it will ask for confirmation on migrations and configurations."
}

# Function to display messages
function echo_message() {
    echo -e "${2}${1}${NC}"
}

# Function to display error messages and exit
function echo_error() {
    echo_message "$1" "${RED}"
    exit 1
}

# Function to prompt for yes/no confirmation
function confirm() {
    while true; do
        read -p "$1 (y/n): " yn
        case $yn in
            [Yy]* ) return 0;;
            [Nn]* ) return 1;;
            * ) echo "Please answer y or n.";;
        esac
    done
}

# Check for SSH key existence
if [ ! -f "$SSH_KEY" ]; then
    echo_error "SSH key not found at $SSH_KEY"
fi

# Confirm migrations
if ! confirm "Have you run the database migrations manually?"; then
    echo_message "Please run the database migrations before deploying." "${YELLOW}"
    exit 1
fi

# Confirm configurations
if ! confirm "Are the server configurations up to date?"; then
    echo_message "Please update the server configurations before deploying." "${YELLOW}"
    exit 1
fi

# Build the executable
echo_message "Building executable for Linux..." "${GREEN}"
env GOOS=linux GOARCH=amd64 go build my_business.go

if [ $? -ne 0 ]; then
    echo_error "Build failed. Please fix the errors and try again."
fi

echo_message "Build successful!" "${GREEN}"

# Deployment function
function deploy() {
    # Backup the current executable
    echo_message "Backing up current executable on server..." "${GREEN}"
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "cd $REMOTE_DIR && if [ -f $EXECUTABLE ]; then mv $EXECUTABLE $BACKUP_EXECUTABLE; fi"
    
    if [ $? -ne 0 ]; then
        echo_error "Failed to backup the current executable. Aborting deployment."
    fi
    
    # Upload the new executable
    echo_message "Uploading new executable..." "${GREEN}"
    scp -i "$SSH_KEY" "$EXECUTABLE" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/$EXECUTABLE"
    
    if [ $? -ne 0 ]; then
        echo_error "Failed to upload the new executable. Attempting to restore backup."
        ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "cd $REMOTE_DIR && if [ -f $BACKUP_EXECUTABLE ]; then mv $BACKUP_EXECUTABLE $EXECUTABLE; fi"
        exit 1
    fi
    
    # Make the executable... executable
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "chmod +x $REMOTE_DIR/$EXECUTABLE"
    
    # Upload all scripts
    echo_message "Uploading scripts..." "${GREEN}"
    scp -i "$SSH_KEY" scripts/* "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/"
    
    if [ $? -ne 0 ]; then
        echo_message "Warning: Some scripts may not have been uploaded correctly." "${YELLOW}"
    fi
    
    # Make scripts executable
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "chmod +x $REMOTE_DIR/*.sh"
    
    echo_message "Deployment completed successfully!" "${GREEN}"
}

# Rollback function
function rollback() {
    echo_message "Rolling back to previous version..." "${YELLOW}"
    ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" "cd $REMOTE_DIR && if [ -f $BACKUP_EXECUTABLE ]; then mv $BACKUP_EXECUTABLE $EXECUTABLE; else echo 'No backup found!'; fi"
    
    if [ $? -ne 0 ]; then
        echo_error "Rollback failed. Manual intervention may be required."
    fi
    
    echo_message "Rollback completed!" "${GREEN}"
}

# Main script logic
if [ "$1" == "rollback" ]; then
    rollback
elif [ "$1" == "help" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    show_help
else
    deploy
fi

echo_message "Done!" "${GREEN}" 