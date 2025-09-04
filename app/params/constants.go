package params

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BaseCoinUnit         = "ubadge"
	AccountAddressPrefix = "bb"
)

var (
	initOnce sync.Once
)

func SetAddressPrefixes() {
	InitSDKConfig()
}

func InitSDKConfig() {
	initOnce.Do(func() {
		// Set prefixes
		accountPubKeyPrefix := AccountAddressPrefix + "pub"
		validatorAddressPrefix := AccountAddressPrefix + "valoper"
		validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
		consNodeAddressPrefix := AccountAddressPrefix + "valcons"
		consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

		// Set config
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
		config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
		config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
		config.Seal()
	})
}
