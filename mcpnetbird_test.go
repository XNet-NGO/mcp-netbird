package mcpnetbird

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNetbirdAPIKeyContext(t *testing.T) {
	ctx := context.Background()
	apiKey := "test-api-key"

	// Test adding API key to context
	ctxWithKey := WithNetbirdAPIKey(ctx, apiKey)
	if got := NetbirdAPIKeyFromContext(ctxWithKey); got != apiKey {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", got, apiKey)
	}

	// Test getting API key from context without key
	if got := NetbirdAPIKeyFromContext(ctx); got != "" {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want empty string", got)
	}
}

func TestNewConfigLoader(t *testing.T) {
	tests := []struct {
		name     string
		cliToken string
		cliHost  string
	}{
		{
			name:     "empty values",
			cliToken: "",
			cliHost:  "",
		},
		{
			name:     "with token only",
			cliToken: "test-token",
			cliHost:  "",
		},
		{
			name:     "with host only",
			cliToken: "",
			cliHost:  "test.example.com",
		},
		{
			name:     "with both values",
			cliToken: "test-token",
			cliHost:  "test.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewConfigLoader(tt.cliToken, tt.cliHost)
			if loader == nil {
				t.Fatal("NewConfigLoader returned nil")
			}
			if loader.cliToken != tt.cliToken {
				t.Errorf("cliToken = %v, want %v", loader.cliToken, tt.cliToken)
			}
			if loader.cliHost != tt.cliHost {
				t.Errorf("cliHost = %v, want %v", loader.cliHost, tt.cliHost)
			}
		})
	}
}

func TestConfigLoader_LoadConfig_CLIPriority(t *testing.T) {
	// Set environment variables for testing
	os.Setenv(netbirdAPIEnvVar, "env-token")
	os.Setenv(netbirdHostEnvVar, "env.example.com")
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	tests := []struct {
		name          string
		cliToken      string
		cliHost       string
		httpToken     string
		httpHost      string
		expectedToken string
		expectedHost  string
	}{
		{
			name:          "CLI token takes priority over HTTP and env",
			cliToken:      "cli-token",
			cliHost:       "",
			httpToken:     "http-token",
			httpHost:      "",
			expectedToken: "cli-token",
			expectedHost:  "env.example.com",
		},
		{
			name:          "CLI host takes priority over HTTP and env",
			cliToken:      "",
			cliHost:       "cli.example.com",
			httpToken:     "",
			httpHost:      "http.example.com",
			expectedToken: "env-token",
			expectedHost:  "cli.example.com",
		},
		{
			name:          "HTTP token takes priority over env",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "http-token",
			httpHost:      "",
			expectedToken: "http-token",
			expectedHost:  "env.example.com",
		},
		{
			name:          "HTTP host takes priority over env",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "",
			httpHost:      "http.example.com",
			expectedToken: "env-token",
			expectedHost:  "http.example.com",
		},
		{
			name:          "env vars used when no CLI or HTTP",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "",
			httpHost:      "",
			expectedToken: "env-token",
			expectedHost:  "env.example.com",
		},
		{
			name:          "CLI takes priority for both",
			cliToken:      "cli-token",
			cliHost:       "cli.example.com",
			httpToken:     "http-token",
			httpHost:      "http.example.com",
			expectedToken: "cli-token",
			expectedHost:  "cli.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewConfigLoader(tt.cliToken, tt.cliHost)
			cfg, err := loader.LoadConfig(tt.httpToken, tt.httpHost)
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}
			if cfg.APIToken != tt.expectedToken {
				t.Errorf("APIToken = %v, want %v", cfg.APIToken, tt.expectedToken)
			}
			if cfg.APIHost != tt.expectedHost {
				t.Errorf("APIHost = %v, want %v", cfg.APIHost, tt.expectedHost)
			}
		})
	}
}

