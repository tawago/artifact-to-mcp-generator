Your workspace directory is a github repository called tawago/artifact-to-mcp-generator. We are working on the issue/PR user referenced so please use github api to directly fetch the issue. Usually our issue tickets and PRs contains subtasks. We only work on a subtask one at a time. You MUST avoid working on the next subtasks because we close a session and resume on the subtasks in following sessions. Only commit, push and make a pull request when the user specifically asked. 

Below is the current description of this repository.
---

# Smart‑Contract MCP server Generator — Project Plan (Go)

## 1. Project's Objective
Build an **autonomous, chain‑agnostic tool that generates typed MCP servers for any deployed smart‑contract**.  
The core architecture (artifact → IR → code‑gen templates) is identical across languages; this variant outlines how to implement it in **Go**.

### 1.1 Understanding MCP (Model Context Protocol)

MCP (Model Context Protocol) is a standardized protocol for communication between AI models and external tools or services. It enables AI models to interact with various systems, databases, and APIs in a consistent way. MCP servers act as intermediaries that expose functionality to AI models through a well-defined interface.

#### Key MCP Concepts:

1. **Tools**: Functions that an AI model can call to perform specific actions
2. **Prompts**: Pre-defined templates that help guide the AI model's interactions
3. **Server**: The MCP server that implements the protocol and handles requests from AI models
- **Local Data Sources**: Computer files, databases, and services that MCP servers can securely access

Our project aims to automatically generate MCP servers that expose smart contract functionality as tools, making blockchain interactions accessible to LLMs.

#### MCP Server Structure:

A typical MCP server consists of:

1. **Server Definition**: Initializes the MCP server and defines available tools
2. **Tool Implementations**: Functions that implement the actual functionality
3. **Schema Definitions**: Defines the input/output formats for tools
4. **Error Handling**: Standardized error responses

#### Example MCP Server (TypeScript):
```typescript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  ToolSchema,
} from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import { zodToJsonSchema } from "zod-to-json-schema";

// Define tool input schemas
const BalanceOfSchema = z.object({
  address: z.string().describe("Wallet address to check balance"),
});

// Create server
const server = new Server(
  {
    name: "erc20-contract-server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// Register tools
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "balanceOf",
        description: "Get token balance for an address",
        inputSchema: zodToJsonSchema(BalanceOfSchema),
      },
      // Other contract functions...
    ],
  };
});

// Implement tool handlers
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;
  
  if (name === "balanceOf") {
    const validatedArgs = BalanceOfSchema.parse(args);
    // Call the actual contract method
    const balance = await contractInstance.balanceOf(validatedArgs.address);
    return {
      content: [{ type: "text", text: `Balance: ${balance.toString()}` }],
    };
  }
  
  // Handle other tools...
});

// Connect to transport
const transport = new StdioServerTransport();
await server.connect(transport);
```

#### Python Example:
```python
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent
from pydantic import BaseModel

# Define tool input models
class BalanceOfInput(BaseModel):
    address: str

# Create server
server = Server(
    name="erc20-contract-server",
    version="1.0.0",
    capabilities={"tools": {}}
)

# Register tools
@server.list_tools()
async def list_tools():
    return [
        Tool(
            name="balanceOf",
            description="Get token balance for an address",
            input_schema=BalanceOfInput.model_json_schema()
        )
    ]

# Implement tool handlers
@server.call_tool()
async def call_tool(name: str, arguments: dict):
    if name == "balanceOf":
        address = arguments.get("address")
        # Call the actual contract method
        balance = await contract_instance.balance_of(address)
        return [TextContent(type="text", text=f"Balance: {balance}")]
    
    # Handle other tools...

# Start server
async with stdio_server() as (read_stream, write_stream):
    await server.serve(read_stream, write_stream)
```

---

## 2. Language‑Specific Rationale

You chose **Go** for its quick learning curve, blazing‑fast compiles, and effortless static binaries suited to CLI distribution.  
Go's standard library already includes powerful text templating and JSON handling, allowing you to ship an MVP rapidly. Concurrency primitives (goroutines/channels) make scatter‑/gather code‑generation trivial.

---

## 3. Initial Discussions & Design Notes

### 3.1 Our Project: Smart-Contract MCP Server Generator

Our project aims to automatically generate MCP servers for any smart contract by:

1. **Parsing Contract Artifacts**: Reading ABI/IDL files that describe the contract's interface
2. **Normalizing to IR**: Converting the contract definition to a unified Intermediate Representation
3. **Generating Code**: Using templates to create a fully functional MCP server

### 3.2 High‑level architecture
```text
┌───────────────┐   compile/verify   ┌───────────┐  normalize  ┌──────────────┐
│ user repo     │ ────────────────▶ │ ABI/IDL   │ ───────────▶ │ IR (JSON/YAML)│
│ (any lang)    │                   │ (chain fmt)│             └──────┬───────┘
└───────────────┘                   └───────────┘                    ▼
                                                            ┌──────────────┐
                                                            │ MCP server templates│
                                                            │  TS / Go …   │
                                                            └──────────────┘
```

### 3.3 Example Generated MCP Server for ERC-20:

The generator will create a TypeScript MCP server that exposes ERC-20 functions like:

- `balanceOf`: Get the token balance of an account
- `totalSupply`: Get the total token supply
- `transfer`: Transfer tokens to another account
- `allowance`: Check the amount of tokens that an owner allowed to a spender
- `approve`: Approve the spender to spend a specific amount of tokens

### 3.4 Recommended libraries / tooling
* **JSON** — `encoding/json` stdlib (reflection‑based).  
* **YAML** — `gopkg.in/yaml.v3`.  
* **Templating** — `text/template`, with `sprig` for helper funcs.  
* **CLI flags** — `cobra` or vanilla `flag` pkg.  
* **Testing** — `go test`, `testify`.

---

## 4. Detailed, Bite‑Sized Work Sessions

### Phase 1 — Core (EVM only)
1. **Define IR schema (JSON) v0.1** in `Go` structs/enums.  
2. **Write EVM importer** — parse ABI JSON → IR.  
3. **Draft TypeScript MCP server template** using `text/template`.  
4. **Generate sample MCP server** for a public ERC‑20; run basic tests.  
5. **Edge cases**: overloads, payable functions, indexed events.  
6. **CLI wrapper**: `generate-mcp --artifact path/to/abi.json --lang ts`.

### Phase 2 — Multi‑runtime support
7. **Build Solana importer** (Anchor IDL → IR).  
8. **Extend IR for account sizes & seeds** (Solana specific).  
9. **Add Move importer** (Aptos/Sui).  
10. **Generalize serializer helpers** (Borsh, BCS, SCALE).  
11. **Enhance MCP server templates**: language selection flag, events streaming abstraction.  

### Phase 3 — UX polish & extras
12. **Auto‑fetch artifact by chain + address** (explorer APIs).  
13. **Generate docs site** (Docusaurus) from IR.  
14. **Transaction simulation / gas estimator plug‑ins** (EIP‑4337, Sui sponsored tx).  
15. **CI pipeline** for regression tests against example contracts.  
16. **Package registry publishing** (`go install / goreleaser`).

---

## 5. Potential Obstacles & Follow‑up Triggers

* Reflection‑based JSON may miss unknown fields silently; add validation.  
* Fewer maintained libraries for Borsh/SCALE — may require CGO or self‑port.  
* Lack of generics prior to Go 1.18 means some helpers need type assertions.

---