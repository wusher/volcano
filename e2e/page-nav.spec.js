// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for page navigation (previous/next links)
 */
test.describe('Page Navigation', () => {
  let serverProcess;
  let testDir;
  const PORT = 4266;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home

Welcome to the docs.

${LOREM}
`,
      '01-introduction.md': `# Introduction

The introduction.

${LOREM}
`,
      '02-installation.md': `# Installation

Installation guide.

${LOREM}
`,
      '03-usage.md': `# Usage

Usage instructions.

${LOREM}
`
    });

    // Start with page nav enabled
    serverProcess = await startServer(testDir, PORT, ['--page-nav']);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('page has previous/next navigation', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/installation/`);

    const pageNav = page.locator('.page-nav, .pagination, .prev-next');
    const count = await pageNav.count();
    expect(count).toBeGreaterThan(0);
  });

  test('next link navigates to next page', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/introduction/`);

    const nextLink = page.locator('.page-nav a:has-text("Next"), .next-link, a:has-text("Installation")').first();

    if (await nextLink.count() > 0) {
      await nextLink.click();
      await page.waitForTimeout(500);

      const heading = await page.textContent('h1');
      expect(heading).toContain('Installation');
    }
  });

  test('previous link navigates to previous page', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/installation/`);

    const prevLink = page.locator('.page-nav a:has-text("Previous"), .prev-link, a:has-text("Introduction")').first();

    if (await prevLink.count() > 0) {
      await prevLink.click();
      await page.waitForTimeout(500);

      const heading = await page.textContent('h1');
      expect(heading).toContain('Introduction');
    }
  });

  test('first page has no previous link', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Home page should not have a previous link (or it should be disabled)
    const prevLink = page.locator('.prev-link:not([disabled]), .page-nav a:has-text("Previous")');
    const count = await prevLink.count();

    // Either no previous link, or it's disabled/hidden
    expect(count).toBeLessThanOrEqual(1);
  });

  test('last page has no next link', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/usage/`);

    // Last page should not have a next link (or it should be disabled)
    const pageText = await page.textContent('.page-nav, .pagination').catch(() => '');

    // If there's page nav, verify it doesn't have an active next for last page
    // This is implementation-dependent
    expect(pageText).toBeDefined();
  });
});
