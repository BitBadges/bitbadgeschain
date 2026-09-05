package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestMigrateV35OwnerValues() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{OverridesFromOutgoingApprovals: true}, GetFullUintRanges())}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	collection.Manager = strings.ToUpper(bob)
	collection.CreatedBy = strings.ToUpper(bob)
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	suite.Require().NoError(k.SetDynamicStoreInStore(suite.ctx, types.DynamicStore{StoreId: sdkmath.OneUint(), CreatedBy: strings.ToUpper(bob), GlobalEnabled: true}))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	collection, err = GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	suite.Require().NoError(k.UniversalValidate(suite.ctx, collection, keeper.UniversalValidationParams{Creator: bob, MustBeManager: true}))
	suite.Require().Error(k.UniversalValidate(suite.ctx, collection, keeper.UniversalValidationParams{Creator: alice, MustBeManager: true}))
	suite.Require().Equal(bob, collection.CreatedBy)
	server := keeper.NewMsgServerImpl(k)
	_, err = server.UpdateDynamicStore(wctx, &types.MsgUpdateDynamicStore{Creator: bob, StoreId: sdkmath.OneUint(), GlobalEnabled: true})
	suite.Require().NoError(err)
	_, err = server.DeleteDynamicStore(wctx, &types.MsgDeleteDynamicStore{Creator: alice, StoreId: sdkmath.OneUint()})
	suite.Require().Error(err)
	_, err = server.DeleteDynamicStore(wctx, &types.MsgDeleteDynamicStore{Creator: bob, StoreId: sdkmath.OneUint()})
	suite.Require().NoError(err)
}

func (suite *TestSuite) TestMigrateV35PartyValues() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true}, GetFullUintRanges())}))
	suite.Require().NoError(TransferTokens(suite, wctx, v35Claim(alice, alice, nil)))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	collection.CollectionApprovals[0].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{{CollectionId: sdkmath.OneUint(), AmountRange: GetOneUintRange()[0], OwnershipTimes: GetFullUintRanges(), TokenIds: GetOneUintRange(), OwnershipCheckParty: strings.ToUpper(bob)}}
	collection.CollectionApprovals[0].ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{{StoreId: sdkmath.OneUint(), OwnershipCheckParty: strings.ToUpper(bob)}}
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	suite.Require().NoError(k.SetDynamicStoreInStore(suite.ctx, types.DynamicStore{StoreId: sdkmath.OneUint(), CreatedBy: bob, GlobalEnabled: true}))
	suite.Require().NoError(k.SetDynamicStoreValueInStore(suite.ctx, sdkmath.OneUint(), alice, true))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	collection, err = GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	approval := collection.CollectionApprovals[0]
	checked := 0
	for _, checker := range k.GetApprovalCriteriaCheckers(approval) {
		if checker.Name() != "MustOwnTokens" && checker.Name() != "DynamicStoreChallenges" {
			continue
		}
		_, err := checker.Check(suite.ctx, approval, collection, alice, "Mint", alice, "collection", "", nil, nil, "", true)
		suite.Require().Error(err, "%s must check bob, not the funded initiator", checker.Name())
		checked++
	}
	suite.Require().Equal(2, checked)
	rawStore := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	aliceKey := append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(alice, sdkmath.OneUint()))...)
	bobKey := append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(bob, sdkmath.OneUint()))...)
	rawStore.Set(bobKey, rawStore.Get(aliceKey))
	suite.Require().NoError(k.SetDynamicStoreValueInStore(suite.ctx, sdkmath.OneUint(), bob, true))
	for _, checker := range k.GetApprovalCriteriaCheckers(approval) {
		if checker.Name() != "MustOwnTokens" && checker.Name() != "DynamicStoreChallenges" {
			continue
		}
		_, err := checker.Check(suite.ctx, approval, collection, alice, "Mint", alice, "collection", "", nil, nil, "", true)
		suite.Require().NoError(err, "%s accepts the designated party's state", checker.Name())
	}
}

