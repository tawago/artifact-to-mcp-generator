package template

import (
        "os"
        "path/filepath"
        "strings"
        "testing"

        "github.com/openhands/mcp-generator/internal/ir"
)

func TestTypeScriptTemplateRenderer(t *testing.T) {
        // Create a sample contract IR for testing
        contract := &ir.ContractIR{
                Metadata: ir.ContractMetadata{
                        Name:        "TestToken",
                        Description: "A test ERC20 token",
                        Address:     "0x1234567890123456789012345678901234567890",
                        Chain:       "ethereum",
                },
                Functions: []ir.Function{
                        {
                                Name:            "balanceOf",
                                Description:     "Get the balance of an account",
                                StateMutability: ir.View,
                                Inputs: []ir.Parameter{
                                        {
                                                Name: "account",
                                                Type: ir.ParameterType{
                                                        BaseType: "address",
                                                },
                                                Description: "The address to query the balance of",
                                        },
                                },
                                Outputs: []ir.Parameter{
                                        {
                                                Name: "balance",
                                                Type: ir.ParameterType{
                                                        BaseType: "uint256",
                                                },
                                        },
                                },
                        },
                        {
                                Name:            "transfer",
                                Description:     "Transfer tokens to another address",
                                StateMutability: ir.Nonpayable,
                                Inputs: []ir.Parameter{
                                        {
                                                Name: "to",
                                                Type: ir.ParameterType{
                                                        BaseType: "address",
                                                },
                                                Description: "The address to transfer to",
                                        },
                                        {
                                                Name: "amount",
                                                Type: ir.ParameterType{
                                                        BaseType: "uint256",
                                                },
                                                Description: "The amount to transfer",
                                        },
                                },
                                Outputs: []ir.Parameter{
                                        {
                                                Name: "success",
                                                Type: ir.ParameterType{
                                                        BaseType: "bool",
                                                },
                                        },
                                },
                        },
                },
        }

        // Create a renderer
        renderer := NewTypeScriptTemplateRenderer()

        // Render the templates
        files, err := renderer.Render(contract)
        if err != nil {
                t.Fatalf("Failed to render templates: %v", err)
        }

        // Check that we have the expected files
        expectedFiles := []string{
                "package.json",
                "tsconfig.json",
                "src/server.ts",
                "README.md",
        }

        for _, file := range expectedFiles {
                if _, ok := files[file]; !ok {
                        t.Errorf("Expected file %s not found in rendered output", file)
                }
        }

        // Check that the package.json contains the contract name
        if !contains(string(files["package.json"]), "testtoken-mcp-server") {
                t.Errorf("package.json does not contain the expected contract name")
        }

        // Check that the server.ts contains the balanceOf function but not the transfer function
        // (since transfer is nonpayable and not exposed as a tool)
        serverTS := string(files["src/server.ts"])
        if !contains(serverTS, "balanceOf") {
                t.Errorf("server.ts does not contain the balanceOf function")
        }
        if contains(serverTS, "ToolName.TRANSFER") {
                t.Errorf("server.ts contains the transfer function as a tool, but it should not (nonpayable)")
        }
}

// TestTypeScriptTemplateRendererWithCustomDir tests the renderer with a custom template directory
func TestTypeScriptTemplateRendererWithCustomDir(t *testing.T) {
        // Create a temporary directory for test templates
        tempDir, err := os.MkdirTemp("", "typescript-templates-*")
        if err != nil {
                t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Create a minimal test template
        testTemplate := `Test template for {{.Metadata.Name}}`
        err = os.WriteFile(filepath.Join(tempDir, "package.json.tmpl"), []byte(testTemplate), 0644)
        if err != nil {
                t.Fatalf("Failed to write test template: %v", err)
        }

        // Create a sample contract IR
        contract := &ir.ContractIR{
                Metadata: ir.ContractMetadata{
                        Name: "TestToken",
                },
        }

        // Create a renderer with the custom directory
        renderer := NewTypeScriptTemplateRenderer().WithTemplateDir(tempDir)

        // Render the templates
        files, err := renderer.Render(contract)
        if err != nil {
                t.Fatalf("Failed to render templates: %v", err)
        }

        // Check that the package.json contains our test template content
        if string(files["package.json"]) != "Test template for TestToken" {
                t.Errorf("package.json does not contain the expected content from custom template")
        }
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
        return strings.Contains(s, substr)
}