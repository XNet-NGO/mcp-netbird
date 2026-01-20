# Implementation Plan: MCP NetBird Improvements

## Overview

This implementation plan addresses critical issues in the mcp-netbird MCP server discovered during live testing. The main focus is fixing policy rules management (format mismatch between request and response) and adding group consolidation workflows. The implementation will be done in Go, following the existing MCP server architecture.

## Tasks

- [x] 1. Implement policy rule formatting and validation
  - [x] 1.1 Create FormatRuleForAPI function in tools/policies.go
    - Convert sources/destinations from object arrays to string arrays
    - Extract ID field from objects if present
    - Preserve all other rule fields unchanged
    - Handle nil and empty arrays gracefully
    - _Requirements: 1.1, 1.4, 2.1_

  - [x] 1.2 Write property test for rule formatting
    - **Property 1: Rule Formatting Preserves Structure**
    - **Validates: Requirements 1.1, 1.4, 2.1**
    - Generate random valid rules with various field combinations
    - Verify all fields preserved except sources/destinations converted to strings
    - Test with 100+ iterations

  - [x] 1.3 Create ValidatePolicyRules function in tools/policies.go
    - Validate required fields: name, enabled, action, bidirectional, protocol
    - Validate action enum: "accept" or "drop"
    - Validate protocol enum: "tcp", "udp", "icmp", "all"
    - Validate port ranges: start <= end
    - Validate at least one source (sources or sourceResource)
    - Validate at least one destination (destinations or destinationResource)
    - Return descriptive errors with rule name/index and field name
    - _Requirements: 1.2, 2.3, 3.1, 3.2, 3.3, 3.4, 3.5_

  - [x] 1.4 Write property tests for validation
    - **Property 2: Invalid Rules Rejected Before API Call**
    - **Validates: Requirements 1.2, 2.3, 3.1, 3.2, 3.4, 3.5**
    - Generate random invalid rules (missing fields, bad values)
    - Verify validation returns error before API call
    - Verify error messages contain rule and field information
    - **Property 3: Port Range Invariant**
    - **Validates: Requirements 3.3**
    - Generate random port ranges
    - Verify validation enforces start <= end constraint

  - [x] 1.5 Update create_netbird_policy tool to use formatting and validation
    - Call ValidatePolicyRules before API call
    - Call FormatRuleForAPI for each rule
    - Return validation errors to user
    - _Requirements: 1.1, 1.2, 1.3_

  - [x] 1.6 Update update_netbird_policy tool to use formatting and validation
    - Call ValidatePolicyRules before API call
    - Call FormatRuleForAPI for each rule
    - Preserve existing policy fields not in update
    - Return validation errors to user
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

  - [x] 1.7 Write unit tests for policy create/update
    - Test create with simple rules
    - Test create with complex rules (port ranges, authorized groups)
    - Test update with partial data
    - Test validation error handling
    - Test API error handling

- [x] 2. Checkpoint - Ensure policy formatting tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 3. Implement group dependency discovery
  - [x] 3.1 Create ListPoliciesByGroup function in tools/groups.go
    - Fetch all policies from NetBird API
    - Iterate through all rules in each policy
    - Check if group ID appears in sources array
    - Check if group ID appears in destinations array
    - Check if group ID appears in authorized_groups map keys
    - Return PolicyReference structs with policy ID, name, rule ID, name, and location
    - _Requirements: 4.1, 4.2, 4.3, 4.5_

  - [x] 3.2 Write property test for group discovery
    - **Property 4: Group Dependency Discovery Completeness**
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.5**
    - Generate random policies with groups in various locations
    - Search for each group
    - Verify all occurrences found in sources, destinations, and authorized_groups
    - Test edge case: group not referenced (empty result)

  - [x] 3.3 Write unit tests for group discovery
    - Test finding group in sources only
    - Test finding group in destinations only
    - Test finding group in authorized_groups only
    - Test finding group in multiple locations
    - Test with no matches (empty result)
    - Test with multiple policies referencing same group

