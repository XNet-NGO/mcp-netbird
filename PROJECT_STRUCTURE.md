# MCP-NetBird Project Structure

**Copyright 2025-2026 XNet Inc.**  
**Lead Developer: Joshua S. Doucette**

---

## Directory Structure

```
mcp-netbird/
├── cmd/
│   └── mcp-netbird/          # Main application entry point
│       └── main.go
├── tools/                     # NetBird API tool implementations
│   ├── account.go            # Account management
│   ├── groups.go             # Group management + helper functions
│   ├── policies.go           # Policy management + validation
│   ├── peers.go              # Peer management
│   ├── networks.go           # Network management
│   ├── network_resources.go # Network resource management
│   ├── network_routers.go   # Network router management
│   ├── nameservers.go       # DNS nameserver management
│   ├── routes.go            # Route management (deprecated)
│   ├── setup_keys.go        # Setup key management
│   ├── users.go             # User management
│   ├── posture_checks.go    # Posture check management
│   ├── ingress_ports.go     # Port allocation management
│   └── *_test.go            # Unit and property-based tests
├── docs/                     # Documentation
│   ├── ERROR_HANDLING_GUIDE.md
│   ├── HELPER_FUNCTIONS_GUIDE.md
│   └── NETBIRD_FULL_SYSTEM_REVIEW.md
├── archive/                  # Historical documents
│   ├── BUGFIX_SUMMARY.md
│   ├── IMPLEMENTATION_SUMMARY.md
│   ├── NETBIRD_CONFIG_REVIEW.md
│   ├── NETBIRD_FIXES_APPLIED.md
│   ├── NETBIRD_MCP_DEBUG_REPORT.md
│   ├── DOCUMENTATION_COMPLETE.md
│   ├── LICENSING_UPDATE_SUMMARY.md
│   └── API.txt
├── .kiro/                    # Kiro AI specs and configuration
│   └── specs/
│       ├── api-alignment-improvements/
│       └── mcp-netbird-improvements/
├── mcpnetbird.go            # Core HTTP client and context
├── tools.go                 # Tool registration
├── README.md                # Main documentation
├── LICENSE                  # Apache License 2.0
├── NOTICE                   # Attribution notices
├── AUTHORS                  # Author and contributor list
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── Makefile                 # Build automation
├── Dockerfile               # Docker build (stdio)
├── Dockerfile.sse           # Docker build (SSE)
└── smithery.yaml            # Smithery configuration
```

---

## Core Components

### 1. Main Application (`cmd/mcp-netbird/main.go`)
- Server initialization
- Transport configuration (stdio/SSE)
- Tool registration

### 2. HTTP Client (`mcpnetbird.go`)
- NetBird API client
- Context management
- Authentication handling

### 3. Tools (`tools/`)
Each tool file provides CRUD operations for a specific NetBird resource:
- List, Get, Create, Update, Delete operations
- Input validation
- Error handling
- MCP tool registration

### 4. Helper Functions (`tools/groups.go`, `tools/policies.go`)
- `ListPoliciesByGroup` - Find policy dependencies
- `ReplaceGroupInPolicies` - Group consolidation
- `DeleteGroupForce` - Force delete with cleanup
- `GetPolicyTemplate` - Policy examples
- `FormatRuleForAPI` - Policy rule formatting
- `ValidatePolicyRules` - Policy validation

---

## Documentation

### User Documentation (`docs/`)
- **ERROR_HANDLING_GUIDE.md** - Comprehensive error handling patterns
- **HELPER_FUNCTIONS_GUIDE.md** - Helper function usage and workflows
- **NETBIRD_FULL_SYSTEM_REVIEW.md** - Complete system review and analysis

### Developer Documentation
- **README.md** - Installation, configuration, usage examples
- **AUTHORS** - Author and contributor information
- **LICENSE** - Apache License 2.0 with attribution
- **NOTICE** - Copyright and attribution notices

### Historical Documentation (`archive/`)
- Bug fix summaries
- Implementation notes
- Debug reports
- Configuration reviews

---

## Testing

### Test Files (`tools/*_test.go`)
- **Unit Tests**: Specific functionality validation
- **Property-Based Tests**: Universal correctness properties (100+ iterations)
- **Total**: 137 tests passing

### Test Coverage
- Account management: 18 tests
- Group management: 17 tests (including helper functions)
- Policy management: 35 tests (including validation)
- Network management: 12 tests
- Network resources: 8 tests
- Network routers: 9 tests
- Peers: 6 tests
- Port allocations: 10 tests
- Setup keys: 6 tests
- Nameservers: 1 test

---

## Build and Deployment

### Build Commands
```bash
# Build binary
make build

# Install to $GOPATH/bin
make install

# Run tests
make test

# Lint code
make lint
```

### Docker
```bash
# Build stdio image
docker build -t mcp-netbird:latest .

# Build SSE image
docker build -t mcp-netbird-sse:latest -f Dockerfile.sse .
```

---

## Configuration

### Environment Variables
- `NETBIRD_API_TOKEN` (required) - NetBird API authentication token
- `NETBIRD_HOST` (optional) - NetBird API host (default: api.netbird.io)

### MCP Client Configuration
See README.md for client-specific configuration examples.

---

## Development Workflow

### Adding New Tools
1. Create new file in `tools/` (e.g., `tools/newresource.go`)
2. Define data structures
3. Implement CRUD operations
4. Add tool registration function
5. Register in `cmd/mcp-netbird/main.go`
6. Write tests in `tools/newresource_test.go`

### Adding Helper Functions
1. Add function to appropriate tool file
2. Expose as MCP tool if needed
3. Document in HELPER_FUNCTIONS_GUIDE.md
4. Write unit and property-based tests

---

## Code Statistics

- **Total Lines**: ~5,000+ (excluding tests and docs)
- **Test Lines**: ~2,000+
- **Documentation Lines**: ~1,200+
- **Go Files**: 30+
- **Test Files**: 15+
- **Tools**: 60+ operations across 13 categories

---

## Dependencies

### Direct Dependencies
- `github.com/mark3labs/mcp-go` - MCP server framework
- `github.com/invopop/jsonschema` - JSON schema generation

### Indirect Dependencies
- `github.com/google/uuid` - UUID generation
- `gopkg.in/yaml.v3` - YAML parsing

---

## License and Attribution

**Copyright 2025-2026 XNet Inc.**  
**Copyright 2025-2026 Joshua S. Doucette**

Licensed under Apache License 2.0

Originally derived from MCP Server for Grafana by Grafana Labs.
Uses MCP Go by Mark III Labs.

See LICENSE, NOTICE, and AUTHORS files for complete information.

---

## Maintenance

**Maintained by**: XNet Inc.  
**Lead Developer**: Joshua S. Doucette  
**Status**: Production Ready  
**Version**: 0.1.0

---

## Quick Links

- [README](README.md) - Getting started
- [Error Handling Guide](docs/ERROR_HANDLING_GUIDE.md) - Error patterns
- [Helper Functions Guide](docs/HELPER_FUNCTIONS_GUIDE.md) - Workflows
- [System Review](docs/NETBIRD_FULL_SYSTEM_REVIEW.md) - Complete analysis
- [LICENSE](LICENSE) - License terms
- [AUTHORS](AUTHORS) - Contributors
