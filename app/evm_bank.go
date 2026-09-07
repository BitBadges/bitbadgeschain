package app

import (
	"context"
	"fmt"

	precisebankkeeper "github.com/cosmos/evm/contrib/x/precisebank/keeper"
	precisebanktypes "github.com/cosmos/evm/contrib/x/precisebank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// evmBankKeeper is the bank keeper handed to x/vm on this 9-decimal chain.
//
// x/vm v0.7 commits an account by writing its whole balance with an absolute
// UncheckedSetBalance in the extended denom. contrib/x/precisebank forwards
// that call to x/bank unchanged, which creates an extended-denom entry that
// precisebank's own reads never consult and nothing backs. x/vm v0.6 instead
// computed the delta against the current balance and minted or burned it
// through the bank wrapper; this type restores that behaviour on top of
// precisebank, so the integer denom plus the fractional store stay the only
// representation of the money.
//
// The 18-decimal redenomination removes precisebank altogether and with it
// the need for this wrapper.
type evmBankKeeper struct {
	precisebankkeeper.Keeper
}

var _ evmtypes.BankKeeper = evmBankKeeper{}

func newEVMBankKeeper(k precisebankkeeper.Keeper) evmBankKeeper {
	return evmBankKeeper{Keeper: k}
}

// UncheckedSetBalance sets the account's extended-denom balance to amt by
// minting or burning the difference through precisebank. Other denoms pass
// through unchanged.
func (k evmBankKeeper) UncheckedSetBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) error {
	if amt.Denom != precisebanktypes.ExtendedCoinDenom() {
		return k.Keeper.UncheckedSetBalance(ctx, addr, amt)
	}
	if amt.Amount.IsNegative() {
		return fmt.Errorf("cannot set negative balance %s for %s", amt, addr)
	}

	current := k.Keeper.GetBalance(ctx, addr, amt.Denom).Amount
	delta := amt.Amount.Sub(current)
	switch delta.Sign() {
	case 1:
		coins := sdk.NewCoins(sdk.NewCoin(amt.Denom, delta))
		if err := k.Keeper.MintCoins(ctx, evmtypes.ModuleName, coins); err != nil {
			return err
		}
		return k.Keeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, addr, coins)
	case -1:
		coins := sdk.NewCoins(sdk.NewCoin(amt.Denom, delta.Neg()))
		if err := k.Keeper.SendCoinsFromAccountToModule(ctx, addr, evmtypes.ModuleName, coins); err != nil {
			return err
		}
		return k.Keeper.BurnCoins(ctx, evmtypes.ModuleName, coins)
	}
	return nil
}

// LockedCoins reports the locked EVM-denom amount in extended units.
//
// x/vm adds the locked amount to the 18-decimal spendable balance without
// scaling it, so on a non-18-decimal chain it has to be scaled here.
func (k evmBankKeeper) LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	locked := k.Keeper.LockedCoins(ctx, addr)
	integerDenom := precisebanktypes.IntegerCoinDenom()
	amt := locked.AmountOf(integerDenom)
	if amt.IsZero() {
		return locked
	}
	return locked.
		Sub(sdk.NewCoin(integerDenom, amt)).
		Add(sdk.NewCoin(integerDenom, amt.Mul(precisebanktypes.ConversionFactor())))
}
