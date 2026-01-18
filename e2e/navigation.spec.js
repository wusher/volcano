// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for site navigation (tree nav, folders, active states)
 */
test.describe("Site Navigation", () => {
  let serverProcess;
  let testDir;
  const PORT = 4257;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Home

Welcome to the site.

${LOREM}
`,
      "about.md": `# About

About page.

${LOREM}
`,
      "guides/index.md": `# Guides

Guide section index.

${LOREM}
`,
      "guides/getting-started.md": `# Getting Started

A getting started guide.

${LOREM}
`,
      "guides/advanced.md": `# Advanced Guide

Advanced topics.

${LOREM}
`,
      "api/endpoints.md": `# API Endpoints

API endpoint docs.

${LOREM}
`,
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("navigation shows all pages", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const navText = await page.textContent(".tree-nav");
    expect(navText).toContain("About");
    expect(navText).toContain("Guides");
  });

  test("current page is highlighted in nav", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/about/`);

    // The about link should have active class
    const aboutLink = page.locator('.tree-nav a[href="/about/"]');
    await expect(aboutLink).toHaveClass(/active/);
  });

  test("folders can be expanded to show children", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Guides folder should exist
    const guidesFolder = page.locator(
      '.tree-nav .folder-header:has-text("Guides")',
    );
    await expect(guidesFolder).toBeVisible();

    // Click to expand
    await guidesFolder.click();
    await page.waitForTimeout(200);

    // Children should be visible
    const navText = await page.textContent(".tree-nav");
    expect(navText).toContain("Getting Started");
    expect(navText).toContain("Advanced");
  });

  test("clicking folder with index navigates to it", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Click on Guides folder link
    await page.click('.tree-nav a[href="/guides/"]');
    await page.waitForTimeout(500);

    expect(page.url()).toContain("/guides/");

    const heading = await page.textContent("h1");
    expect(heading).toContain("Guides");
  });

  test("nested page navigation works", async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // First expand the Guides folder by clicking its header
    const guidesFolder = page.locator(
      '.tree-nav .folder-header:has-text("Guides")',
    );
    await guidesFolder.click();
    await page.waitForTimeout(300);

    // Wait for child link to become visible
    const gettingStartedLink = page.locator(
      '.tree-nav a[href="/guides/getting-started/"]',
    );
    await expect(gettingStartedLink).toBeVisible({ timeout: 5000 });

    // Navigate to getting started
    await gettingStartedLink.click();
    await page.waitForTimeout(500);

    expect(page.url()).toContain("/guides/getting-started/");

    const heading = await page.textContent("h1");
    expect(heading).toContain("Getting Started");
  });
});
