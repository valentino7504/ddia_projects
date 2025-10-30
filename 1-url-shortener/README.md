# Minimal URL Shortener — DDIA Chapter 1 Project

### Overview

This project is part of my practical study of _Designing Data-Intensive Applications (DDIA)_ — **Chapter 1: Reliable, Scalable, and Maintainable Applications**.

The goal was not really to build a production system but to observe how systems behave under different loads.
It has also helped me understand more what reliability, scalability, and maintainability look like through load testing with _hey_

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

Total: 7.6477 seconds  
Slowest: 1.3431 seconds  
Fastest: 0.2060 seconds  
Average: 0.3609 seconds  
Requests/sec: 130.7575

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

Total: 0.0808 secs  
Slowest: 0.0180 secs  
Fastest: 0.0001 secs  
Average: 0.0033 secs  
Requests/sec: 12376.84

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

### Shorten Route - POST /api/shorten

#### Initial Concurrent Test (-c 50)

- A lot of failed writes because the database locked
- SQLite does not allow concurrent writes -> trades scalability for reliability
- Also created .db-journal file, which I obviously did not push lol
  **The results below are for the serialized test ie (-c 1)**

#### Summary

Total: 0.3196 secs  
Slowest: 0.0152 secs  
Fastest: 0.0026 secs  
Average: 0.0032 secs  
Requests/sec: 312.91

#### Latency Distribution

| Percentile | Time(s) |
| ---------- | ------- |
| p50        | 0.0029  |
| p90        | 0.0033  |
| p95        | 0.0040  |
| p99        | 0.0152  |

#### Interpretation

- Sustained 300+ req/s with only 1 writer.
- 3 ms average latency; 100% success rate (\[201]).
- SQLite is very reliable when writes are serialized
- SQLite trade off - strong consistency vs limited concurrency

## Concept Mapping

| Concept             | Things Learnt From This Project                                                                                          |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| Reliability         | No corrupt writes, consistent 200/201 responses under necessary concurrency limits                                       |
| Scalability         | 12k reads per second, write path restricted by sqlite to single-writer throughput (around 300 writes per second)         |
| Maintainability     | Simple Go service + sqlite db, nothing complex, plus logs                                                                |
| Load & load testing | Simulated using [`hey`](https://github.com/rakyll/hey) for concurrent HTTP traffic                                       |
| Throughput          | Measured total requests per second, with distinct profiles for read vs write load                                        |
| Latency             | Median (p50) and other percentile metrics show response time distribution under load                                     |
| Availability        | 0% error in stable configurations. Probably not possible in professional setting.                                        |
| Fault Tolerance     | SQLite journaling prevents corruption due to colliding writes. I just discovered SetMaxOpenConns(1) and have now updated |
| Monitoring          | Observed via request metrics and logs but there is room for more - eg health endpoints                                   |

## Key Takeaways

1. Reliability: The system stayed consistent and error free under expected load
2. Scalability: Reads scale easily; concurrent writes bottleneck instantly due to SQLite’s single-writer design
3. Maintainability: Minimal moving parts make debugging and iteration simple.

Overall the project shows:

- reliability through safe persistence,
- scalability through concurrency limits,
- maintainability through simplicity.

## Conclusion

A single-node, file-based system can be highly reliable and maintainable, but it will inevitably hit scalability limits under concurrent writes.
However, I wasn't trying to avoid it, just observe it.
