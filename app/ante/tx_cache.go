package ante

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type TxCache struct {
	HashMap map[string]TxCacheObject
	ChainId *big.Int
	Cfg     *params.ChainConfig
	Cap     int32
}

func NewTxCache() *TxCache {
	return &TxCache{
		HashMap: make(map[string]TxCacheObject),
		Cap:     5000}
}

type TxCacheObject struct {
	Data   string
	Signer common.Address
}

func NewTxCacheObject(from common.Address) TxCacheObject {
	return TxCacheObject{
		Data:   "",
		Signer: from,
	}
}
