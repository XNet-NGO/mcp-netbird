# NetBird MCP Server

A comprehensive [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for [NetBird](https://netbird.io/) providing 50+ tools for complete VPN infrastructure management through AI assistants.

**Maintained by XNet Inc.**  
**Lead Developer: Joshua S. Doucette**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/XNet-NGO/mcp-netbird)](go.mod)
[![Release](https://img.shields.io/github/v/release/XNet-NGO/mcp-netbird)](https://github.com/XNet-NGO/mcp-netbird/releases)
[![Docker Hub](https://img.shields.io/docker/pulls/xnetadmin/mcp-netbird)](https://hub.docker.com/r/xnetadmin/mcp-netbird)

## About

This MCP server enables AI assistants like Kiro, Claude Desktop, and other MCP clients to programmatically manage NetBird VPN infrastructure. It provides:

- **50+ Management Tools**: Complete CRUD operations for all NetBird resources
- **Multiple Deployment Options**: Local STDIO, Remote SSE, or Docker MCP Gateway
- **Advanced Policy Management**: Validation, dependency tracking, and bulk operations
- **Helper Functions**: Group consolidation, policy templates, and common workflows
- **Production Ready**: Comprehensive error handling, logging, and security features

Originally derived from the MCP Server for Grafana by Grafana Labs, this project has been substantially extended and enhanced for NetBird infrastructure management.

## Quick Start

### Docker (Recommended)

The easiest way to get started is using Docker with the Docker MCP Gateway:

```bash
# Pull the latest image
docker pull xnetadmin/mcp-netbird:latest

# Run in SSE mode for remote access
docker run -d \
  --name mcp-netbird \
  -p 8001:8001 \
  -e NETBIRD_API_TOKEN=your_token_here \
  -e NETBIRD_API_HOST=api.netbird.io \
  xnetadmin/mcp-netbird:latest \
  -t sse -sse-address 0.0.0.0:8001
```

Then configure your MCP client (see [Configuration](#configuration) below).

### Installing from Releases

Download pre-built binaries for your platform from the [releases page](https://github.com/XNet-NGO/mcp-netbird/releases).

**Linux (Debian/Ubuntu)**:
```bash
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_linux_x86_64.deb
sudo dpkg -i mcp-netbird_VERSION_linux_x86_64.deb
```

**Linux (Other)**:
```bash
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Linux_x86_64.tar.gz
tar -xzf mcp-netbird_VERSION_Linux_x86_64.tar.gz
sudo mv mcp-netbird /usr/local/bin/
```

**macOS**:
```bash
# Intel Macs
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Darwin_x86_64.tar.gz
tar -xzf mcp-netbird_VERSION_Darwin_x86_64.tar.gz
sudo mv mcp-netbird /usr/local/bin/

# Apple Silicon
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Darwin_arm64.tar.gz
tar -xzf mcp-netbird_VERSION_Darwin_arm64.tar.gz
sudo mv mcp-netbird /usr/local/bin/
```

**Windows**: Download the ZIP from [releases](https://github.com/XNet-NGO/mcp-netbird/releases), extract, and add to PATH.

### Building from Source

```bash
git clone https://github.com/XNet-NGO/mcp-netbird
cd mcp-netbird
make install
```

Or install directly from GitHub:
```bash
go install github.com/XNet-NGO/mcp-netbird/cmd/mcp-netbird@latest
```

## Configuration

The NetBird MCP server supports three deployment modes. Choose the one that best fits your use case.

### Getting a NetBird API Token

Before configuring, you'll need a NetBird API token:

1. Login to your NetBird dashboard (https://app.netbird.io or your self-hosted instance)
2. Go to **Settings → Access Tokens**
3. Create a **Service User** (recommended for production) or use a personal token
4. Assign appropriate role: **admin** (full access) or **network_admin** (network management only)
5. Generate and copy the API token (starts with `nbp_`)

### Deployment Options

#### Option 1: Docker MCP Gateway (Recommended)

The Docker MCP Gateway provides the best experience for remote MCP servers, handling protocol translation between local STDIO and remote SSE.

**1. Enable the NetBird MCP Server**

Create or update `~/.docker/mcp/config.yaml`:
```yaml
netbird-mcp-server:
  netbird_host: api.netbird.io  # or your self-hosted domain
  enabled: true
```

**2. Set Your API Token**

Create or update `~/.docker/mcp/.env`:
```bash
NETBIRD_API_TOKEN=nbp_your_token_here
```

**3. Configure Your MCP Client**

For **Kiro** (`~/.kiro/settings/mcp.json`):
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

For **Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):
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
      ]
    }
  }
}
```

**4. Verify Setup**

```bash
# Check server is enabled
docker mcp server list | grep netbird

# Restart your MCP client (Kiro/Claude Desktop)
# Tools will appear with mcp_MCP_DOCKER_ prefix
```

#### Option 2: Local STDIO Mode

Run the MCP server locally for development or single-user scenarios.

**For Kiro** (`~/.kiro/settings/mcp.json`):
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
        "NETBIRD_API_TOKEN": "nbp_your_token_here",
        "NETBIRD_API_HOST": "api.netbird.io"
      },
      "disabled": false
    }
  }
}
```

**For Claude Desktop**:
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
        "NETBIRD_API_TOKEN": "nbp_your_token_here",
        "NETBIRD_API_HOST": "api.netbird.io"
      }
    }
  }
}
```

#### Option 3: Remote SSE Server

Deploy the MCP server as a remote service for team collaboration or production use.

**1. Deploy MCP Server**

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  mcp-netbird:
    image: xnetadmin/mcp-netbird:latest
    container_name: mcp-netbird-server
    restart: unless-stopped
    command: ["-t", "sse", "-sse-address", "0.0.0.0:8001"]
    environment:
      - NETBIRD_API_TOKEN=nbp_your_token_here
      - NETBIRD_API_HOST=api.netbird.io
    ports:
      - "8001:8001"
```

