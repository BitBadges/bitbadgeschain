package params

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// BaseCoinUnit is the chain's atomic unit. As of the 18-decimal migration
	// this is abadge (atto-badge, 10^-18 BADGE), not the old 9-decimal ubadge.
	//
	// x/vm requires base == extended for an 18-decimal chain: validateCoinInfo
	// rejects any config where Decimals == 18 and Denom != ExtendedDenom. That
	// is also what lets x/precisebank go away — there is no longer a precision
	// gap for it to bridge.
	BaseCoinUnit         = "abadge"
	AccountAddressPrefix = "bb"

	// LegacyBaseCoinUnit is the pre-migration 9-decimal atomic unit. Retained
	// so the upgrade handler can find and convert balances still denominated
	// in it. Nothing new should be written in this denom.
	LegacyBaseCoinUnit = "ubadge"

	// ExtendedCoinUnit is what the EVM operates in. At 18 decimals it is the
	// same denom as the base unit.
	ExtendedCoinUnit = BaseCoinUnit

	// DisplayCoinUnit is the bank display unit; its exponent (BaseCoinDecimals)
	// is what x/vm LoadEvmCoinInfo reads to determine the chain's decimals.
	DisplayCoinUnit  = "badge"
	BaseCoinDecimals = 18

	// LegacyBaseCoinDecimals is the pre-migration exponent. The conversion
	// factor is 10^(BaseCoinDecimals - LegacyBaseCoinDecimals) = 10^9.
	LegacyBaseCoinDecimals = 9

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

// PowerReduction is how many atomic units make one unit of consensus power.
//
// This MUST move with the decimals. Consensus power is tokens/PowerReduction,
// so redenominating balances by 10^9 while leaving the SDK's default 10^6 in
// place would multiply every validator's power by 10^9. With a supply on the
// order of 10^15 ubadge, post-migration power would land around 10^18 —
// against CometBFT's MaxTotalVotingPower ceiling of ~8.2e18, with no headroom
// for supply growth.
//
// 10^15 = the SDK default 10^6 scaled by the same 10^9, which leaves every
// validator's power numerically identical across the upgrade.
//
// Setting this is a consensus-affecting global, so it takes effect for whatever
// binary is running. That is correct under cosmovisor, which replays each
// height range with the binary that produced it.
var PowerReduction = sdkmath.NewIntWithDecimal(1, 15)

// SetPowerReduction installs PowerReduction as the SDK-wide value. Called from
// the same place as the address prefixes, before any staking math runs.
func SetPowerReduction() {
	sdk.DefaultPowerReduction = PowerReduction
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
	SetPowerReduction()
	config := InitSDKConfigWithoutSeal()
	config.SetCoinType(60) // Ethereum's coin type
	config.SetPurpose(hd.CreateHDPath(60, 0, 0).Purpose)
	config.Seal()
}
