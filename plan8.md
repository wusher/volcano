# PWA Support for Volcano

Make Volcano-generated sites fully installable as Progressive Web Apps with offline support.

## Feature Summary

- **Opt-in**: `--pwa` flag enables PWA generation
- **Offline**: All pages cached for full offline reading
- **Icons**: Auto-generated from existing favicon (PNG only, warn on SVG/ICO)
- **Theme**: Manifest theme-color syncs with accent color or defaults to `#3b82f6`

---

## CLI Interface

```bash
volcano build ./docs -o ./public --pwa
```

---

## Generated Files

```
output/
├── manifest.json          # Web App Manifest
├── sw.js                  # Service Worker (precaches all pages)
├── icon-192.png           # Auto-generated from favicon
├── icon-512.png           # Auto-generated from favicon
└── ...existing files...
```

---

## Implementation Details

### 1. Config Changes

**`cmd/config.go`** - Add field:
```go
PWA bool // Enable PWA manifest and service worker generation
```

**`cmd/build.go`** - Add flag:
```go
fs.BoolVar(&cfg.PWA, "pwa", false, "Enable PWA manifest and service worker for offline support")
```

**`internal/generator/generator.go`** - Add to Config:
```go
PWA bool // Enable PWA manifest and service worker generation
```

---

### 2. New Package: `internal/pwa/`

#### `manifest.go`

```go
type ManifestConfig struct {
    SiteTitle   string // Full site title
    Description string // Site description (from first page or empty)
    ThemeColor  string // From --accent-color or default #3b82f6
    BaseURL     string // Base URL path prefix
    HasIcons    bool   // Whether PWA icons were generated
}

type Manifest struct {
    Name            string         `json:"name"`
    ShortName       string         `json:"short_name"`      // Title truncated to 12 chars
    Description     string         `json:"description,omitempty"`
    StartURL        string         `json:"start_url"`       // "/" or baseURL + "/"
    Scope           string         `json:"scope"`           // "/" or baseURL + "/"
    Display         string         `json:"display"`         // "standalone"
    BackgroundColor string         `json:"background_color"` // "#ffffff"
    ThemeColor      string         `json:"theme_color"`
    Icons           []ManifestIcon `json:"icons,omitempty"`
}

func GenerateManifest(outputDir string, config ManifestConfig) error
func GetManifestLinkTag(baseURL string) string // Returns <link rel="manifest" ...>
```

#### `serviceworker.go`

```go
type ServiceWorkerConfig struct {
    BaseURL   string   // Base URL path prefix
    PageURLs  []string // All page URLs to precache (from tree.AllPages)
    AssetURLs []string // CSS, JS, icon URLs
}

func GenerateServiceWorker(outputDir string, config ServiceWorkerConfig) error
func GetServiceWorkerRegistration(baseURL string) string // Returns JS registration code
```

**Service Worker Strategy:**
- Precache all pages and assets at install time
- Cache version hash = SHA256 of all cached URLs (first 8 hex chars)
- Serve from cache first, fall back to network
- Clean up old caches on activate

#### `icons.go`

```go
var IconSizes = []int{192, 512}

type IconResult struct {
    Generated bool     // Whether icons were successfully generated
    Paths     []string // Paths to generated icon files
    Warning   string   // Warning if source too small or unsupported format
}

func GenerateIcons(faviconPath, outputDir string) (*IconResult, error)
func GetIconURLs(baseURL string) []string
```

