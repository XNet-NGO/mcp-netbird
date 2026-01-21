package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
)

// TestAccountSettingsMarshaling tests JSON marshaling/unmarshaling of NetbirdAccountSettings
func TestAccountSettingsMarshaling(t *testing.T) {
	// Test with all new fields populated
	trueVal := true
	falseVal := false
	dnsResolutionEnabled := true
	dnsDomain := "example.com"
	networkRange := "10.0.0.0/8"
	lazyConnection := false
	autoUpdateVersion := "0.28.0"
	embeddedIDP := true

	settings := NetbirdAccountSettings{
		PeerLoginExpiration:             3600,
		PeerLoginExpirationEnabled:      true,
		PeerInactivityExpiration:        7200,
		PeerInactivityExpirationEnabled: true,
		GroupsPropagationEnabled:        true,
		JWTGroupsEnabled:                false,
		JWTGroupsClaimName:              "groups",
		JWTAllowGroups:                  []string{"admin", "users"},
		RegularUsersViewBlocked:         true,
		RoutingPeerDNSResolutionEnabled: &dnsResolutionEnabled,
		DNSDomain:                       &dnsDomain,
		NetworkRange:                    &networkRange,
		LazyConnectionEnabled:           &lazyConnection,
		AutoUpdateVersion:               &autoUpdateVersion,
		EmbeddedIDPEnabled:              &embeddedIDP,
	}

	// Marshal to JSON
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountSettings
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal settings: %v", err)
	}

	// Verify all fields are preserved
	if decoded.PeerLoginExpiration != settings.PeerLoginExpiration {
		t.Errorf("PeerLoginExpiration mismatch: got %d, want %d", decoded.PeerLoginExpiration, settings.PeerLoginExpiration)
	}
	if decoded.RegularUsersViewBlocked != settings.RegularUsersViewBlocked {
		t.Errorf("RegularUsersViewBlocked mismatch: got %v, want %v", decoded.RegularUsersViewBlocked, settings.RegularUsersViewBlocked)
	}
	if decoded.RoutingPeerDNSResolutionEnabled == nil || *decoded.RoutingPeerDNSResolutionEnabled != *settings.RoutingPeerDNSResolutionEnabled {
		t.Errorf("RoutingPeerDNSResolutionEnabled mismatch: got %v, want %v", decoded.RoutingPeerDNSResolutionEnabled, settings.RoutingPeerDNSResolutionEnabled)
	}
	if decoded.DNSDomain == nil || *decoded.DNSDomain != *settings.DNSDomain {
		t.Errorf("DNSDomain mismatch: got %v, want %v", decoded.DNSDomain, settings.DNSDomain)
	}
	if decoded.NetworkRange == nil || *decoded.NetworkRange != *settings.NetworkRange {
		t.Errorf("NetworkRange mismatch: got %v, want %v", decoded.NetworkRange, settings.NetworkRange)
	}
	if decoded.LazyConnectionEnabled == nil || *decoded.LazyConnectionEnabled != *settings.LazyConnectionEnabled {
		t.Errorf("LazyConnectionEnabled mismatch: got %v, want %v", decoded.LazyConnectionEnabled, settings.LazyConnectionEnabled)
	}
	if decoded.AutoUpdateVersion == nil || *decoded.AutoUpdateVersion != *settings.AutoUpdateVersion {
		t.Errorf("AutoUpdateVersion mismatch: got %v, want %v", decoded.AutoUpdateVersion, settings.AutoUpdateVersion)
	}
	if decoded.EmbeddedIDPEnabled == nil || *decoded.EmbeddedIDPEnabled != *settings.EmbeddedIDPEnabled {
		t.Errorf("EmbeddedIDPEnabled mismatch: got %v, want %v", decoded.EmbeddedIDPEnabled, settings.EmbeddedIDPEnabled)
	}

	// Verify pointer fields are present in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"routing_peer_dns_resolution_enabled",
		"dns_domain",
		"network_range",
		"lazy_connection_enabled",
		"auto_update_version",
		"embedded_idp_enabled",
	}

	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %s to be present in JSON", field)
		}
	}

	// Test with boolean values
	settings.RegularUsersViewBlocked = false
	settings.RoutingPeerDNSResolutionEnabled = &falseVal
	settings.LazyConnectionEnabled = &trueVal
	settings.EmbeddedIDPEnabled = &falseVal

	data, err = json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal settings with false values: %v", err)
	}

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal settings with false values: %v", err)
	}

	if decoded.RegularUsersViewBlocked != false {
		t.Errorf("RegularUsersViewBlocked should be false, got %v", decoded.RegularUsersViewBlocked)
	}
	if decoded.RoutingPeerDNSResolutionEnabled == nil || *decoded.RoutingPeerDNSResolutionEnabled != false {
		t.Errorf("RoutingPeerDNSResolutionEnabled should be false, got %v", decoded.RoutingPeerDNSResolutionEnabled)
	}
	if decoded.LazyConnectionEnabled == nil || *decoded.LazyConnectionEnabled != true {
		t.Errorf("LazyConnectionEnabled should be true, got %v", decoded.LazyConnectionEnabled)
	}
	if decoded.EmbeddedIDPEnabled == nil || *decoded.EmbeddedIDPEnabled != false {
		t.Errorf("EmbeddedIDPEnabled should be false, got %v", decoded.EmbeddedIDPEnabled)
	}
}

