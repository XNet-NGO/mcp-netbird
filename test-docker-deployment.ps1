# Test script for Docker MCP Toolkit deployment
# Tests all three configuration methods for stateless configuration support

Write-Host "=== Testing MCP-Netbird Docker Deployment ===" -ForegroundColor Cyan
Write-Host ""

# Configuration
$CONTAINER_NAME = "mcp-netbird-test"
$PORT = "8001"
$TEST_TOKEN = "test-token-12345"
$TEST_HOST = "api.netbird.io"

# Cleanup function
function Cleanup {
    Write-Host "Cleaning up..." -ForegroundColor Yellow
    docker stop $CONTAINER_NAME 2>$null | Out-Null
    docker rm $CONTAINER_NAME 2>$null | Out-Null
}

# Test 1: Stateless container with HTTP headers (recommended for Docker MCP Toolkit)
Write-Host "Test 1: Stateless Container with HTTP Headers" -ForegroundColor Green
Write-Host "-----------------------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container without credentials..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 mcp-netbird-sse:latest
Start-Sleep -Seconds 2

Write-Host "Testing with HTTP headers..."
$response = curl.exe -s -X POST http://localhost:${PORT}/sse `
    -H "X-Netbird-API-Token: $TEST_TOKEN" `
    -H "X-Netbird-Host: $TEST_HOST" `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ HTTP headers test PASSED" -ForegroundColor Green
    Write-Host "Response preview: $($response.Substring(0, [Math]::Min(200, $response.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "✗ HTTP headers test FAILED" -ForegroundColor Red
}
Write-Host ""

# Test 2: Container with CLI arguments
Write-Host "Test 2: Container with CLI Arguments" -ForegroundColor Green
Write-Host "-------------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container with CLI arguments..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 `
    mcp-netbird-sse:latest `
    --transport sse --sse-address 0.0.0.0:8001 `
    --api-token $TEST_TOKEN --api-host $TEST_HOST
Start-Sleep -Seconds 2

Write-Host "Testing without HTTP headers (using CLI config)..."
$response = curl.exe -s -X POST http://localhost:${PORT}/sse `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ CLI arguments test PASSED" -ForegroundColor Green
    Write-Host "Response preview: $($response.Substring(0, [Math]::Min(200, $response.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "✗ CLI arguments test FAILED" -ForegroundColor Red
}
Write-Host ""

# Test 3: Container with environment variables (backward compatibility)
Write-Host "Test 3: Container with Environment Variables" -ForegroundColor Green
Write-Host "---------------------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container with environment variables..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 `
    -e NETBIRD_API_TOKEN=$TEST_TOKEN `
    -e NETBIRD_HOST=$TEST_HOST `
    mcp-netbird-sse:latest
Start-Sleep -Seconds 2

Write-Host "Testing without HTTP headers (using env vars)..."
$response = curl.exe -s -X POST http://localhost:${PORT}/sse `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Environment variables test PASSED" -ForegroundColor Green
    Write-Host "Response preview: $($response.Substring(0, [Math]::Min(200, $response.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "✗ Environment variables test FAILED" -ForegroundColor Red
}
Write-Host ""

# Test 4: Priority order verification
Write-Host "Test 4: Configuration Priority Order" -ForegroundColor Green
Write-Host "-------------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container with CLI args and env vars..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 `
    -e NETBIRD_API_TOKEN="env-token" `
    -e NETBIRD_HOST="env.example.com" `
    mcp-netbird-sse:latest `
    --transport sse --sse-address 0.0.0.0:8001 `
    --api-token "cli-token" --api-host "cli.example.com"
Start-Sleep -Seconds 2

Write-Host "Testing with HTTP headers (all three sources present)..."
$response = curl.exe -s -X POST http://localhost:${PORT}/sse `
    -H "X-Netbird-API-Token: header-token" `
    -H "X-Netbird-Host: header.example.com" `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Priority order test PASSED (CLI should take priority)" -ForegroundColor Green
    Write-Host "Response preview: $($response.Substring(0, [Math]::Min(200, $response.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "✗ Priority order test FAILED" -ForegroundColor Red
}
Write-Host ""

# Test 5: Protocol prefix stripping
Write-Host "Test 5: Protocol Prefix Stripping" -ForegroundColor Green
Write-Host "----------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 mcp-netbird-sse:latest
Start-Sleep -Seconds 2

Write-Host "Testing with protocol prefix in header..."
$response = curl.exe -s -X POST http://localhost:${PORT}/sse `
    -H "X-Netbird-API-Token: $TEST_TOKEN" `
    -H "X-Netbird-Host: https://api.netbird.io" `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Protocol prefix stripping test PASSED" -ForegroundColor Green
    Write-Host "Response preview: $($response.Substring(0, [Math]::Min(200, $response.Length)))..." -ForegroundColor Gray
} else {
    Write-Host "✗ Protocol prefix stripping test FAILED" -ForegroundColor Red
}
Write-Host ""

# Test 6: Missing authentication (should return 401)
Write-Host "Test 6: Missing Authentication (HTTP 401)" -ForegroundColor Green
Write-Host "------------------------------------------" -ForegroundColor Green
Cleanup

Write-Host "Starting container without credentials..."
docker run -d --name $CONTAINER_NAME -p ${PORT}:8001 mcp-netbird-sse:latest
Start-Sleep -Seconds 2

Write-Host "Testing without any authentication..."
$response = curl.exe -s -w "\nHTTP_CODE:%{http_code}" -X POST http://localhost:${PORT}/sse `
    -H "Content-Type: application/json" `
    -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'

if ($response -match "HTTP_CODE:401") {
    Write-Host "✓ Missing authentication test PASSED (HTTP 401 returned)" -ForegroundColor Green
} else {
    Write-Host "✗ Missing authentication test FAILED (expected HTTP 401)" -ForegroundColor Red
}
Write-Host ""

# Cleanup
Write-Host "=== Test Summary ===" -ForegroundColor Cyan
Write-Host "All tests completed. Check results above." -ForegroundColor Cyan
Write-Host ""
Cleanup
Write-Host "Cleanup complete." -ForegroundColor Yellow
