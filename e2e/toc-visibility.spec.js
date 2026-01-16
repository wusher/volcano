// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for TOC visibility during navigation
 *
 * Bug: When loading a page without TOC, then navigating to a page with TOC,
 * the TOC doesn't become visible because the .has-toc class isn't updated
 * on the .main-wrapper element.
 */
test.describe('TOC Visibility', () => {
  let serverProcess;
  let testDir;
  const PORT = 4260;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home Page

This is a simple page without any headings, so no TOC will be generated.

${LOREM}

Just some basic content here.
`,
      'guide.md': `# Guide Page

This page has multiple sections with h2 headings, which should generate a TOC.

## Section One

${LOREM}

## Section Two

${LOREM}

## Section Three

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('TOC should appear when navigating from page without TOC to page with TOC', async ({ page }) => {
    // Start on a page without TOC
    await page.goto(`http://localhost:${PORT}/`);

    // Verify no TOC on home page
    const hasTocClassBefore = await page.evaluate(() => {
      const mainWrapper = document.querySelector('.main-wrapper');
      return mainWrapper ? mainWrapper.classList.contains('has-toc') : false;
    });
    expect(hasTocClassBefore).toBe(false);

    const tocVisibleBefore = await page.isVisible('.toc-sidebar');
    expect(tocVisibleBefore).toBe(false);

    // Navigate to page with TOC using instant nav
    await page.click('.tree-nav a[href="/guide/"]');
    await page.waitForTimeout(500);

    // Verify we're on the guide page
    const heading = await page.textContent('h1');
    expect(heading).toContain('Guide Page');

    // Verify TOC is now visible
    const hasTocClassAfter = await page.evaluate(() => {
      const mainWrapper = document.querySelector('.main-wrapper');
      return mainWrapper ? mainWrapper.classList.contains('has-toc') : false;
    });
    expect(hasTocClassAfter).toBe(true);

    const tocVisibleAfter = await page.isVisible('.toc-sidebar');
    expect(tocVisibleAfter).toBe(true);

    // Verify TOC contains the expected sections
    const tocContent = await page.textContent('.toc-sidebar');
    expect(tocContent).toContain('Section One');
    expect(tocContent).toContain('Section Two');
    expect(tocContent).toContain('Section Three');
  });

  test('TOC should disappear when navigating from page with TOC to page without TOC', async ({ page }) => {
    // Start on a page with TOC
    await page.goto(`http://localhost:${PORT}/guide/`);

    // Verify TOC is visible on guide page
    const hasTocClassBefore = await page.evaluate(() => {
      const mainWrapper = document.querySelector('.main-wrapper');
      return mainWrapper ? mainWrapper.classList.contains('has-toc') : false;
    });
    expect(hasTocClassBefore).toBe(true);

    const tocVisibleBefore = await page.isVisible('.toc-sidebar');
    expect(tocVisibleBefore).toBe(true);

    // Navigate to page without TOC - use site title link instead of navigation
    await page.click('.site-title');
    await page.waitForTimeout(500);

    // Verify we're on the home page
    const heading = await page.textContent('h1');
    expect(heading).toContain('Home Page');

    // Verify TOC is now hidden
    const hasTocClassAfter = await page.evaluate(() => {
      const mainWrapper = document.querySelector('.main-wrapper');
      return mainWrapper ? mainWrapper.classList.contains('has-toc') : false;
    });
    expect(hasTocClassAfter).toBe(false);

    const tocVisibleAfter = await page.isVisible('.toc-sidebar');
    expect(tocVisibleAfter).toBe(false);
  });

  test('mobile TOC button should only appear on pages with TOC', async ({ page }) => {
    // Set viewport to mobile size
    await page.setViewportSize({ width: 375, height: 667 });

    // Start on a page without TOC
    await page.goto(`http://localhost:${PORT}/`);

    // Verify mobile TOC button is not visible
    const mobileButtonVisibleBefore = await page.isVisible('.mobile-toc-toggle');
    expect(mobileButtonVisibleBefore).toBe(false);

    // Navigate to page with TOC - open drawer first
    await page.click('.mobile-menu-btn');
    await page.waitForTimeout(300);
    await page.click('.tree-nav a[href="/guide/"]');
    await page.waitForTimeout(500);

    // Verify mobile TOC button is now visible
    const mobileButtonVisibleAfter = await page.isVisible('.mobile-toc-toggle');
    expect(mobileButtonVisibleAfter).toBe(true);
  });
});
