# Requirements Document

## Introduction

This document specifies requirements for improving the mcp-netbird MCP server based on live testing feedback. The improvements focus on fixing policy rules management, adding group consolidation workflows, and providing helper functions for common administrative tasks. Live testing revealed that while core functionality (networks, resources, routers, peers) works excellently, policy rules API and group deletion with dependencies need fixes and enhancements.

## Glossary

- **MCP_Server**: The Model Context Protocol server that provides programmatic access to NetBird API
- **Policy**: A NetBird access control policy that defines network access rules between groups
- **Policy_Rule**: A component of a policy that specifies source, destination, protocol, ports, and action
- **Group**: A NetBird group containing peers or resources
- **NetBird_API**: The REST API provided by NetBird for network management
- **Dependent_Policy**: A policy that references a specific group in its rules or configuration
- **Force_Delete**: An operation that removes dependencies before deleting a resource

## Requirements

### Requirement 1: Policy Rules Creation

**User Story:** As a NetBird administrator, I want to create policies with rules through the MCP server, so that I can programmatically configure access control policies.

#### Acceptance Criteria

1. WHEN a user creates a policy with rules parameter, THE MCP_Server SHALL format the rules according to NetBird API specification
2. WHEN the rules parameter contains invalid structure, THE MCP_Server SHALL return a descriptive error before making the API call
3. WHEN a policy is created with valid rules, THE NetBird_API SHALL accept the request and create the policy
4. THE MCP_Server SHALL support all rule fields including sources, destinations, protocol, port_ranges, action, bidirectional, and enabled

### Requirement 2: Policy Rules Updates

**User Story:** As a NetBird administrator, I want to update existing policies with new rules, so that I can modify access control configurations programmatically.

#### Acceptance Criteria

1. WHEN a user updates a policy with rules parameter, THE MCP_Server SHALL format the rules according to NetBird API specification
2. WHEN updating a policy, THE MCP_Server SHALL preserve existing policy properties not specified in the update
3. WHEN the rules parameter contains invalid structure, THE MCP_Server SHALL return a descriptive error before making the API call
4. WHEN a policy update includes valid rules, THE NetBird_API SHALL accept the request and update the policy

### Requirement 3: Policy Rules Validation

**User Story:** As a NetBird administrator, I want the MCP server to validate policy rules before sending to the API, so that I can catch configuration errors early.

#### Acceptance Criteria

1. WHEN validating policy rules, THE MCP_Server SHALL verify that each rule contains required fields
2. WHEN a rule references groups, THE MCP_Server SHALL verify the group references are properly formatted
3. WHEN port ranges are specified, THE MCP_Server SHALL verify that start port is less than or equal to end port
4. WHEN protocol is specified, THE MCP_Server SHALL verify it is a valid protocol value
5. WHEN validation fails, THE MCP_Server SHALL return a descriptive error message indicating which rule and field failed validation

### Requirement 4: Group Dependency Discovery

**User Story:** As a NetBird administrator, I want to list all policies that reference a specific group, so that I can understand dependencies before making changes.

#### Acceptance Criteria

1. WHEN querying policies by group, THE MCP_Server SHALL return all policies where the group appears in sources
2. WHEN querying policies by group, THE MCP_Server SHALL return all policies where the group appears in destinations
3. WHEN querying policies by group, THE MCP_Server SHALL return all policies where the group appears in authorized_groups
4. WHEN a group is not referenced by any policy, THE MCP_Server SHALL return an empty list
5. THE MCP_Server SHALL provide the policy ID, name, and rule details for each dependent policy

### Requirement 5: Group Force Delete

**User Story:** As a NetBird administrator, I want to force delete a group by automatically updating dependent policies, so that I can consolidate groups without manual policy updates.

#### Acceptance Criteria

1. WHEN force deleting a group, THE MCP_Server SHALL first identify all dependent policies
2. WHEN dependent policies exist, THE MCP_Server SHALL remove the group references from those policies before deletion
3. WHEN all dependencies are resolved, THE MCP_Server SHALL delete the group
4. WHEN force delete is not specified and dependencies exist, THE MCP_Server SHALL return an error listing the dependent policies
5. THE MCP_Server SHALL return a summary of policies modified and the deletion status

### Requirement 6: Group Replacement in Policies

**User Story:** As a NetBird administrator, I want to replace one group with another across all policies, so that I can consolidate groups efficiently.

#### Acceptance Criteria

1. WHEN replacing a group in policies, THE MCP_Server SHALL identify all policies referencing the old group
2. WHEN updating policies, THE MCP_Server SHALL replace the old group ID with the new group ID in sources
3. WHEN updating policies, THE MCP_Server SHALL replace the old group ID with the new group ID in destinations
4. WHEN updating policies, THE MCP_Server SHALL replace the old group ID with the new group ID in authorized_groups
5. THE MCP_Server SHALL return a list of policies that were updated with the replacement
6. WHEN the old group ID is not found in any policy, THE MCP_Server SHALL return an empty list of updates

### Requirement 7: Policy Template Provision

**User Story:** As a NetBird administrator, I want to get an example policy structure, so that I can understand the correct format for creating policies with rules.

#### Acceptance Criteria

1. WHEN requesting a policy template, THE MCP_Server SHALL return a complete example policy structure
2. THE template SHALL include examples of all supported rule fields
3. THE template SHALL include comments or descriptions explaining each field
4. THE template SHALL demonstrate both simple and complex rule configurations

### Requirement 8: Integration Testing

**User Story:** As a developer, I want comprehensive integration tests for policy operations, so that I can verify the MCP server works correctly with the NetBird API.

#### Acceptance Criteria

1. THE test suite SHALL include tests for creating policies with various rule configurations
2. THE test suite SHALL include tests for updating policies with rule modifications
3. THE test suite SHALL include tests for group deletion with dependencies
4. THE test suite SHALL include tests for group replacement across policies
5. THE test suite SHALL use test fixtures from real-world configurations
6. WHEN tests run against a live NetBird instance, THE tests SHALL clean up created resources

### Requirement 9: Error Handling Documentation

**User Story:** As a developer using the MCP server, I want clear documentation of error handling patterns, so that I can handle failures appropriately in my code.

#### Acceptance Criteria

1. THE documentation SHALL describe all error types returned by policy operations
2. THE documentation SHALL provide examples of error responses for common failure scenarios
3. THE documentation SHALL explain how to interpret validation errors
4. THE documentation SHALL document retry strategies for transient failures
5. THE documentation SHALL include examples of error handling in client code
