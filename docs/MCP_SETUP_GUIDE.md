# MCP NetBird Server Setup Guide

Complete guide for setting up and using the NetBird MCP (Model Context Protocol) server with AI assistants like Kiro, Claude Desktop, and other MCP clients.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Deployment Options](#deployment-options)
  - [Remote SSE Server](#remote-sse-server)
  - [Local STDIO Mode](#local-stdio-mode)
  - [Docker MCP Gateway](#docker-mcp-gateway)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)

---

## Overview

The NetBird MCP server provides AI assistants with programmatic access to NetBird VPN infrastructure through the Model Context Protocol. It exposes 50+ tools for managing peers, groups, policies, networks, routes, DNS, and more.

### Key Features

- **Complete NetBird API Coverage**: Manage all aspects of your NetBird infrastructure
- **Multiple Transport Modes**: STDIO for local use, SSE for remote access
- **Secure Authentication**: API token-based authentication with NetBird
- **Real-time Operations**: Direct API calls for immediate results
- **MCP Standard Compliant**: Works with any MCP-compatible client

---

## Architecture

### Remote Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                Local Machine (Developer)                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  AI Assistant (Kiro/Claude)                          │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │ MCP Protocol (STDIO)                │
│                       ▼                                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Docker MCP Gateway                                  │  │
│  │  - Protocol Translation (STDIO ↔ SSE)               │  │
│  │  - Multi-server Management                           │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │ HTTPS/SSE                           │
└───────────────────────┼─────────────────────────────────────┘
                        │
                        │ Internet
                        ▼
┌─────────────────────────────────────────────────────────────┐
│           Remote Server (Production/Demo)                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Reverse Proxy (Caddy/Nginx)                         │  │
│  │  - SSL/TLS Termination                               │  │
│  │  - Domain: mcp.example.com                           │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │                                     │
│                       ▼                                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  MCP NetBird Server (Container)                      │  │
│  │  - Mode: SSE                                         │  │
│  │  - Port: 8001                                        │  │
│  │  - Endpoint: /sse                                    │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │ HTTP API                            │
│                       ▼                                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  NetBird Management API                              │  │
│  │  - Dashboard, Management, Signal                     │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Deployment Options

### Remote SSE Server

Deploy the MCP server as a remote service accessible over HTTPS. Best for:
- Production environments
- Team collaboration
- Centralized management
- Demo environments

#### 1. Deploy MCP Server Container

**Docker Compose** (`docker-compose-mcp-server.yml`):

```yaml
version: '3.8'

services:
  mcp-netbird:
    image: xnetadmin/mcp-netbird:latest
    container_name: mcp-netbird-server
    restart: unless-stopped
    command: ["-t", "sse", "-sse-address", "0.0.0.0:8001"]
    environment:
      - NETBIRD_API_TOKEN=your_api_token_here
      - NETBIRD_API_HOST=api.netbird.io
      - NETBIRD_MGMT_API_ENDPOINT=https://api.netbird.io
    ports:
      - "8001:8001"
    networks:
      - netbird

networks:
  netbird:
    external: true
    name: netbird_netbird
```

**Deploy**:
```bash
docker compose -f docker-compose-mcp-server.yml up -d
```

#### 2. Configure Reverse Proxy

**Caddy** (`Caddyfile`):

```caddyfile
mcp.example.com {
    # SSE endpoint for MCP
    reverse_proxy /sse mcp-netbird-server:8001 {
        header_up Host {upstream_hostport}
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
        
        # SSE specific - disable buffering
        flush_interval -1
    }
    
    # Root redirect
    redir / /sse
}
```

**Nginx**:

```nginx
server {
    listen 443 ssl http2;
    server_name mcp.example.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location /sse {
        proxy_pass http://localhost:8001;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # SSE specific
        proxy_buffering off;
        proxy_cache off;
        chunked_transfer_encoding off;
    }
}
```

#### 3. Configure Docker MCP Gateway

The Docker MCP Gateway handles protocol translation between local STDIO and remote SSE.

**Enable the Server** (`~/.docker/mcp/config.yaml`):

```yaml
netbird-mcp-server:
  netbird_host: api.netbird.io
  enabled: true
```

**Set API Token** (`~/.docker/mcp/.env`):

```bash
NETBIRD_API_TOKEN=your_api_token_here
```

**Add to Gateway** (`~/.kiro/settings/mcp.json` for Kiro):

```json
{
  "mcpServers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=netbird-mcp-server"
      ],
      "disabled": false,
      "autoApprove": ["*"]
    }
  }
}
```

**Verify Setup**:

```bash
# Check server is enabled
docker mcp server list | grep netbird

# Test connection
curl https://mcp.example.com/sse -H "Accept: text/event-stream"
```

---

### Local STDIO Mode

Run the MCP server locally for development or single-user scenarios.

#### Option 1: Direct Docker Run

**Kiro Configuration** (`~/.kiro/settings/mcp.json`):

```json
{
  "mcpServers": {
    "netbird": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "xnetadmin/mcp-netbird:latest",
        "-t",
        "stdio"
      ],
      "env": {
        "NETBIRD_API_TOKEN": "your_api_token_here",
        "NETBIRD_API_HOST": "api.netbird.io",
        "NETBIRD_MGMT_API_ENDPOINT": "https://api.netbird.io"
      },
      "disabled": false
    }
  }
}
```

#### Option 2: Claude Desktop

**Configuration** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "netbird": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "xnetadmin/mcp-netbird:latest",
        "-t",
        "stdio"
      ],
      "env": {
        "NETBIRD_API_TOKEN": "your_api_token_here",
        "NETBIRD_API_HOST": "api.netbird.io",
        "NETBIRD_MGMT_API_ENDPOINT": "https://api.netbird.io"
      }
    }
  }
}
```

---

### Docker MCP Gateway

Use the Docker MCP Gateway for advanced scenarios with multiple MCP servers.

#### 1. Add to Catalog

**Create/Update Catalog** (`~/.docker/mcp/catalogs/custom-servers.yaml`):

```yaml
version: 3
name: custom-servers
displayName: Custom MCP Servers
registry:
  netbird-mcp-server:
    description: Model Context Protocol server for Netbird VPN management
    title: Netbird MCP Server
    type: server
    image: xnetadmin/mcp-netbird:latest
    secrets:
      - name: netbird-mcp-server.netbird_api_token
        env: NETBIRD_API_TOKEN
        example: nbp_your_api_token_here
    env:
      - name: NETBIRD_HOST
        value: '{{netbird-mcp-server.netbird_host}}'
    config:
      - name: netbird-mcp-server
        description: Configure the connection to Netbird API
        type: object
        properties:
          netbird_host:
            type: string
            description: Netbird API host (e.g., api.netbird.io)
        required:
          - netbird_host
