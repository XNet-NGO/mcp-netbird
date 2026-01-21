// Copyright 2025-2026 XNet Inc.
// Copyright 2025-2026 Joshua S. Doucette
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Originally derived from MCP Server for Grafana by Grafana Labs.

package mcpnetbird

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultNetbirdHost = "api.netbird.io"
	defaultNetbirdURL  = "https://" + defaultNetbirdHost
	netbirdAPIPath     = "/api"

	netbirdHostEnvVar = "NETBIRD_HOST"
	netbirdAPIEnvVar  = "NETBIRD_API_TOKEN"
)

// Config holds the Netbird API configuration
type Config struct {
	APIToken string
	APIHost  string
}

// ConfigLoader loads configuration from multiple sources with priority order
type ConfigLoader struct {
	cliToken string
	cliHost  string
}

// GlobalConfigLoader is the global configuration loader instance
var GlobalConfigLoader *ConfigLoader

// NewConfigLoader creates a new configuration loader with CLI arguments
func NewConfigLoader(cliToken, cliHost string) *ConfigLoader {
	return &ConfigLoader{
		cliToken: cliToken,
		cliHost:  cliHost,
	}
}

// LoadConfig loads configuration with priority: CLI > HTTP headers > env vars
func (cl *ConfigLoader) LoadConfig(httpToken, httpHost string) (*Config, error) {
	cfg := &Config{}

	// Load API token with priority order
	if cl.cliToken != "" {
		cfg.APIToken = cl.cliToken
	} else if httpToken != "" {
		cfg.APIToken = httpToken
	} else {
		cfg.APIToken = os.Getenv(netbirdAPIEnvVar)
	}

	// Load API host with priority order
	if cl.cliHost != "" {
		cfg.APIHost = stripProtocolPrefix(cl.cliHost)
	} else if httpHost != "" {
		cfg.APIHost = stripProtocolPrefix(httpHost)
	} else {
		envHost := os.Getenv(netbirdHostEnvVar)
		if envHost != "" {
			cfg.APIHost = stripProtocolPrefix(envHost)
		} else {
			cfg.APIHost = defaultNetbirdHost
		}
	}

	return cfg, nil
}

// stripProtocolPrefix removes http:// or https:// prefix from hostname
func stripProtocolPrefix(host string) string {
	// Remove http:// prefix
	if strings.HasPrefix(host, "http://") {
		return strings.TrimPrefix(host, "http://")
	}
	// Remove https:// prefix
	if strings.HasPrefix(host, "https://") {
		return strings.TrimPrefix(host, "https://")
	}
	return host
}

// ValidateConfig validates the configuration and returns descriptive errors
func ValidateConfig(cfg *Config) error {
	// Validate API token
	if cfg.APIToken == "" {
		return fmt.Errorf("API token is required but not provided")
	}
	if strings.TrimSpace(cfg.APIToken) == "" {
		return fmt.Errorf("API token cannot be empty or whitespace-only")
	}

	// Validate API host
	if cfg.APIHost == "" {
		return fmt.Errorf("API host is required but not provided")
	}
	if strings.TrimSpace(cfg.APIHost) == "" {
		return fmt.Errorf("API host cannot be empty or whitespace-only")
	}

	// Check for protocol prefix (should have been stripped already, but validate)
	if strings.HasPrefix(cfg.APIHost, "http://") || strings.HasPrefix(cfg.APIHost, "https://") {
		return fmt.Errorf("API host '%s' should not contain protocol prefix (http:// or https://)", cfg.APIHost)
	}

	// Validate URL format - check for invalid characters and basic structure
	// A valid hostname should not contain spaces, and should have valid characters
	if strings.Contains(cfg.APIHost, " ") {
		return fmt.Errorf("API host '%s' is not a valid URL format: contains spaces", cfg.APIHost)
	}

	// Check for other invalid URL characters
	invalidChars := []string{"<", ">", "\"", "{", "}", "|", "\\", "^", "`"}
	for _, char := range invalidChars {
		if strings.Contains(cfg.APIHost, char) {
			return fmt.Errorf("API host '%s' is not a valid URL format: contains invalid character '%s'", cfg.APIHost, char)
		}
	}

	// Basic check: hostname should have at least one character that's not just special chars
	// Allow alphanumeric, dots, hyphens, colons (for ports), and forward slashes (for paths)
	validHostChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-:/"
	hasValidChar := false
	for _, char := range cfg.APIHost {
		if strings.ContainsRune(validHostChars, char) {
			hasValidChar = true
			break
		}
	}
	if !hasValidChar {
		return fmt.Errorf("API host '%s' is not a valid URL format: no valid hostname characters", cfg.APIHost)
	}

	return nil
}

