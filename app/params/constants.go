package params

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BaseCoinUnit         = "ubadge"
	AccountAddressPrefix = "bb"

	// EVMChainIDMainnet is the EVM chain ID for BitBadges mainnet
	// Chain ID: 50024 (to be claimed in ethereum-lists/chains registry)
	// This should match the chain_id in genesis under app_state.evm.params.chain_config.chain_id
	EVMChainIDMainnet = "50024"

	// EVMChainIDTestnet is the EVM chain ID for BitBadges testnet
	// Chain ID: 50025 (to be claimed in ethereum-lists/chains registry)
	// This should match the chain_id in genesis under app_state.evm.params.chain_config.chain_id
	EVMChainIDTestnet = "50025"
	
	// EVMChainIDLocalDev is the EVM chain ID for local development
	// Chain ID: 90123 (for local development only)
	EVMChainIDLocalDev = "90123"
	
	// CosmosChainIDMainnet is the Cosmos chain ID for BitBadges mainnet
	CosmosChainIDMainnet = "bitbadges-1"

	// CosmosChainIDTestnet is the Cosmos chain ID for BitBadges testnet
	CosmosChainIDTestnet = "bitbadges-2"
)

// Build-time EVM Chain ID set via ldflags
// This allows different binaries to have different chain IDs compiled in
// If not set at build time, defaults to local dev chain ID (90123)
var BuildTimeEVMChainID string

// GetEVMChainID returns the EVM chain ID to use.
// Uses build-time chain ID if set via ldflags, otherwise defaults to local dev (90123)
func GetEVMChainID() string {
	// If build-time chain ID is set, use it (for separate mainnet/testnet binaries)
	if BuildTimeEVMChainID != "" {
		return BuildTimeEVMChainID
	}
	
	// Default to local dev chain ID for local development
	return EVMChainIDLocalDev
}

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
