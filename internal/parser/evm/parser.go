package evm

import (
        "encoding/json"
        "fmt"
        "io"
        "strconv"
        "strings"

        "github.com/openhands/mcp-generator/internal/ir"
)

// ABIParser parses Ethereum ABI JSON into the intermediate representation
type ABIParser struct {
        // Map to track function signatures for handling overloads
        functionSignatures map[string]int
}

// NewABIParser creates a new EVM ABI parser
func NewABIParser() *ABIParser {
        return &ABIParser{
                functionSignatures: make(map[string]int),
        }
}

// Parse parses an EVM ABI from a reader into the intermediate representation
func (p *ABIParser) Parse(reader io.Reader, metadata ir.ContractMetadata) (*ir.ContractIR, error) {
        var abiItems []ABIItem
        decoder := json.NewDecoder(reader)
        if err := decoder.Decode(&abiItems); err != nil {
                return nil, fmt.Errorf("failed to decode ABI JSON: %w", err)
        }

        contract := &ir.ContractIR{
                Metadata:  metadata,
                Functions: []ir.Function{},
                Events:    []ir.Event{},
                Errors:    []ir.ContractError{},
        }

        // Set chain to ethereum if not specified
        if contract.Metadata.Chain == "" {
                contract.Metadata.Chain = "ethereum"
        }

        for _, item := range abiItems {
                switch item.Type {
                case "function":
                        function, err := p.parseFunction(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Functions = append(contract.Functions, function)
                case "event":
                        event, err := p.parseEvent(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Events = append(contract.Events, event)
                case "error":
                        contractError, err := p.parseError(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Errors = append(contract.Errors, contractError)
                case "constructor":
                        function, err := p.parseConstructor(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Functions = append(contract.Functions, function)
                case "fallback":
                        function, err := p.parseFallback(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Functions = append(contract.Functions, function)
                case "receive":
                        function, err := p.parseReceive(item)
                        if err != nil {
                                return nil, err
                        }
                        contract.Functions = append(contract.Functions, function)
                }
        }

        return contract, nil
}

// parseFunction converts an ABI function item to IR Function
func (p *ABIParser) parseFunction(item ABIItem) (ir.Function, error) {
        inputs, err := p.parseParameters(item.Inputs)
        if err != nil {
                return ir.Function{}, fmt.Errorf("failed to parse function inputs: %w", err)
        }

        outputs, err := p.parseParameters(item.Outputs)
        if err != nil {
                return ir.Function{}, fmt.Errorf("failed to parse function outputs: %w", err)
        }

        // Build function signature
        signature := buildFunctionSignature(item.Name, item.Inputs)

        // Handle function overloads
        functionName := item.Name
        if count, exists := p.functionSignatures[item.Name]; exists {
                // This is an overloaded function, append a suffix to make it unique
                p.functionSignatures[item.Name] = count + 1
                functionName = fmt.Sprintf("%s_%d", item.Name, count)
        } else {
                p.functionSignatures[item.Name] = 1
        }

        // Generate a better description based on the function name and inputs
        description := item.Name
        if len(inputs) > 0 {
                description += " - Parameters: "
                for i, input := range inputs {
                        if i > 0 {
                                description += ", "
                        }
                        typeStr := string(input.Type.BaseType)
                        if input.Type.IsArray {
                                if input.Type.ArraySize > 0 {
                                        typeStr += fmt.Sprintf("[%d]", input.Type.ArraySize)
                                } else {
                                        typeStr += "[]"
                                }
                        } else if len(input.Type.Components) > 0 {
                                typeStr = "tuple"
                        }
                        description += input.Name + " (" + typeStr + ")"
                }
        }
        
        if len(outputs) > 0 {
                description += " - Returns: "
                for i, output := range outputs {
                        if i > 0 {
                                description += ", "
                        }
                        outputName := output.Name
                        if outputName == "" {
                                outputName = fmt.Sprintf("output%d", i)
                        }
                        typeStr := string(output.Type.BaseType)
                        if output.Type.IsArray {
                                if output.Type.ArraySize > 0 {
                                        typeStr += fmt.Sprintf("[%d]", output.Type.ArraySize)
                                } else {
                                        typeStr += "[]"
                                }
                        } else if len(output.Type.Components) > 0 {
                                typeStr = "tuple"
                        }
                        description += outputName + " (" + typeStr + ")"
                }
        }

        // Calculate function selector (first 4 bytes of keccak256 hash of the signature)
        // In a real implementation, we would compute this, but for now we'll leave it empty
        selector := ""

        // Determine state mutability
        stateMutability := ir.StateMutability(item.StateMutability)
        if stateMutability == "" {
                // Handle legacy ABI format
                if item.Constant {
                        stateMutability = ir.View
                } else if item.Payable {
                        stateMutability = ir.Payable
                } else {
                        stateMutability = ir.Nonpayable
                }
        }

        // Create chain-specific data
        chainData := make(map[string]interface{})
        if item.Constant {
                chainData["constant"] = item.Constant
        }
        if item.Payable {
                chainData["payable"] = item.Payable
        }
        
        // Store original name and signature for overloaded functions
        if functionName != item.Name {
                chainData["originalName"] = item.Name
                chainData["originalSignature"] = signature
        }

        return ir.Function{
                Name:            functionName,
                Description:     description,
                Signature:       signature,
                Selector:        selector,
                Inputs:          inputs,
                Outputs:         outputs,
                StateMutability: stateMutability,
                Visibility:      ir.Public, // Default to public for EVM functions in ABI
                ChainData:       chainData,
        }, nil
}

// parseEvent converts an ABI event item to IR Event
func (p *ABIParser) parseEvent(item ABIItem) (ir.Event, error) {
        parameters := make([]ir.EventParameter, len(item.Inputs))
        indexedCount := 0
        
        for i, input := range item.Inputs {
                paramType, err := p.parseParameterType(input.Type, input.Components)
                if err != nil {
                        return ir.Event{}, fmt.Errorf("failed to parse event parameter type: %w", err)
                }
                
                // Count indexed parameters (EVM allows up to 3)
                if input.Indexed {
                        indexedCount++
                }

                parameters[i] = ir.EventParameter{
                        Name:    input.Name,
                        Type:    paramType,
                        Indexed: input.Indexed,
                }
        }

        // Build event signature
        signature := buildEventSignature(item.Name, item.Inputs)

        // Create chain-specific data
        chainData := make(map[string]interface{})
        if item.Anonymous {
                chainData["anonymous"] = item.Anonymous
        }
        
        // Add indexed parameters information
        chainData["indexedCount"] = indexedCount
        
        // Generate a better description that includes indexed parameters
        description := fmt.Sprintf("%s event", item.Name)
        if len(parameters) > 0 {
                description += " - Parameters: "
                for i, param := range parameters {
                        if i > 0 {
                                description += ", "
                        }
                        typeStr := string(param.Type.BaseType)
                        if param.Type.IsArray {
                                if param.Type.ArraySize > 0 {
                                        typeStr += fmt.Sprintf("[%d]", param.Type.ArraySize)
                                } else {
                                        typeStr += "[]"
                                }
                        } else if len(param.Type.Components) > 0 {
                                typeStr = "tuple"
                        }
                        
                        indexedStr := ""
                        if param.Indexed {
                                indexedStr = " (indexed)"
                        }
                        
                        description += param.Name + " (" + typeStr + ")" + indexedStr
                }
        }

        return ir.Event{
                Name:        item.Name,
                Description: description,
                Signature:   signature,
                Parameters:  parameters,
                ChainData:   chainData,
        }, nil
}

// parseError converts an ABI error item to IR ContractError
func (p *ABIParser) parseError(item ABIItem) (ir.ContractError, error) {
        parameters, err := p.parseParameters(item.Inputs)
        if err != nil {
                return ir.ContractError{}, fmt.Errorf("failed to parse error parameters: %w", err)
        }

        // Build error signature (not used in IR currently but could be useful for future extensions)
        _ = buildErrorSignature(item.Name, item.Inputs)

        return ir.ContractError{
                Name:        item.Name,
                Description: fmt.Sprintf("%s error", item.Name),
                Parameters:  parameters,
        }, nil
}

// parseConstructor converts an ABI constructor item to IR Function
func (p *ABIParser) parseConstructor(item ABIItem) (ir.Function, error) {
        inputs, err := p.parseParameters(item.Inputs)
        if err != nil {
                return ir.Function{}, fmt.Errorf("failed to parse constructor inputs: %w", err)
        }

        // Determine state mutability
        stateMutability := ir.StateMutability(item.StateMutability)
        if stateMutability == "" {
                // Handle legacy ABI format
                if item.Payable {
                        stateMutability = ir.Payable
                } else {
                        stateMutability = ir.Nonpayable
                }
        }

        return ir.Function{
                Name:            "constructor",
                Description:     "Contract constructor",
                Inputs:          inputs,
                Outputs:         []ir.Parameter{},
                StateMutability: stateMutability,
                IsConstructor:   true,
        }, nil
}

// parseFallback converts an ABI fallback item to IR Function
func (p *ABIParser) parseFallback(item ABIItem) (ir.Function, error) {
        // Determine state mutability
        stateMutability := ir.StateMutability(item.StateMutability)
        if stateMutability == "" {
                // Handle legacy ABI format
                if item.Payable {
                        stateMutability = ir.Payable
                } else {
                        stateMutability = ir.Nonpayable
                }
        }

        return ir.Function{
                Name:            "fallback",
                Description:     "Fallback function",
                Inputs:          []ir.Parameter{},
                Outputs:         []ir.Parameter{},
                StateMutability: stateMutability,
                IsFallback:      true,
        }, nil
}

// parseReceive converts an ABI receive item to IR Function
func (p *ABIParser) parseReceive(item ABIItem) (ir.Function, error) {
        return ir.Function{
                Name:            "receive",
                Description:     "Receive function",
                Inputs:          []ir.Parameter{},
                Outputs:         []ir.Parameter{},
                StateMutability: ir.Payable, // Receive functions are always payable
                IsReceive:       true,
        }, nil
}

// parseParameters converts ABI parameters to IR Parameters
func (p *ABIParser) parseParameters(inputs []ABIInput) ([]ir.Parameter, error) {
        parameters := make([]ir.Parameter, len(inputs))
        for i, input := range inputs {
                paramType, err := p.parseParameterType(input.Type, input.Components)
                if err != nil {
                        return nil, err
                }

                parameters[i] = ir.Parameter{
                        Name: input.Name,
                        Type: paramType,
                }
        }
        return parameters, nil
}

// parseParameterType converts an ABI type string to IR ParameterType
func (p *ABIParser) parseParameterType(typeStr string, components []ABIInput) (ir.ParameterType, error) {
        paramType := ir.ParameterType{}

        // Check if it's an array type
        isArray := strings.HasSuffix(typeStr, "[]") || strings.Contains(typeStr, "[")
        if isArray {
                paramType.IsArray = true
                
                // Check if it's a fixed-size array
                if strings.Contains(typeStr, "[") && !strings.HasSuffix(typeStr, "[]") {
                        // Extract the array size
                        start := strings.Index(typeStr, "[")
                        end := strings.Index(typeStr, "]")
                        if start != -1 && end != -1 {
                                sizeStr := typeStr[start+1 : end]
                                size, err := strconv.Atoi(sizeStr)
                                if err != nil {
                                        return paramType, fmt.Errorf("invalid array size: %s", sizeStr)
                                }
                                paramType.ArraySize = size
                        }
                        
                        // Extract the base type
                        paramType.BaseType = typeStr[:start]
                } else {
                        // Dynamic array
                        paramType.ArraySize = 0 // 0 indicates dynamic size
                        paramType.BaseType = typeStr[:len(typeStr)-2]
                }
        } else {
                paramType.BaseType = typeStr
        }

        // Handle tuple types (structs)
        if paramType.BaseType == "tuple" && components != nil {
                componentParams, err := p.parseParameters(components)
                if err != nil {
                        return paramType, err
                }
                paramType.Components = componentParams
                
                // Add additional metadata for complex types
                chainData := make(map[string]interface{})
                chainData["isTuple"] = true
                
                // Create a type description for the tuple
                typeDesc := "{"
                for i, comp := range componentParams {
                        if i > 0 {
                                typeDesc += ", "
                        }
                        compTypeStr := string(comp.Type.BaseType)
                        if comp.Type.IsArray {
                                if comp.Type.ArraySize > 0 {
                                        compTypeStr += fmt.Sprintf("[%d]", comp.Type.ArraySize)
                                } else {
                                        compTypeStr += "[]"
                                }
                        }
                        typeDesc += comp.Name + ": " + compTypeStr
                }
                typeDesc += "}"
                chainData["typeDescription"] = typeDesc
                
                paramType.ChainData = chainData
        }

        // Handle mapping types (not directly supported in ABI, but we can detect some common patterns)
        if strings.HasPrefix(paramType.BaseType, "mapping(") {
                paramType.IsMap = true
                // Extract key type (simplified, would need more robust parsing in a real implementation)
                keyStart := len("mapping(")
                keyEnd := strings.Index(paramType.BaseType, "=>")
                if keyEnd != -1 {
                        paramType.MapKeyType = strings.TrimSpace(paramType.BaseType[keyStart:keyEnd])
                        // The value type becomes the base type
                        valueStart := keyEnd + 2
                        valueEnd := len(paramType.BaseType) - 1 // Exclude closing parenthesis
                        paramType.BaseType = strings.TrimSpace(paramType.BaseType[valueStart:valueEnd])
                }
        }
        
        // Add array-specific metadata
        if paramType.IsArray {
                if paramType.ChainData == nil {
                        paramType.ChainData = make(map[string]interface{})
                }
                paramType.ChainData["isArray"] = true
                if paramType.ArraySize > 0 {
                        paramType.ChainData["isFixedArray"] = true
                        paramType.ChainData["arraySize"] = paramType.ArraySize
                } else {
                        paramType.ChainData["isDynamicArray"] = true
                }
        }

        return paramType, nil
}

// buildFunctionSignature creates a canonical function signature for EVM
func buildFunctionSignature(name string, inputs []ABIInput) string {
        signature := name + "("
        inputTypes := make([]string, len(inputs))
        for i, input := range inputs {
                inputTypes[i] = input.Type
        }
        signature += strings.Join(inputTypes, ",") + ")"
        return signature
}

// buildEventSignature creates a canonical event signature for EVM
func buildEventSignature(name string, inputs []ABIInput) string {
        signature := name + "("
        inputTypes := make([]string, len(inputs))
        for i, input := range inputs {
                inputTypes[i] = input.Type
        }
        signature += strings.Join(inputTypes, ",") + ")"
        return signature
}

// buildErrorSignature creates a canonical error signature for EVM
func buildErrorSignature(name string, inputs []ABIInput) string {
        signature := name + "("
        inputTypes := make([]string, len(inputs))
        for i, input := range inputs {
                inputTypes[i] = input.Type
        }
        signature += strings.Join(inputTypes, ",") + ")"
        return signature
}

// ABIItem represents an item in the Ethereum ABI
type ABIItem struct {
        Type            string     `json:"type"`
        Name            string     `json:"name"`
        Inputs          []ABIInput `json:"inputs"`
        Outputs         []ABIInput `json:"outputs"`
        StateMutability string     `json:"stateMutability"`
        Anonymous       bool       `json:"anonymous"`
        Constant        bool       `json:"constant"`
        Payable         bool       `json:"payable"`
}

// ABIInput represents an input or output parameter in the Ethereum ABI
type ABIInput struct {
        Name       string     `json:"name"`
        Type       string     `json:"type"`
        Components []ABIInput `json:"components"`
        Indexed    bool       `json:"indexed"`
}