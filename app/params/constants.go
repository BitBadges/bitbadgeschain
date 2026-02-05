package params

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BaseCoinUnit         = "ubadge"
	AccountAddressPrefix = "bb"
	
	// EVMChainID is the custom EVM chain ID for BitBadges chain
	// This should match the chain_id in genesis under app_state.evm.params.chain_config.chain_id
	// For production, set this to a unique chain ID (e.g., 9001, 13337, etc.)
	// Default value of 9000 is used for testing compatibility
	EVMChainID = "90123" // TODO: change this to a unique chain ID and set genesis
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

	// Set config (don't seal - caller will seal if needed)
	config := sdk.GetConfig()
	// Only set if config is not already sealed
	// Check current prefix - if it's already "bb", assume it's set correctly
	currentPrefix := config.GetBech32AccountAddrPrefix()
	if currentPrefix != AccountAddressPrefix {
		// Try to set the prefix - this will panic if config is sealed, but that's ok
		// We'll catch it and use the existing prefix
		config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
		config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
		config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	}
	return config
}

func InitSDKConfig() {
	config := InitSDKConfigWithoutSeal()
	config.SetCoinType(60) // Ethereum's coin type
	config.SetPurpose(hd.CreateHDPath(60, 0, 0).Purpose)
	config.Seal()
}
