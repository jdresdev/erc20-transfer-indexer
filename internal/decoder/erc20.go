package decoder

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jdresdev/erc20-transfer-indexer/internal/storage"
)

var TransferEventSignature = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55aebf5e3f4c7")

func DecodeTransfers(logs []types.Log) ([]storage.Transfer, error) {
	transfers := make([]storage.Transfer, 0)

	for _, log := range logs {
		if len(log.Topics) != 3 {
			continue
		}

		if log.Topics[0] != TransferEventSignature {
			continue
		}

		from := common.HexToAddress(string(log.Topics[1].Hex()))
		to := common.HexToAddress(string(log.Topics[2].Hex()))
		value := new(big.Int).SetBytes(log.Data)

		transfer := storage.Transfer{
			BlockNumber:     log.BlockNumber,
			TxHash:          log.TxHash.Hex(),
			LogIndex:        uint32(log.Index),
			ContractAddress: log.Address.Hex(),
			FromAddress:     from.Hex(),
			ToAddress:       to.Hex(),
			Value:           value.String(),
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
