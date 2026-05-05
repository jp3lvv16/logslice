# logslice

Stream and filter structured log files by time range and field patterns from the terminal.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
logslice [flags] <logfile>
```

### Examples

Filter logs within a time range:
```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log
```

Filter by field pattern:
```bash
logslice --field "level=error" app.log
```

Combine time range and field filters:
```bash
logslice --from "2024-01-15T08:00:00Z" --field "service=api" --field "level=warn" app.log
```

Stream from stdin:
```bash
tail -f app.log | logslice --field "level=error"
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--field` | Filter by field pattern `key=value` (repeatable) |
| `--format` | Input format: `json`, `logfmt` (default: `json`) |

---

## License

MIT © yourname