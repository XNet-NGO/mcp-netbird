package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/aantti/mcp-netbird"
)

func TestListNetbirdNetworkResources(t *testing.T) {
	description := "Test resource description"
	// Mock response data
	mockResp := []NetbirdNetworkResource{
		{
			ID:          "res1",
			Type:        "host",
			Name:        "Test Resource",
			Description: &description,
			Address:     "1.1.1.1",
			Enabled:     true,
			Groups: []NetbirdNetworkResourceGroup{
				{
					ID:             "grp1",
					Name:           "Test Group",
					PeersCount:     2,
					ResourcesCount: 5,
					Issued:         "api",
				},
			},
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	resources, err := listNetbirdNetworkResources(ctx, ListNetbirdNetworkResourcesParams{NetworkID: "net1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	if resources[0].ID != "res1" {
		t.Errorf("expected ID 'res1', got '%s'", resources[0].ID)
	}
	if resources[0].Type != "host" {
		t.Errorf("expected Type 'host', got '%s'", resources[0].Type)
	}
	if resources[0].Name != "Test Resource" {
		t.Errorf("expected Name 'Test Resource', got '%s'", resources[0].Name)
	}
	if resources[0].Address != "1.1.1.1" {
		t.Errorf("expected Address '1.1.1.1', got '%s'", resources[0].Address)
	}
	if !resources[0].Enabled {
		t.Error("expected Enabled to be true")
	}
	if len(resources[0].Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(resources[0].Groups))
	}
}

func TestGetNetbirdNetworkResource(t *testing.T) {
	description := "Test resource description"
	mockResp := NetbirdNetworkResource{
		ID:          "res1",
		Type:        "host",
		Name:        "Test Resource",
		Description: &description,
		Address:     "192.168.1.1",
		Enabled:     true,
		Groups: []NetbirdNetworkResourceGroup{
			{
				ID:             "grp1",
				Name:           "Test Group",
				PeersCount:     3,
				ResourcesCount: 10,
				Issued:         "api",
			},
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources/res1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	resource, err := getNetbirdNetworkResource(ctx, GetNetbirdNetworkResourceParams{
		NetworkID:  "net1",
		ResourceID: "res1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.ID != "res1" {
		t.Errorf("expected ID 'res1', got '%s'", resource.ID)
	}
	if resource.Type != "host" {
		t.Errorf("expected Type 'host', got '%s'", resource.Type)
	}
	if resource.Name != "Test Resource" {
		t.Errorf("expected Name 'Test Resource', got '%s'", resource.Name)
	}
	if resource.Description == nil || *resource.Description != "Test resource description" {
		t.Errorf("expected Description 'Test resource description', got %v", resource.Description)
	}
	if resource.Address != "192.168.1.1" {
		t.Errorf("expected Address '192.168.1.1', got '%s'", resource.Address)
	}
	if !resource.Enabled {
		t.Error("expected Enabled to be true")
	}
	if len(resource.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(resource.Groups))
	}
}

func TestGetNetbirdNetworkResource_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := getNetbirdNetworkResource(ctx, GetNetbirdNetworkResourceParams{
		NetworkID:  "net1",
		ResourceID: "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for non-existent resource, got nil")
	}
}

func TestCreateNetbirdNetworkResource(t *testing.T) {
	description := "Test resource description"
	mockResp := NetbirdNetworkResource{
		ID:          "res1",
		Type:        "host",
		Name:        "Test Resource",
		Description: &description,
		Address:     "10.0.0.1",
		Enabled:     true,
		Groups: []NetbirdNetworkResourceGroup{
			{
				ID:             "grp1",
				Name:           "Test Group",
				PeersCount:     2,
				ResourcesCount: 5,
				Issued:         "api",
			},
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req["name"] != "Test Resource" {
			t.Errorf("expected Name 'Test Resource', got '%v'", req["name"])
		}
		if req["description"] != "Test resource description" {
			t.Errorf("expected Description 'Test resource description', got '%v'", req["description"])
		}
		if req["address"] != "10.0.0.1" {
			t.Errorf("expected Address '10.0.0.1', got '%v'", req["address"])
		}
		if req["enabled"] != true {
			t.Errorf("expected Enabled true, got %v", req["enabled"])
		}
		groups, ok := req["groups"].([]interface{})
		if !ok || len(groups) != 1 {
			t.Errorf("expected 1 group, got %v", req["groups"])
		}
		// Verify network_id is NOT in the request body
		if _, exists := req["network_id"]; exists {
			t.Error("network_id should not be in request body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	resource, err := createNetbirdNetworkResource(ctx, CreateNetbirdNetworkResourceParams{
		NetworkID:   "net1",
		Name:        "Test Resource",
		Description: &description,
		Address:     "10.0.0.1",
		Enabled:     true,
		Groups:      []string{"grp1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.ID != "res1" {
		t.Errorf("expected ID 'res1', got '%s'", resource.ID)
	}
	if resource.Name != "Test Resource" {
		t.Errorf("expected Name 'Test Resource', got '%s'", resource.Name)
	}
	if resource.Description == nil || *resource.Description != "Test resource description" {
		t.Errorf("expected Description 'Test resource description', got %v", resource.Description)
	}
	if resource.Address != "10.0.0.1" {
		t.Errorf("expected Address '10.0.0.1', got '%s'", resource.Address)
	}
	if !resource.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestCreateNetbirdNetworkResource_WithoutDescription(t *testing.T) {
	mockResp := NetbirdNetworkResource{
		ID:          "res2",
		Type:        "subnet",
		Name:        "Minimal Resource",
		Description: nil,
		Address:     "192.168.0.0/24",
		Enabled:     false,
		Groups:      []NetbirdNetworkResourceGroup{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body - description should be omitted
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if _, exists := req["description"]; exists {
			t.Error("expected description to be omitted when nil")
		}
		if req["name"] != "Minimal Resource" {
			t.Errorf("expected Name 'Minimal Resource', got '%v'", req["name"])
		}
		if req["address"] != "192.168.0.0/24" {
			t.Errorf("expected Address '192.168.0.0/24', got '%v'", req["address"])
		}
		if req["enabled"] != false {
			t.Errorf("expected Enabled false, got %v", req["enabled"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	resource, err := createNetbirdNetworkResource(ctx, CreateNetbirdNetworkResourceParams{
		NetworkID: "net1",
		Name:      "Minimal Resource",
		Address:   "192.168.0.0/24",
		Enabled:   false,
		Groups:    []string{"grp1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.ID != "res2" {
		t.Errorf("expected ID 'res2', got '%s'", resource.ID)
	}
	if resource.Description != nil {
		t.Errorf("expected Description to be nil, got %v", resource.Description)
	}
}

func TestCreateNetbirdNetworkResource_WithDomain(t *testing.T) {
	description := "Domain resource"
	mockResp := NetbirdNetworkResource{
		ID:          "res3",
		Type:        "domain",
		Name:        "Domain Resource",
		Description: &description,
		Address:     "*.example.com",
		Enabled:     true,
		Groups:      []NetbirdNetworkResourceGroup{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources" {
			http.NotFound(w, r)
			return
		}

		// Verify request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req["address"] != "*.example.com" {
			t.Errorf("expected Address '*.example.com', got '%v'", req["address"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	resource, err := createNetbirdNetworkResource(ctx, CreateNetbirdNetworkResourceParams{
		NetworkID:   "net1",
		Name:        "Domain Resource",
		Description: &description,
		Address:     "*.example.com",
		Enabled:     true,
		Groups:      []string{"grp1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.Address != "*.example.com" {
		t.Errorf("expected Address '*.example.com', got '%s'", resource.Address)
	}
}

func TestUpdateNetbirdNetworkResource(t *testing.T) {
	newName := "Updated Resource"
	newDescription := "Updated description"
	newAddress := "10.0.0.2"
	newEnabled := false
	mockResp := NetbirdNetworkResource{
		ID:          "res1",
		Type:        "host",
		Name:        newName,
		Description: &newDescription,
		Address:     newAddress,
		Enabled:     newEnabled,
		Groups: []NetbirdNetworkResourceGroup{
			{
				ID:             "grp2",
				Name:           "Updated Group",
				PeersCount:     5,
				ResourcesCount: 10,
				Issued:         "api",
			},
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources/res1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req["name"] != "Updated Resource" {
			t.Errorf("expected Name 'Updated Resource', got '%v'", req["name"])
		}
		if req["description"] != "Updated description" {
			t.Errorf("expected Description 'Updated description', got '%v'", req["description"])
		}
		if req["address"] != "10.0.0.2" {
			t.Errorf("expected Address '10.0.0.2', got '%v'", req["address"])
		}
		if req["enabled"] != false {
			t.Errorf("expected Enabled false, got %v", req["enabled"])
		}
		// Verify network_id and resource_id are NOT in the request body
		if _, exists := req["network_id"]; exists {
			t.Error("network_id should not be in request body")
		}
		if _, exists := req["resource_id"]; exists {
			t.Error("resource_id should not be in request body")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	resource, err := updateNetbirdNetworkResource(ctx, UpdateNetbirdNetworkResourceParams{
		NetworkID:   "net1",
		ResourceID:  "res1",
		Name:        &newName,
		Description: &newDescription,
		Address:     &newAddress,
		Enabled:     &newEnabled,
		Groups:      []string{"grp2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.ID != "res1" {
		t.Errorf("expected ID 'res1', got '%s'", resource.ID)
	}
	if resource.Name != "Updated Resource" {
		t.Errorf("expected Name 'Updated Resource', got '%s'", resource.Name)
	}
	if resource.Description == nil || *resource.Description != "Updated description" {
		t.Errorf("expected Description 'Updated description', got %v", resource.Description)
	}
	if resource.Address != "10.0.0.2" {
		t.Errorf("expected Address '10.0.0.2', got '%s'", resource.Address)
	}
	if resource.Enabled {
		t.Error("expected Enabled to be false")
	}
}

func TestUpdateNetbirdNetworkResource_PartialUpdate(t *testing.T) {
	newName := "Partially Updated Resource"
	mockResp := NetbirdNetworkResource{
		ID:          "res1",
		Type:        "host",
		Name:        newName,
		Description: nil,
		Address:     "10.0.0.1",
		Enabled:     true,
		Groups:      []NetbirdNetworkResourceGroup{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources/res1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body - only name should be present
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if _, exists := req["description"]; exists {
			t.Error("expected description to be omitted when nil")
		}
		if _, exists := req["address"]; exists {
			t.Error("expected address to be omitted when nil")
		}
		if _, exists := req["enabled"]; exists {
			t.Error("expected enabled to be omitted when nil")
		}
		if req["name"] == nil {
			t.Error("expected name to be present")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	resource, err := updateNetbirdNetworkResource(ctx, UpdateNetbirdNetworkResourceParams{
		NetworkID:  "net1",
		ResourceID: "res1",
		Name:       &newName,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected resource, got nil")
	}
	if resource.Name != "Partially Updated Resource" {
		t.Errorf("expected Name 'Partially Updated Resource', got '%s'", resource.Name)
	}
}

func TestDeleteNetbirdNetworkResource(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/resources/res1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	result, err := deleteNetbirdNetworkResource(ctx, DeleteNetbirdNetworkResourceParams{
		NetworkID:  "net1",
		ResourceID: "res1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result["status"] != "deleted" {
		t.Errorf("expected status 'deleted', got '%s'", result["status"])
	}
	if result["network_id"] != "net1" {
		t.Errorf("expected network_id 'net1', got '%s'", result["network_id"])
	}
	if result["resource_id"] != "res1" {
		t.Errorf("expected resource_id 'res1', got '%s'", result["resource_id"])
	}
}

func TestDeleteNetbirdNetworkResource_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := deleteNetbirdNetworkResource(ctx, DeleteNetbirdNetworkResourceParams{
		NetworkID:  "net1",
		ResourceID: "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for non-existent resource, got nil")
	}
}