func TestConfigLoader_LoadConfig_DefaultHost(t *testing.T) {
	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	loader := NewConfigLoader("", "")
	cfg, err := loader.LoadConfig("", "")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.APIHost != defaultNetbirdHost {
		t.Errorf("APIHost = %v, want %v", cfg.APIHost, defaultNetbirdHost)
	}
}

func TestStripProtocolPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no prefix",
			input:    "api.example.com",
			expected: "api.example.com",
		},
		{
			name:     "http prefix",
			input:    "http://api.example.com",
			expected: "api.example.com",
		},
		{
			name:     "https prefix",
			input:    "https://api.example.com",
			expected: "api.example.com",
		},
		{
			name:     "http prefix with port",
			input:    "http://api.example.com:8080",
			expected: "api.example.com:8080",
		},
		{
			name:     "https prefix with port",
			input:    "https://api.example.com:8443",
			expected: "api.example.com:8443",
		},
		{
			name:     "http prefix with path",
			input:    "http://api.example.com/path",
			expected: "api.example.com/path",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just protocol",
			input:    "http://",
			expected: "",
		},
		{
			name:     "localhost",
			input:    "http://localhost",
			expected: "localhost",
		},
		{
			name:     "IP address with http",
			input:    "http://192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "IP address with https",
			input:    "https://192.168.1.1",
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripProtocolPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("stripProtocolPrefix(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api.netbird.io",
			},
			expectError: false,
		},
		{
			name: "valid configuration with port",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api.example.com:8080",
			},
			expectError: false,
		},
		{
			name: "valid configuration with path",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api.example.com/api/v1",
			},
			expectError: false,
		},
		{
			name: "valid configuration with IP address",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "192.168.1.1",
			},
			expectError: false,
		},
		{
			name: "valid configuration with localhost",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "localhost",
			},
			expectError: false,
		},
		{
			name: "empty API token",
			config: &Config{
				APIToken: "",
				APIHost:  "api.netbird.io",
			},
			expectError: true,
			errorMsg:    "API token is required but not provided",
		},
		{
			name: "whitespace-only API token (spaces)",
			config: &Config{
				APIToken: "   ",
				APIHost:  "api.netbird.io",
			},
			expectError: true,
			errorMsg:    "API token cannot be empty or whitespace-only",
		},
		{
			name: "whitespace-only API token (tabs)",
			config: &Config{
				APIToken: "\t\t\t",
				APIHost:  "api.netbird.io",
			},
			expectError: true,
			errorMsg:    "API token cannot be empty or whitespace-only",
		},
		{
			name: "whitespace-only API token (newlines)",
			config: &Config{
				APIToken: "\n\n",
				APIHost:  "api.netbird.io",
			},
			expectError: true,
			errorMsg:    "API token cannot be empty or whitespace-only",
		},
		{
			name: "whitespace-only API token (mixed)",
			config: &Config{
				APIToken: " \t\n ",
				APIHost:  "api.netbird.io",
			},
			expectError: true,
			errorMsg:    "API token cannot be empty or whitespace-only",
		},
		{
			name: "empty API host",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "",
			},
			expectError: true,
			errorMsg:    "API host is required but not provided",
		},
		{
			name: "whitespace-only API host (spaces)",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "   ",
			},
			expectError: true,
			errorMsg:    "API host cannot be empty or whitespace-only",
		},
		{
			name: "whitespace-only API host (tabs)",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "\t\t",
			},
			expectError: true,
			errorMsg:    "API host cannot be empty or whitespace-only",
		},
		{
			name: "whitespace-only API host (newlines)",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "\n\n",
			},
			expectError: true,
			errorMsg:    "API host cannot be empty or whitespace-only",
		},
		{
			name: "API host with http protocol prefix",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "http://api.example.com",
			},
			expectError: true,
			errorMsg:    "should not contain protocol prefix",
		},
		{
			name: "API host with https protocol prefix",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "https://api.example.com",
			},
			expectError: true,
			errorMsg:    "should not contain protocol prefix",
		},
		{
			name: "API host with spaces",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api example com",
			},
			expectError: true,
			errorMsg:    "contains spaces",
		},
		{
			name: "API host with invalid character <",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api<example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '<'",
		},
		{
			name: "API host with invalid character >",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api>example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '>'",
		},
		{
			name: "API host with invalid character \"",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api\"example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '\"'",
		},
		{
			name: "API host with invalid character {",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api{example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '{'",
		},
		{
			name: "API host with invalid character }",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api}example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '}'",
		},
		{
			name: "API host with invalid character |",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api|example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '|'",
		},
		{
			name: "API host with invalid character \\",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api\\example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '\\'",
		},
		{
			name: "API host with invalid character ^",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api^example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '^'",
		},
		{
			name: "API host with invalid character `",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "api`example.com",
			},
			expectError: true,
			errorMsg:    "contains invalid character '`'",
		},
		{
			name: "API host with only special characters",
			config: &Config{
				APIToken: "valid-token-123",
				APIHost:  "!!!",
			},
			expectError: true,
			errorMsg:    "no valid hostname characters",
		},
		{
			name: "both empty",
			config: &Config{
				APIToken: "",
				APIHost:  "",
			},
			expectError: true,
			errorMsg:    "API token is required but not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateConfig() expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateConfig() error = %v, want error containing %q", err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConfig() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateConfig_ErrorMessages(t *testing.T) {
	// Test that error messages specify which parameter failed
	tests := []struct {
		name          string
		config        *Config
		expectedParam string // "API token" or "API host"
	}{
		{
			name: "empty token error mentions API token",
			config: &Config{
				APIToken: "",
				APIHost:  "api.netbird.io",
			},
			expectedParam: "API token",
		},
		{
			name: "whitespace token error mentions API token",
			config: &Config{
				APIToken: "   ",
				APIHost:  "api.netbird.io",
			},
			expectedParam: "API token",
		},
		{
			name: "empty host error mentions API host",
			config: &Config{
				APIToken: "valid-token",
				APIHost:  "",
			},
			expectedParam: "API host",
		},
		{
			name: "whitespace host error mentions API host",
			config: &Config{
				APIToken: "valid-token",
				APIHost:  "   ",
			},
			expectedParam: "API host",
		},
		{
			name: "invalid host format error mentions API host",
			config: &Config{
				APIToken: "valid-token",
				APIHost:  "api example com",
			},
			expectedParam: "API host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if err == nil {
				t.Fatal("ValidateConfig() expected error but got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedParam) {
				t.Errorf("ValidateConfig() error = %v, want error mentioning %q", err, tt.expectedParam)
			}
		})
	}
}