```

#### 2. Enable and Configure

```bash
# Enable the server
docker mcp server enable netbird-mcp-server

# Set API token
docker mcp secret set "netbird-mcp-server.netbird_api_token=your_token_here"

# Configure host (updates ~/.docker/mcp/config.yaml)
# Edit config.yaml manually:
# netbird-mcp-server:
#   netbird_host: api.netbird.io
#   enabled: true
```

#### 3. Add to Gateway

Update your MCP client configuration to include the gateway:

```json
{
  "mcpServers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=netbird-mcp-server"
      ],
      "disabled": false,
      "autoApprove": ["*"]
    }
  }
}
```

---

## Configuration

### Environment Variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `NETBIRD_API_TOKEN` | Yes | NetBird API token | `nbp_abc123...` |
| `NETBIRD_API_HOST` | Yes | NetBird API hostname (without protocol) | `api.netbird.io` |
| `NETBIRD_MGMT_API_ENDPOINT` | No | Full management API URL | `https://api.netbird.io` |

### Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-t, --transport` | `stdio` | Transport type: `stdio` or `sse` |
| `-sse-address` | `localhost:8001` | SSE server address (SSE mode only) |
| `-api-token` | - | NetBird API token (overrides env var) |
| `-api-host` | - | NetBird API host (overrides env var) |

### Getting a NetBird API Token

1. **Login to NetBird Dashboard**
   - Navigate to your NetBird dashboard
   - Go to Settings → Access Tokens

2. **Create Service User** (Recommended)
   - Create a dedicated service user for the MCP server
   - Assign appropriate role (admin or network_admin)
   - Generate API token for the service user

3. **Use Personal Token** (Development Only)
   - Generate a personal access token
   - Not recommended for production use

---

## Usage Examples

### List All Peers

```javascript
// Via MCP tool
mcp_MCP_DOCKER_list_netbird_peers()

// Returns array of peer objects with:
// - id, name, ip, dns_label
// - connected status
// - groups, os, version
// - last_seen, created_at
```

### Create a Group

```javascript
mcp_MCP_DOCKER_create_netbird_group({
  name: "developers",
  peers: []  // Add peer IDs here
})
```

### Create a Policy

```javascript
mcp_MCP_DOCKER_create_netbird_policy({
  name: "dev-to-prod",
  description: "Allow developers to access production",
  enabled: true,
  rules: [{
    action: "accept",
    bidirectional: true,
    description: "Dev to prod access",
    destinations: [{ id: "prod-group-id" }],
    sources: [{ id: "dev-group-id" }],
    enabled: true,
    name: "dev-prod-rule",
    protocol: "all"
  }]
})
```

### Create Setup Key

