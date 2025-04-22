package template

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"
        "strings"
        "text/template"

        "github.com/Masterminds/sprig/v3"
        "github.com/openhands/mcp-generator/internal/ir"
)

// TypeScriptTemplateRenderer renders TypeScript MCP server templates
type TypeScriptTemplateRenderer struct {
        // Template directory path
        templateDir string
}

// NewTypeScriptTemplateRenderer creates a new TypeScript template renderer
func NewTypeScriptTemplateRenderer() *TypeScriptTemplateRenderer {
        // Default to the embedded templates if not specified
        return &TypeScriptTemplateRenderer{
                templateDir: filepath.Join("internal", "template", "typescript"),
        }
}

// WithTemplateDir sets a custom template directory
func (r *TypeScriptTemplateRenderer) WithTemplateDir(dir string) *TypeScriptTemplateRenderer {
        r.templateDir = dir
        return r
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

// loadTemplate loads a template file from the template directory
func (r *TypeScriptTemplateRenderer) loadTemplate(name string) (string, error) {
        templatePath := filepath.Join(r.templateDir, name)
        
        // Check if the file exists
        if _, err := os.Stat(templatePath); os.IsNotExist(err) {
                // Fall back to embedded templates if file doesn't exist
                switch name {
                case "package.json.tmpl":
                        return packageJSONTemplate, nil
                case "tsconfig.json.tmpl":
                        return tsconfigJSONTemplate, nil
                case "server.ts.tmpl":
                        return serverTSTemplate, nil
                case "README.md.tmpl":
                        return readmeTemplate, nil
                case "inspector-e2e/e2e-tests.spec.ts.tmpl":
                        return e2eTestsTemplate, nil
                case "playwright.config.ts.tmpl":
                        return playwrightConfigTemplate, nil
                default:
                        return "", fmt.Errorf("template %s not found", name)
                }
        }
        
        // Read the template file
        content, err := ioutil.ReadFile(templatePath)
        if err != nil {
                return "", fmt.Errorf("failed to read template %s: %w", name, err)
        }
        
        return string(content), nil
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

        // Generate e2e tests
        e2eTests, err := r.renderE2ETests(contract)
        if err != nil {
                return nil, fmt.Errorf("failed to render e2e tests: %w", err)
        }
        files["inspector-e2e/e2e-tests.spec.ts"] = e2eTests

        // Generate playwright config
        playwrightConfig, err := r.renderPlaywrightConfig(contract)
        if err != nil {
                return nil, fmt.Errorf("failed to render playwright config: %w", err)
        }
        files["playwright.config.ts"] = playwrightConfig

        return files, nil
}

// renderPackageJSON generates the package.json file
func (r *TypeScriptTemplateRenderer) renderPackageJSON(contract *ir.ContractIR) ([]byte, error) {
        // Load the template
        templateContent, err := r.loadTemplate("package.json.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Parse the template
        tmpl, err := template.New("package.json").Funcs(getFuncMap()).Parse(templateContent)
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
        // Load the template
        templateContent, err := r.loadTemplate("tsconfig.json.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Parse the template
        tmpl, err := template.New("tsconfig.json").Funcs(getFuncMap()).Parse(templateContent)
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
        // Load the template
        templateContent, err := r.loadTemplate("server.ts.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Parse the template
        tmpl, err := template.New("server.ts").Funcs(getFuncMap()).Parse(templateContent)
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
        // Load the template
        templateContent, err := r.loadTemplate("README.md.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Create a template with sprig functions
        tmpl := template.New("README.md").Funcs(getFuncMap())

        // Parse the template
        tmpl, err = tmpl.Parse(templateContent)
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

// renderE2ETests generates the e2e tests file
func (r *TypeScriptTemplateRenderer) renderE2ETests(contract *ir.ContractIR) ([]byte, error) {
        // Load the template
        templateContent, err := r.loadTemplate("inspector-e2e/e2e-tests.spec.ts.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Create a template with sprig functions
        tmpl := template.New("e2e-tests.spec.ts").Funcs(getFuncMap())

        // Parse the template
        tmpl, err = tmpl.Parse(templateContent)
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

// renderPlaywrightConfig generates the playwright.config.ts file
func (r *TypeScriptTemplateRenderer) renderPlaywrightConfig(contract *ir.ContractIR) ([]byte, error) {
        // Load the template
        templateContent, err := r.loadTemplate("playwright.config.ts.tmpl")
        if err != nil {
                return nil, err
        }
        
        // Create a template with sprig functions
        tmpl := template.New("playwright.config.ts").Funcs(getFuncMap())

        // Parse the template
        tmpl, err = tmpl.Parse(templateContent)
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

// Fallback templates in case the files don't exist
// These are kept for backward compatibility

// e2eTestsTemplate is the template for e2e-tests.spec.ts
const e2eTestsTemplate = `import { test, expect } from '@playwright/test';

test.describe('MCP Server Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Click the Connect button
    await page.getByRole('button', { name: 'Connect' }).click();
    await expect(page.getByText('Connected')).toBeVisible();
    
    // Click the List Tools button
    await page.getByRole('button', { name: 'List Tools' }).click();
  });

  test('should list available tools', async ({ page }) => {
    // Check if tools are listed
    await expect(page.getByText('Tool List')).toBeVisible();
  });
});`

// playwrightConfigTemplate is the template for playwright.config.ts
const playwrightConfigTemplate = `import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './inspector-e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 3 : 2,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL: "http://localhost:6274",
    trace: 'on-first-retry',
    navigationTimeout: 60000,
    actionTimeout: 30000,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: "npx @modelcontextprotocol/inspector node ./dist/server.js",
    url: "http://localhost:6274",
    reuseExistingServer: !process.env.CI,
    timeout: 180000,
  },
});`

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