```markdown
# Apollo Edge Portal

A minimal edge-device monitoring portal built with Go, MQTT (Eclipse Mosquitto), and Docker. The portal subscribes to sensor data over MQTT, keeps the latest reading in memory, and displays it on a single auto-refreshing HTML page with an industrial-style dashboard aesthetic.

## Overview

This project simulates an edge-device gateway. A Go HTTP server subscribes to an MQTT topic, stores the most recent sensor payload in memory, and renders it on a web page that refreshes every five seconds. Everything runs in Docker, so the entire stack starts with a single command.

## Architecture

```text
┌──────────────┐      MQTT       ┌────────────────┐      HTTP       ┌──────────────┐
│ Sensor /     │ ───────────────▶ │ Mosquitto      │ ◀────────────── │ Go Portal    │
│ Publisher    │                 │ Broker         │                 │ Dashboard    │
└──────────────┘                 └────────────────┘                 └──────────────┘
```

- **Mosquitto** provides the MQTT broker.
- **Go portal** subscribes to the sensor topic and serves the dashboard.
- **Dashboard** displays the latest reading and refreshes automatically.

## Running with Docker

Start the stack from the project directory:

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080) in a browser. Stop the services with:

```bash
docker compose down
```

## Publishing a test reading

If the broker exposes its MQTT port on `localhost:1883`, publish a JSON payload with any MQTT client. For example, using `mosquitto_pub`:

```bash
mosquitto_pub -h localhost -p 1883 -t sensors/reading \\
	-m '{"device":"edge-01","temperature":23.4,"humidity":48.2,"status":"online"}'
```

The dashboard displays the new payload on its next refresh.

## Configuration

The portal can be configured through environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `MQTT_BROKER` | `tcp://mosquitto:1883` | MQTT broker address |
| `MQTT_TOPIC` | `sensors/reading` | Topic to subscribe to |
| `HTTP_ADDR` | `:8080` | HTTP listen address |

## Development

Run the Go service locally with a reachable Mosquitto broker:

```bash
go run .
```

Build the service with:

```bash
go build ./...
```

## Notes

- The latest reading is stored in memory and is lost when the portal restarts.
- The dashboard is intentionally lightweight and does not require a database.
- In production, secure MQTT with authentication and TLS before exposing it beyond a trusted network.
```

Both containers run on the same Docker network (`apollo-edge-portal_default`), allowing the Go service to reach Mosquitto by its container name (`mosquitto`) instead of a hardcoded IP.

## Requirements

- Docker Desktop (with Docker Compose)
- A modern web browser
- (Optional, for manual testing) `mosquitto_pub` — already included inside the Mosquitto container, so no local install is needed

## Project Structure

## Installation & Running

1. Clone or download this repository.
2. Open a terminal in the project root (`apollo-edge-portal/`).
3. Run the setup script (recommended):
```bash
   ./setup.sh
```
   This checks that Docker is running, then builds and starts both containers.

   Or start manually with Docker Compose:
```bash
   docker compose up --build
```
4. Open your browser to:

To stop the stack:
```bash
docker compose down
```

## Accessing the Portal from Another Device on the Same Network

The `go-portal` service publishes port 8080 on all network interfaces (`"8080:8080"` in `docker-compose.yml`, not restricted to `127.0.0.1`). Any device on the same local network can reach the dashboard at:
Find the host machine's IP with `ipconfig` (Windows) and look for the IPv4 address on the active network adapter.

## MQTT Topic & Example Message

- **Topic:** `sensors/data`
- **Payload format (JSON):**
```json
  {
    "temperature": 25.5,
    "power": 3.2,
    "current": 12.0
  }
```

### Sending a Test Message (Windows PowerShell)

Because PowerShell handles nested quotes differently from Bash, use the `--%` stop-parsing token when publishing a JSON payload:

```powershell
docker --% exec -it mosquitto mosquitto_pub -t sensors/data -m "{\"temperature\": 25.5, \"power\": 3.2, \"current\": 12.0}"
```

After running this, refresh the browser (or wait up to 5 seconds for auto-refresh) to see the updated values and timestamp.

### Sending a Test Message (Linux/macOS)

```bash
docker exec -it mosquitto mosquitto_pub -t sensors/data -m '{"temperature": 25.5, "power": 3.2, "current": 12.0}'
```

## Health Check

The server exposes a simple health endpoint for monitoring:
Example response:
```json
{"status": "ok", "mqtt_connected": true}
```

## Design Decisions

- **In-memory state:** The latest sensor reading is stored in a package-level variable rather than a database, matching the "minimal edge device" scope of this assignment. Simpler, faster, and sufficient for a single live reading.
- **Auto-refresh via `<meta>` tag:** Rather than adding JavaScript (explicitly out of scope per the assignment), the page uses `<meta http-equiv="refresh" content="5">` to poll the server every 5 seconds.
- **Dynamic device name:** `os.Hostname()` is used to populate the device name shown on the dashboard, so it reflects the actual container/host running the service rather than a hardcoded string.
- **Dynamic connection status:** A `mqttConnected` boolean, set once the MQTT client successfully connects, drives the ONLINE/OFFLINE indicator on the page rather than a static label.
- **Structured logging:** All output uses Go's `log` package (not `fmt`) so every line is automatically timestamped, making it easier to trace events during troubleshooting.
- **Error handling strategy:**
  - Invalid/malformed JSON payloads are logged (with the raw payload) and skipped — the server keeps running.
  - A failure to start the HTTP server is fatal (`log.Fatalf`), since the service has no purpose without it.
  - A failure to connect to MQTT at startup is fatal (`log.Fatal`), since the service cannot function without a data source.
- **Industrial-style CSS:** Dark background with a subtle grid pattern, monospace font, and boxed/bordered data fields — chosen to evoke a real industrial control panel rather than a generic web page.
- **Multi-stage Dockerfile:** Keeps the final image small by building the Go binary in one stage and copying only the compiled binary, templates, and static assets into the final image.

## Troubleshooting

- **Style changes don't appear in the browser:** Static files (`static/`) are copied into the Docker image at build time, not mounted as a live volume. After editing `style.css`, rebuild with `docker compose down && docker compose up --build`, or do a hard refresh (`Ctrl+Shift+R`) if only the browser cache is stale.
- **`docker-compose.yml` breaks after editing in VS Code:** VS Code can sometimes convert YAML indentation to tabs, which breaks the file. If this happens, edit the file in Notepad instead and save it directly with the correct filename.
- **PowerShell errors when publishing MQTT messages with JSON:** Use the `docker --% exec ...` syntax shown above — the `--%` token stops PowerShell from re-interpreting the quotes inside the JSON payload.
- **Page shows `Status: OFFLINE` or zero values on first load:** This is expected before any MQTT message has been published. Send a test message (see above) and refresh.
- **`bash setup.sh` fails with "Docker is not running" even though Docker Desktop is open:** This happens when running the script through WSL without Docker's WSL Integration enabled (Docker Desktop → Settings → Resources → WSL Integration). As a workaround, run `docker compose up --build` directly instead of the script.