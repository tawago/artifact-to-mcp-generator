#!/bin/bash
set -e

# Build the MCP server
echo "Building MCP server..."
cd /workspace/mcp-server
npm run build

# Run the tests
echo "Running tests..."
cd /workspace/tests
npm test

echo "Tests completed!"