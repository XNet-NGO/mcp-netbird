# Netbird MCP Debug Report

**Date:** 2026-01-19  
**Status:** Fixed - Requires Q CLI Restart

## Summary

Debugged the Netbird MCP server and identified two issues:

### ✅ Working Tools (6/7)
- `list_netbird_peers` - Returns 12 peers successfully
- `list_netbird_networks` - Returns 1 network (Core-Net)
- `list_netbird_policies` - Returns 8 policies
- `list_netbird_posture_checks` - Returns empty array (no posture checks configured)
- `list_netbird_nameservers` - Returns 1 nameserver group (XNet-DNS)

### ❌ Broken Tools (2/7)

#### 1. `list_netbird_groups` - JSON Decoding Error (FIXED)

**Error:**
```
json: cannot unmarshal object into Go struct field NetbirdGroup.resources of type string
```

**Root Cause:**  
The `Resources` field in the `NetbirdGroup` struct was defined as `[]string`, but the Netbird API returns an array of objects with `id` and `name` fields (same structure as `Peers`).

**Fix Applied:**  
File: `tools/groups.go` (line 18)

```go
// Before:
Resources      []string             `json:"resources"`

// After:
Resources      []NetbirdGroupMember `json:"resources"`
```

**Status:** Fixed and compiled. Binary updated at `C:\Users\xnet-admin\go\bin\mcp-netbird.exe`

#### 2. `list_netbird_port_allocations` - 404 Not Found

**Error:**
```
unexpected status code: 404, body: 404 page not found
```

**Root Cause:**  
The API endpoint `/peers/{peer_id}/ingress/ports` does not exist or is not available in the current Netbird API version.

**Tested With:** Peer ID `d5n0fu3ngf8s73bjmja0`

**Status:** API endpoint issue - requires Netbird API investigation or documentation review.

## Testing Results

### Successful API Calls

**Peers (12 total):**
- dm2qsqw (Android, connected)
- IndigoStation (Windows 11, disconnected)
- ip-172-31-11-4 (Ubuntu, connected)
- ip-172-31-24-183 (Debian, connected)
- ip-172-31-24-53 (Debian, disconnected)
- ip-172-31-44-163 (Debian, connected)
- kansas_g_sys (Android, login expired)
- OpenWrt (disconnected)
- serv-1 (Debian, disconnected)
- serv-2 (Debian, disconnected)
- xnet-book (Windows 11, disconnected)
- xnet-exit-node (Alpine, disconnected)

**Networks (1 total):**
- Core-Net (4 routing peers, 1 resource, 5 policies)

**Policies (8 total):**
- Default-Access
- Admin-Access
- SSH-Access
- Team-Access
- Exit-Node
- Openwrt
- Global-Exit
- Client-Access

**Nameservers (1 group):**
- XNet-DNS (2 nameservers: 100.105.86.234, 100.22.6.88)

## Next Steps

1. **Restart Q CLI** to reload the fixed MCP server
2. **Test `list_netbird_groups`** to verify the fix
3. **Investigate port allocations API** - check Netbird API documentation or version requirements

## Build Commands

```bash
# Build the MCP server
cd C:\Users\xnet-admin\Repos\mcp-netbird
go build -o mcp-netbird.exe ./cmd/mcp-netbird

# Copy to go bin
copy /Y mcp-netbird.exe C:\Users\xnet-admin\go\bin\mcp-netbird.exe

# Kill old process (if needed)
taskkill /F /IM mcp-netbird.exe
```

## Files Modified

- `tools/groups.go` - Fixed `Resources` field type from `[]string` to `[]NetbirdGroupMember`

## Environment

- **OS:** Windows 11
- **Go Version:** (detected from build)
- **MCP Server:** mcp-netbird
- **Netbird API:** api.netbird.io
- **Q CLI:** Amazon Q Developer CLI
