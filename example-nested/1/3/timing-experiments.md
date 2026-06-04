---
status: pending
priority: medium
issue_id: PROJ-105
tags: [perf, experiment]
dependencies: []
---

# Timing Experiments

Three runs across the candidate routing strategies. Numbers are median page-load times in milliseconds.

| Strategy | Cold | Warm |
|----------|-----:|-----:|
| Static fallback | n/a | n/a |
| Dynamic resolver | 14 | 6 |
| Cached resolver | 12 | 5 |
