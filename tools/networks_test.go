package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
)

func TestListNetbirdNetworks(t *testing.T) {
	// Mock response data
	mockResp := []NetbirdNetwork{
		{
			ID:   "net1",
			Name: "Test Network",
			// Add other fields as needed for your struct
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks" {
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
	networks, err := listNetbirdNetworks(ctx, ListNetbirdNetworksParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(networks) != 1 || networks[0].ID != "net1" {
		t.Errorf("unexpected result: %+v", networks)
	}
}

func TestGetNetbirdNetwork(t *testing.T) {
	description := "Test network description"
	mockResp := NetbirdNetwork{
		ID:                "net1",
		Name:              "Test Network",
		Description:       &description,
		Routers:           []string{"router1", "router2"},
		RoutingPeersCount: 2,
		Resources:         []string{"res1", "res2"},
		Policies:          []string{"pol1"},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1" {
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
	network, err := getNetbirdNetwork(ctx, GetNetbirdNetworkParams{NetworkID: "net1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if network == nil {
		t.Fatal("expected network, got nil")
	}
	if network.ID != "net1" {
		t.Errorf("expected ID 'net1', got '%s'", network.ID)
	}
	if network.Name != "Test Network" {
		t.Errorf("expected Name 'Test Network', got '%s'", network.Name)
	}
	if network.Description == nil || *network.Description != "Test network description" {
		t.Errorf("expected Description 'Test network description', got %v", network.Description)
	}
	if len(network.Routers) != 2 {
		t.Errorf("expected 2 routers, got %d", len(network.Routers))
	}
	if network.RoutingPeersCount != 2 {
		t.Errorf("expected RoutingPeersCount 2, got %d", network.RoutingPeersCount)
	}
	if len(network.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(network.Resources))
	}
	if len(network.Policies) != 1 {
		t.Errorf("expected 1 policy, got %d", len(network.Policies))
	}
}

func TestGetNetbirdNetwork_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := getNetbirdNetwork(ctx, GetNetbirdNetworkParams{NetworkID: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for non-existent network, got nil")
	}
}

func TestCreateNetbirdNetwork(t *testing.T) {
	description := "Test network description"
	mockResp := NetbirdNetwork{
		ID:                "net1",
		Name:              "Test Network",
		Description:       &description,
		Routers:           []string{},
		RoutingPeersCount: 0,
		Resources:         []string{},
		Policies:          []string{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body
		var req CreateNetbirdNetworkParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Name != "Test Network" {
			t.Errorf("expected Name 'Test Network', got '%s'", req.Name)
		}
		if req.Description == nil || *req.Description != "Test network description" {
			t.Errorf("expected Description 'Test network description', got %v", req.Description)
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
	network, err := createNetbirdNetwork(ctx, CreateNetbirdNetworkParams{
		Name:        "Test Network",
		Description: &description,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if network == nil {
		t.Fatal("expected network, got nil")
	}
	if network.ID != "net1" {
		t.Errorf("expected ID 'net1', got '%s'", network.ID)
	}
	if network.Name != "Test Network" {
		t.Errorf("expected Name 'Test Network', got '%s'", network.Name)
	}
	if network.Description == nil || *network.Description != "Test network description" {
		t.Errorf("expected Description 'Test network description', got %v", network.Description)
	}
}

func TestCreateNetbirdNetwork_WithoutDescription(t *testing.T) {
	mockResp := NetbirdNetwork{
		ID:                "net2",
		Name:              "Minimal Network",
		Description:       nil,
		Routers:           []string{},
		RoutingPeersCount: 0,
		Resources:         []string{},
		Policies:          []string{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks" {
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
		if req["name"] != "Minimal Network" {
			t.Errorf("expected Name 'Minimal Network', got '%v'", req["name"])
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
	network, err := createNetbirdNetwork(ctx, CreateNetbirdNetworkParams{
		Name: "Minimal Network",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if network == nil {
		t.Fatal("expected network, got nil")
	}
	if network.ID != "net2" {
		t.Errorf("expected ID 'net2', got '%s'", network.ID)
	}
	if network.Description != nil {
		t.Errorf("expected Description to be nil, got %v", network.Description)
	}
}

func TestUpdateNetbirdNetwork(t *testing.T) {
	newName := "Updated Network"
	newDescription := "Updated description"
	mockResp := NetbirdNetwork{
		ID:                "net1",
		Name:              newName,
		Description:       &newDescription,
		Routers:           []string{"router1"},
		RoutingPeersCount: 1,
		Resources:         []string{},
		Policies:          []string{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body
		var req UpdateNetbirdNetworkParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Name == nil || *req.Name != "Updated Network" {
			t.Errorf("expected Name 'Updated Network', got %v", req.Name)
		}
		if req.Description == nil || *req.Description != "Updated description" {
			t.Errorf("expected Description 'Updated description', got %v", req.Description)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	network, err := updateNetbirdNetwork(ctx, UpdateNetbirdNetworkParams{
		NetworkID:   "net1",
		Name:        &newName,
		Description: &newDescription,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if network == nil {
		t.Fatal("expected network, got nil")
	}
	if network.ID != "net1" {
		t.Errorf("expected ID 'net1', got '%s'", network.ID)
	}
	if network.Name != "Updated Network" {
		t.Errorf("expected Name 'Updated Network', got '%s'", network.Name)
	}
	if network.Description == nil || *network.Description != "Updated description" {
		t.Errorf("expected Description 'Updated description', got %v", network.Description)
	}
}

func TestUpdateNetbirdNetwork_PartialUpdate(t *testing.T) {
	newName := "Partially Updated Network"
	mockResp := NetbirdNetwork{
		ID:                "net1",
		Name:              newName,
		Description:       nil,
		Routers:           []string{},
		RoutingPeersCount: 0,
		Resources:         []string{},
		Policies:          []string{},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1" {
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
	network, err := updateNetbirdNetwork(ctx, UpdateNetbirdNetworkParams{
		NetworkID: "net1",
		Name:      &newName,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if network == nil {
		t.Fatal("expected network, got nil")
	}
	if network.Name != "Partially Updated Network" {
		t.Errorf("expected Name 'Partially Updated Network', got '%s'", network.Name)
	}
}

func TestDeleteNetbirdNetwork(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1" {
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
	result, err := deleteNetbirdNetwork(ctx, DeleteNetbirdNetworkParams{NetworkID: "net1"})
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
}

func TestDeleteNetbirdNetwork_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := deleteNetbirdNetwork(ctx, DeleteNetbirdNetworkParams{NetworkID: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for non-existent network, got nil")
	}
}
