# NTA Backend - Project Summary

## Overview

This is a **production-ready Golang backend** for Network Threat Analysis (NTA) system that manages:

- Network interfaces (state, promiscuous mode, IP configuration)
- System services (Logstash, Suricata, RPort, Rsyslog, Wazuh-Agent)

## Architecture: Feature-Based Structure

```
nta-backend/
├── cmd/server/                 # Application entry point
├── internal/
│   ├── api/                   # API routing
│   ├── config/                # Configuration management
│   ├── interfaces/            # Network interface management feature
│   │   ├── models.go          # Data structures
│   │   ├── repository.go      # System operations
│   │   ├── service.go         # Business logic
│   │   └── handler.go         # HTTP handlers
│   ├── services/              # System service management feature
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   ├── middleware/            # HTTP middlewares
│   ├── logger/                # Logging utilities
│   └── utils/                 # Helper functions
├── docs/                      # Documentation
├── scripts/                   # Setup and utility scripts
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose setup
├── Makefile                   # Build automation
└── .env.example              # Environment configuration template
```

## Core Features Implemented

### 1. Network Interface Management (`internal/interfaces/`)

✅ **List Interfaces** - Get all network interfaces with details
✅ **Interface Details** - State, MAC, IP, speed, duplex, statistics
✅ **State Control** - Bring interfaces up/down
✅ **Promiscuous Mode** - Enable/disable for packet capture
✅ **TCPDump Integration** - Check network and firewall logs
✅ **Netplan Management** - Read and update network configuration
✅ **IP Configuration** - Update addresses, routes, nameservers

**API Endpoints:**

- `GET /api/v1/interfaces` - List all interfaces
- `GET /api/v1/interfaces/:name` - Get interface details
- `PUT /api/v1/interfaces/:name/state` - Update state (up/down)
- `PUT /api/v1/interfaces/:name/promiscuous` - Toggle promiscuous mode
- `GET /api/v1/interfaces/:name/check-network-logs` - Check network logs
- `GET /api/v1/interfaces/:name/check-firewall-logs` - Check firewall logs
- `GET /api/v1/interfaces/netplan` - Get Netplan config
- `PUT /api/v1/interfaces/netplan` - Update Netplan config

### 2. System Service Management (`internal/services/`)

#### Logstash

✅ Service control (start, stop, restart)
✅ Configuration syntax check
✅ Pipeline management (list, enable, disable)
✅ Supported pipelines: network, edr, asset-discovery, identities-access
✅ Log viewing

#### Suricata (NTA)

✅ Service control (start, stop, restart, reload)
✅ Configuration syntax check
✅ Interface management (get, update monitoring interface)
✅ HOME_NET configuration
✅ Rule updates (suricata-update)
✅ Alert viewing
✅ Log viewing

#### OpenRPort

✅ Service control
✅ Connection status check
✅ Log viewing

#### Rsyslog

✅ Service control
✅ Firewall log collection toggle
✅ Port listening check
✅ Firewall log viewing
✅ Service log viewing

#### Wazuh-Agent

✅ Service control
✅ Log viewing

### 3. Production-Ready Features

✅ **Configuration Management** - Environment-based configuration
✅ **Logging** - Structured logging with Logrus
✅ **Error Handling** - Comprehensive error handling and recovery
✅ **CORS** - Configurable CORS support
✅ **Rate Limiting** - Request rate limiting
✅ **Middleware** - Request logging, error handling
✅ **Health Checks** - Health endpoint for monitoring
✅ **Graceful Shutdown** - Proper shutdown handling
✅ **Docker Support** - Dockerfile and docker-compose
✅ **Swagger Documentation** - Auto-generated API docs
✅ **Makefile** - Build automation

## Technology Stack

- **Framework**: Gin Web Framework
- **Logging**: Logrus
- **Configuration**: Viper + godotenv
- **API Documentation**: Swagger/OpenAPI
- **Architecture**: Clean Architecture with Feature-Based Structure
- **Language**: Go 1.21+

## Key Files

### Configuration

- `.env.example` - Environment configuration template
- `internal/config/config.go` - Configuration loader

### Main Application

- `cmd/server/main.go` - Application entry point with dependency injection

### API Layer

- `internal/api/router.go` - API route definitions

### Feature Modules

Each feature follows the same pattern:

