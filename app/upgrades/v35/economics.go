package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v11/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	poolmanager "github.com/bitbadges/bitbadgeschain/x/poolmanager"
)

// EconomicsMigrationResult reports what RescaleEconomics touched.
type EconomicsMigrationResult struct {
	MinterRescaled bool
	FeeMarket      bool
	EscrowDenoms   int
	PoolVolumes    int
}

// RescaleEconomics converts the remaining amount-bearing state outside the
// modules that own user balances directly.
//
// Four separate things, grouped because each is small and all four are easy to
// forget precisely because they are not "balances".
func RescaleEconomics(
	ctx sdk.Context,
	mk mintkeeper.Keeper,
	fmk feemarketkeeper.Keeper,
	tk *transferkeeper.Keeper,
	pmk poolmanager.Keeper,
	ek *evmkeeper.Keeper,
) (EconomicsMigrationResult, error) {
	var res EconomicsMigrationResult

	// --- x/mint ---
	//
	// AnnualProvisions is a token amount, not a rate: the block reward is
	// AnnualProvisions / blocksPerYear. Left unscaled, the chain would mint
	// 10^-9 of the intended inflation — a silent halt to staking rewards that
	// nobody would notice for weeks.
	//
	// Inflation and the Params rates (InflationMax/Min/RateChange, GoalBonded)
	// are dimensionless ratios and are deliberately untouched. MaxSupply is a
	// token amount and does scale; leaving it would cap supply at 10^-9 of the
	// intended ceiling and stop minting entirely once reached.
	// Neither AnnualProvisions nor MaxSupply carries a denom, so a re-run would
	// scale them again. MintDenom is the denom they are quoted in and is moved
	// by MigrateDenomParams at the end of the upgrade, which makes it the gate.
	mintParams, err := mk.Params.Get(ctx)
	if err != nil {
		return res, fmt.Errorf("reading mint params: %w", err)
	}
	if mintParams.MintDenom == legacyDenom() {
		minter, err := mk.Minter.Get(ctx)
		if err != nil {
			return res, fmt.Errorf("reading minter: %w", err)
		}
		if !minter.AnnualProvisions.IsNil() && !minter.AnnualProvisions.IsZero() {
			minter.AnnualProvisions = minter.AnnualProvisions.Mul(conversionDec)
			if err := mk.Minter.Set(ctx, minter); err != nil {
				return res, fmt.Errorf("setting minter: %w", err)
			}
			res.MinterRescaled = true
		}

		if !mintParams.MaxSupply.IsNil() && !mintParams.MaxSupply.IsZero() {
			mintParams.MaxSupply = mintParams.MaxSupply.Mul(ConversionFactor)
			if err := mk.Params.Set(ctx, mintParams); err != nil {
				return res, fmt.Errorf("setting mint params: %w", err)
			}
		}
	} else {
		ctx.Logger().Info("v35: mint already migrated, skipping", "mint_denom", mintParams.MintDenom)
	}

	// --- x/feemarket ---
	//
	// BaseFee and MinGasPrice are prices denominated in the EVM's coin, so they
	// move with it. Skipping them makes gas 10^9 times cheaper in real terms —
	// which is not a cosmetic issue but a spam vector, since the cost of
	// filling blocks collapses.
	//
	// MinGasMultiplier is a ratio and is left alone.
	//
	// Neither price carries a denom either. x/vm's EvmDenom is the denom they
	// are quoted in and MigrateDenomParams moves it, so it is the gate.
	evmDenom := ek.GetParams(ctx).EvmDenom
	if evmDenom == legacyDenom() {
		feeParams := fmk.GetParams(ctx)
		changed := false
		if !feeParams.BaseFee.IsNil() && !feeParams.BaseFee.IsZero() {
			feeParams.BaseFee = feeParams.BaseFee.Mul(conversionDec)
			changed = true
		}
		if !feeParams.MinGasPrice.IsNil() && !feeParams.MinGasPrice.IsZero() {
			feeParams.MinGasPrice = feeParams.MinGasPrice.Mul(conversionDec)
			changed = true
		}
		if changed {
			if err := fmk.SetParams(ctx, feeParams); err != nil {
				return res, fmt.Errorf("setting feemarket params: %w", err)
			}
			res.FeeMarket = true
		}
	} else {
		ctx.Logger().Info("v35: feemarket already migrated, skipping", "evm_denom", evmDenom)
	}

	// --- ibc-go transfer escrow totals ---
	//
	// The escrowed coins themselves live on per-channel escrow accounts and
	// move with the bank migration. This is ibc-go's separate running total per
	// denom, used to enforce that unescrow never exceeds what was escrowed.
	// Leave it stale and that invariant is off by 10^9 — which either blocks
	// legitimate withdrawals or stops catching a real over-withdrawal.
	var escrowed []sdk.Coin
	tk.IterateTokensInEscrow(ctx, []byte(transfertypes.KeyTotalEscrowPrefix), func(coin sdk.Coin) bool {
		if coin.Denom == legacyDenom() {
			escrowed = append(escrowed, coin)
		}
		return false
	})
	for _, coin := range escrowed {
		// Retire the old denom's total and record it under the new one.
		tk.SetTotalEscrowForDenom(ctx, sdk.NewCoin(legacyDenom(), sdkmath.ZeroInt()))
		tk.SetTotalEscrowForDenom(ctx, sdk.NewCoin(newDenom(), coin.Amount.Mul(ConversionFactor)))
		res.EscrowDenoms++
	}

	// --- x/poolmanager tracked volume ---
	//
	// Cumulative per-pool trade volume. Not consensus-critical in the way a
	// balance is, but it feeds fee skimming and any volume-based reporting, so
	// a 10^9 discontinuity in the series is worth avoiding.
	pools, err := pmk.AllPools(ctx)
	if err != nil {
		return res, fmt.Errorf("reading pools for volume rescale: %w", err)
	}
	for _, pool := range pools {
		volume := pmk.GetTotalVolumeForPool(ctx, pool.GetId())
		converted, did := convertCoinsChanged(volume)
		if !did {
			continue
		}
		pmk.SetVolume(ctx, pool.GetId(), converted)
		res.PoolVolumes++
	}

	ctx.Logger().Info(
		"v35: rescaled economics",
		"factor", ConversionFactor.String(),
		"minter", res.MinterRescaled,
		"feemarket", res.FeeMarket,
		"escrow_denoms", res.EscrowDenoms,
		"pool_volumes", res.PoolVolumes,
	)
	return res, nil
}
