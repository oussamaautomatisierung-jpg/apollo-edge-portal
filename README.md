 # Apollo Edge Portal

 A minimal edge device monitoring portal built with Go, MQTT, and Docker. The portal subscribes to sensor readings published over MQTT, keeps the latest values in memory, and serves them on a single auto-refreshing HTML page.

 ## Architecture

 MQTT Publisher → Mosquitto Broker → Go Server (MQTT Subscriber) → In-Memory State → HTTP Handler → Browser (auto-refresh every 5s)

 The Go server does two things concurrently:

 - Subscribes to the `sensors/data` topic on the Mosquitto broker and updates an in-memory struct whenever a new message arrives.
 - Serves an HTTP page at `/` that renders the latest in-memory values as HTML. The page auto-refreshes every 5 seconds without client-side JavaScript.

 A `/health` endpoint reports server and MQTT connection status as JSON, useful for container health checks and manual debugging.

 ## Project Structure

 ```text
 apollo-edge-portal/
 ├── main.go              # Go server: MQTT subscriber + HTTP handler
 ├── go.mod
 ├── go.sum
 ├── templates/
 │   └── index.html       # HTML template rendered at "/"
 ├── static/
 │   └── style.css        # Industrial-themed CSS
 ├── Dockerfile            # Multi-stage build for the Go server
 ├── docker-compose.yml    # Orchestrates go-portal + mosquitto
 ├── setup.sh              # Checks Docker is running, then starts the stack
 └── README.md
 ```

 ## Installation & Running

 ### Prerequisites

 - Docker and Docker Compose installed and running

 ### Quick Start

 ```sh
 ./setup.sh
 ```

 `setup.sh` verifies that the Docker daemon is running, then starts the full stack with `docker compose up`.

 The portal is available at [http://localhost:8080](http://localhost:8080).

 To stop the stack:

 ```sh
 docker compose down
 ```

 ## Testing

 Once the stack is running, publish a test sensor reading:

 ```sh
 docker exec -it mosquitto mosquitto_pub -t sensors/data -m '{"temperature": 25.5, "power": 3.2, "current": 12.0}'
 ```

 Refresh [http://localhost:8080](http://localhost:8080), or wait five seconds, to see the updated values and timestamp.

 Check server and MQTT connection status:

 ```sh
 curl http://localhost:8080/health
 ```

 Expected response:

 ```json
 {"status":"ok","mqtt_connected":true}
 ```

 ## Configuration

 | Environment Variable | Default | Description |
 |---|---|---|
 | `MQTT_BROKER` | `localhost` | Hostname of the MQTT broker, without scheme or port |

 The server connects to `tcp://<MQTT_BROKER>:1883`. Use a bare hostname such as `mosquitto` in Docker Compose, not a full URL.

 ## Design Decisions

 - **In-memory state only:** The latest reading is sufficient at this scale; a restart waits for the next MQTT message.
 - **No JavaScript frameworks:** A plain meta-refresh keeps the frontend dependency-free.
 - **Single MQTT topic:** Temperature, power, and current are sent as one JSON payload on `sensors/data`.
 - **Docker Compose as the run path:** Broker and server versions and networking remain consistent.

 ## Troubleshooting

 **Page shows “Disconnected” or `mqtt_connected: false`**

 - Confirm the broker is running: `docker compose ps`
 - Check server logs: `docker compose logs go-portal`

 **Published messages do not appear**

 - Publish to exactly `sensors/data`.
 - Ensure the payload is valid JSON with `temperature`, `power`, and `current` fields.

 **Port 8080 is already in use**

 Stop the process using port 8080 or change the host-side mapping in `docker-compose.yml`.
