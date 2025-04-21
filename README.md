# Smart-Contract MCP Server Generator

A tool that generates typed MCP (Model Context Protocol) servers for any deployed smart contract.

## Overview

This project provides a chain-agnostic tool that automatically generates MCP servers from smart contract artifacts (ABIs, IDLs, etc.). The generated servers expose smart contract functionality as tools that can be used by Large Language Models (LLMs) through the Model Context Protocol.

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to Large Language Models. It enables LLMs to securely access tools and data sources through a standardized interface.

MCP follows a client-server architecture:
- **MCP Hosts**: Programs like Claude Desktop, IDEs, or AI tools that want to access data through MCP
- **MCP Clients**: Protocol clients that maintain 1:1 connections with servers
- **MCP Servers**: Lightweight programs that expose specific capabilities through the standardized Model Context Protocol

## Features

- Generate typed MCP servers from smart contract ABIs/IDLs
- Support for multiple blockchain platforms (Ethereum, Solana, etc.)
- Customizable templates for different programming languages
- CLI tool for easy integration into development workflows

## Installation

```bash
go install github.com/openhands/mcp-generator/cmd/generate-mcp@latest
```

## Usage

```bash
# Generate a TypeScript MCP server from an Ethereum ABI
generate-mcp --artifact path/to/abi.json --lang ts --output ./my-mcp-server

# Generate a Python MCP server from a Solana IDL
generate-mcp --artifact path/to/idl.json --chain solana --lang python --output ./my-mcp-server

```

## Testing

The project includes end-to-end tests to verify that the generated MCP servers work correctly with the MCP Inspector.

### Running Tests

```bash
# current command to test generate
go run cmd/generate-mcp/main.go -a examples/erc20.json --name "ERC20 Token" --address 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

# Navigate to the tests directory
cd mcp-tests

# Install dependencies
npm install

# Run the tests
./run-tests.sh
```

The tests use Playwright to automate interactions with the MCP Inspector UI and verify that:
1. The MCP server connects successfully
2. All tools are listed correctly
3. Each tool can be executed with valid parameters
4. Error handling works as expected

See the [tests/README.md](tests/README.md) for more details.

## Development Status

This project is currently in active development. See the [Phase 1 issue](https://github.com/openhands/mcp-generator/issues/1) for current progress.

## License

MIT