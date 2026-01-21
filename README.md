# NetBird MCP Server

A comprehensive [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for [NetBird](https://netbird.io/) providing full CRUD operations, policy management, and automation workflows.

**Maintained by XNet Inc.**  
**Lead Developer: Joshua S. Doucette**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/XNet-NGO/mcp-netbird)](go.mod)
[![Release](https://img.shields.io/github/v/release/XNet-NGO/mcp-netbird)](https://github.com/XNet-NGO/mcp-netbird/releases)

## About

This MCP server provides complete management capabilities for NetBird networks, including:
- Full CRUD operations for all NetBird resources
- Advanced policy management with validation
- Group consolidation and dependency workflows
- Helper functions for common administrative tasks
- Comprehensive error handling and documentation

Originally derived from the MCP Server for Grafana by Grafana Labs, this project has been substantially extended and enhanced.

## Installing

### Official Releases

Download pre-built binaries for your platform from the [releases page](https://github.com/XNet-NGO/mcp-netbird/releases).

#### Linux (Debian/Ubuntu)

```bash
# Download the .deb package for your architecture
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_linux_x86_64.deb
sudo dpkg -i mcp-netbird_VERSION_linux_x86_64.deb
```

#### Linux (Other Distributions)

```bash
# Download and extract the tarball
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Linux_x86_64.tar.gz
tar -xzf mcp-netbird_VERSION_Linux_x86_64.tar.gz
sudo mv mcp-netbird /usr/local/bin/
```

#### macOS

```bash
# Intel Macs
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Darwin_x86_64.tar.gz
tar -xzf mcp-netbird_VERSION_Darwin_x86_64.tar.gz
sudo mv mcp-netbird /usr/local/bin/

# Apple Silicon Macs
wget https://github.com/XNet-NGO/mcp-netbird/releases/latest/download/mcp-netbird_VERSION_Darwin_arm64.tar.gz
tar -xzf mcp-netbird_VERSION_Darwin_arm64.tar.gz
sudo mv mcp-netbird /usr/local/bin/
```

#### Windows

Download the ZIP file from the [releases page](https://github.com/XNet-NGO/mcp-netbird/releases), extract it, and add the directory to your PATH.

### Installing from source

#### Clone the repository

```bash
git clone https://github.com/XNet-NGO/mcp-netbird
```

#### Build and install

```bash
cd mcp-netbird && \
make install
```

### Installing from GitHub

```bash
go install github.com/XNet-NGO/mcp-netbird/cmd/mcp-netbird@latest
```

## Configuration

The NetBird MCP server supports three configuration methods with a clear priority order. This enables flexible deployment in various environments, including stateless containers.

### Configuration Methods

The server accepts configuration through three methods (in priority order):

1. **Command-Line Arguments** (highest priority)
2. **HTTP Headers** (SSE mode only, medium priority)
3. **Environment Variables** (lowest priority)

#### 1. Command-Line Arguments

Pass configuration directly via CLI flags:

```bash
mcp-netbird --api-token "your-token" --api-host "api.netbird.io"
```

**Available Flags:**
- `--api-token`: Your Netbird API token
- `--api-host`: The Netbird API host without protocol (e.g., `api.netbird.io`)
- `--transport`: Transport mode - `stdio` (default) or `sse`
- `--sse-address`: Address for SSE mode (default: `:8001`)

**Example - Stdio Mode:**
```bash
mcp-netbird \
  --api-token "nb_1234567890abcdef" \
  --api-host "api.netbird.io"
```

**Example - SSE Mode:**
```bash
mcp-netbird \
  --transport sse \
  --sse-address :8001 \
  --api-token "nb_1234567890abcdef" \
  --api-host "api.netbird.io"
```

#### 2. HTTP Headers (SSE Mode Only)

When running in SSE mode, you can pass configuration via HTTP headers with each request:

**Available Headers:**
- `X-Netbird-API-Token`: Your Netbird API token
- `X-Netbird-Host`: The Netbird API host without protocol

**Example - Using curl:**
```bash
curl -X POST http://localhost:8001/sse \
  -H "X-Netbird-API-Token: nb_1234567890abcdef" \
  -H "X-Netbird-Host: api.netbird.io" \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

**Example - Docker MCP Toolkit:**
```json
{
  "mcpServers": {
    "netbird": {
      "url": "http://localhost:8001/sse",
      "headers": {
        "X-Netbird-API-Token": "nb_1234567890abcdef",
        "X-Netbird-Host": "api.netbird.io"
      }
    }
  }
}
```

**Note:** HTTP headers are only used in SSE mode. In stdio mode, only CLI arguments and environment variables are available.

#### 3. Environment Variables

Set configuration via environment variables (traditional method):

```bash
export NETBIRD_API_TOKEN="nb_1234567890abcdef"
export NETBIRD_HOST="api.netbird.io"
mcp-netbird
```

**Available Variables:**
- `NETBIRD_API_TOKEN`: Your Netbird API token (required)
- `NETBIRD_HOST`: The Netbird API host (default: `api.netbird.io`)

**Example - MCP Client Configuration:**
```json
{
  "mcpServers": {
    "netbird": {
      "command": "mcp-netbird",
      "args": [],
      "env": {
        "NETBIRD_API_TOKEN": "nb_1234567890abcdef",
        "NETBIRD_HOST": "api.netbird.io"
      }
    }
  }
}
```

### Configuration Priority Order

When multiple configuration sources provide the same value, the priority order is:

**CLI Arguments > HTTP Headers > Environment Variables**

**Example:**
```bash
# Environment variable
export NETBIRD_API_TOKEN="token-from-env"

# CLI argument overrides environment variable
mcp-netbird --api-token "token-from-cli"
# Result: Uses "token-from-cli"
```

**Example - SSE Mode with Multiple Sources:**
```bash
# Environment variable
export NETBIRD_API_TOKEN="token-from-env"

# Start server with CLI argument
mcp-netbird --transport sse --api-token "token-from-cli"

# Make request with HTTP header
curl -H "X-Netbird-API-Token: token-from-header" http://localhost:8001/sse
# Result: Uses "token-from-cli" (CLI has highest priority)
```

### Configuration Validation

The server validates all configuration values:

- **API Token**: Must not be empty or whitespace-only
- **API Host**: Must not be empty or whitespace-only
- **Protocol Prefix**: Automatically stripped if present (e.g., `https://api.netbird.io` → `api.netbird.io`)

**Invalid configurations will result in descriptive error messages:**
```
Error: API token is required but not provided
Error: API token cannot be empty or whitespace-only
Error: API host cannot be empty or whitespace-only
Error: API host 'invalid host!' is not a valid URL format
```

### Troubleshooting

#### Common Configuration Errors

**1. Missing API Token**

**Error:** `API token is required but not provided`

**Solution:** Provide the API token via one of the three methods:
```bash
# Option 1: CLI argument
mcp-netbird --api-token "your-token"

# Option 2: Environment variable
export NETBIRD_API_TOKEN="your-token"
mcp-netbird

# Option 3: HTTP header (SSE mode only)
curl -H "X-Netbird-API-Token: your-token" http://localhost:8001/sse
```

**2. Empty or Whitespace-Only Values**

**Error:** `API token cannot be empty or whitespace-only`

**Solution:** Ensure your token/host values are not empty strings or only whitespace:
```bash
# Bad
mcp-netbird --api-token "   "

# Good
mcp-netbird --api-token "nb_1234567890abcdef"
```

**3. Invalid API Host Format**

**Error:** `API host 'invalid host!' is not a valid URL format`

**Solution:** Use a valid hostname without protocol:
```bash
# Bad
mcp-netbird --api-host "https://api.netbird.io"
mcp-netbird --api-host "invalid host!"

# Good
mcp-netbird --api-host "api.netbird.io"
mcp-netbird --api-host "api.yourdomain.com"
```

**Note:** If you include a protocol prefix (`http://` or `https://`), it will be automatically stripped.

**4. HTTP 401 Unauthorized (SSE Mode)**

**Error:** HTTP 401 response when making SSE requests

**Solution:** Ensure API token is provided via CLI argument, HTTP header, or environment variable:
```bash
# Start server with token
mcp-netbird --transport sse --api-token "your-token"

# Or provide token in request header
curl -H "X-Netbird-API-Token: your-token" http://localhost:8001/sse
```

**5. Configuration Not Taking Effect**

**Issue:** Your configuration seems to be ignored

**Solution:** Check the priority order. CLI arguments override HTTP headers, which override environment variables:
```bash
# If you have an environment variable set
export NETBIRD_API_TOKEN="old-token"

# But pass a CLI argument
mcp-netbird --api-token "new-token"

# The CLI argument "new-token" will be used (higher priority)
```

**6. Stdio Mode Not Reading HTTP Headers**

**Issue:** HTTP headers are not being used in stdio mode

**Solution:** HTTP headers are only available in SSE mode. Use CLI arguments or environment variables for stdio mode:
```bash
# Stdio mode (default) - use CLI args or env vars
mcp-netbird --api-token "your-token"

# SSE mode - can use HTTP headers
mcp-netbird --transport sse --api-token "your-token"
```

### Use Cases

#### Stateless Docker Containers

For Docker MCP Toolkit or other stateless container environments, use HTTP headers:

```dockerfile
# Dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o mcp-netbird cmd/mcp-netbird/main.go
EXPOSE 8001
CMD ["./mcp-netbird", "--transport", "sse", "--sse-address", ":8001"]
```

```bash
# Run without environment variables
docker run -p 8001:8001 mcp-netbird-sse:v1

# Pass credentials via HTTP headers
curl -H "X-Netbird-API-Token: your-token" http://localhost:8001/sse
```

#### Development and Testing

Use CLI arguments for quick testing without modifying environment:

```bash
# Test with different tokens quickly
mcp-netbird --api-token "test-token-1"
mcp-netbird --api-token "test-token-2"

# Test with self-hosted instance
mcp-netbird --api-token "your-token" --api-host "api.yourdomain.com"
```

#### Production Deployment

Use environment variables for traditional deployments:

```bash
# Set once in your deployment configuration
export NETBIRD_API_TOKEN="production-token"
export NETBIRD_HOST="api.netbird.io"

# Start server
mcp-netbird --transport sse --sse-address :8001
```

#### Multi-Tenant SSE Server

Use HTTP headers to support multiple tenants with different credentials:

```bash
# Start server without credentials
mcp-netbird --transport sse

# Each tenant provides their own credentials
curl -H "X-Netbird-API-Token: tenant1-token" http://localhost:8001/sse
curl -H "X-Netbird-API-Token: tenant2-token" http://localhost:8001/sse
```

## Features

This server uses the Netbird API to provide LLMs information about Netbird network. Currently it's a 1:1 mapping of select read-only Netbird API resources to tools.

- [x] Uses Netbird API to access configuration and status
- [x] Configurable API endpoint
- [x] Secure token-based authentication for Netbird API

### Tools

The MCP-NetBird server provides comprehensive CRUD operations for all NetBird resources:

#### Resource Management Tools

| Category | Operations | Description |
| --- | --- | --- |
| **Peers** | list, get, update, delete | Manage network peers |
| **Groups** | list, get, create, update, delete | Manage peer groups |
| **Policies** | list, get, create, update, delete | Manage access control policies |
| **Networks** | list, get, create, update, delete | Manage network configurations |
| **Network Resources** | list, get, create, update, delete | Manage network resources (subnets) |
| **Network Routers** | list, get, create, update, delete | Manage network routing peers |
| **Nameservers** | list, get, create, update, delete | Manage DNS nameserver groups |
| **Routes** | list, get, create, update, delete | Manage network routes (deprecated) |
| **Setup Keys** | list, get, create, update, delete | Manage peer setup keys |
| **Users** | list, get, invite, update, delete | Manage user accounts |
| **Posture Checks** | list, get, create, update, delete | Manage security posture checks |
| **Port Allocations** | list, get, create, update, delete | Manage ingress port allocations |
| **Account** | get, update | Manage account settings |

#### Helper Tools

| Tool | Description |
| --- | --- |
| `list_policies_by_group` | Find all policies that reference a specific group |
| `replace_group_in_policies` | Replace one group with another across all policies |
| `get_policy_template` | Get example policy structures with documentation |

### Adding tools

To add new tools:

1. Create a new file in `tools` (e.g., `tools/users.go`), possibly use existing code as a template
2. Add API route and response specifics to the new file
3. Add the tool to `func newServer()` in `cmd/main.go`

## Usage

### Quick Start

1. Get your [Netbird API token](https://docs.netbird.io/api/guides/authentication) from the Netbird management console.

2. Install the `mcp-netbird` binary using one of the installation methods above. Make sure the binary is in your PATH.

3. Configure the server using one of three methods:

   **Method 1: Environment Variables (Traditional)**
   
   Add the server configuration to your client configuration file. E.g., for Codeium Windsurf add the following to `~/.codeium/windsurf/mcp_config.json`:

   ```json
   {
     "mcpServers": {
       "netbird": {
         "command": "mcp-netbird",
         "args": [],
         "env": {
           "NETBIRD_API_TOKEN": "<your-api-token>",
           "NETBIRD_HOST": "api.netbird.io"
         }
       }
     }
   }
   ```

   **Method 2: Command-Line Arguments**
   
   Pass configuration directly via CLI flags:

   ```json
   {
     "mcpServers": {
       "netbird": {
         "command": "mcp-netbird",
         "args": [
           "--api-token", "<your-api-token>",
           "--api-host", "api.netbird.io"
         ]
       }
     }
   }
   ```

   **Method 3: HTTP Headers (SSE Mode Only)**
   
   For Docker MCP Toolkit or other SSE-based clients:

   ```json
   {
     "mcpServers": {
       "netbird": {
         "url": "http://localhost:8001/sse",
         "headers": {
           "X-Netbird-API-Token": "<your-api-token>",
           "X-Netbird-Host": "api.netbird.io"
         }
       }
     }
   }
   ```

   **Note**: `NETBIRD_HOST` (or `--api-host` or `X-Netbird-Host`) defaults to `api.netbird.io` if not specified. For self-hosted instances, set it to your NetBird API host (e.g., `api.yourdomain.com`).

   For more information on how to add a similar configuration to Claude Desktop, see [here](https://modelcontextprotocol.io/quickstart/user).

   > Note: if you see something along the lines of `[netbird] [error] spawn mcp-netbird ENOENT` in Claude Desktop logs, you need to specify the full path to `mcp-netbird`. On macOS Claude Logs are in `~/Library/Logs/Claude`.

4. Try asking questions along the lines of "Can you explain my Netbird peers, groups and policies to me?"
   
![claude-desktop-mcp-netbird](https://github.com/user-attachments/assets/094614cd-9399-4c90-adb3-06ae67c604e4)

### Working with Policies

#### Policy Rule Format

When creating or updating policies, rules must follow this format:

```json
{
  "name": "Rule Name",
  "description": "Optional description",
  "enabled": true,
  "action": "accept",
  "bidirectional": true,
  "protocol": "all",
  "sources": ["group-id-1", "group-id-2"],
  "destinations": ["group-id-3", "group-id-4"],
  "port_ranges": [
    {"start": 80, "end": 80},
    {"start": 443, "end": 443}
  ]
}
```

**Required Fields:**
- `name` (string): Rule name
- `enabled` (boolean): Whether rule is active
- `action` (string): Either "accept" or "drop"
- `bidirectional` (boolean): Whether traffic is allowed in both directions
- `protocol` (string): One of "tcp", "udp", "icmp", or "all"
- At least one source: `sources` (array of group IDs) OR `sourceResource` (object with id and type)
- At least one destination: `destinations` (array of group IDs) OR `destinationResource` (object with id and type)

**Optional Fields:**
- `description` (string): Rule description
- `port_ranges` (array): Port ranges for TCP/UDP protocols
- `authorized_groups` (object): Map of group IDs to arrays of authorized group IDs

**Important Notes:**
- Sources and destinations must be arrays of group ID strings (e.g., `["d535b93ngf8s73892nng"]`)
- Port ranges are only valid for TCP and UDP protocols
- Port range start must be <= end
- All rules are automatically validated before being sent to the API

#### Example: Simple Policy

```json
{
  "name": "Admin SSH Access",
  "description": "Allow admins to SSH to all servers",
  "enabled": true,
  "rules": [
    {
      "name": "SSH Rule",
      "enabled": true,
      "action": "accept",
      "bidirectional": false,
      "protocol": "tcp",
      "sources": ["admin-group-id"],
      "destinations": ["server-group-id"],
      "port_ranges": [{"start": 22, "end": 22}]
    }
  ]
}
```

#### Example: Complex Policy with Multiple Rules

```json
{
  "name": "Infrastructure Access",
  "description": "Complex access control for infrastructure",
  "enabled": true,
  "rules": [
    {
      "name": "HTTP/HTTPS Access",
      "enabled": true,
      "action": "accept",
      "bidirectional": true,
      "protocol": "tcp",
      "sources": ["user-group-id"],
      "destinations": ["web-server-group-id"],
      "port_ranges": [
        {"start": 80, "end": 80},
        {"start": 443, "end": 443}
      ]
    },
    {
      "name": "Database Access",
      "enabled": true,
      "action": "accept",
      "bidirectional": false,
      "protocol": "tcp",
      "sources": ["app-server-group-id"],
      "destinations": ["database-group-id"],
      "port_ranges": [{"start": 5432, "end": 5432}],
      "authorized_groups": {
        "admin-group-id": ["dba-group-id"]
      }
    }
  ]
}
```

#### Getting Policy Templates

Use the `get_policy_template` tool to get example policy structures:

```bash
# Returns examples of simple and complex policy rules with documentation
get_policy_template
```

### Working with Groups

#### Finding Group Dependencies

Before deleting or modifying a group, find which policies reference it:

```bash
# Find all policies that reference a specific group
list_policies_by_group --group_id "d535b93ngf8s73892nng"
```

Returns:
```json
{
  "group_id": "d535b93ngf8s73892nng",
  "group_name": "Admins",
  "policies": [
    {
      "policy_id": "policy-1",
      "policy_name": "Admin Access",
      "rule_id": "rule-1",
      "rule_name": "Admin Rule",
      "location": "sources"
    }
  ]
}
```

#### Replacing Groups Across Policies

Replace one group with another in all policies:

```bash
# Replace old group with new group in all policies
replace_group_in_policies \
  --old_group_id "old-group-id" \
  --new_group_id "new-group-id"
```

Returns:
```json
{
  "old_group_id": "old-group-id",
  "new_group_id": "new-group-id",
  "updated_policies": ["policy-1", "policy-2"],
  "errors": []
}
```

#### Force Deleting Groups

Delete a group and automatically remove it from all dependent policies:

```bash
# Delete group without checking dependencies (removes from all policies)
delete_netbird_group --group_id "group-id" --force true
```

If `force=false` and dependencies exist, returns an error with the list of dependent policies.

## Docker

Build an image and tag it:

```bash
docker build -t mcp-netbird-sse:v1 -f Dockerfile.sse .
```

Run the image with different configuration methods:

**Option 1: Environment Variables (Traditional)**
```bash
docker run --name mcp-netbird -p 8001:8001 \
  -e NETBIRD_API_TOKEN=<your-api-token> \
  mcp-netbird-sse:v1
```

**Option 2: Command-Line Arguments**
```bash
docker run --name mcp-netbird -p 8001:8001 \
  mcp-netbird-sse:v1 \
  --transport sse \
  --sse-address :8001 \
  --api-token <your-api-token> \
  --api-host api.netbird.io
```

**Option 3: HTTP Headers (Stateless Containers)**
```bash
# Run without credentials (stateless)
docker run --name mcp-netbird -p 8001:8001 mcp-netbird-sse:v1

# Pass credentials via HTTP headers with each request
curl -X POST http://localhost:8001/sse \
  -H "X-Netbird-API-Token: <your-api-token>" \
  -H "X-Netbird-Host: api.netbird.io" \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

**Note:** Option 3 is ideal for Docker MCP Toolkit and other stateless container environments where you don't want to bake credentials into the container.

## ToolHive

[ToolHive](https://github.com/StacklokLabs/toolhive) (thv) is a lightweight utility designed to simplify the deployment and management of MCP servers.

You can use ToolHive to deploy and run Netbird MCP as follows:

1. Install `thv` as described in [ToolHive README](https://github.com/StacklokLabs/toolhive#installation).

2. Add Netbird API token to `thv` secrets:

```bash
thv secret set netbird
```

3. Build an SSE image as described in the Docker section [above](#docker)

4. Start Netbird MCP with `thv run` on port 8080:

```bash
thv run --secret netbird,target=NETBIRD_API_TOKEN --transport sse --name thv-mcp-netbird --port 8080 --target-port 8001 mcp-netbird-sse:v1
```

5. When you want to stop the server, use:

```bash
thv stop thv-mcp-netbird
```

## Development

Contributions are welcome! Please open an issue or submit a pull request if you have any suggestions or improvements.

This project is written in Go. Install Go following the instructions for your platform.

To run the server manually, you can use any of the three configuration methods:

**Option 1: Environment Variables**
```bash
export NETBIRD_API_TOKEN=your-token && \
go run cmd/mcp-netbird/main.go
```

**Option 2: Command-Line Arguments**
```bash
go run cmd/mcp-netbird/main.go --api-token your-token --api-host api.netbird.io
```

**Option 3: SSE Mode with HTTP Headers**
```bash
# Start server without credentials
go run cmd/mcp-netbird/main.go --transport sse --sse-address :8001

# In another terminal, make requests with headers
curl -X POST http://localhost:8001/sse \
  -H "X-Netbird-API-Token: your-token" \
  -H "X-Netbird-Host: api.netbird.io" \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/list"}'
```

Or in SSE mode with CLI arguments:

```bash
go run cmd/mcp-netbird/main.go \
  --transport sse \
  --sse-address :8001 \
  --api-token your-token \
  --api-host api.netbird.io
```

### Debugging

The **MCP Inspector** is an interactive developer tool for testing and debugging MCP servers. Read more about it [here](https://modelcontextprotocol.io/docs/tools/inspector).

Here's how to start the MCP Inspector with different configuration methods:

**Option 1: Environment Variables**
```bash
export NETBIRD_API_TOKEN=your-token && \
npx @modelcontextprotocol/inspector
```

**Option 2: Command-Line Arguments**
```bash
npx @modelcontextprotocol/inspector
# Then in the UI, specify: mcp-netbird --api-token your-token --api-host api.netbird.io
```

**Option 3: SSE Mode with HTTP Headers**
```bash
# Start the server
go run cmd/mcp-netbird/main.go --transport sse --sse-address :8001

# In the MCP Inspector UI, configure SSE transport with headers:
# URL: http://localhost:8001/sse
# Headers: X-Netbird-API-Token: your-token
```

Netbird MCP Server can then be tested with either `stdio` or `SSE` transport type. For `stdio` specify the full path to `mcp-netbird` in the UI.

### Testing

**TODO: add more tests**

### Linting

To lint the code, run:

```bash
make lint
```

## License

This project is licensed under the [Apache License, Version 2.0](LICENSE).

**Copyright 2025-2026 XNet Inc.**  
**Copyright 2025-2026 Joshua S. Doucette**

### Attribution

This project was originally derived from the MCP Server for Grafana (https://github.com/grafana/mcp-grafana) developed by Grafana Labs. The current codebase has been substantially modified and extended.

This project uses MCP Go (https://github.com/mark3labs/mcp-go) developed by Mark III Labs.

---

**Maintained by XNet Inc. | Lead Developer: Joshua S. Doucette**