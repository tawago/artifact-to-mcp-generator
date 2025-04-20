#!/bin/bash

# Integration test script for MCP Generator

# Set up error handling
set -e
trap 'echo "Error: Command failed with exit code $?"; exit 1' ERR

# Print header
echo "====================================="
echo "MCP Generator Integration Tests"
echo "====================================="

# Step 1: Generate MCP server
echo -e "\n[Step 1] Generating MCP server..."
cd ..
go run cmd/generate-mcp/main.go --artifact examples/erc20.json --output mcp-server --name "USDC" --address "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

# Step 2: Install dependencies
echo -e "\n[Step 2] Installing dependencies..."
cd mcp-server
npm install

# Step 3: Build the server
echo -e "\n[Step 3] Building the server..."
npm run build

# Step 4: Start the server in the background
echo -e "\n[Step 4] Starting the server..."
export RPC_URL="https://eth.llamarpc.com"
export CONTRACT_ADDRESS="0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
node dist/server.js > server.log 2>&1 &
SERVER_PID=$!

# Give the server time to start
sleep 2

# Step 5: Test basic functionality
echo -e "\n[Step 5] Testing basic functionality..."

# Test name function
echo "Testing name function..."
echo '{"jsonrpc":"2.0","id":"1","method":"mcp.call_tool","params":{"name":"name","arguments":{}}}' | nc -N localhost 52745 | grep -q "USDC" && echo "✅ Name function works" || echo "❌ Name function failed"

# Test symbol function
echo "Testing symbol function..."
echo '{"jsonrpc":"2.0","id":"2","method":"mcp.call_tool","params":{"name":"symbol","arguments":{}}}' | nc -N localhost 52745 | grep -q "USDC" && echo "✅ Symbol function works" || echo "❌ Symbol function failed"

# Test decimals function
echo "Testing decimals function..."
echo '{"jsonrpc":"2.0","id":"3","method":"mcp.call_tool","params":{"name":"decimals","arguments":{}}}' | nc -N localhost 52745 | grep -q "6" && echo "✅ Decimals function works" || echo "❌ Decimals function failed"

# Test totalSupply function
echo "Testing totalSupply function..."
echo '{"jsonrpc":"2.0","id":"4","method":"mcp.call_tool","params":{"name":"totalSupply","arguments":{}}}' | nc -N localhost 52745 | grep -q "\"text\"" && echo "✅ TotalSupply function works" || echo "❌ TotalSupply function failed"

# Test balanceOf function
echo "Testing balanceOf function..."
echo '{"jsonrpc":"2.0","id":"5","method":"mcp.call_tool","params":{"name":"balanceOf","arguments":{"_owner":"0xf584f8728b874a6a5c7a8d4d387c9aae9172d621"}}}' | nc -N localhost 52745 | grep -q "\"text\"" && echo "✅ BalanceOf function works" || echo "❌ BalanceOf function failed"

# Step 6: Clean up
echo -e "\n[Step 6] Cleaning up..."
kill $SERVER_PID

# Print summary
echo -e "\n====================================="
echo "Integration tests completed!"
echo "====================================="

# Return to tests directory
cd ../tests