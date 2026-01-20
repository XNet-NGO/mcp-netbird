# MCP-NetBird Error Handling Guide

This guide explains the error handling patterns used in MCP-NetBird and how to handle common errors.

## Table of Contents

1. [Error Types](#error-types)
2. [Validation Errors](#validation-errors)
3. [API Errors](#api-errors)
4. [Dependency Errors](#dependency-errors)
5. [Retry Strategies](#retry-strategies)
6. [Common Error Scenarios](#common-error-scenarios)

---

## Error Types

MCP-NetBird uses three main categories of errors:

### 1. Validation Errors

Errors that occur before making API calls, when validating input parameters.

**Format:**
```
validation error: [field_name]: [error_description]
```

**Example:**
```
validation error: rule 0 (SSH Rule): action must be 'accept' or 'drop', got 'allow'
```

### 2. API Errors

Errors returned by the NetBird API.

**Format:**
```
unexpected status code: [code], body: [response_body]
```

**Example:**
```
unexpected status code: 422, body: {"message":"nameserver group primary status is false and domains are empty","code":422}
```

### 3. Dependency Errors

Errors that occur when trying to delete resources that have dependencies.

**Format:**
```
cannot delete group: [count] policies depend on it: [policy_list]
```

**Example:**
```
cannot delete group: 3 policies depend on it: Admin-Access, Team-Access, Exit-Node
```

---

## Validation Errors

### Policy Rule Validation

Policy rules are validated before being sent to the API. Validation checks:

#### Required Fields
- `name` - Rule name must be present
- `enabled` - Must be boolean (true/false)
- `action` - Must be "accept" or "drop"
- `bidirectional` - Must be boolean
- `protocol` - Must be "tcp", "udp", "icmp", or "all"

#### Source/Destination Requirements
- At least one source required: `sources` array OR `sourceResource` object
- At least one destination required: `destinations` array OR `destinationResource` object

#### Port Range Validation
- Only valid for TCP and UDP protocols
- Start port must be <= end port
- Ports must be in range 1-65535

### Validation Error Examples

**Missing Required Field:**
```json
{
  "error": "validation error: rule 0 (My Rule): name is required"
}
```

**Invalid Action:**
```json
{
  "error": "validation error: rule 1 (Access Rule): action must be 'accept' or 'drop', got 'allow'"
}
```

**Invalid Protocol:**
```json
{
  "error": "validation error: rule 0 (Network Rule): protocol must be 'tcp', 'udp', 'icmp', or 'all', got 'http'"
}
```

**Invalid Port Range:**
```json
{
  "error": "validation error: rule 2 (Port Rule): port range start (443) must be <= end (80)"
}
```

**Missing Source:**
```json
{
  "error": "validation error: rule 0 (Rule): at least one source (sources or sourceResource) is required"
}
```

**Missing Destination:**
```json
{
  "error": "validation error: rule 1 (Rule): at least one destination (destinations or destinationResource) is required"
}
```

### Handling Validation Errors

Validation errors indicate problems with the input data. To fix:

1. Check the error message for the specific field and rule
2. Verify the field value matches the required format
3. Correct the input and retry

**Example Fix:**
```go
// Before (invalid)
rule := NetbirdPolicyRule{
    Name:          "My Rule",
    Enabled:       true,
    Action:        "allow",  // ❌ Invalid - should be "accept" or "drop"
    Bidirectional: true,
    Protocol:      "tcp",
    Sources:       []string{"group-id"},
    Destinations:  []string{"group-id"},
}

// After (valid)
rule := NetbirdPolicyRule{
    Name:          "My Rule",
    Enabled:       true,
    Action:        "accept",  // ✅ Valid
    Bidirectional: true,
    Protocol:      "tcp",
    Sources:       []string{"group-id"},
    Destinations:  []string{"group-id"},
}
```

---

## API Errors

### Common HTTP Status Codes

| Code | Meaning | Common Causes |
|------|---------|---------------|
| 400 | Bad Request | Invalid JSON, missing required fields |
| 401 | Unauthorized | Invalid or missing API token |
| 404 | Not Found | Resource doesn't exist, invalid ID |
| 422 | Unprocessable Entity | Business logic validation failed |
| 500 | Internal Server Error | Server-side error |

### API Error Examples

**401 Unauthorized:**
```
unexpected status code: 401, body: {"message":"invalid token","code":401}
```

**Fix:** Check that `NETBIRD_API_TOKEN` environment variable is set correctly.

**404 Not Found:**
```
unexpected status code: 404, body: 404 page not found
```

**Fix:** Verify the resource ID exists. For port allocations, this may indicate the feature is not available.

**422 Unprocessable Entity:**
```
unexpected status code: 422, body: {"message":"the list of nameservers should be 1 or 3, got 4","code":422}
```

**Fix:** Follow the API's business rules (in this case, use 1 or 3 nameservers, not 4).

**500 Internal Server Error:**
```
unexpected status code: 500, body: {"message":"internal server error","code":500}
```

**Fix:** This is a server-side issue. Check NetBird API status or try again later.

### Handling API Errors

1. **Parse the status code** to determine the error category
2. **Read the response body** for specific error details
3. **Take appropriate action** based on the error type:
   - 400/422: Fix the request data
   - 401: Check authentication
   - 404: Verify resource exists
   - 500: Retry or report to NetBird

---

## Dependency Errors

### Group Deletion with Dependencies

When attempting to delete a group that is referenced by policies, you'll receive a dependency error.

**Error Format:**
```
cannot delete group: [count] policies depend on it: [policy_names]
```

**Example:**
```
cannot delete group: 3 policies depend on it: Admin-Access (sources, destinations), Team-Access (destinations), Exit-Node (sources)
```

### Handling Dependency Errors

**Option 1: Force Delete**

Use the `force=true` parameter to automatically remove the group from all policies:

```bash
delete_netbird_group --group_id "group-id" --force true
```

This will:
1. Find all policies referencing the group
2. Remove the group from sources, destinations, and authorized_groups
3. Remove rules that become invalid (no sources or destinations)
4. Delete policies that become empty (no rules)
5. Delete the group

**Option 2: Manual Cleanup**

1. Use `list_policies_by_group` to find dependent policies
2. Update each policy to remove the group reference
3. Delete the group

**Example:**
```bash
# Step 1: Find dependencies
list_policies_by_group --group_id "old-group-id"

# Step 2: Replace with another group (if needed)
replace_group_in_policies \
  --old_group_id "old-group-id" \
  --new_group_id "new-group-id"

# Step 3: Delete the group
delete_netbird_group --group_id "old-group-id"
```

---

## Retry Strategies

### When to Retry

Retry on transient errors:
- Network timeouts
- 500 Internal Server Error
- 503 Service Unavailable
- Connection errors

**Do NOT retry on:**
- 400 Bad Request (fix the request first)
- 401 Unauthorized (fix authentication first)
- 404 Not Found (resource doesn't exist)
- 422 Unprocessable Entity (fix business logic first)

### Exponential Backoff

For transient errors, use exponential backoff:

```go
func retryWithBackoff(operation func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = operation()
        if err == nil {
            return nil
        }
        
        // Check if error is retryable
        if !isRetryable(err) {
            return err
        }
        
        // Wait before retry: 1s, 2s, 4s, 8s, ...
        waitTime := time.Duration(1<<uint(i)) * time.Second
        time.Sleep(waitTime)
    }
    return fmt.Errorf("max retries exceeded: %w", err)
}

func isRetryable(err error) bool {
    // Check for network errors, 500, 503, etc.
    if strings.Contains(err.Error(), "status code: 500") {
        return true
    }
    if strings.Contains(err.Error(), "status code: 503") {
        return true
    }
    if strings.Contains(err.Error(), "connection refused") {
        return true
    }
    return false
}
```

### Recommended Retry Configuration

- **Max Retries:** 3
- **Initial Delay:** 1 second
- **Max Delay:** 8 seconds
- **Backoff Factor:** 2x

---

## Common Error Scenarios

### Scenario 1: Policy Creation Fails with Validation Error

**Error:**
```
validation error: rule 0 (SSH Access): port range start (22) must be <= end (21)
```

**Cause:** Port range is invalid (start > end)

**Solution:**
```json
// Fix the port range
"port_ranges": [{"start": 22, "end": 22}]  // or {"start": 21, "end": 22}
```

---

### Scenario 2: Cannot Delete Group

**Error:**
```
cannot delete group: 2 policies depend on it: Admin-Access, Team-Access
```

**Cause:** Group is referenced by policies

**Solution 1 (Force Delete):**
```bash
delete_netbird_group --group_id "group-id" --force true
```

**Solution 2 (Replace Group):**
```bash
# Replace with another group first
replace_group_in_policies \
  --old_group_id "old-group-id" \
  --new_group_id "new-group-id"

# Then delete
delete_netbird_group --group_id "old-group-id"
```

---

### Scenario 3: Nameserver Creation Fails

**Error:**
```
unexpected status code: 422, body: {"message":"nameserver group primary status is false and domains are empty, it should be primary or have at least one domain","code":422}
```

**Cause:** Nameserver must be either primary OR have domains

**Solution:**
```json
// Option 1: Make it primary
{
  "name": "DNS Servers",
  "nameservers": [...],
  "primary": true,
  "domains": []
}

// Option 2: Add domains
{
  "name": "DNS Servers",
  "nameservers": [...],
  "primary": false,
  "domains": ["example.com"]
}
```

---

### Scenario 4: Port Allocation Returns 404

**Error:**
```
unexpected status code: 404, body: 404 page not found
```

**Cause:** Port allocation feature not available on your NetBird instance

**Solution:** This feature may require:
- Specific NetBird version
- Enterprise/Cloud plan
- Feature flag enabled

Check NetBird documentation or contact support.

---

### Scenario 5: Unauthorized Error

**Error:**
```
unexpected status code: 401, body: {"message":"invalid token","code":401}
```

**Cause:** Invalid or missing API token

**Solution:**
1. Verify `NETBIRD_API_TOKEN` environment variable is set
2. Check token is valid in NetBird dashboard
3. Regenerate token if needed
4. Restart MCP server with new token

---

## Error Handling Best Practices

### 1. Always Validate Before API Calls

```go
// Validate rules before creating policy
if err := ValidatePolicyRules(rules); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// Format rules for API
formattedRules := make([]NetbirdPolicyRule, len(rules))
for i, rule := range rules {
    formattedRules[i] = FormatRuleForAPI(rule)
}

// Now make API call
```

### 2. Provide Context in Errors

```go
// Bad
return err

// Good
return fmt.Errorf("failed to create policy %s: %w", policyName, err)
```

### 3. Handle Partial Failures

When updating multiple resources, collect errors and continue:

```go
var errors []error
for _, policy := range policies {
    if err := updatePolicy(policy); err != nil {
        errors = append(errors, fmt.Errorf("policy %s: %w", policy.ID, err))
        continue  // Continue with other policies
    }
}

if len(errors) > 0 {
    return fmt.Errorf("partial failure: %v", errors)
}
```

### 4. Log Errors with Details

```go
log.Printf("API error: status=%d, body=%s, request=%v", 
    statusCode, responseBody, requestData)
```

### 5. Return Structured Error Information

```go
type OperationResult struct {
    Success        bool     `json:"success"`
    UpdatedIDs     []string `json:"updated_ids"`
    Errors         []string `json:"errors"`
    PartialFailure bool     `json:"partial_failure"`
}
```

---

## Getting Help

If you encounter errors not covered in this guide:

1. Check the [NetBird API Documentation](https://docs.netbird.io/api)
2. Review the [NetBird GitHub Issues](https://github.com/netbirdio/netbird/issues)
3. Open an issue in the [MCP-NetBird repository](https://github.com/aantti/mcp-netbird/issues)
4. Include:
   - Full error message
   - Request parameters
   - NetBird version
   - MCP-NetBird version

---

## Summary

- **Validation errors** occur before API calls - fix input data
- **API errors** come from NetBird - check status code and response body
- **Dependency errors** occur when deleting referenced resources - use force delete or manual cleanup
- **Retry** only on transient errors (500, 503, network issues)
- **Always validate** input before making API calls
- **Provide context** in error messages for easier debugging
