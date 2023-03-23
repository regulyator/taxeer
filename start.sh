#!/bin/bash

# Check if enough arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <bot_api_key> <db_password>"
    exit 1
fi

# Set your environment variables from command line arguments
export BOT_API="$1"
export DB_PASSWORD="$2"

docker-compose stop
docker-compose rm -f
# Remove old app image
OLD_APP_IMAGE_ID=$(docker images | grep -E "^taxeer_app\s+latest\s" | awk '{print $3}')
if [ -n "$OLD_APP_IMAGE_ID" ]; then
  docker rmi $OLD_APP_IMAGE_ID
fi
docker-compose build app
docker-compose up -d