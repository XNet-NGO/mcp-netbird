# Requirements Document

## Introduction

This document specifies the requirements for aligning the MCP-NETBIRD codebase with the official NetBird API documentation. The MCP-NETBIRD project is a Model Context Protocol server that provides tools to interact with the NetBird API. Analysis of the official API documentation and current codebase has revealed several gaps in CRUD operations, incomplete data structures, and missing resources that need to be addressed to ensure complete and accurate API coverage.

## Glossary

- **MCP_Server**: The Model Context Protocol server that provides tools to interact with NetBird
- **NetBird_API**: The official NetBird REST API for managing network resources
- **CRUD_Operations**: Create, Read, Update, Delete operations for API resources
- **Data_Structure**: Go struct definitions that represent API request/response objects
- **Tool**: An MCP tool that wraps a NetBird API endpoint
- **Resource**: A NetBird API entity (e.g., peer, group, policy, route)
- **Required_Field**: An API field that must be provided in requests
- **Optional_Field**: An API field that may be omitted in requests
- **Cloud_Only**: Features available only in NetBird Cloud, not self-hosted

## Requirements

### Requirement 1: Complete Account Settings Data Structure

**User Story:** As a developer, I want the account settings structure to include all fields from the API, so that I can manage all account configuration options.

#### Acceptance Criteria

1. WHEN retrieving account settings, THE System SHALL include the regular_users_view_blocked field
2. WHEN retrieving account settings, THE System SHALL include the routing_peer_dns_resolution_enabled field
3. WHEN retrieving account settings, THE System SHALL include the dns_domain field
4. WHEN retrieving account settings, THE System SHALL include the network_range field
5. WHEN retrieving account settings, THE System SHALL include the extra nested object with peer_approval_enabled, user_approval_required, network_traffic_logs_enabled, network_traffic_logs_groups, and network_traffic_packet_counter_enabled fields
6. WHEN retrieving account settings, THE System SHALL include the lazy_connection_enabled field
7. WHEN retrieving account settings, THE System SHALL include the auto_update_version field
8. WHEN retrieving account settings, THE System SHALL include the embedded_idp_enabled field
9. WHEN retrieving account information, THE System SHALL include the onboarding nested object with signup_form_pending and onboarding_flow_pending fields
10. WHEN retrieving account information, THE System SHALL include the domain, domain_category, created_at, and created_by fields

### Requirement 2: Complete Peer Data Structure

**User Story:** As a developer, I want the peer structure to include all fields from the API, so that I can access complete peer information.

#### Acceptance Criteria

1. WHEN retrieving peer information, THE System SHALL include the created_at timestamp field
2. WHEN retrieving peer information, THE System SHALL include the ephemeral boolean field
3. WHEN retrieving peer information, THE System SHALL include the local_flags nested object with rosenpass_enabled, rosenpass_permissive, server_ssh_allowed, disable_client_routes, disable_server_routes, disable_dns, disable_firewall, block_lan_access, block_inbound, and lazy_connection_enabled fields
4. WHEN retrieving peer information, THE System SHALL include the disapproval_reason field for cloud-only deployments

### Requirement 3: Complete Policy Rule Data Structure

**User Story:** As a developer, I want policy rules to include all fields from the API, so that I can configure complete access control policies.

#### Acceptance Criteria

1. WHEN creating or updating policy rules, THE System SHALL support the port_ranges array field with start and end integer fields
2. WHEN creating or updating policy rules, THE System SHALL support the authorized_groups map field
3. WHEN creating or updating policy rules, THE System SHALL support the sourceResource nested object with id and type fields
4. WHEN creating or updating policy rules, THE System SHALL support the destinationResource nested object with id and type fields

### Requirement 4: Complete Posture Check Data Structure

**User Story:** As a developer, I want posture checks to include all check types from the API, so that I can implement comprehensive security policies.

#### Acceptance Criteria

