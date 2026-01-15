// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * E2E tests for root index fallback behavior
 *
 * When no index.md exists, should fall back to:
 * 1. README.md
 * 2. First markdown file (alphabetically)
 */
test.describe('Root Index Fallback', () => {
  test.describe('README.md fallback', () => {
    let serverProcess;
    let testDir;
    const PORT = 4255;

    test.beforeAll(async () => {
      // Create site with README.md but no index.md
      testDir = createTestSite({
        'README.md': `# Welcome to My Project

This is the README file serving as the home page.

${LOREM}
`,
        'guide.md': `# Guide

A guide page.

${LOREM}
`
      });

      serverProcess = await startServer(testDir, PORT);
    });

    test.afterAll(async () => {
      stopServer(serverProcess, testDir);
    });

    test('README.md becomes home page when no index.md', async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const heading = await page.textContent('h1');
      expect(heading).toContain('Welcome to My Project');

      const content = await page.textContent('.prose');
      expect(content).toContain('README file serving as the home page');
    });
  });

  test.describe('First file fallback', () => {
    let serverProcess;
    let testDir;
    const PORT = 4256;

    test.beforeAll(async () => {
      // Create site with no index.md or README.md
      testDir = createTestSite({
        '01-introduction.md': `# Introduction

This should become the home page as the first file.

${LOREM}
`,
        '02-setup.md': `# Setup

Setup instructions.

${LOREM}
`,
        '03-usage.md': `# Usage

Usage guide.

${LOREM}
`
      });

      serverProcess = await startServer(testDir, PORT);
    });

    test.afterAll(async () => {
      stopServer(serverProcess, testDir);
    });

    test('first file becomes home page when no index or README', async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      const heading = await page.textContent('h1');
      expect(heading).toContain('Introduction');

      const content = await page.textContent('.prose');
      expect(content).toContain('should become the home page');
    });

    test('first file is removed from nav (since its the home page)', async ({ page }) => {
      await page.goto(`http://localhost:${PORT}/`);

      // The nav should not show Introduction since it's the home page
      const navText = await page.textContent('.tree-nav');

      // Setup and Usage should be in nav
      expect(navText).toContain('Setup');
      expect(navText).toContain('Usage');
    });
  });
});