// TestAccountSettingsOptionalFieldsOmitted tests that optional fields are omitted when nil
func TestAccountSettingsOptionalFieldsOmitted(t *testing.T) {
	settings := NetbirdAccountSettings{
		PeerLoginExpiration:             3600,
		PeerLoginExpirationEnabled:      true,
		PeerInactivityExpiration:        7200,
		PeerInactivityExpirationEnabled: true,
		GroupsPropagationEnabled:        true,
		JWTGroupsEnabled:                false,
		JWTGroupsClaimName:              "groups",
		JWTAllowGroups:                  []string{"admin"},
		RegularUsersViewBlocked:         false,
		// All pointer fields are nil
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	// Verify optional fields are NOT present in JSON when nil
	optionalFields := []string{
		"routing_peer_dns_resolution_enabled",
		"dns_domain",
		"network_range",
		"lazy_connection_enabled",
		"auto_update_version",
		"embedded_idp_enabled",
	}

	for _, field := range optionalFields {
		if _, exists := jsonMap[field]; exists {
			t.Errorf("field %s should be omitted when nil, but was present", field)
		}
	}

	// Verify required field is present
	if _, exists := jsonMap["regular_users_view_blocked"]; !exists {
		t.Errorf("required field regular_users_view_blocked should be present")
	}
}

// TestGetNetbirdAccount tests the getNetbirdAccount function
func TestGetNetbirdAccount(t *testing.T) {
	dnsDomain := "test.example.com"
	networkRange := "10.0.0.0/16"
	lazyConnection := true

	mockResp := []NetbirdAccount{
		{
			ID: "account1",
			Settings: NetbirdAccountSettings{
				PeerLoginExpiration:             3600,
				PeerLoginExpirationEnabled:      true,
				PeerInactivityExpiration:        7200,
				PeerInactivityExpirationEnabled: false,
				GroupsPropagationEnabled:        true,
				JWTGroupsEnabled:                true,
				JWTGroupsClaimName:              "groups",
				JWTAllowGroups:                  []string{"admin", "users"},
				RegularUsersViewBlocked:         true,
				DNSDomain:                       &dnsDomain,
				NetworkRange:                    &networkRange,
				LazyConnectionEnabled:           &lazyConnection,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	account, err := getNetbirdAccount(ctx, GetNetbirdAccountParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.ID != "account1" {
		t.Errorf("unexpected account ID: got %s, want account1", account.ID)
	}
	if account.Settings.RegularUsersViewBlocked != true {
		t.Errorf("RegularUsersViewBlocked should be true, got %v", account.Settings.RegularUsersViewBlocked)
	}
	if account.Settings.DNSDomain == nil || *account.Settings.DNSDomain != dnsDomain {
		t.Errorf("DNSDomain mismatch: got %v, want %s", account.Settings.DNSDomain, dnsDomain)
	}
	if account.Settings.NetworkRange == nil || *account.Settings.NetworkRange != networkRange {
		t.Errorf("NetworkRange mismatch: got %v, want %s", account.Settings.NetworkRange, networkRange)
	}
	if account.Settings.LazyConnectionEnabled == nil || *account.Settings.LazyConnectionEnabled != lazyConnection {
		t.Errorf("LazyConnectionEnabled mismatch: got %v, want %v", account.Settings.LazyConnectionEnabled, lazyConnection)
	}
}

// TestUpdateNetbirdAccount tests the updateNetbirdAccount function
func TestUpdateNetbirdAccount(t *testing.T) {
	autoUpdateVersion := "0.28.0"
	embeddedIDP := false

	mockResp := []NetbirdAccount{
		{
			ID: "account1",
			Settings: NetbirdAccountSettings{
				PeerLoginExpiration:             7200,
				PeerLoginExpirationEnabled:      false,
				PeerInactivityExpiration:        14400,
				PeerInactivityExpirationEnabled: true,
				GroupsPropagationEnabled:        false,
				JWTGroupsEnabled:                false,
				JWTGroupsClaimName:              "roles",
				JWTAllowGroups:                  []string{"developers"},
				RegularUsersViewBlocked:         false,
				AutoUpdateVersion:               &autoUpdateVersion,
				EmbeddedIDPEnabled:              &embeddedIDP,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	params := UpdateNetbirdAccountParams{
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             7200,
			PeerLoginExpirationEnabled:      false,
			PeerInactivityExpiration:        14400,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        false,
			JWTGroupsEnabled:                false,
			JWTGroupsClaimName:              "roles",
			JWTAllowGroups:                  []string{"developers"},
			RegularUsersViewBlocked:         false,
			AutoUpdateVersion:               &autoUpdateVersion,
			EmbeddedIDPEnabled:              &embeddedIDP,
		},
	}

	account, err := updateNetbirdAccount(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.Settings.AutoUpdateVersion == nil || *account.Settings.AutoUpdateVersion != autoUpdateVersion {
		t.Errorf("AutoUpdateVersion mismatch: got %v, want %s", account.Settings.AutoUpdateVersion, autoUpdateVersion)
	}
	if account.Settings.EmbeddedIDPEnabled == nil || *account.Settings.EmbeddedIDPEnabled != embeddedIDP {
		t.Errorf("EmbeddedIDPEnabled mismatch: got %v, want %v", account.Settings.EmbeddedIDPEnabled, embeddedIDP)
	}
}

// TestAccountExtraMarshaling tests JSON marshaling/unmarshaling of NetbirdAccountExtra (Requirement 1.5)
func TestAccountExtraMarshaling(t *testing.T) {
	extra := NetbirdAccountExtra{
		PeerApprovalEnabled:                true,
		UserApprovalRequired:               false,
		NetworkTrafficLogsEnabled:          true,
		NetworkTrafficLogsGroups:           []string{"group1", "group2"},
		NetworkTrafficPacketCounterEnabled: false,
	}

	// Marshal to JSON
	data, err := json.Marshal(extra)
	if err != nil {
		t.Fatalf("failed to marshal extra: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountExtra
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal extra: %v", err)
	}

	// Verify all fields are preserved
	if decoded.PeerApprovalEnabled != extra.PeerApprovalEnabled {
		t.Errorf("PeerApprovalEnabled mismatch: got %v, want %v", decoded.PeerApprovalEnabled, extra.PeerApprovalEnabled)
	}
	if decoded.UserApprovalRequired != extra.UserApprovalRequired {
		t.Errorf("UserApprovalRequired mismatch: got %v, want %v", decoded.UserApprovalRequired, extra.UserApprovalRequired)
	}
	if decoded.NetworkTrafficLogsEnabled != extra.NetworkTrafficLogsEnabled {
		t.Errorf("NetworkTrafficLogsEnabled mismatch: got %v, want %v", decoded.NetworkTrafficLogsEnabled, extra.NetworkTrafficLogsEnabled)
	}
	if len(decoded.NetworkTrafficLogsGroups) != len(extra.NetworkTrafficLogsGroups) {
		t.Errorf("NetworkTrafficLogsGroups length mismatch: got %d, want %d", len(decoded.NetworkTrafficLogsGroups), len(extra.NetworkTrafficLogsGroups))
	}
	for i, group := range extra.NetworkTrafficLogsGroups {
		if decoded.NetworkTrafficLogsGroups[i] != group {
			t.Errorf("NetworkTrafficLogsGroups[%d] mismatch: got %s, want %s", i, decoded.NetworkTrafficLogsGroups[i], group)
		}
	}
	if decoded.NetworkTrafficPacketCounterEnabled != extra.NetworkTrafficPacketCounterEnabled {
		t.Errorf("NetworkTrafficPacketCounterEnabled mismatch: got %v, want %v", decoded.NetworkTrafficPacketCounterEnabled, extra.NetworkTrafficPacketCounterEnabled)
	}

	// Verify JSON field names
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"peer_approval_enabled",
		"user_approval_required",
		"network_traffic_logs_enabled",
		"network_traffic_logs_groups",
		"network_traffic_packet_counter_enabled",
	}

	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %s to be present in JSON", field)
		}
	}
}

// TestAccountSettingsWithExtra tests NetbirdAccountSettings with Extra field (Requirement 1.5)
func TestAccountSettingsWithExtra(t *testing.T) {
	extra := &NetbirdAccountExtra{
		PeerApprovalEnabled:                true,
		UserApprovalRequired:               true,
		NetworkTrafficLogsEnabled:          false,
		NetworkTrafficLogsGroups:           []string{"admin", "monitoring"},
		NetworkTrafficPacketCounterEnabled: true,
	}

	settings := NetbirdAccountSettings{
		PeerLoginExpiration:             3600,
		PeerLoginExpirationEnabled:      true,
		PeerInactivityExpiration:        7200,
		PeerInactivityExpirationEnabled: true,
		GroupsPropagationEnabled:        true,
		JWTGroupsEnabled:                false,
		JWTGroupsClaimName:              "groups",
		JWTAllowGroups:                  []string{"admin"},
		RegularUsersViewBlocked:         false,
		Extra:                           extra,
	}

	// Marshal to JSON
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountSettings
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal settings: %v", err)
	}

	// Verify Extra field is preserved
	if decoded.Extra == nil {
		t.Fatal("Extra field should not be nil")
	}
	if decoded.Extra.PeerApprovalEnabled != extra.PeerApprovalEnabled {
		t.Errorf("Extra.PeerApprovalEnabled mismatch: got %v, want %v", decoded.Extra.PeerApprovalEnabled, extra.PeerApprovalEnabled)
	}
	if decoded.Extra.UserApprovalRequired != extra.UserApprovalRequired {
		t.Errorf("Extra.UserApprovalRequired mismatch: got %v, want %v", decoded.Extra.UserApprovalRequired, extra.UserApprovalRequired)
	}
	if decoded.Extra.NetworkTrafficLogsEnabled != extra.NetworkTrafficLogsEnabled {
		t.Errorf("Extra.NetworkTrafficLogsEnabled mismatch: got %v, want %v", decoded.Extra.NetworkTrafficLogsEnabled, extra.NetworkTrafficLogsEnabled)
	}
	if len(decoded.Extra.NetworkTrafficLogsGroups) != len(extra.NetworkTrafficLogsGroups) {
		t.Errorf("Extra.NetworkTrafficLogsGroups length mismatch: got %d, want %d", len(decoded.Extra.NetworkTrafficLogsGroups), len(extra.NetworkTrafficLogsGroups))
	}
	if decoded.Extra.NetworkTrafficPacketCounterEnabled != extra.NetworkTrafficPacketCounterEnabled {
		t.Errorf("Extra.NetworkTrafficPacketCounterEnabled mismatch: got %v, want %v", decoded.Extra.NetworkTrafficPacketCounterEnabled, extra.NetworkTrafficPacketCounterEnabled)
	}

	// Verify Extra field is present in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["extra"]; !exists {
		t.Error("extra field should be present in JSON")
	}
}

// TestAccountSettingsExtraOmitted tests that Extra field is omitted when nil (Requirement 1.5)
func TestAccountSettingsExtraOmitted(t *testing.T) {
	settings := NetbirdAccountSettings{
		PeerLoginExpiration:             3600,
		PeerLoginExpirationEnabled:      true,
		PeerInactivityExpiration:        7200,
		PeerInactivityExpirationEnabled: true,
		GroupsPropagationEnabled:        true,
		JWTGroupsEnabled:                false,
		JWTGroupsClaimName:              "groups",
		JWTAllowGroups:                  []string{"admin"},
		RegularUsersViewBlocked:         false,
		Extra:                           nil, // Explicitly nil
	}

	// Marshal to JSON
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal settings: %v", err)
	}

	// Verify Extra field is NOT present in JSON when nil
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["extra"]; exists {
		t.Error("extra field should be omitted when nil")
	}
}

// TestAccountExtraEmptyGroups tests NetbirdAccountExtra with empty NetworkTrafficLogsGroups (Requirement 1.5)
func TestAccountExtraEmptyGroups(t *testing.T) {
	extra := NetbirdAccountExtra{
		PeerApprovalEnabled:                false,
		UserApprovalRequired:               false,
		NetworkTrafficLogsEnabled:          false,
		NetworkTrafficLogsGroups:           []string{}, // Empty slice
		NetworkTrafficPacketCounterEnabled: false,
	}

	// Marshal to JSON
	data, err := json.Marshal(extra)
	if err != nil {
		t.Fatalf("failed to marshal extra: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountExtra
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal extra: %v", err)
	}

	// Verify empty slice is preserved
	if decoded.NetworkTrafficLogsGroups == nil {
		t.Error("NetworkTrafficLogsGroups should not be nil, should be empty slice")
	}
	if len(decoded.NetworkTrafficLogsGroups) != 0 {
		t.Errorf("NetworkTrafficLogsGroups should be empty, got length %d", len(decoded.NetworkTrafficLogsGroups))
	}
}


// TestAccountOnboardingMarshaling tests JSON marshaling/unmarshaling of NetbirdAccountOnboarding (Requirement 1.9)
func TestAccountOnboardingMarshaling(t *testing.T) {
	onboarding := NetbirdAccountOnboarding{
		SignupFormPending:     true,
		OnboardingFlowPending: false,
	}

	// Marshal to JSON
	data, err := json.Marshal(onboarding)
	if err != nil {
		t.Fatalf("failed to marshal onboarding: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountOnboarding
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal onboarding: %v", err)
	}

	// Verify all fields are preserved
	if decoded.SignupFormPending != onboarding.SignupFormPending {
		t.Errorf("SignupFormPending mismatch: got %v, want %v", decoded.SignupFormPending, onboarding.SignupFormPending)
	}
	if decoded.OnboardingFlowPending != onboarding.OnboardingFlowPending {
		t.Errorf("OnboardingFlowPending mismatch: got %v, want %v", decoded.OnboardingFlowPending, onboarding.OnboardingFlowPending)
	}

	// Verify JSON field names
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"signup_form_pending",
		"onboarding_flow_pending",
	}

	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %s to be present in JSON", field)
		}
	}
}

// TestAccountWithOnboarding tests NetbirdAccount with Onboarding field (Requirement 1.9)
func TestAccountWithOnboarding(t *testing.T) {
	onboarding := &NetbirdAccountOnboarding{
		SignupFormPending:     false,
		OnboardingFlowPending: true,
	}

	account := NetbirdAccount{
		ID: "account123",
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             3600,
			PeerLoginExpirationEnabled:      true,
			PeerInactivityExpiration:        7200,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        true,
			JWTGroupsEnabled:                false,
			JWTGroupsClaimName:              "groups",
			JWTAllowGroups:                  []string{"admin"},
			RegularUsersViewBlocked:         false,
		},
		Onboarding: onboarding,
	}

	// Marshal to JSON
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal account: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccount
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal account: %v", err)
	}

	// Verify Onboarding field is preserved
	if decoded.Onboarding == nil {
		t.Fatal("Onboarding field should not be nil")
	}
	if decoded.Onboarding.SignupFormPending != onboarding.SignupFormPending {
		t.Errorf("Onboarding.SignupFormPending mismatch: got %v, want %v", decoded.Onboarding.SignupFormPending, onboarding.SignupFormPending)
	}
	if decoded.Onboarding.OnboardingFlowPending != onboarding.OnboardingFlowPending {
		t.Errorf("Onboarding.OnboardingFlowPending mismatch: got %v, want %v", decoded.Onboarding.OnboardingFlowPending, onboarding.OnboardingFlowPending)
	}

	// Verify Onboarding field is present in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["onboarding"]; !exists {
		t.Error("onboarding field should be present in JSON")
	}
}

