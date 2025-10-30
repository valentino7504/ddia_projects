# Minimal URL Shortener — DDIA Chapter 1 Project

### Overview

This project is part of my practical study of _Designing Data-Intensive Applications (DDIA)_ — **Chapter 1: Reliable, Scalable, and Maintainable Applications**.

The goal was not to build a production system but to _observe system behavior_ under different loads and failure modes — to **feel** what reliability, scalability, and maintainability look like in code.

---

## Functionality

| Endpoint         | Method | Description                  |
| ---------------- | ------ | ---------------------------- |
| `/api/shorten`   | POST   | Create a shortened URL       |
| `/:short_code`   | GET    | Redirect to the original URL |
| `/api/urls/{id}` | GET    | Retrieve stored URL details  |

---

## Tech Stack

- **Language:** Go
- **Database:** SQLite (local file store)
- **Load testing:** [`hey`](https://github.com/rakyll/hey)
- **Environment:** Localhost (single-node)

Optional observability features like `/healthz` and `/ready` are omitted for now (simple to add if deployed).

---

## Load Testing Setup

**Tool:** `hey`  
**Purpose:** Simulate concurrent requests to evaluate throughput, latency, and reliability.

Commands used:

```bash
hey -n 1000 -c 50 http://localhost:8080/<short_code>         # Redirect test
hey -n 1000 -c 50 http://localhost:4000/api/urls/1           # Read test
hey -n 1000 -c 50 -m POST -H "Content-Type: application/json" \
    -d '{"url": "https://example.com"}' http://localhost:4000/api/shorten  # Write test

```

## Results and Analysis

### Redirect Route - GET/{short_code}

#### Summary

> Total: 7.6477 seconds
> Slowest: 1.3431 seconds
> Fastest: 0.2060 seconds
> Average: 0.3609 seconds
> Requests/sec: 130.7575

#### Latency Distribution

| Percentile | Time(s) |
| ---------- | ------- |
| p50        | 0.2808  |
| p90        | 0.6683  |
| p95        | 1.1780  |
| p99        | 1.2694  |

#### Interpretation

- Stable median latency (~300ms) even under 50 concurrent clients.
- Long-tail outliers (~1s) caused by connection setup and redirect overhead - also my network speed is not particularly good lol.
- Zero failed responses (\[200] 1000 responses) → strong reliability.
- Throughput ~138 req/s limited by HTTP round-trip and client redirect behavior.

### Get URL Details Route - GET /api/urls/{short_code}

#### Summary

> Total: 0.0808 secs
> Slowest: 0.0180 secs
> Fastest: 0.0001 secs
> Average: 0.0033 secs
> Requests/sec: 12376.84

#### Latency Distribution

| Percentile | Time(s) |
| ---------- | ------- |
| p50        | 0.0019  |
| p90        | 0.0088  |
| p95        | 0.0109  |
| p99        | 0.0155  |

#### Interpretation

- Extremely fast reads (≈ 3ms average).
- High stability — 99% of requests complete in < 16ms.
- Throughput > 12k req/s — excellent local read scalability.
- Confirms that SQLite read operations scale well on a single node.