- [x] 4. Implement group replacement functionality
  - [x] 4.1 Create ReplaceGroupInPolicies function in tools/groups.go
    - Use ListPoliciesByGroup to find affected policies
    - For each policy, fetch current configuration
    - Iterate through rules and replace group ID in sources arrays
    - Iterate through rules and replace group ID in destinations arrays
    - Iterate through rules and replace group ID in authorized_groups keys
    - Update each policy via PUT request
    - Collect list of updated policy IDs
    - Handle partial failures (continue on error, return summary)
    - _Requirements: 6.2, 6.3, 6.4, 6.5_

  - [x] 4.2 Write property test for group replacement
    - **Property 6: Group Replacement Completeness**
    - **Validates: Requirements 6.2, 6.3, 6.4, 6.5**
    - Generate random policies with group references
    - Replace old group with new group
    - Verify all occurrences replaced in sources, destinations, authorized_groups
    - Verify no old group IDs remain
    - Test edge case: old group not found (empty update list)

  - [x] 4.3 Write unit tests for group replacement
    - Test replacement in sources
    - Test replacement in destinations
    - Test replacement in authorized_groups
    - Test with no matches (no updates)
    - Test with multiple occurrences in same rule
    - Test partial failure handling

- [x] 5. Implement force delete group functionality
  - [x] 5.1 Add force parameter to delete_netbird_group tool
    - Add boolean force parameter to tool definition
    - Update tool description to explain force behavior
    - _Requirements: 5.4_

  - [x] 5.2 Create DeleteGroupForce function in tools/groups.go
    - Call ListPoliciesByGroup to find dependencies
    - For each dependent policy, fetch current configuration
    - Remove group ID from all rules (sources, destinations, authorized_groups)
    - If rule becomes invalid (no sources or destinations), remove the rule
    - If policy becomes empty (no rules), delete the policy
    - Otherwise, update the policy
    - After all dependencies resolved, delete the group
    - Return ForceDeleteResult with modified policies, deletion status, and errors
    - _Requirements: 5.1, 5.2, 5.3, 5.5_

  - [x] 5.3 Update delete_netbird_group tool to support force delete
    - If force=false and dependencies exist, return error with policy list
    - If force=true, call DeleteGroupForce
    - Return summary of operations
    - _Requirements: 5.4, 5.5_

  - [x] 5.4 Write property test for force delete
    - **Property 7: Force Delete Removes Dependencies**
    - **Validates: Requirements 5.2, 5.5**
    - Generate random policies with group dependencies
    - Force delete group
    - Verify group removed from all policies
    - Verify group deleted successfully

  - [x] 5.5 Write unit tests for force delete
    - Test force delete with dependencies
    - Test normal delete with dependencies (should fail)
    - Test delete without dependencies
    - Test cleanup of invalid rules
    - Test cleanup of empty policies

- [x] 6. Checkpoint - Ensure group management tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. Add helper functions and templates
  - [x] 7.1 Create GetPolicyTemplate function in tools/policies.go
    - Return example policy structure with simple rule
    - Return example policy structure with complex rule (port ranges, authorized groups)
    - Include comments explaining each field
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [x] 7.2 Write unit test for policy template
    - Verify template contains all required fields
    - Verify template includes simple and complex examples
    - Verify template is valid (passes validation)

  - [x] 7.3 Expose helper functions as MCP tools
    - Add list_policies_by_group tool
    - Add replace_group_in_policies tool
    - Add get_policy_template tool
    - Update tool descriptions and parameter schemas
    - _Requirements: 4.1, 6.1, 7.1_

- [ ] 8. Add integration tests
  - [ ] 8.1 Write integration test for policy creation with rules
    - Create policy with simple rule
    - Create policy with complex rule
    - Verify policies created successfully
    - Clean up created policies

  - [ ] 8.2 Write integration test for policy update with rules
    - Create policy
    - Update with new rules
    - Verify update successful and unchanged fields preserved
    - Clean up policy

  - [ ] 8.3 Write integration test for group force delete
    - Create test groups
    - Create policies referencing groups
    - Force delete group
    - Verify policies updated and group deleted
    - Clean up policies

  - [ ] 8.4 Write integration test for group replacement
    - Create test groups (old and new)
    - Create policies referencing old group
    - Replace old group with new group
    - Verify all policies updated
    - Clean up policies and groups

- [x] 9. Update documentation
  - [x] 9.1 Document policy rule format in README
    - Document correct format for sources/destinations (string arrays)
    - Document all supported rule fields
    - Provide examples of simple and complex rules
    - _Requirements: 9.1, 9.2_

  - [x] 9.2 Document error handling patterns
    - Document validation error format
    - Document API error handling
    - Document dependency error format
    - Document retry strategy for transient errors
    - Provide error handling examples
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

  - [x] 9.3 Document helper functions
    - Document list_policies_by_group usage
    - Document replace_group_in_policies usage
    - Document get_policy_template usage
    - Provide examples for common workflows (group consolidation)
    - _Requirements: 4.1, 6.1, 7.1_

- [x] 10. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties with 100+ iterations
- Unit tests validate specific examples and edge cases
- Integration tests verify functionality against live NetBird API
- All integration tests must clean up created resources