**Icon Generation:**
- Uses `golang.org/x/image/draw` with CatmullRom resampling
- Supports PNG, JPG, GIF source formats
- SVG/ICO: Skip with warning (can't resize without external deps)
- Warn if source < 512px (icons may be blurry)

---

### 3. Generator Integration

**`internal/generator/generator.go`**

Add to Generator struct:
```go
pwaEnabled bool
```

Add to Generate() method, after generating all pages:
```go
// Step 8: Generate PWA assets if enabled
if g.config.PWA {
    if err := g.generatePWA(site.AllPages, foldersNeedingIndex); err != nil {
        return nil, fmt.Errorf("failed to generate PWA assets: %w", err)
    }
}
```

New method:
```go
func (g *Generator) generatePWA(allPages []*tree.Node, autoIndexFolders []*tree.Node) error {
    g.logger.Verbose("Generating PWA assets...")

    // 1. Generate icons from favicon
    iconResult, err := pwa.GenerateIcons(g.config.FaviconPath, g.config.OutputDir)
    if err != nil {
        return err
    }
    if iconResult.Warning != "" {
        g.logger.Warning(iconResult.Warning)
    }

    // 2. Generate manifest.json
    manifestConfig := pwa.ManifestConfig{
        SiteTitle:  g.config.Title,
        ThemeColor: g.config.AccentColor, // Empty = default in manifest.go
        BaseURL:    g.baseURL,
        HasIcons:   iconResult.Generated,
    }
    if err := pwa.GenerateManifest(g.config.OutputDir, manifestConfig); err != nil {
        return err
    }

    // 3. Collect all page URLs
    pageURLs := collectPageURLs(allPages, autoIndexFolders, g.baseURL)

    // 4. Collect asset URLs
    assetURLs := []string{}
    if g.cssURL != "" {
        assetURLs = append(assetURLs, g.cssURL)
    }
    if g.jsURL != "" {
        assetURLs = append(assetURLs, g.jsURL)
    }
    if iconResult.Generated {
        assetURLs = append(assetURLs, pwa.GetIconURLs(g.baseURL)...)
    }

    // 5. Generate service worker
    swConfig := pwa.ServiceWorkerConfig{
        BaseURL:   g.baseURL,
        PageURLs:  pageURLs,
        AssetURLs: assetURLs,
    }
    if err := pwa.GenerateServiceWorker(g.config.OutputDir, swConfig); err != nil {
        return err
    }

    g.logger.Verbose("  manifest.json")
    g.logger.Verbose("  sw.js")
    if iconResult.Generated {
        g.logger.Verbose("  icon-192.png, icon-512.png")
    }

    return nil
}

func collectPageURLs(allPages []*tree.Node, autoIndexFolders []*tree.Node, baseURL string) []string {
    urls := make([]string, 0, len(allPages)+len(autoIndexFolders)+1)

    // Add root
    if baseURL != "" {
        urls = append(urls, baseURL+"/")
    } else {
        urls = append(urls, "/")
    }

    // Add all pages
    for _, page := range allPages {
        urlPath := tree.GetURLPath(page)
        if baseURL != "" {
            urlPath = baseURL + urlPath
        }
        urls = append(urls, urlPath)
    }

    // Add auto-index folders
    for _, folder := range autoIndexFolders {
        urlPath := "/" + tree.SlugifyPath(folder.Path) + "/"
        if baseURL != "" {
            urlPath = baseURL + urlPath
        }
        urls = append(urls, urlPath)
    }

    return urls
}
```

---

### 4. Template Changes

**`internal/templates/renderer.go`** - Add to PageData:
```go
PWAEnabled bool // Whether PWA is enabled (adds manifest link + SW registration)
```

**`internal/templates/layout.html`** - Add in `<head>` section:
```html
{{if .PWAEnabled}}
<link rel="manifest" href="{{.BaseURL}}/manifest.json">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="default">
<meta name="apple-mobile-web-app-title" content="{{.SiteTitle}}">
{{end}}
```

Add before closing `</body>`:
```html
{{if .PWAEnabled}}
<script>
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('{{.BaseURL}}/sw.js');
}
</script>
{{end}}
```

---

### 5. Dependencies

Add to `go.mod`:
```
golang.org/x/image v0.x.x
```

Run: `go get golang.org/x/image/draw`

---

## Edge Cases

| Scenario | Handling |
|----------|----------|
| No favicon provided | Skip icon generation, omit icons from manifest |
| SVG favicon | Skip with warning (can't resize without rasterization) |
| ICO favicon | Skip with warning (complex container format) |
| Favicon < 512px | Generate icons but warn about quality |
| BaseURL set | Prefix all URLs in manifest, sw.js, and meta tags |
| Very large sites | Still precache all (user opted in) |
| InlineAssets + PWA | Works fine, SW caches pages not external assets |

---

## Files to Modify/Create

| File | Action |
|------|--------|
| `cmd/config.go` | Add `PWA bool` field |
| `cmd/build.go` | Add `--pwa` flag + help text |
| `cmd/generate.go` | Pass `PWA` to generator config |
| `internal/generator/generator.go` | Add `PWA` to Config, add `generatePWA()` method |
| `internal/pwa/manifest.go` | **New** - Generate manifest.json |
| `internal/pwa/serviceworker.go` | **New** - Generate sw.js |
| `internal/pwa/icons.go` | **New** - Resize favicon to PWA icons |
| `internal/templates/renderer.go` | Add `PWAEnabled` to PageData |
| `internal/templates/layout.html` | Add PWA meta tags + SW registration |
| `go.mod` | Add `golang.org/x/image` dependency |

---

## Costs

| Cost | Impact |
|------|--------|
| **New dependency** | `golang.org/x/image` adds ~2MB to module cache, minimal binary size impact |
| **Build time** | Icon resizing adds ~50-100ms per build (only when `--pwa` enabled) |
| **Output size** | +20-50KB per site (manifest.json ~500B, sw.js ~2KB, icons ~20-50KB depending on source) |
| **Complexity** | ~400 lines of new code across 3 files in `internal/pwa/` |
| **Maintenance** | Service worker caching strategies may need updates as web standards evolve |

---

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| **Stale cache** | Medium | Cache version hash changes when any URL changes, forcing SW update |
| **Large precache** | Low | Sites with 1000+ pages may have slow initial SW install; could add `--pwa-lazy` flag later for cache-on-visit |
| **Icon quality** | Medium | Warn when source < 512px; recommend users provide high-res favicon |
| **SVG/ICO not supported** | Low | Clear warning message; document PNG requirement for PWA icons |
| **Service worker bugs** | Low | SW is generated with conservative cache-first strategy; users can clear cache in DevTools |
| **BaseURL edge cases** | Low | Tested with existing baseURL handling in generator; URLs are prefixed consistently |
| **Browser compatibility** | Very Low | Service Workers supported in all modern browsers (Chrome, Firefox, Safari, Edge) |

---

## Verification

1. **Build with PWA flag:**
   ```bash
   volcano build ./example -o ./test-output --pwa --favicon ./example/favicon.png
   ```

2. **Check generated files:**
   ```bash
   ls test-output/manifest.json test-output/sw.js test-output/icon-*.png
   cat test-output/manifest.json
   ```

3. **Serve and test in browser:**
   ```bash
   volcano serve ./test-output -p 8080
   ```
   - Open Chrome DevTools > Application > Manifest (should show installable)
   - Check Service Workers tab (should be registered)
   - Go offline (Network > Offline) and navigate (should work)
   - Click "Install" in address bar (should install as app)

4. **Run tests:**
   ```bash
   go test -race ./internal/pwa/...
   go test -race ./...
   ```
