#!/bin/bash
if [ -z "$GOOGLE_API_KEY" ]
then
    echo "Configure env variable GOOGLE_API_KEY before starting."
    exit 1
fi

# Run test suite
docker-compose up --build tester
docker-compose stop test-db

# Serve 
docker-compose up --build app