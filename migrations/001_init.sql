-- indexer_state
CREATE TABLE IF NOT EXISTS indexer_state (
    indexer_name TEXT PRIMARY KEY,
    last_processed_block BIGINT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- erc20_transfers
CREATE TABLE IF NOT EXISTS erc20_transfers (
    id BIGSERIAL PRIMARY KEY,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    log_index INTEGER NOT NULL,
    contract_address TEXT NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    value NUMERIC NOT NULL
);

-- idempotency
CREATE UNIQUE INDEX erc20_transfer_unique
ON erc20_transfers (tx_hash, log_index);

-- performance indexes
CREATE INDEX erc20_block_idx
ON erc20_transfers (block_number);

CREATE INDEX erc20_contract_idx
ON erc20_transfers (contract_address);

CREATE INDEX erc20_from_idx
ON erc20_transfers (from_address);

CREATE INDEX erc20_to_idx
ON erc20_transfers (to_address);