```javascript
mcp_MCP_DOCKER_create_netbird_setup_key({
  name: "dev-team-key",
  type: "reusable",
  expires_in: 2592000,  // 30 days in seconds
  auto_groups: ["dev-group-id"],
  usage_limit: 50
})
```

### Create Network with Resources

```javascript
// 1. Create network
const network = mcp_MCP_DOCKER_create_netbird_network({
  name: "production-network",
  description: "Production infrastructure network"
})

// 2. Add network resource
mcp_MCP_DOCKER_create_netbird_network_resource({
  network_id: network.id,
  name: "database-subnet",
  address: "10.0.1.0/24",
  enabled: true,
  groups: ["db-group-id"]
})

// 3. Add network router
mcp_MCP_DOCKER_create_netbird_network_router({
  network_id: network.id,
  peer: "router-peer-id",
  metric: 100,
  masquerade: true,
  enabled: true
})
```

### Configure DNS

```javascript
mcp_MCP_DOCKER_create_netbird_nameserver({
  name: "Cloudflare DNS",
  description: "Primary DNS resolver",
  nameservers: [
    { ip: "1.1.1.1", ns_type: "udp", port: 53 },
    { ip: "1.0.0.1", ns_type: "udp", port: 53 }
  ],
  enabled: true,
  groups: ["all-group-id"],
  primary: true,
  domains: [],
  search_domains_enabled: false
})
```

---

## Available Tools

### Account Management
- `get_netbird_account` - Get account information
- `update_netbird_account` - Update account settings

### Peer Management
- `list_netbird_peers` - List all peers
- `get_netbird_peer` - Get peer details
- `update_netbird_peer` - Update peer configuration
- `delete_netbird_peer` - Delete a peer

### Group Management
- `list_netbird_groups` - List all groups
- `get_netbird_group` - Get group details
- `create_netbird_group` - Create a new group
- `update_netbird_group` - Update group configuration
- `delete_netbird_group` - Delete a group

### Policy Management
- `list_netbird_policies` - List all policies
- `get_netbird_policy` - Get policy details
- `create_netbird_policy` - Create a new policy
- `update_netbird_policy` - Update policy configuration
- `delete_netbird_policy` - Delete a policy
- `list_policies_by_group` - List policies referencing a group
- `replace_group_in_policies` - Replace group across all policies

### Network Management
- `list_netbird_networks` - List all networks
- `get_netbird_network` - Get network details
- `create_netbird_network` - Create a new network
- `update_netbird_network` - Update network configuration
- `delete_netbird_network` - Delete a network

### Network Resource Management
- `list_netbird_network_resources` - List network resources
- `get_netbird_network_resource` - Get resource details
- `create_netbird_network_resource` - Create a network resource
- `update_netbird_network_resource` - Update resource configuration
- `delete_netbird_network_resource` - Delete a network resource

### Network Router Management
- `list_netbird_network_routers` - List network routers
- `get_netbird_network_router` - Get router details
- `create_netbird_network_router` - Create a network router
- `update_netbird_network_router` - Update router configuration
- `delete_netbird_network_router` - Delete a network router

### Route Management
- `list_netbird_routes` - List all routes
- `get_netbird_route` - Get route details
- `create_netbird_route` - Create a new route
- `update_netbird_route` - Update route configuration
- `delete_netbird_route` - Delete a route

### DNS Management
- `list_netbird_nameservers` - List all nameservers
- `get_netbird_nameserver` - Get nameserver details
- `create_netbird_nameserver` - Create a nameserver
- `update_netbird_nameserver` - Update nameserver configuration
- `delete_netbird_nameserver` - Delete a nameserver

### Setup Key Management
- `list_netbird_setup_keys` - List all setup keys
- `get_netbird_setup_key` - Get setup key details
- `create_netbird_setup_key` - Create a setup key
- `update_netbird_setup_key` - Update setup key configuration
- `delete_netbird_setup_key` - Delete a setup key

### User Management
- `list_netbird_users` - List all users
- `get_netbird_user` - Get user details
- `invite_netbird_user` - Invite a new user
- `update_netbird_user` - Update user configuration
- `delete_netbird_user` - Delete a user

### Posture Check Management
- `list_netbird_posture_checks` - List posture checks
- `get_netbird_posture_check` - Get posture check details
- `create_netbird_posture_check` - Create a posture check
- `update_netbird_posture_check` - Update posture check configuration
- `delete_netbird_posture_check` - Delete a posture check

### Port Allocation Management
- `list_netbird_port_allocations` - List port allocations
- `get_netbird_port_allocation` - Get port allocation details
- `create_netbird_port_allocation` - Create a port allocation
- `update_netbird_port_allocation` - Update port allocation configuration
- `delete_netbird_port_allocation` - Delete a port allocation

---

## Troubleshooting

### MCP Server Not Starting

