// ERC-20 MCP Server Integration Test
import { ethers } from 'ethers';
import { spawn } from 'child_process';
import { createInterface } from 'readline';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';

// Get the current directory
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// USDC contract address on Ethereum mainnet
const USDC_ADDRESS = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48';

// Test configuration
const config = {
  rpcUrl: 'https://eth.llamarpc.com',
  contractAddress: USDC_ADDRESS,
  testWalletAddress: '0xf584f8728b874a6a5c7a8d4d387c9aae9172d621' // Random whale address with USDC
};

// Helper function to send a request to the MCP server
async function sendMcpRequest(serverProcess, request) {
  return new Promise((resolve, reject) => {
    const requestStr = JSON.stringify(request);
    console.log(`Sending request: ${requestStr}`);
    
    serverProcess.stdin.write(requestStr + '\n');
    
    // Set up a listener for the response
    const responseListener = (data) => {
      try {
        const response = JSON.parse(data.toString());
        console.log(`Received response: ${JSON.stringify(response, null, 2)}`);
        resolve(response);
      } catch (error) {
        console.error('Error parsing response:', error);
        reject(error);
      }
    };
    
    // Listen for one response
    serverProcess.stdout.once('data', responseListener);
  });
}

// Main test function
async function runTests() {
  console.log('Starting ERC-20 MCP Server tests...');
  
  // Step 1: Generate the MCP server
  console.log('Generating MCP server for USDC contract...');
  try {
    // Remove existing mcp-server directory if it exists
    if (fs.existsSync(join(__dirname, '../mcp-server'))) {
      fs.rmSync(join(__dirname, '../mcp-server'), { recursive: true, force: true });
    }
    
    // Generate the MCP server
    const generateProcess = spawn('go', [
      'run', 
      'cmd/generate-mcp/main.go', 
      '--artifact', 'examples/erc20.json',
      '--output', 'mcp-server',
      '--name', 'USDC',
      '--address', USDC_ADDRESS
    ], { cwd: join(__dirname, '..') });
    
    // Wait for the generation to complete
    await new Promise((resolve, reject) => {
      generateProcess.on('close', (code) => {
        if (code === 0) {
          console.log('MCP server generated successfully');
          resolve();
        } else {
          reject(new Error(`MCP server generation failed with code ${code}`));
        }
      });
      
      generateProcess.stdout.on('data', (data) => {
        console.log(`Generator output: ${data}`);
      });
      
      generateProcess.stderr.on('data', (data) => {
        console.error(`Generator error: ${data}`);
      });
    });
    
    // Step 2: Install dependencies
    console.log('Installing MCP server dependencies...');
    const npmProcess = spawn('npm', ['install'], { cwd: join(__dirname, '../mcp-server') });
    
    await new Promise((resolve, reject) => {
      npmProcess.on('close', (code) => {
        if (code === 0) {
          console.log('Dependencies installed successfully');
          resolve();
        } else {
          reject(new Error(`Dependency installation failed with code ${code}`));
        }
      });
    });
    
    // Step 3: Build the MCP server
    console.log('Building MCP server...');
    const buildProcess = spawn('npm', ['run', 'build'], { cwd: join(__dirname, '../mcp-server') });
    
    await new Promise((resolve, reject) => {
      buildProcess.on('close', (code) => {
        if (code === 0) {
          console.log('MCP server built successfully');
          resolve();
        } else {
          reject(new Error(`MCP server build failed with code ${code}`));
        }
      });
    });
    
    // Step 4: Start the MCP server
    console.log('Starting MCP server...');
    const serverProcess = spawn('node', ['dist/server.js'], { 
      cwd: join(__dirname, '../mcp-server'),
      env: {
        ...process.env,
        RPC_URL: config.rpcUrl,
        CONTRACT_ADDRESS: config.contractAddress
      }
    });
    
    // Create readline interface for reading server output
    const rl = createInterface({
      input: serverProcess.stdout,
      crlfDelay: Infinity
    });
    
    // Wait for the server to start
    await new Promise((resolve) => setTimeout(resolve, 2000));
    
    // Step 5: Test listing tools
    console.log('Testing tool listing...');
    const listToolsRequest = {
      jsonrpc: '2.0',
      id: '1',
      method: 'mcp.list_tools',
      params: {}
    };
    
    const listToolsResponse = await sendMcpRequest(serverProcess, listToolsRequest);
    console.log('Available tools:', listToolsResponse.result.tools);
    
    // Step 6: Test view functions
    console.log('Testing view functions...');
    
    // Test name function
    const nameRequest = {
      jsonrpc: '2.0',
      id: '2',
      method: 'mcp.call_tool',
      params: {
        name: 'name',
        arguments: {}
      }
    };
    
    const nameResponse = await sendMcpRequest(serverProcess, nameRequest);
    console.log('Token name:', nameResponse.result.content[0].text);
    
    // Test symbol function
    const symbolRequest = {
      jsonrpc: '2.0',
      id: '3',
      method: 'mcp.call_tool',
      params: {
        name: 'symbol',
        arguments: {}
      }
    };
    
    const symbolResponse = await sendMcpRequest(serverProcess, symbolRequest);
    console.log('Token symbol:', symbolResponse.result.content[0].text);
    
    // Test decimals function
    const decimalsRequest = {
      jsonrpc: '2.0',
      id: '4',
      method: 'mcp.call_tool',
      params: {
        name: 'decimals',
        arguments: {}
      }
    };
    
    const decimalsResponse = await sendMcpRequest(serverProcess, decimalsRequest);
    console.log('Token decimals:', decimalsResponse.result.content[0].text);
    
    // Test totalSupply function
    const totalSupplyRequest = {
      jsonrpc: '2.0',
      id: '5',
      method: 'mcp.call_tool',
      params: {
        name: 'totalSupply',
        arguments: {}
      }
    };
    
    const totalSupplyResponse = await sendMcpRequest(serverProcess, totalSupplyRequest);
    console.log('Total supply:', totalSupplyResponse.result.content[0].text);
    
    // Test balanceOf function
    const balanceOfRequest = {
      jsonrpc: '2.0',
      id: '6',
      method: 'mcp.call_tool',
      params: {
        name: 'balanceOf',
        arguments: {
          _owner: config.testWalletAddress
        }
      }
    };
    
    const balanceOfResponse = await sendMcpRequest(serverProcess, balanceOfRequest);
    console.log(`Balance of ${config.testWalletAddress}:`, balanceOfResponse.result.content[0].text);
    
    // Test allowance function
    const allowanceRequest = {
      jsonrpc: '2.0',
      id: '7',
      method: 'mcp.call_tool',
      params: {
        name: 'allowance',
        arguments: {
          _owner: config.testWalletAddress,
          _spender: '0x0000000000000000000000000000000000000000'
        }
      }
    };
    
    const allowanceResponse = await sendMcpRequest(serverProcess, allowanceRequest);
    console.log('Allowance:', allowanceResponse.result.content[0].text);
    
    // Step 7: Clean up
    console.log('Tests completed successfully!');
    serverProcess.kill();
    
  } catch (error) {
    console.error('Test failed:', error);
    process.exit(1);
  }
}

// Run the tests
runTests();