package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jdresdev/erc20-transfer-indexer/internal/config"
	"github.com/jdresdev/erc20-transfer-indexer/internal/ethclient"
	"github.com/jdresdev/erc20-transfer-indexer/internal/indexer"
	"github.com/jdresdev/erc20-transfer-indexer/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store, err := storage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	client, err := ethclient.New(cfg.RPCURL)
	if err != nil {
		slog.Error("failed to connect to ethereum node", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	idx := indexer.New(client, store, cfg, "erc20_transfer_indexer")

	slog.Info("indexer starting", "start_block", cfg.StartBlock, "batch_size", cfg.BatchSize, "confirmations", cfg.Confirmations)

	if err := idx.Run(ctx); err != nil && err != context.Canceled {
		slog.Error("indexer stopped with error", "error", err)
		os.Exit(1)
	}

	slog.Info("indexer stopped")
}
