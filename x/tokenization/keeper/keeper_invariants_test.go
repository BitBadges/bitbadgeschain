package keeper_test

import (
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	testutilkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TestStoreKeyUint64BoundsDoNotPanic verifies that collection and dynamic store
// lookup/delete operations guard uint64 bounds before encoding keys.
func TestStoreKeyUint64BoundsDoNotPanic(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	overMaxUint64 := sdkmath.NewUintFromBigInt(new(big.Int).Add(new(big.Int).SetUint64(math.MaxUint64), big.NewInt(1)))

	require.NotPanics(t, func() {
		_, found := k.GetCollectionFromStore(ctx, overMaxUint64)
		require.False(t, found)
		k.DeleteCollectionFromStore(ctx, overMaxUint64)
		require.Error(t, k.SetCollectionInStore(ctx, &types.TokenCollection{CollectionId: overMaxUint64}, false))

		_, found = k.GetDynamicStoreFromStore(ctx, overMaxUint64)
		require.False(t, found)
		k.DeleteDynamicStoreFromStore(ctx, overMaxUint64)
		require.Error(t, k.SetDynamicStoreInStore(ctx, types.DynamicStore{StoreId: overMaxUint64}))

		_, found = k.GetDynamicStoreValueFromStore(ctx, overMaxUint64, "bb1test")
		require.False(t, found)
		k.DeleteDynamicStoreValueFromStore(ctx, overMaxUint64, "bb1test")
		k.DeleteAllDynamicStoreValues(ctx, overMaxUint64)

		_, found = k.GetCollectionStatsFromStore(ctx, overMaxUint64)
		require.False(t, found)
		require.Error(t, k.SetCollectionStatsInStore(ctx, overMaxUint64, &types.CollectionStats{}))
	})
}

// TestExecuteEVMQueryRollsBackMutationsOnFallbackKeeper verifies that any KV
// writes attempted inside ExecuteEVMQuery are discarded by the atomic context
// while gas consumption is preserved on the outer context.
func TestExecuteEVMQueryRollsBackMutationsOnFallbackKeeper(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(500_000))

	k.SetEVMKeeper(&contextQueryKeeper{call: func(child sdk.Context, _ *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
		// Write a dynamic store inside the query child context; atomic rollback must revert it.
		err := k.SetDynamicStoreInStore(child, types.DynamicStore{
			StoreId:   sdkmath.NewUint(77),
			CreatedBy: "bb1creator",
		})
		require.NoError(t, err)
		return &evmtypes.MsgEthereumTxResponse{GasUsed: 1234, Ret: []byte{1}}, nil
	}})

	ret, err := k.ExecuteEVMQuery(ctx, gasMeteringContract, []byte{0x01}, 200_000)
	require.NoError(t, err)
	require.Equal(t, []byte{1}, ret)
	require.GreaterOrEqual(t, ctx.GasMeter().GasConsumed(), uint64(1234))

	_, found := k.GetDynamicStoreFromStore(ctx, sdkmath.NewUint(77))
	require.False(t, found, "writes inside ExecuteEVMQuery must be rolled back")
}

// TestGetBalanceOrApplyDefaultReadOnlyDoesNotMutateStore verifies that read-only
// default balance resolution does not persist a balance or increment approval versions.
func TestGetBalanceOrApplyDefaultReadOnlyDoesNotMutateStore(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	addr := sdk.AccAddress([]byte("readonly-balance-usr")).String()
	col := &types.TokenCollection{
		CollectionId: sdkmath.OneUint(),
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{{
				Amount:         sdkmath.NewUint(5),
				TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.NewUint(math.MaxUint64)}},
			}},
			OutgoingApprovals: []*types.UserOutgoingApproval{{
				ApprovalId: "default-out",
				ToListId:   "All",
				Version:    sdkmath.ZeroUint(),
			}},
		},
	}

	bal, appliedDefault, err := k.GetBalanceOrApplyDefaultReadOnly(ctx, col, addr)
	require.NoError(t, err)
	require.True(t, appliedDefault)
	require.Len(t, bal.Balances, 1)

	// Store must remain untouched and version must not have been bumped.
	_, stored := k.GetUserBalanceFromStore(ctx, keeper.ConstructBalanceKey(addr, col.CollectionId))
	require.False(t, stored)
	vKey := keeper.ConstructApprovalVersionKey(col.CollectionId, "outgoing", addr, "default-out")
	v, _ := k.GetApprovalTrackerVersionFromStore(ctx, vKey)
	require.True(t, v.IsZero())
}

