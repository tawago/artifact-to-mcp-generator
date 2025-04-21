import { test, expect } from '@playwright/test';

test.describe('MCP Server Tests', () => {
  test('should connect to the MCP server', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Check if the server is connected
    const serverInfo = await page.getByText('ERC20 Token').first();
    await expect(serverInfo).toBeVisible();
    
    // Take a screenshot of the connected server
    await page.screenshot({ path: 'server-connected.png' });
  });

  test('should list all available tools', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Check if all expected tools are listed
    const toolNames = ['name', 'totalSupply', 'decimals', 'balanceOf', 'symbol', 'allowance'];
    
    for (const toolName of toolNames) {
      const toolElement = await page.getByText(toolName, { exact: true }).first();
      await expect(toolElement).toBeVisible();
    }
    
    // Take a screenshot of the tools list
    await page.screenshot({ path: 'tools-list.png' });
  });

  test('should execute the name tool', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the name tool
    await page.getByText('name', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the name of the token');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text=', { timeout: 10000 });
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'name-tool-result.png' });
  });

  test('should execute the symbol tool', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the symbol tool
    await page.getByText('symbol', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the symbol of the token');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text=', { timeout: 10000 });
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'symbol-tool-result.png' });
  });

  test('should execute the decimals tool', async ({ page }) => {
    // Navigate to the MCP Inspector
    await page.goto('/');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Click on the decimals tool
    await page.getByText('decimals', { exact: true }).first().click();
    
    // Wait for the tool details to load
    await page.waitForSelector('text=Get the number of decimals the token uses');
    
    // Execute the tool
    await page.getByRole('button', { name: 'Execute' }).click();
    
    // Wait for the result
    await page.waitForSelector('text=', { timeout: 10000 });
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'decimals-tool-result.png' });
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
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'totalSupply-tool-result.png' });
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
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'balanceOf-tool-result.png' });
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
    
    // Take a screenshot of the result
    await page.screenshot({ path: 'allowance-tool-result.png' });
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
    
    // Take a screenshot of the error
    await page.screenshot({ path: 'balanceOf-error-result.png' });
  });
});