# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o /gateway

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for TLS connections
RUN apk --no-cache add ca-certificates tzdata

# Create directory for crypto material
RUN mkdir -p /etc/secret-volume

# Copy binary from builder
COPY --from=builder /gateway /app/gateway

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose REST API port
EXPOSE 8080

# Default environment variables
ENV PORT=8080 \
    CHANNEL_NAME=mychannel \
    CHAINCODE_NAME=asset-transfer \
    MSP_ID=Org1MSP \
    PEER_ENDPOINT=test-network-org1-peer1-peer.localho.st:443 \
    PEER_HOST_ALIAS=test-network-org1-peer1-peer.localho.st

# Mount these at runtime:
# - /etc/secret-volume/certPath
# - /etc/secret-volume/keyPath
# - /etc/secret-volume/tlsCertPath

ENTRYPOINT ["/app/gateway"]