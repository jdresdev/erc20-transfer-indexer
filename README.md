# ERC20 Transfer Indexer (Go)

A production-oriented Ethereum ERC20 transfer indexer written in Go.

The service continuously scans Ethereum blocks, extracts ERC20 `Transfer` events, and stores them in PostgreSQL for fast querying.

This project demonstrates backend engineering practices for blockchain infrastructure including batch processing, idempotent indexing, and crash-safe persistence.

---

Main components:

- **Ethereum client** – fetches logs using `FilterLogs`
- **Indexer worker** – processes blocks in ranges
- **Decoder** – extracts ERC20 transfers
- **Storage layer** – persists transfers and cursor state

---

## Features

- ERC20 `Transfer` event indexing
- Batch block processing
- Idempotent inserts (`ON CONFLICT DO NOTHING`)
- Crash-safe cursor tracking
- Reorg protection with confirmations
- Structured logging
- PostgreSQL persistence
- Docker-ready

---

## Database Schema

### indexer_state

Tracks indexing progress.