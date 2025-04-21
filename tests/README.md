# MCP Server Tests

This directory contains end-to-end tests for the generated MCP server using Playwright.

## Prerequisites

- Node.js (v16 or later)
- npm

## Setup

1. Install dependencies:

```bash
npm install
```

2. Install Playwright browsers:

```bash
npx playwright install --with-deps chromium
```

## Running Tests

### Automated Testing

To run the tests automatically:

```bash
./run-tests.sh
```

This script will:
1. Build the MCP server
2. Run the Playwright tests

### Manual Testing

For manual testing, you can run the MCP server and tests separately:

1. Start the MCP server with the inspector:

```bash
./run-mcp-server.sh
```

2. In a separate terminal, run the tests:

```bash
./run-tests-manual.sh
```

This approach is useful for debugging as you can see the MCP Inspector UI in action.

## Test Structure

The tests verify that:

1. The MCP server connects successfully
2. All tools are listed correctly
3. Each tool can be executed with valid parameters
4. Error handling works as expected

## Test Results

Test results are stored in the following locations:

- HTML report: `playwright-report/`
- Screenshots: `./*.png`

## Debugging

To run tests in debug mode:

```bash
npm run test:debug
```

To run tests with the Playwright UI:

```bash
npm run test:ui
```