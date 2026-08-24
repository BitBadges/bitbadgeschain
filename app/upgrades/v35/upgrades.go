package v35

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v11/modules/apps/transfer/keeper"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	poolmanager "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	ratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

const (
	// UpgradeName is the v35 upgrade: the 9 -> 18 decimal migration.
	UpgradeName = "v35"
)

// Keepers is the set this upgrade needs. Passed as a struct because the list is
// long enough that positional arguments stop being readable.
type Keepers struct {
	Account      authkeeper.AccountKeeper
	Authz        authzkeeper.Keeper
	Bank         bankkeeper.BaseKeeper
	FeeGrant     feegrantkeeper.Keeper
	ERC20        erc20keeper.Keeper
	FeeMarket    feemarketkeeper.Keeper
	RateLimit    ratelimitkeeper.Keeper
	Transfer     *transferkeeper.Keeper
	Staking      *stakingkeeper.Keeper
	Mint         mintkeeper.Keeper
	Gov          *govkeeper.Keeper
	Distribution distrkeeper.Keeper
	EVM          *evmkeeper.Keeper
	Gamm         gammkeeper.Keeper
	PoolManager  poolmanager.Keeper
	Tokenization tokenizationkeeper.Keeper
}

// CustomUpgradeHandlerLogic runs the 9 -> 18 decimal migration.
//
// Order is load-bearing; the numbered comments in the body say why. The shape is:
// metadata, then bank, then every module that keeps its own copy of an amount,
// then the keys that embed the denom, then params.
//
// This is in a separate function so it can be run against a snapshot in a test.
func CustomUpgradeHandlerLogic(goCtx context.Context, k Keepers) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Denom metadata. x/vm derives the chain's decimals from the display
	//    unit's exponent here, so everything downstream that reads coin info
	//    must see 18 rather than 9 before it runs.
	MigrateDenomMetadata(ctx, k.Bank)

	// 2. Bank balances and supply. Everything backed by a bank balance —
	//    staking pools, gov's deposit escrow, GAMM pool reserves, IBC escrow
	//    accounts — converts as a side effect of this step.
	if _, err := RedenominateBank(ctx, k.Bank); err != nil {
		return err
	}

	// 3-6. The module-owned accounting that bank balances do not cover. Each of
	//      these holds amounts in its own store, and each would be understated
	//      by 10^9 relative to the coins backing it if skipped.
	if _, err := RescaleStaking(ctx, k.Staking); err != nil {
		return err
	}
	if _, err := RescaleDistribution(ctx, k.Distribution); err != nil {
		return err
	}
	if _, err := RescaleGovDeposits(ctx, k.Gov); err != nil {
		return err
	}
	if _, err := RescaleGammPools(ctx, k.Gamm); err != nil {
		return err
	}
	if _, err := RescaleTokenizationCoinTransfers(ctx, k.Tokenization); err != nil {
		return err
	}
	if _, err := RescaleVestingAccounts(ctx, k.Account); err != nil {
		return err
	}
	if _, err := RescaleGrants(ctx, k.FeeGrant, k.Authz); err != nil {
		return err
	}
	if _, err := RescaleEconomics(ctx, k.Mint, k.FeeMarket, k.Transfer, k.PoolManager); err != nil {
		return err
	}

	// 7. State that names the denom without holding an amount of it — no
	//    number changes, so a value-conservation check cannot catch these.
	if _, err := RepointDenomReferences(ctx, k.ERC20, k.RateLimit); err != nil {
		return err
	}

	// 8. Store keys that embed the denom. Not an amount, but orphaned all the
	//    same by the rename.
	if _, err := RekeyTakerFees(ctx, k.PoolManager); err != nil {
		return err
	}

	// 9. Module params last, so nothing above reads a half-updated denom.
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
