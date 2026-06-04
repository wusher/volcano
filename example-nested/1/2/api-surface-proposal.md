---
status: pending
priority: high
issue_id: PROJ-103
tags: [api, proposal]
dependencies: []
---

# API Surface Proposal

Sketch of the public functions we expect to commit to.

```go
func ServeNested(dir string, port int) error
```

The signature is deliberately small — only the inputs the caller can be expected to know up front.
