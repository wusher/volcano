// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for TOC (Table of Contents) anchor navigation
 *
 * Verifies that clicking TOC links scrolls to the correct section
 * and doesn't jump back to top (bug fix verification).
 */
test.describe('TOC Anchor Navigation', () => {
  let serverProcess;
  let testDir;
  const PORT = 4250;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Test Page

This is intro content.

## First Section

${LOREM}

${LOREM}

## Second Section

${LOREM}

${LOREM}

## Third Section

${LOREM}

${LOREM}

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('clicking TOC link scrolls to section and stays there', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.toc-sidebar');

    const initialScrollY = await page.evaluate(() => window.scrollY);
    expect(initialScrollY).toBeLessThan(100);

    await page.click('.toc-sidebar a[href="#third-section"]');
    await page.waitForTimeout(500);

    const finalScrollY = await page.evaluate(() => window.scrollY);
    expect(finalScrollY).toBeGreaterThan(200);

    const thirdSectionRect = await page.evaluate(() => {
      const el = document.getElementById('third-section');
      return el ? el.getBoundingClientRect().top : null;
    });
    expect(thirdSectionRect).toBeLessThan(150);
    expect(page.url()).toContain('#third-section');
  });

  test('clicking TOC link does not trigger page reload', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.toc-sidebar');

    await page.evaluate(() => { window.__noReloadMarker = true; });
    await page.click('.toc-sidebar a[href="#second-section"]');
    await page.waitForTimeout(300);

    const markerExists = await page.evaluate(() => window.__noReloadMarker === true);
    expect(markerExists).toBe(true);
  });

  test('multiple TOC clicks navigate correctly', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.toc-sidebar');

    await page.click('.toc-sidebar a[href="#first-section"]');
    await page.waitForTimeout(300);
    const scroll1 = await page.evaluate(() => window.scrollY);
    expect(scroll1).toBeGreaterThan(0);

    await page.click('.toc-sidebar a[href="#third-section"]');
    await page.waitForTimeout(300);
    const scroll2 = await page.evaluate(() => window.scrollY);
    expect(scroll2).toBeGreaterThan(scroll1);

    await page.click('.toc-sidebar a[href="#first-section"]');
    await page.waitForTimeout(300);
    const scroll3 = await page.evaluate(() => window.scrollY);
    expect(scroll3).toBeLessThan(scroll2);
  });
});
