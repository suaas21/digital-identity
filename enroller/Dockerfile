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

# Install glibc compatibility layer
RUN apk add --no-cache --virtual .deps \
    wget \
    ca-certificates \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.35-r1/glibc-2.35-r1.apk \
    && apk add --no-cache glibc-2.35-r1.apk \
    && rm glibc-2.35-r1.apk \
    && apk del .deps

# Create required lib64 directory and setup links
RUN mkdir -p /lib64 \
    && ln -sf /usr/glibc-compat/lib/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2

# Install remaining dependencies
RUN apk add --no-cache \
    libstdc++ \
    tzdata

# Install fabric-ca-client
RUN wget https://github.com/hyperledger/fabric-ca/releases/download/v1.5.7/hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz \
    && tar xzf hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz -C /usr/local/bin --strip-components=1 \
    && rm hyperledger-fabric-ca-linux-amd64-1.5.7.tar.gz \
    && chmod +x /usr/local/bin/fabric-ca-client

# Verify installation
RUN fabric-ca-client version

# Create directory for crypto material
RUN mkdir -p /etc/tls
RUN mkdir -p /etc/rcaadmin/msp/

# Copy binary from builder
COPY --from=builder /enroller /app/enroller

# Create non-root user
#RUN adduser -D -g '' appuser
#USER appuser

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