package ir

// Schema defines the Intermediate Representation (IR) for smart contract definitions
// This IR is used to generate MCP servers for different blockchain platforms

// ContractIR represents a smart contract in the intermediate representation
type ContractIR struct {
        // Metadata about the contract
        Metadata ContractMetadata `json:"metadata"`
        
        // Functions defined in the contract
        Functions []Function `json:"functions"`
        
        // Events that can be emitted by the contract
        Events []Event `json:"events"`
        
        // Errors that can be thrown by the contract
        Errors []ContractError `json:"errors,omitempty"`
        
        // Custom types defined in the contract
        Types []CustomType `json:"types,omitempty"`
}

// ContractMetadata contains information about the contract itself
type ContractMetadata struct {
        // Name of the contract
        Name string `json:"name"`
        
        // Description of the contract's purpose
        Description string `json:"description,omitempty"`
        
        // Address where the contract is deployed (if known)
        Address string `json:"address,omitempty"`
        
        // Chain identifier (e.g., "ethereum", "solana")
        Chain string `json:"chain"`
        
        // Chain-specific information
        ChainData map[string]interface{} `json:"chainData,omitempty"`
        
        // Source code information
        Source *SourceInfo `json:"source,omitempty"`
}

// SourceInfo contains information about the contract's source code
type SourceInfo struct {
        // Programming language
        Language string `json:"language"`
        
        // Compiler version
        Compiler string `json:"compiler,omitempty"`
        
        // Source code URL or path
        SourceURL string `json:"sourceUrl,omitempty"`
}

// Function represents a callable function in the contract
type Function struct {
        // Function name
        Name string `json:"name"`
        
        // Human-readable description
        Description string `json:"description,omitempty"`
        
        // Function signature (e.g., "transfer(address,uint256)")
        Signature string `json:"signature,omitempty"`
        
        // Function selector (e.g., "0xa9059cbb" for EVM)
        Selector string `json:"selector,omitempty"`
        
        // Input parameters
        Inputs []Parameter `json:"inputs"`
        
        // Output parameters
        Outputs []Parameter `json:"outputs"`
        
        // Whether the function modifies state
        StateMutability StateMutability `json:"stateMutability"`
        
        // Function visibility
        Visibility Visibility `json:"visibility,omitempty"`
        
        // Whether this is a constructor
        IsConstructor bool `json:"isConstructor,omitempty"`
        
        // Whether this is a fallback function
        IsFallback bool `json:"isFallback,omitempty"`
        
        // Whether this is a receive function (EVM specific)
        IsReceive bool `json:"isReceive,omitempty"`
        
        // Chain-specific function data
        ChainData map[string]interface{} `json:"chainData,omitempty"`
}

// Event represents an event that can be emitted by the contract
type Event struct {
        // Event name
        Name string `json:"name"`
        
        // Human-readable description
        Description string `json:"description,omitempty"`
        
        // Event signature
        Signature string `json:"signature,omitempty"`
        
        // Parameters included in the event
        Parameters []EventParameter `json:"parameters"`
        
        // Chain-specific event data
        ChainData map[string]interface{} `json:"chainData,omitempty"`
}

// EventParameter represents a parameter in an event
type EventParameter struct {
        // Parameter name
        Name string `json:"name"`
        
        // Parameter type
        Type ParameterType `json:"type"`
        
        // Whether the parameter is indexed (for efficient filtering)
        Indexed bool `json:"indexed"`
}

// Parameter represents a function parameter (input or output)
type Parameter struct {
        // Parameter name
        Name string `json:"name"`
        
        // Parameter type
        Type ParameterType `json:"type"`
        
        // Human-readable description
        Description string `json:"description,omitempty"`
}

// ParameterType represents the type of a parameter
type ParameterType struct {
        // Base type (e.g., "uint256", "address", "string")
        BaseType string `json:"baseType"`
        
        // Whether this is an array
        IsArray bool `json:"isArray,omitempty"`
        
        // Fixed array size (0 means dynamic)
        ArraySize int `json:"arraySize,omitempty"`
        
        // Whether this is a map/dictionary
        IsMap bool `json:"isMap,omitempty"`
        
        // If this is a map, the key type
        MapKeyType string `json:"mapKeyType,omitempty"`
        
        // If this is a custom struct type, the fields
        Components []Parameter `json:"components,omitempty"`
        
        // Chain-specific type data
        ChainData map[string]interface{} `json:"chainData,omitempty"`
}

// ContractError represents a custom error that can be thrown by the contract
type ContractError struct {
        // Error name
        Name string `json:"name"`
        
        // Human-readable description
        Description string `json:"description,omitempty"`
        
        // Error parameters
        Parameters []Parameter `json:"parameters,omitempty"`
}

// CustomType represents a custom type defined in the contract
type CustomType struct {
        // Type name
        Name string `json:"name"`
        
        // Human-readable description
        Description string `json:"description,omitempty"`
        
        // Fields in the custom type
        Fields []Parameter `json:"fields"`
}

// StateMutability indicates how a function interacts with contract state
type StateMutability string

const (
        // Pure functions don't read or modify state
        Pure StateMutability = "pure"
        
        // View functions read but don't modify state
        View StateMutability = "view"
        
        // Nonpayable functions modify state but don't accept Ether
        Nonpayable StateMutability = "nonpayable"
        
        // Payable functions modify state and accept Ether
        Payable StateMutability = "payable"
)

// Visibility indicates the visibility of a function
type Visibility string

const (
        // External functions can only be called from outside the contract
        External Visibility = "external"
        
        // Public functions can be called from inside or outside the contract
        Public Visibility = "public"
        
        // Internal functions can only be called from inside the contract or derived contracts
        Internal Visibility = "internal"
        
        // Private functions can only be called from inside the contract
        Private Visibility = "private"
)