#!/usr/bin/env node

// MCP Inspector - A simple tool to inspect and interact with MCP servers
import { createInterface } from 'readline';
import { spawn } from 'child_process';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';

// Get the current directory
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Default configuration
const config = {
  serverPath: join(__dirname, '../mcp-server/dist/server.js'),
  rpcUrl: 'https://eth.llamarpc.com',
  contractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', // USDC
  wsRpcUrl: 'wss://eth.llamarpc.com'
};

// Parse command line arguments
const args = process.argv.slice(2);
for (let i = 0; i < args.length; i++) {
  if (args[i] === '--server' && i + 1 < args.length) {
    config.serverPath = args[i + 1];
    i++;
  } else if (args[i] === '--rpc' && i + 1 < args.length) {
    config.rpcUrl = args[i + 1];
    i++;
  } else if (args[i] === '--ws-rpc' && i + 1 < args.length) {
    config.wsRpcUrl = args[i + 1];
    i++;
  } else if (args[i] === '--contract' && i + 1 < args.length) {
    config.contractAddress = args[i + 1];
    i++;
  } else if (args[i] === '--help') {
    console.log(`
MCP Inspector - A simple tool to inspect and interact with MCP servers

Usage:
  node mcp_inspector.js [options]

Options:
  --server <path>      Path to the MCP server (default: ../mcp-server/dist/server.js)
  --rpc <url>          RPC URL (default: https://eth.llamarpc.com)
  --ws-rpc <url>       WebSocket RPC URL (default: wss://eth.llamarpc.com)
  --contract <address> Contract address (default: USDC address)
  --help               Show this help message

Commands (once running):
  list                 List available tools
  call <tool> [args]   Call a tool with arguments (JSON format)
  subscribe <event>    Subscribe to an event
  exit                 Exit the inspector
`);
    process.exit(0);
  }
}

// Start the MCP server
console.log('Starting MCP server...');
const serverProcess = spawn('node', [config.serverPath], {
  env: {
    ...process.env,
    RPC_URL: config.rpcUrl,
    WS_RPC_URL: config.wsRpcUrl,
    CONTRACT_ADDRESS: config.contractAddress
  }
});

// Create readline interfaces
const serverRl = createInterface({
  input: serverProcess.stdout,
  crlfDelay: Infinity
});

const userRl = createInterface({
  input: process.stdin,
  output: process.stdout,
  prompt: 'mcp> '
});

// Handle server output
serverRl.on('line', (line) => {
  try {
    const data = JSON.parse(line);
    if (data.method === 'mcp.event') {
      console.log('\n🔔 Event received:');
      console.log(`Event: ${data.params.event}`);
      console.log(`Data: ${JSON.stringify(data.params.data, null, 2)}`);
    } else {
      console.log('\n📥 Server response:');
      console.log(JSON.stringify(data, null, 2));
    }
    userRl.prompt();
  } catch (error) {
    console.log(`\n🖥️  Server: ${line}`);
    userRl.prompt();
  }
});

// Handle server errors
serverProcess.stderr.on('data', (data) => {
  console.error(`\n❌ Server error: ${data}`);
  userRl.prompt();
});

// Handle user input
userRl.prompt();
userRl.on('line', (line) => {
  const args = line.trim().split(' ');
  const command = args[0].toLowerCase();
  
  if (command === 'exit') {
    console.log('Exiting...');
    serverProcess.kill();
    process.exit(0);
  } else if (command === 'list') {
    const request = {
      jsonrpc: '2.0',
      id: Date.now().toString(),
      method: 'mcp.list_tools',
      params: {}
    };
    serverProcess.stdin.write(JSON.stringify(request) + '\n');
  } else if (command === 'call' && args.length >= 2) {
    const toolName = args[1];
    let toolArgs = {};
    
    if (args.length > 2) {
      try {
        toolArgs = JSON.parse(args.slice(2).join(' '));
      } catch (error) {
        console.error('Error parsing arguments. Please provide valid JSON.');
        userRl.prompt();
        return;
      }
    }
    
    const request = {
      jsonrpc: '2.0',
      id: Date.now().toString(),
      method: 'mcp.call_tool',
      params: {
        name: toolName,
        arguments: toolArgs
      }
    };
    serverProcess.stdin.write(JSON.stringify(request) + '\n');
  } else if (command === 'subscribe' && args.length >= 2) {
    const eventName = args[1];
    let filter = {};
    
    if (args.length > 2) {
      try {
        filter = JSON.parse(args.slice(2).join(' '));
      } catch (error) {
        console.error('Error parsing filter. Please provide valid JSON.');
        userRl.prompt();
        return;
      }
    }
    
    const request = {
      jsonrpc: '2.0',
      id: Date.now().toString(),
      method: 'mcp.subscribe_events',
      params: {
        event: eventName,
        filter: filter
      }
    };
    serverProcess.stdin.write(JSON.stringify(request) + '\n');
  } else if (command === 'help') {
    console.log(`
Available commands:
  list                 List available tools
  call <tool> [args]   Call a tool with arguments (JSON format)
  subscribe <event>    Subscribe to an event
  exit                 Exit the inspector
  help                 Show this help message

Examples:
  list
  call name
  call balanceOf {"_owner":"0xf584f8728b874a6a5c7a8d4d387c9aae9172d621"}
  subscribe Transfer
`);
  } else {
    console.log('Unknown command. Type "help" for available commands.');
  }
  
  userRl.prompt();
});