#!/bin/bash

# Run the tests without starting the MCP server (assumes it's already running)
npx playwright test --headed --config=playwright.config.manual.ts