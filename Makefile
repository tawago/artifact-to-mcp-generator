.PHONY: build test clean example

# Build the CLI tool
build:
	go build -o bin/generate-mcp ./cmd/generate-mcp

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf mcp-server/

# Generate an example MCP server for ERC-20
example: build
	./bin/generate-mcp --artifact examples/erc20.json --name "ERC20 Token" --address "0x1234567890123456789012345678901234567890" --output ./mcp-server

# Install the tool
install:
	go install ./cmd/generate-mcp

# Run the example MCP server
run-example:
	cd mcp-server && npm install && npm run build && npm start