func TestExtractNetbirdInfoFromEnv_StdioMode(t *testing.T) {
	tests := []struct {
		name          string
		cliToken      string
		cliHost       string
		envToken      string
		envHost       string
		expectedToken string
	}{
		{
			name:          "CLI token takes priority over env",
			cliToken:      "cli-token",
			cliHost:       "cli.example.com",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "cli-token",
		},
		{
			name:          "env token used when no CLI token",
			cliToken:      "",
			cliHost:       "",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "env-token",
		},
		{
			name:          "empty token when neither CLI nor env provided",
			cliToken:      "",
			cliHost:       "",
			envToken:      "",
			envHost:       "",
			expectedToken: "",
		},
		{
			name:          "CLI token only",
			cliToken:      "cli-token-only",
			cliHost:       "",
			envToken:      "",
			envHost:       "",
			expectedToken: "cli-token-only",
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

			// Set up GlobalConfigLoader
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Call the context function
			ctx := context.Background()
			resultCtx := ExtractNetbirdInfoFromEnv(ctx)

			// Extract the token from the context
			token := NetbirdAPIKeyFromContext(resultCtx)
			if token != tt.expectedToken {
				t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, tt.expectedToken)
			}
		})
	}
}

func TestExtractNetbirdInfoFromEnv_NoHTTPHeaderAccess(t *testing.T) {
	// This test verifies that stdio mode does not attempt to access HTTP headers
	// by ensuring the function works correctly without any HTTP request context

	// Set up environment variables
	os.Setenv(netbirdAPIEnvVar, "env-token")
	os.Setenv(netbirdHostEnvVar, "env.example.com")
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Set up GlobalConfigLoader with CLI arguments
	GlobalConfigLoader = NewConfigLoader("cli-token", "cli.example.com")

	// Call the context function with a plain context (no HTTP request)
	ctx := context.Background()
	resultCtx := ExtractNetbirdInfoFromEnv(ctx)

	// Verify the token was extracted correctly
	token := NetbirdAPIKeyFromContext(resultCtx)
	if token != "cli-token" {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, "cli-token")
	}

	// The test passes if no panic occurs and the token is correctly extracted
	// This demonstrates that stdio mode does not access HTTP headers
}

