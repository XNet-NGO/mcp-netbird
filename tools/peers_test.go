package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/aantti/mcp-netbird"
)

func TestListNetbirdPeers(t *testing.T) {
	// Mock response data
	mockResp := []NetbirdPeer{
		{
			ID:   "peer1",
			Name: "Test Peer",
			// Add other fields as needed for your struct
		},
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/peers" {
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
	peers, err := listNetbirdPeers(ctx, ListNetbirdPeersParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(peers) != 1 || peers[0].ID != "peer1" {
		t.Errorf("unexpected result: %+v", peers)
	}
}

func TestNetbirdPeerLocalFlagsMarshaling(t *testing.T) {
	// Test marshaling with LocalFlags present - all fields set to true
	peer := NetbirdPeer{
		ID:   "test-peer",
		Name: "Test Peer",
		LocalFlags: &NetbirdPeerLocalFlags{
			RosenpassEnabled:      true,
			RosenpassPermissive:   true,
			ServerSSHAllowed:      true,
			DisableClientRoutes:   true,
			DisableServerRoutes:   true,
			DisableDNS:            true,
			DisableFirewall:       true,
			BlockLANAccess:        true,
			BlockInbound:          true,
			LazyConnectionEnabled: true,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("failed to marshal peer: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdPeer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal peer: %v", err)
	}

	// Verify LocalFlags is present and all fields are correct
	if decoded.LocalFlags == nil {
		t.Fatal("LocalFlags should not be nil")
	}
	if decoded.LocalFlags.RosenpassEnabled != true {
		t.Errorf("expected RosenpassEnabled=true, got %v", decoded.LocalFlags.RosenpassEnabled)
	}
	if decoded.LocalFlags.RosenpassPermissive != true {
		t.Errorf("expected RosenpassPermissive=true, got %v", decoded.LocalFlags.RosenpassPermissive)
	}
	if decoded.LocalFlags.ServerSSHAllowed != true {
		t.Errorf("expected ServerSSHAllowed=true, got %v", decoded.LocalFlags.ServerSSHAllowed)
	}
	if decoded.LocalFlags.DisableClientRoutes != true {
		t.Errorf("expected DisableClientRoutes=true, got %v", decoded.LocalFlags.DisableClientRoutes)
	}
	if decoded.LocalFlags.DisableServerRoutes != true {
		t.Errorf("expected DisableServerRoutes=true, got %v", decoded.LocalFlags.DisableServerRoutes)
	}
	if decoded.LocalFlags.DisableDNS != true {
		t.Errorf("expected DisableDNS=true, got %v", decoded.LocalFlags.DisableDNS)
	}
	if decoded.LocalFlags.DisableFirewall != true {
		t.Errorf("expected DisableFirewall=true, got %v", decoded.LocalFlags.DisableFirewall)
	}
	if decoded.LocalFlags.BlockLANAccess != true {
		t.Errorf("expected BlockLANAccess=true, got %v", decoded.LocalFlags.BlockLANAccess)
	}
	if decoded.LocalFlags.BlockInbound != true {
		t.Errorf("expected BlockInbound=true, got %v", decoded.LocalFlags.BlockInbound)
	}
	if decoded.LocalFlags.LazyConnectionEnabled != true {
		t.Errorf("expected LazyConnectionEnabled=true, got %v", decoded.LocalFlags.LazyConnectionEnabled)
	}

	// Test marshaling without LocalFlags (should be omitted)
	peerWithoutFlags := NetbirdPeer{
		ID:   "test-peer-2",
		Name: "Test Peer 2",
	}

	data2, err := json.Marshal(peerWithoutFlags)
	if err != nil {
		t.Fatalf("failed to marshal peer without flags: %v", err)
	}

	// Verify local_flags is not in the JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data2, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}
	if _, exists := jsonMap["local_flags"]; exists {
		t.Error("local_flags should be omitted when nil")
	}
}

func TestNetbirdPeerNewFieldsMarshaling(t *testing.T) {
	// Test marshaling with all new fields present
	createdAt := "2024-01-15T10:30:00Z"
	ephemeral := true
	disapprovalReason := "Security policy violation"

	peer := NetbirdPeer{
		ID:                "test-peer",
		Name:              "Test Peer",
		CreatedAt:         &createdAt,
		Ephemeral:         &ephemeral,
		DisapprovalReason: &disapprovalReason,
	}

	// Marshal to JSON
	data, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("failed to marshal peer: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdPeer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal peer: %v", err)
	}

	// Verify CreatedAt is present and correct
	if decoded.CreatedAt == nil {
		t.Fatal("CreatedAt should not be nil")
	}
	if *decoded.CreatedAt != createdAt {
		t.Errorf("expected CreatedAt=%s, got %s", createdAt, *decoded.CreatedAt)
	}

	// Verify Ephemeral is present and correct
	if decoded.Ephemeral == nil {
		t.Fatal("Ephemeral should not be nil")
	}
	if *decoded.Ephemeral != ephemeral {
		t.Errorf("expected Ephemeral=%v, got %v", ephemeral, *decoded.Ephemeral)
	}

	// Verify DisapprovalReason is present and correct
	if decoded.DisapprovalReason == nil {
		t.Fatal("DisapprovalReason should not be nil")
	}
	if *decoded.DisapprovalReason != disapprovalReason {
		t.Errorf("expected DisapprovalReason=%s, got %s", disapprovalReason, *decoded.DisapprovalReason)
	}

	// Test marshaling without new fields (should be omitted)
	peerWithoutNewFields := NetbirdPeer{
		ID:   "test-peer-2",
		Name: "Test Peer 2",
	}

	data2, err := json.Marshal(peerWithoutNewFields)
	if err != nil {
		t.Fatalf("failed to marshal peer without new fields: %v", err)
	}

	// Verify new fields are not in the JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data2, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}
	if _, exists := jsonMap["created_at"]; exists {
		t.Error("created_at should be omitted when nil")
	}
	if _, exists := jsonMap["ephemeral"]; exists {
		t.Error("ephemeral should be omitted when nil")
	}
	if _, exists := jsonMap["disapproval_reason"]; exists {
		t.Error("disapproval_reason should be omitted when nil")
	}
}

func TestNetbirdPeerCompleteMarshaling(t *testing.T) {
	// Test marshaling with all fields including new ones and LocalFlags
	createdAt := "2024-01-15T10:30:00Z"
	ephemeral := false
	disapprovalReason := "Pending approval"

	peer := NetbirdPeer{
		ID:                "test-peer",
		Name:              "Test Peer",
		CreatedAt:         &createdAt,
		Ephemeral:         &ephemeral,
		DisapprovalReason: &disapprovalReason,
		LocalFlags: &NetbirdPeerLocalFlags{
			RosenpassEnabled:      true,
			ServerSSHAllowed:      true,
			LazyConnectionEnabled: true,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("failed to marshal peer: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdPeer
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal peer: %v", err)
	}

	// Verify all new fields are present and correct
	if decoded.CreatedAt == nil || *decoded.CreatedAt != createdAt {
		t.Errorf("CreatedAt mismatch: expected %s, got %v", createdAt, decoded.CreatedAt)
	}
	if decoded.Ephemeral == nil || *decoded.Ephemeral != ephemeral {
		t.Errorf("Ephemeral mismatch: expected %v, got %v", ephemeral, decoded.Ephemeral)
	}
	if decoded.DisapprovalReason == nil || *decoded.DisapprovalReason != disapprovalReason {
		t.Errorf("DisapprovalReason mismatch: expected %s, got %v", disapprovalReason, decoded.DisapprovalReason)
	}
	if decoded.LocalFlags == nil {
		t.Fatal("LocalFlags should not be nil")
	}
	if !decoded.LocalFlags.RosenpassEnabled {
		t.Error("expected RosenpassEnabled=true")
	}
}

func TestNetbirdPeerJSONFieldNames(t *testing.T) {
	// Test that JSON field names are correct (snake_case)
	createdAt := "2024-01-15T10:30:00Z"
	ephemeral := true
	disapprovalReason := "Test reason"

	peer := NetbirdPeer{
		ID:                "test-peer",
		Name:              "Test Peer",
		CreatedAt:         &createdAt,
		Ephemeral:         &ephemeral,
		DisapprovalReason: &disapprovalReason,
		LocalFlags: &NetbirdPeerLocalFlags{
			RosenpassEnabled:      true,
			RosenpassPermissive:   false,
			ServerSSHAllowed:      true,
			DisableClientRoutes:   false,
			DisableServerRoutes:   true,
			DisableDNS:            false,
			DisableFirewall:       true,
			BlockLANAccess:        false,
			BlockInbound:          true,
			LazyConnectionEnabled: false,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("failed to marshal peer: %v", err)
	}

	// Unmarshal to map to check field names
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Verify new field names are in snake_case
	if _, exists := jsonMap["created_at"]; !exists {
		t.Error("expected 'created_at' field in JSON")
	}
	if _, exists := jsonMap["ephemeral"]; !exists {
		t.Error("expected 'ephemeral' field in JSON")
	}
	if _, exists := jsonMap["disapproval_reason"]; !exists {
		t.Error("expected 'disapproval_reason' field in JSON")
	}
	if _, exists := jsonMap["local_flags"]; !exists {
		t.Error("expected 'local_flags' field in JSON")
	}

	// Verify LocalFlags nested field names
	localFlags, ok := jsonMap["local_flags"].(map[string]interface{})
	if !ok {
		t.Fatal("local_flags should be a map")
	}

	expectedLocalFlagsFields := []string{
		"rosenpass_enabled",
		"rosenpass_permissive",
		"server_ssh_allowed",
		"disable_client_routes",
		"disable_server_routes",
		"disable_dns",
		"disable_firewall",
		"block_lan_access",
		"block_inbound",
		"lazy_connection_enabled",
	}

	for _, field := range expectedLocalFlagsFields {
		if _, exists := localFlags[field]; !exists {
			t.Errorf("expected '%s' field in local_flags", field)
		}
	}
}

func TestNetbirdPeerLocalFlagsAllFieldsRoundTrip(t *testing.T) {
	// Test that all LocalFlags fields survive a round-trip through JSON
	original := NetbirdPeerLocalFlags{
		RosenpassEnabled:      true,
		RosenpassPermissive:   false,
		ServerSSHAllowed:      true,
		DisableClientRoutes:   false,
		DisableServerRoutes:   true,
		DisableDNS:            false,
		DisableFirewall:       true,
		BlockLANAccess:        false,
		BlockInbound:          true,
		LazyConnectionEnabled: false,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal LocalFlags: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdPeerLocalFlags
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal LocalFlags: %v", err)
	}

	// Verify all fields match
	if decoded.RosenpassEnabled != original.RosenpassEnabled {
		t.Errorf("RosenpassEnabled mismatch: expected %v, got %v", original.RosenpassEnabled, decoded.RosenpassEnabled)
	}
	if decoded.RosenpassPermissive != original.RosenpassPermissive {
		t.Errorf("RosenpassPermissive mismatch: expected %v, got %v", original.RosenpassPermissive, decoded.RosenpassPermissive)
	}
	if decoded.ServerSSHAllowed != original.ServerSSHAllowed {
		t.Errorf("ServerSSHAllowed mismatch: expected %v, got %v", original.ServerSSHAllowed, decoded.ServerSSHAllowed)
	}
	if decoded.DisableClientRoutes != original.DisableClientRoutes {
		t.Errorf("DisableClientRoutes mismatch: expected %v, got %v", original.DisableClientRoutes, decoded.DisableClientRoutes)
	}
	if decoded.DisableServerRoutes != original.DisableServerRoutes {
		t.Errorf("DisableServerRoutes mismatch: expected %v, got %v", original.DisableServerRoutes, decoded.DisableServerRoutes)
	}
	if decoded.DisableDNS != original.DisableDNS {
		t.Errorf("DisableDNS mismatch: expected %v, got %v", original.DisableDNS, decoded.DisableDNS)
	}
	if decoded.DisableFirewall != original.DisableFirewall {
		t.Errorf("DisableFirewall mismatch: expected %v, got %v", original.DisableFirewall, decoded.DisableFirewall)
	}
	if decoded.BlockLANAccess != original.BlockLANAccess {
		t.Errorf("BlockLANAccess mismatch: expected %v, got %v", original.BlockLANAccess, decoded.BlockLANAccess)
	}
	if decoded.BlockInbound != original.BlockInbound {
		t.Errorf("BlockInbound mismatch: expected %v, got %v", original.BlockInbound, decoded.BlockInbound)
	}
	if decoded.LazyConnectionEnabled != original.LazyConnectionEnabled {
		t.Errorf("LazyConnectionEnabled mismatch: expected %v, got %v", original.LazyConnectionEnabled, decoded.LazyConnectionEnabled)
	}
}
