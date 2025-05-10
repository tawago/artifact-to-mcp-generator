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
        // Get the absolute path to the project root
        projectRoot := filepath.Join("/workspace", "artifact-to-mcp-generator")

        // Default to the embedded templates if not specified
        return &TypeScriptTemplateRenderer{
                templateDir: filepath.Join(projectRoot, "internal", "template", "typescript"),
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
                return "", fmt.Errorf("template %s not found", name)
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