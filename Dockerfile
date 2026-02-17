
# --- Build Stage ---
FROM golang:1.23-alpine AS builder
LABEL stage=builder

RUN apk add --no-cache git make gcc musl-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nta-backend ./cmd/server
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init --generalInfo cmd/server/main.go --output docs/swagger

# --- Runtime Stage ---
FROM alpine:latest
LABEL maintainer="AmplifySec <support@amplifysec.net>"
LABEL org.opencontainers.image.source="https://github.com/sharetrip/nta-backend"
LABEL org.opencontainers.image.description="Production-ready NTA Backend API"

# Install only required runtime dependencies
RUN apk --no-cache add ca-certificates iproute2 tcpdump ethtool curl jq

# Create non-root user (commented out because privileged mode is required)
# RUN addgroup -S nta && adduser -S nta -G nta
# USER nta

WORKDIR /app
COPY --from=builder /app/nta-backend .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Use environment variables for config (no .env file in image)

# Entrypoint
CMD ["./nta-backend"]