Deploy:
```bash
docker compose up -d
```

**2. Configure Reverse Proxy (Optional but Recommended)**

**Caddy** (`Caddyfile`):
```caddyfile
mcp.example.com {
    reverse_proxy /sse localhost:8001 {
        flush_interval -1  # Required for SSE
    }
    redir / /sse
}
```

**Nginx**:
```nginx
server {
    listen 443 ssl http2;
    server_name mcp.example.com;
    
    location /sse {
        proxy_pass http://localhost:8001;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_buffering off;  # Required for SSE
        proxy_cache off;
        chunked_transfer_encoding off;
    }
}
```

**3. Configure Docker MCP Gateway**

Update `~/.docker/mcp/config.yaml`:
```yaml
netbird-mcp-server:
  url: https://mcp.example.com/sse
  transport: sse
  enabled: true
```

Then add to your MCP client configuration as shown in Option 1.

### Configuration Priority

When multiple configuration sources provide the same value:

**CLI Arguments > HTTP Headers > Environment Variables**

Example:
```bash
# Environment variable
export NETBIRD_API_TOKEN="token-from-env"

# CLI argument overrides environment
mcp-netbird --api-token "token-from-cli"
# Result: Uses "token-from-cli"
```

### Self-Hosted NetBird

For self-hosted NetBird instances, set the API host to your domain:

```yaml
# Docker MCP Gateway
netbird-mcp-server:
  netbird_host: api.yourdomain.com
  enabled: true
```

Or for STDIO mode:
```json
{
  "env": {
    "NETBIRD_API_TOKEN": "your_token",
    "NETBIRD_API_HOST": "api.yourdomain.com"
  }
}
```

### Troubleshooting

**Tools not appearing**: Restart your MCP client after configuration changes.

**Connection timeout**: Verify API token is valid and has appropriate permissions.

**401 Unauthorized**: Check that your API token hasn't expired.

**For detailed setup instructions**, see [docs/MCP_SETUP_GUIDE.md](docs/MCP_SETUP_GUIDE.md).

## Features

### Complete NetBird API Coverage

The MCP server provides 50+ tools covering all NetBird resources:

