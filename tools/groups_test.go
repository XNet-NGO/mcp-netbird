package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
)

func TestGroupResourceMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		resource GroupResource
		expected string
	}{
		{
			name: "basic group resource",
			resource: GroupResource{
				ID:   "res-123",
				Type: "network",
			},
			expected: `{"id":"res-123","type":"network"}`,
		},
		{
			name: "group resource with peer type",
			resource: GroupResource{
				ID:   "peer-456",
				Type: "peer",
			},
			expected: `{"id":"peer-456","type":"peer"}`,
		},
		{
			name: "empty group resource",
			resource: GroupResource{
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
			var decoded GroupResource
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

func TestGroupResourceUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected GroupResource
		wantErr  bool
	}{
		{
			name: "valid group resource",
			json: `{"id":"res-123","type":"network"}`,
			expected: GroupResource{
				ID:   "res-123",
				Type: "network",
			},
			wantErr: false,
		},
		{
			name: "empty values",
			json: `{"id":"","type":""}`,
			expected: GroupResource{
				ID:   "",
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "missing fields defaults to empty strings",
			json: `{}`,
			expected: GroupResource{
				ID:   "",
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "partial fields - only id",
			json: `{"id":"res-789"}`,
			expected: GroupResource{
				ID:   "res-789",
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "partial fields - only type",
			json: `{"type":"peer"}`,
			expected: GroupResource{
				ID:   "",
				Type: "peer",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resource GroupResource
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

func TestCreateNetbirdGroupParamsMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		params   CreateNetbirdGroupParams
		expected string
	}{
		{
			name: "with name only",
			params: CreateNetbirdGroupParams{
				Name: "test-group",
			},
			expected: `{"name":"test-group"}`,
		},
		{
			name: "with name and peers",
			params: CreateNetbirdGroupParams{
				Name:  "test-group",
				Peers: &[]string{"peer-1", "peer-2"},
			},
			expected: `{"name":"test-group","peers":["peer-1","peer-2"]}`,
		},
		{
			name: "with name and resources",
			params: CreateNetbirdGroupParams{
				Name: "test-group",
				Resources: &[]GroupResource{
					{ID: "res-1", Type: "network"},
					{ID: "res-2", Type: "peer"},
				},
			},
			expected: `{"name":"test-group","resources":[{"id":"res-1","type":"network"},{"id":"res-2","type":"peer"}]}`,
		},
		{
			name: "with name, peers, and resources",
			params: CreateNetbirdGroupParams{
				Name:  "test-group",
				Peers: &[]string{"peer-1"},
				Resources: &[]GroupResource{
					{ID: "res-1", Type: "network"},
				},
			},
			expected: `{"name":"test-group","peers":["peer-1"],"resources":[{"id":"res-1","type":"network"}]}`,
		},
		{
			name: "with empty peers slice",
			params: CreateNetbirdGroupParams{
				Name:  "test-group",
				Peers: &[]string{},
			},
			expected: `{"name":"test-group","peers":[]}`,
		},
		{
			name: "with empty resources slice",
			params: CreateNetbirdGroupParams{
				Name:      "test-group",
				Resources: &[]GroupResource{},
			},
			expected: `{"name":"test-group","resources":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed: got %s, want %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded CreateNetbirdGroupParams
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.Name != tt.params.Name {
				t.Errorf("name mismatch: got %s, want %s", decoded.Name, tt.params.Name)
			}
			
			// Check peers
			if tt.params.Peers != nil {
				if decoded.Peers == nil {
					t.Error("peers should not be nil")
				} else if len(*decoded.Peers) != len(*tt.params.Peers) {
					t.Errorf("peers length mismatch: got %d, want %d", len(*decoded.Peers), len(*tt.params.Peers))
				}
			}
			
			// Check resources
			if tt.params.Resources != nil {
				if decoded.Resources == nil {
					t.Error("resources should not be nil")
				} else if len(*decoded.Resources) != len(*tt.params.Resources) {
					t.Errorf("resources length mismatch: got %d, want %d", len(*decoded.Resources), len(*tt.params.Resources))
				}
			}
		})
	}
}

func TestCreateNetbirdGroupParamsOmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		params   CreateNetbirdGroupParams
		contains []string
		notContains []string
	}{
		{
			name: "nil peers and resources are omitted",
			params: CreateNetbirdGroupParams{
				Name: "test-group",
			},
			contains: []string{"name"},
			notContains: []string{"peers", "resources"},
		},
		{
			name: "nil peers is omitted but resources is present",
			params: CreateNetbirdGroupParams{
				Name: "test-group",
				Resources: &[]GroupResource{
					{ID: "res-1", Type: "network"},
				},
			},
			contains: []string{"name", "resources"},
			notContains: []string{"peers"},
		},
		{
			name: "nil resources is omitted but peers is present",
			params: CreateNetbirdGroupParams{
				Name: "test-group",
				Peers: &[]string{"peer-1"},
			},
			contains: []string{"name", "peers"},
			notContains: []string{"resources"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			
			jsonStr := string(data)
			
			// Check that expected fields are present
			for _, field := range tt.contains {
				if !strings.Contains(jsonStr, `"`+field+`"`) {
					t.Errorf("expected field %q to be present in JSON: %s", field, jsonStr)
				}
			}
			
			// Check that unexpected fields are not present
			for _, field := range tt.notContains {
				if strings.Contains(jsonStr, `"`+field+`"`) {
					t.Errorf("expected field %q to be omitted from JSON: %s", field, jsonStr)
				}
			}
		})
	}
}

func TestUpdateNetbirdGroupParamsMarshaling(t *testing.T) {
	groupID := "group-123"
	name := "updated-group"
	peers := []string{"peer-1", "peer-2"}
	resources := []GroupResource{
		{ID: "res-1", Type: "network"},
		{ID: "res-2", Type: "peer"},
	}

	tests := []struct {
		name     string
		params   UpdateNetbirdGroupParams
		expected string
	}{
		{
			name: "with group_id only",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
			},
			expected: `{"group_id":"group-123"}`,
		},
		{
			name: "with group_id and name",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
				Name:    &name,
			},
			expected: `{"group_id":"group-123","name":"updated-group"}`,
		},
		{
			name: "with group_id and peers",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
				Peers:   &peers,
			},
			expected: `{"group_id":"group-123","peers":["peer-1","peer-2"]}`,
		},
		{
			name: "with group_id and resources",
			params: UpdateNetbirdGroupParams{
				GroupID:   groupID,
				Resources: &resources,
			},
			expected: `{"group_id":"group-123","resources":[{"id":"res-1","type":"network"},{"id":"res-2","type":"peer"}]}`,
		},
		{
			name: "with all fields",
			params: UpdateNetbirdGroupParams{
				GroupID:   groupID,
				Name:      &name,
				Peers:     &peers,
				Resources: &resources,
			},
			expected: `{"group_id":"group-123","name":"updated-group","peers":["peer-1","peer-2"],"resources":[{"id":"res-1","type":"network"},{"id":"res-2","type":"peer"}]}`,
		},
		{
			name: "with empty peers slice",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
				Peers:   &[]string{},
			},
			expected: `{"group_id":"group-123","peers":[]}`,
		},
		{
			name: "with empty resources slice",
			params: UpdateNetbirdGroupParams{
				GroupID:   groupID,
				Resources: &[]GroupResource{},
			},
			expected: `{"group_id":"group-123","resources":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed: got %s, want %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded UpdateNetbirdGroupParams
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.GroupID != tt.params.GroupID {
				t.Errorf("group_id mismatch: got %s, want %s", decoded.GroupID, tt.params.GroupID)
			}

			// Check name
			if tt.params.Name != nil {
				if decoded.Name == nil {
					t.Error("name should not be nil")
				} else if *decoded.Name != *tt.params.Name {
					t.Errorf("name mismatch: got %s, want %s", *decoded.Name, *tt.params.Name)
				}
			}

			// Check peers
			if tt.params.Peers != nil {
				if decoded.Peers == nil {
					t.Error("peers should not be nil")
				} else if len(*decoded.Peers) != len(*tt.params.Peers) {
					t.Errorf("peers length mismatch: got %d, want %d", len(*decoded.Peers), len(*tt.params.Peers))
				}
			}

			// Check resources
			if tt.params.Resources != nil {
				if decoded.Resources == nil {
					t.Error("resources should not be nil")
				} else if len(*decoded.Resources) != len(*tt.params.Resources) {
					t.Errorf("resources length mismatch: got %d, want %d", len(*decoded.Resources), len(*tt.params.Resources))
				}
			}
		})
	}
}

func TestUpdateNetbirdGroupParamsOmitEmpty(t *testing.T) {
	groupID := "group-123"
	name := "updated-group"
	peers := []string{"peer-1"}
	resources := []GroupResource{{ID: "res-1", Type: "network"}}

	tests := []struct {
		name        string
		params      UpdateNetbirdGroupParams
		contains    []string
		notContains []string
	}{
		{
			name: "nil optional fields are omitted",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
			},
			contains:    []string{"group_id"},
			notContains: []string{"name", "peers", "resources"},
		},
		{
			name: "only name is present",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
				Name:    &name,
			},
			contains:    []string{"group_id", "name"},
			notContains: []string{"peers", "resources"},
		},
		{
			name: "only peers is present",
			params: UpdateNetbirdGroupParams{
				GroupID: groupID,
				Peers:   &peers,
			},
			contains:    []string{"group_id", "peers"},
			notContains: []string{"name", "resources"},
		},
		{
			name: "only resources is present",
			params: UpdateNetbirdGroupParams{
				GroupID:   groupID,
				Resources: &resources,
			},
			contains:    []string{"group_id", "resources"},
			notContains: []string{"name", "peers"},
		},
		{
			name: "name and resources present, peers omitted",
			params: UpdateNetbirdGroupParams{
				GroupID:   groupID,
				Name:      &name,
				Resources: &resources,
			},
			contains:    []string{"group_id", "name", "resources"},
			notContains: []string{"peers"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			jsonStr := string(data)

			// Check that expected fields are present
			for _, field := range tt.contains {
				if !strings.Contains(jsonStr, `"`+field+`"`) {
					t.Errorf("expected field %q to be present in JSON: %s", field, jsonStr)
				}
			}

			// Check that unexpected fields are not present
			for _, field := range tt.notContains {
				if strings.Contains(jsonStr, `"`+field+`"`) {
					t.Errorf("expected field %q to be omitted from JSON: %s", field, jsonStr)
				}
			}
		})
	}
}

// Property Test 4: Group Dependency Discovery Completeness
// Validates: Requirements 4.1, 4.2, 4.3, 4.5
// Feature: mcp-netbird-improvements, Property 4: Group Dependency Discovery Completeness
func TestListPoliciesByGroup_PropertyCompletenessDiscovery(t *testing.T) {
	// Run 50+ iterations with random policies
	for i := 0; i < 50; i++ {
		// Generate random policies with groups in various locations
		targetGroupID := "target-group"
		policies := generateRandomPoliciesWithGroup(i, targetGroupID)
		
		// Count expected references
		expectedCount := 0
		for _, policy := range policies {
			for _, rule := range policy.Rules {
				// Count sources
				for _, source := range rule.Sources {
					if source.ID == targetGroupID {
						expectedCount++
						break
					}
				}
				// Count destinations
				for _, dest := range rule.Destinations {
					if dest.ID == targetGroupID {
						expectedCount++
						break
					}
				}
				// Count authorized_groups
				if rule.AuthorizedGroups != nil {
					for authGroupID := range *rule.AuthorizedGroups {
						if authGroupID == targetGroupID {
							expectedCount++
							break
						}
					}
				}
			}
		}
		
		// Mock the API response
		mockServer := createMockPolicyServer(policies)
		defer mockServer.Close()
		
		mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
		defer func() { mcpnetbird.TestNetbirdClient = nil }()
		
		ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
		
		// Search for the group
		references, err := ListPoliciesByGroup(ctx, targetGroupID)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		
		// Verify all occurrences found
		if len(references) != expectedCount {
			t.Errorf("iteration %d: expected %d references, got %d", i, expectedCount, len(references))
		}
		
		// Verify each reference has correct location
		for _, ref := range references {
			if ref.Location != "sources" && ref.Location != "destinations" && ref.Location != "authorized_groups" {
				t.Errorf("iteration %d: invalid location: %s", i, ref.Location)
			}
		}
	}
}

// Helper function to generate random policies with a specific group
func generateRandomPoliciesWithGroup(seed int, groupID string) []NetbirdPolicy {
	numPolicies := (seed % 3) + 1
	policies := make([]NetbirdPolicy, numPolicies)
	
	for i := 0; i < numPolicies; i++ {
		numRules := (seed % 2) + 1
		rules := make([]NetbirdPolicyRule, numRules)
		
		for j := 0; j < numRules; j++ {
			rule := NetbirdPolicyRule{
				ID:            fmt.Sprintf("rule-%d-%d", i, j),
				Name:          fmt.Sprintf("rule-%d-%d", i, j),
				Enabled:       true,
				Action:        "accept",
				Bidirectional: false,
				Protocol:      "tcp",
			}
			
			// Vary where the group appears based on seed
			switch (seed + i + j) % 4 {
			case 0:
				// Group in sources
				rule.Sources = []NetbirdPeerGroup{
					{ID: groupID, Name: "Target Group"},
					{ID: "other-group", Name: "Other Group"},
				}
				rule.Destinations = []NetbirdPeerGroup{
					{ID: "dest-group", Name: "Dest Group"},
				}
			case 1:
				// Group in destinations
				rule.Sources = []NetbirdPeerGroup{
					{ID: "src-group", Name: "Src Group"},
				}
				rule.Destinations = []NetbirdPeerGroup{
					{ID: groupID, Name: "Target Group"},
					{ID: "other-group", Name: "Other Group"},
				}
			case 2:
				// Group in authorized_groups
				rule.Sources = []NetbirdPeerGroup{
					{ID: "src-group", Name: "Src Group"},
				}
				rule.Destinations = []NetbirdPeerGroup{
					{ID: "dest-group", Name: "Dest Group"},
				}
				authGroups := map[string][]string{
					groupID:      {"user1", "user2"},
					"other-auth": {"user3"},
				}
				rule.AuthorizedGroups = &authGroups
			case 3:
				// Group not present in this rule
				rule.Sources = []NetbirdPeerGroup{
					{ID: "src-group", Name: "Src Group"},
				}
				rule.Destinations = []NetbirdPeerGroup{
					{ID: "dest-group", Name: "Dest Group"},
				}
			}
			
			rules[j] = rule
		}
		
		policies[i] = NetbirdPolicy{
			ID:      fmt.Sprintf("policy-%d", i),
			Name:    fmt.Sprintf("Policy %d", i),
			Enabled: true,
			Rules:   rules,
		}
	}
	
	return policies
}

// Helper function to create a mock server for policies
func createMockPolicyServer(policies []NetbirdPolicy) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/policies" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(policies)
	}))
}

