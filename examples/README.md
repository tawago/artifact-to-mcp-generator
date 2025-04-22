# MCP Server Generator Examples

This directory contains example smart contract artifacts that can be used with the MCP server generator.

## Available Examples

### 1. SimpleVoting.json

A simple voting contract that allows users to vote for choices and view results.

**Usage:**
```bash
# Extract the ABI first
jq '.abi' SimpleVoting.json > extracted/SimpleVoting.abi.json

# Generate MCP server
go run cmd/generate-mcp/main.go -a examples/extracted/SimpleVoting.abi.json -n SimpleVoting
```

### 2. IERC20.json

The standard ERC-20 token interface with functions like `balanceOf`, `transfer`, and `approve`.

**Usage:**
```bash
go run cmd/generate-mcp/main.go -a examples/IERC20.json -n IERC20
```

### 3. ComplexTypes.json

Demonstrates handling of complex types including arrays, fixed-size arrays, tuples (structs), and nested structs.

**Usage:**
```bash
go run cmd/generate-mcp/main.go -a examples/ComplexTypes.json -n ComplexTypes
```

### 4. FunctionOverloads.json

Demonstrates handling of function overloads (multiple functions with the same name but different parameters).

**Usage:**
```bash
go run cmd/generate-mcp/main.go -a examples/FunctionOverloads.json -n FunctionOverloads
```

### 5. PayableFunctions.json

Demonstrates handling of payable functions, receive function, and fallback function.

**Usage:**
```bash
go run cmd/generate-mcp/main.go -a examples/PayableFunctions.json -n PayableFunctions
```

## Notes on ABI Format

The MCP server generator expects an array of ABI items. If your ABI is wrapped in an object (like Hardhat artifacts), you need to extract just the ABI array:

```bash
# Extract ABI from a Hardhat artifact
jq '.abi' YourContract.json > extracted/YourContract.abi.json
```

## Generated MCP Server

The generated MCP server will:

1. Create TypeScript types for all function parameters and return values
2. Handle complex types like arrays and structs
3. Properly handle function overloads with unique names
4. Support payable functions with ETH value passing
5. Include indexed event parameters for efficient filtering

## Running the Generated Server

```bash
cd mcp-server
npm install
npm run build
npm start
```

You can configure the server with environment variables:
- `RPC_URL`: Ethereum RPC URL (default: https://eth.llamarpc.com)
- `CONTRACT_ADDRESS`: Contract address