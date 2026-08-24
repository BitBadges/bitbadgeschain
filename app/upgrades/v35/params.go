package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// MigrateDenomParams repoints every module parameter that names the bond denom
// at the new 18-decimal denom, and rescales the amounts among them.
//
// Missing one of these is not a soft failure. A staking bond_denom still set to
// ubadge would make every delegation fail against a denom with zero supply; an
// evm_denom left behind would leave the EVM reading a denom nobody holds.
func MigrateDenomParams(
	ctx sdk.Context,
	sk *stakingkeeper.Keeper,
	mk mintkeeper.Keeper,
	gk *govkeeper.Keeper,
	ek *evmkeeper.Keeper,
	tk tokenizationkeeper.Keeper,
) error {
	legacy := appparams.LegacyBaseCoinUnit
	updated := appparams.BaseCoinUnit

	// staking: bond denom
	stakingParams, err := sk.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("reading staking params: %w", err)
	}
	if stakingParams.BondDenom == legacy {
		stakingParams.BondDenom = updated
		if err := sk.SetParams(ctx, stakingParams); err != nil {
			return fmt.Errorf("setting staking params: %w", err)
		}
	}

	// mint: mint denom
	mintParams, err := mk.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("reading mint params: %w", err)
	}
	if mintParams.MintDenom == legacy {
		mintParams.MintDenom = updated
		if err := mk.Params.Set(ctx, mintParams); err != nil {
			return fmt.Errorf("setting mint params: %w", err)
		}
	}

	// gov: deposit amounts are denominated *and* scaled
	govParams, err := gk.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("reading gov params: %w", err)
	}
	govParams.MinDeposit = convertCoins(govParams.MinDeposit, legacy, updated)
	govParams.ExpeditedMinDeposit = convertCoins(govParams.ExpeditedMinDeposit, legacy, updated)
	if err := gk.Params.Set(ctx, govParams); err != nil {
		return fmt.Errorf("setting gov params: %w", err)
	}

	// evm: denom and extended denom. At 18 decimals these must be equal —
	// x/vm's validateCoinInfo rejects any 18-decimal config where they differ.
	evmParams := ek.GetParams(ctx)
	evmParams.EvmDenom = updated
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: updated}
	if err := ek.SetParams(ctx, evmParams); err != nil {
		return fmt.Errorf("setting evm params: %w", err)
	}

	// tokenization: AllowedDenoms names the native denom explicitly. Leaving
	// ubadge in the list would allow a denom with zero supply and silently
	// disallow the real one.
	tokParams := tk.GetParams(ctx)
	changed := false
	for i, d := range tokParams.AllowedDenoms {
		if d == legacy {
			tokParams.AllowedDenoms[i] = updated
			changed = true
		}
	}
	if changed {
		tk.SetParams(ctx, tokParams)
	}

	ctx.Logger().Info("v35: repointed module params", "from", legacy, "to", updated)
	return nil
}

// convertCoins rewrites any coin in the legacy denom to the new denom, scaled by
// ConversionFactor. Coins in other denoms are passed through untouched — IBC
// vouchers and pool shares must not move.
func convertCoins(coins sdk.Coins, legacy, updated string) sdk.Coins {
	out := make(sdk.Coins, 0, len(coins))
	for _, c := range coins {
		if c.Denom == legacy {
			out = append(out, sdk.NewCoin(updated, c.Amount.Mul(ConversionFactor)))
			continue
		}
		out = append(out, c)
	}
	return out.Sort()
}
