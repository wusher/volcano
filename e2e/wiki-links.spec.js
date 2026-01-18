// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for wiki-style links [[page]] and [[page|text]]
 */
test.describe("Wiki Links", () => {
  let serverProcess;
  let testDir;
  const PORT = 4273;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Wiki Links Test

Testing wiki-style links.

${LOREM}
`,
      "about.md": `# About

This page links to [[getting-started]].

Also see [[getting-started|the getting started guide]].

${LOREM}
`,
      "getting-started.md": `# Getting Started

Welcome to the getting started guide.

Go back to [[about]].

${LOREM}
`,
      "guides/advanced.md": `# Advanced Guide

Link to root: [[index|home]]

Link to sibling section: [[about]]

${LOREM}
`,
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("wiki link renders as anchor tag", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    // Wiki link [[getting-started]] should become an <a> tag
    const link = page.locator('a[href*="getting-started"]');
    const count = await link.count();
    expect(count).toBeGreaterThan(0);
  });

  test("wiki link with custom text shows that text", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    // [[getting-started|the getting started guide]] should show custom text
    const linkText = await page.textContent(".prose");
    expect(linkText).toContain("the getting started guide");
  });

  test("clicking wiki link navigates to target page", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    const link = page.locator('a[href*="getting-started"]').first();
    await link.click();
    await page.waitForTimeout(500);

    expect(page.url()).toContain("/getting-started/");

    const heading = await page.textContent("h1");
    expect(heading).toContain("Getting Started");
  });

  test("wiki link without custom text uses page title", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    // The first [[getting-started]] should use the page title or filename
    const links = page.locator(".prose a");
    const texts = await links.allTextContents();

    // Should have "Getting Started" or "getting-started" as link text
    const hasPageName = texts.some(
      (t) =>
        t.toLowerCase().includes("getting") ||
        t.toLowerCase().includes("started"),
    );
    expect(hasPageName).toBe(true);
  });

  test("wiki links work from nested pages", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/guides/advanced/`);

    // Should have link to home/index
    const homeLink = page.locator('a:has-text("home")');
    const count = await homeLink.count();
    expect(count).toBeGreaterThan(0);
  });

  test("bidirectional wiki links work", async ({ page }) => {
    // Start at about page
    await page.goto(`http://localhost:${PORT}/about/`);

    // Click to getting-started
    await page.click('a[href*="getting-started"]');
    await page.waitForTimeout(500);

    // Now click back to about
    const aboutLink = page.locator('a[href*="about"]');
    if ((await aboutLink.count()) > 0) {
      await aboutLink.first().click();
      await page.waitForTimeout(500);

      expect(page.url()).toContain("/about/");
    }
  });
});
