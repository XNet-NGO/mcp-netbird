# Design Document: API Alignment Improvements

## Overview

This design document outlines the approach for aligning the MCP-NETBIRD codebase with the official NetBird API documentation. The implementation will focus on completing data structures, adding missing CRUD operations, and ensuring consistency with the API specification while maintaining backward compatibility.

The design follows a systematic approach:
1. Update existing data structures to include all API fields
2. Add missing CRUD operations for existing resources
3. Implement new sub-resource operations (network resources, network routers)
4. Ensure proper field validation (required vs optional)
5. Maintain backward compatibility with existing integrations

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    MCP Server Layer                          │
│  (Tool Registration & Request Handling)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Tools Layer                               │
│  - Account Tools                                             │
│  - Peer Tools                                                │
│  - Group Tools (Enhanced)                                    │
│  - Policy Tools (Enhanced)                                   │
│  - Posture Check Tools (Enhanced)                            │
│  - Network Tools (New)                                       │
│  - Network Resource Tools (New)                              │
│  - Network Router Tools (New)                                │
│  - Nameserver Tools (Enhanced)                               │
│  - Setup Key Tools (Enhanced)                                │
│  - Route Tools (Enhanced)                                    │
│  - Port Allocation Tools (New)                               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                NetbirdClient Layer                           │
│  (HTTP Client with Get/Post/Put/Delete methods)              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                NetBird REST API                              │
│  (api.netbird.io/api)                                        │
└─────────────────────────────────────────────────────────────┘
```

### Design Principles

1. **Completeness**: All API fields must be represented in Go structs
2. **Consistency**: Required fields use jsonschema:"required", optional fields use pointers and omitempty
3. **Backward Compatibility**: Existing field names and types remain unchanged
4. **Modularity**: Each resource type has its own file in the tools/ directory
5. **Testability**: All tools follow the same pattern for easy testing

## Components and Interfaces

### 1. Data Structure Updates

#### Account Settings Structure

```go
type NetbirdAccountExtra struct {
    PeerApprovalEnabled              bool     `json:"peer_approval_enabled"`
    UserApprovalRequired             bool     `json:"user_approval_required"`
    NetworkTrafficLogsEnabled        bool     `json:"network_traffic_logs_enabled"`
    NetworkTrafficLogsGroups         []string `json:"network_traffic_logs_groups"`
    NetworkTrafficPacketCounterEnabled bool   `json:"network_traffic_packet_counter_enabled"`
}

type NetbirdAccountOnboarding struct {
    SignupFormPending      bool `json:"signup_form_pending"`
    OnboardingFlowPending  bool `json:"onboarding_flow_pending"`
}

type NetbirdAccountSettings struct {
    // Existing fields
    PeerLoginExpiration                int      `json:"peer_login_expiration"`
    PeerLoginExpirationEnabled         bool     `json:"peer_login_expiration_enabled"`
    PeerInactivityExpiration           int      `json:"peer_inactivity_expiration"`
    PeerInactivityExpirationEnabled    bool     `json:"peer_inactivity_expiration_enabled"`
    GroupsPropagationEnabled           bool     `json:"groups_propagation_enabled"`
    JWTGroupsEnabled                   bool     `json:"jwt_groups_enabled"`
    JWTGroupsClaimName                 string   `json:"jwt_groups_claim_name"`
    JWTAllowGroups                     []string `json:"jwt_allow_groups"`
    
    // New fields
    RegularUsersViewBlocked            bool     `json:"regular_users_view_blocked"`
    RoutingPeerDNSResolutionEnabled    *bool    `json:"routing_peer_dns_resolution_enabled,omitempty"`
    DNSDomain                          *string  `json:"dns_domain,omitempty"`
    NetworkRange                       *string  `json:"network_range,omitempty"`
    Extra                              *NetbirdAccountExtra `json:"extra,omitempty"`
    LazyConnectionEnabled              *bool    `json:"lazy_connection_enabled,omitempty"`
    AutoUpdateVersion                  *string  `json:"auto_update_version,omitempty"`
    EmbeddedIDPEnabled                 *bool    `json:"embedded_idp_enabled,omitempty"`
}

