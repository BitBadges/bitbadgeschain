package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/council/keeper"
	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

// ---------------------------------------------------------------------------
// Mock TokenizationKeeper
// ---------------------------------------------------------------------------

type mockTokenizationKeeper struct {
	balances    map[string]uint64 // key: "collectionId/tokenId/address"
	totalSupply map[string]uint64 // key: "collectionId/tokenId"
}

func newMockTokenizationKeeper() *mockTokenizationKeeper {
	return &mockTokenizationKeeper{
		balances:    make(map[string]uint64),
		totalSupply: make(map[string]uint64),
	}
}

func mockBalanceKey(collectionId, tokenId uint64, address string) string {
	return fmt.Sprintf("%d/%d/%s", collectionId, tokenId, address)
}

func mockSupplyKey(collectionId, tokenId uint64) string {
	return fmt.Sprintf("%d/%d", collectionId, tokenId)
}

func (m *mockTokenizationKeeper) SetBalance(collectionId, tokenId uint64, address string, balance uint64) {
	key := mockBalanceKey(collectionId, tokenId, address)
	m.balances[key] = balance
}

func (m *mockTokenizationKeeper) SetTotalSupply(collectionId, tokenId uint64, supply uint64) {
	key := mockSupplyKey(collectionId, tokenId)
	m.totalSupply[key] = supply
}

func (m *mockTokenizationKeeper) GetCredentialBalance(_ sdk.Context, collectionId, tokenId uint64, address string) (uint64, error) {
	key := mockBalanceKey(collectionId, tokenId, address)
	return m.balances[key], nil
}

func (m *mockTokenizationKeeper) GetTotalSupply(_ sdk.Context, collectionId, tokenId uint64) (uint64, error) {
	key := mockSupplyKey(collectionId, tokenId)
	return m.totalSupply[key], nil
}

// ---------------------------------------------------------------------------
// Mock MsgRouter (no-op for most tests; execute tests check it's called)
// ---------------------------------------------------------------------------

type mockMsgRouter struct {
	dispatched []sdk.Msg
	failOn     int // if >= 0, fail on the nth dispatch
}

func newMockMsgRouter() *mockMsgRouter {
	return &mockMsgRouter{failOn: -1}
}

func (r *mockMsgRouter) Handler(msg sdk.Msg) keeper.MsgHandler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		r.dispatched = append(r.dispatched, msg)
		return &sdk.Result{}, nil
	}
}

// ---------------------------------------------------------------------------
// Test setup helper
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockTokenizationKeeper) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	mockTK := newMockTokenizationKeeper()
	mockRouter := newMockMsgRouter()

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		log.NewNopLogger(),
		mockTK,
		mockRouter,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).
		WithBlockTime(time.Unix(1000, 0)) // block time = 1000s = 1_000_000 ms

	return k, ctx, mockTK
}

func setupMsgServer(t *testing.T) (keeper.Keeper, sdk.Context, *mockTokenizationKeeper) {
	return setupKeeper(t)
}

const (
	testCreator  = "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"
	testVoter1   = "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"
	testVoter2   = "cosmos1p8s0p6gqc6c9gt77lg6ayjss9fkwqphprka7mr"
	testNoHolder = "cosmos1x8fhpj9nmhqk8n9kx8kvce29fy4ux5jqhpk30j"

	collectionId uint64 = 1
	tokenId      uint64 = 1
)

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateCouncilSuccess(t *testing.T) {
	k, ctx, _ := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	msg := &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         60000, // 60 seconds
		AllowedMsgTypes:        nil,
	}

	resp, err := srv.CreateCouncil(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.CouncilId)

	// Verify stored
	council, found := k.GetCouncil(ctx, 1)
	require.True(t, found)
	require.Equal(t, testCreator, council.Creator)
	require.Equal(t, uint64(51), council.VotingThreshold)
	require.NotEmpty(t, council.AccountAddress)
}

func TestProposeWithCredential(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testCreator, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 1000)

	// Create council first
	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	// Propose
	resp, err := srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testCreator,
		CouncilId:   1,
		MsgTypeUrls: []string{"/cosmos.bank.v1beta1.MsgSend"},
		MsgBytes:    [][]byte{[]byte("test")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000, // +1 day
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.ProposalId)

	proposal, found := k.GetProposal(ctx, 1, 1)
	require.True(t, found)
	require.Equal(t, types.ProposalStatusPending, proposal.Status)
}

func TestProposeWithoutCredential(t *testing.T) {
	k, ctx, _ := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	// Create council (no balance set for testCreator means 0)
	// But we need the creator to NOT need credentials to create council.
	// First set balance so we can create, then test proposing without credentials.
	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	// Propose without holding credential (balance = 0)
	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testNoHolder,
		CouncilId:   1,
		MsgTypeUrls: []string{"/cosmos.bank.v1beta1.MsgSend"},
		MsgBytes:    [][]byte{[]byte("test")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "no credential balance")
}

func TestVoteWithCredentialWeightEqualsBalance(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 250)
	mockTK.SetTotalSupply(collectionId, tokenId, 1000)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter:      testVoter1,
		CouncilId:  1,
		ProposalId: 1,
		VoteYes:    true,
	})
	require.NoError(t, err)

	vote, found := k.GetVote(ctx, 1, 1, testVoter1)
	require.True(t, found)
	require.Equal(t, uint64(250), vote.Weight)

	proposal, found := k.GetProposal(ctx, 1, 1)
	require.True(t, found)
	require.Equal(t, uint64(250), proposal.YesWeight)
}

