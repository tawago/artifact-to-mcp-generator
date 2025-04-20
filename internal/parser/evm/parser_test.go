package evm

import (
	"strings"
	"testing"

	"github.com/openhands/mcp-generator/internal/ir"
	"github.com/stretchr/testify/assert"
)

func TestABIParser_Parse(t *testing.T) {
	// Test basic ERC20 ABI
	abiJSON := `[
		{
			"constant": true,
			"inputs": [{"name": "owner", "type": "address"}],
			"name": "balanceOf",
			"outputs": [{"name": "", "type": "uint256"}],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "value", "type": "uint256"}
			],
			"name": "transfer",
			"outputs": [{"name": "", "type": "bool"}],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "name": "from", "type": "address"},
				{"indexed": true, "name": "to", "type": "address"},
				{"indexed": false, "name": "value", "type": "uint256"}
			],
			"name": "Transfer",
			"type": "event"
		}
	]`

	parser := NewABIParser()
	metadata := ir.ContractMetadata{
		Name:  "TestToken",
		Chain: "ethereum",
	}

	contractIR, err := parser.Parse(strings.NewReader(abiJSON), metadata)
	assert.NoError(t, err)
	assert.NotNil(t, contractIR)

	// Check contract metadata
	assert.Equal(t, "TestToken", contractIR.Metadata.Name)
	assert.Equal(t, "ethereum", contractIR.Metadata.Chain)

	// Check functions
	assert.Len(t, contractIR.Functions, 2)
	
	// Check balanceOf function
	balanceOf := contractIR.Functions[0]
	assert.Equal(t, "balanceOf", balanceOf.Name)
	assert.Equal(t, ir.View, balanceOf.StateMutability)
	assert.Len(t, balanceOf.Inputs, 1)
	assert.Equal(t, "owner", balanceOf.Inputs[0].Name)
	assert.Equal(t, "address", balanceOf.Inputs[0].Type.BaseType)
	assert.Len(t, balanceOf.Outputs, 1)
	assert.Equal(t, "uint256", balanceOf.Outputs[0].Type.BaseType)

	// Check transfer function
	transfer := contractIR.Functions[1]
	assert.Equal(t, "transfer", transfer.Name)
	assert.Equal(t, ir.Nonpayable, transfer.StateMutability)
	assert.Len(t, transfer.Inputs, 2)
	assert.Equal(t, "to", transfer.Inputs[0].Name)
	assert.Equal(t, "address", transfer.Inputs[0].Type.BaseType)
	assert.Equal(t, "value", transfer.Inputs[1].Name)
	assert.Equal(t, "uint256", transfer.Inputs[1].Type.BaseType)
	assert.Len(t, transfer.Outputs, 1)
	assert.Equal(t, "bool", transfer.Outputs[0].Type.BaseType)

	// Check events
	assert.Len(t, contractIR.Events, 1)
	
	// Check Transfer event
	transferEvent := contractIR.Events[0]
	assert.Equal(t, "Transfer", transferEvent.Name)
	assert.Len(t, transferEvent.Parameters, 3)
	assert.Equal(t, "from", transferEvent.Parameters[0].Name)
	assert.Equal(t, "address", transferEvent.Parameters[0].Type.BaseType)
	assert.True(t, transferEvent.Parameters[0].Indexed)
	assert.Equal(t, "to", transferEvent.Parameters[1].Name)
	assert.Equal(t, "address", transferEvent.Parameters[1].Type.BaseType)
	assert.True(t, transferEvent.Parameters[1].Indexed)
	assert.Equal(t, "value", transferEvent.Parameters[2].Name)
	assert.Equal(t, "uint256", transferEvent.Parameters[2].Type.BaseType)
	assert.False(t, transferEvent.Parameters[2].Indexed)
}

