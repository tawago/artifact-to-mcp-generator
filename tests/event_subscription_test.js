// ERC-20 MCP Server Event Subscription Test
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
  wsRpcUrl: 'wss://eth.llamarpc.com' // WebSocket URL for event subscription
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
  console.log('Starting ERC-20 MCP Server Event Subscription tests...');
  
  try {
    // Step 1: Enhance the MCP server with event subscription support
    console.log('Enhancing MCP server with event subscription...');
    
    // Create a modified server.ts file with event subscription support
    const serverTsPath = join(__dirname, '../mcp-server/src/server.ts');
    let serverTsContent = fs.readFileSync(serverTsPath, 'utf8');
    
    // Add event subscription capability
    const eventSubscriptionCode = `
// Event subscription support
server.setRequestHandler(
  {
    method: "mcp.subscribe_events",
    params: z.object({
      event: z.string(),
      filter: z.record(z.any()).optional(),
    }),
  },
  async (request) => {
    const { event, filter } = request.params;
    
    try {
      // Create a WebSocket provider for events
      const wsProvider = new ethers.WebSocketProvider(process.env.WS_RPC_URL || "wss://eth.llamarpc.com");
      const wsContract = new ethers.Contract(contractAddress, contractABI, wsProvider);
      
      console.log(\`Subscribing to \${event} events with filter: \${JSON.stringify(filter)}\`);
      
      // Set up event listener
      wsContract.on(event, (...args) => {
        const eventData = args.slice(0, -1); // Remove the event object
        console.log(\`Event received: \${event}\`, eventData);
        
        // Send event notification
        server.sendNotification("mcp.event", {
          event,
          data: eventData.map(item => 
            typeof item === 'bigint' ? item.toString() : item
          ),
        });
      });
      
      return {
        content: [
          {
            type: "text",
            text: \`Subscribed to \${event} events\`,
          },
        ],
      };
    } catch (error) {
      console.error("Error subscribing to events:", error);
      throw new Error("Error subscribing to " + event + ": " + error.message);
    }
  }
);`;
    
    // Insert the event subscription code before the transport connection
    const transportIndex = serverTsContent.indexOf('// Connect to transport');
    if (transportIndex !== -1) {
      serverTsContent = 
        serverTsContent.slice(0, transportIndex) + 
        eventSubscriptionCode + 
        '\n\n' + 
        serverTsContent.slice(transportIndex);
      
      fs.writeFileSync(serverTsPath, serverTsContent);
      console.log('Added event subscription support to server.ts');
    } else {
      throw new Error('Could not find transport connection in server.ts');
    }
    
    // Step 2: Build the enhanced MCP server
    console.log('Building enhanced MCP server...');
    const buildProcess = spawn('npm', ['run', 'build'], { cwd: join(__dirname, '../mcp-server') });
    
    await new Promise((resolve, reject) => {
      buildProcess.on('close', (code) => {
        if (code === 0) {
          console.log('Enhanced MCP server built successfully');
          resolve();
        } else {
          reject(new Error(`Enhanced MCP server build failed with code ${code}`));
        }
      });
    });
    
    // Step 3: Start the enhanced MCP server
    console.log('Starting enhanced MCP server...');
    const serverProcess = spawn('node', ['dist/server.js'], { 
      cwd: join(__dirname, '../mcp-server'),
      env: {
        ...process.env,
        RPC_URL: config.rpcUrl,
        WS_RPC_URL: config.wsRpcUrl,
        CONTRACT_ADDRESS: config.contractAddress
      }
    });
    
    // Create readline interface for reading server output
    const rl = createInterface({
      input: serverProcess.stdout,
      crlfDelay: Infinity
    });
    
    // Log server output
    rl.on('line', (line) => {
      console.log(`Server: ${line}`);
    });
    
    // Wait for the server to start
    await new Promise((resolve) => setTimeout(resolve, 2000));
    
    // Step 4: Subscribe to Transfer events
    console.log('Subscribing to Transfer events...');
    const subscribeRequest = {
      jsonrpc: '2.0',
      id: '1',
      method: 'mcp.subscribe_events',
      params: {
        event: 'Transfer',
        filter: {}
      }
    };
    
    const subscribeResponse = await sendMcpRequest(serverProcess, subscribeRequest);
    console.log('Subscription response:', subscribeResponse);
    
    // Step 5: Wait for events (in a real test, we would trigger events)
    console.log('Waiting for Transfer events (will timeout after 30 seconds)...');
    console.log('Note: In a real test environment, we would trigger events.');
    console.log('For this test, we are just listening for any Transfer events that happen on the network.');
    
    // Wait for 30 seconds to potentially receive events
    await new Promise((resolve) => setTimeout(resolve, 30000));
    
    // Step 6: Clean up
    console.log('Event subscription test completed');
    serverProcess.kill();
    
  } catch (error) {
    console.error('Test failed:', error);
    process.exit(1);
  }
}

// Run the tests
runTests();