// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for instant navigation (hover prefetch + AJAX page loads)
 */
test.describe('Instant Navigation', () => {
  let serverProcess;
  let testDir;
  const PORT = 4251;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home Page

Welcome to the site.

${LOREM}
`,
      'about.md': `# About Page

This is the about page.

${LOREM}
`,
      'guide.md': `# Guide Page

This is a guide.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('clicking nav link loads page without full reload', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Set marker to detect page reload
    await page.evaluate(() => { window.__instantNavMarker = 'original'; });

    // Click on About link in navigation
    await page.click('.tree-nav a[href="/about/"]');
    await page.waitForTimeout(500);

    // Verify URL changed
    expect(page.url()).toContain('/about/');

    // Verify content changed
    const heading = await page.textContent('h1');
    expect(heading).toContain('About Page');

    // Verify page wasn't fully reloaded
    const marker = await page.evaluate(() => window.__instantNavMarker);
    expect(marker).toBe('original');
  });

  test('browser back button works after instant navigation', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Navigate to about page
    await page.click('.tree-nav a[href="/about/"]');
    await page.waitForTimeout(500);
    expect(page.url()).toContain('/about/');

    // Go back
    await page.goBack();
    await page.waitForTimeout(500);

    // Verify we're back on home page
    expect(page.url()).toBe(`http://localhost:${PORT}/`);
    const heading = await page.textContent('h1');
    expect(heading).toContain('Home Page');
  });

  test('browser forward button works after going back', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Navigate to about
    await page.click('.tree-nav a[href="/about/"]');
    await page.waitForTimeout(500);

    // Go back
    await page.goBack();
    await page.waitForTimeout(500);

    // Go forward
    await page.goForward();
    await page.waitForTimeout(500);

    expect(page.url()).toContain('/about/');
    const heading = await page.textContent('h1');
    expect(heading).toContain('About Page');
  });

  test('page title updates on navigation', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const homeTitle = await page.title();
    expect(homeTitle).toContain('Home Page');

    await page.click('.tree-nav a[href="/about/"]');
    await page.waitForTimeout(500);

    const aboutTitle = await page.title();
    expect(aboutTitle).toContain('About Page');
  });
});
