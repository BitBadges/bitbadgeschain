package v35

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

const (
	// UpgradeName is the v35 upgrade: the 9 -> 18 decimal migration.
	UpgradeName = "v35"
)

// Keepers is the set this upgrade needs. Passed as a struct because the list is
// long enough that positional arguments stop being readable.
type Keepers struct {
	Bank     bankkeeper.BaseKeeper
	Staking  *stakingkeeper.Keeper
	Mint     mintkeeper.Keeper
	Gov      *govkeeper.Keeper
	EVM          *evmkeeper.Keeper
	Tokenization tokenizationkeeper.Keeper
}

// CustomUpgradeHandlerLogic runs the 9 -> 18 decimal migration.
//
// Order matters and is not arbitrary:
//
//  1. Denom metadata first. x/vm derives the chain's decimals from the display
//     unit's exponent here, and anything downstream that reads coin info needs
//     to see 18 rather than 9.
//  2. Bank balances and supply. Everything backed by a bank balance — the
//     staking pools, gov's deposit escrow, GAMM pool reserves, IBC escrow
//     accounts — converts as a side effect of this step.
//  3. Staking's own accounting, which is not bank-backed and so has to be
//     rescaled explicitly.
//  4. Module params, last, so nothing above reads a half-updated denom.
//
// This is in a separate function so it can be run against a snapshot in a test.
func CustomUpgradeHandlerLogic(goCtx context.Context, k Keepers) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	MigrateDenomMetadata(ctx, k.Bank)

	if _, err := RedenominateBank(ctx, k.Bank); err != nil {
		return err
	}

	if _, err := RescaleStaking(ctx, k.Staking); err != nil {
		return err
	}

	if err := MigrateDenomParams(ctx, k.Staking, k.Mint, k.Gov, k.EVM, k.Tokenization); err != nil {
		return err
	}

	return nil
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
