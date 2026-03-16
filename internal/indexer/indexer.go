package indexer

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jdresdev/erc20-transfer-indexer/internal/config"
	"github.com/jdresdev/erc20-transfer-indexer/internal/decoder"
	"github.com/jdresdev/erc20-transfer-indexer/internal/ethereum"
	"github.com/jdresdev/erc20-transfer-indexer/internal/storage"
)

type Indexer struct {
	client      *ethereum.Client
	storage     *storage.Storage
	cfg         *config.Config
	indexerName string
}

func New(client *ethereum.Client, storage *storage.Storage, cfg *config.Config, indexerName string) *Indexer {
	return &Indexer{
		client:      client,
		storage:     storage,
		cfg:         cfg,
		indexerName: indexerName,
	}
}

func (i *Indexer) Run(ctx context.Context) error {
	for {
		// read cursor
		cursor, err := i.storage.GetLastProcessedBlock(ctx, i.indexerName)
		if err != nil {
			return err
		}

		// get chain head
		head, err := i.client.BlockNumber(ctx)
		if err != nil {
			slog.Error("Error fetching block number", "error", err)
			time.Sleep(i.cfg.PollInterval)
		}

		// compute safe head
		safeHead := head - i.cfg.Confirmations

		if cursor >= safeHead {
			time.Sleep(i.cfg.PollInterval)
			continue
		}

		end := min(cursor+i.cfg.BatchSize, safeHead)

		topics := [][]common.Hash{
			{decoder.TransferEventSignature},
		}

		// fetch logs
		logs, err := i.client.GetLogs(ctx, cursor+1, end, topics)
		if err != nil {
			slog.Error("Error fetching logs", "error", err)
			time.Sleep(i.cfg.PollInterval)
			continue
		}

		// decode transfers
		transfers, err := decoder.DecodeTransfers(logs)
		if err != nil {
			slog.Error("Error decoding logs", "error", err)
			time.Sleep(i.cfg.PollInterval)
			continue
		}

		tx, err := i.storage.Pool().Begin(ctx)
		if err != nil {
			return err
		}

		// store transfers
		err = i.storage.InsertTransfers(ctx, tx, transfers)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		// update cursor
		err = i.storage.UpdateLastProcessedBlock(ctx, tx, i.indexerName, end)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		slog.Info("indexed block range", "start", cursor+1, "end", end, "transfers", len(transfers))

		time.Sleep(200 * time.Millisecond)
	}
}
