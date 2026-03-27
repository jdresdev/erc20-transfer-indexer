package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transfer struct {
	BlockNumber     uint64
	TxHash          string
	LogIndex        uint32
	ContractAddress string
	FromAddress     string
	ToAddress       string
	Value           string
}

func (s *Storage) InsertTransfers(ctx context.Context, tx pgx.Tx, transfers []Transfer) error {
	if len(transfers) == 0 {
		return nil
	}

	blockNumbers := make([]uint64, len(transfers))
	txHashes := make([]string, len(transfers))
	logIndexes := make([]uint32, len(transfers))
	contractAddresses := make([]string, len(transfers))
	fromAddresses := make([]string, len(transfers))
	toAddresses := make([]string, len(transfers))
	values := make([]string, len(transfers))

	for idx, t := range transfers {
		blockNumbers[idx] = t.BlockNumber
		txHashes[idx] = t.TxHash
		logIndexes[idx] = t.LogIndex
		contractAddresses[idx] = t.ContractAddress
		fromAddresses[idx] = t.FromAddress
		toAddresses[idx] = t.ToAddress
		values[idx] = t.Value
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO erc20_transfers (block_number, tx_hash, log_index, contract_address, from_address, to_address, value)
		SELECT * FROM UNNEST($1::bigint[], $2::text[], $3::integer[], $4::text[], $5::text[], $6::text[], $7::numeric[])
		ON CONFLICT (tx_hash, log_index) DO NOTHING
	`, blockNumbers, txHashes, logIndexes, contractAddresses, fromAddresses, toAddresses, values)

	return err
}
