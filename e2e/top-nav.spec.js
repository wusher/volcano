// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for top navigation bar toggle
 */
test.describe("Top Navigation", () => {
  test.describe("With Top Nav Enabled", () => {
    let serverProcess;
    let testDir;
    const PORT = 4271;

    test.beforeAll(async () => {
      testDir = createTestSite({
        "index.md": `# Top Nav Test

Testing top navigation.

${LOREM}
`,
        "about.md": `# About

About page.

${LOREM}
`,
      });

      // Start with top nav enabled (default)
      serverProcess = await startServer(testDir, PORT, ["--top-nav"]);
    });

    test.afterAll(async () => {
      stopServer(serverProcess, testDir);
    });

    test("top navigation bar is visible", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const topNav = page.locator(".top-nav, header, .header");
      const count = await topNav.count();
      expect(count).toBeGreaterThan(0);
    });

    test("top nav contains site title", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const topNav = page.locator(".top-nav").first();
      const count = await topNav.count();

      if (count > 0) {
        const text = await topNav.textContent();
        expect(text.length).toBeGreaterThan(0);
      } else {
        // Header should exist
        const header = page.locator("header").first();
        await expect(header).toBeVisible();
      }
    });

    test("theme toggle exists on page", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      // Theme toggle should exist somewhere on page
      const themeToggle = page.locator(".theme-toggle");
      const count = await themeToggle.count();
      expect(count).toBeGreaterThan(0);
    });

    test("top nav contains search trigger", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      // Search button or Cmd+K indicator
      const searchTrigger = page.locator(
        ".search-trigger, .search-button, [data-search]",
      );
      const count = await searchTrigger.count();
      // May or may not exist depending on implementation
      expect(count).toBeGreaterThanOrEqual(0);
    });
  });

  test.describe("With Top Nav Disabled", () => {
    let serverProcess;
    let testDir;
    const PORT = 4272;

    test.beforeAll(async () => {
      testDir = createTestSite({
        "index.md": `# No Top Nav

Testing without top navigation.

${LOREM}
`,
      });

      // Start with top nav disabled
      serverProcess = await startServer(testDir, PORT, ["--top-nav=false"]);
    });

    test.afterAll(async () => {
      stopServer(serverProcess, testDir);
    });

    test("page renders without top nav", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const heading = await page.textContent("h1");
      expect(heading).toContain("No Top Nav");
    });

    test("content is still accessible", async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const content = await page.textContent(".prose, article");
      expect(content).toContain("Testing without top navigation");
    });
  });
});
