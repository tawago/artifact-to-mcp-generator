import { test, expect } from '@playwright/test';

test.describe('MCP Server Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Click the Connect button
    await page.getByRole('button', { name: 'Connect' }).click();
    await expect(page.getByText('Connected')).toBeVisible();
    await expect(page.getByRole('button', { name: 'List Tools' })).toBeVisible();
    // Click the List Tools button (after connection, this button should be available)
    await page.getByRole('button', { name: 'List Tools' }).click();
  });

  test.afterEach(async ({ page }) => {
    await page.close();
  });

  test('should list all available tools', async ({ page }) => {
    // Check if all expected tools are listed (the beforeEach hook already handles connecting and listing tools)
    const toolNames = ['name', 'totalSupply', 'decimals', 'balanceOf', 'symbol', 'allowance'];
    
    for (const toolName of toolNames) {
      await expect(await page.getByText(toolName, { exact: true })).toBeVisible();
    }
  });

  // Test all simple tools without parameters using a loop
  const simpleTools = {
    'name': 'USD Coin',
    'totalSupply': '\\d{17,}',
    'decimals':'\"6\"',
    'symbol': 'USDC',
  };
  
  for (const [tool, value] of Object.entries(simpleTools)) {
    test(`should execute the ${tool} tool`, async ({ page }) => {
      // Click on the tool (the beforeEach hook already handles connecting and listing tools)
      await page.getByText(tool, { exact: true }).click();

      // Execute the tool
      await page.getByRole('button', { name: 'Run Tool' }).click();
      
      // Wait for the result
      await expect(page.getByText('Tool Result: Success')).toBeVisible();
      await expect(page.getByText(RegExp(value))).toBeVisible();
    });
  }

  test('should execute the balanceOf tool with a valid address', async ({ page }) => {
    // Click on the balanceOf tool
    await page.getByText('balanceOf', { exact: true }).click();
    
    // Fill in the address parameter
    await page.getByLabel('_owner').fill('0x1234567890123456789012345678901234567890');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Run Tool' }).click();
    
    // Wait for the result
    await expect(page.getByText('Tool Result: Success')).toBeVisible();
  });

  test('should execute the allowance tool with valid addresses', async ({ page }) => {
    // Click on the allowance tool
    await page.getByText('allowance', { exact: true }).click();
    
    // Fill in the address parameters
    await page.getByLabel('_owner').fill('0x1234567890123456789012345678901234567890');
    await page.getByLabel('_spender').fill('0x1234567890123456789012345678901234567890');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Run Tool' }).click();
    
    // Wait for the result
    await expect(page.getByText('Tool Result: Success')).toBeVisible();
  });

  test('should handle errors for balanceOf with invalid address', async ({ page }) => {
    // Click on the balanceOf tool
    await page.getByText('balanceOf', { exact: true }).click();
    
    // Fill in an invalid address
    await page.getByLabel('_owner').fill('invalid-address');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Run Tool' }).click();
    
    // Wait for the error result
    await expect(page.getByText(/Error/)).toBeVisible();
  });
});