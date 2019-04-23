# Basic API Server

## Setup

First, set your Google Maps API key via:

```
export GOOGLE_API_KEY="<INSERT KEY HERE>"
```

## Running API service

### From Docker Compose

Recommended for automatic management of database, build-args, ports, etc.

```bash
docker-compose up
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

## Testing

### Manual Tests

First, make sure your service is running using instructions above.

#### Create Order

```bash
curl -X POST http://localhost:8080/orders -d '{"origin": ["22.278", "114.185"], "destination": ["22.31", "114.216"]}'
```