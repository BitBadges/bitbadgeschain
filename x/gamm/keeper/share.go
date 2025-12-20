package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func (k Keeper) applyJoinPoolStateChange(ctx sdk.Context, pool poolmanagertypes.PoolI, joiner sdk.AccAddress, numShares osmomath.Int, joinCoins sdk.Coins) error {
	poolAddress := pool.GetAddress()

	err := k.SendCoinsToPoolWithAliasRouting(ctx, joiner, poolAddress, joinCoins)
	if err != nil {
		return err
	}
	err = k.MintPoolShareToAccount(ctx, pool, joiner, numShares)
	if err != nil {
		return err
	}

	err = k.setPool(ctx, pool)
	if err != nil {
		return err
	}

	k.RecordTotalLiquidityIncrease(ctx, joinCoins)

	// Global pool invariant check: ensure pool has enough underlying assets for all recorded liquidity
	// This check happens after liquidity increments/decrements
	err = k.CheckPoolLiquidityInvariant(ctx, pool)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) applyExitPoolStateChange(ctx sdk.Context, pool poolmanagertypes.PoolI, exiter sdk.AccAddress, numShares osmomath.Int, exitCoins sdk.Coins) error {
	poolAddress := pool.GetAddress()

	err := k.SendCoinsFromPoolWithAliasRouting(ctx, poolAddress, exiter, exitCoins)
	if err != nil {
		return err
	}

	err = k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares)
	if err != nil {
		return err
	}

	err = k.setPool(ctx, pool)
	if err != nil {
		return err
	}

	k.RecordTotalLiquidityDecrease(ctx, exitCoins)

	// Global pool invariant check: ensure pool has enough underlying assets for all recorded liquidity
	// This check happens after liquidity increments/decrements
	err = k.CheckPoolLiquidityInvariant(ctx, pool)
	if err != nil {
		return err
	}

	return nil
}

// MintPoolShareToAccount attempts to mint shares of a GAMM denomination to the
// specified address returning an error upon failure. Shares are minted using
// the x/gamm module account.
func (k Keeper) MintPoolShareToAccount(ctx sdk.Context, pool poolmanagertypes.PoolI, addr sdk.AccAddress, amount osmomath.Int) error {
	amt := sdk.NewCoins(sdk.NewCoin(types.GetPoolShareDenom(pool.GetId()), amount))

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amt)
	if err != nil {
		return err
	}

	return nil
}

// BurnPoolShareFromAccount burns `amount` of the given pools shares held by `addr`.
func (k Keeper) BurnPoolShareFromAccount(ctx sdk.Context, pool poolmanagertypes.PoolI, addr sdk.AccAddress, amount osmomath.Int) error {
	amt := sdk.Coins{
		sdk.NewCoin(types.GetPoolShareDenom(pool.GetId()), amount),
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	return nil
}
