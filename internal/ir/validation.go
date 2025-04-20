package ir

import (
	"fmt"
	"strings"
)

// ValidationError represents an error found during IR validation
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate checks if the ContractIR is valid and returns a list of validation errors
func (c *ContractIR) Validate() []ValidationError {
	var errors []ValidationError

	// Validate metadata
	metadataErrors := c.Metadata.Validate()
	errors = append(errors, metadataErrors...)

	// Validate functions
	for i, function := range c.Functions {
		fieldPrefix := fmt.Sprintf("Functions[%d]", i)
		functionErrors := function.Validate()
		for _, err := range functionErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	// Validate events
	for i, event := range c.Events {
		fieldPrefix := fmt.Sprintf("Events[%d]", i)
		eventErrors := event.Validate()
		for _, err := range eventErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	// Validate errors
	for i, contractError := range c.Errors {
		fieldPrefix := fmt.Sprintf("Errors[%d]", i)
		errorErrors := contractError.Validate()
		for _, err := range errorErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	// Validate custom types
	for i, customType := range c.Types {
		fieldPrefix := fmt.Sprintf("Types[%d]", i)
		typeErrors := customType.Validate()
		for _, err := range typeErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	return errors
}

// Validate checks if the ContractMetadata is valid and returns a list of validation errors
func (m *ContractMetadata) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(m.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "contract name is required",
		})
	}

	// Chain is required
	if strings.TrimSpace(m.Chain) == "" {
		errors = append(errors, ValidationError{
			Field:   "Chain",
			Message: "chain identifier is required",
		})
	}

	// Validate source info if present
	if m.Source != nil {
		sourceErrors := m.Source.Validate()
		for _, err := range sourceErrors {
			err.Field = "Source." + err.Field
			errors = append(errors, err)
		}
	}

	return errors
}

// Validate checks if the SourceInfo is valid and returns a list of validation errors
func (s *SourceInfo) Validate() []ValidationError {
	var errors []ValidationError

	// Language is required
	if strings.TrimSpace(s.Language) == "" {
		errors = append(errors, ValidationError{
			Field:   "Language",
			Message: "programming language is required",
		})
	}

	return errors
}

// Validate checks if the Function is valid and returns a list of validation errors
func (f *Function) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(f.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "function name is required",
		})
	}

	// Validate inputs
	for i, input := range f.Inputs {
		fieldPrefix := fmt.Sprintf("Inputs[%d]", i)
		inputErrors := input.Validate()
		for _, err := range inputErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	// Validate outputs
	for i, output := range f.Outputs {
		fieldPrefix := fmt.Sprintf("Outputs[%d]", i)
		outputErrors := output.Validate()
		for _, err := range outputErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	// Validate state mutability
	if err := validateStateMutability(f.StateMutability); err != nil {
		errors = append(errors, ValidationError{
			Field:   "StateMutability",
			Message: err.Error(),
		})
	}

	// Validate visibility if present
	if f.Visibility != "" {
		if err := validateVisibility(f.Visibility); err != nil {
			errors = append(errors, ValidationError{
				Field:   "Visibility",
				Message: err.Error(),
			})
		}
	}

	return errors
}

// Validate checks if the Event is valid and returns a list of validation errors
func (e *Event) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(e.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "event name is required",
		})
	}

	// Validate parameters
	for i, param := range e.Parameters {
		fieldPrefix := fmt.Sprintf("Parameters[%d]", i)
		paramErrors := param.Validate()
		for _, err := range paramErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	return errors
}

// Validate checks if the EventParameter is valid and returns a list of validation errors
func (p *EventParameter) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(p.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "parameter name is required",
		})
	}

	// Validate type
	typeErrors := p.Type.Validate()
	for _, err := range typeErrors {
		err.Field = "Type." + err.Field
		errors = append(errors, err)
	}

	return errors
}

// Validate checks if the Parameter is valid and returns a list of validation errors
func (p *Parameter) Validate() []ValidationError {
	var errors []ValidationError

	// Validate type
	typeErrors := p.Type.Validate()
	for _, err := range typeErrors {
		err.Field = "Type." + err.Field
		errors = append(errors, err)
	}

	return errors
}

// Validate checks if the ParameterType is valid and returns a list of validation errors
func (t *ParameterType) Validate() []ValidationError {
	var errors []ValidationError

	// BaseType is required
	if strings.TrimSpace(t.BaseType) == "" {
		errors = append(errors, ValidationError{
			Field:   "BaseType",
			Message: "base type is required",
		})
	}

	// If it's an array, validate array size
	if t.IsArray && t.ArraySize < 0 {
		errors = append(errors, ValidationError{
			Field:   "ArraySize",
			Message: "array size must be non-negative (0 for dynamic arrays)",
		})
	}

	// If it's a map, validate key type
	if t.IsMap && strings.TrimSpace(t.MapKeyType) == "" {
		errors = append(errors, ValidationError{
			Field:   "MapKeyType",
			Message: "map key type is required for maps",
		})
	}

	// Validate components if present
	if len(t.Components) > 0 {
		for i, component := range t.Components {
			fieldPrefix := fmt.Sprintf("Components[%d]", i)
			componentErrors := component.Validate()
			for _, err := range componentErrors {
				err.Field = fieldPrefix + "." + err.Field
				errors = append(errors, err)
			}
		}
	}

	return errors
}

// Validate checks if the ContractError is valid and returns a list of validation errors
func (e *ContractError) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(e.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "error name is required",
		})
	}

	// Validate parameters if present
	for i, param := range e.Parameters {
		fieldPrefix := fmt.Sprintf("Parameters[%d]", i)
		paramErrors := param.Validate()
		for _, err := range paramErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	return errors
}

// Validate checks if the CustomType is valid and returns a list of validation errors
func (t *CustomType) Validate() []ValidationError {
	var errors []ValidationError

	// Name is required
	if strings.TrimSpace(t.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "Name",
			Message: "type name is required",
		})
	}

	// Validate fields
	for i, field := range t.Fields {
		fieldPrefix := fmt.Sprintf("Fields[%d]", i)
		fieldErrors := field.Validate()
		for _, err := range fieldErrors {
			err.Field = fieldPrefix + "." + err.Field
			errors = append(errors, err)
		}
	}

	return errors
}

// validateStateMutability checks if the state mutability is valid
func validateStateMutability(sm StateMutability) error {
	switch sm {
	case Pure, View, Nonpayable, Payable:
		return nil
	case "":
		return fmt.Errorf("state mutability is required")
	default:
		return fmt.Errorf("invalid state mutability: %s", sm)
	}
}

// validateVisibility checks if the visibility is valid
func validateVisibility(v Visibility) error {
	switch v {
	case External, Public, Internal, Private:
		return nil
	default:
		return fmt.Errorf("invalid visibility: %s", v)
	}
}