1. WHEN creating or updating posture checks, THE System SHALL support the nb_version_check with min_version field
2. WHEN creating or updating posture checks, THE System SHALL support the os_version_check with android, ios, darwin, linux, and windows nested objects
3. WHEN creating or updating posture checks, THE System SHALL support the geo_location_check with locations array and action field
4. WHEN creating or updating posture checks, THE System SHALL support the peer_network_range_check with ranges array and action field
5. WHEN creating or updating posture checks, THE System SHALL support the process_check with processes array containing linux_path, mac_path, and windows_path fields

### Requirement 5: Complete Group CRUD Operations

**User Story:** As a developer, I want full CRUD operations for groups, so that I can manage group resources programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_group tool
2. THE System SHALL provide a get_netbird_group tool
3. THE System SHALL provide an update_netbird_group tool
4. THE System SHALL provide a delete_netbird_group tool
5. THE System SHALL provide a list_netbird_groups tool
6. WHEN creating or updating groups, THE System SHALL support the resources array field with id and type for each resource

### Requirement 6: Complete Policy CRUD Operations

**User Story:** As a developer, I want full CRUD operations for policies, so that I can manage access control policies programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_policy tool
2. THE System SHALL provide a get_netbird_policy tool
3. THE System SHALL provide an update_netbird_policy tool
4. THE System SHALL provide a delete_netbird_policy tool
5. THE System SHALL provide a list_netbird_policies tool

### Requirement 7: Complete Posture Check CRUD Operations

**User Story:** As a developer, I want full CRUD operations for posture checks, so that I can manage security posture requirements programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_posture_check tool
2. THE System SHALL provide a get_netbird_posture_check tool
3. THE System SHALL provide an update_netbird_posture_check tool
4. THE System SHALL provide a delete_netbird_posture_check tool
5. THE System SHALL provide a list_netbird_posture_checks tool

### Requirement 8: Complete Network CRUD Operations

**User Story:** As a developer, I want full CRUD operations for networks, so that I can manage network resources programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_network tool
2. THE System SHALL provide a get_netbird_network tool
3. THE System SHALL provide an update_netbird_network tool
4. THE System SHALL provide a delete_netbird_network tool
5. THE System SHALL provide a list_netbird_networks tool

### Requirement 9: Network Resources Sub-Resource Operations

**User Story:** As a developer, I want to manage network resources within networks, so that I can configure network access control.

#### Acceptance Criteria

1. THE System SHALL provide a list_netbird_network_resources tool that accepts network_id parameter
2. THE System SHALL provide a create_netbird_network_resource tool that accepts network_id and resource parameters
3. THE System SHALL provide a get_netbird_network_resource tool that accepts network_id and resource_id parameters
4. THE System SHALL provide an update_netbird_network_resource tool that accepts network_id and resource_id parameters
5. THE System SHALL provide a delete_netbird_network_resource tool that accepts network_id and resource_id parameters
6. WHEN creating or updating network resources, THE System SHALL support name, description, address, enabled, and groups fields

### Requirement 10: Network Routers Sub-Resource Operations

**User Story:** As a developer, I want to manage network routers within networks, so that I can configure routing peers.

#### Acceptance Criteria

1. THE System SHALL provide a list_netbird_network_routers tool that accepts network_id parameter
2. THE System SHALL provide a create_netbird_network_router tool that accepts network_id and router parameters
3. THE System SHALL provide a get_netbird_network_router tool that accepts network_id and router_id parameters
4. THE System SHALL provide an update_netbird_network_router tool that accepts network_id and router_id parameters
5. THE System SHALL provide a delete_netbird_network_router tool that accepts network_id and router_id parameters
6. WHEN creating or updating network routers, THE System SHALL support peer, peer_groups, metric, masquerade, and enabled fields

### Requirement 11: Complete Nameserver CRUD Operations

**User Story:** As a developer, I want full CRUD operations for nameservers, so that I can manage DNS configuration programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_nameserver tool
2. THE System SHALL provide a get_netbird_nameserver tool
3. THE System SHALL provide an update_netbird_nameserver tool
4. THE System SHALL provide a delete_netbird_nameserver tool
5. THE System SHALL provide a list_netbird_nameservers tool

