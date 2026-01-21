# Changelog

All notable changes to mcp-netbird will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Stateless configuration support with three methods (CLI, HTTP headers, environment variables)
- Configuration priority order: CLI > HTTP headers > Environment variables
- Protocol prefix stripping for API host values
- Comprehensive configuration validation with descriptive error messages
- Docker deployment support with stateless containers
- Multi-platform release support (Linux x86_64/ARM64, Windows x64, macOS x64/ARM64)
- Debian package support for Linux distributions
- GitHub Actions workflow for automated releases
- Docker multi-arch images (linux/amd64, linux/arm64)
- Comprehensive documentation (CONTRIBUTING.md, RELEASE_GUIDE.md)

### Changed
- Updated branding to XNet Inc. and Joshua S. Doucette
- Enhanced README with installation instructions for all platforms
- Improved Docker deployment guide with all configuration methods
- Updated LICENSE with proper copyright notices

### Fixed
- Configuration loading in both stdio and SSE modes
- Context-based configuration for NetbirdClient
- Backward compatibility with environment variables

## [0.1.0] - Initial Release

### Added
- Full CRUD operations for all NetBird API resources
- Peer management (list, get, update, delete)
- Group management (list, get, create, update, delete)
- Policy management (list, get, create, update, delete)
- Network management (list, get, create, update, delete)
- Network resource management (list, get, create, update, delete)
- Network router management (list, get, create, update, delete)
- Nameserver management (list, get, create, update, delete)
- Route management (list, get, create, update, delete)
- Setup key management (list, get, create, update, delete)
- User management (list, get, invite, update, delete)
- Posture check management (list, get, create, update, delete)
- Port allocation management (list, get, create, update, delete)
- Account management (get, update)
- Helper tools:
  - `list_policies_by_group` - Find policies referencing a group
  - `replace_group_in_policies` - Replace group across all policies
  - `get_policy_template` - Get example policy structures
- Support for both stdio and SSE transport modes
- Environment variable configuration
- Docker support with SSE mode
- Comprehensive error handling
- Extensive test coverage (318 tests)

### Attribution
- Originally derived from MCP Server for Grafana by Grafana Labs
- Substantially modified and extended by XNet Inc. and Joshua S. Doucette

---

**Maintained by XNet Inc.**  
**Lead Developer: Joshua S. Doucette**  
**Licensed under Apache License 2.0**
