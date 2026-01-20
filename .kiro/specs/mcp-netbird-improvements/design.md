# Design Document

## Overview

This design addresses critical issues discovered during live testing of the mcp-netbird MCP server, specifically focusing on policy rules management and group consolidation workflows. The primary issue is a format mismatch in the NetBird API: when creating/updating policies, the `sources` and `destinations` fields in rules must be string arrays (group IDs), but the API returns them as object arrays (full group objects). This design provides fixes for policy operations and adds helper functions for common administrative tasks like group consolidation.

## Architecture

The solution follows the existing MCP server architecture with enhancements to three main areas:

1. **Policy Management Layer**: Fix the rules parameter handling in create/update operations
2. **Helper Functions Layer**: Add utility functions for group dependency discovery and bulk operations
3. **Validation Layer**: Add pre-flight validation for policy rules to catch errors before API calls

```
┌─────────────────────────────────────────────────────────┐
│                    MCP Client                            │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              MCP Server Tool Layer                       │
│  ┌──────────────────┐  ┌──────────────────────────┐    │
│  │ Policy Tools     │  │ Group Tools              │    │
│  │ - create_policy  │  │ - delete_group (force)   │    │
│  │ - update_policy  │  │ - list_policies_by_group │    │
│  │ - validate_rules │  │ - replace_group          │    │
│  └──────────────────┘  └──────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              NetBird REST API                            │
└─────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Policy Rule Formatter

**Purpose**: Convert between MCP input format and NetBird API format

**Key Issue**: The NetBird API has an asymmetric format:
- **Request (POST/PUT)**: `sources` and `destinations` are `string[]` (group IDs only)
- **Response (GET)**: `sources` and `destinations` are `object[]` (full group objects)

**Interface**:
```go
// FormatRuleForAPI converts a rule from MCP format to NetBird API request format
func FormatRuleForAPI(rule map[string]interface{}) (map[string]interface{}, error)

// Input format (from MCP client):
// {
//   "sources": ["group-id-1", "group-id-2"],  // Can be strings or objects
//   "destinations": ["group-id-3"],
//   "protocol": "tcp",
//   "action": "accept",
//   ...
// }
//
// Output format (for NetBird API):
// {
//   "sources": ["group-id-1", "group-id-2"],  // Always strings
//   "destinations": ["group-id-3"],
//   "protocol": "tcp",
//   "action": "accept",
//   ...
// }
```

**Implementation Strategy**:
- Check if `sources`/`destinations` are already string arrays
- If they're object arrays, extract the `id` field from each object
- Preserve all other rule fields unchanged
- Handle nil/empty arrays gracefully

### 2. Policy Rule Validator

**Purpose**: Validate policy rules before sending to API to provide early error detection

**Interface**:
```go
// ValidatePolicyRules validates an array of policy rules
func ValidatePolicyRules(rules []map[string]interface{}) error

// Validation checks:
// - Required fields: name, enabled, action, bidirectional, protocol
// - Valid action values: "accept", "drop"
// - Valid protocol values: "tcp", "udp", "icmp", "all"
// - Port ranges: start <= end
// - At least one of: sources, sourceResource
// - At least one of: destinations, destinationResource
```

**Error Messages**:
- Clear indication of which rule failed (by index or name)
- Specific field that caused the validation failure
- Expected format or valid values

### 3. Group Dependency Analyzer

**Purpose**: Identify all policies that reference a specific group

**Interface**:
```go
// ListPoliciesByGroup returns all policies referencing a group
func ListPoliciesByGroup(ctx context.Context, groupID string) ([]PolicyReference, error)

type PolicyReference struct {
    PolicyID    string
    PolicyName  string
    RuleID      string
    RuleName    string
    Location    string  // "sources", "destinations", or "authorized_groups"
}
```

**Algorithm**:
1. Fetch all policies from NetBird API
2. For each policy, iterate through all rules
3. Check if group ID appears in:
   - `sources` array
   - `destinations` array
   - `authorized_groups` map keys
4. Collect all matches with full context

### 4. Group Replacement Engine

**Purpose**: Replace one group with another across all policies

**Interface**:
```go
// ReplaceGroupInPolicies replaces oldGroupID with newGroupID in all policies
func ReplaceGroupInPolicies(ctx context.Context, oldGroupID, newGroupID string) ([]string, error)

