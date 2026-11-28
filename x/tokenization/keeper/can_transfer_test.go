package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestEvaluateCanTransfer_PassingAndZeroStateMutation() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	// Add a paid transfer approval with CoinTransfers + MaxNumTransfers=1 + AutoDeletion
	// so we prove NONE of those side effects leak out of EvaluateCanTransfer!
	feeCoin := sdk.NewInt64Coin("ubadge", 5_000_000)
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0], // keep mint-test for initial mint
		{
			ApprovalId:        "paid-once",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MustPrioritize:                 true,
				CoinTransfers: []*types.CoinTransfer{{
					To:    charlie,
					Coins: []*sdk.Coin{&feeCoin},
				}},
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1),
					AmountTrackerId:        "paid-once-tracker",
				},
				AutoDeletionOptions: &types.AutoDeletionOptions{
					AfterOneUse: true,
				},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Capture state snapshots BEFORE EvaluateCanTransfer.
	bobBankBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	charlieBankBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge")
	bobTokenBefore, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	aliceTokenBefore, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	colBefore, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	transfer := &types.Transfer{
		From:        bob,
		ToAddresses: []string{alice},
		Balances: []*types.Balance{{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		}},
		PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "paid-once",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		}},
	}

	res := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     transfer,
		InitiatedBy:  bob,
	})

	suite.Require().True(res.Allowed, "expected transfer to be allowed, got failures: %+v", res.FailedCriteria)
	suite.Require().Contains(res.MatchedApprovalIds, "paid-once")
	suite.Require().Empty(res.FailedCriteria)
	// RequiredCoins includes the 5,000,000 ubadge payment (+ 0.1% protocol fee = 5,005,000 ubadge)
	suite.Require().True(res.RequiredCoins.AmountOf("ubadge").GTE(sdkmath.NewInt(5_000_000)))

	// CRITICAL: Verify ZERO state mutation occurred in suite.ctx!
	bobBankAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	charlieBankAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge")
	bobTokenAfter, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	aliceTokenAfter, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	colAfter, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	suite.Require().Equal(bobBankBefore, bobBankAfter, "EvaluateCanTransfer mutated bob bank balance")
	suite.Require().Equal(charlieBankBefore, charlieBankAfter, "EvaluateCanTransfer mutated charlie bank balance")
	AssertBalancesEqual(suite, bobTokenBefore.Balances, bobTokenAfter.Balances)
	AssertBalancesEqual(suite, aliceTokenBefore.Balances, aliceTokenAfter.Balances)
	suite.Require().Len(colAfter.CollectionApprovals, len(colBefore.CollectionApprovals), "EvaluateCanTransfer leaked auto-deletion")

	// Calling EvaluateCanTransfer a SECOND time must STILL succeed (proving MaxNumTransfers tracker wasn't incremented).
	res2 := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     transfer,
		InitiatedBy:  bob,
	})
	suite.Require().True(res2.Allowed, "second EvaluateCanTransfer should still succeed since trackers must not mutate")
}

func (suite *TestSuite) TestEvaluateCanTransfer_AtTimestampBoundaryWithoutMutatingBlockTime() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	origBlockTime := suite.ctx.BlockTime()

	windowStart := uint64(origBlockTime.Add(1 * time.Hour).UnixMilli())
	windowEnd := uint64(origBlockTime.Add(2 * time.Hour).UnixMilli())

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
		{
			ApprovalId:        "windowed",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     []*types.UintRange{{Start: sdkmath.NewUint(windowStart), End: sdkmath.NewUint(windowEnd)}},
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
		},
	}

	suite.Require().NoError(CreateCollections(suite, wctx, collectionsToCreate))

	transfer := &types.Transfer{
		From:        bob,
		ToAddresses: []string{alice},
		Balances: []*types.Balance{{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		}},
	}

	// Before window (now) -> rejected with transferTime criterion.
	resBefore := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId:  sdkmath.NewUint(1),
		Transfer:      transfer,
		InitiatedBy:   bob,
		AtTimestampMs: windowStart - 1000,
	})
	suite.Require().False(resBefore.Allowed)
	suite.Require().True(hasFailedCriterion(resBefore.FailedCriteria, "transferTime"), "expected transferTime failure, got %+v", resBefore.FailedCriteria)

	// Inside window -> allowed!
	resInside := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId:  sdkmath.NewUint(1),
		Transfer:      transfer,
		InitiatedBy:   bob,
		AtTimestampMs: windowStart + 1000,
	})
	suite.Require().True(resInside.Allowed, "expected allowed inside window, got %+v", resInside.FailedCriteria)
	suite.Require().Contains(resInside.MatchedApprovalIds, "windowed")

	// After window -> rejected with transferTime criterion.
	resAfter := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId:  sdkmath.NewUint(1),
		Transfer:      transfer,
		InitiatedBy:   bob,
		AtTimestampMs: windowEnd + 1000,
	})
	suite.Require().False(resAfter.Allowed)
	suite.Require().True(hasFailedCriterion(resAfter.FailedCriteria, "transferTime"))

	// Suite block time must remain untouched.
	suite.Require().Equal(origBlockTime, suite.ctx.BlockTime())
}

