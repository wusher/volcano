// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for search/command palette functionality
 */
test.describe('Search Command Palette', () => {
  let serverProcess;
  let testDir;
  const PORT = 4252;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home Page

Welcome to the documentation site.

${LOREM}
`,
      'getting-started.md': `# Getting Started

This guide helps you get started quickly.

${LOREM}
`,
      'api-reference.md': `# API Reference

Complete API documentation.

${LOREM}
`,
      'guides/advanced.md': `# Advanced Guide

Advanced topics and techniques.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('Cmd+K opens command palette', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.tree-nav');

    // Command palette should not be visible initially
    const paletteHidden = await page.isHidden('.command-palette');
    expect(paletteHidden).toBe(true);

    // Press Cmd+K (or Ctrl+K on non-Mac)
    await page.keyboard.press('Meta+k');
    await page.waitForTimeout(300);

    // Command palette should now be visible
    const paletteVisible = await page.isVisible('.command-palette');
    expect(paletteVisible).toBe(true);
  });

  test('typing in search filters results', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Type a search query
    await page.fill('#command-palette-input', 'api');
    await page.waitForTimeout(300);

    // Should show API Reference result
    const results = await page.locator('.command-palette-result').count();
    expect(results).toBeGreaterThan(0);

    const resultText = await page.textContent('.command-palette-results');
    expect(resultText).toContain('API Reference');
  });

  test('clicking search result navigates to page', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for getting started
    await page.fill('#command-palette-input', 'getting');
    await page.waitForTimeout(300);

    // Click the result
    await page.click('.command-palette-result');
    await page.waitForTimeout(500);

    // Should navigate to the page
    expect(page.url()).toContain('/getting-started/');
    const heading = await page.textContent('h1');
    expect(heading).toContain('Getting Started');
  });

  test('Escape closes command palette', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Press Escape
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);

    // Command palette should be closed
    const paletteHidden = await page.isHidden('.command-palette') ||
      !(await page.locator('.command-palette.open').count());
    expect(paletteHidden).toBe(true);
  });
});
