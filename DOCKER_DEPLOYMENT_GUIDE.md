# Docker Deployment Guide - MCP-Netbird

This guide demonstrates how to deploy mcp-netbird with Docker using the new stateless configuration support feature.

## Quick Start

### 1. Build the Docker Image

```bash
docker build -t mcp-netbird-sse:latest -f Dockerfile.sse .
```

### 2. Run the Container

Choose one of three configuration methods:

#### Option A: Stateless Container (Recommended for Docker MCP Toolkit)

Start the container without any credentials:

```bash
docker run -d --name mcp-netbird -p 8001:8001 mcp-netbird-sse:latest
```

Credentials are provided via HTTP headers with each request. This is ideal for:
- Docker MCP Toolkit deployments
- Multi-tenant environments
- Environments where you don't want to bake credentials into containers

#### Option B: CLI Arguments

Start the container with credentials as CLI arguments:

```bash
docker run -d --name mcp-netbird -p 8001:8001 \
  mcp-netbird-sse:latest \
  --transport sse \
  --sse-address 0.0.0.0:8001 \
  --api-token "your-netbird-api-token" \
  --api-host "api.netbird.io"
```

#### Option C: Environment Variables (Backward Compatible)

Start the container with credentials as environment variables:

```bash
docker run -d --name mcp-netbird -p 8001:8001 \
  -e NETBIRD_API_TOKEN="your-netbird-api-token" \
  -e NETBIRD_HOST="api.netbird.io" \
  mcp-netbird-sse:latest
```

## Docker MCP Toolkit Integration

### Configuration File

Create or update your Docker MCP Toolkit configuration file:

```json
{
  "mcpServers": {
    "netbird": {
      "url": "http://localhost:8001",
      "transport": "sse",
      "headers": {
        "X-Netbird-API-Token": "your-netbird-api-token",
        "X-Netbird-Host": "api.netbird.io"
      }
    }
  }
}
```

### Start the Container

```bash
docker run -d --name mcp-netbird -p 8001:8001 mcp-netbird-sse:latest
```

The container runs without any credentials. Each request from the Docker MCP Toolkit will include the credentials via HTTP headers.

## Configuration Priority

When multiple configuration sources are present, the priority order is:

1. **CLI Arguments** (highest priority)
2. **HTTP Headers** (medium priority)
3. **Environment Variables** (lowest priority)

### Example: Priority in Action

```bash
# Start container with environment variables
docker run -d --name mcp-netbird -p 8001:8001 \
  -e NETBIRD_API_TOKEN="env-token" \
  -e NETBIRD_HOST="env.example.com" \
  mcp-netbird-sse:latest \
  --transport sse \
  --sse-address 0.0.0.0:8001 \
  --api-token "cli-token" \
  --api-host "cli.example.com"

# Result: CLI arguments take priority
# Token used: "cli-token"
# Host used: "cli.example.com"
```

## Testing the Deployment

### Test 1: Health Check

```bash
docker logs mcp-netbird
# Expected output: "SSE server listening on 0.0.0.0:8001"
```

### Test 2: Container Status

```bash
docker ps --filter "name=mcp-netbird"
# Should show container running on port 8001
```

### Test 3: Configuration Verification

For stateless containers, verify the server is running without credentials:

```bash
docker inspect mcp-netbird --format '{{range .Config.Env}}{{println .}}{{end}}' | grep NETBIRD
# Should show no NETBIRD_* environment variables
```

## Use Cases

### Use Case 1: Development Environment

Quick testing with CLI arguments:

```bash
docker run --rm -p 8001:8001 \
  mcp-netbird-sse:latest \
  --transport sse \
  --sse-address 0.0.0.0:8001 \
  --api-token "dev-token" \
  --api-host "api.netbird.io"
```

### Use Case 2: Production Deployment

Stateless container with credentials managed externally:

