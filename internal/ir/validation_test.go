package ir

import (
	"testing"
)

func TestContractIRValidation(t *testing.T) {
	tests := []struct {
		name          string
		contract      ContractIR
		expectedError bool
		errorCount    int
	}{
		{
			name: "Valid Contract",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "transfer",
						StateMutability: Nonpayable,
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
							},
							{
								Name: "amount",
								Type: ParameterType{
									BaseType: "uint256",
								},
							},
						},
						Outputs: []Parameter{
							{
								Type: ParameterType{
									BaseType: "bool",
								},
							},
						},
					},
				},
				Events: []Event{
					{
						Name: "Transfer",
						Parameters: []EventParameter{
							{
								Name: "from",
								Type: ParameterType{
									BaseType: "address",
								},
								Indexed: true,
							},
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
								Indexed: true,
							},
							{
								Name: "value",
								Type: ParameterType{
									BaseType: "uint256",
								},
								Indexed: false,
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Missing Contract Name",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "transfer",
						StateMutability: Nonpayable,
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
							},
						},
						Outputs: []Parameter{
							{
								Type: ParameterType{
									BaseType: "bool",
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Missing Chain",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name: "TestContract",
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Function",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						// Missing name
						StateMutability: Nonpayable,
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Parameter Type",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "transfer",
						StateMutability: Nonpayable,
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									// Missing base type
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid State Mutability",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "transfer",
						StateMutability: "invalid",
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Visibility",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "transfer",
						StateMutability: Nonpayable,
						Visibility:     "invalid",
						Inputs: []Parameter{
							{
								Name: "to",
								Type: ParameterType{
									BaseType: "address",
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Event",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Events: []Event{
					{
						// Missing name
						Parameters: []EventParameter{
							{
								Name: "from",
								Type: ParameterType{
									BaseType: "address",
								},
								Indexed: true,
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Event Parameter",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Events: []Event{
					{
						Name: "Transfer",
						Parameters: []EventParameter{
							{
								// Missing name
								Type: ParameterType{
									BaseType: "address",
								},
								Indexed: true,
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Array Size",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "getArray",
						StateMutability: View,
						Outputs: []Parameter{
							{
								Type: ParameterType{
									BaseType:  "uint256",
									IsArray:   true,
									ArraySize: -1, // Invalid negative size
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid Map",
			contract: ContractIR{
				Metadata: ContractMetadata{
					Name:  "TestContract",
					Chain: "ethereum",
				},
				Functions: []Function{
					{
						Name:           "getMap",
						StateMutability: View,
						Outputs: []Parameter{
							{
								Type: ParameterType{
									BaseType: "uint256",
									IsMap:    true,
									// Missing map key type
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Multiple Errors",
			contract: ContractIR{
				Metadata: ContractMetadata{
					// Missing name
					// Missing chain
				},
				Functions: []Function{
					{
						// Missing name
						StateMutability: "invalid",
					},
				},
			},
			expectedError: true,
			errorCount:    4, // Missing name, chain, function name, invalid state mutability
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.contract.Validate()
			
			if tt.expectedError && len(errors) == 0 {
				t.Errorf("Expected validation errors but got none")
			}
			
			if !tt.expectedError && len(errors) > 0 {
				t.Errorf("Expected no validation errors but got %d: %v", len(errors), errors)
			}
			
			if tt.expectedError && tt.errorCount > 0 && len(errors) != tt.errorCount {
				t.Errorf("Expected %d validation errors but got %d: %v", tt.errorCount, len(errors), errors)
			}
		})
	}
}

func TestStateMutabilityValidation(t *testing.T) {
	tests := []struct {
		name          string
		mutability    StateMutability
		expectedError bool
	}{
		{"Pure", Pure, false},
		{"View", View, false},
		{"Nonpayable", Nonpayable, false},
		{"Payable", Payable, false},
		{"Empty", "", true},
		{"Invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStateMutability(tt.mutability)
			
			if tt.expectedError && err == nil {
				t.Errorf("Expected validation error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestVisibilityValidation(t *testing.T) {
	tests := []struct {
		name          string
		visibility    Visibility
		expectedError bool
	}{
		{"External", External, false},
		{"Public", Public, false},
		{"Internal", Internal, false},
		{"Private", Private, false},
		{"Invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVisibility(tt.visibility)
			
			if tt.expectedError && err == nil {
				t.Errorf("Expected validation error but got none")
			}
			
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}