package main

import (
        "fmt"
        "os"
        "path/filepath"

        "github.com/openhands/mcp-generator/internal/ir"
        "github.com/openhands/mcp-generator/internal/parser"
        "github.com/openhands/mcp-generator/internal/template"
        "github.com/spf13/cobra"
)

var (
        artifactPath string
        outputDir    string
        lang         string
        chainType    string
        contractName string
        contractAddr string
)

func main() {
        rootCmd := &cobra.Command{
                Use:   "generate-mcp",
                Short: "Generate MCP servers for smart contracts",
                Long:  `A tool that generates typed MCP (Model Context Protocol) servers for any deployed smart contract.`,
                RunE:  run,
        }

        rootCmd.Flags().StringVarP(&artifactPath, "artifact", "a", "", "Path to the contract artifact (ABI/IDL)")
        rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./mcp-server", "Output directory for the generated MCP server")
        rootCmd.Flags().StringVarP(&lang, "lang", "l", "ts", "Output language (ts, python)")
        rootCmd.Flags().StringVarP(&chainType, "chain", "c", "ethereum", "Blockchain type (ethereum, solana)")
        rootCmd.Flags().StringVarP(&contractName, "name", "n", "", "Contract name")
        rootCmd.Flags().StringVarP(&contractAddr, "address", "d", "", "Contract address")

        rootCmd.MarkFlagRequired("artifact")

        if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

func run(cmd *cobra.Command, args []string) error {
        // Open the artifact file
        file, err := os.Open(artifactPath)
        if err != nil {
                return fmt.Errorf("failed to open artifact file: %w", err)
        }
        defer file.Close()

        // If contract name is not provided, use the filename
        if contractName == "" {
                contractName = filepath.Base(artifactPath)
                contractName = contractName[:len(contractName)-len(filepath.Ext(contractName))]
        }

        // Create metadata
        metadata := ir.ContractMetadata{
                Name:    contractName,
                Chain:   chainType,
                Address: contractAddr,
        }

        // Parse the artifact
        var contractIR *ir.ContractIR
        switch chainType {
        case "ethereum", "evm":
                p := parser.NewEVMABIParser()
                contractIR, err = p.Parse(file, metadata)
                if err != nil {
                        return fmt.Errorf("failed to parse EVM ABI: %w", err)
                }
        case "solana":
                return fmt.Errorf("solana support not implemented yet")
        default:
                return fmt.Errorf("unsupported chain type: %s", chainType)
        }

        // Debug output
        fmt.Println("Parsed contract IR:")
        fmt.Printf("Functions: %d\n", len(contractIR.Functions))
        for i, f := range contractIR.Functions {
                fmt.Printf("  %d. %s (StateMutability: %s)\n", i+1, f.Name, f.StateMutability)
        }
        
        // Generate the MCP server
        var files map[string][]byte
        switch lang {
        case "ts", "typescript":
                r := template.NewTypeScriptTemplateRenderer()
                files, err = r.Render(contractIR)
                if err != nil {
                        return fmt.Errorf("failed to render TypeScript MCP server: %w", err)
                }
        case "python", "py":
                return fmt.Errorf("python support not implemented yet")
        default:
                return fmt.Errorf("unsupported language: %s", lang)
        }

        // Create the output directory
        if err := os.MkdirAll(outputDir, 0755); err != nil {
                return fmt.Errorf("failed to create output directory: %w", err)
        }

        // Create src directory
        if err := os.MkdirAll(filepath.Join(outputDir, "src"), 0755); err != nil {
                return fmt.Errorf("failed to create src directory: %w", err)
        }

        // Write the files
        for path, content := range files {
                fullPath := filepath.Join(outputDir, path)
                
                // Create parent directories if they don't exist
                if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
                        return fmt.Errorf("failed to create directory for %s: %w", path, err)
                }
                
                if err := os.WriteFile(fullPath, content, 0644); err != nil {
                        return fmt.Errorf("failed to write file %s: %w", path, err)
                }
        }

        fmt.Printf("MCP server generated successfully in %s\n", outputDir)
        return nil
}