// Unit tests for ListPoliciesByGroup
func TestListPoliciesByGroup_InSourcesOnly(t *testing.T) {
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "target-group", Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "other-group", Name: "Other Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(references) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(references))
	}
	
	if references[0].Location != "sources" {
		t.Errorf("expected location 'sources', got '%s'", references[0].Location)
	}
}

func TestListPoliciesByGroup_InDestinationsOnly(t *testing.T) {
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "target-group", Name: "Target Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(references) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(references))
	}
	
	if references[0].Location != "destinations" {
		t.Errorf("expected location 'destinations', got '%s'", references[0].Location)
	}
}

func TestListPoliciesByGroup_InAuthorizedGroupsOnly(t *testing.T) {
	authGroups := map[string][]string{
		"target-group": {"user1", "user2"},
	}
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "src-group", Name: "Src Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
					AuthorizedGroups: &authGroups,
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(references) != 1 {
		t.Fatalf("expected 1 reference, got %d", len(references))
	}
	
	if references[0].Location != "authorized_groups" {
		t.Errorf("expected location 'authorized_groups', got '%s'", references[0].Location)
	}
}

func TestListPoliciesByGroup_InMultipleLocations(t *testing.T) {
	authGroups := map[string][]string{
		"target-group": {"user1"},
	}
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "target-group", Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "target-group", Name: "Target Group"},
					},
					AuthorizedGroups: &authGroups,
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should find 3 references (sources, destinations, authorized_groups)
	if len(references) != 3 {
		t.Fatalf("expected 3 references, got %d", len(references))
	}
	
	// Verify all locations are present
	locations := make(map[string]bool)
	for _, ref := range references {
		locations[ref.Location] = true
	}
	
	if !locations["sources"] || !locations["destinations"] || !locations["authorized_groups"] {
		t.Errorf("expected all three locations, got: %v", locations)
	}
}

