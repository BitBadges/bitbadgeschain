package v37

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	poolmanager "github.com/bitbadges/bitbadgeschain/x/poolmanager"
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
	pmk poolmanager.Keeper,
) error {
	legacy := appparams.LegacyBaseCoinUnit
	updated := appparams.BaseCoinUnit

	// staking: bond denom and SDK v0.55 key rotation fee (whose denom must
	// match BondDenom to pass stakingtypes.Params.Validate).
	stakingParams, err := sk.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("reading staking params: %w", err)
	}
	stakingChanged := false
	if stakingParams.BondDenom == legacy {
		stakingParams.BondDenom = updated
		stakingChanged = true
	}
	if stakingParams.KeyRotationFee.Denom == legacy {
		stakingParams.KeyRotationFee = sdk.NewCoin(
			updated,
			stakingParams.KeyRotationFee.Amount.Mul(ConversionFactor),
		)
		stakingChanged = true
	} else if stakingChanged && stakingParams.KeyRotationFee.Denom == "" {
		stakingParams.KeyRotationFee = sdk.NewCoin(updated, sdkmath.ZeroInt())
	}
	if stakingChanged {
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

	// evm: denom, extended denom, persisted EvmCoinInfo, and process-global
	// coin info. At 18 decimals base and extended denom must be equal — x/vm's
	// validateCoinInfo rejects any 18-decimal config where they differ.
	evmParams := ek.GetParams(ctx)
	evmParams.EvmDenom = updated
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: updated}
	if err := ek.SetParams(ctx, evmParams); err != nil {
		return fmt.Errorf("setting evm params: %w", err)
	}
	coinInfo, err := ek.LoadEvmCoinInfo(ctx)
	if err != nil {
		return fmt.Errorf("loading evm coin info: %w", err)
	}
	if err := ek.SetEvmCoinInfo(ctx, coinInfo); err != nil {
		return fmt.Errorf("setting evm coin info: %w", err)
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

	// poolmanager: community-pool swap denom and whitelist (stored in
	// poolmanager's own KV store as of v36).
	pmParams := pmk.GetParams(ctx)
	pmChanged := false
	if pmParams.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo == legacy {
		pmParams.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo = updated
		pmChanged = true
	}
	for i, d := range pmParams.TakerFeeParams.CommunityPoolDenomWhitelist {
		if d == legacy {
			pmParams.TakerFeeParams.CommunityPoolDenomWhitelist[i] = updated
			pmChanged = true
		}
	}
	if pmChanged {
		if err := pmParams.Validate(); err != nil {
			return fmt.Errorf("validating poolmanager params: %w", err)
		}
		pmk.SetParams(ctx, pmParams)
	}

	ctx.Logger().Info("v37: repointed module params", "from", legacy, "to", updated)
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
