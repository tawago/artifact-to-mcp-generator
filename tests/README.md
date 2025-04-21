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

To run the tests:

```bash
./run-tests.sh
```

This script will:
1. Build the MCP server
2. Run the Playwright tests

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