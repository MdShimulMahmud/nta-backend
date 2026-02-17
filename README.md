# NTA Backend

Network Threat Analysis (NTA) Backend - A production-ready Golang backend for managing network interfaces and system services including Logstash, Suricata, RPort, Rsyslog, and Wazuh-Agent.

## Features

### Network Interface Management

- ✅ List all available network interfaces
- ✅ Get detailed interface information (state, MAC, IP, speed, statistics)
- ✅ Bring interfaces up/down
- ✅ Enable/disable promiscuous mode
- ✅ Check network and firewall logs with tcpdump
- ✅ Manage Netplan configuration
- ✅ Update interface IP addresses and routing

### System Service Management

#### Logstash

- ✅ Service control (start, stop, restart)
- ✅ Configuration syntax check
- ✅ List running pipelines
- ✅ Enable/disable pipelines (network, edr, asset-discovery, identities-access)
- ✅ View service logs

#### Suricata (NTA)

- ✅ Service control (start, stop, restart, reload)
- ✅ Configuration syntax check
- ✅ View recent alerts
- ✅ Check running interface
- ✅ Change monitoring interface
- ✅ Set HOME_NET configuration
- ✅ Update rules (suricata-update)
- ✅ View service logs

#### OpenRPort

- ✅ Service control (start, stop, restart)
- ✅ Check connection status to remote server
- ✅ View service logs

#### Rsyslog

- ✅ Service control (start, stop, restart)
- ✅ Enable/disable firewall log collection
- ✅ Check listening port status
- ✅ View firewall logs
- ✅ View service logs

#### Wazuh-Agent

- ✅ Service control (start, stop, restart)
- ✅ View agent logs

## Architecture

This project follows a **Feature-Based Structure** with clean architecture principles:

```
nta-backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/            # Private application code
│   ├── api/            # API routes
│   ├── config/         # Configuration management
│   ├── interfaces/     # Network interface management
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   ├── services/       # System service management
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   ├── middleware/     # HTTP middlewares
│   ├── logger/         # Logging utilities
│   └── utils/          # Helper utilities
├── docs/               # API documentation (Swagger)
├── scripts/            # Helper scripts
├── .env.example        # Environment configuration template
├── Dockerfile          # Docker image definition
├── docker-compose.yml  # Docker Compose configuration
├── Makefile           # Build and development tasks
└── go.mod             # Go module dependencies
```

## Requirements

- Go 1.21+
- Linux system with systemd
- Root/sudo privileges (for network and service management)
- Required system tools:
  - `ip`, `ifconfig`, `ethtool`
  - `tcpdump`
  - `systemctl`
  - `netplan` (for Ubuntu/Debian)

## Installation

### Quick Start

1. **Clone the repository:**

```bash
git clone https://github.com/sharetrip/nta-backend.git
cd nta-backend
```

2. **Copy environment configuration:**

```bash
cp .env.example .env
```

3. **Edit configuration (optional):**

```bash
nano .env
```

4. **Install dependencies:**

```bash
make install
```

5. **Build the application:**

```bash
make build
```

6. **Run the application:**

```bash
sudo ./bin/nta-backend
```

Or run directly:

```bash
sudo make run
```

### Docker Deployment

**Build and run with Docker:**

```bash
make docker-build
make docker-run
```

**Or use Docker Compose:**

```bash
docker-compose up -d
```

## Configuration

All configuration is managed through environment variables. See [.env.example](.env.example) for available options:

```env
# Application
APP_NAME=NTA Backend
APP_ENV=development
APP_PORT=8080
APP_HOST=0.0.0.0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# API
API_PREFIX=/api/v1
API_TIMEOUT=30s

# Network
DEFAULT_NETWORK_INTERFACE=enp3s0
NETPLAN_CONFIG_PATH=/etc/netplan/00-installer-config.yaml

# Services
LOGSTASH_CONFIG_PATH=/etc/logstash
SURICATA_CONFIG_PATH=/etc/suricata/suricata.yaml
RSYSLOG_CONFIG_PATH=/etc/rsyslog.d
WAZUH_CONFIG_PATH=/var/ossec
RPORT_CONFIG_PATH=/etc/rport
```

