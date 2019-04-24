# Basic API Server

## Quick Start

1. set your Google Maps API key via:

    ```bash
    export GOOGLE_API_KEY="<INSERT KEY HERE>"
    ```

2. Run unit + integration tests and then run service

    ```bash
    sh start.sh
    ```

3. Wait for app / db initialization. During start-up, app will wait for up to 30s for DB to be initialized.

4. Query service at `localhost:8080` e.g.

    ```bash
    curl -L -X GET \
    'http://localhost:8080/orders?page=0&limit=3' \
    -H 'Content-Type: application/json'
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
docker-compose down --volumes
```

## Testing

### Automated Tests

#### Integration + Unit Tests

```bash
docker-compose up tester
```

#### Unit Tests

```bash
docker-compose run tester go test ./... -short
```

Alternative, for quicker debugging if local environment is configured:

```bash
cd app
go test ./... -short
```

### Manual Tests

First, make sure your service is running using instructions above.

#### Create Order

```bash
curl -L -X POST \
  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"origin": ["22.278", "114.185"], "destination": ["22.31", "114.212"]}'
```

#### Take Order

```bash
curl -L -X PATCH \
  http://localhost:8080/orders/2 \
  -H 'Content-Type: application/json' \
  -d '{"status": "TAKEN"}'
```

#### Fetch Orders

```bash
curl -L -X GET \
  'http://localhost:8080/orders?page=0&limit=3' \
  -H 'Content-Type: application/json'
```