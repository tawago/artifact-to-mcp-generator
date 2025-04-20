# Smart‑Contract MCP server Generator — Project Plan (Go)

## 1. Project's Objective
Build an **autonomous, chain‑agnostic tool that generates typed MCP servers for any deployed smart‑contract**.  
The core architecture (artifact → IR → code‑gen templates) is identical across languages; this variant outlines how to implement it in **Go**.

### 1.1 Understanding MCP (Model Context Protocol)

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to Large Language Models (LLMs). It enables LLMs to securely access tools and data sources through a standardized interface.

MCP follows a client-server architecture:
- **MCP Hosts**: Programs like Claude Desktop, IDEs, or AI tools that want to access data through MCP
- **MCP Clients**: Protocol clients that maintain 1:1 connections with servers
- **MCP Servers**: Lightweight programs that expose specific capabilities through the standardized Model Context Protocol
- **Local Data Sources**: Computer files, databases, and services that MCP servers can securely access
- **Remote Services**: External systems available over the internet that MCP servers can connect to

Key components of an MCP server include:
1. **Tools**: Functions that LLMs can call to perform actions
2. **Resources**: Data that can be accessed by LLMs
3. **Prompts**: Reusable prompt templates for LLMs

Our project aims to automatically generate MCP servers that expose smart contract functionality as tools, making blockchain interactions accessible to LLMs.

### 1.2 MCP Server Structure and Implementation

MCP servers are typically implemented using the official SDKs (TypeScript, Python, Java, Kotlin, C#). Here's a simplified example of an MCP server structure:

#### TypeScript Example:
```typescript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { CallToolRequestSchema, ToolSchema } from "@modelcontextprotocol/sdk/types.js";
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

### 3.1 Artifact‑first strategy
Same as the master plan: consume ABI/IDL/schema artifacts, normalize into a unified IR, then render MCP servers.

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

### 3.3 Recommended libraries / tooling
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