package indexer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jdresdev/erc20-transfer-indexer/internal/storage"
)

var transferEventSignature = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55aebf5e3f4c7")

func decodeTransfers(logs []types.Log) []storage.Transfer {
	transfers := make([]storage.Transfer, 0, len(logs))

	for _, log := range logs {
		if len(log.Topics) != 3 {
			continue
		}

		if log.Topics[0] != transferEventSignature {
			continue
		}

		from := common.BytesToAddress(log.Topics[1].Bytes())
		to := common.BytesToAddress(log.Topics[2].Bytes())
		value := new(big.Int).SetBytes(log.Data)

		transfers = append(transfers, storage.Transfer{
			BlockNumber:     log.BlockNumber,
			TxHash:          log.TxHash.Hex(),
			LogIndex:        uint32(log.Index),
			ContractAddress: log.Address.Hex(),
			FromAddress:     from.Hex(),
			ToAddress:       to.Hex(),
			Value:           value.String(),
		})
	}

	return transfers
}
