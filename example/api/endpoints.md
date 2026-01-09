# Server Endpoints

When running in serve mode, Volcano provides a simple HTTP server for previewing your generated site.

## Starting the Server

```bash
volcano -s -p 8080 ./output
```

This starts the server at `http://localhost:8080`.

## Endpoints

### Static Files

All files in the output directory are served as static content.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Serves `index.html` |
| GET | `/<path>` | Serves `<path>/index.html` |
| GET | `/<file>.html` | Serves the HTML file |
| GET | `/<asset>` | Serves static assets |

### Clean URLs

The server supports clean URLs automatically:

- `/guides/` serves `/guides/index.html`
- `/about` serves `/about/index.html` or `/about.html`

### 404 Handling

If a file is not found, the server returns `404.html` if it exists, otherwise a plain 404 response.

## HTTP Headers

The server sets the following headers for development:

```http
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
```

This ensures you always see the latest content during development.

## Graceful Shutdown

Press `Ctrl+C` to stop the server. It will gracefully shutdown, finishing any pending requests.

## Example Session

```
$ volcano -s ./output
Serving ./output at http://localhost:1776
Press Ctrl+C to stop

GET  / 200 5ms
GET  /guides/ 200 2ms
GET  /about 200 1ms
^C
Shutting down server...
```