func TestListPoliciesByGroup_NoMatches(t *testing.T) {
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "another-group", Name: "Another Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(references) != 0 {
		t.Errorf("expected 0 references, got %d", len(references))
	}
}

func TestListPoliciesByGroup_MultiplePolicies(t *testing.T) {
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
			Name:    "Policy 1",
			Enabled: true,
			Rules: []NetbirdPolicyRule{
				{
					ID:            "rule-1",
					Name:          "rule-1",
					Enabled:       true,
					Action:        "accept",
					Bidirectional: false,
					Protocol:      "tcp",
					Sources: []NetbirdPeerGroup{
						{ID: "target-group", Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "other-group", Name: "Other Group"},
					},
				},
			},
		},
		{
			ID:      "policy-2",
			Name:    "Policy 2",
			Enabled: true,
			Rules: []NetbirdPolicyRule{
				{
					ID:            "rule-2",
					Name:          "rule-2",
					Enabled:       true,
					Action:        "accept",
					Bidirectional: false,
					Protocol:      "tcp",
					Sources: []NetbirdPeerGroup{
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "target-group", Name: "Target Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockPolicyServer(policies)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	references, err := ListPoliciesByGroup(ctx, "target-group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(references) != 2 {
		t.Fatalf("expected 2 references, got %d", len(references))
	}
	
	// Verify both policies are referenced
	policyIDs := make(map[string]bool)
	for _, ref := range references {
		policyIDs[ref.PolicyID] = true
	}
	
	if !policyIDs["policy-1"] || !policyIDs["policy-2"] {
		t.Errorf("expected both policies to be referenced")
	}
}

// Property Test 6: Group Replacement Completeness
// Validates: Requirements 6.2, 6.3, 6.4, 6.5
// Feature: mcp-netbird-improvements, Property 6: Group Replacement Completeness
func TestReplaceGroupInPolicies_PropertyCompletenessReplacement(t *testing.T) {
	// Run 100+ iterations with random policies
	for i := 0; i < 100; i++ {
		oldGroupID := "old-group"
		newGroupID := "new-group"
		
		// Generate random policies with old group references
		policies := generateRandomPoliciesWithGroup(i, oldGroupID)
		
		// Create mock server that handles both GET and PUT requests
		mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
		defer mockServer.Close()
		
		mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
		defer func() { mcpnetbird.TestNetbirdClient = nil }()
		
		ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
		
		// Replace old group with new group
		result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		
		// Verify no errors occurred
		if len(result.Errors) > 0 {
			t.Errorf("iteration %d: unexpected errors: %v", i, result.Errors)
		}
		
		// Verify all policies with old group were updated
		expectedUpdates := countPoliciesWithGroup(policies, oldGroupID)
		if len(result.UpdatedPolicyIDs) != expectedUpdates {
			t.Errorf("iteration %d: expected %d updated policies, got %d", i, expectedUpdates, len(result.UpdatedPolicyIDs))
		}
	}
}

// Helper function to count unique policies containing a group
func countPoliciesWithGroup(policies []NetbirdPolicy, groupID string) int {
	count := 0
	for _, policy := range policies {
		hasGroup := false
		for _, rule := range policy.Rules {
			// Check sources
			for _, source := range rule.Sources {
				if source.ID == groupID {
					hasGroup = true
					break
				}
			}
			if hasGroup {
				break
			}
			// Check destinations
			for _, dest := range rule.Destinations {
				if dest.ID == groupID {
					hasGroup = true
					break
				}
			}
			if hasGroup {
				break
			}
			// Check authorized_groups
			if rule.AuthorizedGroups != nil {
				for authGroupID := range *rule.AuthorizedGroups {
					if authGroupID == groupID {
						hasGroup = true
						break
					}
				}
			}
			if hasGroup {
				break
			}
		}
		if hasGroup {
			count++
		}
	}
	return count
}

// Helper function to create a mock server that handles replacement operations
func createMockReplacementServer(policies []NetbirdPolicy, oldGroupID, newGroupID string) *httptest.Server {
	// Store updated policies
	updatedPolicies := make(map[string]NetbirdPolicy)
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Handle GET /policies
		if r.Method == "GET" && r.URL.Path == "/policies" {
			_ = json.NewEncoder(w).Encode(policies)
			return
		}
		
		// Handle GET /policies/{id}
		if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			
			// Check if we have an updated version
			if updated, ok := updatedPolicies[policyID]; ok {
				_ = json.NewEncoder(w).Encode(updated)
				return
			}
			
			// Find original policy
			for _, policy := range policies {
				if policy.ID == policyID {
					_ = json.NewEncoder(w).Encode(policy)
					return
				}
			}
			http.NotFound(w, r)
			return
		}
		
		// Handle PUT /policies/{id}
		if r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			
			var updateBody map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			// Find original policy
			var originalPolicy NetbirdPolicy
			found := false
			for _, policy := range policies {
				if policy.ID == policyID {
					originalPolicy = policy
					found = true
					break
				}
			}
			
			if !found {
				http.NotFound(w, r)
				return
			}
			
			// Create updated policy
			updatedPolicy := originalPolicy
			if name, ok := updateBody["name"].(string); ok {
				updatedPolicy.Name = name
			}
			if desc, ok := updateBody["description"].(string); ok {
				updatedPolicy.Description = desc
			}
			if enabled, ok := updateBody["enabled"].(bool); ok {
				updatedPolicy.Enabled = enabled
			}
			
			// Store the updated policy
			updatedPolicies[policyID] = updatedPolicy
			
			_ = json.NewEncoder(w).Encode(updatedPolicy)
			return
		}
		
		http.NotFound(w, r)
	}))
}