func TestABIParser_ComplexTypes(t *testing.T) {
	// Test ABI with complex types (arrays, structs)
	abiJSON := `[
		{
			"inputs": [
				{"name": "values", "type": "uint256[]"},
				{"name": "fixedValues", "type": "uint256[3]"}
			],
			"name": "processArrays",
			"outputs": [{"name": "", "type": "uint256"}],
			"stateMutability": "pure",
			"type": "function"
		},
		{
			"inputs": [
				{
					"name": "person",
					"type": "tuple",
					"components": [
						{"name": "name", "type": "string"},
						{"name": "age", "type": "uint256"},
						{"name": "addresses", "type": "address[]"}
					]
				}
			],
			"name": "processPerson",
			"outputs": [{"name": "", "type": "bool"}],
			"stateMutability": "view",
			"type": "function"
		}
	]`

	parser := NewABIParser()
	metadata := ir.ContractMetadata{
		Name:  "ComplexContract",
		Chain: "ethereum",
	}

	contractIR, err := parser.Parse(strings.NewReader(abiJSON), metadata)
	assert.NoError(t, err)
	assert.NotNil(t, contractIR)

	// Check functions
	assert.Len(t, contractIR.Functions, 2)
	
	// Check processArrays function
	processArrays := contractIR.Functions[0]
	assert.Equal(t, "processArrays", processArrays.Name)
	assert.Equal(t, ir.Pure, processArrays.StateMutability)
	assert.Len(t, processArrays.Inputs, 2)
	
	// Check dynamic array parameter
	assert.Equal(t, "values", processArrays.Inputs[0].Name)
	assert.Equal(t, "uint256", processArrays.Inputs[0].Type.BaseType)
	assert.True(t, processArrays.Inputs[0].Type.IsArray)
	assert.Equal(t, 0, processArrays.Inputs[0].Type.ArraySize) // 0 means dynamic
	
	// Check fixed array parameter
	assert.Equal(t, "fixedValues", processArrays.Inputs[1].Name)
	assert.Equal(t, "uint256", processArrays.Inputs[1].Type.BaseType)
	assert.True(t, processArrays.Inputs[1].Type.IsArray)
	assert.Equal(t, 3, processArrays.Inputs[1].Type.ArraySize)
	
	// Check processPerson function (struct/tuple)
	processPerson := contractIR.Functions[1]
	assert.Equal(t, "processPerson", processPerson.Name)
	assert.Equal(t, ir.View, processPerson.StateMutability)
	assert.Len(t, processPerson.Inputs, 1)
	
	// Check struct parameter
	assert.Equal(t, "person", processPerson.Inputs[0].Name)
	assert.Equal(t, "tuple", processPerson.Inputs[0].Type.BaseType)
	assert.Len(t, processPerson.Inputs[0].Type.Components, 3)
	
	// Check struct components
	components := processPerson.Inputs[0].Type.Components
	assert.Equal(t, "name", components[0].Name)
	assert.Equal(t, "string", components[0].Type.BaseType)
	assert.Equal(t, "age", components[1].Name)
	assert.Equal(t, "uint256", components[1].Type.BaseType)
	assert.Equal(t, "addresses", components[2].Name)
	assert.Equal(t, "address", components[2].Type.BaseType)
	assert.True(t, components[2].Type.IsArray)
	assert.Equal(t, 0, components[2].Type.ArraySize) // Dynamic array
}