// TestAccountOnboardingOmitted tests that Onboarding field is omitted when nil (Requirement 1.9)
func TestAccountOnboardingOmitted(t *testing.T) {
	account := NetbirdAccount{
		ID: "account456",
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             3600,
			PeerLoginExpirationEnabled:      true,
			PeerInactivityExpiration:        7200,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        true,
			JWTGroupsEnabled:                false,
			JWTGroupsClaimName:              "groups",
			JWTAllowGroups:                  []string{"admin"},
			RegularUsersViewBlocked:         false,
		},
		Onboarding: nil, // Explicitly nil
	}

	// Marshal to JSON
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal account: %v", err)
	}

	// Verify Onboarding field is NOT present in JSON when nil
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["onboarding"]; exists {
		t.Error("onboarding field should be omitted when nil")
	}
}

// TestAccountOnboardingBothTrue tests NetbirdAccountOnboarding with both fields true (Requirement 1.9)
func TestAccountOnboardingBothTrue(t *testing.T) {
	onboarding := NetbirdAccountOnboarding{
		SignupFormPending:     true,
		OnboardingFlowPending: true,
	}

	// Marshal to JSON
	data, err := json.Marshal(onboarding)
	if err != nil {
		t.Fatalf("failed to marshal onboarding: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountOnboarding
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal onboarding: %v", err)
	}

	// Verify both fields are true
	if !decoded.SignupFormPending {
		t.Error("SignupFormPending should be true")
	}
	if !decoded.OnboardingFlowPending {
		t.Error("OnboardingFlowPending should be true")
	}
}

// TestAccountOnboardingBothFalse tests NetbirdAccountOnboarding with both fields false (Requirement 1.9)
func TestAccountOnboardingBothFalse(t *testing.T) {
	onboarding := NetbirdAccountOnboarding{
		SignupFormPending:     false,
		OnboardingFlowPending: false,
	}

	// Marshal to JSON
	data, err := json.Marshal(onboarding)
	if err != nil {
		t.Fatalf("failed to marshal onboarding: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccountOnboarding
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal onboarding: %v", err)
	}

	// Verify both fields are false
	if decoded.SignupFormPending {
		t.Error("SignupFormPending should be false")
	}
	if decoded.OnboardingFlowPending {
		t.Error("OnboardingFlowPending should be false")
	}
}

// TestAccountTopLevelFieldsMarshaling tests JSON marshaling/unmarshaling of top-level NetbirdAccount fields (Requirement 1.10)
func TestAccountTopLevelFieldsMarshaling(t *testing.T) {
	domain := "example.netbird.io"
	domainCategory := "private"
	createdAt := "2024-01-15T10:30:00Z"
	createdBy := "user123"

	account := NetbirdAccount{
		ID: "account789",
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             3600,
			PeerLoginExpirationEnabled:      true,
			PeerInactivityExpiration:        7200,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        true,
			JWTGroupsEnabled:                false,
			JWTGroupsClaimName:              "groups",
			JWTAllowGroups:                  []string{"admin"},
			RegularUsersViewBlocked:         false,
		},
		Domain:         &domain,
		DomainCategory: &domainCategory,
		CreatedAt:      &createdAt,
		CreatedBy:      &createdBy,
	}

	// Marshal to JSON
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal account: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccount
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal account: %v", err)
	}

	// Verify all top-level fields are preserved
	if decoded.Domain == nil || *decoded.Domain != domain {
		t.Errorf("Domain mismatch: got %v, want %s", decoded.Domain, domain)
	}
	if decoded.DomainCategory == nil || *decoded.DomainCategory != domainCategory {
		t.Errorf("DomainCategory mismatch: got %v, want %s", decoded.DomainCategory, domainCategory)
	}
	if decoded.CreatedAt == nil || *decoded.CreatedAt != createdAt {
		t.Errorf("CreatedAt mismatch: got %v, want %s", decoded.CreatedAt, createdAt)
	}
	if decoded.CreatedBy == nil || *decoded.CreatedBy != createdBy {
		t.Errorf("CreatedBy mismatch: got %v, want %s", decoded.CreatedBy, createdBy)
	}

	// Verify fields are present in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	expectedFields := []string{
		"domain",
		"domain_category",
		"created_at",
		"created_by",
	}

	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("expected field %s to be present in JSON", field)
		}
	}
}