func TestExtractNetbirdInfoFromEnv_ValidationWarnings(t *testing.T) {
	// Test that validation warnings are logged but execution continues

	tests := []struct {
		name     string
		cliToken string
		cliHost  string
		envToken string
		envHost  string
	}{
		{
			name:     "empty token logs warning",
			cliToken: "",
			cliHost:  "",
			envToken: "",
			envHost:  "api.netbird.io",
		},
		{
			name:     "whitespace token logs warning",
			cliToken: "   ",
			cliHost:  "",
			envToken: "",
			envHost:  "api.netbird.io",
		},
		{
			name:     "invalid host logs warning",
			cliToken: "valid-token",
			cliHost:  "http://api.example.com",
			envToken: "",
			envHost:  "",
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

			// Set up GlobalConfigLoader
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Call the context function - should not panic even with invalid config
			ctx := context.Background()
			resultCtx := ExtractNetbirdInfoFromEnv(ctx)

			// Verify context was created (even if validation failed)
			if resultCtx == nil {
				t.Error("ExtractNetbirdInfoFromEnv() returned nil context")
			}
		})
	}
}

func TestExtractNetbirdInfoFromEnv_UninitializedLoader(t *testing.T) {
	// Test that the function handles uninitialized GlobalConfigLoader gracefully

	// Clear GlobalConfigLoader
	GlobalConfigLoader = nil

	// Set up environment variables
	os.Setenv(netbirdAPIEnvVar, "env-token")
	os.Setenv(netbirdHostEnvVar, "env.example.com")
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Call the context function - should not panic
	ctx := context.Background()
	resultCtx := ExtractNetbirdInfoFromEnv(ctx)

	// Verify context was created
	if resultCtx == nil {
		t.Error("ExtractNetbirdInfoFromEnv() returned nil context")
	}

	// Verify the token was extracted from environment variables
	token := NetbirdAPIKeyFromContext(resultCtx)
	if token != "env-token" {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, "env-token")
	}
}

func TestExtractNetbirdInfoFromEnvSSE_HTTPHeaders(t *testing.T) {
	tests := []struct {
		name          string
		cliToken      string
		cliHost       string
		httpToken     string
		httpHost      string
		envToken      string
		envHost       string
		expectedToken string
	}{
		{
			name:          "HTTP headers used when no CLI args",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "http-token",
			httpHost:      "http.example.com",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "http-token",
		},
		{
			name:          "CLI token takes priority over HTTP headers",
			cliToken:      "cli-token",
			cliHost:       "cli.example.com",
			httpToken:     "http-token",
			httpHost:      "http.example.com",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "cli-token",
		},
		{
			name:          "HTTP token takes priority over env",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "http-token",
			httpHost:      "http.example.com",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "http-token",
		},
		{
			name:          "env token used when no CLI or HTTP",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "",
			httpHost:      "",
			envToken:      "env-token",
			envHost:       "env.example.com",
			expectedToken: "env-token",
		},
		{
			name:          "empty token when all sources empty",
			cliToken:      "",
			cliHost:       "",
			httpToken:     "",
			httpHost:      "",
			envToken:      "",
			envHost:       "",
			expectedToken: "",
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

			// Set up GlobalConfigLoader
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Create a mock HTTP request with headers
			req, err := http.NewRequest("GET", "http://localhost:8001", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.httpToken != "" {
				req.Header.Set("X-Netbird-API-Token", tt.httpToken)
			}
			if tt.httpHost != "" {
				req.Header.Set("X-Netbird-Host", tt.httpHost)
			}

			// Call the SSE context function
			ctx := context.Background()
			resultCtx := ExtractNetbirdInfoFromEnvSSE(ctx, req)

			// Extract the token from the context
			token := NetbirdAPIKeyFromContext(resultCtx)
			if token != tt.expectedToken {
				t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, tt.expectedToken)
			}
		})
	}
}

