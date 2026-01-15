// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for 404 page handling
 */
test.describe('404 Page', () => {
  let serverProcess;
  let testDir;
  const PORT = 4264;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home

Welcome to the site.

${LOREM}
`,
      'about.md': `# About

About page.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('non-existent page shows 404 content', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/this-page-does-not-exist/`);

    const bodyText = await page.textContent('body');
    // Should show some indication of not found
    expect(bodyText.toLowerCase()).toMatch(/not found|404|doesn't exist/);
  });

  test('404 page still has site navigation', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/nonexistent-page/`);

    // Navigation should still be present
    const nav = page.locator('.tree-nav, nav');
    const count = await nav.count();
    expect(count).toBeGreaterThan(0);
  });

  test('404 page has link back to home', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/fake-page/`);

    // Should have a way to get back to home
    const homeLink = page.locator('a[href="/"], a[href="./"]');
    const bodyText = await page.textContent('body');

    // Either has home link or navigation is present
    const hasNavigation = await page.locator('.tree-nav').count() > 0;
    expect(hasNavigation || await homeLink.count() > 0).toBe(true);
  });

  test('deeply nested non-existent path shows 404', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/a/b/c/d/nonexistent/`);

    const bodyText = await page.textContent('body');
    expect(bodyText.toLowerCase()).toMatch(/not found|404|doesn't exist/);
  });
});
