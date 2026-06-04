# How-To: OpenTelemetry Observability (APM Tracing)

This guide provides a practical recipe explaining how to spin up a local observability stack (such as Jaeger) and configure the Northstar backend to export live, distributed trace spans.

---

## Running a Local Jaeger Observability Stack

The easiest way to view traces locally is using **Jaeger**, an open-source distributed tracing APM platform. Jaeger supports standard OpenTelemetry OTLP ingestion out of the box.

### Prerequisites
Ensure you have **Docker** and **docker-compose** installed on your system.

### Step 1: Spin up Jaeger via Docker
Run the following command to start a local Jaeger container exposing the OTLP HTTP receiver on port `4318` and the search UI on port `16686`:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

Once running, access the Jaeger Search UI at **`http://localhost:16686`**.

---

## Configuring the Northstar Backend

By default, Northstar attempts to send HTTP/JSON OTLP spans to `localhost:4318`. If no collector is running, the background OTel engine will gracefully log a message and continue without blocking.

To customize or activate full tracing, configure the following environment variables:

| Environment Variable | Recommended Value | Purpose |
|----------------------|-------------------|---------|
| `OTEL_ENABLED` | `true` | Explicitly enables OpenTelemetry trace exporting (disabled by default to keep local terminals silent). |
| `OTEL_SERVICE_NAME` | `northstar-cmdb` | Identifies this service inside your APM dashboard. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4318` | The OTLP HTTP receiver host and port. |
| `OTEL_EXPORTER_OTLP_INSECURE` | `true` | Allows unencrypted HTTP connections for local testing. |

---

## Starting the App with Tracing Enabled

If you are using our `justfile` shortcuts, start the backend with the required environment variables:

```bash
OTEL_ENABLED="true" \
OTEL_SERVICE_NAME="northstar-cmdb" \
OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4318" \
OTEL_EXPORTER_OTLP_INSECURE="true" \
just backend
```

Now, navigate to the Nuxt frontend at `http://localhost:3000`, log in, and click around the dashboard or register a new asset.

### Verifying Traces in Jaeger

1. Open **`http://localhost:16686`** in your browser.
2. Under the **Service** dropdown, select **`northstar-cmdb`**.
3. Click **Find Traces**.
4. Click on any trace to inspect:
    * **HTTP Incoming Spans:** View execution times and status codes of REST routes.
    * **GORM Database Spans:** Drill down to see the exact SQL query strings (such as `SELECT * FROM categories`) executed against the SQLite database file!