// Unit tests for ReplaceGroupInPolicies

func TestReplaceGroupInPolicies_InSources(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: oldGroupID, Name: "Old Group"},
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(result.UpdatedPolicyIDs) != 1 {
		t.Errorf("expected 1 updated policy, got %d", len(result.UpdatedPolicyIDs))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
	
	if result.UpdatedPolicyIDs[0] != "policy-1" {
		t.Errorf("expected policy-1 to be updated, got %s", result.UpdatedPolicyIDs[0])
	}
}

func TestReplaceGroupInPolicies_InDestinations(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "src-group", Name: "Src Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: oldGroupID, Name: "Old Group"},
						{ID: "other-group", Name: "Other Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(result.UpdatedPolicyIDs) != 1 {
		t.Errorf("expected 1 updated policy, got %d", len(result.UpdatedPolicyIDs))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestReplaceGroupInPolicies_InAuthorizedGroups(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	authGroups := map[string][]string{
		oldGroupID:   {"user1", "user2"},
		"other-auth": {"user3"},
	}
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "src-group", Name: "Src Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
					AuthorizedGroups: &authGroups,
				},
			},
		},
	}
	
	mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(result.UpdatedPolicyIDs) != 1 {
		t.Errorf("expected 1 updated policy, got %d", len(result.UpdatedPolicyIDs))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestReplaceGroupInPolicies_NoMatches(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "src-group", Name: "Src Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(result.UpdatedPolicyIDs) != 0 {
		t.Errorf("expected 0 updated policies, got %d", len(result.UpdatedPolicyIDs))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestReplaceGroupInPolicies_MultipleOccurrencesInSameRule(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	authGroups := map[string][]string{
		oldGroupID: {"user1"},
	}
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: oldGroupID, Name: "Old Group"},
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: oldGroupID, Name: "Old Group"},
					},
					AuthorizedGroups: &authGroups,
				},
			},
		},
	}
	
	mockServer := createMockReplacementServer(policies, oldGroupID, newGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should only update the policy once, even though group appears multiple times
	if len(result.UpdatedPolicyIDs) != 1 {
		t.Errorf("expected 1 updated policy, got %d", len(result.UpdatedPolicyIDs))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestReplaceGroupInPolicies_PartialFailure(t *testing.T) {
	oldGroupID := "old-group"
	newGroupID := "new-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
			Name:    "Test Policy 1",
			Enabled: true,
			Rules: []NetbirdPolicyRule{
				{
					ID:            "rule-1",
					Name:          "test-rule-1",
					Enabled:       true,
					Action:        "accept",
					Bidirectional: false,
					Protocol:      "tcp",
					Sources: []NetbirdPeerGroup{
						{ID: oldGroupID, Name: "Old Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
		{
			ID:      "policy-2",
			Name:    "Test Policy 2",
			Enabled: true,
			Rules: []NetbirdPolicyRule{
				{
					ID:            "rule-2",
					Name:          "test-rule-2",
					Enabled:       true,
					Action:        "accept",
					Bidirectional: false,
					Protocol:      "tcp",
					Sources: []NetbirdPeerGroup{
						{ID: oldGroupID, Name: "Old Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	// Create a mock server that fails for policy-2
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == "GET" && r.URL.Path == "/policies" {
			_ = json.NewEncoder(w).Encode(policies)
			return
		}
		
		if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			for _, policy := range policies {
				if policy.ID == policyID {
					_ = json.NewEncoder(w).Encode(policy)
					return
				}
			}
			http.NotFound(w, r)
			return
		}
		
		if r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			
			// Fail for policy-2
			if policyID == "policy-2" {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			
			// Succeed for policy-1
			for _, policy := range policies {
				if policy.ID == policyID {
					_ = json.NewEncoder(w).Encode(policy)
					return
				}
			}
			http.NotFound(w, r)
			return
		}
		
		http.NotFound(w, r)
	}))
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := ReplaceGroupInPolicies(ctx, oldGroupID, newGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should have updated policy-1 successfully
	if len(result.UpdatedPolicyIDs) != 1 {
		t.Errorf("expected 1 updated policy, got %d", len(result.UpdatedPolicyIDs))
	}
	
	// Should have error for policy-2
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	
	if _, ok := result.Errors["policy-2"]; !ok {
		t.Errorf("expected error for policy-2")
	}
}

// Property Test 7: Force Delete Removes Dependencies
// Validates: Requirements 5.2, 5.5
// Feature: mcp-netbird-improvements, Property 7: Force Delete Removes Dependencies
func TestDeleteGroupForce_PropertyRemovesDependencies(t *testing.T) {
	// Run 100+ iterations with random policies
	for i := 0; i < 100; i++ {
		targetGroupID := "target-group"
		
		// Generate random policies with group dependencies
		policies := generateRandomPoliciesWithGroup(i, targetGroupID)
		
		// Create mock server that handles force delete operations
		mockServer := createMockForceDeleteServer(policies, targetGroupID)
		defer mockServer.Close()
		
		mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
		defer func() { mcpnetbird.TestNetbirdClient = nil }()
		
		ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
		
		// Force delete the group
		result, err := DeleteGroupForce(ctx, targetGroupID)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		
		// Verify group was deleted
		if !result.Deleted {
			t.Errorf("iteration %d: expected group to be deleted", i)
		}
		
		// Verify policies were modified
		expectedModifications := countPoliciesWithGroup(policies, targetGroupID)
		if len(result.PoliciesModified) != expectedModifications {
			t.Errorf("iteration %d: expected %d modified policies, got %d", i, expectedModifications, len(result.PoliciesModified))
		}
	}
}

// Helper function to create a mock server for force delete operations
func createMockForceDeleteServer(policies []NetbirdPolicy, groupID string) *httptest.Server {
	// Track updated and deleted policies
	updatedPolicies := make(map[string]NetbirdPolicy)
	deletedPolicies := make(map[string]bool)
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Handle GET /policies
		if r.Method == "GET" && r.URL.Path == "/policies" {
			_ = json.NewEncoder(w).Encode(policies)
			return
		}
		
		// Handle GET /policies/{id}
		if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			
			// Check if deleted
			if deletedPolicies[policyID] {
				http.NotFound(w, r)
				return
			}
			
			// Check if updated
			if updated, ok := updatedPolicies[policyID]; ok {
				_ = json.NewEncoder(w).Encode(updated)
				return
			}
			
			// Find original policy
			for _, policy := range policies {
				if policy.ID == policyID {
					_ = json.NewEncoder(w).Encode(policy)
					return
				}
			}
			http.NotFound(w, r)
			return
		}
		
		// Handle PUT /policies/{id}
		if r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			
			var updateBody map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			// Find original policy
			var originalPolicy NetbirdPolicy
			found := false
			for _, policy := range policies {
				if policy.ID == policyID {
					originalPolicy = policy
					found = true
					break
				}
			}
			
			if !found {
				http.NotFound(w, r)
				return
			}
			
			// Create updated policy
			updatedPolicy := originalPolicy
			if name, ok := updateBody["name"].(string); ok {
				updatedPolicy.Name = name
			}
			if desc, ok := updateBody["description"].(string); ok {
				updatedPolicy.Description = desc
			}
			if enabled, ok := updateBody["enabled"].(bool); ok {
				updatedPolicy.Enabled = enabled
			}
			
			// Store the updated policy
			updatedPolicies[policyID] = updatedPolicy
			
			_ = json.NewEncoder(w).Encode(updatedPolicy)
			return
		}
		
		// Handle DELETE /policies/{id}
		if r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/policies/") {
			policyID := strings.TrimPrefix(r.URL.Path, "/policies/")
			deletedPolicies[policyID] = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		// Handle DELETE /groups/{id}
		if r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/groups/") {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		http.NotFound(w, r)
	}))
}

// Unit tests for DeleteGroupForce and delete_netbird_group with force parameter

func TestDeleteGroupForce_WithDependencies(t *testing.T) {
	targetGroupID := "target-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: targetGroupID, Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockForceDeleteServer(policies, targetGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := DeleteGroupForce(ctx, targetGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Deleted {
		t.Error("expected group to be deleted")
	}
	
	if len(result.PoliciesModified) != 1 {
		t.Errorf("expected 1 modified policy, got %d", len(result.PoliciesModified))
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("unexpected errors: %v", result.Errors)
	}
}

func TestDeleteNetbirdGroup_NormalDeleteWithDependencies(t *testing.T) {
	targetGroupID := "target-group"
	
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: targetGroupID, Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockForceDeleteServer(policies, targetGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	// Try normal delete (force=false) with dependencies - should fail
	_, err := deleteNetbirdGroup(ctx, DeleteNetbirdGroupParams{
		GroupID: targetGroupID,
		Force:   false,
	})
	
	if err == nil {
		t.Fatal("expected error when deleting group with dependencies")
	}
	
	if !strings.Contains(err.Error(), "referenced by") {
		t.Errorf("expected error message to mention dependencies, got: %v", err)
	}
}

func TestDeleteNetbirdGroup_WithoutDependencies(t *testing.T) {
	targetGroupID := "target-group"
	
	// No policies reference the target group
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: "other-group", Name: "Other Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: "dest-group", Name: "Dest Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockForceDeleteServer(policies, targetGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	// Normal delete should succeed
	result, err := deleteNetbirdGroup(ctx, DeleteNetbirdGroupParams{
		GroupID: targetGroupID,
		Force:   false,
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	status, ok := result["status"].(string)
	if !ok || status != "deleted" {
		t.Errorf("expected status 'deleted', got %v", result["status"])
	}
}

func TestDeleteGroupForce_CleanupInvalidRules(t *testing.T) {
	targetGroupID := "target-group"
	
	// Policy with a rule that will become invalid after removing the group
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: targetGroupID, Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: targetGroupID, Name: "Target Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockForceDeleteServer(policies, targetGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := DeleteGroupForce(ctx, targetGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Deleted {
		t.Error("expected group to be deleted")
	}
	
	// Policy should be marked for deletion since rule becomes invalid
	if len(result.PoliciesModified) != 1 {
		t.Errorf("expected 1 modified policy, got %d", len(result.PoliciesModified))
	}
}

func TestDeleteGroupForce_CleanupEmptyPolicies(t *testing.T) {
	targetGroupID := "target-group"
	
	// Policy with only one rule that references the target group
	policies := []NetbirdPolicy{
		{
			ID:      "policy-1",
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
						{ID: targetGroupID, Name: "Target Group"},
					},
					Destinations: []NetbirdPeerGroup{
						{ID: targetGroupID, Name: "Target Group"},
					},
				},
			},
		},
	}
	
	mockServer := createMockForceDeleteServer(policies, targetGroupID)
	defer mockServer.Close()
	
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(mockServer.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()
	
	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	
	result, err := DeleteGroupForce(ctx, targetGroupID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Deleted {
		t.Error("expected group to be deleted")
	}
	
	// Policy should be deleted since it becomes empty
	if len(result.PoliciesModified) != 1 {
		t.Errorf("expected 1 modified policy, got %d", len(result.PoliciesModified))
	}
}
