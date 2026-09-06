package v35

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"

	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

// EIP-712 acceptance change (no store migration).
//
// v35 also changes which Cosmos-tx signatures the ante handler accepts, so
// every validator must run the same binary from the upgrade height:
//
//   - go-ethereum moves to cosmos/go-ethereum v1.17.2-cosmos-1, whose
//     EncodeType renders an empty struct type as `Name()` (upstream
//     go-ethereum#33702). EIP-712-signed txs whose typed-data tree contains an
//     empty sub-struct now verify; before, the chain-side hash never matched
//     the wallet's.
//   - The app now initialises cosmos/evm's EIP-712 verifier codecs
//     (eip712.SetEncodingConfig in app/evm.go) and registers the tokenization
//     and managersplitter amino msg types on the app's legacy amino codec, so
//     EIP-712 signatures over those msgs verify at all.
//
// Regression coverage: app/eip712_sign_test.go.

const (
	UpgradeName = "v35"

	// BlockMaxGas bounds the gas a single block (and therefore a single
	// transaction) may consume. The chain has run with the CometBFT default of
	// -1 (unbounded) since genesis; v35 aligns mainnet with the value the
	// shipped dev/testnet configs already use.
	BlockMaxGas int64 = 100_000_000

	// MinGasPrice is the fee-market floor, in ubadge per gas, applied to both
	// EVM and Cosmos transactions through cosmos/evm's ante handlers. The
	// live value is 0. A 200k-gas transaction at this floor costs 0.002 BADGE.
	MinGasPrice = "10.0"
)

// Keepers is the set of keepers the v35 handler needs. Held by value/pointer
// exactly as app.App stores them.
type Keepers struct {
	Account         authkeeper.AccountKeeper
	ConsensusParams consensuskeeper.Keeper
	FeeMarket       feemarketkeeper.Keeper
	IBCRateLimit    ibcratelimitkeeper.Keeper
	Tokenization    *tokenizationkeeper.Keeper
}

// CustomUpgradeHandlerLogic runs the chain-specific v35 work. It is separate
// from the handler so it can be exercised directly in tests.
func CustomUpgradeHandlerLogic(ctx context.Context, k Keepers) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := MigrateV35PreciseBankModulePermissions(sdkCtx, k.Account); err != nil {
		return fmt.Errorf("v35: precisebank permissions: %w", err)
	}

	if err := setBlockMaxGas(ctx, k.ConsensusParams); err != nil {
		return fmt.Errorf("v35: block max gas: %w", err)
	}
	if err := setMinGasPrice(sdkCtx, k.FeeMarket); err != nil {
		return fmt.Errorf("v35: min gas price: %w", err)
	}
	if err := k.IBCRateLimit.MigrateV35WindowsToBlockTime(sdkCtx); err != nil {
		return fmt.Errorf("v35: rate-limit windows: %w", err)
	}
	if err := runTokenizationMigrations(sdkCtx, k.Tokenization); err != nil {
		return fmt.Errorf("v35: tokenization: %w", err)
	}
	return nil
}

func setBlockMaxGas(ctx context.Context, ck consensuskeeper.Keeper) error {
	params, err := ck.ParamsStore.Get(ctx)
	if err != nil {
		return err
	}
	if params.Block == nil {
		return fmt.Errorf("consensus params have no block section")
	}
	params.Block.MaxGas = BlockMaxGas
	return ck.ParamsStore.Set(ctx, params)
}

func setMinGasPrice(ctx sdk.Context, fk feemarketkeeper.Keeper) error {
	params := fk.GetParams(ctx)
	minGasPrice, err := sdkmath.LegacyNewDecFromStr(MinGasPrice)
	if err != nil {
		return err
	}
	params.MinGasPrice = minGasPrice
	return fk.SetParams(ctx, params)
}

// runTokenizationMigrations applies the tokenization state migrations that
// accompany v35's validation changes. Order matters: address canonicalisation
// first, so every later walk sees one key per account.
func runTokenizationMigrations(ctx sdk.Context, tk *tokenizationkeeper.Keeper) error {
	return tk.MigrateV35CanonicalAddresses(ctx)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k Keepers,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if err := CustomUpgradeHandlerLogic(ctx, k); err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
