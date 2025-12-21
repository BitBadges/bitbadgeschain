package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	GetBalance(context.Context, sdk.AccAddress, string) sdk.Coin
	GetAllBalances(context.Context, sdk.AccAddress) sdk.Coins
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SendCoins(context.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error
	MintCoins(context.Context, string, sdk.Coins) error
	BurnCoins(context.Context, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromAccountToModule(context.Context, sdk.AccAddress, string, sdk.Coins) error
	// Methods imported from bank should be defined here
}

// DistributionKeeper defines the expected interface for the Distribution module.
type DistributionKeeper interface {
	FundCommunityPool(context.Context, sdk.Coins, sdk.AccAddress) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