// NetbirdClient provides methods to interact with the Netbird API
type NetbirdClient struct {
	baseURL string
	client  *http.Client
}

// Single global variable for testing
var TestNetbirdClient *NetbirdClient

// NewNetbirdClient creates a new NetbirdClient with configuration from context.
// If context doesn't contain API host, falls back to environment variable for backward compatibility.
func NewNetbirdClient(ctx context.Context) *NetbirdClient {
	// Try to get host from context first
	host := NetbirdAPIHostFromContext(ctx)
	
	// Fall back to environment variable for backward compatibility
	if host == "" {
		host = os.Getenv(netbirdHostEnvVar)
	}
	
	// Use default if still empty
	if host == "" {
		host = defaultNetbirdHost
	}

	baseURL := "https://" + host + netbirdAPIPath
	return &NetbirdClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func NewNetbirdClientWithBaseURL(baseURL string) *NetbirdClient {
	return &NetbirdClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// do performs an HTTP request to the Netbird API
func (c *NetbirdClient) do(ctx context.Context, method, path string, body, v any) error {
	token := NetbirdAPIKeyFromContext(ctx)
	if token == "" {
		return fmt.Errorf("netbird API token not found in context")
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request to the Netbird API
func (c *NetbirdClient) Get(ctx context.Context, path string, v any) error {
	return c.do(ctx, http.MethodGet, path, nil, v)
}

// Post performs a POST request to the Netbird API
func (c *NetbirdClient) Post(ctx context.Context, path string, body, v any) error {
	return c.do(ctx, http.MethodPost, path, body, v)
}

// Put performs a PUT request to the Netbird API
func (c *NetbirdClient) Put(ctx context.Context, path string, body, v any) error {
	return c.do(ctx, http.MethodPut, path, body, v)
}

// Delete performs a DELETE request to the Netbird API
func (c *NetbirdClient) Delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

type netbirdAPIKeyKey struct{}
type netbirdAPIHostKey struct{}

// ExtractNetbirdInfoFromEnv is a StdioContextFunc that extracts Netbird configuration
// from CLI arguments and environment variables (no HTTP headers in stdio mode).
var ExtractNetbirdInfoFromEnv server.StdioContextFunc = func(ctx context.Context) context.Context {
	// Ensure GlobalConfigLoader is initialized
	if GlobalConfigLoader == nil {
		log.Printf("Warning: GlobalConfigLoader not initialized, using empty CLI arguments")
		GlobalConfigLoader = NewConfigLoader("", "")
	}

	// Load configuration from CLI arguments and environment variables only
	// No HTTP headers in stdio mode (httpToken and httpHost are empty strings)
	cfg, err := GlobalConfigLoader.LoadConfig("", "")
	if err != nil {
		log.Printf("Warning: Failed to load configuration: %v", err)
		return WithNetbirdConfig(ctx, "", "")
	}

	// Validate the configuration
	if err := ValidateConfig(cfg); err != nil {
		log.Printf("Warning: Configuration validation failed: %v", err)
		// Still inject the configuration even if validation fails, to maintain backward compatibility
		// The actual API calls will fail with authentication errors if the token is invalid
	} else {
		log.Printf("Successfully loaded and validated Netbird configuration from CLI arguments and environment variables")
	}

	// Inject validated configuration into context
	return WithNetbirdConfig(ctx, cfg.APIToken, cfg.APIHost)
}

// ExtractNetbirdInfoFromEnvSSE is an SSEContextFunc that extracts Netbird configuration
// from CLI arguments, HTTP headers, and environment variables.
var ExtractNetbirdInfoFromEnvSSE server.SSEContextFunc = func(ctx context.Context, req *http.Request) context.Context {
	// Ensure GlobalConfigLoader is initialized
	if GlobalConfigLoader == nil {
		log.Printf("SSE MODE - Warning: GlobalConfigLoader not initialized, using empty CLI arguments")
		GlobalConfigLoader = NewConfigLoader("", "")
	}

	// Extract HTTP headers
	httpToken := req.Header.Get("X-Netbird-API-Token")
	httpHost := req.Header.Get("X-Netbird-Host")

	// Load configuration from CLI arguments, HTTP headers, and environment variables
	cfg, err := GlobalConfigLoader.LoadConfig(httpToken, httpHost)
	if err != nil {
		log.Printf("SSE MODE - Failed to load configuration: %v", err)
		return WithNetbirdConfig(ctx, "", "")
	}

	// Validate the configuration
	if err := ValidateConfig(cfg); err != nil {
		// Return HTTP 401 for missing API token
		if cfg.APIToken == "" || strings.TrimSpace(cfg.APIToken) == "" {
			log.Printf("SSE MODE - Missing API token from all sources")
			// Note: We can't directly return HTTP 401 from here, but we inject empty token
			// which will cause authentication to fail at the API call level
			return WithNetbirdConfig(ctx, "", "")
		}
		log.Printf("SSE MODE - Configuration validation failed: %v", err)
		// For other validation errors, still inject the configuration
		// The actual API calls will fail with appropriate errors
	} else {
		log.Printf("SSE MODE - Successfully loaded and validated Netbird configuration")
	}

	// Inject validated configuration into context
	return WithNetbirdConfig(ctx, cfg.APIToken, cfg.APIHost)
}

// WithNetbirdConfig adds the Netbird API token and host to the context.
func WithNetbirdConfig(ctx context.Context, apiKey, apiHost string) context.Context {
	ctx = context.WithValue(ctx, netbirdAPIKeyKey{}, apiKey)
	ctx = context.WithValue(ctx, netbirdAPIHostKey{}, apiHost)
	return ctx
}

// WithNetbirdAPIKey adds the Netbird API key to the context.
// Deprecated: Use WithNetbirdConfig instead for full configuration support.
func WithNetbirdAPIKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, netbirdAPIKeyKey{}, apiKey)
}

// NetbirdAPIKeyFromContext extracts the Netbird API key from the context.
func NetbirdAPIKeyFromContext(ctx context.Context) string {
	if v := ctx.Value(netbirdAPIKeyKey{}); v != nil {
		return v.(string)
	}
	return ""
}

// NetbirdAPIHostFromContext extracts the Netbird API host from the context.
func NetbirdAPIHostFromContext(ctx context.Context) string {
	if v := ctx.Value(netbirdAPIHostKey{}); v != nil {
		return v.(string)
	}
	return ""
}

// ComposeStdioContextFuncs composes multiple StdioContextFuncs into a single one.
func ComposeStdioContextFuncs(funcs ...server.StdioContextFunc) server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
		for _, f := range funcs {
			ctx = f(ctx)
		}
		return ctx
	}
}

// ComposeSSEContextFuncs composes multiple SSEContextFuncs into a single one.
func ComposeSSEContextFuncs(funcs ...server.SSEContextFunc) server.SSEContextFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		for _, f := range funcs {
			ctx = f(ctx, req)
		}
		return ctx
	}
}

// ComposedStdioContextFunc is a StdioContextFunc that comprises all predefined StdioContextFuncs.
var ComposedStdioContextFunc = ComposeStdioContextFuncs(
	ExtractNetbirdInfoFromEnv,
)

// ComposedSSEContextFunc is an SSEContextFunc that comprises all predefined SSEContextFuncs.
var ComposedSSEContextFunc = ComposeSSEContextFuncs(
	ExtractNetbirdInfoFromEnvSSE,
)
