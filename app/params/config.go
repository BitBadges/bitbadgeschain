package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

func MakeEncodingConfig(cdc codec.Codec) EncodingConfig {
	return EncodingConfig{
		TxConfig: tx.NewTxConfig(cdc, tx.DefaultSignModes),
	}
}
