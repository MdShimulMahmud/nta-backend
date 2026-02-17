# Quick Start Guide

## Prerequisites

- Go 1.21+ installed
- Linux system (Ubuntu/Debian preferred)
- sudo/root access
- System tools: `ip`, `tcpdump`, `systemctl`, `netplan`

## 5-Minute Setup

### 1. Clone and Navigate

```bash
git clone https://github.com/sharetrip/nta-backend.git
cd nta-backend
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit if needed: nano .env
```

### 3. Install Dependencies and Build

```bash
make install
make build
```

### 4. Run the Application

```bash
sudo ./bin/nta-backend
```

The API will be available at: `http://localhost:8080`

## Quick Test

### Check Health

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "healthy",
  "message": "NTA Backend is running"
}
```

### List Network Interfaces

```bash
curl http://localhost:8080/api/v1/interfaces
```

### Get Service Status

```bash
curl http://localhost:8080/api/v1/services/logstash/status
```

## API Documentation

View full API documentation at:

```
http://localhost:8080/swagger/index.html
```

## Common Commands

### Development Mode (with hot reload)

```bash
sudo make dev
```

### Run Tests

```bash
make test
```

### Build Only

```bash
make build
```

### Clean Build Artifacts

```bash
make clean
```

### Format Code

```bash
make format
```

### Generate API Docs

```bash
make swagger
```

## Docker Setup

### Using Docker Compose

```bash
docker-compose up -d
```

### Manual Docker Build and Run

```bash
docker build -t nta-backend .
docker run -p 8080:8080 --privileged nta-backend
```

## Production Setup

### 1. Run Setup Script

```bash
sudo bash scripts/setup.sh
```

### 2. Enable Systemd Service

```bash
sudo systemctl enable nta-backend
sudo systemctl start nta-backend
```

### 3. Check Status

```bash
sudo systemctl status nta-backend
```

### 4. View Logs

```bash
sudo journalctl -u nta-backend -f
```

## Example API Calls

### Bring Interface Up

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/eth0/state \
  -H "Content-Type: application/json" \
  -d '{"state": "up"}'
```

### Enable Promiscuous Mode

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/eth0/promiscuous \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### Restart Suricata Service

```bash
curl -X POST http://localhost:8080/api/v1/services/suricata/action \
  -H "Content-Type: application/json" \
  -d '{"action": "restart"}'
```

### Get Suricata Alerts

```bash
curl "http://localhost:8080/api/v1/services/suricata/alerts?limit=10"
```

## Troubleshooting

### Permission Errors

Always run with sudo:

```bash
sudo ./bin/nta-backend
```

### Port Already in Use

Change port in `.env`:

```env
APP_PORT=8081
```

### Build Errors

```bash
make clean
make install
make build
```

## Next Steps

1. **Read full documentation**: Check [README.md](../README.md)
2. **Explore API**: Visit swagger docs at `/swagger/index.html`
3. **Development guide**: See [docs/DEVELOPMENT.md](DEVELOPMENT.md)
4. **API examples**: Check [docs/API_EXAMPLES.md](API_EXAMPLES.md)

## Support

- GitHub Issues: [Report a problem](https://github.com/sharetrip/nta-backend/issues)
- Email: support@amplifysec.net
