package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
)

// BadgesAliasDenomRouter is an adapter that makes badges keeper implement AliasDenomRouter
type BadgesAliasDenomRouter struct {
	badgesKeeper badgeskeeper.Keeper
}

// NewBadgesAliasDenomRouter creates a new adapter for badges keeper
func NewBadgesAliasDenomRouter(badgesKeeper badgeskeeper.Keeper) *BadgesAliasDenomRouter {
	return &BadgesAliasDenomRouter{
		badgesKeeper: badgesKeeper,
	}
}

// CheckIsAliasDenom implements AliasDenomRouter interface
func (r *BadgesAliasDenomRouter) CheckIsAliasDenom(ctx sdk.Context, denom string) bool {
	return r.badgesKeeper.CheckIsAliasDenom(ctx, denom)
}

// SendNativeTokensViaAliasDenom implements AliasDenomRouter interface
func (r *BadgesAliasDenomRouter) SendNativeTokensViaAliasDenom(ctx sdk.Context, recipientAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	return r.badgesKeeper.SendNativeTokensViaAliasDenom(ctx, recipientAddress, toAddress, denom, amount)
}

// FundCommunityPoolViaAliasDenom implements AliasDenomRouter interface
func (r *BadgesAliasDenomRouter) FundCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error {
	return r.badgesKeeper.FundCommunityPoolViaAliasDenom(ctx, fromAddress, toAddress, denom, amount)
}

// GetBalanceWithAliasRouting implements AliasDenomRouter interface
// Uses getMaxWrappableAmount flow via GetSpendableCoinAmountWithAliasRouting
// This function is only called when the prefix matches, so we can assume it's a badgeslp: denom
func (r *BadgesAliasDenomRouter) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// Calculate from badge balances using getMaxWrappableAmount flow
	amount, err := r.badgesKeeper.GetSpendableCoinAmountBadgesLPOnly(ctx, address, denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(denom, amount), nil
}
