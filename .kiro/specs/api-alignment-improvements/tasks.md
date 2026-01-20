# Implementation Plan: API Alignment Improvements

## Overview

This implementation plan breaks down the API alignment improvements into discrete, incremental tasks. Each task builds on previous work and includes validation through code and tests. The plan follows a systematic approach: update data structures first, then add missing operations, and finally ensure consistency and backward compatibility.

## Tasks

- [ ] 1. Update Account data structures
  - [x] 1.1 Add missing fields to NetbirdAccountSettings struct
    - Add RegularUsersViewBlocked, RoutingPeerDNSResolutionEnabled, DNSDomain, NetworkRange fields
    - Add LazyConnectionEnabled, AutoUpdateVersion, EmbeddedIDPEnabled fields
    - Use pointer types with omitempty for optional fields
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.6, 1.7, 1.8_
  
  - [x] 1.2 Create NetbirdAccountExtra nested struct
    - Define struct with PeerApprovalEnabled, UserApprovalRequired, NetworkTrafficLogsEnabled fields
    - Add NetworkTrafficLogsGroups, NetworkTrafficPacketCounterEnabled fields
    - Add Extra field to NetbirdAccountSettings as pointer with omitempty
    - _Requirements: 1.5_
  
  - [x] 1.3 Create NetbirdAccountOnboarding nested struct
    - Define struct with SignupFormPending, OnboardingFlowPending fields
    - Add Onboarding field to NetbirdAccount as pointer with omitempty
    - _Requirements: 1.9_
  
  - [x] 1.4 Add top-level fields to NetbirdAccount struct
    - Add Domain, DomainCategory, CreatedAt, CreatedBy fields as pointers with omitempty
    - _Requirements: 1.10_
  
  - [x] 1.5 Write unit tests for account data structures
    - Test JSON marshaling/unmarshaling of all new fields
    - Test that optional fields are omitted when nil
    - Test that required fields are always present
    - _Requirements: 1.1-1.10_

- [ ] 2. Update Peer data structures
  - [x] 2.1 Create NetbirdPeerLocalFlags nested struct
    - Define struct with all 10 boolean fields (rosenpass_enabled, server_ssh_allowed, etc.)
    - Add LocalFlags field to NetbirdPeer as pointer with omitempty
    - _Requirements: 2.3_
  
  - [x] 2.2 Add missing fields to NetbirdPeer struct
    - Add CreatedAt, Ephemeral, DisapprovalReason fields as pointers with omitempty
    - _Requirements: 2.1, 2.2, 2.4_
  
  - [x] 2.3 Write unit tests for peer data structures
    - Test JSON marshaling/unmarshaling of all new fields
    - Test LocalFlags nested object handling
    - _Requirements: 2.1-2.4_

- [ ] 3. Update Policy data structures
  - [x] 3.1 Create PortRange struct
    - Define struct with Start and End integer fields
    - _Requirements: 3.1_
  
  - [x] 3.2 Create ResourceReference struct
    - Define struct with ID and Type string fields
    - _Requirements: 3.3, 3.4_
  
  - [x] 3.3 Add missing fields to NetbirdPolicyRule struct
    - Add PortRanges field as pointer to slice of PortRange with omitempty
    - Add AuthorizedGroups field as pointer to map[string][]string with omitempty
    - Add SourceResource and DestinationResource fields as pointers to ResourceReference with omitempty
    - _Requirements: 3.1, 3.2, 3.3, 3.4_
  
  - [x] 3.4 Write unit tests for policy data structures
    - Test JSON marshaling/unmarshaling of all new fields
    - Test PortRange and ResourceReference nested objects
    - _Requirements: 3.1-3.4_

- [x] 4. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Add Group resources field support
  - [x] 5.1 Create GroupResource struct
    - Define struct with ID and Type string fields
    - _Requirements: 5.6_
  
  - [x] 5.2 Update CreateNetbirdGroupParams struct
    - Add Resources field as pointer to slice of GroupResource with omitempty
    - _Requirements: 5.6_
  
  - [x] 5.3 Update UpdateNetbirdGroupParams struct
    - Add Resources field as pointer to slice of GroupResource with omitempty
    - _Requirements: 5.6_
  
  - [x] 5.4 Write unit tests for group resources field
    - Test JSON marshaling of resources field
    - Test that resources field is omitted when nil
    - _Requirements: 5.6_

