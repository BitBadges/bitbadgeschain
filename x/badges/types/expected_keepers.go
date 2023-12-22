package types

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	// GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// SetAccount(ctx sdk.Context, acc types.AccountI)
	// NewAccount(ctx sdk.Context, acc types.AccountI) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	// SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}
