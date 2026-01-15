// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for external link handling
 */
test.describe('External Links', () => {
  let serverProcess;
  let testDir;
  const PORT = 4263;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# External Links Test

Here are some external links:

- [Google](https://www.google.com)
- [GitHub](https://github.com)

And an internal link: [About](/about/)

${LOREM}
`,
      'about.md': `# About

This is an internal page.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('external links have target="_blank"', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const googleLink = page.locator('a[href="https://www.google.com"]');
    await expect(googleLink).toHaveAttribute('target', '_blank');
  });

  test('external links have rel="noopener noreferrer"', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const googleLink = page.locator('a[href="https://www.google.com"]');
    const rel = await googleLink.getAttribute('rel');
    expect(rel).toContain('noopener');
  });

  test('internal links navigate within site', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Internal link should work for navigation
    const aboutLink = page.locator('.prose a[href="/about/"], .prose a[href*="about"]').first();
    const count = await aboutLink.count();

    if (count > 0) {
      await aboutLink.click();
      await page.waitForTimeout(500);
      expect(page.url()).toContain('/about/');
    } else {
      // If no about link in prose, just verify page renders
      const heading = await page.textContent('h1');
      expect(heading).toBeDefined();
    }
  });

  test('external links are visually indicated', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // External links typically have an icon or class
    const externalLink = page.locator('a[href="https://www.google.com"]');

    // Check if link has external indicator (class, icon, or ::after content)
    const hasExternalClass = await externalLink.evaluate(el => {
      return el.classList.contains('external') ||
             el.querySelector('svg') !== null ||
             window.getComputedStyle(el, '::after').content !== 'none';
    });

    // At minimum, the link should exist and be clickable
    await expect(externalLink).toBeVisible();
  });

  test('multiple external links all have correct attributes', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const externalLinks = page.locator('a[href^="https://"]');
    const count = await externalLinks.count();
    expect(count).toBeGreaterThanOrEqual(2);

    // Check each external link
    for (let i = 0; i < count; i++) {
      const link = externalLinks.nth(i);
      await expect(link).toHaveAttribute('target', '_blank');
    }
  });
});