// Returns: list of policy IDs that were updated
```

**Algorithm**:
1. Use `ListPoliciesByGroup` to find all affected policies
2. For each policy:
   a. Fetch current policy configuration
   b. Iterate through rules and replace group ID in:
      - `sources` arrays
      - `destinations` arrays
      - `authorized_groups` map keys
   c. Update policy via PUT request
3. Collect and return list of updated policy IDs

**Error Handling**:
- If any update fails, continue with remaining policies
- Return partial success with list of updated policies and error details
- Log which policies failed and why

### 5. Force Delete Group

**Purpose**: Delete a group by first removing it from all dependent policies

**Interface**:
```go
// DeleteGroupForce deletes a group after removing dependencies
func DeleteGroupForce(ctx context.Context, groupID string) (*ForceDeleteResult, error)

type ForceDeleteResult struct {
    GroupID          string
    PoliciesModified []string
    Deleted          bool
    Errors           []string
}
```

**Algorithm**:
1. Call `ListPoliciesByGroup` to find dependencies
2. For each dependent policy:
   a. Fetch current policy
   b. Remove group ID from all rules (sources, destinations, authorized_groups)
   c. If a rule becomes invalid (no sources or destinations), remove the rule
   d. If policy becomes empty (no rules), delete the policy
   e. Otherwise, update the policy
3. After all dependencies resolved, delete the group
4. Return summary of operations

## Data Models

### Policy Rule Structure (NetBird API)

```go
type PolicyRule struct {
    ID                  string                 `json:"id,omitempty"`
    Name                string                 `json:"name"`
    Description         string                 `json:"description,omitempty"`
    Enabled             bool                   `json:"enabled"`
    Action              string                 `json:"action"`              // "accept" or "drop"
    Bidirectional       bool                   `json:"bidirectional"`
    Protocol            string                 `json:"protocol"`            // "tcp", "udp", "icmp", "all"
    Ports               []string               `json:"ports,omitempty"`
    PortRanges          []PortRange            `json:"port_ranges,omitempty"`
    AuthorizedGroups    map[string][]string    `json:"authorized_groups,omitempty"`
    Sources             []string               `json:"sources,omitempty"`             // Group IDs (request)
    SourceResource      *ResourceReference     `json:"sourceResource,omitempty"`
    Destinations        []string               `json:"destinations,omitempty"`        // Group IDs (request)
    DestinationResource *ResourceReference     `json:"destinationResource,omitempty"`
}

type PortRange struct {
    Start int `json:"start"`
    End   int `json:"end"`
}

