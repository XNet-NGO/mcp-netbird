package mcpnetbird

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// TestIntegration_StdioMode_WithCLIArguments tests that the server can start
// successfully in stdio mode with CLI arguments providing configuration.
// This validates Requirements 1.3 (multi-mode CLI support).
func TestIntegration_StdioMode_WithCLIArguments(t *testing.T) {
	// Clear environment variables to ensure we're only testing CLI arguments
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	tests := []struct {
		name      string
		cliToken  string
		cliHost   string
		wantToken string
		wantHost  string
	}{
		{
			name:      "valid CLI arguments",
			cliToken:  "cli-test-token",
			cliHost:   "cli.test.example.com",
			wantToken: "cli-test-token",
			wantHost:  "cli.test.example.com",
		},
		{
			name:      "CLI with protocol prefix stripped",
			cliToken:  "cli-test-token",
			cliHost:   "https://cli.test.example.com",
			wantToken: "cli-test-token",
			wantHost:  "cli.test.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up GlobalConfigLoader with CLI arguments
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Create a context using the stdio context function
			ctx := ExtractNetbirdInfoFromEnv(context.Background())

			// Verify the configuration was loaded correctly
			gotToken := NetbirdAPIKeyFromContext(ctx)
			gotHost := NetbirdAPIHostFromContext(ctx)

			if gotToken != tt.wantToken {
				t.Errorf("Token = %v, want %v", gotToken, tt.wantToken)
			}
			if gotHost != tt.wantHost {
				t.Errorf("Host = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

// TestIntegration_SSEMode_WithCLIArguments tests that the server can start
// successfully in SSE mode with CLI arguments providing configuration.
// This validates Requirements 1.3, 2.3 (multi-mode CLI support).
func TestIntegration_SSEMode_WithCLIArguments(t *testing.T) {
	// Clear environment variables to ensure we're only testing CLI arguments
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	tests := []struct {
		name      string
		cliToken  string
		cliHost   string
		wantToken string
		wantHost  string
	}{
		{
			name:      "valid CLI arguments",
			cliToken:  "cli-test-token",
			cliHost:   "cli.test.example.com",
			wantToken: "cli-test-token",
			wantHost:  "cli.test.example.com",
		},
		{
			name:      "CLI with protocol prefix stripped",
			cliToken:  "cli-test-token",
			cliHost:   "https://cli.test.example.com",
			wantToken: "cli-test-token",
			wantHost:  "cli.test.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up GlobalConfigLoader with CLI arguments
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Create a mock HTTP request (no headers, relying on CLI args)
			req := httptest.NewRequest("GET", "/", nil)

			// Create a context using the SSE context function
			ctx := ExtractNetbirdInfoFromEnvSSE(context.Background(), req)

			// Verify the configuration was loaded correctly
			gotToken := NetbirdAPIKeyFromContext(ctx)
			gotHost := NetbirdAPIHostFromContext(ctx)

			if gotToken != tt.wantToken {
				t.Errorf("Token = %v, want %v", gotToken, tt.wantToken)
			}
			if gotHost != tt.wantHost {
				t.Errorf("Host = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

// TestIntegration_SSEMode_ConfigurationChanges tests that configuration can
// change between SSE requests using HTTP headers.
// This validates Requirement 2.3 (SSE mode with HTTP headers).
func TestIntegration_SSEMode_ConfigurationChanges(t *testing.T) {
	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Set up GlobalConfigLoader with no CLI arguments (to test header priority)
	GlobalConfigLoader = NewConfigLoader("", "")

	// First request with one set of headers
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.Header.Set("X-Netbird-API-Token", "token-request-1")
	req1.Header.Set("X-Netbird-Host", "host1.example.com")

	ctx1 := ExtractNetbirdInfoFromEnvSSE(context.Background(), req1)
	token1 := NetbirdAPIKeyFromContext(ctx1)
	host1 := NetbirdAPIHostFromContext(ctx1)

	if token1 != "token-request-1" {
		t.Errorf("Request 1: Token = %v, want token-request-1", token1)
	}
	if host1 != "host1.example.com" {
		t.Errorf("Request 1: Host = %v, want host1.example.com", host1)
	}

	// Second request with different headers
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("X-Netbird-API-Token", "token-request-2")
	req2.Header.Set("X-Netbird-Host", "host2.example.com")

	ctx2 := ExtractNetbirdInfoFromEnvSSE(context.Background(), req2)
	token2 := NetbirdAPIKeyFromContext(ctx2)
	host2 := NetbirdAPIHostFromContext(ctx2)

	if token2 != "token-request-2" {
		t.Errorf("Request 2: Token = %v, want token-request-2", token2)
	}
	if host2 != "host2.example.com" {
		t.Errorf("Request 2: Host = %v, want host2.example.com", host2)
	}

	// Third request with protocol prefix in header
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("X-Netbird-API-Token", "token-request-3")
	req3.Header.Set("X-Netbird-Host", "https://host3.example.com")

	ctx3 := ExtractNetbirdInfoFromEnvSSE(context.Background(), req3)
	token3 := NetbirdAPIKeyFromContext(ctx3)
	host3 := NetbirdAPIHostFromContext(ctx3)

	if token3 != "token-request-3" {
		t.Errorf("Request 3: Token = %v, want token-request-3", token3)
	}
	if host3 != "host3.example.com" {
		t.Errorf("Request 3: Host = %v, want host3.example.com (protocol stripped)", host3)
	}
}

// TestIntegration_BothModes_AllConfigurationSources tests that both stdio and
// SSE modes work correctly with all configuration sources (CLI, headers, env vars).
// This validates Requirements 1.3, 2.3, 3.1-3.4 (priority order).
func TestIntegration_BothModes_AllConfigurationSources(t *testing.T) {
	tests := []struct {
		name        string
		cliToken    string
		cliHost     string
		envToken    string
		envHost     string
		httpToken   string
		httpHost    string
		mode        string // "stdio" or "sse"
		wantToken   string
		wantHost    string
		description string
	}{
		{
			name:        "stdio: CLI overrides env",
			cliToken:    "cli-token",
			cliHost:     "cli.example.com",
			envToken:    "env-token",
			envHost:     "env.example.com",
			mode:        "stdio",
			wantToken:   "cli-token",
			wantHost:    "cli.example.com",
			description: "CLI arguments should take priority over environment variables in stdio mode",
		},
		{
			name:        "stdio: env fallback when no CLI",
			envToken:    "env-token",
			envHost:     "env.example.com",
			mode:        "stdio",
			wantToken:   "env-token",
			wantHost:    "env.example.com",
			description: "Environment variables should be used when no CLI arguments provided in stdio mode",
		},
		{
			name:        "sse: CLI overrides headers and env",
			cliToken:    "cli-token",
			cliHost:     "cli.example.com",
			envToken:    "env-token",
			envHost:     "env.example.com",
			httpToken:   "http-token",
			httpHost:    "http.example.com",
			mode:        "sse",
			wantToken:   "cli-token",
			wantHost:    "cli.example.com",
			description: "CLI arguments should take priority over HTTP headers and environment variables in SSE mode",
		},
		{
			name:        "sse: headers override env",
			envToken:    "env-token",
			envHost:     "env.example.com",
			httpToken:   "http-token",
			httpHost:    "http.example.com",
			mode:        "sse",
			wantToken:   "http-token",
			wantHost:    "http.example.com",
			description: "HTTP headers should take priority over environment variables in SSE mode",
		},
		{
			name:        "sse: env fallback when no CLI or headers",
			envToken:    "env-token",
			envHost:     "env.example.com",
			mode:        "sse",
			wantToken:   "env-token",
			wantHost:    "env.example.com",
			description: "Environment variables should be used when no CLI arguments or HTTP headers provided in SSE mode",
		},
		{
			name:        "sse: partial CLI with header fallback",
			cliToken:    "cli-token",
			httpHost:    "http.example.com",
			mode:        "sse",
			wantToken:   "cli-token",
			wantHost:    "http.example.com",
			description: "Should use CLI token and HTTP header host when partially configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			if tt.envToken != "" {
				os.Setenv(netbirdAPIEnvVar, tt.envToken)
			} else {
				os.Unsetenv(netbirdAPIEnvVar)
			}
			if tt.envHost != "" {
				os.Setenv(netbirdHostEnvVar, tt.envHost)
			} else {
				os.Unsetenv(netbirdHostEnvVar)
			}
			defer func() {
				os.Unsetenv(netbirdAPIEnvVar)
				os.Unsetenv(netbirdHostEnvVar)
			}()

			// Set up GlobalConfigLoader with CLI arguments
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			var ctx context.Context
			if tt.mode == "stdio" {
				// Stdio mode: no HTTP headers
				ctx = ExtractNetbirdInfoFromEnv(context.Background())
			} else {
				// SSE mode: with HTTP headers
				req := httptest.NewRequest("GET", "/", nil)
				if tt.httpToken != "" {
					req.Header.Set("X-Netbird-API-Token", tt.httpToken)
				}
				if tt.httpHost != "" {
					req.Header.Set("X-Netbird-Host", tt.httpHost)
				}
				ctx = ExtractNetbirdInfoFromEnvSSE(context.Background(), req)
			}

			// Verify the configuration
			gotToken := NetbirdAPIKeyFromContext(ctx)
			gotHost := NetbirdAPIHostFromContext(ctx)

			if gotToken != tt.wantToken {
				t.Errorf("%s: Token = %v, want %v", tt.description, gotToken, tt.wantToken)
			}
			if gotHost != tt.wantHost {
				t.Errorf("%s: Host = %v, want %v", tt.description, gotHost, tt.wantHost)
			}
		})
	}
}

// TestIntegration_SSEMode_ConcurrentRequests tests that multiple concurrent
// SSE requests with different configurations are handled correctly.
// This validates that configuration is properly isolated per request.
func TestIntegration_SSEMode_ConcurrentRequests(t *testing.T) {
	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Set up GlobalConfigLoader with no CLI arguments
	GlobalConfigLoader = NewConfigLoader("", "")

	// Create multiple requests with different configurations
	requests := []struct {
		token string
		host  string
	}{
		{"token-1", "host1.example.com"},
		{"token-2", "host2.example.com"},
		{"token-3", "host3.example.com"},
		{"token-4", "host4.example.com"},
		{"token-5", "host5.example.com"},
	}

	// Process requests concurrently
	done := make(chan bool, len(requests))
	for i, reqData := range requests {
		go func(idx int, token, host string) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Netbird-API-Token", token)
			req.Header.Set("X-Netbird-Host", host)

			ctx := ExtractNetbirdInfoFromEnvSSE(context.Background(), req)

			gotToken := NetbirdAPIKeyFromContext(ctx)
			gotHost := NetbirdAPIHostFromContext(ctx)

			if gotToken != token {
				t.Errorf("Request %d: Token = %v, want %v", idx, gotToken, token)
			}
			if gotHost != host {
				t.Errorf("Request %d: Host = %v, want %v", idx, gotHost, host)
			}

			done <- true
		}(i, reqData.token, reqData.host)
	}

	// Wait for all requests to complete with timeout
	timeout := time.After(5 * time.Second)
	for i := 0; i < len(requests); i++ {
		select {
		case <-done:
			// Request completed successfully
		case <-timeout:
			t.Fatal("Timeout waiting for concurrent requests to complete")
		}
	}
}

// TestIntegration_ComposedContextFunctions tests that the composed context
// functions work correctly in both modes.
func TestIntegration_ComposedContextFunctions(t *testing.T) {
	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	t.Run("stdio composed function", func(t *testing.T) {
		// Set up GlobalConfigLoader
		GlobalConfigLoader = NewConfigLoader("composed-cli-token", "composed.example.com")

		// Use the composed stdio context function
		ctx := ComposedStdioContextFunc(context.Background())

		// Verify configuration
		gotToken := NetbirdAPIKeyFromContext(ctx)
		gotHost := NetbirdAPIHostFromContext(ctx)

		if gotToken != "composed-cli-token" {
			t.Errorf("Token = %v, want composed-cli-token", gotToken)
		}
		if gotHost != "composed.example.com" {
			t.Errorf("Host = %v, want composed.example.com", gotHost)
		}
	})

	t.Run("sse composed function", func(t *testing.T) {
		// Set up GlobalConfigLoader
		GlobalConfigLoader = NewConfigLoader("", "")

		// Create request with headers
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Netbird-API-Token", "composed-http-token")
		req.Header.Set("X-Netbird-Host", "composed-http.example.com")

		// Use the composed SSE context function
		ctx := ComposedSSEContextFunc(context.Background(), req)

		// Verify configuration
		gotToken := NetbirdAPIKeyFromContext(ctx)
		gotHost := NetbirdAPIHostFromContext(ctx)

		if gotToken != "composed-http-token" {
			t.Errorf("Token = %v, want composed-http-token", gotToken)
		}
		if gotHost != "composed-http.example.com" {
			t.Errorf("Host = %v, want composed-http.example.com", gotHost)
		}
	})
}

// TestIntegration_ServerCreation tests that an MCP server can be created
// and configured with both transport modes.
func TestIntegration_ServerCreation(t *testing.T) {
	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	t.Run("stdio server creation", func(t *testing.T) {
		// Set up GlobalConfigLoader
		GlobalConfigLoader = NewConfigLoader("stdio-server-token", "stdio.server.example.com")

		// Create a basic MCP server
		s := server.NewMCPServer("test-server", "1.0.0")

		// Create stdio server with context function
		srv := server.NewStdioServer(s)
		srv.SetContextFunc(ComposedStdioContextFunc)

		// Verify server was created successfully
		if srv == nil {
			t.Fatal("Failed to create stdio server")
		}

		// Note: We don't actually start the server in tests as it would block
		// The important part is that the server can be created and configured
	})

	t.Run("sse server creation", func(t *testing.T) {
		// Set up GlobalConfigLoader
		GlobalConfigLoader = NewConfigLoader("sse-server-token", "sse.server.example.com")

		// Create a basic MCP server
		s := server.NewMCPServer("test-server", "1.0.0")

		// Create SSE server with context function
		srv := server.NewSSEServer(s,
			server.WithSSEContextFunc(ComposedSSEContextFunc),
		)

		// Verify server was created successfully
		if srv == nil {
			t.Fatal("Failed to create SSE server")
		}

		// Note: We don't actually start the server in tests as it would block
		// The important part is that the server can be created and configured
	})
}
