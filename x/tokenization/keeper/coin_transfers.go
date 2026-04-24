package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const (
	RoyaltyDivisor = 10000
	// ProtocolFeeDenominator is the divisor for the inclusive protocol fee
	// (0.1% = amount / 1000). Applied per-coin before royalty/recipient split,
	// so the payer sends exactly the quoted amount.
	ProtocolFeeDenominator = 1000
)

// formatDenomForDisplay formats a denom for display in error messages
// Shows "BADGE" for "ubadge" and prints others as-is
func formatDenomForDisplay(denom string) string {
	if denom == "ubadge" {
		return "BADGE"
	}
	return denom
}

func (k Keeper) ExecuteCoinTransfers(
	ctx sdk.Context,
	coinTransfers []*types.CoinTransfer,
	transferMetadata TransferMetadata,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.TokenCollection,
	userApprovalSettings *types.UserApprovalSettings,
) (string, error) {
	initiatedBy := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	if len(coinTransfers) == 0 {
		return "", nil
	}

	// Enforce UserApprovalSettings from collection-level for user-level coin transfers
	if userApprovalSettings != nil && (approvalLevel == "incoming" || approvalLevel == "outgoing") {
		if userApprovalSettings.DisableUserCoinTransfers {
			detErrMsg := "user-level coin transfers are disabled by collection approval settings"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
		if len(userApprovalSettings.AllowedDenoms) > 0 {
			for _, ct := range coinTransfers {
				for _, coin := range ct.Coins {
					allowed := false
					for _, denom := range userApprovalSettings.AllowedDenoms {
						if coin.Denom == denom {
							allowed = true
							break
						}
					}
					if !allowed {
						detErrMsg := fmt.Sprintf("denom %s is not allowed by collection approval settings (allowed: %v)", coin.Denom, userApprovalSettings.AllowedDenoms)
						return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
					}
				}
			}
		}
	}

	// Extract royalties from UserApprovalSettings
	royalties := &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(0),
		PayoutAddress: "",
	}
	if userApprovalSettings != nil && userApprovalSettings.UserRoyalties != nil {
		royalties = userApprovalSettings.UserRoyalties
	}

	if collection == nil {
		detErrMsg := "collection is nil"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	royaltyPercentage := royalties.Percentage
	royaltyPayoutAddress := royalties.PayoutAddress

	// Validate royalty percentage doesn't exceed 100% (RoyaltyDivisor = 10000)
	// This prevents panic when royaltyAmount > coin.Amount in the subtraction below
	if royaltyPercentage.GT(sdkmath.NewUint(RoyaltyDivisor)) {
		detErrMsg := fmt.Sprintf("royalty percentage cannot exceed %d (100%%), got %s", RoyaltyDivisor, royaltyPercentage.String())
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	if royaltyPercentage.GT(sdkmath.NewUint(0)) {
		if royaltyPayoutAddress == "" {
			detErrMsg := "payout address is required when royalty percentage is greater than 0"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
		}

		_, err := sdk.AccAddressFromBech32(royaltyPayoutAddress)
		if err != nil {
			return "", sdkerrors.Wrapf(err, "invalid payout address %s", royaltyPayoutAddress)
		}
	}

	for _, coinTransfer := range coinTransfers {
		if coinTransfer == nil {
			detErrMsg := "coin transfer is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
		if len(coinTransfer.Coins) == 0 {
			detErrMsg := "coin transfer cannot have empty coins slice"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Reject coinTransfers with badgeslp: denoms that reference the same collection.
	// Alias-routed transfers create a nested MsgTransferTokens which reads/writes balances
	// to the store, but the outer HandleTransfer holds stale in-memory copies that overwrite
	// those changes. Cross-collection or non-alias denoms (e.g., USDC) are unaffected.
	//
	// We can probably enable this in the future. Just be mindful of the stale state bug.
	for _, coinTransfer := range coinTransfers {
		for _, coin := range coinTransfer.Coins {
			if CheckStartsWithAliasDenom(coin.Denom) {
				denomCollectionId, err := ParseDenomCollectionId(coin.Denom)
				if err == nil && sdkmath.NewUint(denomCollectionId).Equal(collection.CollectionId) {
					detErrMsg := fmt.Sprintf("coinTransfer denom %s references the same collection %s — use an intermediate denom (e.g., USDC) instead", coin.Denom, collection.CollectionId.String())
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}
			}
		}
	}

	// Enforce allowed denoms from params. If the allowlist is non-empty, every coin denom
	// must either appear in the list or use the badgeslp: alias prefix (always allowed).
	allowedDenoms := k.GetAllowedDenoms(ctx)
	if len(allowedDenoms) > 0 {
		allowedSet := make(map[string]bool, len(allowedDenoms))
		for _, d := range allowedDenoms {
			allowedSet[d] = true
		}
		for _, coinTransfer := range coinTransfers {
			for _, coin := range coinTransfer.Coins {
				if !allowedSet[coin.Denom] && !CheckStartsWithAliasDenom(coin.Denom) {
					detErrMsg := fmt.Sprintf("denom %s is not in the allowed denoms list for coin transfers", coin.Denom)
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}
			}
		}
	}

	// Execute coin transfers directly - if they fail, the cached context will rollback
	for _, coinTransfer := range coinTransfers {
		coinsToTransfer := coinTransfer.Coins

		to := coinTransfer.To
		if coinTransfer.OverrideToWithInitiator {
			to = initiatedBy
		} else if types.IsMintAddress(to) {
			// "Mint" in the to field auto-resolves to the collection's mint escrow address.
			// ValidateBasic allows Mint for coin transfer targets; the same resolution applies
			// for collection, incoming, and outgoing approval execution paths.
			to = collection.MintEscrowAddress
		}

		toAddressAcc, err := sdk.AccAddressFromBech32(to)
		if err != nil {
			detErrMsg := fmt.Sprintf("invalid to address: %s", to)
			return detErrMsg, sdkerrors.Wrapf(err, "invalid to address: %s", to)
		}
		fromAddressAcc, err := sdk.AccAddressFromBech32(initiatedBy)
		if err != nil {
			detErrMsg := fmt.Sprintf("invalid initiatedBy address: %s", initiatedBy)
			return detErrMsg, sdkerrors.Wrapf(err, "invalid initiatedBy address: %s", initiatedBy)
		}
		if coinTransfer.OverrideFromWithApproverAddress {
			// collection-level
			if approverAddress == "" && approvalLevel == "collection" {
				approverAddress = collection.MintEscrowAddress
			}

			if approverAddress == "" {
				detErrMsg := "approver address is required when overrideFromWithApproverAddress is true"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
			}

			fromAddressAcc, err = sdk.AccAddressFromBech32(approverAddress)
			if err != nil {
				detErrMsg := fmt.Sprintf("invalid approver address: %s", approverAddress)
				return detErrMsg, sdkerrors.Wrapf(err, "invalid approver address: %s", approverAddress)
			}
		}

		for _, coin := range coinsToTransfer {
			// Inclusive protocol fee: taken out of the payer's gross amount first,
			// so the payer sends exactly coin.Amount (not amount + fee on top).
			coinAmountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			protocolFeeUint := coinAmountUint.Quo(sdkmath.NewUint(ProtocolFeeDenominator))
			protocolFeeInt := sdkmath.NewIntFromBigInt(protocolFeeUint.BigInt())

			royaltyAmountUint := coinAmountUint.Mul(royaltyPercentage).Quo(sdkmath.NewUint(RoyaltyDivisor))
			royaltyAmountInt := sdkmath.NewIntFromBigInt(royaltyAmountUint.BigInt())

			// Both fees come off the gross. With royalty capped at 100% the two could add
			// to more than gross; reject rather than silently shortchange the recipient.
			if protocolFeeInt.Add(royaltyAmountInt).GT(coin.Amount) {
				detErrMsg := fmt.Sprintf("royalty %s + protocol fee %s exceeds transfer amount %s for denom %s", royaltyAmountInt.String(), protocolFeeInt.String(), coin.Amount.String(), formatDenomForDisplay(coin.Denom))
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}

			remainingAmount := coin.Amount.Sub(royaltyAmountInt).Sub(protocolFeeInt)

			if err := k.sendProtocolFee(ctx, coin, protocolFeeInt, fromAddressAcc, coinTransfersUsed); err != nil {
				detErrMsg := fmt.Sprintf("insufficient %s balance to cover protocol fee", formatDenomForDisplay(coin.Denom))
				return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, err.Error())
			}

			err := k.sendCoinWithRoyalty(
				ctx,
				coin,
				royaltyAmountInt,
				remainingAmount,
				fromAddressAcc,
				toAddressAcc,
				royaltyPayoutAddress,
				coinTransfersUsed,
			)
			if err != nil {
				// Extract a more descriptive error message if possible
				detErrMsg := fmt.Sprintf("insufficient %s balance to complete transfer", formatDenomForDisplay(coin.Denom))
				if err.Error() != "" {
					// If the error already has a descriptive message, use it
					return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, err.Error())
				}
				return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, detErrMsg)
			}
		}
	}

	return "", nil
}

// sendProtocolFee routes the inclusive protocol fee from the payer to the community pool.
// No-op when the fee rounds to zero (amounts below ProtocolFeeDenominator).
func (k Keeper) sendProtocolFee(
	ctx sdk.Context,
	coin *sdk.Coin,
	feeAmount sdkmath.Int,
	fromAddressAcc sdk.AccAddress,
	coinTransfersUsed *[]CoinTransfers,
) error {
	if feeAmount.IsZero() {
		return nil
	}

	feeCoins := sdk.NewCoins(sdk.NewCoin(coin.Denom, feeAmount))
	if err := k.sendManagerKeeper.FundCommunityPoolWithAliasRouting(ctx, fromAddressAcc, feeCoins); err != nil {
		return sdkerrors.Wrapf(err, "error funding community pool with protocol fee: %s", feeCoins)
	}

	*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
		From:          fromAddressAcc.String(),
		To:            authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
		Amount:        feeAmount.String(),
		Denom:         coin.Denom,
		IsProtocolFee: true,
	})

	return nil
}

