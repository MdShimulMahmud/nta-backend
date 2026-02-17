#!/bin/bash

# setup.sh - Setup script for NTA Backend

set -e

echo "================================"
echo "NTA Backend Setup Script"
echo "================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "This script requires root privileges. Please run with sudo."
  exit 1
fi

echo "[1/6] Checking system requirements..."

# Check for required commands
REQUIRED_COMMANDS="ip ifconfig tcpdump systemctl netplan"
MISSING_COMMANDS=""

for cmd in $REQUIRED_COMMANDS; do
  if ! command -v $cmd &> /dev/null; then
    MISSING_COMMANDS="$MISSING_COMMANDS $cmd"
  fi
done

if [ ! -z "$MISSING_COMMANDS" ]; then
  echo "Missing required commands:$MISSING_COMMANDS"
  echo "Please install them before continuing."
  exit 1
fi

echo "✓ All required commands are available"

echo ""
echo "[2/6] Checking Go installation..."

if ! command -v go &> /dev/null; then
  echo "Go is not installed. Installing Go 1.21..."
  wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
  rm go1.21.0.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
  echo "✓ Go installed successfully"
else
  echo "✓ Go is already installed ($(go version))"
fi

echo ""
echo "[3/6] Installing dependencies..."
make install

echo ""
echo "[4/6] Building application..."
make build

echo ""
echo "[5/6] Setting up environment configuration..."

if [ ! -f .env ]; then
  cp .env.example .env
  echo "✓ Created .env file from template"
  echo "  Please edit .env to configure your environment"
else
  echo "✓ .env file already exists"
fi

echo ""
echo "[6/6] Creating systemd service..."

cat > /etc/systemd/system/nta-backend.service << EOF
[Unit]
Description=NTA Backend Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/bin/nta-backend
Restart=on-failure
RestartSec=10
EnvironmentFile=$(pwd)/.env

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
echo "✓ Systemd service created"

echo ""
echo "================================"
echo "Setup Complete!"
echo "================================"
echo ""
echo "Next steps:"
echo "1. Edit .env file to configure your settings:"
echo "   nano .env"
echo ""
echo "2. Enable and start the service:"
echo "   sudo systemctl enable nta-backend"
echo "   sudo systemctl start nta-backend"
echo ""
echo "3. Check service status:"
echo "   sudo systemctl status nta-backend"
echo ""
echo "4. View logs:"
echo "   sudo journalctl -u nta-backend -f"
echo ""
echo "5. Access API documentation:"
echo "   http://localhost:8080/swagger/index.html"
echo ""
echo "For manual run:"
echo "   sudo ./bin/nta-backend"
echo ""
