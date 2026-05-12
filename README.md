# logslice

Fast log filtering and aggregation tool with regex pipelines and structured output formats.

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build ./...
```

## Usage

```bash
# Filter logs by pattern and output as JSON
logslice --input app.log --match "ERROR|WARN" --format json

# Chain multiple regex filters in a pipeline
logslice --input app.log --match "timeout" --exclude "healthcheck" --format table

# Aggregate log levels from stdin
cat app.log | logslice --aggregate level --format json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Input log file (or stdin) | stdin |
| `--match` | Regex pattern to include | — |
| `--exclude` | Regex pattern to exclude | — |
| `--aggregate` | Field to aggregate counts by | — |
| `--format` | Output format: `json`, `table`, `text` | `text` |

## Example Output

```json
[
  { "level": "ERROR", "count": 42, "message": "connection timeout" },
  { "level": "WARN",  "count": 17, "message": "retry attempt exceeded" }
]
```

## Requirements

- Go 1.21+

## License

MIT © yourname