```bash
# Start container
docker run -d --name mcp-netbird -p 8001:8001 mcp-netbird-sse:latest

# Credentials provided by Docker MCP Toolkit via HTTP headers
# No credentials stored in container or environment
```

### Use Case 3: Multi-Tenant Environment

Single container serving multiple tenants:

```bash
# Start one container
docker run -d --name mcp-netbird -p 8001:8001 mcp-netbird-sse:latest

# Each tenant provides their own credentials via HTTP headers
# Tenant 1: X-Netbird-API-Token: tenant1-token
# Tenant 2: X-Netbird-API-Token: tenant2-token
```

### Use Case 4: Backward Compatible Deployment

Traditional deployment with environment variables:

```bash
docker run -d --name mcp-netbird -p 8001:8001 \
  -e NETBIRD_API_TOKEN="production-token" \
  -e NETBIRD_HOST="api.netbird.io" \
  mcp-netbird-sse:latest
```

## Docker Compose

### Stateless Configuration

```yaml
version: '3.8'

services:
  mcp-netbird:
    image: mcp-netbird-sse:latest
    ports:
      - "8001:8001"
    restart: unless-stopped
    # No environment variables - credentials via HTTP headers
```

### With Environment Variables

```yaml
version: '3.8'

services:
  mcp-netbird:
    image: mcp-netbird-sse:latest
    ports:
      - "8001:8001"
    environment:
      - NETBIRD_API_TOKEN=${NETBIRD_API_TOKEN}
      - NETBIRD_HOST=${NETBIRD_HOST:-api.netbird.io}
    restart: unless-stopped
```

### With CLI Arguments

```yaml
version: '3.8'

services:
  mcp-netbird:
    image: mcp-netbird-sse:latest
    command:
      - --transport
      - sse
      - --sse-address
      - 0.0.0.0:8001
      - --api-token
      - ${NETBIRD_API_TOKEN}
      - --api-host
      - ${NETBIRD_HOST:-api.netbird.io}
    ports:
      - "8001:8001"
    restart: unless-stopped
```

## Troubleshooting

### Container Won't Start

Check the logs:

```bash
docker logs mcp-netbird
```

Common issues:
- Port 8001 already in use: Change the port mapping `-p 8002:8001`
- Invalid CLI arguments: Check the command syntax

### Configuration Not Working

Verify the configuration priority:

```bash
# Check environment variables
docker inspect mcp-netbird --format '{{range .Config.Env}}{{println .}}{{end}}'

# Check command arguments
docker inspect mcp-netbird --format '{{.Config.Cmd}}'
```

### HTTP Headers Not Being Used

Ensure you're running in SSE mode and the container is stateless:

```bash
# Container should not have NETBIRD_* environment variables
docker inspect mcp-netbird --format '{{range .Config.Env}}{{println .}}{{end}}' | grep NETBIRD

# Should return nothing for stateless containers
```

## Security Best Practices

1. **Use Stateless Containers**: Don't bake credentials into container images
2. **Use HTTP Headers**: Provide credentials per-request via HTTP headers
3. **Rotate Credentials**: Easy to rotate when using HTTP headers
4. **Least Privilege**: Use API tokens with minimal required permissions
5. **Network Isolation**: Run containers in isolated networks
6. **TLS/HTTPS**: Use reverse proxy with TLS for production

## Performance Considerations

- **Stateless containers**: No performance overhead, credentials loaded per-request
- **CLI arguments**: Configuration loaded once at startup
- **Environment variables**: Configuration loaded once at startup
- **HTTP headers**: Minimal overhead, only in SSE mode

## Cleanup

Stop and remove the container:

```bash
docker stop mcp-netbird
docker rm mcp-netbird
```

Remove the image:

```bash
docker rmi mcp-netbird-sse:latest
```

## Next Steps

1. Deploy to your Docker MCP Toolkit environment
2. Configure your MCP client with the appropriate configuration method
3. Test with your Netbird API token
4. Monitor logs for any issues

For more information, see the main [README.md](README.md) file.
