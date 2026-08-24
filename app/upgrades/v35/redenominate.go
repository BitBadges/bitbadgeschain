package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// ConversionFactor is 10^(18-9): the multiplier taking a ubadge amount to the
// equivalent abadge amount. One ubadge becomes 10^9 abadge, so nobody's holdings
// change value — only their representation.
var ConversionFactor = sdkmath.NewIntWithDecimal(1, appparams.BaseCoinDecimals-appparams.LegacyBaseCoinDecimals)

// redenominationModule is the module account borrowed to move supply around.
// It needs both Minter and Burner permissions; x/vm's is the natural choice on
// this chain since the denom being moved is the EVM's gas token. Its own
// balance is restored exactly, so borrowing it has no lasting effect.
const redenominationModule = evmtypes.ModuleName

// legacyDenom and newDenom name the two sides of the migration. Small helpers
// rather than bare constant references so every module's migration reads the
// same way and there is one place to look when auditing which denom moved.
func legacyDenom() string { return appparams.LegacyBaseCoinUnit }
func newDenom() string    { return appparams.BaseCoinUnit }

// RedenominationResult reports what the migration did.
type RedenominationResult struct {
	Holders      int
	LegacySupply sdkmath.Int
	NewSupply    sdkmath.Int
}

// RedenominateBank converts every ubadge balance to abadge at 10^9, and moves
// the supply with it.
//
// Why a redenomination at all: cosmos/evm v0.7 treats 18-decimal gas tokens as
// the only fully supported configuration. A 9-decimal chain works today only
// because x/precisebank bridges the gap, and precisebank now lives in contrib
// as a compatibility shim. Going to 18 decimals removes the shim, unblocks
// EnableVirtualFeeCollection (which panics below 18 decimals, and gates the
// block-STM fee path), and stops every future cosmos/evm bump from being a
// coin-flip on whether non-18-decimal support survived.
//
// Balances are rewritten with UncheckedSetBalance because there is no transfer
// happening — this is a representation change, not a movement of value, and
// routing it through SendCoins would trip blocked-address checks on module
// accounts and fire transfer hooks for what is not a transfer.
//
// UncheckedSetBalance deliberately does not touch supply, so supply is moved
// separately: the legacy supply is burned in one operation and the new supply
// minted in one operation, through a module account whose own balance is
// restored exactly. Doing it per-holder would be O(n) mint calls and would
// leave supply briefly inconsistent at every step.
func RedenominateBank(ctx sdk.Context, bk bankkeeper.BaseKeeper) (RedenominationResult, error) {
	legacyDenom := appparams.LegacyBaseCoinUnit
	newDenom := appparams.BaseCoinUnit

	res := RedenominationResult{
		LegacySupply: sdkmath.ZeroInt(),
		NewSupply:    sdkmath.ZeroInt(),
	}

	// 1. Snapshot every legacy holder. Collected up front because the writes
	//    below invalidate the iterator's view.
	type holding struct {
		addr   sdk.AccAddress
		amount sdkmath.Int
	}
	var holders []holding
	summed := sdkmath.ZeroInt()

	bk.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom != legacyDenom || !coin.Amount.IsPositive() {
			return false
		}
		holders = append(holders, holding{addr: addr, amount: coin.Amount})
		summed = summed.Add(coin.Amount)
		return false
	})
	res.Holders = len(holders)

	legacySupply := bk.GetSupply(ctx, legacyDenom).Amount
	res.LegacySupply = legacySupply

	// A mismatch means some legacy supply is unaccounted for. Converting under
	// that assumption would silently mint or destroy value, so refuse.
	if !summed.Equal(legacySupply) {
		return res, fmt.Errorf(
			"%s balances sum to %s but supply is %s: refusing to redenominate with unaccounted supply",
			legacyDenom, summed, legacySupply,
		)
	}

	if legacySupply.IsZero() {
		return res, nil
	}

	newSupply := legacySupply.Mul(ConversionFactor)
	res.NewSupply = newSupply

	moduleAddr := authtypes.NewModuleAddress(redenominationModule)

	// 2. Zero every legacy balance.
	for _, h := range holders {
		if err := bk.UncheckedSetBalance(ctx, h.addr, sdk.NewCoin(legacyDenom, sdkmath.ZeroInt())); err != nil {
			return res, fmt.Errorf("clearing %s balance of %s: %w", legacyDenom, h.addr, err)
		}
	}

	// 3. Retire the legacy supply. Park the whole amount on the module account
	//    and burn it in one call, which zeroes both that balance and the supply.
	legacyTotal := sdk.NewCoins(sdk.NewCoin(legacyDenom, legacySupply))
	if err := bk.UncheckedSetBalance(ctx, moduleAddr, sdk.NewCoin(legacyDenom, legacySupply)); err != nil {
		return res, fmt.Errorf("staging legacy supply for burn: %w", err)
	}
	if err := bk.BurnCoins(ctx, redenominationModule, legacyTotal); err != nil {
		return res, fmt.Errorf("burning legacy supply: %w", err)
	}

	// 4. Mint the new supply in one call.
	newTotal := sdk.NewCoins(sdk.NewCoin(newDenom, newSupply))
	if err := bk.MintCoins(ctx, redenominationModule, newTotal); err != nil {
		return res, fmt.Errorf("minting redenominated supply: %w", err)
	}

	// 5. Hand it out at the converted amounts. If the module account was itself
	//    a holder this overwrites the freshly minted balance with its correct
	//    share, which is why step 6 only zeroes it when it was not.
	moduleWasHolder := false
	for _, h := range holders {
		converted := h.amount.Mul(ConversionFactor)
		if err := bk.UncheckedSetBalance(ctx, h.addr, sdk.NewCoin(newDenom, converted)); err != nil {
			return res, fmt.Errorf("setting %s balance of %s: %w", newDenom, h.addr, err)
		}
		if h.addr.Equals(moduleAddr) {
			moduleWasHolder = true
		}
	}
	if !moduleWasHolder {
		if err := bk.UncheckedSetBalance(ctx, moduleAddr, sdk.NewCoin(newDenom, sdkmath.ZeroInt())); err != nil {
			return res, fmt.Errorf("clearing staging balance: %w", err)
		}
	}

	ctx.Logger().Info(
		"v35: redenominated bank balances",
		"from", legacyDenom,
		"to", newDenom,
		"factor", ConversionFactor.String(),
		"holders", res.Holders,
		"legacy_supply", legacySupply.String(),
		"new_supply", newSupply.String(),
	)

	return res, nil
}

// MigrateDenomMetadata replaces the 9-decimal ubadge metadata with 18-decimal
// abadge metadata.
//
// This is load-bearing, not cosmetic: x/vm's LoadEvmCoinInfo derives the chain's
// decimals from the exponent of the display unit in this metadata. Leaving the
// old entry in place would keep the EVM believing the chain is 9-decimal.
func MigrateDenomMetadata(ctx sdk.Context, bk bankkeeper.BaseKeeper) {
	bk.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: "The native token of BitBadges Chain",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: appparams.BaseCoinUnit, Exponent: 0},
			{Denom: appparams.DisplayCoinUnit, Exponent: appparams.BaseCoinDecimals},
		},
		Base:    appparams.BaseCoinUnit,
		Display: appparams.DisplayCoinUnit,
		Name:    "Badge",
		Symbol:  "BADGE",
	})

	ctx.Logger().Info(
		"v35: set denom metadata",
		"base", appparams.BaseCoinUnit,
		"display", appparams.DisplayCoinUnit,
		"exponent", appparams.BaseCoinDecimals,
	)
}
