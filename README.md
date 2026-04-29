# cronwrap

Lightweight wrapper for cron jobs that adds structured logging and failure alerting.

## Installation

```bash
go install github.com/yourusername/cronwrap@latest
```

## Usage

Wrap any cron job command by prefixing it with `cronwrap`:

```bash
# Basic usage
cronwrap -- /usr/local/bin/backup.sh

# With alerting via webhook
cronwrap --alert-url https://hooks.example.com/notify -- /usr/local/bin/backup.sh

# With a custom job name for cleaner log output
cronwrap --name "nightly-backup" --alert-url https://hooks.example.com/notify -- /usr/local/bin/backup.sh
```

Output is structured JSON, making it easy to pipe into log aggregators like Loki or Datadog:

```json
{
  "level": "info",
  "job": "nightly-backup",
  "exit_code": 0,
  "duration_ms": 3421,
  "message": "job completed successfully",
  "timestamp": "2024-11-01T02:00:03Z"
}
```

Add it to your crontab as a drop-in replacement:

```cron
0 2 * * * cronwrap --name "nightly-backup" -- /usr/local/bin/backup.sh
```

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Human-readable job name | command basename |
| `--alert-url` | Webhook URL to POST on failure | none |
| `--timeout` | Kill job after duration (e.g. `30m`) | none |
| `--quiet` | Suppress stdout on success | false |

## Why cronwrap?

Standard cron provides no structured output, no alerting, and no timeout control. `cronwrap` fills that gap with zero configuration overhead and a single binary.

## License

MIT © 2024 yourusername