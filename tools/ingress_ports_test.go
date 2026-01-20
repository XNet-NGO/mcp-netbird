package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIngressPortRangeMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		portRange IngressPortRange
		expected string
	}{
		{
			name: "basic port range",
			portRange: IngressPortRange{
				Start:    80,
				End:      320,
				Protocol: "tcp",
			},
			expected: `{"start":80,"end":320,"protocol":"tcp"}`,
		},
		{
			name: "udp port range",
			portRange: IngressPortRange{
				Start:    5000,
				End:      5100,
				Protocol: "udp",
			},
			expected: `{"start":5000,"end":5100,"protocol":"udp"}`,
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
			var decoded IngressPortRange
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.Start != tt.portRange.Start || decoded.End != tt.portRange.End || decoded.Protocol != tt.portRange.Protocol {
				t.Errorf("unmarshaling failed: got %+v, want %+v", decoded, tt.portRange)
			}
		})
	}
}

func TestDirectPortMarshaling(t *testing.T) {
	tests := []struct {
		name       string
		directPort DirectPort
		expected   string
	}{
		{
			name: "basic direct port",
			directPort: DirectPort{
				Count:    5,
				Protocol: "udp",
			},
			expected: `{"count":5,"protocol":"udp"}`,
		},
		{
			name: "tcp direct port",
			directPort: DirectPort{
				Count:    10,
				Protocol: "tcp",
			},
			expected: `{"count":10,"protocol":"tcp"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.directPort)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("marshaling failed: got %s, want %s", string(data), tt.expected)
			}

			// Test unmarshaling
			var decoded DirectPort
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.Count != tt.directPort.Count || decoded.Protocol != tt.directPort.Protocol {
				t.Errorf("unmarshaling failed: got %+v, want %+v", decoded, tt.directPort)
			}
		})
	}
}

func TestCreateNetbirdPortAllocationParamsMarshaling(t *testing.T) {
	portRanges := []IngressPortRange{
		{Start: 80, End: 320, Protocol: "tcp"},
	}
	directPort := DirectPort{Count: 5, Protocol: "udp"}

	tests := []struct {
		name     string
		params   CreateNetbirdPortAllocationParams
		expected string
	}{
		{
			name: "with required fields only",
			params: CreateNetbirdPortAllocationParams{
				PeerID:  "peer-123",
				Name:    "test-allocation",
				Enabled: true,
			},
			expected: `{"peer_id":"peer-123","name":"test-allocation","enabled":true}`,
		},
		{
			name: "with port_ranges",
			params: CreateNetbirdPortAllocationParams{
				PeerID:     "peer-123",
				Name:       "test-allocation",
				Enabled:    true,
				PortRanges: &portRanges,
			},
			expected: `{"peer_id":"peer-123","name":"test-allocation","enabled":true,"port_ranges":[{"start":80,"end":320,"protocol":"tcp"}]}`,
		},
		{
			name: "with direct_port",
			params: CreateNetbirdPortAllocationParams{
				PeerID:     "peer-123",
				Name:       "test-allocation",
				Enabled:    true,
				DirectPort: &directPort,
			},
			expected: `{"peer_id":"peer-123","name":"test-allocation","enabled":true,"direct_port":{"count":5,"protocol":"udp"}}`,
		},
		{
			name: "with all fields",
			params: CreateNetbirdPortAllocationParams{
				PeerID:     "peer-123",
				Name:       "test-allocation",
				Enabled:    true,
				PortRanges: &portRanges,
				DirectPort: &directPort,
			},
			expected: `{"peer_id":"peer-123","name":"test-allocation","enabled":true,"port_ranges":[{"start":80,"end":320,"protocol":"tcp"}],"direct_port":{"count":5,"protocol":"udp"}}`,
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
			var decoded CreateNetbirdPortAllocationParams
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.PeerID != tt.params.PeerID {
				t.Errorf("peer_id mismatch: got %s, want %s", decoded.PeerID, tt.params.PeerID)
			}
			if decoded.Name != tt.params.Name {
				t.Errorf("name mismatch: got %s, want %s", decoded.Name, tt.params.Name)
			}
			if decoded.Enabled != tt.params.Enabled {
				t.Errorf("enabled mismatch: got %v, want %v", decoded.Enabled, tt.params.Enabled)
			}
		})
	}
}

func TestCreateNetbirdPortAllocationParamsOmitEmpty(t *testing.T) {
	portRanges := []IngressPortRange{
		{Start: 80, End: 320, Protocol: "tcp"},
	}

	tests := []struct {
		name        string
		params      CreateNetbirdPortAllocationParams
		contains    []string
		notContains []string
	}{
		{
			name: "nil port_ranges and direct_port are omitted",
			params: CreateNetbirdPortAllocationParams{
				PeerID:  "peer-123",
				Name:    "test-allocation",
				Enabled: true,
			},
			contains:    []string{"peer_id", "name", "enabled"},
			notContains: []string{"port_ranges", "direct_port"},
		},
		{
			name: "port_ranges is present, direct_port is omitted",
			params: CreateNetbirdPortAllocationParams{
				PeerID:     "peer-123",
				Name:       "test-allocation",
				Enabled:    true,
				PortRanges: &portRanges,
			},
			contains:    []string{"peer_id", "name", "enabled", "port_ranges"},
			notContains: []string{"direct_port"},
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

func TestUpdateNetbirdPortAllocationParamsMarshaling(t *testing.T) {
	name := "updated-allocation"
	enabled := false
	portRanges := []IngressPortRange{
		{Start: 80, End: 320, Protocol: "tcp"},
	}
	directPort := DirectPort{Count: 5, Protocol: "udp"}

	tests := []struct {
		name     string
		params   UpdateNetbirdPortAllocationParams
		expected string
	}{
		{
			name: "with required fields only",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
			},
			expected: `{"peer_id":"peer-123","allocation_id":"alloc-456"}`,
		},
		{
			name: "with name",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
				Name:         &name,
			},
			expected: `{"peer_id":"peer-123","allocation_id":"alloc-456","name":"updated-allocation"}`,
		},
		{
			name: "with enabled",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
				Enabled:      &enabled,
			},
			expected: `{"peer_id":"peer-123","allocation_id":"alloc-456","enabled":false}`,
		},
		{
			name: "with all fields",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
				Name:         &name,
				Enabled:      &enabled,
				PortRanges:   &portRanges,
				DirectPort:   &directPort,
			},
			expected: `{"peer_id":"peer-123","allocation_id":"alloc-456","name":"updated-allocation","enabled":false,"port_ranges":[{"start":80,"end":320,"protocol":"tcp"}],"direct_port":{"count":5,"protocol":"udp"}}`,
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
			var decoded UpdateNetbirdPortAllocationParams
			err = json.Unmarshal([]byte(tt.expected), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if decoded.PeerID != tt.params.PeerID {
				t.Errorf("peer_id mismatch: got %s, want %s", decoded.PeerID, tt.params.PeerID)
			}
			if decoded.AllocationID != tt.params.AllocationID {
				t.Errorf("allocation_id mismatch: got %s, want %s", decoded.AllocationID, tt.params.AllocationID)
			}
		})
	}
}

func TestUpdateNetbirdPortAllocationParamsOmitEmpty(t *testing.T) {
	name := "updated-allocation"

	tests := []struct {
		name        string
		params      UpdateNetbirdPortAllocationParams
		contains    []string
		notContains []string
	}{
		{
			name: "nil optional fields are omitted",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
			},
			contains:    []string{"peer_id", "allocation_id"},
			notContains: []string{"name", "enabled", "port_ranges", "direct_port"},
		},
		{
			name: "only name is present",
			params: UpdateNetbirdPortAllocationParams{
				PeerID:       "peer-123",
				AllocationID: "alloc-456",
				Name:         &name,
			},
			contains:    []string{"peer_id", "allocation_id", "name"},
			notContains: []string{"enabled", "port_ranges", "direct_port"},
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

func TestNetbirdPortAllocationsMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected NetbirdPortAllocations
	}{
		{
			name: "complete port allocation",
			json: `{"id":"alloc-123","name":"test-allocation","ingress_peer_id":"peer-456","region":"germany","enabled":true,"ingress_ip":"192.34.0.123","port_range_mappings":[{"translated_start":80,"translated_end":320,"ingress_start":1080,"ingress_end":1320,"protocol":"tcp"}]}`,
			expected: NetbirdPortAllocations{
				ID:            "alloc-123",
				Name:          "test-allocation",
				IngressPeerID: "peer-456",
				Region:        "germany",
				Enabled:       true,
				IngressIP:     "192.34.0.123",
				PortRangeMappings: []PortRangeMapping{
					{
						TranslatedStart: 80,
						TranslatedEnd:   320,
						IngressStart:    1080,
						IngressEnd:      1320,
						Protocol:        "tcp",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded NetbirdPortAllocations
			err := json.Unmarshal([]byte(tt.json), &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if decoded.ID != tt.expected.ID {
				t.Errorf("id mismatch: got %s, want %s", decoded.ID, tt.expected.ID)
			}
			if decoded.Name != tt.expected.Name {
				t.Errorf("name mismatch: got %s, want %s", decoded.Name, tt.expected.Name)
			}
			if decoded.Enabled != tt.expected.Enabled {
				t.Errorf("enabled mismatch: got %v, want %v", decoded.Enabled, tt.expected.Enabled)
			}
			if len(decoded.PortRangeMappings) != len(tt.expected.PortRangeMappings) {
				t.Errorf("port_range_mappings length mismatch: got %d, want %d", len(decoded.PortRangeMappings), len(tt.expected.PortRangeMappings))
			}
		})
	}
}
