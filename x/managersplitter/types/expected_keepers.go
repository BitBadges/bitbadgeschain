package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// TokenizationKeeper defines the expected interface for the Tokenization module.
type TokenizationKeeper interface {
	// GetNextCollectionId gets the next collection ID
	GetNextCollectionId(ctx sdk.Context) sdkmath.Uint
	// Keeper methods needed for UniversalUpdateCollection
	// Note: UniversalUpdateCollection is accessed through msgServer, not directly from keeper
}