- [ ] 6. Implement Network CRUD operations
  - [x] 6.1 Create tools/networks.go file
    - Define NetbirdNetwork struct with id, name, description, routers, resources, policies fields
    - _Requirements: 8.1-8.5_
  
  - [x] 6.2 Implement list_netbird_networks tool
    - Create ListNetbirdNetworksParams struct (empty)
    - Implement listNetbirdNetworks function using GET /networks
    - Register tool with MCP server
    - _Requirements: 8.5_
  
  - [x] 6.3 Implement get_netbird_network tool
    - Create GetNetbirdNetworkParams struct with network_id field
    - Implement getNetbirdNetwork function using GET /networks/{networkId}
    - Register tool with MCP server
    - _Requirements: 8.2_
  
  - [x] 6.4 Implement create_netbird_network tool
    - Create CreateNetbirdNetworkParams struct with name (required) and description (optional) fields
    - Implement createNetbirdNetwork function using POST /networks
    - Register tool with MCP server
    - _Requirements: 8.1_
  
  - [x] 6.5 Implement update_netbird_network tool
    - Create UpdateNetbirdNetworkParams struct with network_id (required), name and description (optional) fields
    - Implement updateNetbirdNetwork function using PUT /networks/{networkId}
    - Register tool with MCP server
    - _Requirements: 8.3_
  
  - [x] 6.6 Implement delete_netbird_network tool
    - Create DeleteNetbirdNetworkParams struct with network_id field
    - Implement deleteNetbirdNetwork function using DELETE /networks/{networkId}
    - Register tool with MCP server
    - _Requirements: 8.4_
  
  - [x] 6.7 Create AddNetbirdNetworkTools function
    - Register all network tools with MCP server
    - _Requirements: 8.1-8.5_
  
  - [x] 6.8 Write unit tests for network tools
    - Test all CRUD operations with mock HTTP client
    - Test error handling for invalid network IDs
    - _Requirements: 8.1-8.5_

- [ ] 7. Implement Network Resource sub-resource operations
  - [x] 7.1 Create tools/network_resources.go file
    - Define NetbirdNetworkResource struct with id, type, name, description, address, enabled, groups fields
    - _Requirements: 9.1-9.6_
  
  - [x] 7.2 Implement list_netbird_network_resources tool
    - Create ListNetbirdNetworkResourcesParams struct with network_id field
    - Implement listNetbirdNetworkResources function using GET /networks/{networkId}/resources
    - Register tool with MCP server
    - _Requirements: 9.1_
  
  - [x] 7.3 Implement create_netbird_network_resource tool
    - Create CreateNetbirdNetworkResourceParams struct with network_id, name, address, enabled, groups (required) and description (optional) fields
    - Implement createNetbirdNetworkResource function using POST /networks/{networkId}/resources
    - Register tool with MCP server
    - _Requirements: 9.2, 9.6_
  
  - [x] 7.4 Implement get_netbird_network_resource tool
    - Create GetNetbirdNetworkResourceParams struct with network_id and resource_id fields
    - Implement getNetbirdNetworkResource function using GET /networks/{networkId}/resources/{resourceId}
    - Register tool with MCP server
    - _Requirements: 9.3_
  
  - [x] 7.5 Implement update_netbird_network_resource tool
    - Create UpdateNetbirdNetworkResourceParams struct with network_id, resource_id (required) and all resource fields (optional)
    - Implement updateNetbirdNetworkResource function using PUT /networks/{networkId}/resources/{resourceId}
    - Register tool with MCP server
    - _Requirements: 9.4, 9.6_
  
  - [x] 7.6 Implement delete_netbird_network_resource tool
    - Create DeleteNetbirdNetworkResourceParams struct with network_id and resource_id fields
    - Implement deleteNetbirdNetworkResource function using DELETE /networks/{networkId}/resources/{resourceId}
    - Register tool with MCP server
    - _Requirements: 9.5_
  
  - [x] 7.7 Create AddNetbirdNetworkResourceTools function
    - Register all network resource tools with MCP server
    - _Requirements: 9.1-9.6_
  
  - [x] 7.8 Write unit tests for network resource tools
    - Test all CRUD operations with mock HTTP client
    - Test error handling for invalid network/resource IDs
    - _Requirements: 9.1-9.6_

- [x] 8. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Implement Network Router sub-resource operations
  - [x] 9.1 Create tools/network_routers.go file
    - Define NetbirdNetworkRouter struct with id, peer, peer_groups, metric, masquerade, enabled fields
    - _Requirements: 10.1-10.6_
  
  - [x] 9.2 Implement list_netbird_network_routers tool
    - Create ListNetbirdNetworkRoutersParams struct with network_id field
    - Implement listNetbirdNetworkRouters function using GET /networks/{networkId}/routers
    - Register tool with MCP server
    - _Requirements: 10.1_
  
  - [x] 9.3 Implement create_netbird_network_router tool
    - Create CreateNetbirdNetworkRouterParams struct with network_id, metric, masquerade, enabled (required) and peer, peer_groups (optional) fields
    - Implement createNetbirdNetworkRouter function using POST /networks/{networkId}/routers
    - Register tool with MCP server
    - _Requirements: 10.2, 10.6_
  
  - [x] 9.4 Implement get_netbird_network_router tool
    - Create GetNetbirdNetworkRouterParams struct with network_id and router_id fields
    - Implement getNetbirdNetworkRouter function using GET /networks/{networkId}/routers/{routerId}
    - Register tool with MCP server
    - _Requirements: 10.3_
  
  - [x] 9.5 Implement update_netbird_network_router tool
    - Create UpdateNetbirdNetworkRouterParams struct with network_id, router_id (required) and all router fields (optional)
    - Implement updateNetbirdNetworkRouter function using PUT /networks/{networkId}/routers/{routerId}
    - Register tool with MCP server
    - _Requirements: 10.4, 10.6_
  
  - [x] 9.6 Implement delete_netbird_network_router tool
    - Create DeleteNetbirdNetworkRouterParams struct with network_id and router_id fields
    - Implement deleteNetbirdNetworkRouter function using DELETE /networks/{networkId}/routers/{routerId}
    - Register tool with MCP server
    - _Requirements: 10.5_
  
  - [x] 9.7 Create AddNetbirdNetworkRouterTools function
    - Register all network router tools with MCP server
    - _Requirements: 10.1-10.6_
  
  - [x] 9.8 Write unit tests for network router tools
    - Test all CRUD operations with mock HTTP client
    - Test error handling for invalid network/router IDs
    - Test mutual exclusivity of peer and peer_groups fields
    - _Requirements: 10.1-10.6_