### Requirement 12: Complete Setup Key CRUD Operations

**User Story:** As a developer, I want full CRUD operations for setup keys, so that I can manage peer registration keys programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_setup_key tool
2. THE System SHALL provide a get_netbird_setup_key tool
3. THE System SHALL provide an update_netbird_setup_key tool
4. THE System SHALL provide a delete_netbird_setup_key tool
5. THE System SHALL provide a list_netbird_setup_keys tool
6. WHEN creating setup keys, THE System SHALL support the allow_extra_dns_labels field

### Requirement 13: Complete Route CRUD Operations

**User Story:** As a developer, I want full CRUD operations for routes, so that I can manage network routing programmatically.

#### Acceptance Criteria

1. THE System SHALL provide a create_netbird_route tool
2. THE System SHALL provide a get_netbird_route tool
3. THE System SHALL provide an update_netbird_route tool
4. THE System SHALL provide a delete_netbird_route tool
5. THE System SHALL provide a list_netbird_routes tool

### Requirement 14: Ingress Port Allocations (Cloud Only)

**User Story:** As a cloud user, I want to manage ingress port allocations, so that I can configure port forwarding for peers.

#### Acceptance Criteria

1. THE System SHALL provide a list_netbird_port_allocations tool that accepts peer_id parameter
2. THE System SHALL provide a create_netbird_port_allocation tool that accepts peer_id and allocation parameters
3. THE System SHALL provide a get_netbird_port_allocation tool that accepts peer_id and allocation_id parameters
4. THE System SHALL provide an update_netbird_port_allocation tool that accepts peer_id and allocation_id parameters
5. THE System SHALL provide a delete_netbird_port_allocation tool that accepts peer_id and allocation_id parameters
6. WHEN creating or updating port allocations, THE System SHALL support name, enabled, port_ranges, and direct_port fields

### Requirement 15: Required and Optional Field Consistency

**User Story:** As a developer, I want consistent handling of required and optional fields, so that API requests are validated correctly.

#### Acceptance Criteria

1. WHEN a field is marked as required in the API documentation, THE System SHALL mark it as required in the jsonschema tag
2. WHEN a field is marked as optional in the API documentation, THE System SHALL use a pointer type and omitempty in the json tag
3. WHEN validating requests, THE System SHALL reject requests missing required fields
4. WHEN processing responses, THE System SHALL handle missing optional fields gracefully

### Requirement 16: API Endpoint Path Accuracy

**User Story:** As a developer, I want all API endpoint paths to match the official documentation, so that API calls succeed.

#### Acceptance Criteria

1. WHEN making API calls, THE System SHALL use the exact endpoint paths from the API documentation
2. WHEN making API calls, THE System SHALL use the correct HTTP methods (GET, POST, PUT, DELETE)
3. WHEN processing responses, THE System SHALL handle both array and single object responses correctly
4. WHEN constructing URLs, THE System SHALL properly concatenate base URL and endpoint paths

### Requirement 17: Backward Compatibility

**User Story:** As an existing user, I want my current integrations to continue working, so that updates don't break my workflows.

#### Acceptance Criteria

1. WHEN updating data structures, THE System SHALL maintain existing field names and types
2. WHEN adding new fields, THE System SHALL use optional fields with pointer types
3. WHEN updating tools, THE System SHALL maintain existing tool names and parameter structures
4. WHEN processing API responses, THE System SHALL handle both old and new field formats

### Requirement 18: Setup Key Missing Fields

**User Story:** As a developer, I want setup keys to include all fields from the API, so that I can configure peer registration completely.

#### Acceptance Criteria

1. WHEN creating setup keys, THE System SHALL support the allow_extra_dns_labels boolean field
2. WHEN retrieving setup keys, THE System SHALL include the allow_extra_dns_labels field in responses