func (suite *TestSuite) TestMigrateV35ApprovalAddressValues() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{OverridesFromOutgoingApprovals: true}, GetFullUintRanges())}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	upper := strings.ToUpper(bob)
	a := collection.CollectionApprovals[0]
	a.ToListId = "!(AllWithout" + upper + ":Mint)"
	a.InitiatedByListId = upper + ":" + alice
	a.ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{{To: upper}}
	a.ApprovalCriteria.UserApprovalSettings = &types.UserApprovalSettings{UserRoyalties: &types.UserRoyalties{PayoutAddress: upper, Percentage: sdkmath.OneUint()}}
	a.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{{ProposalId: "proposal", Voters: []*types.Voter{{Address: upper, Weight: sdkmath.OneUint()}, {Address: bob, Weight: sdkmath.NewUint(2)}}, QuorumThreshold: sdkmath.OneUint(), DelayAfterQuorum: sdkmath.ZeroUint()}}
	votingBefore := suite.app.AppCodec().MustMarshal(a.ApprovalCriteria.VotingChallenges[0])
	collection.CollectionPermissions.CanUpdateCollectionApprovals = []*types.CollectionApprovalPermission{{FromListId: upper, ToListId: "All", InitiatedByListId: "NamedGroup"}}
	collection.DefaultBalances = &types.UserBalanceStore{
		IncomingApprovals: []*types.UserIncomingApproval{{FromListId: "!" + upper, InitiatedByListId: "Mint", Version: sdkmath.ZeroUint(), ApprovalCriteria: &types.IncomingApprovalCriteria{MustOwnTokens: []*types.MustOwnTokens{{CollectionId: sdkmath.OneUint(), OwnershipCheckParty: upper}}}}},
		OutgoingApprovals: []*types.UserOutgoingApproval{{ToListId: upper, InitiatedByListId: "All", Version: sdkmath.ZeroUint(), ApprovalCriteria: &types.OutgoingApprovalCriteria{DynamicStoreChallenges: []*types.DynamicStoreChallenge{{StoreId: sdkmath.OneUint(), OwnershipCheckParty: upper}}}}},
		UserPermissions:   &types.UserPermissions{CanUpdateIncomingApprovals: []*types.UserIncomingApprovalPermission{{FromListId: upper}}, CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{{ToListId: upper}}},
	}
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	rawStore := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	key := append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(alice, sdkmath.OneUint()))...)
	rawStore.Set(key, suite.app.AppCodec().MustMarshal(collection.DefaultBalances))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	collection, err = GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	a = collection.CollectionApprovals[0]
	suite.Require().Equal("!(AllWithout"+bob+":Mint)", a.ToListId)
	suite.Require().Equal(bob+":"+alice, a.InitiatedByListId)
	suite.Require().Equal(bob, a.ApprovalCriteria.CoinTransfers[0].To)
	suite.Require().Equal(bob, a.ApprovalCriteria.UserApprovalSettings.UserRoyalties.PayoutAddress)
	suite.Require().Equal(votingBefore, suite.app.AppCodec().MustMarshal(a.ApprovalCriteria.VotingChallenges[0]))
	suite.Require().Equal(bob, collection.CollectionPermissions.CanUpdateCollectionApprovals[0].FromListId)
	suite.Require().Equal("NamedGroup", collection.CollectionPermissions.CanUpdateCollectionApprovals[0].InitiatedByListId)
	var stored types.UserBalanceStore
	suite.app.AppCodec().MustUnmarshal(rawStore.Get(key), &stored)
	for _, b := range []*types.UserBalanceStore{collection.DefaultBalances, &stored} {
		suite.Require().Equal("!"+bob, b.IncomingApprovals[0].FromListId)
		suite.Require().Equal("Mint", b.IncomingApprovals[0].InitiatedByListId)
		suite.Require().Equal(bob, b.IncomingApprovals[0].ApprovalCriteria.MustOwnTokens[0].OwnershipCheckParty)
		suite.Require().Equal(bob, b.OutgoingApprovals[0].ApprovalCriteria.DynamicStoreChallenges[0].OwnershipCheckParty)
		suite.Require().Equal(bob, b.UserPermissions.CanUpdateIncomingApprovals[0].FromListId)
		suite.Require().Equal(bob, b.UserPermissions.CanUpdateOutgoingApprovals[0].ToListId)
	}
	before := rawStore.Get(key)
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	suite.Require().Equal(before, rawStore.Get(key), "migration is idempotent")
}

func (suite *TestSuite) TestMigrateV35PreservesLegacyVotingEntries() {
	k := suite.app.TokenizationKeeper
	upper := strings.ToUpper(bob)
	oldKey := keeper.ConstructVotingTrackerKey(sdkmath.OneUint(), upper, "incoming", "approval", "proposal", upper)
	newKey := keeper.ConstructVotingTrackerKey(sdkmath.OneUint(), bob, "incoming", "approval", "proposal", bob)
	vote := &types.VoteProof{ProposalId: "proposal", Voter: upper, YesWeight: sdkmath.NewUint(100), VotedAt: sdkmath.OneUint()}
	suite.Require().NoError(k.SetVoteInStore(suite.ctx, oldKey, vote))
	vote.Voter = bob
	vote.YesWeight = sdkmath.ZeroUint()
	suite.Require().NoError(k.SetVoteInStore(suite.ctx, newKey, vote))
	oldTracker := keeper.ConstructVotingChallengeTrackerKey(sdkmath.OneUint(), upper, "incoming", "approval", "proposal")
	newTracker := keeper.ConstructVotingChallengeTrackerKey(sdkmath.OneUint(), bob, "incoming", "approval", "proposal")
	suite.Require().NoError(k.SetVotingChallengeTrackerInStore(suite.ctx, oldTracker, &types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.OneUint()}))
	suite.Require().NoError(k.SetVotingChallengeTrackerInStore(suite.ctx, newTracker, &types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.NewUint(2)}))
	rawStore := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	keys := [][]byte{
		append(append([]byte{}, keeper.VotingTrackerKey...), []byte(oldKey)...),
		append(append([]byte{}, keeper.VotingTrackerKey...), []byte(newKey)...),
		append(append([]byte{}, keeper.VotingChallengeTrackerKey...), []byte(oldTracker)...),
		append(append([]byte{}, keeper.VotingChallengeTrackerKey...), []byte(newTracker)...),
	}
	before := make([][]byte, len(keys))
	for i, key := range keys {
		before[i] = append([]byte{}, rawStore.Get(key)...)
	}
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	for i, key := range keys {
		suite.Require().Equal(before[i], rawStore.Get(key))
	}
}