type NetbirdAccount struct {
    ID             string                     `json:"id"`
    Settings       NetbirdAccountSettings     `json:"settings"`
    Domain         *string                    `json:"domain,omitempty"`
    DomainCategory *string                    `json:"domain_category,omitempty"`
    CreatedAt      *string                    `json:"created_at,omitempty"`
    CreatedBy      *string                    `json:"created_by,omitempty"`
    Onboarding     *NetbirdAccountOnboarding  `json:"onboarding,omitempty"`
}
```

#### Peer Structure

```go
type NetbirdPeerLocalFlags struct {
    RosenpassEnabled       bool `json:"rosenpass_enabled"`
    RosenpassPermissive    bool `json:"rosenpass_permissive"`
    ServerSSHAllowed       bool `json:"server_ssh_allowed"`
    DisableClientRoutes    bool `json:"disable_client_routes"`
    DisableServerRoutes    bool `json:"disable_server_routes"`
    DisableDNS             bool `json:"disable_dns"`
    DisableFirewall        bool `json:"disable_firewall"`
    BlockLANAccess         bool `json:"block_lan_access"`
    BlockInbound           bool `json:"block_inbound"`
    LazyConnectionEnabled  bool `json:"lazy_connection_enabled"`
}

type NetbirdPeer struct {
    // Existing fields...
    ID                      string              `json:"id"`
    Name                    string              `json:"name"`
    IP                      string              `json:"ip"`
    // ... other existing fields ...
    
    // New fields
    CreatedAt               *string                  `json:"created_at,omitempty"`
    Ephemeral               *bool                    `json:"ephemeral,omitempty"`
    LocalFlags              *NetbirdPeerLocalFlags   `json:"local_flags,omitempty"`
    DisapprovalReason       *string                  `json:"disapproval_reason,omitempty"`
}
```

#### Policy Rule Structure

```go
type PortRange struct {
    Start int `json:"start"`
    End   int `json:"end"`
}

type ResourceReference struct {
    ID   string `json:"id"`
    Type string `json:"type"`
}

type NetbirdPolicyRule struct {
    // Existing fields...
    Action        string   `json:"action"`
    Bidirectional bool     `json:"bidirectional"`
    Protocol      string   `json:"protocol"`
    Ports         []string `json:"ports,omitempty"`
    
    // New fields
    PortRanges           *[]PortRange                `json:"port_ranges,omitempty"`
    AuthorizedGroups     *map[string][]string        `json:"authorized_groups,omitempty"`
    SourceResource       *ResourceReference          `json:"sourceResource,omitempty"`
    DestinationResource  *ResourceReference          `json:"destinationResource,omitempty"`
}
```

### 2. New Tool Implementations

#### Network Tools (tools/networks.go)

```go
// List Networks
type ListNetbirdNetworksParams struct{}
func listNetbirdNetworks(ctx context.Context, args ListNetbirdNetworksParams) ([]NetbirdNetwork, error)

// Get Network
type GetNetbirdNetworkParams struct {
    NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}
func getNetbirdNetwork(ctx context.Context, args GetNetbirdNetworkParams) (*NetbirdNetwork, error)

// Create Network
type CreateNetbirdNetworkParams struct {
    Name        string  `json:"name" jsonschema:"required,description=Network name"`
    Description *string `json:"description,omitempty" jsonschema:"description=Network description"`
}
func createNetbirdNetwork(ctx context.Context, args CreateNetbirdNetworkParams) (*NetbirdNetwork, error)

// Update Network
type UpdateNetbirdNetworkParams struct {
    NetworkID   string  `json:"network_id" jsonschema:"required,description=The ID of the network to update"`
    Name        *string `json:"name,omitempty" jsonschema:"description=Network name"`
    Description *string `json:"description,omitempty" jsonschema:"description=Network description"`
}
func updateNetbirdNetwork(ctx context.Context, args UpdateNetbirdNetworkParams) (*NetbirdNetwork, error)

