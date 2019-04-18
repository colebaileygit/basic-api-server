# Build stage compiles golang binary
FROM golang:1.12 as builder

WORKDIR /basic-api-server
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o serve

# Final stage prepares minimized image for serving
FROM alpine:latest

WORKDIR /app/
COPY --from=builder /basic-api-server/serve .

EXPOSE 8080

CMD ["./serve"]