| Resource | Operations | Description |
|----------|-----------|-------------|
| **Peers** | list, get, update, delete | Manage network peers and their configuration |
| **Groups** | list, get, create, update, delete | Organize peers into logical groups |
| **Policies** | list, get, create, update, delete | Control network access between groups |
| **Networks** | list, get, create, update, delete | Manage network configurations |
| **Network Resources** | list, get, create, update, delete | Define network subnets and resources |
| **Network Routers** | list, get, create, update, delete | Configure routing peers for networks |
| **Nameservers** | list, get, create, update, delete | Manage DNS nameserver groups |
| **Routes** | list, get, create, update, delete | Configure network routes (legacy) |
| **Setup Keys** | list, get, create, update, delete | Generate peer enrollment keys |
| **Users** | list, get, invite, update, delete | Manage user accounts and permissions |
| **Posture Checks** | list, get, create, update, delete | Define security posture requirements |
| **Port Allocations** | list, get, create, update, delete | Manage ingress port forwarding |
| **Account** | get, update | Configure account-wide settings |

### Helper Tools

Advanced tools for common administrative workflows:

- **list_policies_by_group**: Find all policies referencing a specific group
- **replace_group_in_policies**: Bulk replace groups across all policies
- **get_policy_template**: Get example policy structures with documentation

### Key Capabilities

- ✅ **Full CRUD Operations**: Create, read, update, and delete all NetBird resources
- ✅ **Policy Validation**: Automatic validation of policy rules before API submission
- ✅ **Dependency Tracking**: Find and manage resource dependencies
- ✅ **Bulk Operations**: Perform operations across multiple resources
- ✅ **Error Handling**: Comprehensive error messages and recovery suggestions
- ✅ **Production Ready**: Secure authentication, logging, and monitoring support

## Usage Examples

### Basic Operations

**List all peers**:
```javascript
mcp_MCP_DOCKER_list_netbird_peers()
// Returns: Array of peer objects with IP, status, groups, etc.
```

**Create a group**:
```javascript
mcp_MCP_DOCKER_create_netbird_group({
  name: "developers",
  peers: []  // Add peer IDs here
})
```

**Create a setup key**:
```javascript
mcp_MCP_DOCKER_create_netbird_setup_key({
  name: "dev-team-key",
  type: "reusable",
  expires_in: 2592000,  // 30 days
  auto_groups: ["dev-group-id"],
  usage_limit: 50
})
```

### Policy Management

**Create a simple policy**:
```javascript
mcp_MCP_DOCKER_create_netbird_policy({
  name: "Admin SSH Access",
  description: "Allow admins to SSH to servers",
  enabled: true,
  rules: [{
    name: "SSH Rule",
    enabled: true,
    action: "accept",
    bidirectional: false,
    protocol: "tcp",
    sources: ["admin-group-id"],
    destinations: ["server-group-id"],
    port_ranges: [{ start: 22, end: 22 }]
  }]
})
```

**Find policies using a group**:
```javascript
mcp_MCP_DOCKER_list_policies_by_group({
  group_id: "d535b93ngf8s73892nng"
})
// Returns: List of policies referencing this group
```

**Replace a group across all policies**:
```javascript
mcp_MCP_DOCKER_replace_group_in_policies({
  old_group_id: "old-group-id",
  new_group_id: "new-group-id"
})
// Updates all policies to use the new group
```

### Network Configuration

**Create a network**:
```javascript
const network = mcp_MCP_DOCKER_create_netbird_network({
  name: "production-network",
  description: "Production infrastructure"
})
```

**Add network resource**:
```javascript
mcp_MCP_DOCKER_create_netbird_network_resource({
  network_id: network.id,
  name: "database-subnet",
  address: "10.0.1.0/24",
  enabled: true,
  groups: ["db-group-id"]
})
```

**Configure network router**:
```javascript
mcp_MCP_DOCKER_create_netbird_network_router({
  network_id: network.id,
  peer: "router-peer-id",
  metric: 100,
  masquerade: true,
  enabled: true
})
```

### DNS Configuration

**Add nameserver group**:
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

### Complete Workflow Example

Here's a complete example of setting up a new environment:

