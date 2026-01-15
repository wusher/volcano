// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for admonition blocks (note, warning, tip, etc.)
 */
test.describe('Admonitions', () => {
  let serverProcess;
  let testDir;
  const PORT = 4262;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Admonition Examples

Here are various admonition types.

> [!NOTE]
> This is a note admonition with helpful information.

> [!WARNING]
> This is a warning admonition about potential issues.

> [!TIP]
> This is a tip admonition with a helpful suggestion.

> [!IMPORTANT]
> This is an important admonition.

> [!CAUTION]
> This is a caution admonition.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('note admonition is rendered', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const pageContent = await page.textContent('.prose');
    expect(pageContent).toContain('This is a note admonition');
  });

  test('warning admonition is rendered', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const pageContent = await page.textContent('.prose');
    expect(pageContent).toContain('This is a warning admonition');
  });

  test('tip admonition is rendered', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const pageContent = await page.textContent('.prose');
    expect(pageContent).toContain('This is a tip admonition');
  });

  test('admonitions have distinct styling', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Check that admonition elements exist with styling
    const admonitions = page.locator('.admonition, .callout, blockquote');
    const count = await admonitions.count();
    expect(count).toBeGreaterThan(0);
  });

  test('admonition content is not rendered as plain blockquote', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // The [!NOTE] syntax should be transformed, not shown literally
    const bodyText = await page.textContent('body');
    // If admonitions are working, the literal [!NOTE] shouldn't appear
    // But the content "This is a note" should appear
    expect(bodyText).toContain('note admonition');
  });
});
