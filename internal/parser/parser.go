package parser

import (
	"io"

	"github.com/openhands/mcp-generator/internal/ir"
	"github.com/openhands/mcp-generator/internal/parser/evm"
)

// Parser is the interface for all contract artifact parsers
type Parser interface {
	// Parse parses a contract artifact into the intermediate representation
	Parse(reader io.Reader, metadata ir.ContractMetadata) (*ir.ContractIR, error)
}

// NewEVMABIParser creates a new EVM ABI parser
func NewEVMABIParser() Parser {
	return evm.NewABIParser()
}