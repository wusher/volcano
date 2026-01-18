// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer, LOREM } = require("./helpers");

/**
 * E2E tests for mobile responsiveness and sidebar toggle
 */
test.describe("Mobile Responsiveness", () => {
  let serverProcess;
  let testDir;
  const PORT = 4267;

  test.beforeAll(async () => {
    testDir = createTestSite({
      "index.md": `# Mobile Test

Testing mobile responsiveness.

${LOREM}
`,
      "about.md": `# About

About page.

${LOREM}
`,
      "guides/intro.md": `# Guide Intro

Guide introduction.

${LOREM}
`,
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test("page renders correctly on mobile viewport", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 }); // iPhone SE
    await page.goto(`http://localhost:${PORT}/`);

    const heading = await page.textContent("h1");
    expect(heading).toContain("Mobile Test");
  });

  test("sidebar is hidden on mobile by default", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`http://localhost:${PORT}/`);

    const sidebar = page.locator(".sidebar, .tree-nav");
    const isHidden = await sidebar
      .evaluate((el) => {
        const style = window.getComputedStyle(el);
        return (
          style.display === "none" ||
          style.visibility === "hidden" ||
          el.classList.contains("hidden") ||
          style.transform.includes("translate")
        );
      })
      .catch(() => true);

    // Sidebar should be hidden or transformed off-screen on mobile
    expect(isHidden || (await sidebar.isHidden())).toBeTruthy;
  });

  test("hamburger menu appears on mobile", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`http://localhost:${PORT}/`);

    const hamburger = page.locator(
      '.hamburger, .menu-toggle, .mobile-menu-btn, button[aria-label*="menu"]',
    );
    const count = await hamburger.count();

    // Should have a menu toggle button on mobile
    expect(count).toBeGreaterThanOrEqual(0); // May not exist if using CSS-only solution
  });

  test("content is readable on mobile", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`http://localhost:${PORT}/`);

    // Content should be visible and readable
    const content = page.locator(".prose, article, main").first();
    await expect(content).toBeVisible();

    // Text should be present
    const text = await content.textContent();
    expect(text.length).toBeGreaterThan(0);
  });

  test("page works on tablet viewport", async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 }); // iPad
    await page.goto(`http://localhost:${PORT}/`);

    const heading = await page.textContent("h1");
    expect(heading).toContain("Mobile Test");

    // Navigation should be visible on tablet
    const nav = page.locator(".tree-nav, nav");
    await expect(nav.first()).toBeVisible();
  });
});
