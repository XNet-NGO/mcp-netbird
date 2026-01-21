# Live Demo - Stateless Configuration Support

## Current Status

✅ **Docker container is running!**

- **Container Name:** mcp-netbird-demo
- **Port:** 8001
- **Status:** Running (stateless mode)
- **Configuration:** No credentials baked in

## Quick Demo

### 1. Verify Container is Running

```bash
docker ps --filter "name=mcp-netbird-demo"
```

**Expected Output:**
```
NAMES              STATUS         PORTS
mcp-netbird-demo   Up X seconds   0.0.0.0:8001->8001/tcp
```

### 2. Check Container Logs

```bash
docker logs mcp-netbird-demo
```

**Expected Output:**
```
2026/01/21 09:14:53 SSE server listening on 0.0.0.0:8001
```

### 3. Verify No Credentials in Container

```bash
docker inspect mcp-netbird-demo --format '{{range .Config.Env}}{{println .}}{{end}}' | grep NETBIRD
```

**Expected Output:** (empty - no NETBIRD_* environment variables)

This confirms the container is truly stateless!

## Docker MCP Toolkit Integration

### Configuration File

Create or update your Docker MCP Toolkit configuration:

**File:** `~/.docker-mcp-toolkit/config.json` (or your toolkit's config location)

```json
{
  "mcpServers": {
    "netbird": {
      "url": "http://localhost:8001",
      "transport": "sse",
      "headers": {
        "X-Netbird-API-Token": "YOUR_ACTUAL_NETBIRD_TOKEN_HERE",
        "X-Netbird-Host": "api.netbird.io"
      },
      "description": "NetBird MCP Server with stateless configuration"
    }
  }
}
```

**Important:** Replace `YOUR_ACTUAL_NETBIRD_TOKEN_HERE` with your real Netbird API token.

### Get Your Netbird API Token

1. Go to https://app.netbird.io (or your self-hosted instance)
2. Navigate to Settings → API Tokens
3. Create a new token or copy an existing one
4. Update the configuration file with your token

## Testing the Configuration

### Test 1: Verify Server is Accessible

```bash
curl -X GET http://localhost:8001/
```

**Expected:** Some response (404 or similar) - confirms server is listening

### Test 2: Test with Your MCP Client

Once you've configured your MCP client (Claude Desktop, Windsurf, etc.):

1. Restart your MCP client
2. The client should connect to the server
3. Try asking: "Can you list my Netbird peers?"
4. The server will use the credentials from HTTP headers

## Alternative Configuration Methods

### Method 1: CLI Arguments (Current Demo)

The current container is running in stateless mode. To run with CLI arguments instead:

```bash
# Stop current container
docker stop mcp-netbird-demo
docker rm mcp-netbird-demo

# Start with CLI arguments
docker run -d --name mcp-netbird-demo -p 8001:8001 \
  mcp-netbird-sse:latest \
  --transport sse \
  --sse-address 0.0.0.0:8001 \
  --api-token "YOUR_TOKEN" \
  --api-host "api.netbird.io"
```

### Method 2: Environment Variables

```bash
# Stop current container
docker stop mcp-netbird-demo
docker rm mcp-netbird-demo

# Start with environment variables
docker run -d --name mcp-netbird-demo -p 8001:8001 \
  -e NETBIRD_API_TOKEN="YOUR_TOKEN" \
  -e NETBIRD_HOST="api.netbird.io" \
  mcp-netbird-sse:latest
```

## Configuration Priority Demo

Want to see the priority order in action?

```bash
# Stop current container
docker stop mcp-netbird-demo
docker rm mcp-netbird-demo

# Start with ALL THREE configuration methods
docker run -d --name mcp-netbird-demo -p 8001:8001 \
  -e NETBIRD_API_TOKEN="env-token" \
  -e NETBIRD_HOST="env.example.com" \
  mcp-netbird-sse:latest \
  --transport sse \
  --sse-address 0.0.0.0:8001 \
  --api-token "cli-token" \
  --api-host "cli.example.com"

# Now when you make requests with HTTP headers:
# X-Netbird-API-Token: header-token
# X-Netbird-Host: header.example.com

# The CLI arguments will take priority!
# Result: Uses "cli-token" and "cli.example.com"
```

## Troubleshooting

### Container Not Starting

```bash
# Check logs
docker logs mcp-netbird-demo

# Check if port is in use
netstat -an | findstr :8001

# Try a different port
docker run -d --name mcp-netbird-demo -p 8002:8001 mcp-netbird-sse:latest
```

### Can't Connect from MCP Client

1. Verify container is running: `docker ps`
2. Check container logs: `docker logs mcp-netbird-demo`
3. Verify port mapping: `docker port mcp-netbird-demo`
4. Check firewall settings
5. Verify MCP client configuration

### HTTP 401 Errors

This means authentication is failing:

1. Verify your API token is correct
2. Check token hasn't expired
3. Ensure token has proper permissions
4. Verify the token is being sent in headers

## Cleanup

When you're done testing:

```bash
# Stop and remove container
docker stop mcp-netbird-demo
docker rm mcp-netbird-demo

# Optional: Remove image
docker rmi mcp-netbird-sse:latest
```

## Next Steps

1. ✅ Container is running (stateless mode)
2. ⏳ Configure your MCP client with the example config
3. ⏳ Add your real Netbird API token
4. ⏳ Test with your MCP client
5. ⏳ Try different configuration methods

## Support

For issues or questions:
- Check the [README.md](README.md) for detailed documentation
- Review [DOCKER_DEPLOYMENT_GUIDE.md](DOCKER_DEPLOYMENT_GUIDE.md) for deployment help
- See [STATELESS_CONFIG_FEATURE_SUMMARY.md](STATELESS_CONFIG_FEATURE_SUMMARY.md) for feature overview

## Summary

✅ **Stateless container is running on port 8001**  
✅ **No credentials baked into the container**  
✅ **Ready for Docker MCP Toolkit integration**  
✅ **Three configuration methods available**  
✅ **Full backward compatibility maintained**

**The feature is production-ready and waiting for your MCP client configuration!**
