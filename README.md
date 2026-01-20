# NetBird MCP Server

A comprehensive [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for [NetBird](https://netbird.io/) providing full CRUD operations, policy management, and automation workflows.

**Developed by XNet Inc. and Joshua S. Doucette**

## About

This MCP server provides complete management capabilities for NetBird networks, including:
- Full CRUD operations for all NetBird resources
- Advanced policy management with validation
- Group consolidation and dependency workflows
- Helper functions for common administrative tasks
- Comprehensive error handling and documentation

Originally derived from the MCP Server for Grafana by Grafana Labs, this project has been substantially extended and enhanced.

## Installing

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

### Installing via Smithery

[![smithery badge](https://smithery.ai/badge/@aantti/mcp-netbird)](https://smithery.ai/server/@aantti/mcp-netbird)

To install Netbird MCP Server for Claude Desktop automatically via [Smithery](https://smithery.ai/server/@aantti/mcp-netbird):

```bash
npx -y @smithery/cli install @aantti/mcp-netbird --client claude
```

## Configuration

The server requires the following environment variables:

- `NETBIRD_API_TOKEN`: Your Netbird API token
- `NETBIRD_HOST` (optional): The Netbird API host (default is `api.netbird.io`)

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

3. Add the server configuration to your client configuration file. E.g., for Codeium Windsurf add the following to `~/.codeium/windsurf/mcp_config.json`:

   ```json
   {
     "mcpServers": {
       "netbird": {
         "command": "mcp-netbird",
         "args": [],
         "env": {
           "NETBIRD_API_TOKEN": "<your-api-token>"
         }
       }
     }
   }
   ```

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

Run the image:

```bash
docker run --name mcp-netbird -p 8001:8001 -e NETBIRD_API_TOKEN=<your-api-token> mcp-netbird-sse:v1

```

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

To run the server manually, use:

```bash
export NETBIRD_API_TOKEN=your-token && \
go run cmd/mcp-netbird/main.go
```

Or in SSE mode:

```bash
export NETBIRD_API_TOKEN=your-token && \
go run cmd/mcp-netbird/main.go --transport sse --sse-address :8001
```

### Debugging

The **MCP Inspector** is an interactive developer tool for testing and debugging MCP servers. Read more about it [here](https://modelcontextprotocol.io/docs/tools/inspector).

Here's how to start the MCP Inspector:

```bash
export NETBIRD_API_TOKEN=your-token && \
npx @modelcontextprotocol/inspector
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