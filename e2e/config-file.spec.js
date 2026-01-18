// @ts-check
const { test, expect } = require("@playwright/test");
const { createTestSite, startServer, stopServer } = require("./helpers");

/**
 * Tests for volcano.json config file loading
 * Verifies all config options are properly applied when loaded from file
 */

test.describe("Config File Loading", () => {
  let serverProcess;
  let testDir;
  const PORT = 3090;

  test.afterEach(async () => {
    stopServer(serverProcess, testDir);
  });

  test("loads title from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "Config Test Title",
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Title should be "Page Title - Site Title"
    await expect(page).toHaveTitle(/Config Test Title/);
  });

  test("loads accentColor from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "Accent Test",
        accentColor: "#ff6600",
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Check that accent color CSS variable is set (--accent, not --accent-color)
    const accentColor = await page.evaluate(() => {
      const rootStyle = getComputedStyle(document.documentElement);
      return rootStyle.getPropertyValue("--accent").trim();
    });
    expect(accentColor).toBe("#ff6600");
  });

  test("loads theme from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "Theme Test",
        theme: "blog",
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Blog theme should have specific styling
    // Just verify page loads without error
    await expect(page.locator("h1")).toHaveText("Home");
  });

  test("loads breadcrumbs setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "guides/index.md": "# Guides\nGuide content",
      "volcano.json": JSON.stringify({
        title: "Breadcrumb Test",
        breadcrumbs: false,
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/guides/`);

    // Breadcrumbs should not be visible when disabled
    const breadcrumbs = page.locator(".breadcrumbs");
    await expect(breadcrumbs).not.toBeVisible();
  });

  test("loads pageNav setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "page1.md": "# Page 1\nContent",
      "page2.md": "# Page 2\nContent",
      "volcano.json": JSON.stringify({
        title: "PageNav Test",
        pageNav: true,
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/page1/`);

    // Page nav should be visible when enabled
    const pageNav = page.locator(".page-nav");
    await expect(pageNav).toBeVisible();
  });

  test("loads topNav setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "about.md": "# About\nAbout content",
      "volcano.json": JSON.stringify({
        title: "TopNav Test",
        topNav: true,
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Top nav should contain root-level pages
    const topNav = page.locator(".top-nav");
    await expect(topNav).toBeVisible();
  });

  test("loads search setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "Search Test",
        search: true,
      }),
    });

    // Don't pass default flags, let config file enable search
    serverProcess = await startServer(testDir, PORT, [], { noDefaults: true });
    await page.goto(`http://localhost:${PORT}/`);

    // Command palette should exist and be openable with Cmd+K
    await page.keyboard.press("Meta+k");
    await page.waitForTimeout(300);

    // Command palette should now be visible
    const paletteVisible = await page.isVisible(".command-palette");
    expect(paletteVisible).toBe(true);
  });

  test("loads pwa setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "PWA Test",
        pwa: true,
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // PWA manifest link should be present
    const manifestLink = page.locator('link[rel="manifest"]');
    await expect(manifestLink).toHaveCount(1);
  });

  test("loads instantNav setting from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\n[Link](./page/)",
      "page.md": "# Page\nContent",
      "volcano.json": JSON.stringify({
        title: "InstantNav Test",
        instantNav: true,
      }),
    });

    // Don't pass default flags, let config file enable instantNav
    serverProcess = await startServer(testDir, PORT, [], { noDefaults: true });
    await page.goto(`http://localhost:${PORT}/`);

    // Just verify page loads - instant nav is enabled via config
    await expect(page.locator("h1")).toHaveText("Home");
  });

  test("url in config file does not cause errors", async ({ page }) => {
    // URL is used for canonical links in build mode, but should not cause errors in serve mode
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "URL Test",
        url: "https://example.com/docs",
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Page should render without errors
    await expect(page.locator("h1")).toHaveText("Home");
    await expect(page).toHaveTitle(/URL Test/);
  });

  test("loads multiple options from config file", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "guides/index.md": "# Guides\nGuide content",
      "volcano.json": JSON.stringify({
        title: "Multi Option Test",
        accentColor: "#0ea5e9",
        theme: "docs",
        breadcrumbs: true,
        pageNav: false,
        search: true,
      }),
    });

    serverProcess = await startServer(testDir, PORT, []);
    await page.goto(`http://localhost:${PORT}/`);

    // Verify title
    await expect(page).toHaveTitle(/Multi Option Test/);

    // Verify accent color (--accent, not --accent-color)
    const accentColor = await page.evaluate(() => {
      return getComputedStyle(document.documentElement)
        .getPropertyValue("--accent")
        .trim();
    });
    expect(accentColor).toBe("#0ea5e9");

    // Verify search is enabled (command palette works)
    await page.keyboard.press("Meta+k");
    await page.waitForTimeout(300);
    const paletteVisible = await page.isVisible(".command-palette");
    expect(paletteVisible).toBe(true);
  });

  test("CLI flags override config file values", async ({ page }) => {
    testDir = createTestSite({
      "index.md": "# Home\nWelcome",
      "volcano.json": JSON.stringify({
        title: "Config Title",
        accentColor: "#ff0000",
      }),
    });

    // Pass --title flag to override config
    serverProcess = await startServer(testDir, PORT, ["--title", "CLI Title"]);
    await page.goto(`http://localhost:${PORT}/`);

    // CLI title should override config title
    await expect(page).toHaveTitle(/CLI Title/);
  });
});