func TestVoteWithoutCredential(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 1000)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Vote without credential (testNoHolder has 0 balance)
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter:      testNoHolder,
		CouncilId:  1,
		ProposalId: 1,
		VoteYes:    true,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "no credential balance")
}

func TestVoteReachesThreshold(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	// Total supply = 100, threshold = 51%, so 51 yes weight needed
	mockTK.SetBalance(collectionId, tokenId, testVoter1, 30)
	mockTK.SetBalance(collectionId, tokenId, testVoter2, 25)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         60000,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Voter1 votes yes with weight 30 — not enough
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)
	proposal, _ := k.GetProposal(ctx, 1, 1)
	require.Equal(t, types.ProposalStatusPending, proposal.Status)

	// Voter2 votes yes with weight 25 — total 55 >= 51 → passed
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter2, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)
	proposal, _ = k.GetProposal(ctx, 1, 1)
	require.Equal(t, types.ProposalStatusPassed, proposal.Status)
	require.NotZero(t, proposal.PassedAt)
}

func TestExecuteBeforeDelay(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        50,
		ExecutionDelay:         120000, // 120 seconds
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Vote to pass
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)

	// Execute immediately — should fail (delay not met)
	_, err = srv.ExecuteProposal(ctx, &types.MsgExecuteProposal{
		Sender: testVoter1, CouncilId: 1, ProposalId: 1,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "execution delay has not elapsed")
}

func TestExecuteAfterDelay(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        50,
		ExecutionDelay:         60000, // 60 seconds
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Vote to pass
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)

	// Advance time past execution delay (60s)
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(120 * time.Second))

	_, err = srv.ExecuteProposal(ctx, &types.MsgExecuteProposal{
		Sender: testVoter1, CouncilId: 1, ProposalId: 1,
	})
	require.NoError(t, err)

	proposal, _ := k.GetProposal(ctx, 1, 1)
	require.Equal(t, types.ProposalStatusExecuted, proposal.Status)
}

func TestExecuteWithWrongStatus(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 10)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Vote yes with only 10/100 — threshold 51% not met, still pending
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)

	// Try to execute a pending proposal — should fail
	_, err = srv.ExecuteProposal(ctx, &types.MsgExecuteProposal{
		Sender: testVoter1, CouncilId: 1, ProposalId: 1,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "proposal has not passed")
}

func TestExecuteAlreadyExecuted(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        50,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)

	// First execute — should succeed
	_, err = srv.ExecuteProposal(ctx, &types.MsgExecuteProposal{
		Sender: testVoter1, CouncilId: 1, ProposalId: 1,
	})
	require.NoError(t, err)

	// Second execute — should fail
	_, err = srv.ExecuteProposal(ctx, &types.MsgExecuteProposal{
		Sender: testVoter1, CouncilId: 1, ProposalId: 1,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "already executed")
}

func TestProposeDisallowedMsgType(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 100)
	mockTK.SetTotalSupply(collectionId, tokenId, 1000)

	// Council only allows MsgUpdateParams
	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
		AllowedMsgTypes:        []string{"/cosmos.gov.v1.MsgUpdateParams"},
	})
	require.NoError(t, err)

	// Propose with disallowed type
	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/cosmos.bank.v1beta1.MsgSend"},
		MsgBytes:    [][]byte{[]byte("test")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "not allowed")
}

func TestRevote(t *testing.T) {
	k, ctx, mockTK := setupMsgServer(t)
	srv := keeper.NewMsgServerImpl(k)

	mockTK.SetBalance(collectionId, tokenId, testVoter1, 60)
	mockTK.SetTotalSupply(collectionId, tokenId, 100)

	_, err := srv.CreateCouncil(ctx, &types.MsgCreateCouncil{
		Creator:                testCreator,
		CredentialCollectionId: collectionId,
		CredentialTokenId:      tokenId,
		VotingThreshold:        51,
		ExecutionDelay:         0,
	})
	require.NoError(t, err)

	_, err = srv.Propose(ctx, &types.MsgPropose{
		Proposer:    testVoter1,
		CouncilId:   1,
		MsgTypeUrls: []string{"/test.MsgFoo"},
		MsgBytes:    [][]byte{[]byte("foo")},
		Deadline:    ctx.BlockTime().UnixMilli() + 86400000,
	})
	require.NoError(t, err)

	// Vote no first
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: false,
	})
	require.NoError(t, err)

	proposal, _ := k.GetProposal(ctx, 1, 1)
	require.Equal(t, uint64(0), proposal.YesWeight)
	require.Equal(t, uint64(60), proposal.NoWeight)

	// Re-vote yes — should switch
	_, err = srv.Vote(ctx, &types.MsgVote{
		Voter: testVoter1, CouncilId: 1, ProposalId: 1, VoteYes: true,
	})
	require.NoError(t, err)

	proposal, _ = k.GetProposal(ctx, 1, 1)
	require.Equal(t, uint64(60), proposal.YesWeight)
	require.Equal(t, uint64(0), proposal.NoWeight)
	// 60% >= 51% → passed
	require.Equal(t, types.ProposalStatusPassed, proposal.Status)
}
