package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BaseCoinUnit         = "ubadge"
	AccountAddressPrefix = "bb"
)

func SetAddressPrefixes() {
	InitSDKConfigWithoutSeal()
}

func InitSDKConfigWithoutSeal() *sdk.Config {
	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	return config
}

func InitSDKConfig() {
	config := InitSDKConfigWithoutSeal()
	config.Seal()
}
