package keeper_functions

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/ai_test/testutil"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
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

// TestExecuteTransferTokens_IntermediateSenderIsCreator verifies that the intermediate
// sender (IBC-derived address) is used as the MsgTransferTokens Creator, NOT the original
// sender from the source chain. The error message should reference the intermediate address.
func (suite *ExecuteHookTestSuite) TestExecuteTransferTokens_IntermediateSenderIsCreator() {
	// Use two different addresses: intermediary (sender) vs original sender
	intermediateSender := suite.TestAccs[0]
	originalSender := suite.TestAccs[1].String()

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000))

	// Use a non-existent collection so the transfer fails with a predictable error
	// that includes the sender address in the message
	transfersJSON := json.RawMessage(`[{
		"from": "Mint",
		"to_addresses": ["bb1recipient"],
		"balances": [{"amount": "1", "ownership_times": [{"start": "1", "end": "18446744073709551615"}], "token_ids": [{"start": "1", "end": "1"}]}],
		"merkle_proofs": [],
		"eth_signature_proofs": [],
		"prioritized_approvals": []
	}]`)

	action := &customhookstypes.TransferTokensAction{
		CollectionId: "999999999", // non-existent collection
		Transfers:    transfersJSON,
		FailOnError:  true,
	}

	hookData := &customhookstypes.HookData{
		TransferTokens: action,
	}

	ack := suite.Keeper.ExecuteHook(suite.Ctx, intermediateSender, hookData, tokenIn, originalSender)
	suite.Require().False(ack.Success(), "should fail for non-existent collection")

	// Error should mention the collection not being found
	ackErr := testutil.GetAckError(ack)
	suite.Require().Contains(ackErr, "collection not found", "error should indicate collection not found")
}

// TestExecuteTransferTokens_IntermediateSenderGetsAutoApproval verifies that auto-approval
// flags are set for the intermediate sender address on the target collection, not for
// the original sender. This is critical because the intermediate address acts as the
// Creator in MsgTransferTokens and needs approval to transfer on behalf of the IBC sender.
func (suite *ExecuteHookTestSuite) TestExecuteTransferTokens_IntermediateSenderGetsAutoApproval() {
	intermediateSender := suite.TestAccs[0]
	originalSender := suite.TestAccs[2].String() // deliberately different address

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000))

	// Create a minimal collection with ID 1 so the hook can find it
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
	}
	err := suite.App.TokenizationKeeper.SetCollectionInStore(suite.Ctx, collection, true)
	suite.Require().NoError(err, "should create collection 1")

	// The transfer will fail (no mint approvals set up for this sender),
	// but auto-approval flags are set BEFORE the transfer attempt (on main ctx, not cached).
	transfersJSON := json.RawMessage(`[{
		"from": "Mint",
		"to_addresses": ["bb1recipient"],
		"balances": [{"amount": "1", "ownership_times": [{"start": "1", "end": "18446744073709551615"}], "token_ids": [{"start": "1", "end": "1"}]}],
		"merkle_proofs": [],
		"eth_signature_proofs": [],
		"prioritized_approvals": [],
		"only_check_prioritized_collection_approvals": false,
		"only_check_prioritized_incoming_approvals": false,
		"only_check_prioritized_outgoing_approvals": false
	}]`)

	action := &customhookstypes.TransferTokensAction{
		CollectionId: "1",
		Transfers:    transfersJSON,
		FailOnError:  true,
	}

	hookData := &customhookstypes.HookData{
		TransferTokens: action,
	}

	// Execute — may fail on the actual transfer, but auto-approval should be set regardless
	suite.Keeper.ExecuteHook(suite.Ctx, intermediateSender, hookData, tokenIn, originalSender)

	// Verify auto-approval flags were set for the INTERMEDIATE sender, not the original sender
	collection, found := suite.App.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, sdkmath.NewUint(1))
	suite.Require().True(found, "collection 1 should exist")

	intermediateBalances, _, err := suite.App.TokenizationKeeper.GetBalanceOrApplyDefault(suite.Ctx, collection, intermediateSender.String())
	suite.Require().NoError(err)
	suite.Require().True(intermediateBalances.AutoApproveAllIncomingTransfers, "intermediate sender should have auto-approve incoming")
	suite.Require().True(intermediateBalances.AutoApproveSelfInitiatedOutgoingTransfers, "intermediate sender should have auto-approve self-initiated outgoing")
	suite.Require().True(intermediateBalances.AutoApproveSelfInitiatedIncomingTransfers, "intermediate sender should have auto-approve self-initiated incoming")
}

