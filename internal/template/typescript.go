package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/openhands/mcp-generator/internal/ir"
)

// TypeScriptTemplateRenderer renders TypeScript MCP server templates
type TypeScriptTemplateRenderer struct{}

// NewTypeScriptTemplateRenderer creates a new TypeScript template renderer
func NewTypeScriptTemplateRenderer() *TypeScriptTemplateRenderer {
	return &TypeScriptTemplateRenderer{}
}

// Render generates a TypeScript MCP server from the IR
func (r *TypeScriptTemplateRenderer) Render(contract *ir.ContractIR) (map[string][]byte, error) {
	files := make(map[string][]byte)

	// Generate package.json
	packageJSON, err := r.renderPackageJSON(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to render package.json: %w", err)
	}
	files["package.json"] = packageJSON

	// Generate tsconfig.json
	tsconfigJSON, err := r.renderTSConfigJSON(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to render tsconfig.json: %w", err)
	}
	files["tsconfig.json"] = tsconfigJSON

	// Generate main server file
	serverTS, err := r.renderServerTS(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to render server.ts: %w", err)
	}
	files["src/server.ts"] = serverTS

	// Generate README.md
	readme, err := r.renderReadme(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to render README.md: %w", err)
	}
	files["README.md"] = readme

	return files, nil
}

