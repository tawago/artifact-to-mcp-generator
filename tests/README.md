# MCP Generator Tests

This directory contains tests for the MCP Generator.

## Test Files

- `erc20_mcp_test.js`: Tests the generation and functionality of an MCP server for an ERC-20 contract.
- `event_subscription_test.js`: Tests event subscription for an ERC-20 contract.
- `run_integration_tests.sh`: Shell script for running basic integration tests.
- `mcp_inspector.js`: Interactive tool for inspecting and interacting with MCP servers.

## Running the Tests

1. Install dependencies:
   ```bash
   npm install
   ```

2. Run the basic functionality test:
   ```bash
   npm test
   ```

3. Run the event subscription test:
   ```bash
   npm run test:events
   ```

4. Run the shell-based integration tests:
   ```bash
   ./run_integration_tests.sh
   ```

5. Use the MCP Inspector for interactive testing:
   ```bash
   node mcp_inspector.js
   ```

## Test Details

### ERC-20 MCP Test

This test:
1. Generates an MCP server for the USDC contract
2. Installs dependencies
3. Builds the server
4. Tests all view functions (name, symbol, decimals, totalSupply, balanceOf, allowance)

### Event Subscription Test

This test:
1. Enhances the generated MCP server with event subscription support
2. Builds the enhanced server
3. Subscribes to Transfer events
4. Waits for events (in a real test environment, we would trigger events)

## MCP Inspector

The MCP Inspector is an interactive tool for testing and debugging MCP servers. It allows you to:

1. List available tools
2. Call tools with arguments
3. Subscribe to events
4. View server responses and events in real-time

### Usage

```bash
node mcp_inspector.js [options]
```

### Options

- `--server <path>`: Path to the MCP server (default: ../mcp-server/dist/server.js)
- `--rpc <url>`: RPC URL (default: https://eth.llamarpc.com)
- `--ws-rpc <url>`: WebSocket RPC URL (default: wss://eth.llamarpc.com)
- `--contract <address>`: Contract address (default: USDC address)
- `--help`: Show help message

### Commands

Once the inspector is running, you can use the following commands:

- `list`: List available tools
- `call <tool> [args]`: Call a tool with arguments (JSON format)
- `subscribe <event>`: Subscribe to an event
- `exit`: Exit the inspector
- `help`: Show help message

### Examples

```
list
call name
call balanceOf {"_owner":"0xf584f8728b874a6a5c7a8d4d387c9aae9172d621"}
subscribe Transfer
```

## Configuration

The tests use the following configuration:
- RPC URL: https://eth.llamarpc.com
- WebSocket RPC URL: wss://eth.llamarpc.com
- Contract Address: USDC on Ethereum mainnet (0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48)
- Test Wallet Address: 0xf584f8728b874a6a5c7a8d4d387c9aae9172d621 (a random whale address with USDC)

You can modify these values in the test files if needed.