```javascript
// 1. Create groups
const adminGroup = mcp_MCP_DOCKER_create_netbird_group({
  name: "admins",
  peers: []
})

const serverGroup = mcp_MCP_DOCKER_create_netbird_group({
  name: "servers",
  peers: []
})

// 2. Create setup keys
const adminKey = mcp_MCP_DOCKER_create_netbird_setup_key({
  name: "admin-key",
  type: "reusable",
  expires_in: 2592000,
  auto_groups: [adminGroup.id],
  usage_limit: 10
})

const serverKey = mcp_MCP_DOCKER_create_netbird_setup_key({
  name: "server-key",
  type: "reusable",
  expires_in: 2592000,
  auto_groups: [serverGroup.id],
  usage_limit: 100
})

// 3. Create policy
mcp_MCP_DOCKER_create_netbird_policy({
  name: "Admin Access",
  description: "Allow admins to access servers",
  enabled: true,
  rules: [{
    name: "Admin to Servers",
    enabled: true,
    action: "accept",
    bidirectional: true,
    protocol: "all",
    sources: [adminGroup.id],
    destinations: [serverGroup.id]
  }]
})

// 4. Configure DNS
mcp_MCP_DOCKER_create_netbird_nameserver({
  name: "Primary DNS",
  nameservers: [
    { ip: "1.1.1.1", ns_type: "udp", port: 53 }
  ],
  enabled: true,
  groups: [adminGroup.id, serverGroup.id],
  primary: true
})

// Setup keys can now be used to enroll peers
console.log("Admin key:", adminKey.key)
console.log("Server key:", serverKey.key)
```

### Using with AI Assistants

Once configured, you can interact with NetBird through natural language:

**Kiro/Claude Desktop**:
- "Can you explain my NetBird peers, groups and policies?"
- "Create a new group called 'developers' and generate a setup key for it"
- "Show me all policies that reference the admin group"
- "Configure DNS to use Cloudflare for all peers"

For more examples and detailed documentation, see [docs/MCP_SETUP_GUIDE.md](docs/MCP_SETUP_GUIDE.md).

## Advanced Usage

### Policy Rule Format

When creating or updating policies, rules must follow this structure:

```json
{
  "name": "Rule Name",
  "description": "Optional description",
  "enabled": true,
  "action": "accept",
  "bidirectional": true,
  "protocol": "all",
  "sources": ["group-id-1"],
  "destinations": ["group-id-2"],
  "port_ranges": [
    {"start": 80, "end": 80},
    {"start": 443, "end": 443}
  ]
}
```

**Required Fields**:
- `name`, `enabled`, `action`, `bidirectional`, `protocol`
- At least one source (array of group IDs or sourceResource object)
- At least one destination (array of group IDs or destinationResource object)

**Optional Fields**:
- `description`, `port_ranges` (TCP/UDP only), `authorized_groups`

**Get policy templates**:
```javascript
mcp_MCP_DOCKER_get_policy_template()
// Returns example policy structures with documentation
```

### Group Management

**Find group dependencies**:
```javascript
mcp_MCP_DOCKER_list_policies_by_group({
  group_id: "d535b93ngf8s73892nng"
})
// Returns all policies referencing this group
```

**Force delete a group** (removes from all policies):
```javascript
mcp_MCP_DOCKER_delete_netbird_group({
  group_id: "group-id",
  force: true
})
```

### Production Deployment

For production environments, deploy the MCP server as a remote SSE service:

