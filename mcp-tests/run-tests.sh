#!/bin/bash
set -e

# Build the MCP server
echo "Building MCP server..."
cd ../mcp-server
npm run build

# Run the tests
echo "Running tests..."
cd ../tests
npm test

echo "Tests completed!"