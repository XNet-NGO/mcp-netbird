# Design Document: Stateless Configuration Support

## Overview

This design adds stateless configuration support to mcp-netbird by enabling configuration via CLI arguments and HTTP headers, in addition to the existing environment variable support. The implementation maintains backward compatibility while providing a clear priority order: CLI arguments > HTTP headers > environment variables.

The design introduces a centralized configuration loading mechanism that validates inputs and provides them to both stdio and SSE transport modes. This enables deployment in Docker MCP Toolkit environments where environment variables are not suitable.

## Architecture

### Configuration Flow

```
┌─────────────────┐
│   CLI Parser    │ (highest priority)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ HTTP Headers    │ (SSE mode only, medium priority)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Environment Vars│ (lowest priority)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Configuration  │
│    Validator    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Context      │
│   Functions     │
└─────────────────┘
```

### Component Interaction

The configuration system consists of three main components:

1. **Configuration Loader**: Aggregates configuration from multiple sources with priority order
2. **Configuration Validator**: Validates configuration values and returns descriptive errors
3. **Context Functions**: Inject validated configuration into request context for both stdio and SSE modes

## Components and Interfaces

### 1. Configuration Structure

```go
// Config holds the Netbird API configuration
type Config struct {
    APIToken string
    APIHost  string
}
```

This structure represents the validated configuration that will be injected into the context.

### 2. Configuration Loader

```go
// ConfigLoader loads configuration from multiple sources with priority order
type ConfigLoader struct {
    cliToken  string
    cliHost   string
}

// NewConfigLoader creates a new configuration loader with CLI arguments
func NewConfigLoader(cliToken, cliHost string) *ConfigLoader

// LoadConfig loads configuration with priority: CLI > HTTP headers > env vars
func (cl *ConfigLoader) LoadConfig(httpToken, httpHost string) (*Config, error)
```

The ConfigLoader is initialized once at startup with CLI arguments, then used by context functions to load configuration for each request (incorporating HTTP headers in SSE mode).

### 3. Configuration Validator

```go
// ValidateConfig validates the configuration and returns descriptive errors
func ValidateConfig(cfg *Config) error
```