// TestAccountTopLevelFieldsOmitted tests that top-level fields are omitted when nil (Requirement 1.10)
func TestAccountTopLevelFieldsOmitted(t *testing.T) {
	account := NetbirdAccount{
		ID: "account999",
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             3600,
			PeerLoginExpirationEnabled:      true,
			PeerInactivityExpiration:        7200,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        true,
			JWTGroupsEnabled:                false,
			JWTGroupsClaimName:              "groups",
			JWTAllowGroups:                  []string{"admin"},
			RegularUsersViewBlocked:         false,
		},
		// All top-level pointer fields are nil
		Domain:         nil,
		DomainCategory: nil,
		CreatedAt:      nil,
		CreatedBy:      nil,
	}

	// Marshal to JSON
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal account: %v", err)
	}

	// Verify top-level fields are NOT present in JSON when nil
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	optionalFields := []string{
		"domain",
		"domain_category",
		"created_at",
		"created_by",
	}

	for _, field := range optionalFields {
		if _, exists := jsonMap[field]; exists {
			t.Errorf("field %s should be omitted when nil, but was present", field)
		}
	}

	// Verify required fields are still present
	if _, exists := jsonMap["id"]; !exists {
		t.Error("required field id should be present")
	}
	if _, exists := jsonMap["settings"]; !exists {
		t.Error("required field settings should be present")
	}
}

