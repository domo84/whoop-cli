---
name: whoop-cli
description: >
  This skill should be used when the user asks about their WHOOP fitness data,
  recovery scores, sleep metrics, workout history, strain, HRV, resting heart rate,
  or any physiological data from their WHOOP wearable. Trigger phrases include:
  "show my recovery", "how did I sleep", "what's my HRV", "list my workouts",
  "WHOOP data", "my strain today", "recovery score", "sleep performance",
  "my cycles", "whoop stats", "fitness data".
allowed-tools: Bash
---

## Accessing WHOOP Data

Check authentication status before running any data command:

```sh
whoop auth status
```

If not authenticated, instruct the user to run `whoop auth setup` then `whoop auth login`.

## Mapping User Intent to Commands

| User asks about | Command |
|---|---|
| Recovery score, HRV, resting heart rate, SPO2 | `whoop recovery list` / `whoop recovery get <cycleId>` |
| Sleep, naps, sleep performance, sleep stages | `whoop sleep list` / `whoop sleep get <id>` |
| Workouts, strain, sport, duration | `whoop workout list` / `whoop workout get <id>` |
| Cycles, daily strain, energy | `whoop cycle list` / `whoop cycle get <id>` |
| Profile, body measurements | `whoop user profile` / `whoop user measurements` |

## Fetching Data

Use `--output json` when you need precise values for calculations or comparisons. Use `--output table` for summaries shown directly to the user.

```sh
# Last 5 recovery scores
whoop recovery list --limit 5 --output json

# Sleep records for a date range
whoop sleep list --start 2026-03-20 --end 2026-03-27 --output json

# All workouts (full history)
whoop workout list --all --output json

# Recovery for a specific cycle
whoop recovery get 1001 --output json
```

Date filters accept `YYYY-MM-DD` or RFC3339 format. Use `--all` only when the user asks for full history — it may be slow.

## Interpreting Key Fields

- `score_state`: `SCORED` = data ready, `PENDING_SCORE` = still processing, `UNSCORABLE` = insufficient data
- `recovery_score`: 0–100 percentage — higher is better
- `hrv_rmssd_milli`: heart rate variability in milliseconds
- `resting_heart_rate`: beats per minute
- `strain`: 0–21 scale (WHOOP strain score)
- `sleep_performance_percentage`: actual sleep vs needed sleep (%)
- Duration fields are in milliseconds (`_milli` suffix) — divide by 3,600,000 for hours

## Responding to the User

Surface the most relevant metric(s) for the question asked. For recovery questions lead with `recovery_score`, HRV, and resting heart rate. For sleep questions lead with `sleep_performance_percentage` and total sleep time. Avoid dumping raw JSON — summarize the key numbers in plain language.
