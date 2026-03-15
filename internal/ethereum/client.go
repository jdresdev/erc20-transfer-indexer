package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth *ethclient.Client
}

func New(rpcURL string) (*Client, error) {
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		eth: c,
	}, nil
}

func (c *Client) Close() {
	c.eth.Close()
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

func (c *Client) GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, topics [][]common.Hash) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Topics:    topics,
	}

	return c.eth.FilterLogs(ctx, query)
}