## API Documentation

Once the server is running, access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

### Main Endpoints

#### Network Interfaces

- `GET /api/v1/interfaces` - List all interfaces
- `GET /api/v1/interfaces/:name` - Get interface details
- `PUT /api/v1/interfaces/:name/state` - Update interface state (up/down)
- `PUT /api/v1/interfaces/:name/promiscuous` - Toggle promiscuous mode
- `GET /api/v1/interfaces/:name/check-network-logs` - Check network logs
- `GET /api/v1/interfaces/:name/check-firewall-logs` - Check firewall logs
- `GET /api/v1/interfaces/netplan` - Get Netplan configuration
- `PUT /api/v1/interfaces/netplan` - Update Netplan configuration

#### System Services

- `GET /api/v1/services/:name/status` - Get service status
- `POST /api/v1/services/:name/action` - Perform action (start/stop/restart)

#### Logstash

- `GET /api/v1/services/logstash/config/check` - Check configuration
- `GET /api/v1/services/logstash/pipelines` - List pipelines
- `POST /api/v1/services/logstash/pipelines/toggle` - Toggle pipeline
- `GET /api/v1/services/logstash/logs` - Get logs

#### Suricata

- `GET /api/v1/services/suricata/config/check` - Check configuration
- `GET /api/v1/services/suricata/interface` - Get monitoring interface
- `PUT /api/v1/services/suricata/interface` - Update monitoring interface
- `GET /api/v1/services/suricata/home-net` - Get HOME_NET
- `PUT /api/v1/services/suricata/home-net` - Update HOME_NET
- `POST /api/v1/services/suricata/rules/update` - Update rules
- `GET /api/v1/services/suricata/alerts` - Get recent alerts
- `GET /api/v1/services/suricata/logs` - Get logs

#### RPort

- `GET /api/v1/services/rport/connection` - Check connection status
- `GET /api/v1/services/rport/logs` - Get logs

#### Rsyslog

- `GET /api/v1/services/rsyslog/config/check` - Check configuration
- `POST /api/v1/services/rsyslog/config/toggle` - Enable/disable firewall logs
- `GET /api/v1/services/rsyslog/firewall-logs` - Get firewall logs
- `GET /api/v1/services/rsyslog/logs` - Get service logs

#### Wazuh

- `GET /api/v1/services/wazuh/logs` - Get agent logs

## Development

### Running in Development Mode

With hot reload (requires [air](https://github.com/air-verse/air)):

```bash
make dev
```

### Running Tests

```bash
make test
```

### View Test Coverage

```bash
make coverage
```

### Code Formatting

```bash
make format
```

### Linting

```bash
make lint
```

### Generate Swagger Documentation

```bash
make swagger
```

## Production Deployment

### Systemd Service

Create a systemd service file `/etc/systemd/system/nta-backend.service`:

```ini
[Unit]
Description=NTA Backend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/nta-backend
ExecStart=/opt/nta-backend/nta-backend
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable nta-backend
sudo systemctl start nta-backend
sudo systemctl status nta-backend
```

### Security Considerations

1. **Run with appropriate privileges**: The application requires root/sudo access for network and service management.
2. **Use HTTPS**: In production, place behind a reverse proxy (nginx/Apache) with SSL/TLS.
3. **Firewall**: Restrict access to the API port (8080) to authorized IPs only.
4. **Authentication**: Implement JWT or API key authentication (extend the middleware).
5. **Rate Limiting**: Enabled by default in production mode.

## Troubleshooting

### Permission Denied Errors

Make sure to run the application with sudo:

```bash
sudo ./bin/nta-backend
```

### Service Not Found

Ensure the service is installed on your system:

```bash
systemctl list-unit-files | grep <service-name>
```

### Network Interface Not Found

List available interfaces:

```bash
ip link show
```

### Port Already in Use

Change the port in `.env`:

```env
APP_PORT=8081
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:

- Create an issue on GitHub
- Email: support@amplifysec.net

## Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- Logging with [Logrus](https://github.com/sirupsen/logrus)
- Configuration with [Viper](https://github.com/spf13/viper)
- API documentation with [Swagger](https://swagger.io/)
