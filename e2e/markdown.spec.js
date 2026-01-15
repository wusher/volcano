// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for markdown rendering features
 */
test.describe('Markdown Features', () => {
  let serverProcess;
  let testDir;
  const PORT = 4268;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Markdown Features

Testing various markdown elements.

## Tables

| Name | Age | City |
|------|-----|------|
| Alice | 30 | NYC |
| Bob | 25 | LA |

## Lists

### Unordered
- Item one
- Item two
- Item three

### Ordered
1. First
2. Second
3. Third

## Blockquotes

> This is a blockquote.
> It can span multiple lines.

## Images

![Alt text](https://via.placeholder.com/150 "Placeholder")

## Emphasis

This is **bold** and this is *italic* and this is ~~strikethrough~~.

## Horizontal Rule

---

## Links

[Internal link](#tables)

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('tables are rendered correctly', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const table = page.locator('table');
    await expect(table).toBeVisible();

    const rows = page.locator('table tr');
    const count = await rows.count();
    expect(count).toBeGreaterThanOrEqual(3); // Header + 2 data rows
  });

  test('table has proper header', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const th = page.locator('table th');
    const count = await th.count();
    expect(count).toBe(3); // Name, Age, City

    const headerText = await page.textContent('table thead, table tr:first-child');
    expect(headerText).toContain('Name');
    expect(headerText).toContain('Age');
  });

  test('unordered lists render with bullets', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const ul = page.locator('ul');
    const count = await ul.count();
    expect(count).toBeGreaterThan(0);

    const items = page.locator('ul li');
    expect(await items.count()).toBeGreaterThanOrEqual(3);
  });

  test('ordered lists render with numbers', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const ol = page.locator('ol');
    await expect(ol).toBeVisible();

    const items = page.locator('ol li');
    expect(await items.count()).toBeGreaterThanOrEqual(3);
  });

  test('blockquotes are styled', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const blockquote = page.locator('blockquote');
    await expect(blockquote.first()).toBeVisible();

    const text = await blockquote.first().textContent();
    expect(text).toContain('blockquote');
  });

  test('bold and italic text renders correctly', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const bold = page.locator('strong');
    await expect(bold.first()).toBeVisible();
    expect(await bold.first().textContent()).toBe('bold');

    const italic = page.locator('em');
    await expect(italic.first()).toBeVisible();
    expect(await italic.first().textContent()).toBe('italic');
  });

  test('horizontal rule renders', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const hr = page.locator('hr');
    const count = await hr.count();
    expect(count).toBeGreaterThan(0);
  });

  test('headings have proper hierarchy', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const h1 = page.locator('h1');
    const h2 = page.locator('h2');
    const h3 = page.locator('h3');

    expect(await h1.count()).toBeGreaterThan(0);
    expect(await h2.count()).toBeGreaterThan(0);
    expect(await h3.count()).toBeGreaterThan(0);
  });
});
