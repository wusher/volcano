---
status: pending
priority: high
issue_id: PROJ-202
tags: [launch, runbook]
dependencies: []
---

# Rollback Runbook

The short version: keep the previous binary on disk, watchman the health check, swap if 1% errors for 60 seconds.

> Don't reuse this runbook for anything that isn't a single-binary swap.
