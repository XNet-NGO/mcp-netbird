# MCP-NetBird Helper Functions Guide

This guide explains the helper functions available in MCP-NetBird for common workflows, particularly group consolidation and policy management.

## Table of Contents

1. [Overview](#overview)
2. [Policy Template Helper](#policy-template-helper)
3. [Group Dependency Discovery](#group-dependency-discovery)
4. [Group Replacement](#group-replacement)
5. [Force Delete Group](#force-delete-group)
6. [Common Workflows](#common-workflows)

---

## Overview

MCP-NetBird provides several helper functions to simplify common administrative tasks:

| Function | Purpose | Use Case |
|----------|---------|----------|
| `get_policy_template` | Get example policy structures | Learning policy format, creating new policies |
| `list_policies_by_group` | Find policies using a group | Before deleting/modifying groups |
| `replace_group_in_policies` | Replace group across all policies | Group consolidation, migration |
| `delete_netbird_group` (force) | Delete group and clean dependencies | Removing unused groups |

---

## Policy Template Helper

### Function: `get_policy_template`

Returns example policy structures with documentation to help you create valid policies.

### Usage

```bash
get_policy_template
```

### Returns

```json
{
  "simple_example": {
    "name": "Simple Access Policy",
    "description": "Basic policy with one rule",
    "enabled": true,
    "rules": [
      {
        "name": "Allow Access",
        "description": "Allow group A to access group B",
        "enabled": true,
        "action": "accept",
        "bidirectional": true,
        "protocol": "all",
        "sources": ["source-group-id"],
        "destinations": ["destination-group-id"]
      }
    ]
  },
  "complex_example": {
    "name": "Complex Access Policy",
    "description": "Policy with multiple rules and port restrictions",
    "enabled": true,
    "rules": [
      {
        "name": "HTTP/HTTPS Access",
        "description": "Allow web traffic",
        "enabled": true,
        "action": "accept",
        "bidirectional": true,
        "protocol": "tcp",
        "sources": ["user-group-id"],
        "destinations": ["web-server-group-id"],
        "port_ranges": [
          {"start": 80, "end": 80},
          {"start": 443, "end": 443}
        ]
      },
      {
        "name": "Database Access",
        "description": "Restricted database access",
        "enabled": true,
        "action": "accept",
        "bidirectional": false,
        "protocol": "tcp",
        "sources": ["app-server-group-id"],
        "destinations": ["database-group-id"],
        "port_ranges": [{"start": 5432, "end": 5432}],
        "authorized_groups": {
          "admin-group-id": ["dba-group-id"]
        }
      }
    ]
  },
  "field_descriptions": {
    "name": "Policy name (required)",
    "description": "Policy description (optional)",
    "enabled": "Whether policy is active (required, boolean)",
    "rules": "Array of policy rules (required)",
    "rule.name": "Rule name (required)",
    "rule.enabled": "Whether rule is active (required, boolean)",
    "rule.action": "Either 'accept' or 'drop' (required)",
    "rule.bidirectional": "Allow traffic in both directions (required, boolean)",
    "rule.protocol": "One of: 'tcp', 'udp', 'icmp', 'all' (required)",
    "rule.sources": "Array of source group IDs (required if no sourceResource)",
    "rule.destinations": "Array of destination group IDs (required if no destinationResource)",
    "rule.port_ranges": "Array of port ranges for TCP/UDP (optional)",
    "rule.authorized_groups": "Map of group IDs to authorized group arrays (optional)"
  }
}
```

### Example: Creating a Policy from Template

```bash
# Get template
template=$(get_policy_template)

# Modify for your needs
policy='{
  "name": "My Custom Policy",
  "description": "Based on simple template",
  "enabled": true,
  "rules": [
    {
      "name": "Custom Rule",
      "enabled": true,
      "action": "accept",
      "bidirectional": true,
      "protocol": "tcp",
      "sources": ["my-source-group-id"],
      "destinations": ["my-dest-group-id"],
      "port_ranges": [{"start": 8080, "end": 8080}]
    }
  ]
}'

# Create policy
create_netbird_policy --policy "$policy"
```

---

## Group Dependency Discovery

### Function: `list_policies_by_group`

Finds all policies that reference a specific group in sources, destinations, or authorized_groups.

### Usage

```bash
list_policies_by_group --group_id "d535b93ngf8s73892nng"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `group_id` | string | Yes | The ID of the group to search for |

### Returns

```json
{
  "group_id": "d535b93ngf8s73892nng",
  "group_name": "Admins",
  "references": [
    {
      "policy_id": "d54anfrngf8s73892op0",
      "policy_name": "Admin-Access",
      "rule_id": "d54anfrngf8s73892op0",
      "rule_name": "Admin Rule",
      "location": "sources"
    },
    {
      "policy_id": "d54gamjngf8s73892p2g",
      "policy_name": "Team-Access",
      "rule_id": "d54gamjngf8s73892p2g",
      "rule_name": "Team Rule",
      "location": "destinations"
    },
    {
      "policy_id": "d54gkc3ngf8s73892p4g",
      "policy_name": "Exit-Node",
      "rule_id": "d54gkc3ngf8s73892p4g",
      "rule_name": "Exit Rule",
      "location": "authorized_groups"
    }
  ],
  "total_references": 3
}
```

### Location Values

- `sources`: Group appears in rule's sources array
- `destinations`: Group appears in rule's destinations array
- `authorized_groups`: Group appears as a key in authorized_groups map

### Use Cases

1. **Before Deleting a Group**
   ```bash
   # Check if group is used
   list_policies_by_group --group_id "old-group-id"
   
   # If references found, either:
   # - Use force delete
   # - Replace with another group
   # - Manually update policies
   ```

2. **Auditing Group Usage**
   ```bash
   # Find all policies using a specific group
   list_policies_by_group --group_id "admin-group-id"
   ```

3. **Impact Analysis**
   ```bash
   # Before modifying a group, see what policies will be affected
   list_policies_by_group --group_id "infrastructure-group-id"
   ```

---

## Group Replacement

### Function: `replace_group_in_policies`

Replaces one group with another across all policies. Updates sources, destinations, and authorized_groups.

### Usage

```bash
replace_group_in_policies \
  --old_group_id "d535b93ngf8s73892nn0" \
  --new_group_id "d53i0drngf8s73892ocg"
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `old_group_id` | string | Yes | The group ID to replace |
| `new_group_id` | string | Yes | The group ID to replace with |

### Returns

```json
{
  "old_group_id": "d535b93ngf8s73892nn0",
  "new_group_id": "d53i0drngf8s73892ocg",
  "updated_policies": [
    "d54anfrngf8s73892op0",
    "d54gamjngf8s73892p2g",
    "d54gkc3ngf8s73892p4g"
  ],
  "errors": [],
  "summary": {
    "total_policies_checked": 8,
    "policies_updated": 3,
    "policies_failed": 0
  }
}
```

### Behavior

1. **Finds all references** to old group using `list_policies_by_group`
2. **Fetches each policy** that references the old group
3. **Replaces group ID** in:
   - `sources` arrays
   - `destinations` arrays
   - `authorized_groups` map keys
4. **Updates each policy** via PUT request
5. **Continues on errors** - collects errors but processes all policies
6. **Returns summary** of updated policies and any errors

### Use Cases

#### 1. Group Consolidation

Merge two groups by replacing one with the other:

```bash
# Step 1: Find what will be affected
list_policies_by_group --group_id "old-infrastructure-group"

# Step 2: Replace in all policies
replace_group_in_policies \
  --old_group_id "old-infrastructure-group" \
  --new_group_id "new-infrastructure-group"

# Step 3: Delete old group
delete_netbird_group --group_id "old-infrastructure-group"
```

#### 2. Group Migration

Move from old naming scheme to new:

```bash
# Replace "Servers" with "Infrastructure"
replace_group_in_policies \
  --old_group_id "servers-group-id" \
  --new_group_id "infrastructure-group-id"
```

#### 3. Fixing Misconfigurations

Correct policies that reference the wrong group:

```bash
# Replace incorrect group with correct one
replace_group_in_policies \
  --old_group_id "wrong-group-id" \
  --new_group_id "correct-group-id"
```

### Error Handling

If some policies fail to update:

```json
{
  "old_group_id": "old-id",
  "new_group_id": "new-id",
  "updated_policies": ["policy-1", "policy-2"],
  "errors": [
    "policy-3: unexpected status code: 422",
    "policy-4: validation error: rule 0: missing destination"
  ],
  "summary": {
    "total_policies_checked": 5,
    "policies_updated": 2,
    "policies_failed": 2
  }
}
```

**Action:** Review errors and manually fix failed policies.

---

## Force Delete Group

### Function: `delete_netbird_group` with `force=true`

Deletes a group and automatically removes it from all dependent policies.

### Usage

```bash
# Normal delete (fails if dependencies exist)
delete_netbird_group --group_id "group-id"

# Force delete (removes from all policies first)
delete_netbird_group --group_id "group-id" --force true
```

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `group_id` | string | Yes | The group ID to delete |
| `force` | boolean | No | If true, removes from all policies first (default: false) |

### Returns (force=false, with dependencies)

```json
{
  "error": "cannot delete group: 3 policies depend on it",
  "group_id": "d535b93ngf8s73892nn0",
  "dependent_policies": [
    {
      "policy_id": "policy-1",
      "policy_name": "Admin-Access",
      "locations": ["sources", "destinations"]
    },
    {
      "policy_id": "policy-2",
      "policy_name": "Team-Access",
      "locations": ["destinations"]
    }
  ]
}
```

### Returns (force=true)

```json
{
  "status": "deleted",
  "group_id": "d535b93ngf8s73892nn0",
  "cleanup_summary": {
    "policies_modified": 3,
    "rules_removed": 2,
    "policies_deleted": 1
  },
  "modified_policies": [
    {
      "policy_id": "policy-1",
      "policy_name": "Admin-Access",
      "action": "updated",
      "changes": "removed from sources and destinations"
    },
    {
      "policy_id": "policy-2",
      "policy_name": "Team-Access",
      "action": "updated",
      "changes": "removed from destinations"
    },
    {
      "policy_id": "policy-3",
      "policy_name": "Old-Policy",
      "action": "deleted",
      "changes": "policy became empty after removing group"
    }
  ]
}
```

### Behavior (force=true)

1. **Finds dependencies** using `list_policies_by_group`
2. **For each dependent policy:**
   - Fetches current configuration
   - Removes group from sources arrays
   - Removes group from destinations arrays
   - Removes group from authorized_groups keys
3. **Validates rules:**
   - If rule has no sources or destinations, removes the rule
4. **Validates policies:**
   - If policy has no rules, deletes the policy
   - Otherwise, updates the policy
5. **Deletes the group** after all dependencies resolved
6. **Returns summary** of all operations

### Use Cases

#### 1. Clean Removal of Unused Groups

```bash
# Remove group and all references
delete_netbird_group --group_id "deprecated-group" --force true
```

#### 2. Safe Deletion Check

```bash
# First check what would be affected (force=false)
delete_netbird_group --group_id "test-group"

# Review dependencies, then force delete if appropriate
delete_netbird_group --group_id "test-group" --force true
```

#### 3. Cleanup After Migration

```bash
# After replacing group in policies, delete old group
replace_group_in_policies \
  --old_group_id "old-group" \
  --new_group_id "new-group"

delete_netbird_group --group_id "old-group" --force true
```

### Safety Considerations

⚠️ **Force delete is destructive:**
- Removes group from all policies
- May delete entire policies if they become empty
- Cannot be undone

✅ **Best practices:**
1. Always run without `force` first to see dependencies
2. Review the list of dependent policies
3. Consider using `replace_group_in_policies` instead if you want to preserve policies
4. Backup your configuration before force deleting

---

## Common Workflows

### Workflow 1: Group Consolidation

**Scenario:** You have duplicate groups (Infrastructure and Core-Net) and want to consolidate them.

```bash
# Step 1: Check what policies use each group
list_policies_by_group --group_id "infrastructure-group-id"
list_policies_by_group --group_id "core-net-group-id"

# Step 2: Decide which group to keep (e.g., Infrastructure)
# Replace Core-Net with Infrastructure in all policies
replace_group_in_policies \
  --old_group_id "core-net-group-id" \
  --new_group_id "infrastructure-group-id"

# Step 3: Verify replacement was successful
list_policies_by_group --group_id "core-net-group-id"
# Should return no references

# Step 4: Delete the old group
delete_netbird_group --group_id "core-net-group-id"
```

### Workflow 2: Safe Group Deletion

**Scenario:** You want to delete a group but aren't sure what depends on it.

```bash
# Step 1: Check dependencies
list_policies_by_group --group_id "old-group-id"

# Step 2: Review the output
# If no dependencies, safe to delete:
delete_netbird_group --group_id "old-group-id"

# If dependencies exist, decide:
# Option A: Force delete (removes from all policies)
delete_netbird_group --group_id "old-group-id" --force true

# Option B: Replace with another group first
replace_group_in_policies \
  --old_group_id "old-group-id" \
  --new_group_id "replacement-group-id"
delete_netbird_group --group_id "old-group-id"
```

### Workflow 3: Policy Creation from Template

**Scenario:** You need to create a new policy and want to use the correct format.

```bash
# Step 1: Get template
get_policy_template

# Step 2: Copy the appropriate example (simple or complex)
# Step 3: Modify with your group IDs and requirements
# Step 4: Create the policy

create_netbird_policy --policy '{
  "name": "My New Policy",
  "description": "Based on template",
  "enabled": true,
  "rules": [
    {
      "name": "Access Rule",
      "enabled": true,
      "action": "accept",
      "bidirectional": true,
      "protocol": "all",
      "sources": ["source-group-id"],
      "destinations": ["dest-group-id"]
    }
  ]
}'
```

### Workflow 4: Group Renaming (via Replacement)

**Scenario:** You want to "rename" a group by creating a new one and migrating policies.

```bash
# Step 1: Create new group with desired name
create_netbird_group \
  --name "New Group Name" \
  --peers '["peer-1", "peer-2"]'

# Step 2: Get new group ID from response
new_group_id="<new-group-id>"

# Step 3: Replace old group with new in all policies
replace_group_in_policies \
  --old_group_id "old-group-id" \
  --new_group_id "$new_group_id"

# Step 4: Delete old group
delete_netbird_group --group_id "old-group-id"
```

### Workflow 5: Audit Group Usage

**Scenario:** You want to understand how a group is being used across your policies.

```bash
# Get comprehensive view of group usage
list_policies_by_group --group_id "admin-group-id"

# Output shows:
# - Which policies reference the group
# - Where in each policy (sources, destinations, authorized_groups)
# - Rule names for context

# Use this information to:
# - Understand access patterns
# - Plan changes
# - Document configuration
```

---

## Best Practices

### 1. Always Check Dependencies First

Before modifying or deleting groups:
```bash
list_policies_by_group --group_id "group-id"
```

### 2. Use Templates for Consistency

Start with templates when creating policies:
```bash
get_policy_template
```

### 3. Test Changes in Non-Production First

If possible, test group replacements and deletions in a test environment.

### 4. Document Your Changes

Keep track of group consolidations and policy changes:
```bash
# Before
list_policies_by_group --group_id "old-group" > before.json

# Make changes
replace_group_in_policies --old_group_id "old-group" --new_group_id "new-group"

# After
list_policies_by_group --group_id "new-group" > after.json
```

### 5. Handle Errors Gracefully

Check return values and handle partial failures:
```bash
result=$(replace_group_in_policies --old_group_id "old" --new_group_id "new")
if echo "$result" | grep -q "errors"; then
    echo "Some policies failed to update"
    echo "$result" | jq '.errors'
fi
```

---

## Troubleshooting

### Issue: Group replacement doesn't update all policies

**Cause:** Some policies may have failed to update due to validation errors.

**Solution:**
1. Check the `errors` array in the response
2. Manually fix the failed policies
3. Re-run the replacement for those specific policies

### Issue: Force delete removes more than expected

**Cause:** Group was referenced in many places, including as the only source/destination in some rules.

**Solution:**
1. Always run without `force` first to see what will be affected
2. Consider using `replace_group_in_policies` instead to preserve policies
3. Review the cleanup summary after force delete

### Issue: Cannot find group ID

**Cause:** Need to get group ID from group name.

**Solution:**
```bash
# List all groups and find the ID
list_netbird_groups | jq '.[] | select(.name=="Group Name") | .id'
```

---

## Summary

Helper functions simplify common administrative tasks:

- **`get_policy_template`**: Learn policy format and create valid policies
- **`list_policies_by_group`**: Discover dependencies before making changes
- **`replace_group_in_policies`**: Consolidate or migrate groups across policies
- **`delete_netbird_group --force`**: Clean removal of groups with dependencies

Use these functions together for safe, efficient group and policy management.
