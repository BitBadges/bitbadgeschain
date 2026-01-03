package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/ai_test/testutil"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
)

type ExecuteHookTestSuite struct {
	testutil.AITestSuite
}

func TestExecuteHookTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteHookTestSuite))
}

func (suite *ExecuteHookTestSuite) TestExecuteHook_NilHookData() {
	sender := suite.TestAccs[0]
	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000))

	// Nil hook data should return success acknowledgement
	ack := suite.Keeper.ExecuteHook(suite.Ctx, sender, nil, tokenIn, sender.String())
	suite.Require().True(ack.Success(), "nil hook data should return success")
}

func (suite *ExecuteHookTestSuite) TestExecuteHook_EmptyHookData() {
	sender := suite.TestAccs[0]
	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000))

	hookData := &customhookstypes.HookData{
		SwapAndAction: nil,
	}

	// Empty hook data should return success acknowledgement
	ack := suite.Keeper.ExecuteHook(suite.Ctx, sender, hookData, tokenIn, sender.String())
	suite.Require().True(ack.Success(), "empty hook data should return success")
}

func (suite *ExecuteHookTestSuite) TestExecuteHook_InvalidSwap() {
	sender := suite.TestAccs[0]
	
	// Create an invalid swap that will fail
	invalidSwap := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     "999999", // Non-existent pool
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: suite.Bob,
			},
		},
	}

	hookData := &customhookstypes.HookData{
		SwapAndAction: invalidSwap,
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000))
	
	// Execute hook - should fail
	ack := suite.Keeper.ExecuteHook(suite.Ctx, sender, hookData, tokenIn, sender.String())
	suite.Require().False(ack.Success(), "hook should fail for invalid swap")
}

