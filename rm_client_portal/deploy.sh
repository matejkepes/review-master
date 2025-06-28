#!/bin/bash

# Hardcoded defaults
SSH_KEY_PATH="/Users/matejkepes/Downloads/google-reviews-client-portal-key-pair-london-1.pem"
SERVER_ADDRESS="ubuntu@ec2-18-130-28-184.eu-west-2.compute.amazonaws.com"
REMOTE_PATH="/home/ubuntu/Documents/code/golang/rm_client_portal"

# Allow overriding defaults with command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --key=*)
      SSH_KEY_PATH="${1#*=}"
      shift
      ;;
    --server=*)
      SERVER_ADDRESS="${1#*=}"
      shift
      ;;
    --path=*)
      REMOTE_PATH="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown parameter: $1"
      exit 1
      ;;
  esac
done

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored status messages
print_status() {
  local color=$1
  local message=$2
  echo -e "${color}${message}${NC}"
}

# Function to check if command succeeded
check_status() {
  if [ $? -eq 0 ]; then
    print_status "$GREEN" "✅ Success"
  else
    print_status "$RED" "❌ Failed"
    exit 1
  fi
}

# Current timestamp for backup files
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Step 1: Build Go application for Linux
print_status "$YELLOW" "📦 Building Go application for Linux..."
env GOOS=linux GOARCH=amd64 go build rm_client_portal.go
check_status

# Step 2: Backup and deploy the executable
print_status "$YELLOW" "🔄 Backing up existing executable on server..."
ssh -i "$SSH_KEY_PATH" -o ServerAliveInterval=30 -o ServerAliveCountMax=3 "$SERVER_ADDRESS" "if [ -f $REMOTE_PATH/rm_client_portal ]; then mv $REMOTE_PATH/rm_client_portal $REMOTE_PATH/rm_client_portal.bkp.$TIMESTAMP; fi"
check_status

print_status "$YELLOW" "📤 Uploading new executable to server..."
rsync -avz --progress -e "ssh -i $SSH_KEY_PATH -o ServerAliveInterval=30 -o ServerAliveCountMax=3" rm_client_portal "$SERVER_ADDRESS:$REMOTE_PATH/rm_client_portal"
check_status

# Step 3: Restart the service
print_status "$YELLOW" "🔄 Restarting the service..."
ssh -i "$SSH_KEY_PATH" "$SERVER_ADDRESS" "sudo systemctl restart rm_client_portal.service"
check_status

# Step 4: Check service status
print_status "$YELLOW" "📊 Checking service status..."
ssh -i "$SSH_KEY_PATH" "$SERVER_ADDRESS" "sudo systemctl status rm_client_portal.service"

print_status "$GREEN" "✅ Deployment completed successfully!"   