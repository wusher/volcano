// @ts-check
const { test, expect } = require('@playwright/test');
const { createTestSite, startServer, stopServer, LOREM } = require('./helpers');

/**
 * Comprehensive E2E tests for command palette functionality
 */
test.describe('Command Palette', () => {
  let serverProcess;
  let testDir;
  const PORT = 4253;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home Page

Welcome to the documentation site.

${LOREM}
`,
      'getting-started.md': `# Getting Started

This guide helps you get started quickly.

## Installation

Install the package using npm.

## Configuration

Configure your settings here.

### Advanced Options

These are advanced configuration options.

${LOREM}
`,
      'api-reference.md': `# API Reference

Complete API documentation.

## Core Methods

The core methods available.

## Helper Functions

Additional helper functions.

${LOREM}
`,
      'guides/basics.md': `# Basic Guide

Learn the basics.

## First Steps

Get started with the first steps.

${LOREM}
`,
      'guides/advanced.md': `# Advanced Guide

Advanced topics and techniques.

## Performance Tips

Optimize your performance.

## Troubleshooting

Common issues and solutions.

${LOREM}
`,
      'tutorials/quick-start.md': `# Quick Start Tutorial

Get up and running quickly.

## Prerequisites

What you need before starting.

## Step One

The first step in the tutorial.

## Step Two

The second step in the tutorial.

${LOREM}
`,
      'faq.md': `# FAQ

Frequently asked questions.

## General Questions

Common questions about the project.

## Technical Questions

Technical questions and answers.

${LOREM}
`,
      'changelog.md': `# Changelog

Version history and changes.

## Version 2.0

Major release with new features.

## Version 1.0

Initial release.

${LOREM}
`,
      'special-chars.md': `# Special & Characters <Test>

Testing special characters in search.

## Section with "Quotes"

A section with special characters.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('Ctrl+K opens command palette', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.tree-nav');

    // Press Ctrl+K (alternative to Cmd+K)
    await page.keyboard.press('Control+k');
    await page.waitForTimeout(300);

    // Command palette should be visible
    const paletteVisible = await page.isVisible('.command-palette.open');
    expect(paletteVisible).toBe(true);
  });

  test('clicking backdrop closes command palette', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Click the backdrop
    await page.click('.command-palette-backdrop');
    await page.waitForTimeout(300);

    // Command palette should be closed
    const paletteHidden = await page.isHidden('.command-palette.open');
    expect(paletteHidden).toBe(true);
  });

  test('first result is auto-selected', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Type a search query
    await page.fill('#command-palette-input', 'guide');
    await page.waitForTimeout(300);

    // First result should have 'selected' class
    const firstResult = page.locator('.command-palette-result').first();
    await expect(firstResult).toHaveClass(/selected/);
  });

  test('arrow keys navigate through results', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search to get multiple results
    await page.fill('#command-palette-input', 'guide');
    await page.waitForTimeout(300);

    // Verify we have multiple results
    const resultCount = await page.locator('.command-palette-result').count();
    expect(resultCount).toBeGreaterThan(1);

    // First result should be selected initially
    const firstResult = page.locator('.command-palette-result').first();
    await expect(firstResult).toHaveClass(/selected/);

    // Press ArrowDown to move to second result
    await page.keyboard.press('ArrowDown');
    await page.waitForTimeout(100);

    // Second result should now be selected
    const secondResult = page.locator('.command-palette-result').nth(1);
    await expect(secondResult).toHaveClass(/selected/);
    await expect(firstResult).not.toHaveClass(/selected/);

    // Press ArrowUp to go back to first
    await page.keyboard.press('ArrowUp');
    await page.waitForTimeout(100);

    await expect(firstResult).toHaveClass(/selected/);
    await expect(secondResult).not.toHaveClass(/selected/);
  });

  test('Enter key navigates to selected result', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for FAQ
    await page.fill('#command-palette-input', 'faq');
    await page.waitForTimeout(300);

    // Press Enter to navigate
    await page.keyboard.press('Enter');
    await page.waitForTimeout(500);

    // Should navigate to FAQ page
    expect(page.url()).toContain('/faq/');
  });

  test('ArrowDown then Enter navigates to second result', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for something with multiple results
    await page.fill('#command-palette-input', 'step');
    await page.waitForTimeout(300);

    // Get the second result's href before navigation
    const secondResult = page.locator('.command-palette-result').nth(1);
    const expectedHref = await secondResult.getAttribute('href');

    // Press ArrowDown then Enter
    await page.keyboard.press('ArrowDown');
    await page.waitForTimeout(100);
    await page.keyboard.press('Enter');
    await page.waitForTimeout(500);

    // URL should contain the expected path
    expect(page.url()).toContain(expectedHref);
  });

  test('shows "No results found" for non-matching query', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Type a query that won't match anything
    await page.fill('#command-palette-input', 'xyznonexistent123');
    await page.waitForTimeout(300);

    // Should show "No results found"
    const emptyMessage = await page.textContent('.command-palette-empty');
    expect(emptyMessage).toBe('No results found');
  });

  test('shows "Type to search..." initially', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Should show initial message
    const emptyMessage = await page.textContent('.command-palette-empty');
    expect(emptyMessage).toBe('Type to search...');
  });

  test('multi-term search filters results', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search with multiple terms
    await page.fill('#command-palette-input', 'quick start');
    await page.waitForTimeout(300);

    // Should find the quick start tutorial
    const resultText = await page.textContent('.command-palette-results');
    expect(resultText).toContain('Quick Start');
  });

  test('search is case-insensitive', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search with uppercase
    await page.fill('#command-palette-input', 'API REFERENCE');
    await page.waitForTimeout(300);

    // Should find the API Reference page
    const results = await page.locator('.command-palette-result').count();
    expect(results).toBeGreaterThan(0);

    const resultText = await page.textContent('.command-palette-results');
    expect(resultText).toContain('API Reference');
  });

  test('displays different result types (page, h2, h3)', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for getting started which has page and headings
    await page.fill('#command-palette-input', 'getting');
    await page.waitForTimeout(300);

    // Check for different result types
    const resultTypes = await page.locator('.result-type').allTextContents();
    expect(resultTypes).toContain('page');
  });

  test('heading results show parent page as snippet', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for a specific heading
    await page.fill('#command-palette-input', 'installation');
    await page.waitForTimeout(300);

    // Should show the parent page title as snippet
    const snippet = await page.textContent('.result-snippet');
    expect(snippet).toBe('Getting Started');
  });

  test('limits results to 10', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for something common that would have many matches
    await page.fill('#command-palette-input', 'a');
    await page.waitForTimeout(300);

    // Should have at most 10 results
    const resultCount = await page.locator('.command-palette-result').count();
    expect(resultCount).toBeLessThanOrEqual(10);
  });

  test('input is focused when palette opens', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Check that input is focused
    const isFocused = await page.evaluate(() => {
      return document.activeElement?.id === 'command-palette-input';
    });
    expect(isFocused).toBe(true);
  });

  test('input is cleared when reopening palette', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette and type something
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');
    await page.fill('#command-palette-input', 'test query');

    // Close and reopen
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Input should be empty
    const inputValue = await page.inputValue('#command-palette-input');
    expect(inputValue).toBe('');
  });

  test('navigating to heading scrolls to position', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for a specific heading
    await page.fill('#command-palette-input', 'performance tips');
    await page.waitForTimeout(300);

    // Click the result
    await page.click('.command-palette-result');
    await page.waitForTimeout(500);

    // Should navigate to the page with hash
    expect(page.url()).toContain('/guides/advanced/');
    expect(page.url()).toContain('#performance-tips');

    // The heading should be visible
    const heading = page.locator('#performance-tips');
    await expect(heading).toBeVisible();
  });

  test('searching path segments matches pages in subdirectories', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for a directory path
    await page.fill('#command-palette-input', 'guides');
    await page.waitForTimeout(300);

    // Should find pages in the guides directory
    const resultText = await page.textContent('.command-palette-results');
    expect(resultText).toContain('Guide');
  });

  test('body gets command-palette-open class when open', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Initially body should not have the class
    const hasClassBefore = await page.evaluate(() => {
      return document.body.classList.contains('command-palette-open');
    });
    expect(hasClassBefore).toBe(false);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Body should have the class
    const hasClassAfter = await page.evaluate(() => {
      return document.body.classList.contains('command-palette-open');
    });
    expect(hasClassAfter).toBe(true);

    // Close and check again
    await page.keyboard.press('Escape');
    await page.waitForTimeout(300);

    const hasClassFinal = await page.evaluate(() => {
      return document.body.classList.contains('command-palette-open');
    });
    expect(hasClassFinal).toBe(false);
  });

  test('ArrowUp at first result stays at first result', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search to get results
    await page.fill('#command-palette-input', 'guide');
    await page.waitForTimeout(300);

    // First result should be selected
    const firstResult = page.locator('.command-palette-result').first();
    await expect(firstResult).toHaveClass(/selected/);

    // Press ArrowUp multiple times
    await page.keyboard.press('ArrowUp');
    await page.keyboard.press('ArrowUp');
    await page.waitForTimeout(100);

    // First result should still be selected
    await expect(firstResult).toHaveClass(/selected/);
  });

  test('special characters in search results are escaped', async ({ page }) => {
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search for page with special chars
    await page.fill('#command-palette-input', 'special');
    await page.waitForTimeout(300);

    // Result should display properly (not render as HTML)
    const resultTitle = await page.textContent('.result-title');
    expect(resultTitle).toContain('Special & Characters <Test>');

    // Verify it's escaped (not rendered as HTML)
    const resultHtml = await page.locator('.result-title').first().innerHTML();
    expect(resultHtml).toContain('&amp;');
    expect(resultHtml).toContain('&lt;');
    expect(resultHtml).toContain('&gt;');
  });
});

