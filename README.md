# driftwatch

Lightweight daemon that detects configuration drift between running containers and their declared compose specs.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your Compose file and let it watch for drift:

```bash
driftwatch --compose docker-compose.yml --interval 30s
```

When drift is detected, `driftwatch` logs a diff to stdout and optionally sends an alert:

```
[DRIFT] service=api expected_image=myapp:1.4 actual_image=myapp:1.3
[DRIFT] service=worker expected_env=LOG_LEVEL=info actual_env=LOG_LEVEL=debug
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--compose` | `docker-compose.yml` | Path to Compose spec |
| `--interval` | `60s` | How often to check for drift |
| `--notify` | `` | Webhook URL for drift alerts |
| `--strict` | `false` | Exit non-zero on first drift detected |

---

## How It Works

`driftwatch` reads your Compose spec, queries the Docker daemon for running containers, and compares environment variables, image tags, port bindings, and volume mounts. Any mismatch is reported as drift.

---

## Requirements

- Go 1.21+
- Docker Engine with socket access (`/var/run/docker.sock`)

---

## License

MIT © 2024 yourusername