**Check logs**:
```bash
docker logs mcp-netbird-server
```

**Common issues**:
- Missing or invalid API token
- Incorrect API host
- Port 8001 already in use
- Network connectivity issues

**Solutions**:
```bash
# Verify environment variables
docker exec mcp-netbird-server env | grep NETBIRD

# Test API connectivity
curl -H "Authorization: Bearer $NETBIRD_API_TOKEN" \
  https://api.netbird.io/api/accounts
```

### SSE Connection Timeouts

**Symptoms**:
- MCP client shows "Request timed out"
- No response from server after 60 seconds

**Check**:
```bash
# Test SSE endpoint
curl https://mcp.example.com/sse -H "Accept: text/event-stream"

# Should return:
# event: endpoint
# data: /message?sessionId=...
```

**Common causes**:
- Reverse proxy buffering enabled
- Firewall blocking SSE connections
- SSL/TLS certificate issues

**Solutions**:
- Disable buffering in reverse proxy (see configuration examples)
- Check firewall rules
- Verify SSL certificate is valid

### Authentication Errors (401)

**Error**: `{"message":"token invalid","code":401}`

**Causes**:
- Expired API token
- Invalid API token
- Wrong API host

**Solutions**:
```bash
# Verify token works directly
curl -H "Authorization: Bearer $NETBIRD_API_TOKEN" \
  https://api.netbird.io/api/accounts

# Regenerate token if needed
# Update configuration with new token
```

### Permission Denied (403)

**Error**: `{"message":"permission denied","code":403}`

**Causes**:
- Insufficient permissions for operation
- Service user role too restrictive

**Solutions**:
- Verify user role (admin or network_admin required for most operations)
- Check specific operation requirements
- Use admin token for setup keys and user management

### Docker MCP Gateway Issues

**Server not listed**:
```bash
# Check if server is enabled
docker mcp server list | grep netbird

# Enable if needed
docker mcp server enable netbird-mcp-server
```

**Configuration not loading**:
```bash
# Verify config file
cat ~/.docker/mcp/config.yaml | grep netbird

# Check secrets
cat ~/.docker/mcp/.env | grep NETBIRD

# Restart gateway (restart your MCP client)
```

### Tools Not Appearing in Client

**Kiro**:
- Check MCP Servers view in sidebar
- Verify gateway is running
- Restart Kiro to reload configuration

**Claude Desktop**:
- Check Developer Tools console for errors
- Verify configuration file syntax
- Restart Claude Desktop

**General**:
```bash
# Test server directly
docker run --rm -i xnetadmin/mcp-netbird:latest -t stdio \
  <<< '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

---

## Security Best Practices

### API Token Management

1. **Use Service Users**
   - Create dedicated service users for MCP servers
   - Don't use personal tokens in production

2. **Rotate Tokens Regularly**
   - Set expiration dates on tokens
   - Rotate tokens every 90 days

3. **Limit Permissions**
   - Use network_admin role when possible
   - Only use admin role when necessary

4. **Secure Storage**
   - Store tokens in environment variables or secrets management
   - Never commit tokens to version control
   - Use `.env` files with proper permissions (600)

### Network Security

1. **Use HTTPS**
   - Always use SSL/TLS for remote SSE servers
   - Use valid certificates (Let's Encrypt recommended)

2. **Firewall Rules**
   - Restrict access to MCP server port
   - Use VPN or IP whitelisting for production

3. **Reverse Proxy**
   - Use reverse proxy for SSL termination
   - Enable rate limiting
   - Add authentication layer if needed

### Monitoring

1. **Log Access**
   - Monitor MCP server logs
   - Track API usage
   - Alert on authentication failures

2. **Audit Operations**
   - Review NetBird audit logs
   - Track changes made via MCP
   - Monitor for suspicious activity

---

## Performance Optimization

### Local STDIO Mode
- **Pros**: Lowest latency, no network overhead
- **Cons**: Container startup time per session
- **Best for**: Development, single-user scenarios

### Remote SSE Mode
- **Pros**: Persistent connection, shared access
- **Cons**: Network latency, requires infrastructure
- **Best for**: Production, team collaboration

### Caching Strategies

The MCP server makes direct API calls without caching. For high-frequency operations:

1. **Batch Operations**: Group multiple operations together
2. **Client-Side Caching**: Cache results in your application
3. **Rate Limiting**: Respect NetBird API rate limits

---

## Support

- **Issues**: https://github.com/XNet-NGO/mcp-netbird/issues
- **Discussions**: https://github.com/XNet-NGO/mcp-netbird/discussions
- **NetBird Docs**: https://docs.netbird.io
- **MCP Specification**: https://modelcontextprotocol.io

---

## License

Apache License 2.0 - See [LICENSE](../LICENSE) for details.
