package testutil

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
)

// GenerateSwapAndAction generates a basic SwapAndAction for testing
func GenerateSwapAndAction(poolID string, denomIn, denomOut string, minAmount string) *customhookstypes.SwapAndAction {
	return &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     poolID,
						DenomIn:  denomIn,
						DenomOut: denomOut,
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  denomOut,
				Amount: minAmount,
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", // Bob
			},
		},
	}
}

// GenerateSwapAndActionWithIBCTransfer generates a SwapAndAction with IBC transfer
func GenerateSwapAndActionWithIBCTransfer(poolID string, denomIn, denomOut string, minAmount string, channelID, recipient string) *customhookstypes.SwapAndAction {
	swapAndAction := GenerateSwapAndAction(poolID, denomIn, denomOut, minAmount)
	swapAndAction.PostSwapAction = &customhookstypes.PostSwapAction{
		IBCTransfer: &customhookstypes.IBCTransferInfo{
			IBCInfo: &customhookstypes.IBCInfo{
				SourceChannel: channelID,
				Receiver:      recipient,
			},
		},
	}
	return swapAndAction
}

// GenerateHookData generates a HookData for testing
func GenerateHookData(swapAndAction *customhookstypes.SwapAndAction) *customhookstypes.HookData {
	return &customhookstypes.HookData{
		SwapAndAction: swapAndAction,
	}
}

// GenerateTokenCoin generates a token coin for testing
func GenerateTokenCoin(denom string, amount int64) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewInt(amount))
}
