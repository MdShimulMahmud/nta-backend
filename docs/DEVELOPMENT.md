# Development Guide

## Setting Up Development Environment

### Prerequisites

- Go 1.21 or higher
- Linux system with systemd
- Root/sudo privileges
- Git
- Make

### Initial Setup

1. Clone the repository:

```bash
git clone https://github.com/sharetrip/nta-backend.git
cd nta-backend
```

2. Install dependencies:

```bash
make install
```

3. Copy environment configuration:

```bash
cp .env.example .env
```

4. Run the application:

```bash
sudo make run
```

## Project Structure

```
nta-backend/
├── cmd/
│   └── server/              # Main application entry point
│       └── main.go
├── internal/                # Private application code
│   ├── api/                # API routing
│   │   └── router.go
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── interfaces/         # Network interface feature
│   │   ├── models.go       # Data models
│   │   ├── repository.go   # Data access layer
│   │   ├── service.go      # Business logic
│   │   └── handler.go      # HTTP handlers
│   ├── services/           # System service feature
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   ├── middleware/         # HTTP middlewares
│   │   ├── logger.go
│   │   ├── error.go
│   │   ├── cors.go
│   │   └── rate_limiter.go
│   ├── logger/             # Logging utilities
│   │   └── logger.go
│   └── utils/              # Helper utilities
│       ├── response.go
│       └── system.go
├── docs/                   # Documentation
├── scripts/                # Helper scripts
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── go.sum
```

## Architecture Principles

### Feature-Based Structure

Each feature is organized into its own package with the following layers:

1. **Models** (`models.go`): Data structures and types
2. **Repository** (`repository.go`): Data access and system operations
3. **Service** (`service.go`): Business logic
4. **Handler** (`handler.go`): HTTP request/response handling

### Dependency Flow

```
Handler -> Service -> Repository
           ↓
        Config, Logger
```

### Clean Architecture Benefits

- **Separation of Concerns**: Each layer has a specific responsibility
- **Testability**: Easy to unit test each layer independently
- **Maintainability**: Changes in one layer don't affect others
- **Scalability**: Easy to add new features following the same pattern

## Development Workflow

### Hot Reload Development

Install Air for hot reload:

```bash
go install github.com/air-verse/air@latest
```

Run with hot reload:

```bash
make dev
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test
go test -v ./internal/interfaces/...
```

### Code Quality

#### Formatting

```bash
make format
```

#### Linting

```bash
make lint
```

### Generating API Documentation

```bash
make swagger
```

Then access at: http://localhost:8080/swagger/index.html

## Adding New Features

### Example: Adding a New Service

1. **Create models** in `internal/newfeature/models.go`:

```go
package newfeature

type NewFeature struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

2. **Create repository** in `internal/newfeature/repository.go`:

```go
package newfeature

type Repository struct{}

func NewRepository() *Repository {
    return &Repository{}
}

func (r *Repository) GetData() (*NewFeature, error) {
    // Implementation
    return nil, nil
}
```

3. **Create service** in `internal/newfeature/service.go`:

```go
package newfeature

import "github.com/sharetrip/nta-backend/internal/config"

type Service struct {
    repo   *Repository
    config *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
    return &Service{repo: repo, config: cfg}
}

func (s *Service) GetFeature() (*NewFeature, error) {
    return s.repo.GetData()
}
```

4. **Create handler** in `internal/newfeature/handler.go`:

```go
package newfeature

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) GetFeature(c *gin.Context) {
    data, err := h.service.GetFeature()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
```

5. **Register routes** in `internal/api/router.go`:

```go
featureGroup := v1.Group("/newfeature")
{
    featureGroup.GET("", newFeatureHandler.GetFeature)
}
```

6. **Wire up in main** (`cmd/server/main.go`):

```go
newFeatureRepo := newfeature.NewRepository()
newFeatureService := newfeature.NewService(newFeatureRepo, cfg)
newFeatureHandler := newfeature.NewHandler(newFeatureService)
```

## Debugging

### Logging

Logs are configured via environment variables:

```env
LOG_LEVEL=debug  # debug, info, warn, error
LOG_FORMAT=json  # json or text
```

### Viewing Logs

```bash
# Application logs
sudo journalctl -u nta-backend -f

# Or if running directly
sudo ./bin/nta-backend
```

### Debug Mode

Set in `.env`:

```env
APP_ENV=development
```

## Testing

### Unit Tests

```bash
go test ./internal/interfaces/...
go test ./internal/services/...
```

### Integration Tests

```bash
# Run integration tests (requires sudo)
sudo go test -tags=integration ./tests/...
```

### Manual API Testing

Use the provided examples in `docs/API_EXAMPLES.md`

## Best Practices

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use meaningful variable and function names
- Write comments for exported functions
- Keep functions small and focused

### Error Handling

```go
if err != nil {
    logger.Errorf("Failed to perform action: %v", err)
    return fmt.Errorf("action failed: %w", err)
}
```

### Logging

```go
logger.Info("Starting operation")
logger.Debugf("Processing item: %s", item)
logger.Errorf("Operation failed: %v", err)
```

### Configuration

- Use environment variables for configuration
- Never hardcode sensitive data
- Provide sensible defaults

### API Response Format

```go
// Success
c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data": data,
})

// Error
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Error message",
})
```

## Common Issues

### Permission Denied

Always run with sudo for system operations:

```bash
sudo ./bin/nta-backend
```

### Port Already in Use

Change port in `.env`:

```env
APP_PORT=8081
```

### Service Not Found

Ensure services are installed:

```bash
systemctl list-unit-files | grep suricata
```

## Contributing

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes following the architecture
3. Write tests for new features
4. Run tests and linting: `make check`
5. Update documentation
6. Commit with clear messages
7. Push and create a Pull Request

## Resources

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [systemd Documentation](https://www.freedesktop.org/software/systemd/man/)
