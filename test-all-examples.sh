#!/bin/bash

# Create extracted directory if it doesn't exist
mkdir -p extracted

# Test all examples
echo "Testing all examples..."


# ERC20
echo -e "\n\033[1;34mTesting IERC20 example...\033[0m"
rm -rf mcp-server
jq '.abi' examples/ERC20.json > examples/extracted/ERC20.abi.json
go run cmd/generate-mcp/main.go -a examples/extracted/ERC20.abi.json -n ERC20 --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && npm install && npm run build && npm run test:headless && cd ../

# SimpleVoting
echo -e "\n\033[1;34mTesting SimpleVoting example...\033[0m"
rm -rf mcp-server
jq '.abi' examples/SimpleVoting.json > examples/extracted/SimpleVoting.abi.json
go run cmd/generate-mcp/main.go -a examples/extracted/SimpleVoting.abi.json -n SimpleVoting --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && pwd && npm install && npm run build && npm run test:headless && cd ../

# IERC20
echo -e "\n\033[1;34mTesting IERC20 example...\033[0m"
rm -rf mcp-server
go run cmd/generate-mcp/main.go -a examples/IERC20.json -n IERC20 --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && npm install && npm run build && npm run test:headless && cd ../

# ComplexTypes
echo -e "\n\033[1;34mTesting ComplexTypes example...\033[0m"
rm -rf mcp-server
go run cmd/generate-mcp/main.go -a examples/ComplexTypes.json -n ComplexTypes --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && npm install && npm run build && npm run test:headless && cd ../

# FunctionOverloads
echo -e "\n\033[1;34mTesting FunctionOverloads example...\033[0m"
rm -rf mcp-server
go run cmd/generate-mcp/main.go -a examples/FunctionOverloads.json -n FunctionOverloads --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && npm install && npm run build && npm run test:headless && cd ../

# PayableFunctions
echo -e "\n\033[1;34mTesting PayableFunctions example...\033[0m"
rm -rf mcp-server
go run cmd/generate-mcp/main.go -a examples/PayableFunctions.json -n PayableFunctions --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
cd mcp-server && npm install && npm run build && npm run test:headless && cd ../

echo -e "\n\033[1;32mAll examples tested successfully!\033[0m"