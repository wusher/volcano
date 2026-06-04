---
status: pending
priority: low
issue_id: PROJ-106
tags: [perf, notes]
dependencies: []
---

# Observations

- Cold-cache penalty matters more than warm performance for this user.
- Skipping `os.Stat` on known-bad paths cuts ~3ms.
