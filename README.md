# Basic API Server

## Quick Start

1. set your Google Maps API key via:

    ```bash
    export GOOGLE_API_KEY="<INSERT KEY HERE>"
    ```

2. Run startup script including all unit and integration tests:

    ```bash
    sh start.sh
    ```

## Running API service

### From Docker Compose

Recommended for automatic management of database, build-args, ports, etc.

```bash
docker-compose up --build app
```

### From Docker

```bash
docker build ./app -t basic-api --build-arg GOOGLE_API_KEY=$GOOGLE_API_KEY --build-arg PORT=8080
docker run --rm -p 8080:8080 basic-api
```

### From Local Environment

```bash
go get
go run main.go
```

## Database Migrations

Migration are run automatically to latest version on app start.

To reset / drop database for testing on a clean slate, run

```bash
docker-compose down
docker volume rm basic-api-server_api-db
```

## Testing

### Manual Tests

First, make sure your service is running using instructions above.

#### Create Order

```bash
curl -X POST http://localhost:8080/orders -d '{"origin": ["22.278", "114.185"], "destination": ["22.31", "114.216"]}'
```