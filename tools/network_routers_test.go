package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/aantti/mcp-netbird"
)

func TestListNetbirdNetworkRouters(t *testing.T) {
	peer := "peer1"
	peerGroups := []string{"grp1", "grp2"}
	// Mock response data
	mockResp := []NetbirdNetworkRouter{
		{
			ID:         "rtr1",
			Peer:       &peer,
			PeerGroups: nil,
			Metric:     100,
			Masquerade: true,
			Enabled:    true,
		},
		{
			ID:         "rtr2",
			Peer:       nil,
			PeerGroups: &peerGroups,
			Metric:     200,
			Masquerade: false,
			Enabled:    false,
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers" {
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
	routers, err := listNetbirdNetworkRouters(ctx, ListNetbirdNetworkRoutersParams{NetworkID: "net1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(routers) != 2 {
		t.Fatalf("expected 2 routers, got %d", len(routers))
	}
	if routers[0].ID != "rtr1" {
		t.Errorf("expected ID 'rtr1', got '%s'", routers[0].ID)
	}
	if routers[0].Peer == nil || *routers[0].Peer != "peer1" {
		t.Errorf("expected Peer 'peer1', got %v", routers[0].Peer)
	}
	if routers[0].PeerGroups != nil {
		t.Errorf("expected PeerGroups to be nil, got %v", routers[0].PeerGroups)
	}
	if routers[0].Metric != 100 {
		t.Errorf("expected Metric 100, got %d", routers[0].Metric)
	}
	if !routers[0].Masquerade {
		t.Error("expected Masquerade to be true")
	}
	if !routers[0].Enabled {
		t.Error("expected Enabled to be true")
	}

	// Check second router
	if routers[1].ID != "rtr2" {
		t.Errorf("expected ID 'rtr2', got '%s'", routers[1].ID)
	}
	if routers[1].Peer != nil {
		t.Errorf("expected Peer to be nil, got %v", routers[1].Peer)
	}
	if routers[1].PeerGroups == nil || len(*routers[1].PeerGroups) != 2 {
		t.Errorf("expected PeerGroups with 2 items, got %v", routers[1].PeerGroups)
	}
	if routers[1].Metric != 200 {
		t.Errorf("expected Metric 200, got %d", routers[1].Metric)
	}
	if routers[1].Masquerade {
		t.Error("expected Masquerade to be false")
	}
	if routers[1].Enabled {
		t.Error("expected Enabled to be false")
	}
}

func TestGetNetbirdNetworkRouter(t *testing.T) {
	peer := "peer1"
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr1",
		Peer:       &peer,
		PeerGroups: nil,
		Metric:     150,
		Masquerade: true,
		Enabled:    true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers/rtr1" {
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
	router, err := getNetbirdNetworkRouter(ctx, GetNetbirdNetworkRouterParams{
		NetworkID: "net1",
		RouterID:  "rtr1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.ID != "rtr1" {
		t.Errorf("expected ID 'rtr1', got '%s'", router.ID)
	}
	if router.Peer == nil || *router.Peer != "peer1" {
		t.Errorf("expected Peer 'peer1', got %v", router.Peer)
	}
	if router.PeerGroups != nil {
		t.Errorf("expected PeerGroups to be nil, got %v", router.PeerGroups)
	}
	if router.Metric != 150 {
		t.Errorf("expected Metric 150, got %d", router.Metric)
	}
	if !router.Masquerade {
		t.Error("expected Masquerade to be true")
	}
	if !router.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestGetNetbirdNetworkRouter_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := getNetbirdNetworkRouter(ctx, GetNetbirdNetworkRouterParams{
		NetworkID: "net1",
		RouterID:  "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for non-existent router, got nil")
	}
}

func TestCreateNetbirdNetworkRouter_WithPeer(t *testing.T) {
	peer := "peer1"
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr1",
		Peer:       &peer,
		PeerGroups: nil,
		Metric:     100,
		Masquerade: true,
		Enabled:    true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers" {
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
		if req["peer"] != "peer1" {
			t.Errorf("expected Peer 'peer1', got '%v'", req["peer"])
		}
		if _, exists := req["peer_groups"]; exists {
			t.Error("expected peer_groups to be omitted when using peer")
		}
		if req["metric"] != float64(100) {
			t.Errorf("expected Metric 100, got %v", req["metric"])
		}
		if req["masquerade"] != true {
			t.Errorf("expected Masquerade true, got %v", req["masquerade"])
		}
		if req["enabled"] != true {
			t.Errorf("expected Enabled true, got %v", req["enabled"])
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
	router, err := createNetbirdNetworkRouter(ctx, CreateNetbirdNetworkRouterParams{
		NetworkID:  "net1",
		Peer:       &peer,
		PeerGroups: nil,
		Metric:     100,
		Masquerade: true,
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.ID != "rtr1" {
		t.Errorf("expected ID 'rtr1', got '%s'", router.ID)
	}
	if router.Peer == nil || *router.Peer != "peer1" {
		t.Errorf("expected Peer 'peer1', got %v", router.Peer)
	}
	if router.PeerGroups != nil {
		t.Errorf("expected PeerGroups to be nil, got %v", router.PeerGroups)
	}
}

func TestCreateNetbirdNetworkRouter_WithPeerGroups(t *testing.T) {
	peerGroups := []string{"grp1", "grp2"}
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr2",
		Peer:       nil,
		PeerGroups: &peerGroups,
		Metric:     200,
		Masquerade: false,
		Enabled:    false,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers" {
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
		if _, exists := req["peer"]; exists {
			t.Error("expected peer to be omitted when using peer_groups")
		}
		groups, ok := req["peer_groups"].([]interface{})
		if !ok || len(groups) != 2 {
			t.Errorf("expected peer_groups with 2 items, got %v", req["peer_groups"])
		}
		if req["metric"] != float64(200) {
			t.Errorf("expected Metric 200, got %v", req["metric"])
		}
		if req["masquerade"] != false {
			t.Errorf("expected Masquerade false, got %v", req["masquerade"])
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
	router, err := createNetbirdNetworkRouter(ctx, CreateNetbirdNetworkRouterParams{
		NetworkID:  "net1",
		Peer:       nil,
		PeerGroups: &peerGroups,
		Metric:     200,
		Masquerade: false,
		Enabled:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.ID != "rtr2" {
		t.Errorf("expected ID 'rtr2', got '%s'", router.ID)
	}
	if router.Peer != nil {
		t.Errorf("expected Peer to be nil, got %v", router.Peer)
	}
	if router.PeerGroups == nil || len(*router.PeerGroups) != 2 {
		t.Errorf("expected PeerGroups with 2 items, got %v", router.PeerGroups)
	}
}

func TestCreateNetbirdNetworkRouter_WithoutPeerOrPeerGroups(t *testing.T) {
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr3",
		Peer:       nil,
		PeerGroups: nil,
		Metric:     300,
		Masquerade: true,
		Enabled:    true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers" {
			http.NotFound(w, r)
			return
		}

		// Verify request body - neither peer nor peer_groups should be present
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if _, exists := req["peer"]; exists {
			t.Error("expected peer to be omitted when nil")
		}
		if _, exists := req["peer_groups"]; exists {
			t.Error("expected peer_groups to be omitted when nil")
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
	router, err := createNetbirdNetworkRouter(ctx, CreateNetbirdNetworkRouterParams{
		NetworkID:  "net1",
		Peer:       nil,
		PeerGroups: nil,
		Metric:     300,
		Masquerade: true,
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.ID != "rtr3" {
		t.Errorf("expected ID 'rtr3', got '%s'", router.ID)
	}
}

func TestUpdateNetbirdNetworkRouter(t *testing.T) {
	newPeer := "peer2"
	newMetric := 250
	newMasquerade := false
	newEnabled := false
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr1",
		Peer:       &newPeer,
		PeerGroups: nil,
		Metric:     newMetric,
		Masquerade: newMasquerade,
		Enabled:    newEnabled,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers/rtr1" {
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
		if req["peer"] != "peer2" {
			t.Errorf("expected Peer 'peer2', got '%v'", req["peer"])
		}
		if req["metric"] != float64(250) {
			t.Errorf("expected Metric 250, got %v", req["metric"])
		}
		if req["masquerade"] != false {
			t.Errorf("expected Masquerade false, got %v", req["masquerade"])
		}
		if req["enabled"] != false {
			t.Errorf("expected Enabled false, got %v", req["enabled"])
		}
		// Verify network_id and router_id are NOT in the request body
		if _, exists := req["network_id"]; exists {
			t.Error("network_id should not be in request body")
		}
		if _, exists := req["router_id"]; exists {
			t.Error("router_id should not be in request body")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	router, err := updateNetbirdNetworkRouter(ctx, UpdateNetbirdNetworkRouterParams{
		NetworkID:  "net1",
		RouterID:   "rtr1",
		Peer:       &newPeer,
		Metric:     &newMetric,
		Masquerade: &newMasquerade,
		Enabled:    &newEnabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.ID != "rtr1" {
		t.Errorf("expected ID 'rtr1', got '%s'", router.ID)
	}
	if router.Peer == nil || *router.Peer != "peer2" {
		t.Errorf("expected Peer 'peer2', got %v", router.Peer)
	}
	if router.Metric != 250 {
		t.Errorf("expected Metric 250, got %d", router.Metric)
	}
	if router.Masquerade {
		t.Error("expected Masquerade to be false")
	}
	if router.Enabled {
		t.Error("expected Enabled to be false")
	}
}

func TestUpdateNetbirdNetworkRouter_PartialUpdate(t *testing.T) {
	newMetric := 350
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr1",
		Peer:       nil,
		PeerGroups: nil,
		Metric:     newMetric,
		Masquerade: true,
		Enabled:    true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers/rtr1" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify request body - only metric should be present
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if _, exists := req["peer"]; exists {
			t.Error("expected peer to be omitted when nil")
		}
		if _, exists := req["peer_groups"]; exists {
			t.Error("expected peer_groups to be omitted when nil")
		}
		if _, exists := req["masquerade"]; exists {
			t.Error("expected masquerade to be omitted when nil")
		}
		if _, exists := req["enabled"]; exists {
			t.Error("expected enabled to be omitted when nil")
		}
		if req["metric"] == nil {
			t.Error("expected metric to be present")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	router, err := updateNetbirdNetworkRouter(ctx, UpdateNetbirdNetworkRouterParams{
		NetworkID: "net1",
		RouterID:  "rtr1",
		Metric:    &newMetric,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.Metric != 350 {
		t.Errorf("expected Metric 350, got %d", router.Metric)
	}
}

func TestUpdateNetbirdNetworkRouter_SwitchFromPeerToPeerGroups(t *testing.T) {
	peerGroups := []string{"grp1", "grp2"}
	mockResp := NetbirdNetworkRouter{
		ID:         "rtr1",
		Peer:       nil,
		PeerGroups: &peerGroups,
		Metric:     100,
		Masquerade: true,
		Enabled:    true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers/rtr1" {
			http.NotFound(w, r)
			return
		}

		// Verify request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if _, exists := req["peer"]; exists {
			t.Error("expected peer to be omitted when switching to peer_groups")
		}
		groups, ok := req["peer_groups"].([]interface{})
		if !ok || len(groups) != 2 {
			t.Errorf("expected peer_groups with 2 items, got %v", req["peer_groups"])
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	router, err := updateNetbirdNetworkRouter(ctx, UpdateNetbirdNetworkRouterParams{
		NetworkID:  "net1",
		RouterID:   "rtr1",
		PeerGroups: &peerGroups,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected router, got nil")
	}
	if router.Peer != nil {
		t.Errorf("expected Peer to be nil, got %v", router.Peer)
	}
	if router.PeerGroups == nil || len(*router.PeerGroups) != 2 {
		t.Errorf("expected PeerGroups with 2 items, got %v", router.PeerGroups)
	}
}

func TestDeleteNetbirdNetworkRouter(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/networks/net1/routers/rtr1" {
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
	result, err := deleteNetbirdNetworkRouter(ctx, DeleteNetbirdNetworkRouterParams{
		NetworkID: "net1",
		RouterID:  "rtr1",
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
	if result["router_id"] != "rtr1" {
		t.Errorf("expected router_id 'rtr1', got '%s'", result["router_id"])
	}
}

func TestDeleteNetbirdNetworkRouter_NotFound(t *testing.T) {
	// Create mock HTTP server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set the test client
	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	_, err := deleteNetbirdNetworkRouter(ctx, DeleteNetbirdNetworkRouterParams{
		NetworkID: "net1",
		RouterID:  "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for non-existent router, got nil")
	}
}