Validation rules:
- API token must not be empty or whitespace-only
- API host must not be empty or whitespace-only
- API host must not contain protocol prefix (http:// or https://)
- If protocol prefix exists, it is stripped automatically

### 4. Context Functions

#### Stdio Mode Context Function

```go
// ExtractNetbirdInfoFromEnv is a StdioContextFunc that extracts Netbird configuration
// from CLI arguments and environment variables (no HTTP headers in stdio mode)
var ExtractNetbirdInfoFromEnv server.StdioContextFunc = func(ctx context.Context) context.Context
```

In stdio mode, the context function:
1. Loads configuration from CLI arguments and environment variables
2. Validates the configuration
3. Injects it into the context
4. Logs warnings for missing configuration

#### SSE Mode Context Function

```go
// ExtractNetbirdInfoFromEnvSSE is an SSEContextFunc that extracts Netbird configuration
// from CLI arguments, HTTP headers, and environment variables
var ExtractNetbirdInfoFromEnvSSE server.SSEContextFunc = func(ctx context.Context, req *http.Request) context.Context
```

In SSE mode, the context function:
1. Extracts HTTP headers (X-Netbird-API-Token, X-Netbird-Host)
2. Loads configuration from CLI arguments, HTTP headers, and environment variables
3. Validates the configuration
4. Injects it into the context
5. Returns HTTP 401 error if API token is missing

### 5. CLI Flag Additions

The main.go file will be modified to accept new flags:

```go
var (
    transport  string
    addr       string
    apiToken   string  // new flag
    apiHost    string  // new flag
)

flag.StringVar(&apiToken, "api-token", "", "Netbird API token")
flag.StringVar(&apiHost, "api-host", "", "Netbird API host (without protocol)")
```

These flags are parsed at startup and passed to the ConfigLoader.

## Data Models

### Configuration Priority

The priority order is implemented through a waterfall pattern:

```
1. Check CLI argument value
   ├─ If present and valid → use it
   └─ If absent → continue to step 2

2. Check HTTP header value (SSE mode only)
   ├─ If present and valid → use it
   └─ If absent → continue to step 3

3. Check environment variable value
   ├─ If present and valid → use it
   └─ If absent → use default (for host) or error (for token)
```

### Default Values

- API Host: `api.netbird.io` (if not provided by any source)
- API Token: No default (must be provided by at least one source)

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: CLI Configuration Loading
*For any* valid API token and API host provided via CLI arguments, the configuration loader should successfully load and use those values.
**Validates: Requirements 1.1, 1.2**

### Property 2: HTTP Header Configuration Loading
*For any* valid API token and API host provided via HTTP headers in SSE mode, the configuration loader should successfully load and use those values.
**Validates: Requirements 2.1, 2.2**

### Property 3: Protocol Prefix Normalization
*For any* hostname with or without protocol prefix (http:// or https://), the configuration loader should strip the protocol and produce the same normalized hostname.
**Validates: Requirements 1.6, 2.5**

### Property 4: Whitespace Validation
*For any* string composed entirely of whitespace characters (spaces, tabs, newlines), validation should reject it as an invalid API token or API host and return a descriptive error.
**Validates: Requirements 1.4, 1.5, 5.1, 5.2**

### Property 5: CLI Priority Over Headers
*For any* configuration where both CLI arguments and HTTP headers provide values, the CLI argument values should always be used (higher priority).
**Validates: Requirements 3.1, 3.3**

### Property 6: Header Priority Over Environment
*For any* configuration where both HTTP headers and environment variables provide values, the HTTP header values should always be used (higher priority).
**Validates: Requirements 3.2, 3.4**

### Property 7: Environment Variable Fallback
*For any* valid API token and API host provided only via environment variables (no CLI args or headers), the configuration loader should successfully load and use those values.
**Validates: Requirements 4.1, 4.2**

### Property 8: Stdio Mode Header Isolation
*For any* stdio mode execution, the context function should not attempt to access HTTP request headers and should only use CLI arguments and environment variables.
**Validates: Requirements 2.6**

### Property 9: Valid Configuration Acceptance
*For any* configuration with non-empty API token and valid API host format, validation should succeed and return the validated configuration.
**Validates: Requirements 5.5**

### Property 10: Invalid URL Format Rejection
*For any* API host that is not a valid URL format (e.g., contains invalid characters, malformed structure), validation should reject it and return a descriptive error.
**Validates: Requirements 5.3**

### Property 11: Descriptive Error Messages
*For any* validation failure, the error message should specify which configuration parameter failed (API token or API host) and provide a clear reason for the failure.
**Validates: Requirements 5.4**

### Property 12: Multi-Mode CLI Support
*For any* valid CLI configuration (both --api-token and --api-host provided), the server should successfully start and function in both stdio mode and SSE mode.
**Validates: Requirements 1.3, 2.3**

## Error Handling

### Validation Errors

The system provides specific error messages for each validation failure:

1. **Missing API Token**: "API token is required but not provided"
2. **Empty API Token**: "API token cannot be empty or whitespace-only"
3. **Missing API Host**: "API host is required but not provided"
4. **Empty API Host**: "API host cannot be empty or whitespace-only"
5. **Invalid URL Format**: "API host '{host}' is not a valid URL format: {reason}"
6. **Protocol Prefix Warning**: Automatically stripped with log message "Stripped protocol prefix from API host: {original} -> {normalized}"

### SSE Mode Authentication Errors

When running in SSE mode without a valid API token:
- HTTP Status: 401 Unauthorized
- Response Body: JSON with error message
- Log Entry: "SSE MODE - Missing API token from all sources"

### Stdio Mode Errors

When running in stdio mode without a valid API token:
- Log Warning: "Warning: API token not found in CLI arguments or environment variables"
- Behavior: Continue execution but API calls will fail with authentication errors

### Error Propagation

Validation errors are propagated through the context functions:
1. ConfigLoader.LoadConfig() returns error
2. Context function logs the error
3. In SSE mode: HTTP 401 response sent to client
4. In stdio mode: Warning logged, context created with empty token

## Testing Strategy

### Dual Testing Approach

The implementation will use both unit tests and property-based tests:

**Unit Tests** focus on:
- Specific examples of valid configurations
- Edge cases (empty strings, whitespace variations)
- Error message content verification
- Integration between components
- Backward compatibility with existing tests

**Property-Based Tests** focus on:
- Universal properties across all inputs
- Configuration priority order with random values
- Protocol stripping with various hostname formats
- Validation behavior across input space
- Mode-specific behavior (stdio vs SSE)

### Property-Based Testing Configuration

We will use the `gopter` library for property-based testing in Go:
- Minimum 100 iterations per property test
- Each test tagged with: `// Feature: stateless-config-support, Property N: {property_text}`
- Random generation of:
  - API tokens (various lengths and character sets)
  - API hostnames (valid and invalid formats)
  - Whitespace strings (spaces, tabs, newlines, combinations)
  - Protocol prefixes (http://, https://, mixed case)

### Test Coverage Requirements

1. **CLI Argument Tests**:
   - Valid token and host combinations
   - Empty and whitespace-only values
   - Protocol prefix stripping
   - Priority over other sources

2. **HTTP Header Tests**:
   - Valid header values in SSE mode
   - Missing headers with fallback to env vars
   - Priority over environment variables
   - Header isolation in stdio mode

3. **Environment Variable Tests**:
   - Backward compatibility with existing behavior
   - Fallback when no CLI args or headers present
   - All 137 existing tests must pass

4. **Priority Order Tests**:
   - All combinations of CLI, headers, and env vars
   - Verification of waterfall pattern
   - Edge cases with partial configuration

5. **Validation Tests**:
   - Empty and whitespace-only inputs
   - Invalid URL formats
   - Error message descriptiveness
   - Success cases with valid inputs

### Integration Testing

Integration tests will verify:
1. Server startup with CLI arguments in both modes
2. SSE requests with HTTP headers
3. Stdio mode with environment variables only
4. Configuration changes between requests (SSE mode)
5. Error responses in SSE mode for missing auth

## Implementation Notes

### Backward Compatibility

The implementation maintains 100% backward compatibility:
- Existing environment variable behavior unchanged
- All existing tests pass without modification
- Default values remain the same
- No breaking changes to public APIs

### Performance Considerations

- Configuration loading happens once per request (SSE) or once at startup (stdio)
- Validation is lightweight (string checks and basic URL parsing)
- No performance impact on existing deployments
- HTTP header extraction is only performed in SSE mode

### Security Considerations

- API tokens are never logged in full (only length is logged)
- Configuration validation prevents injection attacks
- Protocol stripping prevents SSRF via protocol smuggling
- Clear error messages don't leak sensitive information

### Logging Strategy

The implementation provides detailed logging for debugging:
- Configuration source used (CLI, header, env var)
- Validation failures with specific reasons
- Mode-specific behavior (stdio vs SSE)
- Protocol stripping operations
- Missing configuration warnings

Log levels:
- INFO: Successful configuration loading
- WARN: Missing configuration with fallback
- ERROR: Validation failures, authentication errors

