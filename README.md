# Designing Data-Intensive Applications — Chapter Guide with Projects & MIT 6.824 Videos

---

## Part I — Foundations of Data Systems

| Chapter Number | Chapter Name | MIT 6.824 Video(s) | Project Title |
|---:|---|---|---|
| 1 | Reliable, Scalable, and Maintainable Applications | Lecture 1 — Introduction | Build a Minimal URL Shortener Service |
| 2 | Data Models and Query Languages | None directly applicable | Multi-Model Blog System |
| 3 | Storage and Retrieval | None | Build a Log-Structured Key-Value Store |

---

### Chapter 1: Reliable, Scalable, and Maintainable Applications

**Key Concepts:** Reliability, scalability, maintainability, load, throughput, latency, availability, fault tolerance, monitoring.

**MIT 6.824 Video:** Lecture 1 — Introduction

**Project: Build a Minimal URL Shortener Service**

**Component — Description**
- **Functionality:** Shorten and redirect URLs; basic CRUD.
- **Infrastructure:** Run locally. Add simulated load (e.g., `hey` or `wrk`).
- **Observability:** Implement logging, health checks, uptime endpoints.
- **Optional:** Add rate-limiting, graceful restarts, simple caching.

**Goals:** Understand what reliability and scalability mean in code, how simple designs fail under load, and how monitoring supports maintainability.

---

### Chapter 2: Data Models and Query Languages

**Key Concepts:** Relational, document, and graph data models; normalization; query expressiveness; ACID vs BASE thinking.

**MIT 6.824 Video:** None directly applicable.

**Project: Multi-Model Blog System**

**Component — Description**
- **Scope:** Implement same data schema (Posts, Comments, Users) in three stores: PostgreSQL (relational), MongoDB (document), in-memory graph (adjacency list).
- **Querying:** Write identical queries: fetch post with comments, count posts per user, traverse relationships.
- **Comparison:** Benchmark response time and ease of expression.

**Goals:** See how data modeling changes between paradigms. Learn how query expressiveness affects system design decisions.

---

### Chapter 3: Storage and Retrieval

**Key Concepts:** B-trees, LSM-trees, SSTables, compaction, write amplification, log-structured storage, indexes.

**MIT 6.824 Video:** None.

**Project: Build a Log-Structured Key-Value Store**

**Component — Description**
- **Core:** Append-only log; hash index in memory.
- **Maintenance:** Add compaction to remove old entries.
- **Optional:** Add sorted-string tables (SSTables) and merges.
- **Testing:** Benchmark random read/write throughput.

**Goals:** Experience how low-level storage engines work — this forms the basis for databases, caches, and distributed stores later.

---

## Part II — Distributed Data

| Chapter Number | Chapter Name | MIT 6.824 Video(s) | Project Title |
|---:|---|---|---|
| 4 | Encoding and Evolution | Lecture 2 — RPC and Threads | Schema Evolution RPC Demo |
| 5 | Replication | Lecture 4 — Primary-Backup Replication; Lecture 9 — More Replication, CRAQ; Lecture 10 — Cloud Replicated DB: Aurora | Replicated Key-Value Store |
| 6 | Partitioning | Lecture 3 — Google File System (GFS) | Partitioned Hash Table Service |
| 7 | Transactions | Lecture 12 — Distributed Transactions; Lecture 13 — Spanner; Lecture 14 — Optimistic Concurrency Control | Transaction Manager Simulation |
| 8 | The Trouble with Distributed Systems | Lecture 6 — Raft (1); Lecture 7 — Raft (2) | Distributed Failure Simulator |
| 9 | Consistency and Consensus | Lecture 5 — Go, Threads, and Raft; Lecture 6 & 7 — Fault Tolerance (Raft 1 & 2); Lecture 8 — ZooKeeper | Simplified Raft Consensus Implementation |

---

### Chapter 4: Encoding and Evolution

**Key Concepts:** Data serialization, schemas, versioning, forward/backward compatibility (Protobuf, Avro, Thrift).

**MIT 6.824 Video:** Lecture 2 — RPC and Threads

**Project: Schema Evolution RPC Demo**

**Component — Description**
- **RPC Setup:** Implement simple gRPC service with a data schema.
- **Evolution:** Modify schema (add/remove field, rename) across versions.
- **Compatibility:** Demonstrate client/server interaction across versions.

**Goals:** Grasp the cost of schema evolution, version negotiation, and backward compatibility — vital for evolving APIs or microservices.

---

### Chapter 5: Replication

**Key Concepts:** Leader-based, follower-based, synchronous/asynchronous replication, read-your-writes, quorum consistency.

**MIT 6.824 Videos:** Lecture 4 — Primary-Backup Replication; Lecture 9 — More Replication, CRAQ; Lecture 10 — Cloud Replicated DB: Aurora

**Project: Replicated Key-Value Store**

**Component — Description**
- **Replication:** Leader + 2 followers. Implement write propagation.
- **Modes:** Support synchronous and asynchronous replication.
- **Faults:** Simulate follower crashes and recovery.
- **Optional:** Try multi-leader or leaderless replication.

**Goals:** Understand the tradeoff between latency and consistency, and how real databases replicate data under failures.

---

### Chapter 6: Partitioning

