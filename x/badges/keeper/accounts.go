package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetOrCreateAccountNumberForAccAddressBech32(ctx sdk.Context, address sdk.AccAddress) uint64 {
 	account := k.accountKeeper.GetAccount(ctx, address)
	if account == nil {
		account = k.accountKeeper.NewAccountWithAddress(ctx, address)
		k.accountKeeper.SetAccount(ctx, account)
	}
	return account.GetAccountNumber()
}

func (k Keeper) MustGetAccountNumberForAddressString(ctx sdk.Context, address string) uint64 {
	acc_address := sdk.MustAccAddressFromBech32(address)
	return k.GetOrCreateAccountNumberForAccAddressBech32(ctx, acc_address)
}

//TODO: update this to be compatible with v0.46.0 AccountKeeper and test it
func (k Keeper) AssertAccountNumbersAreValid(ctx sdk.Context, accountNums []uint64) error {
	// for _, accountNum := range accountNums {
	// 	address := k.accountKeeper.GetAccountAddressByID(ctx, accountNum)
	// 	if address == "" {
	// 		return ErrAccountsAreNotRegistered
	// 	}
	// }

	return nil
}

