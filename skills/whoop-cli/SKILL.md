---
name: whoop-cli
description: Access WHOOP fitness wearable data — recovery scores, sleep, workouts, cycles, HRV, strain, and body metrics. Use when the user asks about their WHOOP data, fitness metrics, recovery, sleep performance, or workout history.
allowed-tools: Bash(whoop:*)
---

# WHOOP Data with whoop-cli

## Quick start

```bash
# check auth status
whoop auth status
# list latest recovery scores
whoop recovery list
# get today's sleep
whoop sleep list --limit 1
# list recent workouts
whoop workout list --limit 5
```

## Commands

### Auth

```bash
whoop auth setup           # interactively configure config.yaml (client ID, secret)
whoop auth login           # open browser for OAuth2 authentication
whoop auth logout          # revoke token and delete local credentials
whoop auth status          # show whether token is valid or expired
```

### User

```bash
whoop user profile         # name, email, user ID
whoop user measurements    # height, weight, max heart rate
```

### Recovery

```bash
whoop recovery list
whoop recovery list --limit 7                         # last 7 days
whoop recovery list --start 2026-03-01 --end 2026-03-27
whoop recovery list --all                             # full history
whoop recovery get <cycleId>                          # single recovery by cycle ID
```

### Sleep

```bash
whoop sleep list
whoop sleep list --limit 7
whoop sleep list --start 2026-03-01 --end 2026-03-27
whoop sleep list --all
whoop sleep get <id>                                  # single sleep record by ID
```

### Cycle

```bash
whoop cycle list
whoop cycle list --limit 7
whoop cycle list --start 2026-03-01 --end 2026-03-27
whoop cycle list --all
whoop cycle get <id>                                  # single cycle by numeric ID
whoop cycle sleep <cycleId>                           # sleep record linked to a cycle
```

### Workout

```bash
whoop workout list
whoop workout list --limit 10
whoop workout list --start 2026-03-01 --end 2026-03-27
whoop workout list --all
whoop workout get <id>                                # single workout by ID
```

## Global flags

```bash
--output table   # default — human-readable table
--output json    # full data, use for precise values and calculations
--limit N        # records per page (max 25)
--start DATE     # YYYY-MM-DD or RFC3339
--end DATE       # YYYY-MM-DD or RFC3339
--all            # fetch all pages (slow on large history)
```

## Example: Weekly recovery summary

```bash
whoop recovery list --start 2026-03-20 --end 2026-03-27 --output json
```

## Example: Last night's sleep

```bash
whoop sleep list --limit 1 --output json
```

## Example: This week's workouts

```bash
whoop workout list --start 2026-03-20 --end 2026-03-27 --output json
```

## Example: Full day picture (cycle + recovery + sleep)

```bash
whoop cycle list --limit 1 --output json
# use cycle ID from above
whoop recovery get <cycleId> --output json
whoop cycle sleep <cycleId> --output json
```

## Key fields

- `recovery_score` — 0–100%, higher is better
- `hrv_rmssd_milli` — HRV in milliseconds
- `resting_heart_rate` — bpm
- `strain` — 0–21 daily strain score
- `sleep_performance_percentage` — actual vs needed sleep
- `score_state` — `SCORED` ready, `PENDING_SCORE` processing, `UNSCORABLE` no data
- `_milli` fields — milliseconds; divide by 3,600,000 for hours
