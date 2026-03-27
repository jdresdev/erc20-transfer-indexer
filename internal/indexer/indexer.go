package indexer

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jdresdev/erc20-transfer-indexer/internal/config"
	"github.com/jdresdev/erc20-transfer-indexer/internal/ethclient"
	"github.com/jdresdev/erc20-transfer-indexer/internal/storage"
)

var topics = [][]common.Hash{{transferEventSignature}}

type Indexer struct {
	client      *ethclient.Client
	storage     *storage.Storage
	cfg         *config.Config
	indexerName string
}

func New(client *ethclient.Client, storage *storage.Storage, cfg *config.Config, indexerName string) *Indexer {
	return &Indexer{
		client:      client,
		storage:     storage,
		cfg:         cfg,
		indexerName: indexerName,
	}
}

func (i *Indexer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := i.processNextBatch(ctx); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			slog.Error("batch failed", "error", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(i.cfg.PollInterval):
			}
		}
	}
}

func (i *Indexer) processNextBatch(ctx context.Context) error {
	cursor, err := i.storage.GetLastProcessedBlock(ctx, i.indexerName, i.cfg.StartBlock)
	if err != nil {
		return err
	}

	safeHead, err := i.getSafeHead(ctx)
	if err != nil {
		return err
	}

	if cursor >= safeHead {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(i.cfg.PollInterval):
		}
		return nil
	}

	start, end := i.computeRange(cursor, safeHead)

	logs, err := i.client.GetLogs(ctx, start, end, topics)
	if err != nil {
		return err
	}

	transfers := decodeTransfers(logs)

	if err := i.persistBatch(ctx, end, transfers); err != nil {
		return err
	}

	slog.Info("batch indexed", "start", start, "end", end, "transfers", len(transfers))
	return nil
}

func (i *Indexer) persistBatch(ctx context.Context, end uint64, transfers []storage.Transfer) error {
	tx, err := i.storage.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := i.storage.InsertTransfers(ctx, tx, transfers); err != nil {
		return err
	}

	if err := i.storage.UpdateLastProcessedBlock(ctx, tx, i.indexerName, end); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (i *Indexer) getSafeHead(ctx context.Context) (uint64, error) {
	head, err := i.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return head - i.cfg.Confirmations, nil
}

func (i *Indexer) computeRange(cursor, safeHead uint64) (uint64, uint64) {
	start := cursor + 1
	end := cursor + i.cfg.BatchSize

	if end > safeHead {
		end = safeHead
	}

	return start, end
}
