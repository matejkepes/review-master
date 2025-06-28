#!/bin/bash

# Help message
show_help() {
  echo "Usage: ./deploy.sh [options]"
  echo ""
  echo "Build and deploy the Quasar application to the production server."
  echo ""
  echo "Options:"
  echo "  -h, --help     Display this help message and exit"
  echo "  -b, --build    Only build the application, don't deploy"
  echo "  -d, --deploy   Only deploy the application, don't build"
  echo ""
  echo "Without options, the script will both build and deploy the application."
}

# Parse command line arguments
BUILD=true
DEPLOY=true

for arg in "$@"
do
  case $arg in
    -h|--help)
      show_help
      exit 0
      ;;
    -b|--build)
      BUILD=true
      DEPLOY=false
      ;;
    -d|--deploy)
      BUILD=false
      DEPLOY=true
      ;;
    *)
      echo "Unknown option: $arg"
      show_help
      exit 1
      ;;
  esac
done

# Constants
SERVER="ubuntu@ec2-3-10-143-220.eu-west-2.compute.amazonaws.com"
KEY_PATH="/Users/matejkepes/Downloads/google-reviews-fe-key-pair-london-1.pem"
REMOTE_PATH="~/Documents/code/golang/google_reviews_ui/static"
LOCAL_DIST="./dist"

# Function to check if a command was successful
check_status() {
  if [ $? -eq 0 ]; then
    echo "✅ $1"
  else
    echo "❌ $1"
    exit 1
  fi
}

# Build the application
if [ "$BUILD" = true ]; then
  echo "🔨 Building the application..."
  npx quasar build
  check_status "Application build"
  echo ""
fi

# Deploy the application
if [ "$DEPLOY" = true ]; then
  echo "🚀 Deploying to server: $SERVER"

  # Backup and remove existing dist folder on server
  echo "📦 Creating backup of existing dist folder on server..."
  ssh -i "$KEY_PATH" "$SERVER" "if [ -d $REMOTE_PATH/dist ]; then cp -r $REMOTE_PATH/dist $REMOTE_PATH/dist_bkp && rm -rf $REMOTE_PATH/dist; fi"
  check_status "Backup creation and old dist removal"

  # Copy new dist folder to server
  echo "📤 Copying new dist folder to server..."
  scp -r -i "$KEY_PATH" "$LOCAL_DIST" "$SERVER:$REMOTE_PATH/"
  check_status "File transfer"

  echo "🎉 Deployment completed successfully!"
fi

exit 0