// Delete Network
type DeleteNetbirdNetworkParams struct {
    NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network to delete"`
}
func deleteNetbirdNetwork(ctx context.Context, args DeleteNetbirdNetworkParams) (map[string]string, error)
```

#### Network Resource Tools (tools/network_resources.go)

```go
// List Network Resources
type ListNetbirdNetworkResourcesParams struct {
    NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}
func listNetbirdNetworkResources(ctx context.Context, args ListNetbirdNetworkResourcesParams) ([]NetbirdNetworkResource, error)

// Create Network Resource
type CreateNetbirdNetworkResourceParams struct {
    NetworkID   string    `json:"network_id" jsonschema:"required,description=The ID of the network"`
    Name        string    `json:"name" jsonschema:"required,description=Resource name"`
    Description *string   `json:"description,omitempty" jsonschema:"description=Resource description"`
    Address     string    `json:"address" jsonschema:"required,description=Resource address (IP, subnet, or domain)"`
    Enabled     bool      `json:"enabled" jsonschema:"required,description=Resource status"`
    Groups      []string  `json:"groups" jsonschema:"required,description=Group IDs containing the resource"`
}
func createNetbirdNetworkResource(ctx context.Context, args CreateNetbirdNetworkResourceParams) (*NetbirdNetworkResource, error)

// Similar patterns for Get, Update, Delete
```

#### Network Router Tools (tools/network_routers.go)

```go
// List Network Routers
type ListNetbirdNetworkRoutersParams struct {
    NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}
func listNetbirdNetworkRouters(ctx context.Context, args ListNetbirdNetworkRoutersParams) ([]NetbirdNetworkRouter, error)

// Create Network Router
type CreateNetbirdNetworkRouterParams struct {
    NetworkID   string    `json:"network_id" jsonschema:"required,description=The ID of the network"`
    Peer        *string   `json:"peer,omitempty" jsonschema:"description=Peer ID (cannot be used with peer_groups)"`
    PeerGroups  *[]string `json:"peer_groups,omitempty" jsonschema:"description=Peer group IDs (cannot be used with peer)"`
    Metric      int       `json:"metric" jsonschema:"required,description=Route metric (1-9999)"`
    Masquerade  bool      `json:"masquerade" jsonschema:"required,description=Enable masquerading"`
    Enabled     bool      `json:"enabled" jsonschema:"required,description=Router status"`
}
func createNetbirdNetworkRouter(ctx context.Context, args CreateNetbirdNetworkRouterParams) (*NetbirdNetworkRouter, error)

// Similar patterns for Get, Update, Delete
```

### 3. Enhanced Existing Tools

#### Groups (tools/groups.go)

Add support for resources field:

```go
type GroupResource struct {
    ID   string `json:"id"`
    Type string `json:"type"`
}

