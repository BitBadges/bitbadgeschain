package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

// TokenizationAliasDenomRouter is an adapter that makes tokenization keeper implement AliasDenomRouter
type TokenizationAliasDenomRouter struct {
	tokenizationKeeper tokenizationkeeper.Keeper
}

// NewTokenizationAliasDenomRouter creates a new adapter for tokenization keeper
func NewTokenizationAliasDenomRouter(tokenizationKeeper tokenizationkeeper.Keeper) *TokenizationAliasDenomRouter {
	return &TokenizationAliasDenomRouter{
		tokenizationKeeper: tokenizationKeeper,
	}
}

// CheckIsAliasDenom implements AliasDenomRouter interface
func (r *TokenizationAliasDenomRouter) CheckIsAliasDenom(ctx sdk.Context, denom string) bool {
	return r.tokenizationKeeper.CheckIsAliasDenom(ctx, denom)
}

// SendNativeTokensViaAliasDenom implements AliasDenomRouter interface
func (r *TokenizationAliasDenomRouter) SendNativeTokensViaAliasDenom(ctx sdk.Context, recipientAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	return r.tokenizationKeeper.SendNativeTokensViaAliasDenom(ctx, recipientAddress, toAddress, denom, amount)
}

// FundCommunityPoolViaAliasDenom implements AliasDenomRouter interface
func (r *TokenizationAliasDenomRouter) FundCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	return r.tokenizationKeeper.FundCommunityPoolViaAliasDenom(ctx, fromAddress, toAddress, denom, amount)
}

// SpendFromCommunityPoolViaAliasDenom implements AliasDenomRouter interface
func (r *TokenizationAliasDenomRouter) SpendFromCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	return r.tokenizationKeeper.SpendFromCommunityPoolViaAliasDenom(ctx, fromAddress, toAddress, denom, amount)
}

// SendFromModuleToAccountViaAliasDenom implements AliasDenomRouter interface
// For tokenization keeper, this is just a standard send
func (r *TokenizationAliasDenomRouter) SendFromModuleToAccountViaAliasDenom(ctx sdk.Context, moduleAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	// For tokenization keeper, module-to-account is just a standard send
	return r.tokenizationKeeper.SendNativeTokensViaAliasDenom(ctx, moduleAddress, toAddress, denom, amount)
}

// SendFromAccountToModuleViaAliasDenom implements AliasDenomRouter interface
// For tokenization keeper, this is just a standard send
func (r *TokenizationAliasDenomRouter) SendFromAccountToModuleViaAliasDenom(ctx sdk.Context, fromAddress string, moduleAddress string, denom string, amount sdkmath.Uint) error {
	// For tokenization keeper, account-to-module is just a standard send
	return r.tokenizationKeeper.SendNativeTokensViaAliasDenom(ctx, fromAddress, moduleAddress, denom, amount)
}

// GetBalanceWithAliasRouting implements AliasDenomRouter interface
// Uses getMaxWrappableAmount flow via GetSpendableCoinAmountWithAliasRouting
// This function is only called when the prefix matches, so we can assume it's a tokenizationlp: denom
func (r *TokenizationAliasDenomRouter) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// Calculate from tokenization balances using getMaxWrappableAmount flow
	amount, err := r.tokenizationKeeper.GetSpendableCoinAmountBadgesLPOnly(ctx, address, denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(denom, amount), nil
}