// TestDeleteDynamicStoreRemovesAllValues verifies that deleting a dynamic store
// also clears every stored DynamicStoreValue entry belonging to that store.
func (suite *TestSuite) TestDeleteDynamicStoreRemovesAllValues() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	server := keeper.NewMsgServerImpl(k)

	createResp, err := server.CreateDynamicStore(wctx, &types.MsgCreateDynamicStore{
		Creator:      alice,
		DefaultValue: false,
	})
	suite.Require().NoError(err)
	storeID := createResp.StoreId

	for _, addr := range []string{alice, bob, charlie} {
		_, err = server.SetDynamicStoreValue(wctx, &types.MsgSetDynamicStoreValue{
			Creator: alice,
			StoreId: storeID,
			Address: addr,
			Value:   true,
		})
		suite.Require().NoError(err)
	}

	_, err = server.DeleteDynamicStore(wctx, &types.MsgDeleteDynamicStore{
		Creator: alice,
		StoreId: storeID,
	})
	suite.Require().NoError(err)

	_, found := k.GetDynamicStoreFromStore(suite.ctx, storeID)
	suite.Require().False(found)
	for _, addr := range []string{alice, bob, charlie} {
		_, valFound := k.GetDynamicStoreValueFromStore(suite.ctx, storeID, addr)
		suite.Require().False(valFound, "dynamic store value for %s must be removed on delete", addr)
	}
}

// TestPurgeApprovalsNoOpAndBoundaryTime verifies that MsgPurgeApprovals is a
// no-op (no state mutation or events) when nothing qualifies and treats
// TransferTimes.End == currentTime as active (not yet expired).
func (suite *TestSuite) TestPurgeApprovalsNoOpAndBoundaryTime() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	server := keeper.NewMsgServerImpl(k)

	now := uint64(suite.ctx.BlockTime().UnixMilli())
	resp, err := server.UniversalUpdateCollection(wctx, &types.MsgUniversalUpdateCollection{
		Creator:               alice,
		CollectionId:          sdkmath.ZeroUint(),
		UpdateValidTokenIds:   true,
		ValidTokenIds:         GetFullUintRanges(),
		CollectionPermissions: &types.CollectionPermissions{},
	})
	suite.Require().NoError(err)

	// Set a user outgoing approval on alice ending at exact current block time `now`.
	_, err = server.SetOutgoingApproval(wctx, &types.MsgSetOutgoingApproval{
		Creator:      alice,
		CollectionId: resp.CollectionId,
		Approval: &types.UserOutgoingApproval{
			ApprovalId:        "boundary-approval",
			ToListId:          "All",
			InitiatedByListId: "All",
			TransferTimes:     []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.NewUint(now)}},
			TokenIds:          GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			Version:           sdkmath.ZeroUint(),
			ApprovalCriteria: &types.OutgoingApprovalCriteria{
				AutoDeletionOptions: &types.AutoDeletionOptions{AllowPurgeIfExpired: true},
			},
		},
	})
	suite.Require().NoError(err)

	// At exact boundary End == now, approval is still active so purge must be a no-op (0 purged, no events).
	emCtx := suite.ctx.WithEventManager(sdk.NewEventManager())
	purgeResp, err := server.PurgeApprovals(sdk.WrapSDKContext(emCtx), &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               resp.CollectionId,
		ApproverAddress:            alice,
		PurgeCounterpartyApprovals: true,
		PurgeExpired:               true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalLevel:   "outgoing",
			ApproverAddress: alice,
			ApprovalId:      "boundary-approval",
			Version:         sdkmath.ZeroUint(),
		}},
	})
	suite.Require().NoError(err)
	suite.Require().True(purgeResp.NumPurged.IsZero())
	suite.Require().Empty(emCtx.EventManager().Events(), "no-op purge must not emit events")

	// Advance 1ms past End: approval is now expired and purge removes it.
	laterCtx := suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Millisecond))
	purgeResp, err = server.PurgeApprovals(sdk.WrapSDKContext(laterCtx), &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               resp.CollectionId,
		ApproverAddress:            alice,
		PurgeCounterpartyApprovals: true,
		PurgeExpired:               true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalLevel:   "outgoing",
			ApproverAddress: alice,
			ApprovalId:      "boundary-approval",
			Version:         sdkmath.ZeroUint(),
		}},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.OneUint(), purgeResp.NumPurged)
}

