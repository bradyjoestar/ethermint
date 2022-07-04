package ante

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tharsis/ethermint/x/evm/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"math/big"
)

type TxCache struct {
	HashMap         map[string]*TxCacheObject
	ChainId         *big.Int
	EthCfg          *params.ChainConfig
	TxCap           int32
	AnteHandlerStep int
	Params          evmtypes.Params
}

type TxCacheObject struct {
	Data   string
	Signer common.Address
	TxData evmtypes.TxData
}

func NewTxCache() *TxCache {
	return &TxCache{
		HashMap:         make(map[string]*TxCacheObject),
		TxCap:           5000,
		AnteHandlerStep: 0,
	}
}

func NewTxCacheObject(from common.Address) *TxCacheObject {
	return &TxCacheObject{
		Data:   "",
		Signer: from,
	}
}

func (txCache *TxCache) AsMessage(msg *evmtypes.MsgEthereumTx, signer ethtypes.Signer, baseFee *big.Int) (core.Message, error) {
	txData, err := types.UnpackTxData(msg.Data)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTx(txData.AsEthereumData())
	var from common.Address
	if value, ok := txCache.HashMap[msg.Hash]; ok {
		from = value.Signer
	} else {
		from, err = ethtypes.Sender(signer, tx)
		if err != nil {
			return nil, err
		}
	}

	var gasPrice, gasTipCap, gasFeeCap *big.Int
	gasPrice = new(big.Int).Set(tx.GasPrice())
	gasTipCap = new(big.Int).Set(tx.GasTipCap())
	gasFeeCap = new(big.Int).Set(tx.GasFeeCap())
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	if baseFee != nil {
		gasPrice = math.BigMin(gasPrice.Add(gasTipCap, baseFee), gasFeeCap)
	}

	ethMsg := ethtypes.NewMessage(from, tx.To(), tx.Nonce(), tx.Value(),
		tx.Gas(), gasPrice, new(big.Int).Set(tx.GasFeeCap()),
		new(big.Int).Set(tx.GasTipCap()), tx.Data(), tx.AccessList(), false)

	return ethMsg, err
}
