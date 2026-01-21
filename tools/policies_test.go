package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
)

func TestListNetbirdPolicies(t *testing.T) {
	// Mock response data
	mockResp := []NetbirdPolicy{
		{
			ID:   "policy1",
			Name: "Test Policy",
			// Add other fields as needed for your struct
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	policies, err := listNetbirdPolicies(ctx, ListNetbirdPoliciesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 1 || policies[0].ID != "policy1" {
		t.Errorf("unexpected result: %+v", policies)
	}
}

func TestPortRangeMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		portRange PortRange
		expected string
	}{
		{
			name: "basic port range",
			portRange: PortRange{
				Start: 80,
				End:   443,
			},
			expected: `{"start":80,"end":443}`,
		},
		{
			name: "single port range",
			portRange: PortRange{
				Start: 8080,
				End:   8080,
			},
			expected: `{"start":8080,"end":8080}`,
		},
		{
			name: "wide port range",
			portRange: PortRange{
				Start: 1024,
				End:   65535,
			},
			expected: `{"start":1024,"end":65535}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.portRange)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed: got %s, want %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded PortRange
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.Start != tt.portRange.Start || decoded.End != tt.portRange.End {
				t.Errorf("unmarshaling failed: got %+v, want %+v", decoded, tt.portRange)
			}
		})
	}
}

func TestPortRangeUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected PortRange
		wantErr  bool
	}{
		{
			name: "valid port range",
			json: `{"start":80,"end":443}`,
			expected: PortRange{
				Start: 80,
				End:   443,
			},
			wantErr: false,
		},
		{
			name: "zero values",
			json: `{"start":0,"end":0}`,
			expected: PortRange{
				Start: 0,
				End:   0,
			},
			wantErr: false,
		},
		{
			name: "missing fields defaults to zero",
			json: `{}`,
			expected: PortRange{
				Start: 0,
				End:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var portRange PortRange
			err := json.Unmarshal([]byte(tt.json), &portRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if portRange.Start != tt.expected.Start || portRange.End != tt.expected.End {
					t.Errorf("got %+v, want %+v", portRange, tt.expected)
				}
			}
		})
	}
}

func TestResourceReferenceMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		resource ResourceReference
		expected string
	}{
		{
			name: "basic resource reference",
			resource: ResourceReference{
				ID:   "res-123",
				Type: "network",
			},
			expected: `{"id":"res-123","type":"network"}`,
		},
		{
			name: "resource reference with different type",
			resource: ResourceReference{
				ID:   "peer-456",
				Type: "peer",
			},
			expected: `{"id":"peer-456","type":"peer"}`,
		},
		{
			name: "empty resource reference",
			resource: ResourceReference{
				ID:   "",
				Type: "",
			},
			expected: `{"id":"","type":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.resource)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed: got %s, want %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded ResourceReference
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.ID != tt.resource.ID || decoded.Type != tt.resource.Type {
				t.Errorf("unmarshaling failed: got %+v, want %+v", decoded, tt.resource)
			}
		})
	}
}

func TestResourceReferenceUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected ResourceReference
		wantErr  bool
	}{
		{
			name: "valid resource reference",
			json: `{"id":"res-123","type":"network"}`,
			expected: ResourceReference{
				ID:   "res-123",
				Type: "network",
			},
			wantErr: false,
		},
		{
			name: "empty values",
			json: `{"id":"","type":""}`,
			expected: ResourceReference{
				ID:   "",
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "missing fields defaults to empty strings",
			json: `{}`,
			expected: ResourceReference{
				ID:   "",
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "partial fields",
			json: `{"id":"res-789"}`,
			expected: ResourceReference{
				ID:   "res-789",
				Type: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resource ResourceReference
			err := json.Unmarshal([]byte(tt.json), &resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resource.ID != tt.expected.ID || resource.Type != tt.expected.Type {
					t.Errorf("got %+v, want %+v", resource, tt.expected)
				}
			}
		})
	}
}

func TestNetbirdPolicyRuleNewFieldsMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		rule     NetbirdPolicyRule
		expected string
	}{
		{
			name: "rule with port ranges",
			rule: NetbirdPolicyRule{
				Action:        "accept",
				Bidirectional: true,
				Protocol:      "tcp",
				PortRanges: &[]PortRange{
					{Start: 80, End: 80},
					{Start: 443, End: 443},
				},
			},
			expected: `{"action":"accept","bidirectional":true,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"tcp","sources":null,"port_ranges":[{"start":80,"end":80},{"start":443,"end":443}]}`,
		},
		{
			name: "rule with authorized groups",
			rule: NetbirdPolicyRule{
				Action:        "accept",
				Bidirectional: false,
				Protocol:      "udp",
				AuthorizedGroups: &map[string][]string{
					"group1": {"user1", "user2"},
					"group2": {"user3"},
				},
			},
			expected: `{"action":"accept","bidirectional":false,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"udp","sources":null,"authorized_groups":{"group1":["user1","user2"],"group2":["user3"]}}`,
		},
		{
			name: "rule with source resource",
			rule: NetbirdPolicyRule{
				Action:        "accept",
				Bidirectional: true,
				Protocol:      "tcp",
				SourceResource: &ResourceReference{
					ID:   "res-123",
					Type: "network",
				},
			},
			expected: `{"action":"accept","bidirectional":true,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"tcp","sources":null,"sourceResource":{"id":"res-123","type":"network"}}`,
		},
		{
			name: "rule with destination resource",
			rule: NetbirdPolicyRule{
				Action:        "accept",
				Bidirectional: false,
				Protocol:      "tcp",
				DestinationResource: &ResourceReference{
					ID:   "res-456",
					Type: "peer",
				},
			},
			expected: `{"action":"accept","bidirectional":false,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"tcp","sources":null,"destinationResource":{"id":"res-456","type":"peer"}}`,
		},
		{
			name: "rule with all new fields",
			rule: NetbirdPolicyRule{
				Action:        "accept",
				Bidirectional: true,
				Protocol:      "tcp",
				PortRanges: &[]PortRange{
					{Start: 8080, End: 8090},
				},
				AuthorizedGroups: &map[string][]string{
					"admin": {"admin1"},
				},
				SourceResource: &ResourceReference{
					ID:   "src-789",
					Type: "network",
				},
				DestinationResource: &ResourceReference{
					ID:   "dst-012",
					Type: "peer",
				},
			},
			expected: `{"action":"accept","bidirectional":true,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"tcp","sources":null,"port_ranges":[{"start":8080,"end":8090}],"authorized_groups":{"admin":["admin1"]},"sourceResource":{"id":"src-789","type":"network"},"destinationResource":{"id":"dst-012","type":"peer"}}`,
		},
		{
			name: "rule without new fields (omitempty test)",
			rule: NetbirdPolicyRule{
				Action:        "deny",
				Bidirectional: false,
				Protocol:      "icmp",
			},
			expected: `{"action":"deny","bidirectional":false,"description":"","destinations":null,"enabled":false,"id":"","name":"","protocol":"icmp","sources":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.rule)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed:\ngot:  %s\nwant: %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded NetbirdPolicyRule
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			// Verify basic fields
			if decoded.Action != tt.rule.Action {
				t.Errorf("Action mismatch: got %s, want %s", decoded.Action, tt.rule.Action)
			}
			if decoded.Bidirectional != tt.rule.Bidirectional {
				t.Errorf("Bidirectional mismatch: got %v, want %v", decoded.Bidirectional, tt.rule.Bidirectional)
			}
			if decoded.Protocol != tt.rule.Protocol {
				t.Errorf("Protocol mismatch: got %s, want %s", decoded.Protocol, tt.rule.Protocol)
			}

			// Verify PortRanges
			if tt.rule.PortRanges != nil {
				if decoded.PortRanges == nil {
					t.Error("PortRanges should not be nil")
				} else if len(*decoded.PortRanges) != len(*tt.rule.PortRanges) {
					t.Errorf("PortRanges length mismatch: got %d, want %d", len(*decoded.PortRanges), len(*tt.rule.PortRanges))
				} else {
					for i, pr := range *tt.rule.PortRanges {
						decodedPR := (*decoded.PortRanges)[i]
						if decodedPR.Start != pr.Start || decodedPR.End != pr.End {
							t.Errorf("PortRange[%d] mismatch: got %+v, want %+v", i, decodedPR, pr)
						}
					}
				}
			} else if decoded.PortRanges != nil {
				t.Error("PortRanges should be nil")
			}

			// Verify AuthorizedGroups
			if tt.rule.AuthorizedGroups != nil {
				if decoded.AuthorizedGroups == nil {
					t.Error("AuthorizedGroups should not be nil")
				} else {
					for key, values := range *tt.rule.AuthorizedGroups {
						decodedValues, ok := (*decoded.AuthorizedGroups)[key]
						if !ok {
							t.Errorf("AuthorizedGroups missing key: %s", key)
							continue
						}
						if len(decodedValues) != len(values) {
							t.Errorf("AuthorizedGroups[%s] length mismatch: got %d, want %d", key, len(decodedValues), len(values))
							continue
						}
						for i, v := range values {
							if decodedValues[i] != v {
								t.Errorf("AuthorizedGroups[%s][%d] mismatch: got %s, want %s", key, i, decodedValues[i], v)
							}
						}
					}
				}
			} else if decoded.AuthorizedGroups != nil {
				t.Error("AuthorizedGroups should be nil")
			}

			// Verify SourceResource
			if tt.rule.SourceResource != nil {
				if decoded.SourceResource == nil {
					t.Error("SourceResource should not be nil")
				} else if decoded.SourceResource.ID != tt.rule.SourceResource.ID || decoded.SourceResource.Type != tt.rule.SourceResource.Type {
					t.Errorf("SourceResource mismatch: got %+v, want %+v", decoded.SourceResource, tt.rule.SourceResource)
				}
			} else if decoded.SourceResource != nil {
				t.Error("SourceResource should be nil")
			}

			// Verify DestinationResource
			if tt.rule.DestinationResource != nil {
				if decoded.DestinationResource == nil {
					t.Error("DestinationResource should not be nil")
				} else if decoded.DestinationResource.ID != tt.rule.DestinationResource.ID || decoded.DestinationResource.Type != tt.rule.DestinationResource.Type {
					t.Errorf("DestinationResource mismatch: got %+v, want %+v", decoded.DestinationResource, tt.rule.DestinationResource)
				}
			} else if decoded.DestinationResource != nil {
				t.Error("DestinationResource should be nil")
			}
		})
	}
}

func TestNetbirdPolicyRuleOmitEmptyBehavior(t *testing.T) {
	// Test that nil pointer fields are omitted from JSON
	rule := NetbirdPolicyRule{
		Action:        "accept",
		Bidirectional: true,
		Protocol:      "tcp",
		// All new fields are nil
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify that the new fields are not present in the JSON
	jsonStr := string(data)
	if contains(jsonStr, "port_ranges") {
		t.Error("port_ranges should be omitted when nil")
	}
	if contains(jsonStr, "authorized_groups") {
		t.Error("authorized_groups should be omitted when nil")
	}
	if contains(jsonStr, "sourceResource") {
		t.Error("sourceResource should be omitted when nil")
	}
	if contains(jsonStr, "destinationResource") {
		t.Error("destinationResource should be omitted when nil")
	}
}

func TestNetbirdPolicyRuleEmptySliceAndMapBehavior(t *testing.T) {
	// Test that empty slices and maps are included in JSON (not nil)
	emptyPortRanges := []PortRange{}
	emptyAuthorizedGroups := map[string][]string{}

	rule := NetbirdPolicyRule{
		Action:           "accept",
		Bidirectional:    true,
		Protocol:         "tcp",
		PortRanges:       &emptyPortRanges,
		AuthorizedGroups: &emptyAuthorizedGroups,
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "port_ranges") {
		t.Error("port_ranges should be included when empty slice")
	}
	if !contains(jsonStr, "authorized_groups") {
		t.Error("authorized_groups should be included when empty map")
	}

	// Verify unmarshaling preserves empty collections
	var decoded NetbirdPolicyRule
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.PortRanges == nil {
		t.Error("PortRanges should not be nil after unmarshaling empty slice")
	} else if len(*decoded.PortRanges) != 0 {
		t.Errorf("PortRanges should be empty: got %d elements", len(*decoded.PortRanges))
	}

	if decoded.AuthorizedGroups == nil {
		t.Error("AuthorizedGroups should not be nil after unmarshaling empty map")
	} else if len(*decoded.AuthorizedGroups) != 0 {
		t.Errorf("AuthorizedGroups should be empty: got %d elements", len(*decoded.AuthorizedGroups))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Property Test 1: Rule Formatting Preserves Structure
// Validates: Requirements 1.1, 1.4, 2.1
// Feature: mcp-netbird-improvements, Property 1: Rule Formatting Preserves Structure
func TestFormatRuleForAPI_PropertyPreservesStructure(t *testing.T) {
	// Run 100+ iterations with random valid rules
	for i := 0; i < 100; i++ {
		// Generate random rule with various field combinations
		rule := generateRandomValidRule(i)
		
		// Format the rule
		formatted, err := FormatRuleForAPI(rule)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		
		// Verify all fields preserved except sources/destinations converted to strings
		for key, value := range rule {
			if key == "sources" || key == "destinations" {
				// These should be converted to string arrays
				formattedValue, ok := formatted[key]
				if !ok {
					t.Errorf("iteration %d: field %s missing in formatted rule", i, key)
					continue
				}
				
				// Verify it's a string array
				stringArray, ok := formattedValue.([]string)
				if !ok {
					t.Errorf("iteration %d: field %s not converted to string array, got type %T", i, key, formattedValue)
					continue
				}
				
				// Verify the IDs match
				originalIDs := extractIDsFromValue(value)
				if len(stringArray) != len(originalIDs) {
					t.Errorf("iteration %d: field %s length mismatch: got %d, want %d", i, key, len(stringArray), len(originalIDs))
				}
				for j, id := range originalIDs {
					if j < len(stringArray) && stringArray[j] != id {
						t.Errorf("iteration %d: field %s[%d] mismatch: got %s, want %s", i, key, j, stringArray[j], id)
					}
				}
			} else {
				// All other fields should be preserved unchanged
				formattedValue, ok := formatted[key]
				if !ok {
					t.Errorf("iteration %d: field %s missing in formatted rule", i, key)
					continue
				}
				
				// Deep comparison for complex types
				if !deepEqual(value, formattedValue) {
					t.Errorf("iteration %d: field %s not preserved: got %v, want %v", i, key, formattedValue, value)
				}
			}
		}
	}
}

// Helper function to generate random valid rules for property testing
func generateRandomValidRule(seed int) map[string]interface{} {
	rule := make(map[string]interface{})
	
	// Required fields
	rule["name"] = fmt.Sprintf("rule-%d", seed)
	rule["enabled"] = seed%2 == 0
	rule["action"] = []string{"accept", "drop"}[seed%2]
	rule["bidirectional"] = seed%3 == 0
	rule["protocol"] = []string{"tcp", "udp", "icmp", "all"}[seed%4]
	
	// Optional description
	if seed%3 == 0 {
		rule["description"] = fmt.Sprintf("description-%d", seed)
	}
	
	// Sources - vary between string arrays and object arrays
	sourcesCount := (seed % 3) + 1
	if seed%2 == 0 {
		// String array format
		sources := make([]string, sourcesCount)
		for i := 0; i < sourcesCount; i++ {
			sources[i] = fmt.Sprintf("group-src-%d-%d", seed, i)
		}
		rule["sources"] = sources
	} else {
		// Object array format (as returned by API)
		sources := make([]interface{}, sourcesCount)
		for i := 0; i < sourcesCount; i++ {
			sources[i] = map[string]interface{}{
				"id":   fmt.Sprintf("group-src-%d-%d", seed, i),
				"name": fmt.Sprintf("Source Group %d-%d", seed, i),
			}
		}
		rule["sources"] = sources
	}
	
	// Destinations - vary between string arrays and object arrays
	destCount := (seed % 3) + 1
	if seed%3 == 0 {
		// String array format
		destinations := make([]string, destCount)
		for i := 0; i < destCount; i++ {
			destinations[i] = fmt.Sprintf("group-dst-%d-%d", seed, i)
		}
		rule["destinations"] = destinations
	} else {
		// Object array format (as returned by API)
		destinations := make([]interface{}, destCount)
		for i := 0; i < destCount; i++ {
			destinations[i] = map[string]interface{}{
				"id":   fmt.Sprintf("group-dst-%d-%d", seed, i),
				"name": fmt.Sprintf("Dest Group %d-%d", seed, i),
			}
		}
		rule["destinations"] = destinations
	}
	
	// Optional port ranges
	if seed%4 == 0 {
		portRanges := make([]interface{}, (seed%2)+1)
		for i := 0; i < len(portRanges); i++ {
			portRanges[i] = map[string]interface{}{
				"start": 80 + i*100,
				"end":   80 + i*100 + 50,
			}
		}
		rule["port_ranges"] = portRanges
	}
	
	// Optional authorized groups
	if seed%5 == 0 {
		rule["authorized_groups"] = map[string]interface{}{
			fmt.Sprintf("group-%d", seed): []string{"user1", "user2"},
		}
	}
	
	return rule
}

// Helper function to extract IDs from various value formats
func extractIDsFromValue(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		ids := make([]string, 0, len(v))
		for _, item := range v {
			switch itemVal := item.(type) {
			case string:
				ids = append(ids, itemVal)
			case map[string]interface{}:
				if id, ok := itemVal["id"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
		return ids
	default:
		return []string{}
	}
}

// Helper function for deep equality comparison
func deepEqual(a, b interface{}) bool {
	// Simple implementation for common types
	switch aVal := a.(type) {
	case string:
		bVal, ok := b.(string)
		return ok && aVal == bVal
	case bool:
		bVal, ok := b.(bool)
		return ok && aVal == bVal
	case int:
		bVal, ok := b.(int)
		return ok && aVal == bVal
	case float64:
		bVal, ok := b.(float64)
		return ok && aVal == bVal
	case []string:
		bVal, ok := b.([]string)
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i, v := range aVal {
			if bVal[i] != v {
				return false
			}
		}
		return true
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for key, aValue := range aVal {
			bValue, ok := bVal[key]
			if !ok || !deepEqual(aValue, bValue) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i, aValue := range aVal {
			if !deepEqual(aValue, bVal[i]) {
				return false
			}
		}
		return true
	default:
		// For other types, use simple equality (but avoid comparing slices directly)
		// This is a fallback for types we haven't explicitly handled
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
}

// Unit tests for FormatRuleForAPI
func TestFormatRuleForAPI_StringArraySources(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources":       []string{"group-1", "group-2"},
		"destinations":  []string{"group-3"},
	}
	
	formatted, err := FormatRuleForAPI(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	sources, ok := formatted["sources"].([]string)
	if !ok {
		t.Fatalf("sources not a string array: %T", formatted["sources"])
	}
	
	if len(sources) != 2 || sources[0] != "group-1" || sources[1] != "group-2" {
		t.Errorf("sources not preserved: %v", sources)
	}
}

func TestFormatRuleForAPI_ObjectArraySources(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources": []interface{}{
			map[string]interface{}{"id": "group-1", "name": "Group 1"},
			map[string]interface{}{"id": "group-2", "name": "Group 2"},
		},
		"destinations": []interface{}{
			map[string]interface{}{"id": "group-3", "name": "Group 3"},
		},
	}
	
	formatted, err := FormatRuleForAPI(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	sources, ok := formatted["sources"].([]string)
	if !ok {
		t.Fatalf("sources not a string array: %T", formatted["sources"])
	}
	
	if len(sources) != 2 || sources[0] != "group-1" || sources[1] != "group-2" {
		t.Errorf("sources not converted correctly: %v", sources)
	}
	
	destinations, ok := formatted["destinations"].([]string)
	if !ok {
		t.Fatalf("destinations not a string array: %T", formatted["destinations"])
	}
	
	if len(destinations) != 1 || destinations[0] != "group-3" {
		t.Errorf("destinations not converted correctly: %v", destinations)
	}
}

func TestFormatRuleForAPI_MixedFormats(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources": []interface{}{
			"group-1",
			map[string]interface{}{"id": "group-2", "name": "Group 2"},
		},
		"destinations": []string{"group-3"},
	}
	
	formatted, err := FormatRuleForAPI(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	sources, ok := formatted["sources"].([]string)
	if !ok {
		t.Fatalf("sources not a string array: %T", formatted["sources"])
	}
	
	if len(sources) != 2 || sources[0] != "group-1" || sources[1] != "group-2" {
		t.Errorf("sources not converted correctly: %v", sources)
	}
}

func TestFormatRuleForAPI_PreservesOtherFields(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"description":   "test description",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources":       []string{"group-1"},
		"destinations":  []string{"group-2"},
		"port_ranges": []interface{}{
			map[string]interface{}{"start": 80, "end": 80},
			map[string]interface{}{"start": 443, "end": 443},
		},
		"authorized_groups": map[string]interface{}{
			"admin-group": []string{"user1", "user2"},
		},
	}
	
	formatted, err := FormatRuleForAPI(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify all fields are preserved
	if formatted["name"] != "test-rule" {
		t.Errorf("name not preserved")
	}
	if formatted["description"] != "test description" {
		t.Errorf("description not preserved")
	}
	if formatted["enabled"] != true {
		t.Errorf("enabled not preserved")
	}
	if formatted["action"] != "accept" {
		t.Errorf("action not preserved")
	}
	if formatted["bidirectional"] != false {
		t.Errorf("bidirectional not preserved")
	}
	if formatted["protocol"] != "tcp" {
		t.Errorf("protocol not preserved")
	}
	
	// Verify port_ranges preserved
	portRanges, ok := formatted["port_ranges"].([]interface{})
	if !ok || len(portRanges) != 2 {
		t.Errorf("port_ranges not preserved correctly")
	}
	
	// Verify authorized_groups preserved
	_, ok = formatted["authorized_groups"].(map[string]interface{})
	if !ok {
		t.Errorf("authorized_groups not preserved correctly")
	}
}

func TestFormatRuleForAPI_NilRule(t *testing.T) {
	_, err := FormatRuleForAPI(nil)
	if err == nil {
		t.Error("expected error for nil rule")
	}
}

func TestFormatRuleForAPI_EmptyArrays(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources":       []string{},
		"destinations":  []interface{}{},
	}
	
	formatted, err := FormatRuleForAPI(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	sources, ok := formatted["sources"].([]string)
	if !ok {
		t.Fatalf("sources not a string array: %T", formatted["sources"])
	}
	if len(sources) != 0 {
		t.Errorf("sources should be empty: %v", sources)
	}
	
	destinations, ok := formatted["destinations"].([]string)
	if !ok {
		t.Fatalf("destinations not a string array: %T", formatted["destinations"])
	}
	if len(destinations) != 0 {
		t.Errorf("destinations should be empty: %v", destinations)
	}
}

func TestFormatRuleForAPI_MissingIDField(t *testing.T) {
	rule := map[string]interface{}{
		"name":          "test-rule",
		"enabled":       true,
		"action":        "accept",
		"bidirectional": false,
		"protocol":      "tcp",
		"sources": []interface{}{
			map[string]interface{}{"name": "Group 1"}, // Missing "id" field
		},
		"destinations": []string{"group-2"},
	}
	
	_, err := FormatRuleForAPI(rule)
	if err == nil {
		t.Error("expected error for object missing 'id' field")
	}
	if err != nil && !contains(err.Error(), "missing 'id' field") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// Property Test 2: Invalid Rules Rejected Before API Call
// Validates: Requirements 1.2, 2.3, 3.1, 3.2, 3.4, 3.5
// Feature: mcp-netbird-improvements, Property 2: Invalid Rules Rejected Before API Call
func TestValidatePolicyRules_PropertyRejectsInvalidRules(t *testing.T) {
	// Run 100+ iterations with random invalid rules
	for i := 0; i < 100; i++ {
		// Generate random invalid rule
		invalidRule := generateRandomInvalidRule(i)
		rules := []map[string]interface{}{invalidRule}
		
		// Validate the rule
		err := ValidatePolicyRules(rules)
		
		// Should return an error
		if err == nil {
			t.Errorf("iteration %d: expected validation error for invalid rule: %+v", i, invalidRule)
			continue
		}
		
		// Error message should contain rule identifier and field information
		errMsg := err.Error()
		if !contains(errMsg, "rule") {
			t.Errorf("iteration %d: error message should contain 'rule': %s", i, errMsg)
		}
	}
}

// Property Test 3: Port Range Invariant
// Validates: Requirements 3.3
// Feature: mcp-netbird-improvements, Property 3: Port Range Invariant
func TestValidatePolicyRules_PropertyPortRangeInvariant(t *testing.T) {
	// Run 100+ iterations with random port ranges
	for i := 0; i < 100; i++ {
		start := i*10 + 1
		end := start - (i % 5) - 1 // Ensure start > end for invalid cases
		
		if end < start {
			// Invalid port range (start > end)
			rule := map[string]interface{}{
				"name":          fmt.Sprintf("rule-%d", i),
				"enabled":       true,
				"action":        "accept",
				"bidirectional": false,
				"protocol":      "tcp",
				"sources":       []string{"group-1"},
				"destinations":  []string{"group-2"},
				"port_ranges": []interface{}{
					map[string]interface{}{
						"start": start,
						"end":   end,
					},
				},
			}
			
			rules := []map[string]interface{}{rule}
			err := ValidatePolicyRules(rules)
			
			// Should return an error
			if err == nil {
				t.Errorf("iteration %d: expected validation error for invalid port range (start=%d > end=%d)", i, start, end)
			} else if !contains(err.Error(), "start") || !contains(err.Error(), "end") {
				t.Errorf("iteration %d: error message should mention start and end: %s", i, err.Error())
			}
		}
	}
	
	// Also test valid port ranges (start <= end)
	for i := 0; i < 50; i++ {
		start := i * 10
		end := start + (i % 10)
		
		rule := map[string]interface{}{
			"name":          fmt.Sprintf("rule-%d", i),
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
			"port_ranges": []interface{}{
				map[string]interface{}{
					"start": start,
					"end":   end,
				},
			},
		}
		
		rules := []map[string]interface{}{rule}
		err := ValidatePolicyRules(rules)
		
		// Should not return an error
		if err != nil {
			t.Errorf("iteration %d: unexpected validation error for valid port range (start=%d <= end=%d): %v", i, start, end, err)
		}
	}
}

// Helper function to generate random invalid rules for property testing
func generateRandomInvalidRule(seed int) map[string]interface{} {
	rule := make(map[string]interface{})
	
	// Vary which field is missing or invalid based on seed
	switch seed % 10 {
	case 0:
		// Missing "name" field
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 1:
		// Missing "enabled" field
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 2:
		// Missing "action" field
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 3:
		// Invalid action value
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "invalid-action"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 4:
		// Invalid protocol value
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "invalid-protocol"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 5:
		// Missing sources and sourceResource
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["destinations"] = []string{"group-2"}
	case 6:
		// Missing destinations and destinationResource
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
	case 7:
		// Invalid port range (start > end)
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
		rule["port_ranges"] = []interface{}{
			map[string]interface{}{
				"start": 443,
				"end":   80,
			},
		}
	case 8:
		// Missing "bidirectional" field
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["protocol"] = "tcp"
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	case 9:
		// Missing "protocol" field
		rule["name"] = fmt.Sprintf("rule-%d", seed)
		rule["enabled"] = true
		rule["action"] = "accept"
		rule["bidirectional"] = false
		rule["sources"] = []string{"group-1"}
		rule["destinations"] = []string{"group-2"}
	}
	
	return rule
}

// Unit tests for ValidatePolicyRules
func TestValidatePolicyRules_ValidRules(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error for valid rules: %v", err)
	}
}

func TestValidatePolicyRules_MissingName(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for missing name field")
	}
	if err != nil && !contains(err.Error(), "name") {
		t.Errorf("error should mention 'name': %v", err)
	}
}

func TestValidatePolicyRules_MissingEnabled(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for missing enabled field")
	}
	if err != nil && !contains(err.Error(), "enabled") {
		t.Errorf("error should mention 'enabled': %v", err)
	}
}

func TestValidatePolicyRules_InvalidAction(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "invalid",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for invalid action")
	}
	if err != nil && !contains(err.Error(), "action") {
		t.Errorf("error should mention 'action': %v", err)
	}
}

func TestValidatePolicyRules_InvalidProtocol(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "invalid",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for invalid protocol")
	}
	if err != nil && !contains(err.Error(), "protocol") {
		t.Errorf("error should mention 'protocol': %v", err)
	}
}

func TestValidatePolicyRules_InvalidPortRange(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
			"port_ranges": []interface{}{
				map[string]interface{}{
					"start": 443,
					"end":   80,
				},
			},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for invalid port range (start > end)")
	}
	if err != nil && (!contains(err.Error(), "start") || !contains(err.Error(), "end")) {
		t.Errorf("error should mention 'start' and 'end': %v", err)
	}
}

func TestValidatePolicyRules_ValidPortRange(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
			"port_ranges": []interface{}{
				map[string]interface{}{
					"start": 80,
					"end":   443,
				},
			},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error for valid port range: %v", err)
	}
}

func TestValidatePolicyRules_MissingSource(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"destinations":  []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for missing source")
	}
	if err != nil && !contains(err.Error(), "source") {
		t.Errorf("error should mention 'source': %v", err)
	}
}

func TestValidatePolicyRules_MissingDestination(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for missing destination")
	}
	if err != nil && !contains(err.Error(), "destination") {
		t.Errorf("error should mention 'destination': %v", err)
	}
}

func TestValidatePolicyRules_WithSourceResource(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sourceResource": map[string]interface{}{
				"id":   "res-123",
				"type": "network",
			},
			"destinations": []string{"group-2"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error with sourceResource: %v", err)
	}
}

func TestValidatePolicyRules_WithDestinationResource(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "test-rule",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinationResource": map[string]interface{}{
				"id":   "res-456",
				"type": "peer",
			},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error with destinationResource: %v", err)
	}
}

func TestValidatePolicyRules_NilRules(t *testing.T) {
	err := ValidatePolicyRules(nil)
	if err != nil {
		t.Errorf("unexpected error for nil rules: %v", err)
	}
}

func TestValidatePolicyRules_EmptyRules(t *testing.T) {
	rules := []map[string]interface{}{}
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error for empty rules: %v", err)
	}
}

func TestValidatePolicyRules_MultipleRules(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "rule-1",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
		{
			"name":          "rule-2",
			"enabled":       false,
			"action":        "drop",
			"bidirectional": true,
			"protocol":      "udp",
			"sources":       []string{"group-3"},
			"destinations":  []string{"group-4"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("unexpected error for multiple valid rules: %v", err)
	}
}

func TestValidatePolicyRules_SecondRuleInvalid(t *testing.T) {
	rules := []map[string]interface{}{
		{
			"name":          "rule-1",
			"enabled":       true,
			"action":        "accept",
			"bidirectional": false,
			"protocol":      "tcp",
			"sources":       []string{"group-1"},
			"destinations":  []string{"group-2"},
		},
		{
			"name":          "rule-2",
			"enabled":       false,
			"action":        "invalid",
			"bidirectional": true,
			"protocol":      "udp",
			"sources":       []string{"group-3"},
			"destinations":  []string{"group-4"},
		},
	}
	
	err := ValidatePolicyRules(rules)
	if err == nil {
		t.Error("expected error for invalid second rule")
	}
	if err != nil && !contains(err.Error(), "rule-2") {
		t.Errorf("error should mention 'rule-2': %v", err)
	}
}

// Unit tests for createNetbirdPolicy with validation and formatting
func TestCreateNetbirdPolicy_WithSimpleRules(t *testing.T) {
	// Mock response data
	mockResp := NetbirdPolicy{
		ID:      "policy-123",
		Name:    "Test Policy",
		Enabled: true,
		Rules: []NetbirdPolicyRule{
			{
				ID:            "rule-1",
				Name:          "test-rule",
				Enabled:       true,
				Action:        "accept",
				Bidirectional: false,
				Protocol:      "tcp",
				Sources: []NetbirdPeerGroup{
					{ID: "group-1", Name: "Group 1"},
				},
				Destinations: []NetbirdPeerGroup{
					{ID: "group-2", Name: "Group 2"},
				},
			},
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies" || r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		
		// Verify request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Verify rules are formatted correctly (sources/destinations as string arrays)
		if rules, ok := requestBody["rules"].([]interface{}); ok && len(rules) > 0 {
			rule := rules[0].(map[string]interface{})
			
			// Check sources is a string array
			if sources, ok := rule["sources"].([]interface{}); ok {
				for _, src := range sources {
					if _, ok := src.(string); !ok {
						http.Error(w, "sources must be string array", http.StatusBadRequest)
						return
					}
				}
			}
			
			// Check destinations is a string array
			if destinations, ok := rule["destinations"].([]interface{}); ok {
				for _, dst := range destinations {
					if _, ok := dst.(string); !ok {
						http.Error(w, "destinations must be string array", http.StatusBadRequest)
						return
					}
				}
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	enabled := true
	rules := []NetbirdPolicyRule{
		{
			Name:          "test-rule",
			Enabled:       true,
			Action:        "accept",
			Bidirectional: false,
			Protocol:      "tcp",
			Sources: []NetbirdPeerGroup{
				{ID: "group-1", Name: "Group 1"},
			},
			Destinations: []NetbirdPeerGroup{
				{ID: "group-2", Name: "Group 2"},
			},
		},
	}
	
	policy, err := createNetbirdPolicy(ctx, CreateNetbirdPolicyParams{
		Name:    "Test Policy",
		Enabled: &enabled,
		Rules:   &rules,
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.ID != "policy-123" {
		t.Errorf("unexpected policy ID: %s", policy.ID)
	}
}

func TestCreateNetbirdPolicy_WithComplexRules(t *testing.T) {
	// Mock response data
	mockResp := NetbirdPolicy{
		ID:      "policy-456",
		Name:    "Complex Policy",
		Enabled: true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies" || r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	enabled := true
	portRanges := []PortRange{
		{Start: 80, End: 80},
		{Start: 443, End: 443},
	}
	authorizedGroups := map[string][]string{
		"admin-group": {"user1", "user2"},
	}
	rules := []NetbirdPolicyRule{
		{
			Name:          "complex-rule",
			Enabled:       true,
			Action:        "accept",
			Bidirectional: false,
			Protocol:      "tcp",
			Sources: []NetbirdPeerGroup{
				{ID: "group-1", Name: "Group 1"},
			},
			Destinations: []NetbirdPeerGroup{
				{ID: "group-2", Name: "Group 2"},
			},
			PortRanges:       &portRanges,
			AuthorizedGroups: &authorizedGroups,
		},
	}
	
	policy, err := createNetbirdPolicy(ctx, CreateNetbirdPolicyParams{
		Name:    "Complex Policy",
		Enabled: &enabled,
		Rules:   &rules,
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.ID != "policy-456" {
		t.Errorf("unexpected policy ID: %s", policy.ID)
	}
}

func TestCreateNetbirdPolicy_ValidationError(t *testing.T) {
	// Create mock HTTP server (should not be called)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("API should not be called when validation fails")
		http.Error(w, "should not reach here", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	enabled := true
	rules := []NetbirdPolicyRule{
		{
			Name:          "invalid-rule",
			Enabled:       true,
			Action:        "invalid-action", // Invalid action
			Bidirectional: false,
			Protocol:      "tcp",
			Sources: []NetbirdPeerGroup{
				{ID: "group-1", Name: "Group 1"},
			},
			Destinations: []NetbirdPeerGroup{
				{ID: "group-2", Name: "Group 2"},
			},
		},
	}
	
	_, err := createNetbirdPolicy(ctx, CreateNetbirdPolicyParams{
		Name:    "Test Policy",
		Enabled: &enabled,
		Rules:   &rules,
	})
	
	if err == nil {
		t.Error("expected validation error")
	}
	if err != nil && !contains(err.Error(), "validation error") {
		t.Errorf("error should mention validation: %v", err)
	}
}

func TestUpdateNetbirdPolicy_WithRules(t *testing.T) {
	// Mock response data
	mockResp := NetbirdPolicy{
		ID:      "policy-789",
		Name:    "Updated Policy",
		Enabled: true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies/policy-789" || r.Method != "PUT" {
			http.NotFound(w, r)
			return
		}
		
		// Verify request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Verify rules are formatted correctly
		if rules, ok := requestBody["rules"].([]interface{}); ok && len(rules) > 0 {
			rule := rules[0].(map[string]interface{})
			
			// Check sources is a string array
			if sources, ok := rule["sources"].([]interface{}); ok {
				for _, src := range sources {
					if _, ok := src.(string); !ok {
						http.Error(w, "sources must be string array", http.StatusBadRequest)
						return
					}
				}
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	name := "Updated Policy"
	rules := []NetbirdPolicyRule{
		{
			Name:          "updated-rule",
			Enabled:       true,
			Action:        "accept",
			Bidirectional: false,
			Protocol:      "tcp",
			Sources: []NetbirdPeerGroup{
				{ID: "group-3", Name: "Group 3"},
			},
			Destinations: []NetbirdPeerGroup{
				{ID: "group-4", Name: "Group 4"},
			},
		},
	}
	
	policy, err := updateNetbirdPolicy(ctx, UpdateNetbirdPolicyParams{
		PolicyID: "policy-789",
		Name:     &name,
		Rules:    &rules,
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.ID != "policy-789" {
		t.Errorf("unexpected policy ID: %s", policy.ID)
	}
}

func TestUpdateNetbirdPolicy_ValidationError(t *testing.T) {
	// Create mock HTTP server (should not be called)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("API should not be called when validation fails")
		http.Error(w, "should not reach here", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	rules := []NetbirdPolicyRule{
		{
			Name:          "invalid-rule",
			Enabled:       true,
			Action:        "accept",
			Bidirectional: false,
			Protocol:      "invalid-protocol", // Invalid protocol
			Sources: []NetbirdPeerGroup{
				{ID: "group-1", Name: "Group 1"},
			},
			Destinations: []NetbirdPeerGroup{
				{ID: "group-2", Name: "Group 2"},
			},
		},
	}
	
	_, err := updateNetbirdPolicy(ctx, UpdateNetbirdPolicyParams{
		PolicyID: "policy-123",
		Rules:    &rules,
	})
	
	if err == nil {
		t.Error("expected validation error")
	}
	if err != nil && !contains(err.Error(), "validation error") {
		t.Errorf("error should mention validation: %v", err)
	}
}

func TestUpdateNetbirdPolicy_WithoutRules(t *testing.T) {
	// Mock response data
	mockResp := NetbirdPolicy{
		ID:      "policy-999",
		Name:    "Updated Name Only",
		Enabled: true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies/policy-999" || r.Method != "PUT" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	name := "Updated Name Only"
	
	policy, err := updateNetbirdPolicy(ctx, UpdateNetbirdPolicyParams{
		PolicyID: "policy-999",
		Name:     &name,
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.Name != "Updated Name Only" {
		t.Errorf("unexpected policy name: %s", policy.Name)
	}
}

// Unit test for GetPolicyTemplate
func TestGetPolicyTemplate(t *testing.T) {
	template := GetPolicyTemplate()
	
	// Verify template contains all required top-level fields
	if template["name"] == nil {
		t.Error("template missing 'name' field")
	}
	if template["description"] == nil {
		t.Error("template missing 'description' field")
	}
	if template["enabled"] == nil {
		t.Error("template missing 'enabled' field")
	}
	if template["rules"] == nil {
		t.Error("template missing 'rules' field")
	}
	
	// Verify rules is an array
	rules, ok := template["rules"].([]map[string]interface{})
	if !ok {
		t.Fatal("rules field is not an array of maps")
	}
	
	// Verify template includes multiple examples
	if len(rules) < 2 {
		t.Errorf("expected at least 2 example rules, got %d", len(rules))
	}
	
	// Verify each rule has required fields
	for i, rule := range rules {
		requiredFields := []string{"name", "enabled", "action", "bidirectional", "protocol"}
		for _, field := range requiredFields {
			if rule[field] == nil {
				t.Errorf("rule %d missing required field '%s'", i, field)
			}
		}
		
		// Verify at least one source
		hasSources := rule["sources"] != nil
		hasSourceResource := rule["sourceResource"] != nil
		if !hasSources && !hasSourceResource {
			t.Errorf("rule %d missing sources or sourceResource", i)
		}
		
		// Verify at least one destination
		hasDestinations := rule["destinations"] != nil
		hasDestinationResource := rule["destinationResource"] != nil
		if !hasDestinations && !hasDestinationResource {
			t.Errorf("rule %d missing destinations or destinationResource", i)
		}
	}
	
	// Verify template is valid (passes validation)
	err := ValidatePolicyRules(rules)
	if err != nil {
		t.Errorf("template rules failed validation: %v", err)
	}
	
	// Verify template includes simple example (with ports)
	hasSimpleExample := false
	for _, rule := range rules {
		if rule["ports"] != nil {
			hasSimpleExample = true
			break
		}
	}
	if !hasSimpleExample {
		t.Error("template missing simple example with ports")
	}
	
	// Verify template includes complex example (with port_ranges or authorized_groups)
	hasComplexExample := false
	for _, rule := range rules {
		if rule["port_ranges"] != nil || rule["authorized_groups"] != nil {
			hasComplexExample = true
			break
		}
	}
	if !hasComplexExample {
		t.Error("template missing complex example with port_ranges or authorized_groups")
	}
	
	// Verify template includes resource reference example
	hasResourceExample := false
	for _, rule := range rules {
		if rule["sourceResource"] != nil || rule["destinationResource"] != nil {
			hasResourceExample = true
			break
		}
	}
	if !hasResourceExample {
		t.Error("template missing resource reference example")
	}
}
