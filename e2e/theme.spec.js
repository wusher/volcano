// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for theme toggle (dark/light mode)
 */
test.describe('Theme Toggle', () => {
  let serverProcess;
  let testDir;
  const PORT = 4253;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Theme Test

Test page for theme switching.

${LOREM}
`,
      'other.md': `# Other Page

Another page.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('clicking theme toggle switches theme', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Get initial theme
    const initialTheme = await page.evaluate(() => document.documentElement.dataset.theme);

    // Click theme toggle
    await page.click('.theme-toggle');
    await page.waitForTimeout(300);

    // Theme should have changed
    const newTheme = await page.evaluate(() => document.documentElement.dataset.theme);
    expect(newTheme).not.toBe(initialTheme);
  });

  test('theme persists after page navigation', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Set to dark mode
    const isDark = await page.evaluate(() => document.documentElement.dataset.theme === 'dark');
    if (!isDark) {
      await page.click('.theme-toggle');
      await page.waitForTimeout(200);
    }

    // Verify it's dark
    const themeBefore = await page.evaluate(() => document.documentElement.dataset.theme);
    expect(themeBefore).toBe('dark');

    // Navigate to another page
    await page.click('.tree-nav a[href="/other/"]');
    await page.waitForTimeout(500);

    // Theme should still be dark
    const themeAfter = await page.evaluate(() => document.documentElement.dataset.theme);
    expect(themeAfter).toBe('dark');
  });

  test('theme persists after page reload', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Toggle to dark mode
    const initialTheme = await page.evaluate(() => document.documentElement.dataset.theme);
    if (initialTheme !== 'dark') {
      await page.click('.theme-toggle');
      await page.waitForTimeout(200);
    }

    // Reload page
    await page.reload();
    await page.waitForSelector('.theme-toggle');

    // Theme should still be dark
    const themeAfter = await page.evaluate(() => document.documentElement.dataset.theme);
    expect(themeAfter).toBe('dark');
  });

  test('theme toggle updates visual appearance', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Get background color in current theme
    const bgBefore = await page.evaluate(() =>
      getComputedStyle(document.body).backgroundColor
    );

    // Toggle theme
    await page.click('.theme-toggle');
    await page.waitForTimeout(300);

    // Background color should have changed
    const bgAfter = await page.evaluate(() =>
      getComputedStyle(document.body).backgroundColor
    );

    expect(bgAfter).not.toBe(bgBefore);
  });
});
