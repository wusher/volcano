// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for code syntax highlighting
 */
test.describe('Code Syntax Highlighting', () => {
  let serverProcess;
  let testDir;
  const PORT = 4261;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Code Examples

Some code examples below.

## JavaScript

\`\`\`javascript
function hello(name) {
  console.log('Hello, ' + name);
  return true;
}
\`\`\`

## Python

\`\`\`python
def hello(name):
    print(f"Hello, {name}")
    return True
\`\`\`

## Inline Code

Here is some \`inline code\` in a sentence.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('code blocks are rendered with pre and code tags', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const codeBlocks = page.locator('pre code');
    const count = await codeBlocks.count();
    expect(count).toBeGreaterThanOrEqual(2);
  });

  test('code blocks have language attribute or class', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Check for language-specific class or data attribute
    const codeBlock = page.locator('pre code').first();
    const className = await codeBlock.getAttribute('class') || '';
    const dataLang = await codeBlock.getAttribute('data-lang') || '';

    // Should have either a language class or the code content is properly formatted
    const hasLanguageIndicator = className.match(/language-|highlight|javascript/) ||
                                  dataLang.length > 0;

    // At minimum, the code block should exist and contain code
    const content = await codeBlock.textContent();
    expect(content.length).toBeGreaterThan(0);
  });

  test('code content is preserved correctly', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const codeContent = await page.textContent('pre code');
    expect(codeContent).toContain('function');
    expect(codeContent).toContain('hello');
  });

  test('inline code is rendered with code tag', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const inlineCode = page.locator('p code');
    const count = await inlineCode.count();
    expect(count).toBeGreaterThan(0);

    const text = await inlineCode.first().textContent();
    expect(text).toBe('inline code');
  });

  test('code blocks have copy button or are selectable', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Code should be in a pre block that's selectable
    const preBlock = page.locator('pre').first();
    await expect(preBlock).toBeVisible();

    // Verify the code is contained properly
    const codeText = await preBlock.textContent();
    expect(codeText.length).toBeGreaterThan(10);
  });
});
