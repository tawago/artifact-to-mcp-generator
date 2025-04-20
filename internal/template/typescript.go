package template

import (
        "bytes"
        "fmt"
        "strings"
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

// getFuncMap returns a template FuncMap with custom functions
func getFuncMap() template.FuncMap {
        funcMap := sprig.FuncMap()
        
        // Add custom functions
        funcMap["sub"] = func(a, b int) int {
                return a - b
        }
        
        funcMap["eq"] = func(a, b interface{}) bool {
                return a == b
        }
        
        funcMap["upper"] = func(s string) string {
                return strings.ToUpper(s)
        }
        
        funcMap["title"] = func(s string) string {
                if len(s) == 0 {
                        return s
                }
                return strings.ToTitle(string(s[0])) + s[1:]
        }
        
        return funcMap
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
        tmpl, err := template.New("package.json").Funcs(getFuncMap()).Parse(packageJSONTemplate)
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
        tmpl, err := template.New("tsconfig.json").Funcs(getFuncMap()).Parse(tsconfigJSONTemplate)
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
        tmpl, err := template.New("server.ts").Funcs(getFuncMap()).Parse(serverTSTemplate)
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
        // Implementation moved to readme.go
        const readmeTemplate = `# {{.Metadata.Name}} MCP Server

This is a Model Context Protocol (MCP) server for interacting with the {{.Metadata.Name}} smart contract.

## Available Functions

{{- range $funcIndex, $func := .Functions -}}
{{- if not $func.IsConstructor -}}
{{- if not $func.IsFallback -}}
{{- if not $func.IsReceive -}}
{{- if or (eq (printf "%s" $func.StateMutability) "view") (eq (printf "%s" $func.StateMutability) "pure") -}}
- **{{$func.Name}}**{{if $func.Inputs}} - Parameters: {{range $paramIndex, $param := $func.Inputs}}{{if $paramIndex}}, {{end}}{{$param.Name}} ({{$param.Type}}){{end}}{{end}}{{if $func.Outputs}} - Returns: {{range $outputIndex, $output := $func.Outputs}}{{if $outputIndex}}, {{end}}{{if $output.Name}}{{$output.Name}}{{else}}output{{$outputIndex}}{{end}} ({{$output.Type}}){{end}}{{end}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end }}

## Installation

` + "```bash" + `
npm install
` + "```" + `

## Building

` + "```bash" + `
npm run build
` + "```" + `

## Running

` + "```bash" + `
npm start
` + "```" + `
`

        // Create a template with sprig functions
        tmpl := template.New("README.md").Funcs(getFuncMap())

        // Parse the template
        tmpl, err := tmpl.Parse(readmeTemplate)
        if err != nil {
                return nil, err
        }

        // Execute the template
        var buf bytes.Buffer
        if err := tmpl.Execute(&buf, contract); err != nil {
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
{{- range $funcIndex, $func := .Functions -}}
{{- if not $func.IsConstructor -}}
{{- if not $func.IsFallback -}}
{{- if not $func.IsReceive -}}
{{- if or (eq (printf "%s" $func.StateMutability) "view") (eq (printf "%s" $func.StateMutability) "pure") }}
  {{$func.Name | upper}} = "{{$func.Name}}",
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end }}
}

// Define input schemas for each function
{{- range $funcIndex, $func := .Functions -}}
{{- if not $func.IsConstructor -}}
{{- if not $func.IsFallback -}}
{{- if not $func.IsReceive }}
const {{$func.Name | title}}Schema = z.object({
{{- range $func.Inputs}}
  {{.Name}}: z.string().describe("{{.Description}}"),
{{- end}}
});
{{end -}}
{{- end -}}
{{- end -}}
{{- end }}

