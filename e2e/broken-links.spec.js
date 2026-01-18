// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for broken links warning in dev server
 *
 * Verifies that broken internal links show an inline warning banner
 * instead of blocking the page.
 */
test.describe("Broken Links Warning", () => {
  let serverProcess;
  let testDir;
  const PORT = 4254;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Home Page

Welcome to the site.

${LOREM}
`,
      "page-with-broken-link.md": `# Page With Broken Link

This page has a [[nonexistent-page|broken wiki link]].

It also has a [broken markdown link](/does-not-exist/).

${LOREM}
`,
      "valid-page.md": `# Valid Page

This page has valid links.

See the [home page](/).

${LOREM}
`,
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("page with broken links shows warning banner", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/page-with-broken-link/`);

    // Page should still render (not blocked)
    const heading = await page.textContent("h1");
    expect(heading).toContain("Page With Broken Link");

    // Warning banner should be visible
    const warningText = await page.textContent("body");
    expect(warningText).toContain("Broken Link");
  });

  test("warning banner shows link details", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/page-with-broken-link/`);

    // Should show the broken URLs
    const pageContent = await page.textContent("body");
    expect(pageContent).toMatch(/nonexistent-page|does-not-exist/);
  });

  test("page content is still visible below warning", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/page-with-broken-link/`);

    // The main content should be present
    const content = await page.textContent(".prose");
    expect(content).toContain("This page has a");
    expect(content).toContain("broken wiki link");
  });

  test("page without broken links has no warning", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/valid-page/`);

    // Page should render without warning
    const heading = await page.textContent("h1");
    expect(heading).toContain("Valid Page");

    // No warning banner - check the prose content area
    const proseContent = await page.textContent(".prose");
    expect(proseContent).not.toContain("Broken Link");
    expect(proseContent).not.toContain("⚠️");
  });

  test("home page renders normally without warnings", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const heading = await page.textContent("h1");
    expect(heading).toContain("Home Page");

    // No warning banner - check the prose content area
    const proseContent = await page.textContent(".prose");
    expect(proseContent).not.toContain("Broken Link");
    expect(proseContent).not.toContain("⚠️");
  });
});
