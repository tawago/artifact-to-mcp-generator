---
name: Phase 1 - Core EVM Implementation
about: Parent issue for tracking Phase 1 (EVM only) implementation
title: 'Phase 1: Core EVM Implementation'
labels: 'enhancement, parent-issue'
assignees: ''

---

# Phase 1: Core EVM Implementation

This is a parent issue for tracking the implementation of Phase 1 of the Smart-Contract MCP Server Generator, focusing on Ethereum Virtual Machine (EVM) support.

## Overview

The goal of Phase 1 is to build a working prototype that can generate MCP servers from EVM smart contract ABIs. This will allow LLMs to interact with Ethereum-compatible smart contracts through the Model Context Protocol.

## Tasks

### 1. Define IR Schema (JSON) v0.1 in Go structs/enums

- [ ] Define the Intermediate Representation (IR) schema as Go structs
- [ ] Include support for contract metadata (name, address, chain)
- [ ] Define function representation (name, inputs, outputs, mutability)
- [ ] Define event representation (name, parameters, indexed fields)
- [ ] Add validation for the IR schema
- [ ] Write tests for IR schema serialization/deserialization

### 2. Write EVM Importer - Parse ABI JSON → IR

- [ ] Implement ABI JSON parser
- [ ] Map ABI function definitions to IR
- [ ] Map ABI event definitions to IR
- [ ] Handle function visibility and mutability
- [ ] Support constructor and fallback functions
- [ ] Write tests with sample ABIs (ERC20, ERC721, etc.)

### 3. Draft TypeScript MCP Server Template

- [ ] Create base template structure using text/template
- [ ] Implement tool definitions based on contract functions
- [ ] Add proper type definitions for inputs/outputs
- [ ] Include web3 provider configuration
- [ ] Add error handling for contract interactions
- [ ] Support read-only and state-changing functions

### 4. Generate Sample MCP Server for a Public ERC-20

- [ ] Select a well-known ERC-20 contract (e.g., USDC, DAI)
- [ ] Generate MCP server from its ABI
- [ ] Implement basic connection to an Ethereum provider
- [ ] Test read functions (balanceOf, totalSupply, etc.)
- [ ] Test event listening capabilities
- [ ] Document the generated server usage

### 5. Handle Edge Cases

- [ ] Support function overloads
- [ ] Handle payable functions
- [ ] Support indexed event parameters
- [ ] Add gas estimation for state-changing functions
- [ ] Handle complex parameter types (structs, arrays)
- [ ] Support contract inheritance

### 6. CLI Wrapper

- [ ] Implement CLI using cobra
- [ ] Add artifact path flag
- [ ] Add language selection flag (initially just TypeScript)
- [ ] Add output directory flag
- [ ] Support configuration file for default settings
- [ ] Add verbose/debug mode
- [ ] Write usage documentation

## Definition of Done

- All tasks are completed and tested
- The tool can generate a working MCP server from any valid EVM ABI
- Generated servers correctly expose contract functions as MCP tools
- CLI provides a user-friendly interface for generating servers
- Documentation is complete and accurate

## Resources

- [MCP Protocol Specification](https://github.com/anthropics/model-context-protocol)
- [Ethereum ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [Example MCP Servers](https://github.com/anthropics/model-context-protocol/tree/main/src)