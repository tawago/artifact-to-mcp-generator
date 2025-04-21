#!/bin/bash
set -e

# Build the MCP server
echo "Building MCP server..."
cd ../mcp-server
npm install
npm run build

# Run the tests
echo "Running tests..."
cd ../mcp-tests
npm run test:manual

echo "Tests completed!"