// Initialize the contract
async function initializeContract() {
  // Connect to provider (replace with your preferred provider)
  const provider = new ethers.JsonRpcProvider(process.env.RPC_URL || "https://eth.llamarpc.com");
  
  // Contract address (can be overridden with environment variable)
  const contractAddress = process.env.CONTRACT_ADDRESS || "{{.Metadata.Address}}";
  
  // Contract ABI
  const contractABI = [
    {{- range $funcIndex, $func := .Functions }}
    {
      "name": "{{$func.Name}}",
      "type": "function",
      "inputs": [
        {{- range $index, $param := $func.Inputs -}}
        {{if $index}},{{end}}
        {
          "name": "{{$param.Name}}",
          "type": "{{$param.Type.BaseType}}"
          {{- if $param.Type.IsArray -}}
          ,
          "components": []
          {{- end -}}
        }
        {{- end -}}
      ],
      "outputs": [
        {{- range $index, $param := $func.Outputs -}}
        {{if $index}},{{end}}
        {
          "name": "{{$param.Name}}",
          "type": "{{$param.Type.BaseType}}"
          {{- if $param.Type.IsArray -}}
          ,
          "components": []
          {{- end -}}
        }
        {{- end -}}
      ],
      "stateMutability": "{{$func.StateMutability}}"
    }{{if not (eq $funcIndex (sub (len $.Functions) 1))}},{{end}}
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
      name: "{{.Metadata.Name}}-mcp-server",
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
    const tools = [
      {{- range $funcIndex, $func := .Functions -}}
      {{- if not $func.IsConstructor -}}
      {{- if not $func.IsFallback -}}
      {{- if not $func.IsReceive -}}
      {{- if or (eq (printf "%s" $func.StateMutability) "view") (eq (printf "%s" $func.StateMutability) "pure") }}
      {
        name: ToolName.{{$func.Name | upper}},
        description: "{{$func.Description}}",
        inputSchema: zodToJsonSchema({{$func.Name | title}}Schema),
      },
      {{- end -}}
      {{- end -}}
      {{- end -}}
      {{- end -}}
      {{- end }}
    ];
    
    return { tools };
  });
  
  // Handle tool calls
  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    const { name, arguments: args } = request.params;
    
    try {
      switch (name) {
        {{- range $funcIndex, $func := .Functions -}}
        {{- if not $func.IsConstructor -}}
        {{- if not $func.IsFallback -}}
        {{- if not $func.IsReceive -}}
        {{- if or (eq (printf "%s" $func.StateMutability) "view") (eq (printf "%s" $func.StateMutability) "pure") }}
        case ToolName.{{$func.Name | upper}}:
          const {{$func.Name}}Args = {{$func.Name | title}}Schema.parse(args);
          const {{$func.Name}}Result = await contract.{{$func.Name}}(
            {{- range $index, $param := $func.Inputs -}}
            {{if $index}}, {{end}}{{$func.Name}}Args.{{$param.Name}}
            {{- end -}}
          );
          return {
            content: [
              {
                type: "text",
                text: JSON.stringify({{$func.Name}}Result, (key, value) => {
                  // Handle BigInt conversion
                  if (typeof value === 'bigint') {
                    return value.toString();
                  }
                  return value;
                }, 2),
              },
            ],
          };
        {{- end -}}
        {{- end -}}
        {{- end -}}
        {{- end -}}
        {{- end }}
        
        default:
          throw new Error("Unknown tool: " + name);
      }
    } catch (error) {
      console.error("Error calling tool:", error);
      throw new Error("Error calling " + name + ": " + error.message);
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
`

// readmeTemplate is the template for README.md
const readmeTemplate = `# {{.Metadata.Name}} MCP Server

This is an MCP (Model Context Protocol) server for the {{.Metadata.Name}} smart contract.

## Overview

This server provides LLM access to the {{.Metadata.Name}} smart contract through the Model Context Protocol. It exposes the following contract functions as tools:

{{range $funcIndex, $func := .Functions}}
{{if not $func.IsConstructor}}
{{if not $func.IsFallback}}
{{if not $func.IsReceive}}
{{if or (eq $func.StateMutability "view") (eq $func.StateMutability "pure")}}
- **{{$func.Name}}**: {{$func.Description}}
{{end}}
{{end}}
{{end}}
{{end}}
{{end}}

## Installation

1. Clone this repository
2. Install dependencies:
   ` + "```bash" + `
   npm install
   ` + "```" + `
3. Build the server:
   ` + "```bash" + `
   npm run build
   ` + "```" + `

## Configuration

Set the following environment variables:

- RPC_URL: Ethereum RPC URL (default: https://eth.llamarpc.com)
- CONTRACT_ADDRESS: Contract address (default: {{.Metadata.Address}})

## Usage

Start the server:

` + "```bash" + `
npm start
` + "```" + `

The server uses stdio for communication with MCP clients.

## Contract Information

- **Name**: {{.Metadata.Name}}
- **Chain**: {{.Metadata.Chain}}
- **Address**: {{.Metadata.Address}}

## Available Functions
{{range $funcIndex, $func := .Functions}}
{{if not $func.IsConstructor}}
{{if not $func.IsFallback}}
{{if not $func.IsReceive}}
{{if or (eq (printf "%s" $func.StateMutability) "view") (eq (printf "%s" $func.StateMutability) "pure")}}
- **{{$func.Name}}**{{if $func.Inputs}} - Parameters: {{range $paramIndex, $param := $func.Inputs}}{{if $paramIndex}}, {{end}}{{$param.Name}} ({{$param.Type}}){{end}}{{end}}{{if $func.Outputs}} - Returns: {{range $outputIndex, $output := $func.Outputs}}{{if $outputIndex}}, {{end}}{{if $output.Name}}{{$output.Name}}{{else}}output{{$outputIndex}}{{end}} ({{$output.Type}}){{end}}{{end}}
{{end}}
{{end}}
{{end}}
{{end}}
{{end}}

## License

MIT
`