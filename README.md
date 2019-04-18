# Basic API Server

## Setup

TODO

## Running API service

### From Docker

```bash
docker build . -t basic-api -f build/package/Dockerfile
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
curl -X POST http://localhost:8080/orders
```