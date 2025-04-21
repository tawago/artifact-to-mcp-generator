import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false, // Already set to false, which is good
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 3 : 2, // Add a retry for non-CI environments
  workers: 1, // Already set to 1, which ensures sequential execution
  reporter: 'html',
  use: {
    baseURL: `http://localhost:6274`, // Using the default port from the logs
    trace: 'on-first-retry',
    // Increase timeouts for better reliability
    navigationTimeout: 60000,
    actionTimeout: 30000,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: `npx @modelcontextprotocol/inspector node ../mcp-server/dist/server.js`,
    url: `http://localhost:6274`, // Using the default port from the logs
    reuseExistingServer: !process.env.CI,
    timeout: 180000, // Increase timeout to 3 minutes
  },
});