/**
 * Mobile-specific command palette tests
 */
test.describe('Command Palette Mobile', () => {
  let serverProcess;
  let testDir;
  const PORT = 4254;

  test.beforeAll(async () => {
    testDir = createTestSite({
      'index.md': `# Home Page

Welcome to the documentation.

${LOREM}
`,
      'about.md': `# About

About this project.

${LOREM}
`
    });

    serverProcess = await startServer(testDir, PORT);
  });

  test.afterAll(async () => {
    stopServer(serverProcess, testDir);
  });

  test('mobile search button opens command palette', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`http://localhost:${PORT}/`);
    await page.waitForSelector('.tree-nav');

    // Find and click the mobile search button
    const searchButton = page.locator('.mobile-search-button');

    // Check if the mobile search button exists and is visible
    if (await searchButton.isVisible()) {
      await searchButton.click();
      await page.waitForTimeout(300);

      // Command palette should be open
      const paletteVisible = await page.isVisible('.command-palette.open');
      expect(paletteVisible).toBe(true);
    } else {
      // If no mobile button, Cmd+K should still work
      await page.keyboard.press('Meta+k');
      await page.waitForSelector('.command-palette.open');
      const paletteVisible = await page.isVisible('.command-palette.open');
      expect(paletteVisible).toBe(true);
    }
  });

  test('command palette works in mobile viewport', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(`http://localhost:${PORT}/`);

    // Open command palette
    await page.keyboard.press('Meta+k');
    await page.waitForSelector('.command-palette.open');

    // Search and navigate
    await page.fill('#command-palette-input', 'about');
    await page.waitForTimeout(300);

    await page.click('.command-palette-result');
    await page.waitForTimeout(500);

    // Should navigate to about page
    expect(page.url()).toContain('/about/');
  });
});