// TestPoolWithdrawalWorksWithLockedDefaultUserPermissions verifies that pool
// withdrawals succeed even when DefaultBalances.UserPermissions forbids users
// from modifying their own outgoing approvals.
func (suite *TestSuite) TestPoolWithdrawalWorksWithLockedDefaultUserPermissions() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	server := keeper.NewMsgServerImpl(k)

	full := GetFullUintRanges()
	lockedPerms := &types.UserPermissions{
		CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{{
			ApprovalId:                "All",
			ToListId:                  "All",
			InitiatedByListId:         "All",
			TransferTimes:             full,
			TokenIds:                  full,
			OwnershipTimes:            full,
			PermanentlyForbiddenTimes: full,
		}},
	}
	resp, err := server.UniversalUpdateCollection(wctx, &types.MsgUniversalUpdateCollection{
		Creator:               alice,
		CollectionId:          sdkmath.ZeroUint(),
		UpdateValidTokenIds:   true,
		ValidTokenIds:         []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
		CollectionPermissions: &types.CollectionPermissions{},
		DefaultBalances: &types.UserBalanceStore{
			UserPermissions: lockedPerms,
		},
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "mint",
				FromListId:        types.MintAddress,
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TransferTimes:     full,
				TokenIds:          []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
				OwnershipTimes:    full,
				Version:           sdkmath.ZeroUint(),
				ApprovalCriteria:  &types.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
			},
			{
				ApprovalId:        "transfer",
				FromListId:        "AllWithoutMint",
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TransferTimes:     full,
				TokenIds:          []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
				OwnershipTimes:    full,
				Version:           sdkmath.ZeroUint(),
				ApprovalCriteria:  &types.ApprovalCriteria{OverridesToIncomingApprovals: true},
			},
		},
		AliasPathsToAdd: []*types.AliasPathAddObject{{
			Denom:  "lpunit",
			Symbol: "LPU",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{Amount: sdkmath.OneUint()},
				SideB: []*types.Balance{{
					Amount:         sdkmath.OneUint(),
					TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
					OwnershipTimes: full,
				}},
			},
			DenomUnits: []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "lpunit", IsDefaultDisplay: true}},
		}},
		Invariants: &types.InvariantsAddObject{NoCustomOwnershipTimes: true},
	})
	suite.Require().NoError(err)

	poolAddr := charlie
	suite.Require().NoError(k.SetReservedProtocolAddressInStore(suite.ctx, poolAddr, true))
	k.SetPoolAddressInCache(suite.ctx, poolAddr, 1)

	// Mint tokens directly to poolAddr
	_, err = server.TransferTokens(wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: resp.CollectionId,
		Transfers: []*types.Transfer{{
			From:        types.MintAddress,
			ToAddresses: []string{poolAddr},
			Balances: []*types.Balance{{
				Amount:         sdkmath.NewUint(100),
				TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
				OwnershipTimes: full,
			}},
		}},
	})
	suite.Require().NoError(err)

	collection, found := k.GetCollectionFromStore(suite.ctx, resp.CollectionId)
	suite.Require().True(found)
	aliasDenom := keeper.AliasDenomPrefix + collection.CollectionId.String() + ":" + collection.AliasPaths[0].Denom

	// Pool withdrawal must succeed despite locked DefaultBalances.UserPermissions
	err = k.SendNativeTokensFromAddressWithPoolApprovals(suite.ctx, poolAddr, bob, aliasDenom, sdkmath.NewUint(25))
	suite.Require().NoError(err)
}

