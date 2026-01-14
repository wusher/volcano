package pwa

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ServiceWorkerConfig holds configuration for service worker generation.
type ServiceWorkerConfig struct {
	BaseURL   string   // Base URL path prefix
	PageURLs  []string // All page URLs to precache (from tree.AllPages)
	AssetURLs []string // CSS, JS, icon URLs
}

// GenerateServiceWorker creates and writes the sw.js file.
func GenerateServiceWorker(outputDir string, config ServiceWorkerConfig) error {
	// Collect all URLs to cache
	allURLs := make([]string, 0, len(config.PageURLs)+len(config.AssetURLs))
	allURLs = append(allURLs, config.PageURLs...)
	allURLs = append(allURLs, config.AssetURLs...)

	// Sort for consistent hashing
	sort.Strings(allURLs)

	// Generate cache version from hash of all URLs
	cacheVersion := generateCacheVersion(allURLs)
	cacheName := "volcano-cache-" + cacheVersion

	// Build the service worker JavaScript
	sw := buildServiceWorkerJS(cacheName, allURLs)

	// Write to file
	return os.WriteFile(filepath.Join(outputDir, "sw.js"), []byte(sw), 0644)
}

// generateCacheVersion creates a short hash from all cached URLs.
func generateCacheVersion(urls []string) string {
	h := sha256.New()
	for _, url := range urls {
		h.Write([]byte(url))
	}
	return hex.EncodeToString(h.Sum(nil))[:8]
}

// buildServiceWorkerJS generates the service worker JavaScript code.
func buildServiceWorkerJS(cacheName string, urls []string) string {
	// Build URL array string
	var urlsJS strings.Builder
	urlsJS.WriteString("[\n")
	for i, url := range urls {
		urlsJS.WriteString(fmt.Sprintf("  %q", url))
		if i < len(urls)-1 {
			urlsJS.WriteString(",")
		}
		urlsJS.WriteString("\n")
	}
	urlsJS.WriteString("]")

	return fmt.Sprintf(`// Service Worker for offline support
// Cache version: %s

const CACHE_NAME = %q;
const URLS_TO_CACHE = %s;

// Install event - precache all pages and assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(URLS_TO_CACHE))
      .then(() => self.skipWaiting())
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name.startsWith('volcano-cache-') && name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      );
    }).then(() => self.clients.claim())
  );
});

// Fetch event - serve from cache first, fall back to network
self.addEventListener('fetch', (event) => {
  // Only handle GET requests
  if (event.request.method !== 'GET') {
    return;
  }

  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        // Return cached response if found
        if (response) {
          return response;
        }

        // Otherwise fetch from network
        return fetch(event.request).then((response) => {
          // Don't cache non-successful responses
          if (!response || response.status !== 200 || response.type !== 'basic') {
            return response;
          }

          // Clone the response for caching
          const responseToCache = response.clone();

          caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, responseToCache);
          });

          return response;
        });
      })
  );
});
`, cacheName, cacheName, urlsJS.String())
}

// GetServiceWorkerRegistration returns the JS code to register the service worker.
func GetServiceWorkerRegistration(baseURL string) string {
	swPath := "/sw.js"
	if baseURL != "" {
		swPath = baseURL + "/sw.js"
	}
	return fmt.Sprintf(`if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('%s');
}`, swPath)
}

// BuildServiceWorker creates sw.js content in memory and returns it as a string.
func BuildServiceWorker(config ServiceWorkerConfig) string {
	// Collect all URLs to cache
	allURLs := make([]string, 0, len(config.PageURLs)+len(config.AssetURLs))
	allURLs = append(allURLs, config.PageURLs...)
	allURLs = append(allURLs, config.AssetURLs...)

	// Sort for consistent hashing
	sort.Strings(allURLs)

	// Generate cache version from hash of all URLs
	cacheVersion := generateCacheVersion(allURLs)
	cacheName := "volcano-cache-" + cacheVersion

	// Build and return the service worker JavaScript
	return buildServiceWorkerJS(cacheName, allURLs)
}
