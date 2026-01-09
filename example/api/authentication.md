# Authentication

> **Note:** Authentication features are planned for a future release. This page documents the planned API.

## Overview

Volcano's development server is intended for local preview only and does not include authentication. For production deployments, use a proper web server like nginx or a CDN.

## Planned Features

### Basic Auth

Future versions may support optional basic authentication for the development server:

```bash
volcano -s --auth user:password ./output
```

### Token-based Auth

For API access, we're considering JWT-based authentication:

```go
// Planned API
type AuthConfig struct {
    Enabled  bool
    Secret   string
    TokenTTL time.Duration
}
```

## Current Recommendations

For now, if you need to protect your preview site:

1. **Use a reverse proxy** like nginx with basic auth
2. **Deploy to a private server** with firewall rules
3. **Use a VPN** to restrict access

## Example nginx Config

```nginx
server {
    listen 80;
    server_name docs.internal.example.com;

    auth_basic "Documentation";
    auth_basic_user_file /etc/nginx/.htpasswd;

    location / {
        proxy_pass http://localhost:1776;
    }
}
```

## Feedback

If you have specific authentication needs, please [open an issue](https://github.com/example/volcano/issues) to let us know!
