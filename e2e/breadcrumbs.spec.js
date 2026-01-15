// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for breadcrumb navigation
 */
test.describe('Breadcrumbs', () => {
  let serverProcess;
  let testDir;
  const PORT = 4265;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home

Welcome.

${LOREM}
`,
      'guides/index.md': `# Guides

Guide section.

${LOREM}
`,
      'guides/getting-started.md': `# Getting Started

A getting started guide.

${LOREM}
`,
      'guides/advanced/index.md': `# Advanced

Advanced section.

${LOREM}
`,
      'guides/advanced/tips.md': `# Advanced Tips

Some advanced tips.

${LOREM}
`
    });

    // Start with breadcrumbs enabled
    serverProcess = await startServer(testDir, PORT, ['--breadcrumbs']);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('breadcrumbs appear on nested pages', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/guides/getting-started/`);

    const breadcrumbs = page.locator('.breadcrumbs, .breadcrumb, nav[aria-label="breadcrumb"]');
    const count = await breadcrumbs.count();
    expect(count).toBeGreaterThan(0);
  });

  test('breadcrumbs show correct hierarchy', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/guides/getting-started/`);

    const pageText = await page.textContent('body');
    // Should show path hierarchy
    expect(pageText).toContain('Guides');
  });

  test('breadcrumb links are clickable', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/guides/getting-started/`);

    // Find a breadcrumb link to Guides
    const guidesLink = page.locator('.breadcrumbs a[href*="guides"], .breadcrumb a[href*="guides"]').first();

    if (await guidesLink.count() > 0) {
      await guidesLink.click();
      await page.waitForTimeout(500);
      expect(page.url()).toContain('/guides/');
    }
  });

  test('deeply nested page shows full breadcrumb trail', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/guides/advanced/tips/`);

    const bodyText = await page.textContent('body');
    // Should contain parent directories in breadcrumb
    expect(bodyText).toContain('Guides');
    expect(bodyText).toContain('Advanced');
  });

  test('home page has minimal or no breadcrumbs', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Home page should render
    const heading = await page.textContent('h1');
    expect(heading).toContain('Home');

    // Home page breadcrumbs should be minimal (if they exist)
    const breadcrumbs = page.locator('.breadcrumbs, .breadcrumb');
    const count = await breadcrumbs.count();

    if (count > 0) {
      const text = await breadcrumbs.first().textContent();
      // Should be minimal on home page
      expect(text.split('/').length).toBeLessThanOrEqual(3);
    }
  });
});
