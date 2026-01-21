# Requirements Document

## Introduction

This document specifies requirements for adding stateless configuration support to mcp-netbird, enabling deployment in Docker MCP Toolkit environments where environment variables are not suitable. The feature will support multiple configuration methods with a clear priority order while maintaining backward compatibility.

## Glossary

- **MCP_Server**: The Model Context Protocol server implementation (mcp-netbird)
- **Configuration_Loader**: Component responsible for loading and validating configuration from multiple sources
- **Context_Function**: Function that injects configuration into request context (StdioContextFunc or SSEContextFunc)
- **SSE_Mode**: Server-Sent Events transport mode for HTTP-based access (--transport sse)
- **Stdio_Mode**: Standard input/output transport mode for MCP clients (--transport stdio, default)
- **API_Token**: Authentication token for Netbird API access (NETBIRD_API_TOKEN)
- **API_Host**: Netbird API hostname without protocol (NETBIRD_HOST, default: api.netbird.io)
- **HTTP_Header**: Request header in SSE mode containing configuration
- **Environment_Variable**: System environment variable containing configuration
- **Stateless_Container**: Docker container that does not persist configuration between restarts
- **NetbirdClient**: HTTP client that makes authenticated requests to Netbird API

## Requirements

### Requirement 1: Command-Line Argument Support

**User Story:** As a developer, I want to pass API configuration via command-line arguments, so that I can configure the server without environment variables.

#### Acceptance Criteria

1. WHEN the MCP_Server starts with --api-token flag, THE Configuration_Loader SHALL use the provided API_Token value
2. WHEN the MCP_Server starts with --api-host flag, THE Configuration_Loader SHALL use the provided API_Host value
3. WHEN both --api-token and --api-host flags are provided, THE MCP_Server SHALL function in both Stdio_Mode and SSE_Mode
4. WHEN the --api-token flag value is empty or whitespace-only, THE MCP_Server SHALL return a descriptive error message
5. WHEN the --api-host flag value is empty or whitespace-only, THE MCP_Server SHALL return a descriptive error message
6. WHEN the --api-host flag contains a protocol prefix (http:// or https://), THE Configuration_Loader SHALL strip the protocol and use only the hostname

### Requirement 2: HTTP Header Support for SSE Mode

**User Story:** As a Docker MCP Toolkit operator, I want to pass API configuration via HTTP headers, so that I can run stateless containers without pre-configured credentials.

#### Acceptance Criteria

1. WHEN an SSE_Mode request includes X-Netbird-API-Token header, THE Configuration_Loader SHALL use the header value for API_Token
2. WHEN an SSE_Mode request includes X-Netbird-Host header, THE Configuration_Loader SHALL use the header value for API_Host
3. WHEN both required HTTP_Headers are present in SSE_Mode, THE MCP_Server SHALL process the request successfully
4. WHEN X-Netbird-API-Token header is missing in SSE_Mode and no other source provides it, THE MCP_Server SHALL return HTTP 401 with descriptive error message
5. WHEN X-Netbird-Host header contains a protocol prefix (http:// or https://), THE Configuration_Loader SHALL strip the protocol and use only the hostname
6. WHILE in Stdio_Mode, THE Context_Function SHALL not attempt to read HTTP_Headers

### Requirement 3: Configuration Priority Order

**User Story:** As a system administrator, I want a clear configuration priority order, so that I can override defaults predictably.

#### Acceptance Criteria

1. WHEN multiple configuration sources provide API_Token, THE Configuration_Manager SHALL use CLI_Parser value over HTTP_Header value
2. WHEN multiple configuration sources provide API_Token, THE Configuration_Manager SHALL use HTTP_Header value over Environment_Variable value
3. WHEN multiple configuration sources provide API_Host, THE Configuration_Manager SHALL use CLI_Parser value over HTTP_Header value
4. WHEN multiple configuration sources provide API_Host, THE Configuration_Manager SHALL use HTTP_Header value over Environment_Variable value
5. THE Configuration_Manager SHALL apply the priority order: CLI arguments > HTTP headers > Environment variables

### Requirement 4: Backward Compatibility

**User Story:** As an existing user, I want my current environment variable configuration to continue working, so that I don't need to change my deployment.

#### Acceptance Criteria

1. WHEN only NETBIRD_API_TOKEN environment variable is set, THE MCP_Server SHALL use it for API_Token
2. WHEN only NETBIRD_HOST environment variable is set, THE MCP_Server SHALL use it for API_Host
3. WHEN the MCP_Server runs with environment variables and no CLI arguments or HTTP headers, THE MCP_Server SHALL function identically to the previous version
4. WHEN all existing tests run against the new implementation, THE MCP_Server SHALL pass all 137 tests

### Requirement 5: Configuration Validation

**User Story:** As a developer, I want clear validation errors, so that I can quickly identify configuration problems.

#### Acceptance Criteria

1. WHEN API_Token is empty or whitespace-only, THE Configuration_Manager SHALL return an error indicating missing API token
2. WHEN API_Host is empty or whitespace-only, THE Configuration_Manager SHALL return an error indicating missing API host
3. WHEN API_Host is not a valid URL format, THE Configuration_Manager SHALL return an error indicating invalid URL format
4. WHEN configuration validation fails, THE error message SHALL specify which configuration parameter failed and why
5. WHEN configuration validation succeeds, THE Configuration_Manager SHALL provide the validated configuration to the MCP_Server

### Requirement 6: Documentation Updates

**User Story:** As a new user, I want comprehensive documentation, so that I can understand all configuration methods.

#### Acceptance Criteria

1. THE documentation SHALL include examples of CLI argument usage for both Stdio_Mode and SSE_Mode
2. THE documentation SHALL include examples of HTTP_Header usage for SSE_Mode
3. THE documentation SHALL include examples of Environment_Variable usage
4. THE documentation SHALL explain the configuration priority order with examples
5. THE documentation SHALL include troubleshooting guidance for common configuration errors

### Requirement 7: Test Coverage

**User Story:** As a maintainer, I want comprehensive test coverage, so that I can ensure all configuration methods work correctly.

#### Acceptance Criteria

1. WHEN tests run for CLI argument configuration, THE test suite SHALL verify both valid and invalid inputs
2. WHEN tests run for HTTP_Header configuration, THE test suite SHALL verify both valid and invalid headers
3. WHEN tests run for Environment_Variable configuration, THE test suite SHALL verify both valid and invalid values
4. WHEN tests run for configuration priority, THE test suite SHALL verify all priority combinations
5. WHEN tests run for backward compatibility, THE test suite SHALL verify all existing functionality remains unchanged