// TestAccountWithAllFields tests NetbirdAccount with all fields populated (Requirement 1.10)
func TestAccountWithAllFields(t *testing.T) {
	domain := "test.netbird.cloud"
	domainCategory := "cloud"
	createdAt := "2024-02-20T15:45:30Z"
	createdBy := "admin@example.com"
	dnsDomain := "internal.example.com"
	networkRange := "10.10.0.0/16"

	onboarding := &NetbirdAccountOnboarding{
		SignupFormPending:     true,
		OnboardingFlowPending: false,
	}

	extra := &NetbirdAccountExtra{
		PeerApprovalEnabled:                true,
		UserApprovalRequired:               false,
		NetworkTrafficLogsEnabled:          true,
		NetworkTrafficLogsGroups:           []string{"monitoring"},
		NetworkTrafficPacketCounterEnabled: true,
	}

	account := NetbirdAccount{
		ID: "full-account",
		Settings: NetbirdAccountSettings{
			PeerLoginExpiration:             3600,
			PeerLoginExpirationEnabled:      true,
			PeerInactivityExpiration:        7200,
			PeerInactivityExpirationEnabled: true,
			GroupsPropagationEnabled:        true,
			JWTGroupsEnabled:                true,
			JWTGroupsClaimName:              "groups",
			JWTAllowGroups:                  []string{"admin", "users"},
			RegularUsersViewBlocked:         true,
			DNSDomain:                       &dnsDomain,
			NetworkRange:                    &networkRange,
			Extra:                           extra,
		},
		Onboarding:     onboarding,
		Domain:         &domain,
		DomainCategory: &domainCategory,
		CreatedAt:      &createdAt,
		CreatedBy:      &createdBy,
	}

	// Marshal to JSON
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("failed to marshal account: %v", err)
	}

	// Unmarshal back
	var decoded NetbirdAccount
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal account: %v", err)
	}

	// Verify all fields are preserved
	if decoded.ID != account.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, account.ID)
	}
	if decoded.Domain == nil || *decoded.Domain != domain {
		t.Errorf("Domain mismatch: got %v, want %s", decoded.Domain, domain)
	}
	if decoded.DomainCategory == nil || *decoded.DomainCategory != domainCategory {
		t.Errorf("DomainCategory mismatch: got %v, want %s", decoded.DomainCategory, domainCategory)
	}
	if decoded.CreatedAt == nil || *decoded.CreatedAt != createdAt {
		t.Errorf("CreatedAt mismatch: got %v, want %s", decoded.CreatedAt, createdAt)
	}
	if decoded.CreatedBy == nil || *decoded.CreatedBy != createdBy {
		t.Errorf("CreatedBy mismatch: got %v, want %s", decoded.CreatedBy, createdBy)
	}
	if decoded.Onboarding == nil {
		t.Error("Onboarding should not be nil")
	}
	if decoded.Settings.Extra == nil {
		t.Error("Settings.Extra should not be nil")
	}
	if decoded.Settings.DNSDomain == nil || *decoded.Settings.DNSDomain != dnsDomain {
		t.Errorf("Settings.DNSDomain mismatch: got %v, want %s", decoded.Settings.DNSDomain, dnsDomain)
	}
}

