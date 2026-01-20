# MCP-NetBird Server Improvements Needed

## Issues Discovered During Live Testing

### 1. Policy Update/Create with Rules - API Format Mismatch ⚠️

**Issue**: 
- `create_netbird_policy` and `update_netbird_policy` fail when passing `rules` parameter
- Error: "couldn't parse JSON request" (400 Bad Request)

**Root Cause**:
The NetBird API likely expects a different JSON structure for rules than what we're sending. The current implementation sends the rules as a nested array of objects, but the API might expect:
1. A different field name
2. A different nesting structure
3. Rules to be sent separately via a different endpoint

**Testing Performed**:
```
Attempted to update policy d54fd4rngf8s73892otg with rules parameter
Result: 400 Bad Request - "couldn't parse JSON request"
```

**Recommendation**:
1. Investigate the actual NetBird API documentation for policy creation/update
2. Test with curl/Postman to determine the correct JSON format
3. Consider if rules need to be added via a separate endpoint after policy creation
4. Add integration tests for policy creation with rules

**Workaround**:
Currently, policies must be created through the NetBird Dashboard UI.

---

### 2. Group Deletion - No Cascade Option ⚠️

**Issue**:
Cannot delete a group that's referenced in policies, even though we want to consolidate groups.

**Error Message**:
```
400 Bad Request: "group has been linked to policy: Default-Access"
```

**Root Cause**:
The API protects against orphaning policy references but doesn't provide:
1. A force delete option
2. A cascade delete option
3. A way to get all policies referencing a group

**Testing Performed**:
```
Attempted to delete group d53i0drngf8s73892ocg (Core-Net)
Result: 400 Bad Request - group linked to policies
```

**Recommendation**:
1. Add a `list_policies_by_group` helper function
2. Add a `force` parameter to `delete_netbird_group` that:
   - Lists all policies referencing the group
   - Updates those policies to remove the group reference
   - Then deletes the group
3. Add proper error handling with helpful messages

**Workaround**:
Manually update all policies through the UI before deleting the group.

---

### 3. Missing Helper Functions 📝

**Needed Functions**:

1. **`list_policies_by_group(group_id)`**
   - Returns all policies that reference a specific group
   - Useful for understanding dependencies before deletion

2. **`replace_group_in_policies(old_group_id, new_group_id)`**
   - Bulk update all policies to replace one group with another
   - Useful for group consolidation

3. **`validate_policy_rules(rules)`**
   - Validate rule structure before sending to API
   - Provide helpful error messages about what's wrong

4. **`get_policy_template()`**
   - Return an example policy structure
   - Help users understand the correct format

---

## Successful Operations ✅

The following operations worked perfectly during testing:

1. ✅ **Network Operations**
   - `list_netbird_networks`
   - `create_netbird_network`
   - `get_netbird_network`

2. ✅ **Network Resource Operations**
   - `create_netbird_network_resource`
   - `list_netbird_network_resources`
   - All CRUD operations worked flawlessly

3. ✅ **Network Router Operations**
   - `create_netbird_network_router`
   - `list_netbird_network_routers`
   - `get_netbird_network_router`

4. ✅ **Route Operations**
   - `list_netbird_routes`
   - `delete_netbird_route`
   - Successfully migrated from deprecated routes to networks

5. ✅ **Group Operations** (except delete with dependencies)
   - `list_netbird_groups`
   - `get_netbird_group`
   - `update_netbird_group`
   - `create_netbird_group`

6. ✅ **Peer Operations**
   - `list_netbird_peers`
   - `get_netbird_peer`
   - `update_netbird_peer`
   - `delete_netbird_peer`

7. ✅ **Policy Operations** (except create/update with rules)
   - `list_netbird_policies`
   - `get_netbird_policy`
   - `delete_netbird_policy`

---

## Real-World Use Cases Tested ✅

### Migration from Network Routes to Networks
**Status**: ✅ Successful

Successfully migrated 3 network routes to the new Networks feature:
1. Created Core-Net network with 3 resources
2. Created OpenWrt LAN network with 1 resource
3. Configured routing peers for each network
4. Deleted old deprecated routes

**Code Quality**: Excellent - all operations worked as expected

---

### Network Configuration Management
**Status**: ✅ Successful

Successfully managed complex network configuration:
1. Added new exit node (ip-172-31-11-4)
2. Removed old exit node (xnet-exit-node)
3. Updated group memberships
4. Configured SSH access
5. Created network routers with proper metrics

**Code Quality**: Excellent - intuitive API design

---

### Group and Policy Review
**Status**: ⚠️ Partially Successful

Successfully reviewed configuration but hit limitations:
1. ✅ Listed all groups and policies
2. ✅ Identified configuration issues
3. ⚠️ Could not programmatically fix policy issues
4. ⚠️ Could not consolidate duplicate groups

**Code Quality**: Good for read operations, needs improvement for complex updates

---

## Priority Recommendations

### High Priority 🔴
1. **Fix Policy Rules API** - Critical for automation
   - Research correct API format
   - Add integration tests
   - Document the correct structure

2. **Add Group Deletion with Cascade** - Important for maintenance
   - Implement force delete option
   - Add policy dependency checking
   - Provide clear error messages

### Medium Priority 🟡
3. **Add Helper Functions** - Improves usability
   - `list_policies_by_group`
   - `replace_group_in_policies`
   - `validate_policy_rules`

4. **Improve Error Messages** - Better developer experience
   - Parse API errors and provide context
   - Suggest fixes for common issues
   - Add examples to error messages

### Low Priority 🟢
5. **Add Policy Templates** - Nice to have
   - Common policy patterns
   - Example configurations
   - Best practices documentation

---

## Testing Recommendations

### Integration Tests Needed
1. Policy creation with rules
2. Policy update with rules
3. Group deletion with policy dependencies
4. Network migration scenarios
5. Complex group membership updates

### Test Data
Use the current live configuration as test fixtures:
- 11 peers across 10 groups
- 2 networks with multiple resources
- 8 policies with various rules
- Real-world complexity

---

## Documentation Improvements

### Add Examples For:
1. Creating policies with complex rules
2. Migrating from network routes to networks
3. Managing group dependencies
4. Bulk operations (updating multiple policies)
5. Error handling and recovery

### API Reference
1. Document all parameter formats
2. Show request/response examples
3. List common error codes
4. Provide troubleshooting guide

---

## Conclusion

The mcp-netbird server is **production-ready for most operations** but needs improvements for:
1. Policy management with rules
2. Group consolidation workflows
3. Complex bulk operations

The core functionality (networks, resources, routers, peers, groups) works excellently and successfully handled a real-world migration from deprecated features to new ones.

**Overall Assessment**: 8/10 - Excellent foundation, needs polish for advanced use cases.
