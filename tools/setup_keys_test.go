package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCreateNetbirdSetupKeyParamsMarshaling(t *testing.T) {
	allowExtraDNSLabelsTrue := true
	allowExtraDNSLabelsFalse := false
	autoGroups := []string{"group-1", "group-2"}
	usageLimit := 10
	ephemeralTrue := true

	tests := []struct {
		name     string
		params   CreateNetbirdSetupKeyParams
		expected string
	}{
		{
			name: "with required fields only",
			params: CreateNetbirdSetupKeyParams{
				Name:      "test-key",
				Type:      "reusable",
				ExpiresIn: 3600,
			},
			expected: `{"name":"test-key","type":"reusable","expires_in":3600}`,
		},
		{
			name: "with allow_extra_dns_labels true",
			params: CreateNetbirdSetupKeyParams{
				Name:                "test-key",
				Type:                "reusable",
				ExpiresIn:           3600,
				AllowExtraDNSLabels: &allowExtraDNSLabelsTrue,
			},
			expected: `{"name":"test-key","type":"reusable","expires_in":3600,"allow_extra_dns_labels":true}`,
		},
		{
			name: "with allow_extra_dns_labels false",
			params: CreateNetbirdSetupKeyParams{
				Name:                "test-key",
				Type:                "reusable",
				ExpiresIn:           3600,
				AllowExtraDNSLabels: &allowExtraDNSLabelsFalse,
			},
			expected: `{"name":"test-key","type":"reusable","expires_in":3600,"allow_extra_dns_labels":false}`,
		},
		{
			name: "with all fields",
			params: CreateNetbirdSetupKeyParams{
				Name:                "test-key",
				Type:                "reusable",
				ExpiresIn:           3600,
				AutoGroups:          &autoGroups,
				UsageLimit:          &usageLimit,
				Ephemeral:           &ephemeralTrue,
				AllowExtraDNSLabels: &allowExtraDNSLabelsTrue,
			},
			expected: `{"name":"test-key","type":"reusable","expires_in":3600,"auto_groups":["group-1","group-2"],"usage_limit":10,"ephemeral":true,"allow_extra_dns_labels":true}`,
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
			var decoded CreateNetbirdSetupKeyParams
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.Name != tt.params.Name {
				t.Errorf("name mismatch: got %s, want %s", decoded.Name, tt.params.Name)
			}
			if decoded.Type != tt.params.Type {
				t.Errorf("type mismatch: got %s, want %s", decoded.Type, tt.params.Type)
			}
			if decoded.ExpiresIn != tt.params.ExpiresIn {
				t.Errorf("expires_in mismatch: got %d, want %d", decoded.ExpiresIn, tt.params.ExpiresIn)
			}

			// Check AllowExtraDNSLabels
			if tt.params.AllowExtraDNSLabels != nil {
				if decoded.AllowExtraDNSLabels == nil {
					t.Error("allow_extra_dns_labels should not be nil")
				} else if *decoded.AllowExtraDNSLabels != *tt.params.AllowExtraDNSLabels {
					t.Errorf("allow_extra_dns_labels mismatch: got %v, want %v", *decoded.AllowExtraDNSLabels, *tt.params.AllowExtraDNSLabels)
				}
			}
		})
	}
}

func TestCreateNetbirdSetupKeyParamsOmitEmpty(t *testing.T) {
	allowExtraDNSLabelsTrue := true

	tests := []struct {
		name        string
		params      CreateNetbirdSetupKeyParams
		contains    []string
		notContains []string
	}{
		{
			name: "nil allow_extra_dns_labels is omitted",
			params: CreateNetbirdSetupKeyParams{
				Name:      "test-key",
				Type:      "reusable",
				ExpiresIn: 3600,
			},
			contains:    []string{"name", "type", "expires_in"},
			notContains: []string{"allow_extra_dns_labels", "auto_groups", "usage_limit", "ephemeral"},
		},
		{
			name: "allow_extra_dns_labels is present when set",
			params: CreateNetbirdSetupKeyParams{
				Name:                "test-key",
				Type:                "reusable",
				ExpiresIn:           3600,
				AllowExtraDNSLabels: &allowExtraDNSLabelsTrue,
			},
			contains:    []string{"name", "type", "expires_in", "allow_extra_dns_labels"},
			notContains: []string{"auto_groups", "usage_limit", "ephemeral"},
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

func TestNetbirdSetupKeyMarshaling(t *testing.T) {
	allowExtraDNSLabelsTrue := true
	allowExtraDNSLabelsFalse := false

	tests := []struct {
		name     string
		json     string
		expected NetbirdSetupKey
	}{
		{
			name: "with allow_extra_dns_labels true",
			json: `{"id":"key-123","key":"abc123","name":"test-key","expires":"2024-01-01T00:00:00Z","type":"reusable","valid":true,"revoked":false,"used_times":0,"last_used":"0001-01-01T00:00:00Z","state":"valid","auto_groups":[],"updated_at":"2024-01-01T00:00:00Z","usage_limit":0,"ephemeral":false,"allow_extra_dns_labels":true}`,
			expected: NetbirdSetupKey{
				AllowExtraDNSLabels: &allowExtraDNSLabelsTrue,
			},
		},
		{
			name: "with allow_extra_dns_labels false",
			json: `{"id":"key-123","key":"abc123","name":"test-key","expires":"2024-01-01T00:00:00Z","type":"reusable","valid":true,"revoked":false,"used_times":0,"last_used":"0001-01-01T00:00:00Z","state":"valid","auto_groups":[],"updated_at":"2024-01-01T00:00:00Z","usage_limit":0,"ephemeral":false,"allow_extra_dns_labels":false}`,
			expected: NetbirdSetupKey{
				AllowExtraDNSLabels: &allowExtraDNSLabelsFalse,
			},
		},
		{
			name: "without allow_extra_dns_labels",
			json: `{"id":"key-123","key":"abc123","name":"test-key","expires":"2024-01-01T00:00:00Z","type":"reusable","valid":true,"revoked":false,"used_times":0,"last_used":"0001-01-01T00:00:00Z","state":"valid","auto_groups":[],"updated_at":"2024-01-01T00:00:00Z","usage_limit":0,"ephemeral":false}`,
			expected: NetbirdSetupKey{
				AllowExtraDNSLabels: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded NetbirdSetupKey
			err := json.Unmarshal([]byte(tt.json), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			// Check AllowExtraDNSLabels
			if tt.expected.AllowExtraDNSLabels == nil {
				if decoded.AllowExtraDNSLabels != nil {
					t.Errorf("allow_extra_dns_labels should be nil, got %v", *decoded.AllowExtraDNSLabels)
				}
			} else {
				if decoded.AllowExtraDNSLabels == nil {
					t.Error("allow_extra_dns_labels should not be nil")
				} else if *decoded.AllowExtraDNSLabels != *tt.expected.AllowExtraDNSLabels {
					t.Errorf("allow_extra_dns_labels mismatch: got %v, want %v", *decoded.AllowExtraDNSLabels, *tt.expected.AllowExtraDNSLabels)
				}
			}
		})
	}
}