- [ ] 10. Add Setup Key allow_extra_dns_labels field
  - [x] 10.1 Update CreateNetbirdSetupKeyParams struct
    - Add AllowExtraDNSLabels field as pointer to bool with omitempty
    - _Requirements: 12.6, 18.1_
  
  - [x] 10.2 Update NetbirdSetupKey struct
    - Add AllowExtraDNSLabels field to response struct
    - _Requirements: 18.2_
  
  - [x] 10.3 Write unit tests for setup key extra DNS labels
    - Test JSON marshaling of allow_extra_dns_labels field
    - Test that field is omitted when nil
    - _Requirements: 12.6, 18.1, 18.2_

- [ ] 11. Implement Ingress Port Allocation operations (Cloud Only)
  - [x] 11.1 Update tools/ingress_ports.go file
    - Verify NetbirdPortAllocations struct has all required fields
    - _Requirements: 14.1-14.6_
  
  - [x] 11.2 Implement create_netbird_port_allocation tool
    - Create CreateNetbirdPortAllocationParams struct with peer_id, name, enabled (required) and port_ranges, direct_port (optional) fields
    - Implement createNetbirdPortAllocation function using POST /peers/{peerId}/ingress/ports
    - Register tool with MCP server
    - _Requirements: 14.2, 14.6_
  
  - [x] 11.3 Implement get_netbird_port_allocation tool
    - Create GetNetbirdPortAllocationParams struct with peer_id and allocation_id fields
    - Implement getNetbirdPortAllocation function using GET /peers/{peerId}/ingress/ports/{allocationId}
    - Register tool with MCP server
    - _Requirements: 14.3_
  
  - [x] 11.4 Implement update_netbird_port_allocation tool
    - Create UpdateNetbirdPortAllocationParams struct with peer_id, allocation_id (required) and all allocation fields (optional)
    - Implement updateNetbirdPortAllocation function using PUT /peers/{peerId}/ingress/ports/{allocationId}
    - Register tool with MCP server
    - _Requirements: 14.4, 14.6_
  
  - [x] 11.5 Implement delete_netbird_port_allocation tool
    - Create DeleteNetbirdPortAllocationParams struct with peer_id and allocation_id fields
    - Implement deleteNetbirdPortAllocation function using DELETE /peers/{peerId}/ingress/ports/{allocationId}
    - Register tool with MCP server
    - _Requirements: 14.5_
  
  - [x] 11.6 Update AddNetbirdPortAllocationTools function
    - Register all new port allocation tools with MCP server
    - _Requirements: 14.1-14.6_
  
  - [x] 11.7 Write unit tests for port allocation tools
    - Test all CRUD operations with mock HTTP client
    - Test error handling for invalid peer/allocation IDs
    - _Requirements: 14.1-14.6_

- [x] 12. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 13. Update tool registration in main.go
  - [x] 13.1 Add network tools registration
    - Import tools package
    - Call tools.AddNetbirdNetworkTools(mcp) in main function
    - _Requirements: 8.1-8.5_
  
  - [x] 13.2 Add network resource tools registration
    - Call tools.AddNetbirdNetworkResourceTools(mcp) in main function
    - _Requirements: 9.1-9.6_
  
  - [x] 13.3 Add network router tools registration
    - Call tools.AddNetbirdNetworkRouterTools(mcp) in main function
    - _Requirements: 10.1-10.6_
  
  - [x] 13.4 Add port allocation tools registration
    - Call tools.AddNetbirdPortAllocationTools(mcp) in main function
    - _Requirements: 14.1-14.6_

- [ ] 14. Verify backward compatibility
  - [x] 14.1 Run existing test suite
    - Execute all existing unit tests
    - Verify no tests fail due to data structure changes
    - _Requirements: 17.1, 17.2, 17.3, 17.4_
  
  - [x] 14.2 Test existing tool functionality
    - Manually test existing tools (account, peer, group, policy, etc.)
    - Verify they still work with updated data structures
    - _Requirements: 17.1, 17.2, 17.3, 17.4_
  
  - [x] 14.3 Write backward compatibility tests
    - Test that old JSON payloads can still be unmarshaled
    - Test that new fields don't break existing code
    - _Requirements: 17.1, 17.2, 17.3, 17.4_

- [x] 15. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- All new tools follow the existing pattern for consistency
- Backward compatibility is maintained throughout
