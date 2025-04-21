import { test, expect } from '@playwright/test';

test.describe('MCP Server Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Click the Connect button
    await page.getByRole('button', { name: 'Connect' }).click();
    
    // Click the List Tools button (after connection, this button should be available)
    await page.getByRole('button', { name: 'List Tools' }).click();
  });

  test('should list all available tools', async ({ page }) => {
    // Check if all expected tools are listed (the beforeEach hook already handles connecting and listing tools)
    const toolNames = ['name', 'totalSupply', 'decimals', 'balanceOf', 'symbol', 'allowance'];
    
    for (const toolName of toolNames) {
      await expect(await page.getByText(toolName, { exact: true })).toBeVisible();
    }
  });

  test('should execute the name tool', async ({ page }) => {
    // Click on the name tool (the beforeEach hook already handles connecting and listing tools)
    await page.getByText('name', { exact: true }).click();
    
    // Execute the tool
    await page.getByRole('button', { name: 'Run Tool' }).click();
    
    // Wait for the result
    await expect(page.getByText('Tool Result: Success')).toBeVisible();
    await expect(page.getByText('USD Coin')).toBeVisible();
  });

  test('should execute the symbol tool', async ({ page }) => {
    // Click on the symbol tool (the beforeEach hook already handles connecting and listing tools)
    await page.getByText('symbol', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the symbol of the token');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text=', { timeout: 10000 });
    
    
    
  });

  test('should execute the decimals tool', async ({ page }) => {
    // Click on the decimals tool (the beforeEach hook already handles connecting and listing tools)
    await page.getByText('decimals', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the number of decimals the token uses');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text=', { timeout: 10000 });
  });

  test('should execute the totalSupply tool', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the totalSupply tool
    await page.getByText('totalSupply', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the total supply of the token');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result (should be a large number)
    await page.waitForSelector('text="', { timeout: 10000 });
  });

  test('should execute the balanceOf tool with a valid address', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the balanceOf tool
    await page.getByText('balanceOf', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the token balance of an account');
    
    // Fill in the address parameter
    await page.getByLabel('_owner').fill('0x1234567890123456789012345678901234567890');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text="', { timeout: 10000 });
  });

  test('should execute the allowance tool with valid addresses', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the allowance tool
    await page.getByText('allowance', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the amount of tokens that an owner allowed to a spender');
    
    // Fill in the address parameters
    await page.getByLabel('_owner').fill('0x1234567890123456789012345678901234567890');
    await page.getByLabel('_spender').fill('0x1234567890123456789012345678901234567890');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text="', { timeout: 10000 });
  });

  test('should handle errors for balanceOf with invalid address', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the balanceOf tool
    await page.getByText('balanceOf', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the token balance of an account');
    
    // Fill in an invalid address
    await page.getByLabel('_owner').fill('invalid-address');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the error result
    await page.waitForSelector('text=Error', { timeout: 10000 });
  });
});