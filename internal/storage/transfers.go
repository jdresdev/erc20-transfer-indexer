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

	rows := make([][]any, 0, len(transfers))

	for _, t := range transfers {
		rows = append(rows, []any{
			t.BlockNumber,
			t.TxHash,
			t.LogIndex,
			t.ContractAddress,
			t.FromAddress,
			t.ToAddress,
			t.Value,
		})
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"erc20_transfers"},
		[]string{
			"block_number",
			"tx_hash",
			"log_index",
			"contract_address",
			"from_address",
			"to_address",
			"value",
		},
		pgx.CopyFromRows(rows),
	)

	return err
}