type ResourceReference struct {
    ID   string `json:"id"`
    Type string `json:"type"`  // "host", "domain", "subnet"
}
```

### Policy Template

```go
func GetPolicyTemplate() map[string]interface{} {
    return map[string]interface{}{
        "name":        "example-policy",
        "description": "Example policy with rules",
        "enabled":     true,
        "rules": []map[string]interface{}{
            {
                "name":          "allow-web-traffic",
                "description":   "Allow HTTP/HTTPS from dev group to prod group",
                "enabled":       true,
                "action":        "accept",
                "bidirectional": false,
                "protocol":      "tcp",
                "ports":         []string{"80", "443"},
                "sources":       []string{"group-id-dev"},
                "destinations":  []string{"group-id-prod"},
            },
            {
                "name":          "allow-ssh-with-auth",
                "description":   "Allow SSH with user authorization",
                "enabled":       true,
                "action":        "accept",
                "bidirectional": false,
                "protocol":      "tcp",
                "port_ranges": []map[string]interface{}{
                    {"start": 22, "end": 22},
                },
                "sources":      []string{"group-id-admins"},
                "destinations": []string{"group-id-servers"},
                "authorized_groups": map[string][]string{
                    "group-id-admins": {"user1", "user2"},
                },
            },
        },
    }
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: Rule Formatting Preserves Structure

*For any* policy rule with valid fields, formatting the rule for the NetBird API should preserve all fields except converting sources/destinations from object arrays to string arrays (extracting IDs).

**Validates: Requirements 1.1, 1.4, 2.1**

### Property 2: Invalid Rules Rejected Before API Call

*For any* policy rule missing required fields or containing invalid values, validation should return a descriptive error before making any API call.

**Validates: Requirements 1.2, 2.3, 3.1, 3.2, 3.4, 3.5**

### Property 3: Port Range Invariant

*For any* port range in a policy rule, the start port must be less than or equal to the end port.

**Validates: Requirements 3.3**

### Property 4: Group Dependency Discovery Completeness

*For any* group ID, querying policies by that group should return all policies where the group appears in sources, destinations, or authorized_groups of any rule.

**Validates: Requirements 4.1, 4.2, 4.3, 4.5**

### Property 5: Partial Policy Update Preservation

*For any* policy update that specifies only a subset of fields, all unspecified fields should remain unchanged after the update.

**Validates: Requirements 2.2**

### Property 6: Group Replacement Completeness

*For any* old group ID and new group ID, replacing the old group across all policies should result in the new group ID appearing in all locations where the old group ID previously appeared (sources, destinations, authorized_groups).

**Validates: Requirements 6.2, 6.3, 6.4, 6.5**

### Property 7: Force Delete Removes Dependencies

*For any* group with dependent policies, force deleting the group should result in the group being removed from all policy rules before the group is deleted.

**Validates: Requirements 5.2, 5.5**

## Error Handling

### Error Categories

1. **Validation Errors**: Caught before API calls
   - Missing required fields
   - Invalid field values (protocol, action, etc.)
   - Invalid port ranges (start > end)
   - Malformed group IDs
   - Example: `"validation error: rule 'allow-web' missing required field 'action'"`

2. **API Errors**: Returned from NetBird API
   - 400 Bad Request: Malformed request body
   - 401 Unauthorized: Invalid or expired token
   - 404 Not Found: Policy or group doesn't exist
   - 409 Conflict: Resource already exists or has dependencies
   - 500 Internal Server Error: NetBird server error

3. **Dependency Errors**: Group deletion conflicts
   - Group referenced in policies (when force=false)
   - Example: `"cannot delete group 'dev-team': referenced by 3 policies [policy-1, policy-2, policy-3]"`

4. **Partial Failure Errors**: Bulk operations
   - Some policies updated successfully, others failed
   - Return list of successes and failures
   - Example: `"updated 5 of 7 policies; failures: [policy-6: 404 not found, policy-7: 400 invalid rule]"`

### Error Response Format

```go
type ErrorResponse struct {
    Error   string            `json:"error"`
    Code    string            `json:"code"`
    Details map[string]string `json:"details,omitempty"`
}

// Examples:
// Validation error:
// {
//   "error": "validation failed",
//   "code": "VALIDATION_ERROR",
//   "details": {
//     "rule": "allow-web",
//     "field": "action",
//     "message": "missing required field"
//   }
// }

// Dependency error:
// {
//   "error": "group has dependencies",
//   "code": "DEPENDENCY_ERROR",
//   "details": {
//     "group_id": "dev-team",
//     "policy_count": "3",
//     "policies": "policy-1,policy-2,policy-3"
//   }
// }
```

### Retry Strategy

- **Transient Errors (5xx, network timeouts)**: Retry with exponential backoff
  - Max retries: 3
  - Initial delay: 1 second
  - Backoff multiplier: 2x
  
- **Client Errors (4xx)**: Do not retry
  - These indicate problems with the request that won't be fixed by retrying
  - Exception: 429 Rate Limit - retry after delay specified in response

- **Bulk Operations**: Continue on failure
  - Don't stop entire operation if one item fails
  - Collect all errors and return summary

## Testing Strategy

### Unit Testing

Unit tests focus on specific functions and edge cases:

1. **Rule Formatting Tests**
   - Test with sources/destinations as strings
   - Test with sources/destinations as objects (extract IDs)
   - Test with mixed formats
   - Test with empty arrays
   - Test with nil values
   - Test preservation of all other fields

2. **Validation Tests**
   - Test each required field validation
   - Test protocol enum validation
   - Test action enum validation
   - Test port range validation (start <= end)
   - Test error message format and content

3. **Group Search Tests**
   - Test finding groups in sources
   - Test finding groups in destinations
   - Test finding groups in authorized_groups
   - Test with no matches (empty result)
   - Test with multiple matches in same policy

4. **Group Replacement Tests**
   - Test replacement in sources
   - Test replacement in destinations
   - Test replacement in authorized_groups
   - Test with no matches (no updates)
   - Test with multiple occurrences in same rule

### Property-Based Testing

Property tests verify universal correctness across randomized inputs. Each test should run minimum 100 iterations.

1. **Property Test: Rule Formatting Idempotence**
   - Generate random valid rules
   - Format for API
   - Verify all fields preserved except sources/destinations converted to strings
   - Tag: `Feature: mcp-netbird-improvements, Property 1: Rule Formatting Preserves Structure`

2. **Property Test: Validation Rejects Invalid Rules**
   - Generate random invalid rules (missing fields, bad values)
   - Verify validation returns error
   - Verify error message contains rule and field info
   - Tag: `Feature: mcp-netbird-improvements, Property 2: Invalid Rules Rejected Before API Call`

3. **Property Test: Port Range Invariant**
   - Generate random port ranges
   - Verify validation enforces start <= end
   - Tag: `Feature: mcp-netbird-improvements, Property 3: Port Range Invariant`

4. **Property Test: Group Discovery Completeness**
   - Generate random policies with groups in various locations
   - Search for each group
   - Verify all occurrences found
   - Tag: `Feature: mcp-netbird-improvements, Property 4: Group Dependency Discovery Completeness`

5. **Property Test: Partial Update Preservation**
   - Generate random policy
   - Generate random partial update
   - Verify unspecified fields unchanged
   - Tag: `Feature: mcp-netbird-improvements, Property 5: Partial Policy Update Preservation`

6. **Property Test: Group Replacement Completeness**
   - Generate random policies with group references
   - Replace group
   - Verify all occurrences replaced
   - Verify no old group IDs remain
   - Tag: `Feature: mcp-netbird-improvements, Property 6: Group Replacement Completeness`

7. **Property Test: Force Delete Cleanup**
   - Generate random policies with group dependencies
   - Force delete group
   - Verify group removed from all policies
   - Verify group deleted
   - Tag: `Feature: mcp-netbird-improvements, Property 7: Force Delete Removes Dependencies`

### Integration Testing

Integration tests verify the MCP server works correctly with the actual NetBird API:

1. **Policy Creation with Rules**
   - Create policy with simple rule (single source/destination)
   - Create policy with complex rule (multiple sources, port ranges, authorized groups)
   - Verify policies created successfully
   - Clean up: delete created policies

2. **Policy Update with Rules**
   - Create policy
   - Update with new rules
   - Verify update successful
   - Verify unchanged fields preserved
   - Clean up: delete policy

3. **Group Force Delete**
   - Create test groups
   - Create policies referencing groups
   - Force delete group
   - Verify policies updated (group removed)
   - Verify group deleted
   - Clean up: delete policies

4. **Group Replacement**
   - Create test groups (old and new)
   - Create policies referencing old group
   - Replace old group with new group
   - Verify all policies updated
   - Clean up: delete policies and groups

### Test Configuration

- **Property Test Library**: Use `testing/quick` (Go standard library) or `gopter` for more advanced features
- **Iterations**: Minimum 100 per property test
- **Test Fixtures**: Include real-world policy configurations from live testing
- **Cleanup**: All integration tests must clean up created resources
- **Isolation**: Each test should be independent and not rely on other tests

### Test Data Generators

For property-based testing, implement generators for:

```go
// Generate random valid policy rules
func GenerateValidRule() map[string]interface{}

// Generate random invalid policy rules (for validation testing)
func GenerateInvalidRule() map[string]interface{}

// Generate random group IDs
func GenerateGroupID() string

// Generate random policies with group references
func GeneratePolicyWithGroups(groupIDs []string) map[string]interface{}
```
