// @ts-check
const { test, expect } = require("@playwright/test");
const fs = require("fs");
const path = require("path");

/**
 * Extract CSS selectors from a CSS file
 * @param {string} cssContent - The CSS file content
 * @returns {Set<string>} - Set of unique selectors
 */
function extractSelectors(cssContent) {
  const selectors = new Set();

  // Remove comments (both single-line and multi-line)
  let css = cssContent.replace(/\/\*[\s\S]*?\*\//g, "");

  // Remove @keyframes blocks entirely (they contain 'from', 'to', percentages that aren't selectors)
  css = css.replace(
    /@keyframes\s+[\w-]+\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}/g,
    "",
  );

  // Remove @media/@supports wrapper but keep content (simplified approach)
  // We'll process everything and filter non-selectors

  // Match selectors: everything before a { that's not inside a @rule
  // This regex captures selector blocks
  const selectorRegex = /([^{}@]+)\{[^{}]*\}/g;

  let match;
  while ((match = selectorRegex.exec(css)) !== null) {
    const selectorBlock = match[1].trim();

    // Skip @rules, empty strings, and @keyframes percentages/keywords
    if (selectorBlock.startsWith("@") || selectorBlock === "") {
      continue;
    }

    // Skip @keyframes keywords and percentages
    if (/^(from|to|\d+%)$/.test(selectorBlock)) {
      continue;
    }

    // Split by comma for grouped selectors and normalize
    const individualSelectors = selectorBlock
      .split(",")
      .map((s) => s.trim())
      .filter((s) => s);

    for (const selector of individualSelectors) {
      // Normalize whitespace
      const normalized = selector.replace(/\s+/g, " ").trim();

      // Skip if it's a @keyframes keyword, percentage, or starts with @
      if (
        !normalized ||
        normalized.startsWith("@") ||
        /^(from|to|\d+%)$/.test(normalized)
      ) {
        continue;
      }

      selectors.add(normalized);
    }
  }

  return selectors;
}

/**
 * E2E test for CSS selector coverage
 * Ensures vanilla.css contains all selectors from docs.css and blog.css
 */
test.describe("CSS Selector Coverage", () => {
  const themesDir = path.join(__dirname, "..", "internal", "styles", "themes");

  test("vanilla.css contains all selectors from docs.css and blog.css", async () => {
    // Read CSS files
    const docsCSS = fs.readFileSync(path.join(themesDir, "docs.css"), "utf-8");
    const blogCSS = fs.readFileSync(path.join(themesDir, "blog.css"), "utf-8");
    const vanillaCSS = fs.readFileSync(
      path.join(themesDir, "vanilla.css"),
      "utf-8",
    );

    // Extract selectors from each theme
    const docsSelectors = extractSelectors(docsCSS);
    const blogSelectors = extractSelectors(blogCSS);
    const vanillaSelectors = extractSelectors(vanillaCSS);

    // Create union of docs and blog selectors
    const docsBlogUnion = new Set([...docsSelectors, ...blogSelectors]);

    // Find selectors missing from vanilla
    const missingSelectors = [];
    for (const selector of docsBlogUnion) {
      if (!vanillaSelectors.has(selector)) {
        missingSelectors.push(selector);
      }
    }

    // Report results
    if (missingSelectors.length > 0) {
      console.log("\nSelectors missing from vanilla.css:");
      missingSelectors.sort().forEach((s) => console.log(`  - ${s}`));
      console.log(`\nTotal: ${missingSelectors.length} missing selectors`);
    }

    expect(
      missingSelectors,
      `Missing selectors in vanilla.css:\n${missingSelectors.join("\n")}`,
    ).toHaveLength(0);
  });

  test("extracts correct number of selectors from each theme", async () => {
    // Read CSS files
    const docsCSS = fs.readFileSync(path.join(themesDir, "docs.css"), "utf-8");
    const blogCSS = fs.readFileSync(path.join(themesDir, "blog.css"), "utf-8");
    const vanillaCSS = fs.readFileSync(
      path.join(themesDir, "vanilla.css"),
      "utf-8",
    );

    // Extract selectors
    const docsSelectors = extractSelectors(docsCSS);
    const blogSelectors = extractSelectors(blogCSS);
    const vanillaSelectors = extractSelectors(vanillaCSS);

    console.log(`\nSelector counts:`);
    console.log(`  docs.css: ${docsSelectors.size} selectors`);
    console.log(`  blog.css: ${blogSelectors.size} selectors`);
    console.log(`  vanilla.css: ${vanillaSelectors.size} selectors`);

    // Each theme should have a reasonable number of selectors
    expect(docsSelectors.size).toBeGreaterThan(50);
    expect(blogSelectors.size).toBeGreaterThan(50);
    expect(vanillaSelectors.size).toBeGreaterThan(50);
  });
});
