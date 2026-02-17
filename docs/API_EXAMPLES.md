# API Usage Examples

## Network Interfaces

### List all interfaces

```bash
curl -X GET http://localhost:8080/api/v1/interfaces
```

### Get specific interface details

```bash
curl -X GET http://localhost:8080/api/v1/interfaces/enp3s0
```

### Bring interface up

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/enp3s0/state \
  -H "Content-Type: application/json" \
  -d '{"state": "up"}'
```

### Bring interface down

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/enp3s0/state \
  -H "Content-Type: application/json" \
  -d '{"state": "down"}'
```

### Enable promiscuous mode

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/enp3s0/promiscuous \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### Disable promiscuous mode

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/enp3s0/promiscuous \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

### Check network logs with tcpdump

```bash
curl -X GET http://localhost:8080/api/v1/interfaces/enp3s0/check-network-logs
```

### Check firewall logs

```bash
curl -X GET http://localhost:8080/api/v1/interfaces/enp2s0/check-firewall-logs
```

### Get Netplan configuration

```bash
curl -X GET http://localhost:8080/api/v1/interfaces/netplan
```

### Update Netplan configuration

```bash
curl -X PUT http://localhost:8080/api/v1/interfaces/netplan \
  -H "Content-Type: application/json" \
  -d '{
    "interface": "enp2s0",
    "addresses": ["10.10.4.22/23"],
    "dhcp4": false,
    "dhcp6": false,
    "gateway4": "10.10.5.254",
    "nameservers": {
      "addresses": ["10.1.2.16"],
      "search": ["amplifysec.net"]
    },
    "routes": [
      {
        "to": "default",
        "via": "10.10.5.254"
      }
    ]
  }'
```

## System Services

### Get service status

```bash
curl -X GET http://localhost:8080/api/v1/services/logstash/status
curl -X GET http://localhost:8080/api/v1/services/suricata/status
curl -X GET http://localhost:8080/api/v1/services/rport/status
curl -X GET http://localhost:8080/api/v1/services/rsyslog/status
curl -X GET http://localhost:8080/api/v1/services/wazuh-agent/status
```

### Start a service

```bash
curl -X POST http://localhost:8080/api/v1/services/suricata/action \
  -H "Content-Type: application/json" \
  -d '{"action": "start"}'
```

### Stop a service

```bash
curl -X POST http://localhost:8080/api/v1/services/suricata/action \
  -H "Content-Type: application/json" \
  -d '{"action": "stop"}'
```

### Restart a service

```bash
curl -X POST http://localhost:8080/api/v1/services/suricata/action \
  -H "Content-Type: application/json" \
  -d '{"action": "restart"}'
```

## Logstash

### Check Logstash configuration

```bash
curl -X GET http://localhost:8080/api/v1/services/logstash/config/check
```

### Get Logstash pipelines

```bash
curl -X GET http://localhost:8080/api/v1/services/logstash/pipelines
```

### Enable a pipeline

```bash
curl -X POST http://localhost:8080/api/v1/services/logstash/pipelines/toggle \
  -H "Content-Type: application/json" \
  -d '{
    "pipeline_id": "network",
    "enabled": true
  }'
```

### Disable a pipeline

```bash
curl -X POST http://localhost:8080/api/v1/services/logstash/pipelines/toggle \
  -H "Content-Type: application/json" \
  -d '{
    "pipeline_id": "edr",
    "enabled": false
  }'
```

### Get Logstash logs (last 100 lines)

```bash
curl -X GET "http://localhost:8080/api/v1/services/logstash/logs?lines=100"
```

## Suricata

### Check Suricata configuration

```bash
curl -X GET http://localhost:8080/api/v1/services/suricata/config/check
```

### Get current monitoring interface

```bash
curl -X GET http://localhost:8080/api/v1/services/suricata/interface
```

### Update monitoring interface

```bash
curl -X PUT http://localhost:8080/api/v1/services/suricata/interface \
  -H "Content-Type: application/json" \
  -d '{"interface": "enp3s0"}'
```

### Get HOME_NET configuration

```bash
curl -X GET http://localhost:8080/api/v1/services/suricata/home-net
```

### Update HOME_NET configuration

```bash
curl -X PUT http://localhost:8080/api/v1/services/suricata/home-net \
  -H "Content-Type: application/json" \
  -d '{
    "networks": ["192.168.1.0/24", "10.0.0.0/8"]
  }'
```

### Update Suricata rules

```bash
curl -X POST http://localhost:8080/api/v1/services/suricata/rules/update
```

### Get Suricata alerts (last 50)

```bash
curl -X GET "http://localhost:8080/api/v1/services/suricata/alerts?limit=50"
```

### Get Suricata logs

```bash
curl -X GET "http://localhost:8080/api/v1/services/suricata/logs?lines=50"
```

## RPort

### Check RPort connection status

```bash
curl -X GET http://localhost:8080/api/v1/services/rport/connection
```

### Get RPort logs

```bash
curl -X GET "http://localhost:8080/api/v1/services/rport/logs?lines=50"
```

## Rsyslog

### Check Rsyslog configuration status

```bash
curl -X GET http://localhost:8080/api/v1/services/rsyslog/config/check
```

### Enable firewall log collection

```bash
curl -X POST http://localhost:8080/api/v1/services/rsyslog/config/toggle \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### Disable firewall log collection

```bash
curl -X POST http://localhost:8080/api/v1/services/rsyslog/config/toggle \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

### Get firewall logs

```bash
curl -X GET "http://localhost:8080/api/v1/services/rsyslog/firewall-logs?lines=100"
```

### Get Rsyslog service logs

```bash
curl -X GET "http://localhost:8080/api/v1/services/rsyslog/logs?lines=50"
```

## Wazuh

### Get Wazuh agent logs

```bash
curl -X GET "http://localhost:8080/api/v1/services/wazuh/logs?lines=50"
```

## Error Responses

All endpoints return JSON responses. Error responses follow this format:

```json
{
  "error": "Error message description"
}
```

Success responses include:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully"
}
```