// renderPackageJSON generates the package.json file
func (r *TypeScriptTemplateRenderer) renderPackageJSON(contract *ir.ContractIR) ([]byte, error) {
	tmpl, err := template.New("package.json").Funcs(sprig.FuncMap()).Parse(packageJSONTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, contract)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// renderTSConfigJSON generates the tsconfig.json file
func (r *TypeScriptTemplateRenderer) renderTSConfigJSON(contract *ir.ContractIR) ([]byte, error) {
	tmpl, err := template.New("tsconfig.json").Funcs(sprig.FuncMap()).Parse(tsconfigJSONTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, contract)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// renderServerTS generates the main server.ts file
func (r *TypeScriptTemplateRenderer) renderServerTS(contract *ir.ContractIR) ([]byte, error) {
	tmpl, err := template.New("server.ts").Funcs(sprig.FuncMap()).Parse(serverTSTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, contract)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// renderReadme generates the README.md file
func (r *TypeScriptTemplateRenderer) renderReadme(contract *ir.ContractIR) ([]byte, error) {
	tmpl, err := template.New("README.md").Funcs(sprig.FuncMap()).Parse(readmeTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, contract)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// packageJSONTemplate is the template for package.json
const packageJSONTemplate = `{
  "name": "{{ .Metadata.Name | lower | replace " " "-" }}-mcp-server",
  "version": "1.0.0",
  "description": "MCP server for {{ .Metadata.Name }} smart contract",
  "main": "dist/server.js",
  "type": "module",
  "scripts": {
    "build": "tsc",
    "start": "node dist/server.js",
    "dev": "tsc -w"
  },
  "dependencies": {
    "@modelcontextprotocol/sdk": "^0.1.0",
    "ethers": "^6.7.1",
    "zod": "^3.22.2",
    "zod-to-json-schema": "^3.21.4"
  },
  "devDependencies": {
    "@types/node": "^20.5.9",
    "typescript": "^5.2.2"
  }
}
`

// tsconfigJSONTemplate is the template for tsconfig.json
const tsconfigJSONTemplate = `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "esModuleInterop": true,
    "strict": true,
    "outDir": "dist",
    "sourceMap": true,
    "declaration": true,
    "skipLibCheck": true
  },
  "include": ["src/**/*"]
}
`

// serverTSTemplate is the template for server.ts
const serverTSTemplate = `import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { 
  CallToolRequestSchema, 
  ListToolsRequestSchema, 
  Tool, 
  ToolInput, 
  TextContent 
} from "@modelcontextprotocol/sdk/types.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/transports/stdio.js";
import { ethers } from "ethers";
import { z } from "zod";
import { zodToJsonSchema } from "zod-to-json-schema";

// Define tool names enum
enum ToolName {
{{- range .Functions }}
{{- if not .IsConstructor }}
{{- if not .IsFallback }}
{{- if not .IsReceive }}
  {{ .Name | upper }} = "{{ .Name }}",
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

// Define input schemas for each function
{{- range .Functions }}
{{- if not .IsConstructor }}
{{- if not .IsFallback }}
{{- if not .IsReceive }}
const {{ .Name | title }}Schema = z.object({
{{- range .Inputs }}
  {{ .Name }}: z.{{ template "zodType" .Type }}.describe("{{ .Description }}"),
{{- end }}
});

{{- end }}
{{- end }}
{{- end }}
{{- end }}

// Initialize the contract
async function initializeContract() {
  // Connect to provider (replace with your preferred provider)
  const provider = new ethers.JsonRpcProvider(process.env.RPC_URL || "https://eth.llamarpc.com");
  
  // Contract address (can be overridden with environment variable)
  const contractAddress = process.env.CONTRACT_ADDRESS || "{{ .Metadata.Address }}";
  
  // Contract ABI
  const contractABI = [
    {{- range .Functions }}
    {
      "name": "{{ .Name }}",
      "type": "function",
      "inputs": [
        {{- range .Inputs }}
        {
          "name": "{{ .Name }}",
          "type": "{{ .Type.BaseType }}"
          {{- if .Type.IsArray }},
          "components": []
          {{- end }}
        }{{ if not (last $.Functions.Inputs .) }},{{ end }}
        {{- end }}
      ],
      "outputs": [
        {{- range .Outputs }}
        {
          "name": "{{ .Name }}",
          "type": "{{ .Type.BaseType }}"
          {{- if .Type.IsArray }},
          "components": []
          {{- end }}
        }{{ if not (last $.Functions.Outputs .) }},{{ end }}
        {{- end }}
      ],
      "stateMutability": "{{ .StateMutability }}"
    }{{ if not (last $.Functions .) }},{{ end }}
    {{- end }}
  ];
  
  // Create contract instance
  return new ethers.Contract(contractAddress, contractABI, provider);
}

async function main() {
  // Initialize the contract
  const contract = await initializeContract();
  
  // Create MCP server
  const server = new Server(
    {
      name: "{{ .Metadata.Name }}-mcp-server",
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
    const tools: Tool[] = [
      {{- range .Functions }}
      {{- if not .IsConstructor }}
      {{- if not .IsFallback }}
      {{- if not .IsReceive }}
      {{- if eq .StateMutability "view" "pure" }}
      {
        name: ToolName.{{ .Name | upper }},
        description: "{{ .Description }}",
        inputSchema: zodToJsonSchema({{ .Name | title }}Schema) as ToolInput,
      },
      {{- end }}
      {{- end }}
      {{- end }}
      {{- end }}
      {{- end }}
    ];
    
    return { tools };
  });
  
  // Handle tool calls
  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    const { name, arguments: args } = request.params;
    
    try {
      switch (name) {
        {{- range .Functions }}
        {{- if not .IsConstructor }}
        {{- if not .IsFallback }}
        {{- if not .IsReceive }}
        {{- if eq .StateMutability "view" "pure" }}
        case ToolName.{{ .Name | upper }}:
          const {{ .Name }}Args = {{ .Name | title }}Schema.parse(args);
          const {{ .Name }}Result = await contract.{{ .Name }}(
            {{- range $index, $param := .Inputs }}
            {{- if $index }}, {{ end }}{{ .Name }}Args.{{ .Name }}
            {{- end }}
          );
          return {
            content: [
              {
                type: "text",
                text: JSON.stringify({{ .Name }}Result, (key, value) => {
                  // Handle BigInt conversion
                  if (typeof value === 'bigint') {
                    return value.toString();
                  }
                  return value;
                }, 2),
              } as TextContent,
            ],
          };
        {{- end }}
        {{- end }}
        {{- end }}
        {{- end }}
        {{- end }}
        
        default:
          throw new Error(\`Unknown tool: \${name}\`);
      }
    } catch (error) {
      console.error("Error calling tool:", error);
      throw new Error(\`Error calling \${name}: \${error.message}\`);
    }
  });
  
  // Connect to transport
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});

{{- define "zodType" }}
{{- if eq .BaseType "uint256" "uint8" "uint16" "uint32" "uint64" "uint128" "int256" "int8" "int16" "int32" "int64" "int128" }}
number()
{{- else if eq .BaseType "bool" }}
boolean()
{{- else if eq .BaseType "string" "address" "bytes" "bytes32" }}
string()
{{- else if .IsArray }}
array(z.{{ template "zodType" (dict "BaseType" .BaseType) }})
{{- else }}
string()
{{- end }}
{{- end }}
`

// readmeTemplate is the template for README.md
const readmeTemplate = `# {{ .Metadata.Name }} MCP Server

This is an MCP (Model Context Protocol) server for the {{ .Metadata.Name }} smart contract.

## Overview

This server provides LLM access to the {{ .Metadata.Name }} smart contract through the Model Context Protocol. It exposes the following contract functions as tools:

{{- range .Functions }}
{{- if not .IsConstructor }}
{{- if not .IsFallback }}
{{- if not .IsReceive }}
{{- if eq .StateMutability "view" "pure" }}
- **{{ .Name }}**: {{ .Description }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

## Installation

1. Clone this repository
2. Install dependencies:
   \`\`\`bash
   npm install
   \`\`\`
3. Build the server:
   \`\`\`bash
   npm run build
   \`\`\`

## Configuration

Set the following environment variables:

- \`RPC_URL\`: Ethereum RPC URL (default: https://eth.llamarpc.com)
- \`CONTRACT_ADDRESS\`: Contract address (default: {{ .Metadata.Address }})

## Usage

Start the server:

\`\`\`bash
npm start
\`\`\`

The server uses stdio for communication with MCP clients.

## Contract Information

- **Name**: {{ .Metadata.Name }}
- **Chain**: {{ .Metadata.Chain }}
- **Address**: {{ .Metadata.Address }}

## License

MIT
`