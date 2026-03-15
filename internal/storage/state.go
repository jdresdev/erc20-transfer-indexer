package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetLastProcessedBlock(ctx context.Context, indexerName string) (uint64, error) {
	var block uint64

	err := s.pool.QueryRow(ctx, `
	SELECT last_processed_block
	FROM indexer_state
	WHERE indexer_name = $1
	`, indexerName).Scan(&block)

	if err == pgx.ErrNoRows {
		_, err = s.pool.Exec(ctx, `
			INSERT INTO indexer_state (indexer_name, last_processed_block)
			VALUES ($1, 0)
		`, indexerName)
		if err != nil {
			return 0, err
		}

		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return block, nil
}

func (s *Storage) UpdateLastProcessedBlock(ctx context.Context, tx pgx.Tx, indexerName string, block uint64) error {
	_, err := tx.Exec(ctx, `
		UPDATE indexer_state
		SET last_processed_block = $1,
			updated_at = NOW()
		WHERE indexer_name = $2
	`, block, indexerName)

	return err
}
