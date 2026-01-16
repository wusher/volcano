# Examples

## Project Types

- [[documentation-site]] — Technical docs with organized sections
- [[blog]] — Date-ordered posts with blog theme
- [[knowledge-base]] — Wiki-style interconnected notes

## Common Commands

**Production build:**
```bash
volcano ./docs -o ./public --title="My Site" --url="https://example.com" --search
```

**Development server:**
```bash
volcano -s ./docs
```

**Custom theme:**
```bash
volcano css -o theme.css
volcano ./docs --css ./theme.css
```