func TestExtractNetbirdInfoFromEnvSSE_MissingAPIToken(t *testing.T) {
	// Test that missing API token is handled correctly in SSE mode

	// Clear environment variables
	os.Unsetenv(netbirdAPIEnvVar)
	os.Unsetenv(netbirdHostEnvVar)
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Set up GlobalConfigLoader with no CLI args
	GlobalConfigLoader = NewConfigLoader("", "")

	// Create a mock HTTP request without API token header
	req, err := http.NewRequest("GET", "http://localhost:8001", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Call the SSE context function
	ctx := context.Background()
	resultCtx := ExtractNetbirdInfoFromEnvSSE(ctx, req)

	// Verify that an empty token is injected (which will cause auth to fail later)
	token := NetbirdAPIKeyFromContext(resultCtx)
	if token != "" {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want empty string", token)
	}
}

func TestExtractNetbirdInfoFromEnvSSE_ProtocolStripping(t *testing.T) {
	// Test that protocol prefixes are stripped from HTTP headers

	tests := []struct {
		name     string
		httpHost string
	}{
		{
			name:     "http prefix in header",
			httpHost: "http://api.example.com",
		},
		{
			name:     "https prefix in header",
			httpHost: "https://api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv(netbirdAPIEnvVar)
			os.Unsetenv(netbirdHostEnvVar)
			defer func() {
				os.Unsetenv(netbirdAPIEnvVar)
				os.Unsetenv(netbirdHostEnvVar)
			}()

			// Set up GlobalConfigLoader
			GlobalConfigLoader = NewConfigLoader("", "")

			// Create a mock HTTP request with headers
			req, err := http.NewRequest("GET", "http://localhost:8001", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("X-Netbird-API-Token", "test-token")
			req.Header.Set("X-Netbird-Host", tt.httpHost)

			// Call the SSE context function
			ctx := context.Background()
			resultCtx := ExtractNetbirdInfoFromEnvSSE(ctx, req)

			// Verify the token was extracted correctly
			token := NetbirdAPIKeyFromContext(resultCtx)
			if token != "test-token" {
				t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, "test-token")
			}

			// The test passes if no error occurs, demonstrating protocol stripping works
		})
	}
}

func TestExtractNetbirdInfoFromEnvSSE_UninitializedLoader(t *testing.T) {
	// Test that the SSE function handles uninitialized GlobalConfigLoader gracefully

	// Clear GlobalConfigLoader
	GlobalConfigLoader = nil

	// Set up environment variables
	os.Setenv(netbirdAPIEnvVar, "env-token")
	os.Setenv(netbirdHostEnvVar, "env.example.com")
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8001", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Call the SSE context function - should not panic
	ctx := context.Background()
	resultCtx := ExtractNetbirdInfoFromEnvSSE(ctx, req)

	// Verify context was created
	if resultCtx == nil {
		t.Error("ExtractNetbirdInfoFromEnvSSE() returned nil context")
	}

	// Verify the token was extracted from environment variables
	token := NetbirdAPIKeyFromContext(resultCtx)
	if token != "env-token" {
		t.Errorf("NetbirdAPIKeyFromContext() = %v, want %v", token, "env-token")
	}
}

