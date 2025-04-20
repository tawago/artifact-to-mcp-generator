package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/openhands/mcp-generator/internal/ir"
)

// EVMABIParser parses Ethereum ABI JSON into the intermediate representation
type EVMABIParser struct{}

// NewEVMABIParser creates a new EVM ABI parser
func NewEVMABIParser() *EVMABIParser {
	return &EVMABIParser{}
}

// Parse parses an EVM ABI from a reader into the intermediate representation
func (p *EVMABIParser) Parse(reader io.Reader, metadata ir.ContractMetadata) (*ir.ContractIR, error) {
	var abiItems []EVMABIItem
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&abiItems); err != nil {
		return nil, fmt.Errorf("failed to decode ABI JSON: %w", err)
	}

	contract := &ir.ContractIR{
		Metadata: metadata,
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
func (p *EVMABIParser) parseFunction(item EVMABIItem) (ir.Function, error) {
	inputs, err := p.parseParameters(item.Inputs)
	if err != nil {
		return ir.Function{}, fmt.Errorf("failed to parse function inputs: %w", err)
	}

	outputs, err := p.parseParameters(item.Outputs)
	if err != nil {
		return ir.Function{}, fmt.Errorf("failed to parse function outputs: %w", err)
	}

	// Build function signature
	signature := item.Name + "("
	inputTypes := make([]string, len(item.Inputs))
	for i, input := range item.Inputs {
		inputTypes[i] = input.Type
	}
	signature += strings.Join(inputTypes, ",") + ")"

	return ir.Function{
		Name:            item.Name,
		Description:     fmt.Sprintf("%s function", item.Name),
		Signature:       signature,
		Inputs:          inputs,
		Outputs:         outputs,
		StateMutability: ir.StateMutability(item.StateMutability),
		Visibility:      ir.Public, // Default to public for EVM functions in ABI
	}, nil
}

// parseEvent converts an ABI event item to IR Event
func (p *EVMABIParser) parseEvent(item EVMABIItem) (ir.Event, error) {
	parameters := make([]ir.EventParameter, len(item.Inputs))
	for i, input := range item.Inputs {
		paramType, err := p.parseParameterType(input.Type, input.Components)
		if err != nil {
			return ir.Event{}, fmt.Errorf("failed to parse event parameter type: %w", err)
		}

		parameters[i] = ir.EventParameter{
			Name:    input.Name,
			Type:    paramType,
			Indexed: input.Indexed,
		}
	}

	// Build event signature
	signature := item.Name + "("
	inputTypes := make([]string, len(item.Inputs))
	for i, input := range item.Inputs {
		inputTypes[i] = input.Type
	}
	signature += strings.Join(inputTypes, ",") + ")"

	return ir.Event{
		Name:        item.Name,
		Description: fmt.Sprintf("%s event", item.Name),
		Signature:   signature,
		Parameters:  parameters,
	}, nil
}

// parseError converts an ABI error item to IR ContractError
func (p *EVMABIParser) parseError(item EVMABIItem) (ir.ContractError, error) {
	parameters, err := p.parseParameters(item.Inputs)
	if err != nil {
		return ir.ContractError{}, fmt.Errorf("failed to parse error parameters: %w", err)
	}

	return ir.ContractError{
		Name:        item.Name,
		Description: fmt.Sprintf("%s error", item.Name),
		Parameters:  parameters,
	}, nil
}

// parseConstructor converts an ABI constructor item to IR Function
func (p *EVMABIParser) parseConstructor(item EVMABIItem) (ir.Function, error) {
	inputs, err := p.parseParameters(item.Inputs)
	if err != nil {
		return ir.Function{}, fmt.Errorf("failed to parse constructor inputs: %w", err)
	}

	return ir.Function{
		Name:            "constructor",
		Description:     "Contract constructor",
		Inputs:          inputs,
		Outputs:         []ir.Parameter{},
		StateMutability: ir.StateMutability(item.StateMutability),
		IsConstructor:   true,
	}, nil
}

// parseFallback converts an ABI fallback item to IR Function
func (p *EVMABIParser) parseFallback(item EVMABIItem) (ir.Function, error) {
	return ir.Function{
		Name:            "fallback",
		Description:     "Fallback function",
		Inputs:          []ir.Parameter{},
		Outputs:         []ir.Parameter{},
		StateMutability: ir.StateMutability(item.StateMutability),
		IsFallback:      true,
	}, nil
}

// parseReceive converts an ABI receive item to IR Function
func (p *EVMABIParser) parseReceive(item EVMABIItem) (ir.Function, error) {
	return ir.Function{
		Name:            "receive",
		Description:     "Receive function",
		Inputs:          []ir.Parameter{},
		Outputs:         []ir.Parameter{},
		StateMutability: ir.StateMutability(item.StateMutability),
		IsReceive:       true,
	}, nil
}

// parseParameters converts ABI parameters to IR Parameters
func (p *EVMABIParser) parseParameters(inputs []EVMABIInput) ([]ir.Parameter, error) {
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
func (p *EVMABIParser) parseParameterType(typeStr string, components []EVMABIInput) (ir.ParameterType, error) {
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
				// TODO: Parse the size string to int
				paramType.ArraySize = 0 // For now, default to 0 (dynamic)
			}
			
			// Extract the base type
			paramType.BaseType = typeStr[:start]
		} else {
			// Dynamic array
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
	}

	return paramType, nil
}

// EVMABIItem represents an item in the Ethereum ABI
type EVMABIItem struct {
	Type            string        `json:"type"`
	Name            string        `json:"name"`
	Inputs          []EVMABIInput `json:"inputs"`
	Outputs         []EVMABIInput `json:"outputs"`
	StateMutability string        `json:"stateMutability"`
	Anonymous       bool          `json:"anonymous"`
	Constant        bool          `json:"constant"`
	Payable         bool          `json:"payable"`
}

// EVMABIInput represents an input or output parameter in the Ethereum ABI
type EVMABIInput struct {
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Components []EVMABIInput `json:"components"`
	Indexed    bool          `json:"indexed"`
}