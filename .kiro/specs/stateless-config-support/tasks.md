# Implementation Plan: Stateless Configuration Support

## Overview

This implementation plan breaks down the stateless configuration support feature into discrete coding tasks. The approach follows an incremental pattern: create core configuration structures, add CLI support, add HTTP header support, update context functions, add validation, and finally integrate everything with comprehensive testing.

## Tasks

- [x] 1. Create configuration structures and loader
  - Create `Config` struct with `APIToken` and `APIHost` fields
  - Create `ConfigLoader` struct with CLI argument fields
  - Implement `NewConfigLoader(cliToken, cliHost string)` constructor
  - Implement `LoadConfig(httpToken, httpHost string) (*Config, error)` method with priority order logic
  - Add helper function to strip protocol prefixes from hostnames
  - _Requirements: 1.1, 1.2, 1.6, 2.1, 2.2, 2.5, 3.1, 3.2, 3.3, 3.4_

- [ ]* 1.1 Write property test for configuration loading
  - **Property 1: CLI Configuration Loading**
  - **Validates: Requirements 1.1, 1.2**

- [ ]* 1.2 Write property test for protocol prefix normalization
  - **Property 3: Protocol Prefix Normalization**
  - **Validates: Requirements 1.6, 2.5**

- [ ]* 1.3 Write property test for priority order
  - **Property 5: CLI Priority Over Headers**
  - **Property 6: Header Priority Over Environment**
  - **Validates: Requirements 3.1, 3.2, 3.3, 3.4**

- [x] 2. Implement configuration validation
  - Create `ValidateConfig(cfg *Config) error` function
  - Add validation for empty/whitespace-only API token
  - Add validation for empty/whitespace-only API host
  - Add validation for invalid URL format in API host
  - Ensure error messages specify which parameter failed and why
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ]* 2.1 Write property test for whitespace validation
  - **Property 4: Whitespace Validation**
  - **Validates: Requirements 1.4, 1.5, 5.1, 5.2**

- [ ]* 2.2 Write property test for URL format validation
  - **Property 10: Invalid URL Format Rejection**
  - **Validates: Requirements 5.3**

- [ ]* 2.3 Write property test for error message quality
  - **Property 11: Descriptive Error Messages**
  - **Validates: Requirements 5.4**

- [ ]* 2.4 Write property test for valid configuration acceptance
  - **Property 9: Valid Configuration Acceptance**
  - **Validates: Requirements 5.5**

- [x] 3. Add CLI flags to main.go
  - Add `--api-token` flag with empty default
  - Add `--api-host` flag with empty default
  - Parse flags and store values in variables
  - Create global `ConfigLoader` instance with CLI flag values
  - _Requirements: 1.1, 1.2_

- [ ]* 3.1 Write unit tests for CLI flag parsing
  - Test flag parsing with valid values
  - Test flag parsing with empty values
  - Test flag parsing with whitespace values
  - _Requirements: 1.1, 1.2, 1.4, 1.5_

- [x] 4. Update stdio mode context function
  - Modify `ExtractNetbirdInfoFromEnv` to use `ConfigLoader`
  - Load configuration from CLI arguments and environment variables only
  - Call `ValidateConfig` and log warnings for validation errors
  - Inject validated configuration into context
  - Ensure no HTTP header access in stdio mode
  - _Requirements: 1.1, 1.2, 2.6, 4.1, 4.2_

- [ ]* 4.1 Write property test for stdio mode header isolation
  - **Property 8: Stdio Mode Header Isolation**
  - **Validates: Requirements 2.6**

- [ ]* 4.2 Write property test for environment variable fallback
  - **Property 7: Environment Variable Fallback**
  - **Validates: Requirements 4.1, 4.2**

- [x] 5. Update SSE mode context function
  - Modify `ExtractNetbirdInfoFromEnvSSE` to use `ConfigLoader`
  - Extract `X-Netbird-API-Token` and `X-Netbird-Host` headers from request
  - Load configuration from CLI arguments, HTTP headers, and environment variables
  - Call `ValidateConfig` and return HTTP 401 for missing API token
  - Inject validated configuration into context
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ]* 5.1 Write property test for HTTP header loading
  - **Property 2: HTTP Header Configuration Loading**
  - **Validates: Requirements 2.1, 2.2**

- [ ]* 5.2 Write unit tests for SSE authentication errors
  - Test HTTP 401 response when API token is missing
  - Test error message content
  - Test successful processing with valid headers
  - _Requirements: 2.3, 2.4_

- [x] 6. Checkpoint - Ensure all tests pass
  - Run all existing tests to verify backward compatibility
  - Run all new property tests
  - Run all new unit tests
  - Ensure all tests pass, ask the user if questions arise

- [x] 7. Update NetbirdClient to use context configuration
  - Modify `NewNetbirdClient()` to accept context parameter
  - Extract API host from context instead of environment variable
  - Update baseURL construction to use context-provided host
  - Ensure backward compatibility with existing code
  - _Requirements: 4.1, 4.2, 4.3_

- [ ]* 7.1 Write unit tests for NetbirdClient configuration
  - Test client creation with context configuration
  - Test baseURL construction with various hosts
  - Test backward compatibility with environment variables
  - _Requirements: 4.1, 4.2_

- [x] 8. Integration and multi-mode support
  - Verify server starts successfully with CLI arguments in stdio mode
  - Verify server starts successfully with CLI arguments in SSE mode
  - Test configuration changes between SSE requests
  - Ensure both modes work with all configuration sources
  - _Requirements: 1.3, 2.3_

- [ ]* 8.1 Write property test for multi-mode CLI support
  - **Property 12: Multi-Mode CLI Support**
  - **Validates: Requirements 1.3, 2.3**

- [ ]* 8.2 Write integration tests for both transport modes
  - Test stdio mode with CLI arguments
  - Test SSE mode with HTTP headers
  - Test SSE mode with CLI arguments
  - Test configuration priority in SSE mode
  - _Requirements: 1.3, 2.3, 3.1, 3.2, 3.3, 3.4_

- [x] 9. Update documentation
  - Add CLI argument examples to README.md
  - Add HTTP header examples for SSE mode
  - Document configuration priority order
  - Add troubleshooting section for common errors
  - Include examples for all three configuration methods
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 10. Final checkpoint - Verify all requirements
  - Run complete test suite (existing + new tests)
  - Verify all 137 existing tests still pass
  - Verify all new property tests pass (minimum 100 iterations each)
  - Verify all new unit tests pass
  - Test backward compatibility scenarios manually
  - Ensure all tests pass, ask the user if questions arise

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests use `gopter` library with minimum 100 iterations
- All property tests are tagged with feature name and property number
- Backward compatibility is critical - all existing tests must pass
- Configuration priority order: CLI > HTTP headers > Environment variables
- Protocol prefixes are automatically stripped from API host values
- SSE mode returns HTTP 401 for missing authentication
- Stdio mode logs warnings but continues execution