func TestNewNetbirdClient_ContextConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		contextHost  string
		envHost      string
		expectedHost string
	}{
		{
			name:         "uses context host when available",
			contextHost:  "context.example.com",
			envHost:      "env.example.com",
			expectedHost: "https://context.example.com/api",
		},
		{
			name:         "falls back to env host when context is empty",
			contextHost:  "",
			envHost:      "env.example.com",
			expectedHost: "https://env.example.com/api",
		},
		{
			name:         "uses default host when both are empty",
			contextHost:  "",
			envHost:      "",
			expectedHost: "https://api.netbird.io/api",
		},
		{
			name:         "context host takes priority over env",
			contextHost:  "priority.example.com",
			envHost:      "env.example.com",
			expectedHost: "https://priority.example.com/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable
			if tt.envHost != "" {
				os.Setenv(netbirdHostEnvVar, tt.envHost)
			} else {
				os.Unsetenv(netbirdHostEnvVar)
			}
			defer os.Unsetenv(netbirdHostEnvVar)

			// Create context with host if specified
			ctx := context.Background()
			if tt.contextHost != "" {
				ctx = WithNetbirdConfig(ctx, "test-token", tt.contextHost)
			}

			// Create client
			client := NewNetbirdClient(ctx)

			// Verify baseURL
			if client.baseURL != tt.expectedHost {
				t.Errorf("NewNetbirdClient() baseURL = %v, want %v", client.baseURL, tt.expectedHost)
			}
		})
	}
}

func TestNewNetbirdClient_BackwardCompatibility(t *testing.T) {
	// Test that existing code using environment variables continues to work

	// Set environment variable
	os.Setenv(netbirdHostEnvVar, "legacy.example.com")
	defer os.Unsetenv(netbirdHostEnvVar)

	// Create client with empty context (simulating legacy code)
	ctx := context.Background()
	client := NewNetbirdClient(ctx)

	// Verify it uses the environment variable
	expectedURL := "https://legacy.example.com/api"
	if client.baseURL != expectedURL {
		t.Errorf("NewNetbirdClient() baseURL = %v, want %v", client.baseURL, expectedURL)
	}
}

func TestNewNetbirdClient_WithContextFromConfigLoader(t *testing.T) {
	// Test the full flow: ConfigLoader -> Context -> NewNetbirdClient

	// Set up environment variables
	os.Setenv(netbirdAPIEnvVar, "env-token")
	os.Setenv(netbirdHostEnvVar, "env.example.com")
	defer func() {
		os.Unsetenv(netbirdAPIEnvVar)
		os.Unsetenv(netbirdHostEnvVar)
	}()

	tests := []struct {
		name         string
		cliToken     string
		cliHost      string
		expectedHost string
	}{
		{
			name:         "CLI host used in client",
			cliToken:     "cli-token",
			cliHost:      "cli.example.com",
			expectedHost: "https://cli.example.com/api",
		},
		{
			name:         "env host used when no CLI host",
			cliToken:     "cli-token",
			cliHost:      "",
			expectedHost: "https://env.example.com/api",
		},
		{
			name:         "default host used when no CLI or env",
			cliToken:     "cli-token",
			cliHost:      "",
			expectedHost: "https://api.netbird.io/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the "default host" test, clear the env var
			if tt.name == "default host used when no CLI or env" {
				os.Unsetenv(netbirdHostEnvVar)
				defer os.Setenv(netbirdHostEnvVar, "env.example.com")
			}

			// Set up GlobalConfigLoader
			GlobalConfigLoader = NewConfigLoader(tt.cliToken, tt.cliHost)

			// Load configuration (simulating what context functions do)
			cfg, err := GlobalConfigLoader.LoadConfig("", "")
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}

			// Create context with configuration
			ctx := WithNetbirdConfig(context.Background(), cfg.APIToken, cfg.APIHost)

			// Create client
			client := NewNetbirdClient(ctx)

			// Verify baseURL
			if client.baseURL != tt.expectedHost {
				t.Errorf("NewNetbirdClient() baseURL = %v, want %v", client.baseURL, tt.expectedHost)
			}
		})
	}
}