1. `models.go` - Data structures
2. `repository.go` - System/data access
3. `service.go` - Business logic
4. `handler.go` - HTTP handlers

### Infrastructure

- `internal/middleware/` - HTTP middleware (logging, errors, CORS, rate limiting)
- `internal/logger/` - Logging setup
- `internal/utils/` - Helper functions

### Build & Deploy

- `Makefile` - Build commands (install, build, run, test, lint, etc.)
- `Dockerfile` - Multi-stage Docker build
- `docker-compose.yml` - Docker Compose configuration
- `scripts/setup.sh` - Production setup script
- `scripts/health_check.sh` - Health check script

### Documentation

- `README.md` - Comprehensive project documentation
- `docs/QUICKSTART.md` - Quick start guide
- `docs/DEVELOPMENT.md` - Development guide
- `docs/API_EXAMPLES.md` - API usage examples

## Getting Started

### Quick Start (Development)

```bash
# 1. Install dependencies
make install

# 2. Build
make build

# 3. Run
sudo ./bin/nta-backend
```

### Production Setup

```bash
sudo bash scripts/setup.sh
sudo systemctl enable nta-backend
sudo systemctl start nta-backend
```

### Docker

```bash
docker-compose up -d
```

## API Access

- **Base URL**: `http://localhost:8080`
- **API Prefix**: `/api/v1`
- **Health Check**: `http://localhost:8080/health`
- **Swagger Docs**: `http://localhost:8080/swagger/index.html`

## Security Features

✅ Requires root/sudo privileges for system operations
✅ Rate limiting enabled in production
✅ CORS configuration
✅ Error recovery middleware
✅ Request logging
✅ Configurable via environment variables
⚠️ Note: Add authentication middleware for production use

## Makefile Commands

```bash
make help          # Show available commands
make install       # Install dependencies
make build         # Build application
make run           # Run application
make dev           # Run with hot reload (air)
make test          # Run tests
make coverage      # Generate coverage report
make lint          # Run linter
make format        # Format code
make swagger       # Generate API docs
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

## Environment Configuration

Key environment variables (see `.env.example`):

- `APP_PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `NETPLAN_CONFIG_PATH` - Path to netplan config
- `SURICATA_CONFIG_PATH` - Path to Suricata config
- `LOGSTASH_CONFIG_PATH` - Path to Logstash config

## Testing

### Health Check

```bash
curl http://localhost:8080/health
```

### List Interfaces

```bash
curl http://localhost:8080/api/v1/interfaces
```

### Get Service Status

```bash
curl http://localhost:8080/api/v1/services/suricata/status
```

## Next Steps

1. **Generate Swagger docs**: `make swagger`
2. **Run tests**: `make test`
3. **Deploy**: Use `scripts/setup.sh` or Docker
4. **Customize**: Edit `.env` for your environment
5. **Extend**: Add authentication, more features following the same pattern

## Clean Code Principles

✅ **Separation of Concerns** - Each layer has single responsibility
✅ **Dependency Injection** - Dependencies passed explicitly
✅ **Interface Segregation** - Small, focused interfaces
✅ **Error Handling** - Proper error propagation and logging
✅ **Configuration** - External configuration via environment
✅ **Logging** - Structured logging throughout
✅ **Testing** - Testable architecture with mockable dependencies

## Production Readiness

✅ Graceful shutdown
✅ Health checks
✅ Structured logging
✅ Error recovery
✅ Rate limiting
✅ CORS support
✅ Docker support
✅ Systemd service file
✅ Setup scripts
✅ Comprehensive documentation
✅ API documentation (Swagger)

## Project Statistics

- **Language**: Go
- **Framework**: Gin
- **Architecture**: Feature-Based + Clean Architecture
- **Total Packages**: 8 internal packages
- **Features**: 2 main features (interfaces, services)
- **API Endpoints**: 30+ endpoints
- **Documentation**: 4 detailed docs files
- **Scripts**: 2 helper scripts

## Support

- **Documentation**: Comprehensive README and docs/
- **Examples**: Complete API examples in docs/API_EXAMPLES.md
- **Scripts**: Automated setup and health check scripts
- **Error Messages**: Descriptive error messages throughout

This is a **production-ready, well-architected Golang backend** following industry best practices and clean architecture principles! 🚀
