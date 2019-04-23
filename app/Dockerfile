# Build stage compiles golang binary
FROM golang:1.12 as builder
ARG GOOGLE_API_KEY

WORKDIR /basic-api-server
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go get -d -v ./...

COPY . .

ENV GOOGLE_API_KEY ${GOOGLE_API_KEY}

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o serve

# Final stage prepares minimized image for serving
FROM alpine:latest
ARG GOOGLE_API_KEY
ARG PORT

RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /basic-api-server/serve .
COPY --from=builder /basic-api-server/migrations ./migrations

ENV GIN_MODE "release"
ENV GOOGLE_API_KEY ${GOOGLE_API_KEY}
ENV PORT ${PORT}

EXPOSE ${PORT}

CMD ["./serve"]