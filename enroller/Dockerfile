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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o /enroller

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install required dependencies
RUN apk add --no-cache \
    curl \
    wget \
    ca-certificates \
    tzdata \

# Install fabric-ca-client
RUN wget https://github.com/hyperledger/fabric-ca/releases/download/v1.5.7/hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz && \
    tar xzf hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz -C /usr/local/bin --strip-components=1 && \
    rm hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz

# Create directory for crypto material
RUN mkdir -p /etc/tls
RUN mkdir -p /etc/rcaadmin/msp/

# Copy binary from builder
COPY --from=builder /enroller /app/enroller

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose REST API port
EXPOSE 8080

# Default environment variables
ENV PORT=8080 \
    ORG=org1 \
    MSP_ID=Org1MSP

# Mount these at runtime:
# - /etc/tls
# - /etc/rcaadmin/msp/

ENTRYPOINT ["/app/enroller"]