// TestBackedCollectionMaxSupplyPerIdInvariant enforces MaxSupplyPerId on
// cosmos-coin-backed collections as well as standard collections.
func (suite *TestSuite) TestBackedCollectionMaxSupplyPerIdInvariant() {
	k := suite.app.TokenizationKeeper
	full := GetFullUintRanges()
	collectionID := sdkmath.NewUint(99)
	collection := &types.TokenCollection{
		CollectionId:  collectionID,
		ValidTokenIds: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
		Invariants: &types.CollectionInvariants{
			NoCustomOwnershipTimes: true,
			MaxSupplyPerId:         sdkmath.NewUint(50),
			CosmosCoinBackedPath: &types.CosmosCoinBackedPath{
				Address: alice,
				Conversion: &types.Conversion{
					SideA: &types.ConversionSideAWithDenom{
						Denom:  "ubadge",
						Amount: sdkmath.OneUint(),
					},
					SideB: []*types.Balance{{
						Amount:         sdkmath.OneUint(),
						TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
						OwnershipTimes: full,
					}},
				},
			},
		},
	}

	// Setting CollectionStats.Balances > MaxSupplyPerId (50) must be caught by post-transfer invariants.
	suite.Require().NoError(k.SetCollectionStatsInStore(suite.ctx, collectionID, &types.CollectionStats{
		Balances: []*types.Balance{{
			Amount:         sdkmath.NewUint(60),
			TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
			OwnershipTimes: full,
		}},
	}))
	suite.Require().Error(k.CheckPostTransferInvariants(suite.ctx, collection, types.MintAddress, []string{alice}, alice))

	// Within MaxSupplyPerId (40 <= 50) passes.
	suite.Require().NoError(k.SetCollectionStatsInStore(suite.ctx, collectionID, &types.CollectionStats{
		Balances: []*types.Balance{{
			Amount:         sdkmath.NewUint(40),
			TokenIds:       []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
			OwnershipTimes: full,
		}},
	}))
	suite.Require().NoError(k.CheckPostTransferInvariants(suite.ctx, collection, types.MintAddress, []string{alice}, alice))
}

// TestICQAddressCanonicalization verifies that ICQ queries canonicalize uppercase
// bech32 addresses and reject oversized hex strings.
func (suite *TestSuite) TestICQAddressCanonicalization() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	server := keeper.NewMsgServerImpl(k)

	full := GetFullUintRanges()
	resp, err := server.UniversalUpdateCollection(wctx, &types.MsgUniversalUpdateCollection{
		Creator:                   alice,
		CollectionId:              sdkmath.ZeroUint(),
		UpdateValidTokenIds:       true,
		ValidTokenIds:             []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}},
		CollectionPermissions:     &types.CollectionPermissions{},
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{{
			ApprovalId: "mint", FromListId: types.MintAddress, ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
			TransferTimes: full, TokenIds: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}}, OwnershipTimes: full,
			Version: sdkmath.ZeroUint(), ApprovalCriteria: &types.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
		}},
	})
	suite.Require().NoError(err)

	_, err = server.TransferTokens(wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: resp.CollectionId,
		Transfers: []*types.Transfer{{
			From: types.MintAddress, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(7), TokenIds: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}}, OwnershipTimes: full}},
		}},
	})
	suite.Require().NoError(err)

	// Uppercase bech32 resolves to the same balance as canonical lowercase bech32
	qResp := k.ProcessOwnershipQuery(suite.ctx, &types.OwnershipQueryPacket{
		QueryId:       "q-upper",
		CollectionId:  resp.CollectionId.String(),
		Address:       strings.ToUpper(alice),
		TokenId:       "1",
		OwnershipTime: "1",
	})
	suite.Require().Empty(qResp.Error)
	suite.Require().True(qResp.OwnsTokens)
	suite.Require().Equal(sdkmath.NewUint(7), qResp.TotalAmount)

	// Oversized 0x hex (> 40 hex chars) must not truncate into alice's 20-byte address
	aliceHex := common.BytesToAddress(sdk.MustAccAddressFromBech32(alice)).Hex()
	qOversized := k.ProcessOwnershipQuery(suite.ctx, &types.OwnershipQueryPacket{
		QueryId:       "q-longhex",
		CollectionId:  resp.CollectionId.String(),
		Address:       aliceHex + "00",
		TokenId:       "1",
		OwnershipTime: "1",
	})
	suite.Require().False(qOversized.OwnsTokens)
}

// TestParamsValidationRejectsExcessiveAffiliatePercentage checks Params.Validate
// bounds on AffiliatePercentage and AllowedDenoms.
func TestParamsValidationRejectsExcessiveAffiliatePercentage(t *testing.T) {
	p := types.DefaultParams()
	require.NoError(t, p.Validate())

	p.AffiliatePercentage = sdkmath.NewUint(10001)
	require.Error(t, p.Validate())

	p = types.DefaultParams()
	p.AllowedDenoms = []string{"bad denom!"}
	require.Error(t, p.Validate())
}

// Compile-time check for statedb import in package keeper_test
var _ *statedb.StateDB