func (suite *TestSuite) TestEvaluateCanTransfer_DiagnosticCategories() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Collection 1: base transferable collection with token 1 minted to bob.
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	hugeCoin := sdk.NewInt64Coin("ubadge", 999_999_999_999_999) // exceeds bob's 100 BADGE
	leafHash := sha256.Sum256([]byte("secret-code"))

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
		{
			ApprovalId:        "needs-coins",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MustPrioritize:                 true,
				CoinTransfers: []*types.CoinTransfer{{
					To:    charlie,
					Coins: []*sdk.Coin{&hugeCoin},
				}},
			},
		},
		{
			ApprovalId:        "needs-merkle",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MustPrioritize:                 true,
				MerkleChallenges: []*types.MerkleChallenge{{
					Root:                hex.EncodeToString(leafHash[:]),
					ExpectedProofLength: sdkmath.NewUint(0),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
					ChallengeTrackerId:  "mc-diag",
				}},
			},
		},
		{
			ApprovalId:        "needs-must-own",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MustPrioritize:                 true,
				MustOwnTokens: []*types.MustOwnTokens{{
					CollectionId:           sdkmath.NewUint(1),
					TokenIds:               GetTwoUintRanges(),
					OwnershipTimes:         GetFullUintRanges(),
					AmountRange:            &types.UintRange{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(100)},
					MustSatisfyForAllAssets: true,
				}},
			},
		},
		{
			ApprovalId:        "needs-user-approvals",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: false,
				OverridesToIncomingApprovals:   false,
				MustPrioritize:                 true,
			},
		},
	}

	suite.Require().NoError(CreateCollections(suite, wctx, collectionsToCreate))

	baseTransfer := func(approvalId string) *types.Transfer {
		return &types.Transfer{
			From:        bob,
			ToAddresses: []string{alice},
			Balances: []*types.Balance{{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			}},
			OnlyCheckPrioritizedCollectionApprovals: true,
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{
				ApprovalId:      approvalId,
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			}},
		}
	}

	// 1. CoinTransfer failure reports "coinTransfer" AND surfaces RequiredCoins.
	coinRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     baseTransfer("needs-coins"),
		InitiatedBy:  bob,
	})
	suite.Require().False(coinRes.Allowed)
	suite.Require().True(hasFailedCriterion(coinRes.FailedCriteria, "coinTransfer"), "got %+v", coinRes.FailedCriteria)
	suite.Require().Equal(hugeCoin.Amount, coinRes.RequiredCoins.AmountOf("ubadge"))

	// 2. MerkleChallenge failure reports "merkleChallenge" and does NOT leak unrelated coins from "needs-coins".
	merkleRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     baseTransfer("needs-merkle"),
		InitiatedBy:  bob,
	})
	suite.Require().False(merkleRes.Allowed)
	suite.Require().True(hasFailedCriterion(merkleRes.FailedCriteria, "merkleChallenge"), "got %+v", merkleRes.FailedCriteria)
	suite.Require().True(merkleRes.RequiredCoins.IsZero(), "prioritized-only merkle failure must not leak requiredCoins from unrelated approvals")

	// 3. MustOwnTokens failure reports "mustOwnTokens".
	ownRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     baseTransfer("needs-must-own"),
		InitiatedBy:  bob,
	})
	suite.Require().False(ownRes.Allowed)
	suite.Require().True(hasFailedCriterion(ownRes.FailedCriteria, "mustOwnTokens"), "got %+v", ownRes.FailedCriteria)

	// 4. User outgoing approval rejection reports "userOutgoingApproval".
	suite.Require().NoError(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                                   bob,
		CollectionId:                              sdkmath.NewUint(1),
		UpdateOutgoingApprovals:                   true,
		OutgoingApprovals:                         []*types.UserOutgoingApproval{},
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers: false,
	}))
	outRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     baseTransfer("needs-user-approvals"),
		InitiatedBy:  bob,
	})
	suite.Require().False(outRes.Allowed)
	suite.Require().True(hasFailedCriterion(outRes.FailedCriteria, "userOutgoingApproval"), "got %+v", outRes.FailedCriteria)

	// 5. User incoming approval rejection reports "userIncomingApproval".
	suite.Require().NoError(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                                  alice,
		CollectionId:                             sdkmath.NewUint(1),
		UpdateIncomingApprovals:                  true,
		IncomingApprovals:                        []*types.UserIncomingApproval{},
		UpdateAutoApproveAllIncomingTransfers:    true,
		AutoApproveAllIncomingTransfers:          false,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedIncomingTransfers: false,
	}))
	// Bob auto-approves self-initiated outgoing, so outgoing passes and incoming fails.
	suite.Require().NoError(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                                   bob,
		CollectionId:                              sdkmath.NewUint(1),
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers: true,
	}))
	incRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     baseTransfer("needs-user-approvals"),
		InitiatedBy:  bob,
	})
	suite.Require().False(incRes.Allowed)
	suite.Require().True(hasFailedCriterion(incRes.FailedCriteria, "userIncomingApproval"), "got %+v", incRes.FailedCriteria)

	// 6. OwnershipTimes mismatch reports "ownershipTimes".
	wrongTimeTransfer := baseTransfer("needs-coins")
	wrongTimeTransfer.Balances = []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetTwoUintRanges(), // approval only allows token ID 1
		OwnershipTimes: GetFullUintRanges(),
	}}
	otRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     wrongTimeTransfer,
		InitiatedBy:  bob,
	})
	suite.Require().False(otRes.Allowed)
	suite.Require().True(hasFailedCriterion(otRes.FailedCriteria, "ownershipTimes"), "got %+v", otRes.FailedCriteria)

	// 7. MaxNumTransfers exhausted after real transfer reports "maxNumTransfers".
	suite.Require().NoError(UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{{
			ApprovalId:        "once-only",
			FromListId:        bob,
			ToListId:          alice,
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MustPrioritize:                 true,
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1),
					AmountTrackerId:        "once-tracker",
				},
			},
		}},
	}))
	onceTransfer := baseTransfer("once-only")
	// Consume the 1 allowed transfer via real MsgTransferTokens.
	suite.Require().NoError(TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers:    []*types.Transfer{onceTransfer},
	}))
	maxRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     onceTransfer,
		InitiatedBy:  bob,
	})
	suite.Require().False(maxRes.Allowed)
	suite.Require().True(hasFailedCriterion(maxRes.FailedCriteria, "maxNumTransfers"), "got %+v", maxRes.FailedCriteria)

	// 8. Missing collection / nil transfer safety checks.
	missingRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(9999),
		Transfer:     baseTransfer("needs-coins"),
		InitiatedBy:  bob,
	})
	suite.Require().False(missingRes.Allowed)
	suite.Require().True(hasFailedCriterion(missingRes.FailedCriteria, "collectionNotFound"))

	nilRes := suite.app.TokenizationKeeper.EvaluateCanTransfer(suite.ctx, keeper.CanTransferRequest{
		CollectionId: sdkmath.NewUint(1),
		Transfer:     nil,
		InitiatedBy:  bob,
	})
	suite.Require().False(nilRes.Allowed)
	suite.Require().True(hasFailedCriterion(nilRes.FailedCriteria, "transfer"))
}

func hasFailedCriterion(failures []keeper.FailedCriterion, criterion string) bool {
	for _, f := range failures {
		if f.Criterion == criterion {
			return true
		}
	}
	return false
}