type CreateNetbirdGroupParams struct {
    Name      string           `json:"name" jsonschema:"required,description=Group name"`
    Peers     *[]string        `json:"peers,omitempty" jsonschema:"description=Peer IDs"`
    Resources *[]GroupResource `json:"resources,omitempty" jsonschema:"description=Resource references"`
}
```

#### Setup Keys (tools/setup_keys.go)

Add allow_extra_dns_labels field:

```go
type CreateNetbirdSetupKeyParams struct {
    // Existing fields...
    Name                  string    `json:"name" jsonschema:"required,description=Setup key name"`
    Type                  string    `json:"type" jsonschema:"required,description=Key type"`
    ExpiresIn             int       `json:"expires_in" jsonschema:"required,description=Expiration in seconds"`
    AutoGroups            *[]string `json:"auto_groups,omitempty" jsonschema:"description=Auto-assign groups"`
    UsageLimit            *int      `json:"usage_limit,omitempty" jsonschema:"description=Usage limit"`
    Ephemeral             *bool     `json:"ephemeral,omitempty" jsonschema:"description=Ephemeral peer"`
    
    // New field
    AllowExtraDNSLabels   *bool     `json:"allow_extra_dns_labels,omitempty" jsonschema:"description=Allow extra DNS labels"`
}
```

## Data Models

### Core Data Models

All data models follow these conventions:
- Required fields: No pointer, no omitempty
- Optional fields: Pointer type with omitempty
- Nested objects: Separate struct definitions
- Arrays: Use slice types
- Maps: Use map[string]T types

### Validation Rules

1. **Required Field Validation**: Enforced by jsonschema:"required" tag
2. **Type Validation**: Enforced by Go type system
3. **Range Validation**: Documented in jsonschema description (e.g., "1-9999")
4. **Mutual Exclusivity**: Documented in description (e.g., "cannot be used with peer_groups")

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Account Settings Structure Completeness

*For any* account settings object retrieved from the API, all required fields from the API specification (regular_users_view_blocked, routing_peer_dns_resolution_enabled, dns_domain, network_range, extra, lazy_connection_enabled, auto_update_version, embedded_idp_enabled) should be present in the Go struct and properly unmarshaled.

**Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8**

### Property 2: Account Information Completeness

*For any* account object retrieved from the API, all top-level fields (domain, domain_category, created_at, created_by, onboarding) should be present in the Go struct and properly unmarshaled.

**Validates: Requirements 1.9, 1.10**

### Property 3: Peer Structure Completeness

*For any* peer object retrieved from the API, all required fields (created_at, ephemeral, local_flags with all subfields, disapproval_reason) should be present in the Go struct and properly unmarshaled.

**Validates: Requirements 2.1, 2.2, 2.3, 2.4**

### Property 4: Policy Rule Structure Completeness

*For any* policy rule object created or updated, all required fields (port_ranges, authorized_groups, sourceResource, destinationResource) should be supported in the Go struct and properly marshaled.

**Validates: Requirements 3.1, 3.2, 3.3, 3.4**

### Property 5: Posture Check Structure Completeness

*For any* posture check object created or updated, all check types (nb_version_check, os_version_check, geo_location_check, peer_network_range_check, process_check) should be supported in the Go struct with all required subfields.

**Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5**

### Property 6: Group Resources Field Support

*For any* group create or update operation, the resources field with id and type subfields should be properly marshaled and sent to the API.

**Validates: Requirement 5.6**

### Property 7: Network Resource Field Support

*For any* network resource create or update operation, all required fields (name, description, address, enabled, groups) should be properly marshaled and sent to the API.

**Validates: Requirement 9.6**

### Property 8: Network Router Field Support

*For any* network router create or update operation, all required fields (peer, peer_groups, metric, masquerade, enabled) should be properly marshaled and sent to the API.

**Validates: Requirement 10.6**

### Property 9: Setup Key Extra DNS Labels Field

*For any* setup key create operation, the allow_extra_dns_labels field should be properly marshaled and sent to the API when provided.

**Validates: Requirements 12.6, 18.1, 18.2**

### Property 10: Port Allocation Field Support

*For any* port allocation create or update operation, all required fields (name, enabled, port_ranges, direct_port) should be properly marshaled and sent to the API.

**Validates: Requirement 14.6**

### Property 11: Required Field Validation

*For any* API request with required fields, omitting a required field should result in a validation error before the request is sent.

**Validates: Requirement 15.1, 15.3**

### Property 12: Optional Field Handling

*For any* API request with optional fields, omitting an optional field should result in that field being excluded from the JSON payload (not sent as null).

**Validates: Requirement 15.2, 15.4**

### Property 13: API Endpoint Path Correctness

*For any* API call, the HTTP method and endpoint path should exactly match the official API documentation.

**Validates: Requirements 16.1, 16.2, 16.4**

### Property 14: Response Type Handling

*For any* API response, the system should correctly handle both array responses (e.g., list operations) and single object responses (e.g., get operations).

**Validates: Requirement 16.3**

### Property 15: Backward Compatibility

*For any* existing tool or data structure, updating to include new fields should not break existing code that uses the old structure.

**Validates: Requirements 17.1, 17.2, 17.3, 17.4**

## Error Handling

### Error Categories

1. **Validation Errors**: Missing required fields, invalid field values
2. **HTTP Errors**: Network failures, timeout errors
3. **API Errors**: 4xx and 5xx responses from NetBird API
4. **Marshaling Errors**: JSON encoding/decoding failures

### Error Handling Strategy

```go
// All tool functions follow this pattern:
func toolFunction(ctx context.Context, args Params) (*Result, error) {
    client := mcpnetbird.NewNetbirdClient()
    
    // Validation happens automatically via jsonschema tags
    
    var result Result
    if err := client.Method(ctx, "/endpoint", args, &result); err != nil {
        // Error is already wrapped with context by NetbirdClient
        return nil, err
    }
    
    return &result, nil
}
```

### Error Messages

All errors should include:
- Operation being performed
- Resource type and ID (if applicable)
- HTTP status code (for API errors)
- Original error message

Example: `"making request: unexpected status code: 404, body: {\"message\":\"route not found\"}"`

## Testing Strategy

### Unit Testing

Unit tests will verify:
1. **Data Structure Marshaling**: JSON encoding/decoding of all structs
2. **Field Presence**: All required fields are present in structs
3. **Optional Field Handling**: Optional fields use pointers and omitempty
4. **Tool Registration**: All tools are properly registered with the MCP server

Example unit test:

```go
func TestAccountSettingsMarshaling(t *testing.T) {
    settings := NetbirdAccountSettings{
        PeerLoginExpiration: 3600,
        RegularUsersViewBlocked: true,
        Extra: &NetbirdAccountExtra{
            PeerApprovalEnabled: true,
        },
    }
    
    data, err := json.Marshal(settings)
    require.NoError(t, err)
    
    var decoded NetbirdAccountSettings
    err = json.Unmarshal(data, &decoded)
    require.NoError(t, err)
    
    assert.Equal(t, settings.PeerLoginExpiration, decoded.PeerLoginExpiration)
    assert.Equal(t, settings.RegularUsersViewBlocked, decoded.RegularUsersViewBlocked)
    assert.NotNil(t, decoded.Extra)
    assert.Equal(t, settings.Extra.PeerApprovalEnabled, decoded.Extra.PeerApprovalEnabled)
}
```

### Property-Based Testing

Property tests will verify universal correctness properties across all inputs. Each property test will:
- Run minimum 100 iterations
- Generate random valid inputs
- Verify the property holds for all inputs
- Reference the design document property

Example property test:

```go
// Feature: api-alignment-improvements, Property 1: Account Settings Structure Completeness
func TestProperty_AccountSettingsCompleteness(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // Generate random account settings
        settings := generateRandomAccountSettings(t)
        
        // Marshal to JSON
        data, err := json.Marshal(settings)
        require.NoError(t, err)
        
        // Unmarshal back
        var decoded NetbirdAccountSettings
        err = json.Unmarshal(data, &decoded)
        require.NoError(t, err)
        
        // Verify all fields are preserved
        assert.Equal(t, settings.RegularUsersViewBlocked, decoded.RegularUsersViewBlocked)
        if settings.RoutingPeerDNSResolutionEnabled != nil {
            assert.Equal(t, *settings.RoutingPeerDNSResolutionEnabled, *decoded.RoutingPeerDNSResolutionEnabled)
        }
        // ... verify all other fields
    })
}
```

### Integration Testing

Integration tests will verify:
1. **End-to-End Flows**: Complete CRUD operations against a test API
2. **Error Handling**: Proper handling of API errors
3. **Backward Compatibility**: Existing tests continue to pass

### Test Configuration

- Property tests: Minimum 100 iterations per test
- Test tags: Each test references its design property
- Test isolation: Each test uses a fresh client instance
- Mock server: Use httptest for unit tests, real API for integration tests

### Dual Testing Approach

The testing strategy uses both unit tests and property-based tests:

- **Unit tests** verify specific examples, edge cases, and error conditions
- **Property tests** verify universal properties across all inputs
- Together they provide comprehensive coverage: unit tests catch concrete bugs, property tests verify general correctness

Unit tests should focus on:
- Specific examples that demonstrate correct behavior
- Integration points between components
- Edge cases and error conditions

Property tests should focus on:
- Universal properties that hold for all inputs
- Comprehensive input coverage through randomization