// TestGetNetbirdAccountWithTopLevelFields tests getNetbirdAccount with top-level fields (Requirement 1.10)
func TestGetNetbirdAccountWithTopLevelFields(t *testing.T) {
	domain := "api.netbird.io"
	domainCategory := "public"
	createdAt := "2024-03-01T08:00:00Z"
	createdBy := "system"

	mockResp := []NetbirdAccount{
		{
			ID: "account-with-fields",
			Settings: NetbirdAccountSettings{
				PeerLoginExpiration:             3600,
				PeerLoginExpirationEnabled:      true,
				PeerInactivityExpiration:        7200,
				PeerInactivityExpirationEnabled: false,
				GroupsPropagationEnabled:        true,
				JWTGroupsEnabled:                true,
				JWTGroupsClaimName:              "groups",
				JWTAllowGroups:                  []string{"admin"},
				RegularUsersViewBlocked:         false,
			},
			Domain:         &domain,
			DomainCategory: &domainCategory,
			CreatedAt:      &createdAt,
			CreatedBy:      &createdBy,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	mcpnetbird.TestNetbirdClient = mcpnetbird.NewNetbirdClientWithBaseURL(server.URL)
	defer func() { mcpnetbird.TestNetbirdClient = nil }()

	ctx := mcpnetbird.WithNetbirdAPIKey(context.Background(), "test-token")
	account, err := getNetbirdAccount(ctx, GetNetbirdAccountParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.ID != "account-with-fields" {
		t.Errorf("unexpected account ID: got %s, want account-with-fields", account.ID)
	}
	if account.Domain == nil || *account.Domain != domain {
		t.Errorf("Domain mismatch: got %v, want %s", account.Domain, domain)
	}
	if account.DomainCategory == nil || *account.DomainCategory != domainCategory {
		t.Errorf("DomainCategory mismatch: got %v, want %s", account.DomainCategory, domainCategory)
	}
	if account.CreatedAt == nil || *account.CreatedAt != createdAt {
		t.Errorf("CreatedAt mismatch: got %v, want %s", account.CreatedAt, createdAt)
	}
	if account.CreatedBy == nil || *account.CreatedBy != createdBy {
		t.Errorf("CreatedBy mismatch: got %v, want %s", account.CreatedBy, createdBy)
	}
}