1. **Deploy with Docker Compose** (see [Configuration](#configuration))
2. **Configure reverse proxy** with SSL/TLS (Caddy or Nginx)
3. **Set up monitoring** and logging
4. **Use service user** with appropriate permissions
5. **Rotate API tokens** regularly

See [docs/MCP_SETUP_GUIDE.md](docs/MCP_SETUP_GUIDE.md) for detailed production deployment guide.

## Docker

### Building Custom Images

Build and tag an image:
```bash
docker build -t mcp-netbird-sse:v1 -f Dockerfile.sse .
```

Run with different configuration methods:

**Environment Variables**:
```bash
docker run --name mcp-netbird -p 8001:8001 \
  -e NETBIRD_API_TOKEN=your_token \
  -e NETBIRD_API_HOST=api.netbird.io \
  mcp-netbird-sse:v1
```

**Command-Line Arguments**:
```bash
docker run --name mcp-netbird -p 8001:8001 \
  mcp-netbird-sse:v1 \
  --transport sse \
  --sse-address :8001 \
  --api-token your_token \
  --api-host api.netbird.io
```

**Stateless (HTTP Headers)**:
```bash
# Run without credentials
docker run --name mcp-netbird -p 8001:8001 mcp-netbird-sse:v1

# Pass credentials per request
curl -X POST http://localhost:8001/sse \
  -H "X-Netbird-API-Token: your_token" \
  -H "X-Netbird-Host: api.netbird.io" \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

### Using ToolHive

[ToolHive](https://github.com/StacklokLabs/toolhive) simplifies MCP server deployment:

```bash
# Install thv
# See: https://github.com/StacklokLabs/toolhive#installation

# Add NetBird API token
thv secret set netbird

# Build SSE image
docker build -t mcp-netbird-sse:v1 -f Dockerfile.sse .

# Start server
thv run --secret netbird,target=NETBIRD_API_TOKEN \
  --transport sse \
  --name thv-mcp-netbird \
  --port 8080 \
  --target-port 8001 \
  mcp-netbird-sse:v1

# Stop server
thv stop thv-mcp-netbird
```

## Development

Contributions are welcome! Please open an issue or submit a pull request.

### Prerequisites

- Go 1.21 or later
- Docker (for testing SSE mode)
- Make (optional, for convenience commands)

### Running Locally

**STDIO mode**:
```bash
export NETBIRD_API_TOKEN=your_token
go run cmd/mcp-netbird/main.go
```

**SSE mode**:
```bash
go run cmd/mcp-netbird/main.go \
  --transport sse \
  --sse-address :8001 \
  --api-token your_token \
  --api-host api.netbird.io
```

### Debugging with MCP Inspector

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is an interactive tool for testing MCP servers:

```bash
# Install and run
npx @modelcontextprotocol/inspector

# In the UI, configure:
# - Transport: stdio or sse
# - Command: mcp-netbird --api-token your_token
# - Or URL: http://localhost:8001/sse (for SSE mode)
```

### Testing

Run the test suite:
```bash
make test
```

Run specific tests:
```bash
go test ./tools -v -run TestListPeers
```

### Linting

```bash
make lint
```

### Adding New Tools

1. Create a new file in `tools/` (e.g., `tools/new_resource.go`)
2. Implement the tool functions following existing patterns
3. Add the tool to `func newServer()` in `cmd/mcp-netbird/main.go`
4. Add tests in `tools/new_resource_test.go`
5. Update documentation

### Project Structure

```
mcp-netbird/
├── cmd/mcp-netbird/     # Main application entry point
├── tools/               # MCP tool implementations
│   ├── peers.go
│   ├── groups.go
│   ├── policies.go
│   └── ...
├── docs/                # Documentation
│   └── MCP_SETUP_GUIDE.md
├── Dockerfile           # STDIO mode container
├── Dockerfile.sse       # SSE mode container
└── Makefile            # Build and development commands
```

## Documentation

- **[MCP Setup Guide](docs/MCP_SETUP_GUIDE.md)**: Comprehensive setup guide with all deployment options
- **[NetBird API Documentation](https://docs.netbird.io/api)**: Official NetBird API reference
- **[Model Context Protocol](https://modelcontextprotocol.io)**: MCP specification and documentation

## Support and Community

- **Issues**: [GitHub Issues](https://github.com/XNet-NGO/mcp-netbird/issues)
- **Discussions**: [GitHub Discussions](https://github.com/XNet-NGO/mcp-netbird/discussions)
- **NetBird Community**: [NetBird Slack](https://netbird.io/community)

## License

This project is licensed under the [Apache License, Version 2.0](LICENSE).

**Copyright 2025-2026 XNet Inc.**  
**Copyright 2025-2026 Joshua S. Doucette**

### Attribution

This project was originally derived from the [MCP Server for Grafana](https://github.com/grafana/mcp-grafana) developed by Grafana Labs. The current codebase has been substantially modified and extended for NetBird infrastructure management.

This project uses [MCP Go](https://github.com/mark3labs/mcp-go) developed by Mark III Labs.

---

**Maintained by XNet Inc. | Lead Developer: Joshua S. Doucette**

For questions or support, please open an issue on GitHub or reach out through our community channels.