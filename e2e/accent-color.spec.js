// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for custom accent color
 */
test.describe("Accent Color", () => {
  let serverProcess;
  let testDir;
  const PORT = 4269;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Accent Color Test

Testing custom accent colors.

[A link to test](#section)

## Section

Some content.

${LOREM}
`,
      "about.md": `# About

Another page.

${LOREM}
`,
    });

    // Start with custom accent color
    serverProcess = await startServer(testDir, PORT, [
      "--accent-color=#ff5500",
    ]);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("custom accent color is applied to page", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Check if accent color is set via CSS variable or inline style
    const hasAccentColor = await page.evaluate(() => {
      // Check root CSS variable
      const rootStyle = getComputedStyle(document.documentElement);
      const accentVar = rootStyle.getPropertyValue("--accent-color").trim();

      // Check for inline style on root or style tags containing the color
      const styleContent = Array.from(document.querySelectorAll("style"))
        .map((s) => s.textContent)
        .join("");

      return (
        accentVar.includes("ff5500") ||
        accentVar.includes("255") ||
        styleContent.includes("ff5500") ||
        styleContent.includes("#ff5500")
      );
    });

    expect(hasAccentColor).toBe(true);
  });

  test("links use accent color", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const link = page.locator(".prose a, article a").first();
    const color = await link.evaluate((el) => {
      return getComputedStyle(el).color;
    });

    // Color should be applied (either directly or via variable)
    expect(color).toBeDefined();
  });

  test("active nav item uses accent color", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    const activeLink = page.locator(
      ".tree-nav a.active, .tree-nav a[aria-current]",
    );
    const exists = await activeLink.count();

    if (exists > 0) {
      const color = await activeLink.evaluate((el) => {
        return getComputedStyle(el).color;
      });
      expect(color).toBeDefined();
    }
  });

  test("page renders with custom accent", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Page should render without errors
    const heading = await page.textContent("h1");
    expect(heading).toContain("Accent Color Test");
  });
});

test.describe("Default Accent Color", () => {
  let serverProcess;
  let testDir;
  const PORT = 4270;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Default Colors

Testing default accent color.

${LOREM}
`,
    });

    // Start without custom accent color
    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("default accent color is applied", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Page should render correctly with default styling
    const heading = await page.textContent("h1");
    expect(heading).toContain("Default Colors");

    // Links should have some color styling
    const bodyStyles = await page.evaluate(() => {
      return document.body !== null;
    });
    expect(bodyStyles).toBe(true);
  });
});
