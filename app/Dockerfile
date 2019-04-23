# Build stage compiles golang binary
FROM golang:1.12 as builder
ARG GOOGLE_API_KEY

WORKDIR /basic-api-server
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

ENV GOOGLE_API_KEY ${GOOGLE_API_KEY}

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o serve

# Create certs needed to query external APIs
FROM alpine:latest as certs
RUN apk --no-cache add ca-certificates

# Final stage prepares minimized image for serving
FROM scratch
ARG GOOGLE_API_KEY
ARG PORT

WORKDIR /app/
COPY --from=builder /basic-api-server/serve .
COPY --from=builder /basic-api-server/migrations ./migrations
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV GIN_MODE "release"
ENV GOOGLE_API_KEY ${GOOGLE_API_KEY}
ENV PORT ${PORT}

EXPOSE ${PORT}

CMD ["./serve"]