func TestABIParser_FunctionOverloads(t *testing.T) {
	// Test ABI with function overloads
	abiJSON := `[
		{
			"inputs": [{"name": "value", "type": "uint256"}],
			"name": "setValue",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [{"name": "value", "type": "string"}],
			"name": "setValue",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	parser := NewABIParser()
	metadata := ir.ContractMetadata{
		Name:  "OverloadContract",
		Chain: "ethereum",
	}

	contractIR, err := parser.Parse(strings.NewReader(abiJSON), metadata)
	assert.NoError(t, err)
	assert.NotNil(t, contractIR)

	// Check functions
	assert.Len(t, contractIR.Functions, 2)
	
	// First overload should keep original name
	assert.Equal(t, "setValue", contractIR.Functions[0].Name)
	assert.Equal(t, "uint256", contractIR.Functions[0].Inputs[0].Type.BaseType)
	
	// Second overload should have a suffix
	assert.Equal(t, "setValue_1", contractIR.Functions[1].Name)
	assert.Equal(t, "string", contractIR.Functions[1].Inputs[0].Type.BaseType)
	
	// Signatures should be different
	assert.NotEqual(t, contractIR.Functions[0].Signature, contractIR.Functions[1].Signature)
}

func TestABIParser_Errors(t *testing.T) {
	// Test ABI with custom errors
	abiJSON := `[
		{
			"inputs": [
				{"name": "available", "type": "uint256"},
				{"name": "required", "type": "uint256"}
			],
			"name": "InsufficientBalance",
			"type": "error"
		},
		{
			"inputs": [{"name": "addr", "type": "address"}],
			"name": "Unauthorized",
			"type": "error"
		}
	]`

	parser := NewABIParser()
	metadata := ir.ContractMetadata{
		Name:  "ErrorContract",
		Chain: "ethereum",
	}

	contractIR, err := parser.Parse(strings.NewReader(abiJSON), metadata)
	assert.NoError(t, err)
	assert.NotNil(t, contractIR)

	// Check errors
	assert.Len(t, contractIR.Errors, 2)
	
	// Check InsufficientBalance error
	insufficientBalance := contractIR.Errors[0]
	assert.Equal(t, "InsufficientBalance", insufficientBalance.Name)
	assert.Len(t, insufficientBalance.Parameters, 2)
	assert.Equal(t, "available", insufficientBalance.Parameters[0].Name)
	assert.Equal(t, "uint256", insufficientBalance.Parameters[0].Type.BaseType)
	assert.Equal(t, "required", insufficientBalance.Parameters[1].Name)
	assert.Equal(t, "uint256", insufficientBalance.Parameters[1].Type.BaseType)
	
	// Check Unauthorized error
	unauthorized := contractIR.Errors[1]
	assert.Equal(t, "Unauthorized", unauthorized.Name)
	assert.Len(t, unauthorized.Parameters, 1)
	assert.Equal(t, "addr", unauthorized.Parameters[0].Name)
	assert.Equal(t, "address", unauthorized.Parameters[0].Type.BaseType)
}

func TestABIParser_SpecialFunctions(t *testing.T) {
	// Test ABI with constructor, fallback, and receive functions
	abiJSON := `[
		{
			"inputs": [
				{"name": "name", "type": "string"},
				{"name": "symbol", "type": "string"}
			],
			"stateMutability": "nonpayable",
			"type": "constructor"
		},
		{
			"stateMutability": "payable",
			"type": "fallback"
		},
		{
			"stateMutability": "payable",
			"type": "receive"
		}
	]`

	parser := NewABIParser()
	metadata := ir.ContractMetadata{
		Name:  "SpecialContract",
		Chain: "ethereum",
	}

	contractIR, err := parser.Parse(strings.NewReader(abiJSON), metadata)
	assert.NoError(t, err)
	assert.NotNil(t, contractIR)

	// Check functions
	assert.Len(t, contractIR.Functions, 3)
	
	// Check constructor
	constructor := contractIR.Functions[0]
	assert.Equal(t, "constructor", constructor.Name)
	assert.True(t, constructor.IsConstructor)
	assert.Equal(t, ir.Nonpayable, constructor.StateMutability)
	assert.Len(t, constructor.Inputs, 2)
	assert.Equal(t, "name", constructor.Inputs[0].Name)
	assert.Equal(t, "string", constructor.Inputs[0].Type.BaseType)
	assert.Equal(t, "symbol", constructor.Inputs[1].Name)
	assert.Equal(t, "string", constructor.Inputs[1].Type.BaseType)
	
	// Check fallback
	fallback := contractIR.Functions[1]
	assert.Equal(t, "fallback", fallback.Name)
	assert.True(t, fallback.IsFallback)
	assert.Equal(t, ir.Payable, fallback.StateMutability)
	
	// Check receive
	receive := contractIR.Functions[2]
	assert.Equal(t, "receive", receive.Name)
	assert.True(t, receive.IsReceive)
	assert.Equal(t, ir.Payable, receive.StateMutability)
}