**Key Concepts:** Sharding, consistent hashing, rebalancing, range partitioning, hotspots, secondary indexes.

**MIT 6.824 Video:** Lecture 3 — Google File System (GFS)

**Project: Partitioned Hash Table Service**

**Component — Description**
- **Hashing:** Implement consistent hashing for data distribution.
- **Cluster:** Simulate multiple nodes as separate processes.
- **Rebalancing:** Add/Remove nodes dynamically.
- **Optional:** Handle range queries efficiently.

**Goals:** See how data is distributed in scalable systems and how rebalancing or hotspots can affect performance.

---

### Chapter 7: Transactions

**Key Concepts:** ACID properties, isolation levels, 2PL, 2PC, OCC, anomalies (lost update, dirty read, phantom read).

**MIT 6.824 Videos:** Lecture 12 — Distributed Transactions; Lecture 13 — Spanner; Lecture 14 — Optimistic Concurrency Control

**Project: Transaction Manager Simulation**

**Component — Description**
- **Isolation:** Implement read committed & serializable isolation.
- **Protocol:** Build two-phase locking and optimistic concurrency.
- **Tests:** Simulate concurrent clients performing conflicting writes.

**Goals:** Understand how real databases enforce consistency and how distributed transactions coordinate changes.

---

### Chapter 8: The Trouble with Distributed Systems

**Key Concepts:** Clocks, partial failures, split-brain, byzantine faults, unreliable networks, CAP tradeoffs.

**MIT 6.824 Videos:** Lecture 6 — Raft (1); Lecture 7 — Raft (2)

**Project: Distributed Failure Simulator**

**Component — Description**
- **Network:** Simulate packet loss, delay, and reorder.
- **Nodes:** Run distributed algorithm (e.g., counter or election).
- **Observation:** Log timing issues and inconsistencies.

**Goals:** Comprehend why distributed systems are hard — time drift, message loss, race conditions — before diving into consensus.

---

### Chapter 9: Consistency and Consensus

**Key Concepts:** Raft, Paxos, quorum writes, leader election, safety and liveness.

**MIT 6.824 Videos:** Lecture 5 — Go, Threads, and Raft; Lecture 6 & 7 — Fault Tolerance (Raft 1 & 2); Lecture 8 — ZooKeeper

**Project: Simplified Raft Consensus Implementation**

**Component — Description**
- **Core:** Leader election + log replication.
- **Cluster:** 3–5 nodes.
- **Faults:** Simulate leader failure, network partition.
- **Optional:** Build a replicated key-value store on top.

**Goals:** Master consensus — the cornerstone of reliability in distributed systems (etcd, Consul, Spanner, etc.).

---

## Part III — Derived Data

| Chapter Number | Chapter Name | MIT 6.824 Video(s) | Project Title |
|---:|---|---|---|
| 10 | Batch Processing | Lecture 15 — Big Data: Spark | Single-Machine MapReduce Framework |
| 11 | Stream Processing | None directly applicable | Streaming Analytics System |
| 12 | The Future of Data Systems | Lecture 11 — Cache Consistency: Frangipiani; Lecture 16 — Memcache at Facebook; Lecture 17 — COPS & Causal Consistency; Lecture 18 — Fork Consistency / Cert Transparency; Lecture 19 — Bitcoin; Lecture 20 — Blockstack | Hybrid Data System Architecture |

---

### Chapter 10: Batch Processing

**Key Concepts:** MapReduce, DAGs, dataflow, fault-tolerant computation, immutability.

**MIT 6.824 Video:** Lecture 15 — Big Data: Spark

**Project: Single-Machine MapReduce Framework**

**Component — Description**
- **Core:** Implement map/reduce with task scheduling.
- **Examples:** Word count, inverted index.
- **Reliability:** Handle worker failure (retry).

**Goals:** Learn large-scale data processing pipelines and why immutability simplifies distributed computing.

---

### Chapter 11: Stream Processing

**Key Concepts:** Windowing, event time vs processing time, exactly-once vs at-least-once semantics, backpressure.

**MIT 6.824 Video:** None directly applicable.

**Project: Streaming Analytics System**

**Component — Description**
- **Stream:** Build producer-consumer pipeline with message queue (e.g., Redis or Kafka clone).
- **Logic:** Implement windowed aggregations and joins.
- **Reliability:** Demonstrate replay and deduplication.

**Goals:** Experience low-latency dataflow design and tradeoffs in event processing systems (Kafka Streams, Flink, etc.).

---

### Chapter 12: The Future of Data Systems

**Key Concepts:** Hybrid systems, lambda/kappa architecture, eventual consistency, CRDTs, blockchain trust.

**MIT 6.824 Videos:** Lecture 11 — Cache Consistency: Frangipiani; Lecture 16 — Memcache at Facebook; Lecture 17 — COPS & Causal Consistency; Lecture 18 — Fork Consistency / Cert Transparency; Lecture 19 — Bitcoin; Lecture 20 — Blockstack

**Project: Hybrid Data System Architecture**

**Component — Description**
- **Layers:** Batch layer (historical), speed layer (real-time), serving layer (query).
- **Example:** E-commerce recommendations or social feed.
- **Optional:** Add blockchain for data integrity.

**Goals:** Integrate everything — combine historical and real-time systems into one consistent design.
