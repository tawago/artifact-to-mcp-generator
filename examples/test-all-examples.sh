#!/bin/bash

# Create extracted directory if it doesn't exist
mkdir -p extracted

# Extract the ABI from SimpleVoting.json
echo "Extracting ABI from SimpleVoting.json..."
jq '.abi' SimpleVoting.json > extracted/SimpleVoting.abi.json

# Test all examples
echo "Testing all examples..."

# SimpleVoting
echo -e "\n\033[1;34mTesting SimpleVoting example...\033[0m"
rm -rf ../mcp-server
go run ../cmd/generate-mcp/main.go -a extracted/SimpleVoting.abi.json -n SimpleVoting -o ../mcp-server

# IERC20
echo -e "\n\033[1;34mTesting IERC20 example...\033[0m"
rm -rf ../mcp-server
go run ../cmd/generate-mcp/main.go -a IERC20.json -n IERC20 -o ../mcp-server

# ComplexTypes
echo -e "\n\033[1;34mTesting ComplexTypes example...\033[0m"
rm -rf ../mcp-server
go run ../cmd/generate-mcp/main.go -a ComplexTypes.json -n ComplexTypes -o ../mcp-server

# FunctionOverloads
echo -e "\n\033[1;34mTesting FunctionOverloads example...\033[0m"
rm -rf ../mcp-server
go run ../cmd/generate-mcp/main.go -a FunctionOverloads.json -n FunctionOverloads -o ../mcp-server

# PayableFunctions
echo -e "\n\033[1;34mTesting PayableFunctions example...\033[0m"
rm -rf ../mcp-server
go run ../cmd/generate-mcp/main.go -a PayableFunctions.json -n PayableFunctions -o ../mcp-server

echo -e "\n\033[1;32mAll examples tested successfully!\033[0m"