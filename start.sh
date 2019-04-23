#!/bin/bash
if [ -z "$GOOGLE_API_KEY" ]
then
    echo "Configure env variable GOOGLE_API_KEY before starting."
    exit 1
fi

# Run unit tests and prepare app image
docker-compose build app

docker-compose up -d app

# TODO: Run integration tests