// sendCoinWithRoyalty handles sending a coin with royalty deduction
// It sends the royalty to the payout address and the remaining amount to the recipient
func (k Keeper) sendCoinWithRoyalty(
	ctx sdk.Context,
	coin *sdk.Coin,
	royaltyAmountInt sdkmath.Int,
	remainingAmount sdkmath.Int,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	royaltyPayoutAddress string,
	coinTransfersUsed *[]CoinTransfers,
) error {
	// Send royalty to payout address
	if !royaltyAmountInt.IsZero() {
		payoutAddressAcc, err := sdk.AccAddressFromBech32(royaltyPayoutAddress)
		if err != nil {
			return sdkerrors.Wrapf(err, "invalid royalty payout address: %s", royaltyPayoutAddress)
		}
		royaltyCoin := sdk.NewCoin(coin.Denom, royaltyAmountInt)

		err = k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, fromAddressAcc, payoutAddressAcc, &royaltyCoin)
		if err != nil {
			return sdkerrors.Wrapf(err, "error sending royalty to payout address")
		}

		*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
			From:   fromAddressAcc.String(),
			To:     payoutAddressAcc.String(),
			Amount: royaltyAmountInt.String(),
			Denom:  coin.Denom,
		})
	}

	// Send remaining amount to recipient
	if !remainingAmount.IsZero() {
		remainingCoin := sdk.NewCoin(coin.Denom, remainingAmount)

		err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, fromAddressAcc, toAddressAcc, &remainingCoin)
		if err != nil {
			return sdkerrors.Wrapf(err, "error sending remaining amount to recipient")
		}

		*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
			From:   fromAddressAcc.String(),
			To:     toAddressAcc.String(),
			Amount: remainingAmount.String(),
			Denom:  coin.Denom,
		})
	}

	return nil
}

// scaleCoinTransfers multiplies all coin amounts by the given integer multiplier.
// Used by allowAmountScaling to scale payments proportionally with transfer quantity.
func scaleCoinTransfers(base []*types.CoinTransfer, multiplier sdkmath.Uint) []*types.CoinTransfer {
	if multiplier.IsZero() {
		return []*types.CoinTransfer{}
	}
	
	multiplierInt := sdkmath.NewIntFromBigInt(multiplier.BigInt())
	scaled := make([]*types.CoinTransfer, len(base))
	for i, ct := range base {
		scaledCoins := make([]*sdk.Coin, len(ct.Coins))
		for j, coin := range ct.Coins {
			scaledCoin := sdk.NewCoin(coin.Denom, coin.Amount.Mul(multiplierInt))
			scaledCoins[j] = &scaledCoin
		}
		scaled[i] = &types.CoinTransfer{
			To:                             ct.To,
			Coins:                          scaledCoins,
			OverrideFromWithApproverAddress: ct.OverrideFromWithApproverAddress,
			OverrideToWithInitiator:         ct.OverrideToWithInitiator,
		}
	}
	return scaled
}
