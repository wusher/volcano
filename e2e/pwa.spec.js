// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for PWA functionality (service worker, manifest, offline)
 */
test.describe('PWA Support', () => {
  let serverProcess;
  let testDir;
  const PORT = 4260;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# PWA Test Site

Welcome to the PWA test.

${LOREM}
`,
      'about.md': `# About

About page for offline testing.

${LOREM}
`
    });

    // Start server with PWA enabled
    serverProcess = await startServer(testDir, PORT, ['--pwa']);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('manifest.json is served', async ({ page }) => {
    const response = await page.goto(`http://localhost:${PORT}/manifest.json`);
    expect(response.status()).toBe(200);

    const manifest = await response.json();
    expect(manifest).toHaveProperty('name');
    expect(manifest).toHaveProperty('start_url');
    expect(manifest).toHaveProperty('display');
  });

  test('service worker is registered', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Wait for service worker registration
    const swRegistered = await page.evaluate(async () => {
      if (!('serviceWorker' in navigator)) return false;
      const registration = await navigator.serviceWorker.getRegistration();
      return !!registration;
    });

    expect(swRegistered).toBe(true);
  });

  test('manifest link is in HTML head', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    const manifestLink = page.locator('link[rel="manifest"]');
    await expect(manifestLink).toHaveAttribute('href', '/manifest.json');
  });

  test('PWA meta tags are present', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // PWA should have viewport meta
    const viewport = page.locator('meta[name="viewport"]');
    await expect(viewport).toHaveCount(1);

    // Page should be mobile-friendly
    const heading = await page.textContent('h1');
    expect(heading).toContain('PWA Test');
  });

  test('manifest link exists in head', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Manifest link should be present when PWA is enabled
    const manifestLink = page.locator('link[rel="manifest"]');
    const count = await manifestLink.count();
    expect(count).toBe(1);
  });
});
