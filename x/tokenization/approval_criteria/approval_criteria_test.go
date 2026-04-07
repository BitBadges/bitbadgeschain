package approval_criteria

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock services
// ============================================================

type mockAddressCheckService struct {
	isEvmContract      map[string]bool
	isLiquidityPool    map[string]bool
	isReservedProtocol map[string]bool
	evmContractErr     error
	liquidityPoolErr   error
}

func newMockAddressCheckService() *mockAddressCheckService {
	return &mockAddressCheckService{
		isEvmContract:      make(map[string]bool),
		isLiquidityPool:    make(map[string]bool),
		isReservedProtocol: make(map[string]bool),
	}
}

func (m *mockAddressCheckService) IsEVMContract(ctx sdk.Context, address string) (bool, error) {
	if m.evmContractErr != nil {
		return false, m.evmContractErr
	}
	return m.isEvmContract[address], nil
}

func (m *mockAddressCheckService) IsLiquidityPool(ctx sdk.Context, address string) (bool, error) {
	if m.liquidityPoolErr != nil {
		return false, m.liquidityPoolErr
	}
	return m.isLiquidityPool[address], nil
}

func (m *mockAddressCheckService) IsAddressReservedProtocol(ctx sdk.Context, address string) bool {
	return m.isReservedProtocol[address]
}

type mockCollectionService struct {
	collections map[uint64]*types.TokenCollection
	balances    map[string]*types.UserBalanceStore
	balanceErr  error
}

func newMockCollectionService() *mockCollectionService {
	return &mockCollectionService{
		collections: make(map[uint64]*types.TokenCollection),
		balances:    make(map[string]*types.UserBalanceStore),
	}
}

func (m *mockCollectionService) GetCollection(ctx sdk.Context, collectionId sdkmath.Uint) (*types.TokenCollection, bool) {
	c, ok := m.collections[collectionId.Uint64()]
	return c, ok
}

func (m *mockCollectionService) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool, error) {
	if m.balanceErr != nil {
		return nil, false, m.balanceErr
	}
	key := fmt.Sprintf("%d-%s", collection.CollectionId.Uint64(), userAddress)
	b, ok := m.balances[key]
	if !ok {
		return &types.UserBalanceStore{Balances: []*types.Balance{}}, false, nil
	}
	return b, true, nil
}

type mockDynamicStoreService struct {
	stores map[uint64]*types.DynamicStore
	values map[string]*types.DynamicStoreValue // key: "storeId-address"
}

func newMockDynamicStoreService() *mockDynamicStoreService {
	return &mockDynamicStoreService{
		stores: make(map[uint64]*types.DynamicStore),
		values: make(map[string]*types.DynamicStoreValue),
	}
}

func (m *mockDynamicStoreService) GetDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) (*types.DynamicStore, bool) {
	s, ok := m.stores[storeId.Uint64()]
	return s, ok
}

func (m *mockDynamicStoreService) GetDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) (*types.DynamicStoreValue, bool) {
	key := fmt.Sprintf("%d-%s", storeId.Uint64(), address)
	v, ok := m.values[key]
	return v, ok
}

type mockVotingServiceFull struct {
	votes    map[string]*types.VoteProof
	trackers map[string]*types.VotingChallengeTracker
}

func newMockVotingServiceFull() *mockVotingServiceFull {
	return &mockVotingServiceFull{
		votes:    make(map[string]*types.VoteProof),
		trackers: make(map[string]*types.VotingChallengeTracker),
	}
}

func (m *mockVotingServiceFull) GetVoteFromStore(ctx sdk.Context, key string) (*types.VoteProof, bool) {
	v, ok := m.votes[key]
	return v, ok
}

func (m *mockVotingServiceFull) GetVotingChallengeTrackerFromStore(ctx sdk.Context, key string) (*types.VotingChallengeTracker, bool) {
	t, ok := m.trackers[key]
	return t, ok
}

func (m *mockVotingServiceFull) ConstructVoteKey(collectionId sdkmath.Uint, approverAddress, approvalLevel, approvalId, proposalId, voterAddress string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s", collectionId.String(), approverAddress, approvalLevel, approvalId, proposalId, voterAddress)
}

func (m *mockVotingServiceFull) ConstructChallengeTrackerKey(collectionId sdkmath.Uint, approverAddress, approvalLevel, approvalId, proposalId string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", collectionId.String(), approverAddress, approvalLevel, approvalId, proposalId)
}

// ============================================================
// Helpers
// ============================================================

func ctxWithTime(t time.Time) sdk.Context {
	return sdk.Context{}.WithBlockTime(t)
}

func baseApproval() *types.CollectionApproval {
	return &types.CollectionApproval{
		ApprovalId: "test",
		ApprovalCriteria: &types.ApprovalCriteria{},
	}
}

func baseCollection() *types.TokenCollection {
	return &types.TokenCollection{
		CollectionId:      sdkmath.NewUint(1),
		MintEscrowAddress: "bb1mintescrow",
	}
}

// ============================================================
// AddressChecksChecker tests
// ============================================================

func TestAddressChecksChecker_NilChecks(t *testing.T) {
	svc := newMockAddressCheckService()
	checker := NewAddressChecksChecker(svc, nil, "sender")
	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestAddressChecksChecker_Name(t *testing.T) {
	svc := newMockAddressCheckService()
	require.Equal(t, "AddressChecks-sender", NewAddressChecksChecker(svc, &types.AddressChecks{}, "sender").Name())
	require.Equal(t, "AddressChecks-recipient", NewAddressChecksChecker(svc, &types.AddressChecks{}, "recipient").Name())
	require.Equal(t, "AddressChecks-initiator", NewAddressChecksChecker(svc, &types.AddressChecks{}, "initiator").Name())
}

func TestAddressChecksChecker_CheckTypeRouting(t *testing.T) {
	tests := []struct {
		checkType            string
		to, from, init       string
		expectedAddr         string
	}{
		{"sender", "toAddr", "fromAddr", "initAddr", "fromAddr"},
		{"recipient", "toAddr", "fromAddr", "initAddr", "toAddr"},
		{"initiator", "toAddr", "fromAddr", "initAddr", "initAddr"},
	}

	for _, tc := range tests {
		t.Run(tc.checkType, func(t *testing.T) {
			svc := newMockAddressCheckService()
			svc.isEvmContract[tc.expectedAddr] = true
			checks := &types.AddressChecks{MustBeEvmContract: true}
			checker := NewAddressChecksChecker(svc, checks, tc.checkType)
			msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), tc.to, tc.from, tc.init, "collection", "approver", nil, nil, "", false)
			require.NoError(t, err)
			require.Empty(t, msg)
		})
	}
}

func TestAddressChecksChecker_InvalidCheckType(t *testing.T) {
	svc := newMockAddressCheckService()
	checks := &types.AddressChecks{MustBeEvmContract: true}
	checker := NewAddressChecksChecker(svc, checks, "invalid")
	_, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
}

func TestAddressChecksChecker_MustBeEvmContract(t *testing.T) {
	svc := newMockAddressCheckService()
	checks := &types.AddressChecks{MustBeEvmContract: true}
	checker := NewAddressChecksChecker(svc, checks, "sender")

	// Fail: not a contract
	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "must be an EVM contract")

	// Pass: is a contract
	svc.isEvmContract["from"] = true
	msg, err = checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestAddressChecksChecker_MustNotBeEvmContract(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.isEvmContract["from"] = true
	checks := &types.AddressChecks{MustNotBeEvmContract: true}
	checker := NewAddressChecksChecker(svc, checks, "sender")

	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "must not be an EVM contract")
}

func TestAddressChecksChecker_MustBeLiquidityPool(t *testing.T) {
	svc := newMockAddressCheckService()
	checks := &types.AddressChecks{MustBeLiquidityPool: true}
	checker := NewAddressChecksChecker(svc, checks, "sender")

	// Fail
	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "must be a liquidity pool")

	// Pass
	svc.isLiquidityPool["from"] = true
	msg, err = checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestAddressChecksChecker_MustNotBeLiquidityPool(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.isLiquidityPool["from"] = true
	checks := &types.AddressChecks{MustNotBeLiquidityPool: true}
	checker := NewAddressChecksChecker(svc, checks, "sender")

	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "must not be a liquidity pool")
}

func TestAddressChecksChecker_ServiceError(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.evmContractErr = fmt.Errorf("rpc failure")
	checks := &types.AddressChecks{MustBeEvmContract: true}
	checker := NewAddressChecksChecker(svc, checks, "sender")

	_, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
}

// ============================================================
// Address equality checkers
// ============================================================

func TestRequireFromDoesNotEqualInitiatedBy(t *testing.T) {
	checker := NewRequireFromDoesNotEqualInitiatedByChecker()
	require.Equal(t, "RequireFromDoesNotEqualInitiatedBy", checker.Name())

	// Nil criteria -> pass
	msg, err := checker.Check(mockContext(), &types.CollectionApproval{}, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)

	approval := baseApproval()
	approval.ApprovalCriteria.RequireFromDoesNotEqualInitiatedBy = true

	// Different addresses -> pass
	msg, err = checker.Check(mockContext(), approval, baseCollection(), "to", "alice", "bob", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)

	// Same addresses -> fail
	msg, err = checker.Check(mockContext(), approval, baseCollection(), "to", "alice", "alice", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "from address equals initiated by")
}

func TestRequireFromEqualsInitiatedBy(t *testing.T) {
	checker := NewRequireFromEqualsInitiatedByChecker()

	approval := baseApproval()
	approval.ApprovalCriteria.RequireFromEqualsInitiatedBy = true

	// Same -> pass
	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "alice", "alice", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)

	// Different -> fail
	msg, err = checker.Check(mockContext(), approval, baseCollection(), "to", "alice", "bob", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "from address does not equal initiated by")
}

func TestRequireToDoesNotEqualInitiatedBy(t *testing.T) {
	checker := NewRequireToDoesNotEqualInitiatedByChecker()

	approval := baseApproval()
	approval.ApprovalCriteria.RequireToDoesNotEqualInitiatedBy = true

	// Different -> pass
	msg, err := checker.Check(mockContext(), approval, baseCollection(), "alice", "from", "bob", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)

	// Same -> fail
	msg, err = checker.Check(mockContext(), approval, baseCollection(), "alice", "from", "alice", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "to address equals initiated by")
}

func TestRequireToEqualsInitiatedBy(t *testing.T) {
	checker := NewRequireToEqualsInitiatedByChecker()

	approval := baseApproval()
	approval.ApprovalCriteria.RequireToEqualsInitiatedBy = true

	// Same -> pass
	msg, err := checker.Check(mockContext(), approval, baseCollection(), "alice", "from", "alice", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)

	// Different -> fail
	msg, err = checker.Check(mockContext(), approval, baseCollection(), "alice", "from", "bob", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "to address does not equal initiated by")
}

// ============================================================
// NoForcefulPostMintTransfersChecker tests
// ============================================================

func TestNoForcefulPostMintTransfers_NilCollection(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	require.Equal(t, "NoForcefulPostMintTransfers", checker.Name())

	msg, err := checker.Check(mockContext(), baseApproval(), nil, "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestNoForcefulPostMintTransfers_InvariantDisabled(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	col := baseCollection()
	col.Invariants = &types.CollectionInvariants{NoForcefulPostMintTransfers: false}

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	msg, err := checker.Check(mockContext(), approval, col, "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestNoForcefulPostMintTransfers_FromMint_Allowed(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	col := baseCollection()
	col.Invariants = &types.CollectionInvariants{NoForcefulPostMintTransfers: true}

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// From Mint -> always allowed even with override
	msg, err := checker.Check(mockContext(), approval, col, "to", "Mint", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestNoForcefulPostMintTransfers_OverrideOutgoing_Blocked(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	col := baseCollection()
	col.Invariants = &types.CollectionInvariants{NoForcefulPostMintTransfers: true}

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	msg, err := checker.Check(mockContext(), approval, col, "to", "alice", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "bypass user outgoing approvals")
}

func TestNoForcefulPostMintTransfers_OverrideIncoming_Blocked(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	col := baseCollection()
	col.Invariants = &types.CollectionInvariants{NoForcefulPostMintTransfers: true}

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesToIncomingApprovals = true

	msg, err := checker.Check(mockContext(), approval, col, "to", "alice", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "bypass user incoming approvals")
}

func TestNoForcefulPostMintTransfers_NoOverrides_Allowed(t *testing.T) {
	checker := NewNoForcefulPostMintTransfersChecker()
	col := baseCollection()
	col.Invariants = &types.CollectionInvariants{NoForcefulPostMintTransfers: true}

	approval := baseApproval()
	// No overrides set
	msg, err := checker.Check(mockContext(), approval, col, "to", "alice", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

// ============================================================
// ReservedProtocolAddressChecker tests
// ============================================================

func TestReservedProtocolAddress_NilCriteria(t *testing.T) {
	svc := newMockAddressCheckService()
	checker := NewReservedProtocolAddressChecker(svc)
	require.Equal(t, "ReservedProtocolAddress", checker.Name())

	approval := &types.CollectionApproval{ApprovalId: "test"}
	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestReservedProtocolAddress_NoOverride(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.isReservedProtocol["pool1"] = true
	checker := NewReservedProtocolAddressChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = false

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "pool1", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestReservedProtocolAddress_FromEqualsInitiator_Bypass(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.isReservedProtocol["pool1"] = true
	checker := NewReservedProtocolAddressChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// from == initiator -> bypass
	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "pool1", "pool1", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestReservedProtocolAddress_Blocked(t *testing.T) {
	svc := newMockAddressCheckService()
	svc.isReservedProtocol["pool1"] = true
	checker := NewReservedProtocolAddressChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "pool1", "attacker", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "reserved protocol addresses")
}

func TestReservedProtocolAddress_NotReserved_Allowed(t *testing.T) {
	svc := newMockAddressCheckService()
	checker := NewReservedProtocolAddressChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "regularAddr", "attacker", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

// ============================================================
// TransferTimesChecker tests
// ============================================================

func TestTransferTimes_NilApproval(t *testing.T) {
	checker := NewTransferTimesChecker()
	require.Equal(t, "TransferTimes", checker.Name())

	_, err := checker.Check(mockContext(), nil, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
}

func TestTransferTimes_InRange(t *testing.T) {
	checker := NewTransferTimesChecker()
	blockTime := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	ctx := ctxWithTime(blockTime)
	currMs := sdkmath.NewUint(uint64(blockTime.UnixMilli()))

	approval := baseApproval()
	approval.TransferTimes = []*types.UintRange{
		{Start: currMs.Sub(sdkmath.NewUint(1000)), End: currMs.Add(sdkmath.NewUint(1000))},
	}

	msg, err := checker.Check(ctx, approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestTransferTimes_OutOfRange(t *testing.T) {
	checker := NewTransferTimesChecker()
	blockTime := time.Date(2026, 4, 7, 12, 0, 0, 0, time.UTC)
	ctx := ctxWithTime(blockTime)

	approval := baseApproval()
	// Range far in the past
	approval.TransferTimes = []*types.UintRange{
		{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)},
	}

	msg, err := checker.Check(ctx, approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "transfer time not in range")
}

// ============================================================
// AltTimeChecksChecker tests
// ============================================================

func TestAltTimeChecks_Nil(t *testing.T) {
	checker := NewAltTimeChecksChecker(nil)
	require.Equal(t, "AltTimeChecks", checker.Name())

	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestAltTimeChecks_OfflineHour_Blocked(t *testing.T) {
	// Block hour 14
	altChecks := &types.AltTimeChecks{
		OfflineHours:          []*types.UintRange{{Start: sdkmath.NewUint(14), End: sdkmath.NewUint(14)}},
		TimezoneOffsetMinutes: sdkmath.NewUint(0),
	}
	checker := NewAltTimeChecksChecker(altChecks)

	blockTime := time.Date(2026, 4, 7, 14, 30, 0, 0, time.UTC) // Hour 14
	ctx := ctxWithTime(blockTime)

	msg, err := checker.Check(ctx, baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "offline hours")
}

func TestAltTimeChecks_OfflineHour_Allowed(t *testing.T) {
	altChecks := &types.AltTimeChecks{
		OfflineHours:          []*types.UintRange{{Start: sdkmath.NewUint(14), End: sdkmath.NewUint(14)}},
		TimezoneOffsetMinutes: sdkmath.NewUint(0),
	}
	checker := NewAltTimeChecksChecker(altChecks)

	blockTime := time.Date(2026, 4, 7, 10, 30, 0, 0, time.UTC) // Hour 10, not blocked
	ctx := ctxWithTime(blockTime)

	msg, err := checker.Check(ctx, baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestAltTimeChecks_OfflineDay_Blocked(t *testing.T) {
	// Block Sunday (0)
	altChecks := &types.AltTimeChecks{
		OfflineDays:           []*types.UintRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(0)}},
		TimezoneOffsetMinutes: sdkmath.NewUint(0),
	}
	checker := NewAltTimeChecksChecker(altChecks)

	blockTime := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC) // April 5, 2026 is a Sunday
	ctx := ctxWithTime(blockTime)

	msg, err := checker.Check(ctx, baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "offline days")
}

func TestAltTimeChecks_OfflineDay_Allowed(t *testing.T) {
	altChecks := &types.AltTimeChecks{
		OfflineDays:           []*types.UintRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(0)}},
		TimezoneOffsetMinutes: sdkmath.NewUint(0),
	}
	checker := NewAltTimeChecksChecker(altChecks)

	blockTime := time.Date(2026, 4, 7, 10, 0, 0, 0, time.UTC) // April 7, 2026 is Tuesday (2)
	ctx := ctxWithTime(blockTime)

	msg, err := checker.Check(ctx, baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

// ============================================================
// DynamicStoreChallengesChecker tests
// ============================================================

func TestDynamicStoreChallenges_NoChallenges(t *testing.T) {
	svc := newMockDynamicStoreService()
	checker := NewDynamicStoreChallengesChecker(svc)
	require.Equal(t, "DynamicStoreChallenges", checker.Name())

	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestDynamicStoreChallenges_StoreNotFound(t *testing.T) {
	svc := newMockDynamicStoreService()
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
		{StoreId: sdkmath.NewUint(999), OwnershipCheckParty: "initiator"},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "dynamic store not found")
}

func TestDynamicStoreChallenges_GlobalDisabled(t *testing.T) {
	svc := newMockDynamicStoreService()
	svc.stores[1] = &types.DynamicStore{StoreId: sdkmath.NewUint(1), GlobalEnabled: false, DefaultValue: true}
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
		{StoreId: sdkmath.NewUint(1), OwnershipCheckParty: "initiator"},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "globally disabled")
}

func TestDynamicStoreChallenges_ValueTrue_Pass(t *testing.T) {
	svc := newMockDynamicStoreService()
	svc.stores[1] = &types.DynamicStore{StoreId: sdkmath.NewUint(1), GlobalEnabled: true, DefaultValue: false}
	svc.values["1-init"] = &types.DynamicStoreValue{StoreId: sdkmath.NewUint(1), Address: "init", Value: true}
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
		{StoreId: sdkmath.NewUint(1), OwnershipCheckParty: "initiator"},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestDynamicStoreChallenges_ValueFalse_Fail(t *testing.T) {
	svc := newMockDynamicStoreService()
	svc.stores[1] = &types.DynamicStore{StoreId: sdkmath.NewUint(1), GlobalEnabled: true, DefaultValue: false}
	svc.values["1-init"] = &types.DynamicStoreValue{StoreId: sdkmath.NewUint(1), Address: "init", Value: false}
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
		{StoreId: sdkmath.NewUint(1), OwnershipCheckParty: "initiator"},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "does not have permission")
}

func TestDynamicStoreChallenges_DefaultValue_Used(t *testing.T) {
	svc := newMockDynamicStoreService()
	svc.stores[1] = &types.DynamicStore{StoreId: sdkmath.NewUint(1), GlobalEnabled: true, DefaultValue: true}
	// No per-address value set -- should use default (true)
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
		{StoreId: sdkmath.NewUint(1), OwnershipCheckParty: "sender"},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestDynamicStoreChallenges_NilChallenge(t *testing.T) {
	svc := newMockDynamicStoreService()
	checker := NewDynamicStoreChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{nil}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "challenge is nil")
}

func TestDynamicStoreChallenges_PartyRouting(t *testing.T) {
	tests := []struct {
		party    string
		expected string
	}{
		{"initiator", "init"},
		{"sender", "from"},
		{"recipient", "to"},
		{"Mint", "bb1mintescrow"},
		{"", "init"}, // defaults to initiator
	}

	for _, tc := range tests {
		t.Run(tc.party, func(t *testing.T) {
			svc := newMockDynamicStoreService()
			svc.stores[1] = &types.DynamicStore{StoreId: sdkmath.NewUint(1), GlobalEnabled: true, DefaultValue: false}
			svc.values[fmt.Sprintf("1-%s", tc.expected)] = &types.DynamicStoreValue{Value: true}
			checker := NewDynamicStoreChallengesChecker(svc)

			approval := baseApproval()
			approval.ApprovalCriteria.DynamicStoreChallenges = []*types.DynamicStoreChallenge{
				{StoreId: sdkmath.NewUint(1), OwnershipCheckParty: tc.party},
			}

			msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
			require.NoError(t, err, "party=%s expected=%s", tc.party, tc.expected)
			require.Empty(t, msg)
		})
	}
}

// ============================================================
// VotingChallengesChecker tests
// ============================================================

func TestVotingChallenges_NoChallenges(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)
	require.Equal(t, "VotingChallenges", checker.Name())

	msg, err := checker.Check(mockContext(), baseApproval(), baseCollection(), "to", "from", "init", "", "", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestVotingChallenges_QuorumMet(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(50),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters: []*types.Voter{
				{Address: "voter1", Weight: sdkmath.NewUint(60)},
				{Address: "voter2", Weight: sdkmath.NewUint(40)},
			},
		},
	}

	col := baseCollection()

	// voter1 votes 100% yes (weight 60), voter2 not voting
	// totalYes = 60*100/100 = 60, totalPossible = 100, pct = 60%
	voteKey := svc.ConstructVoteKey(col.CollectionId, "approver", "collection", "test", "prop1", "voter1")
	svc.votes[voteKey] = &types.VoteProof{ProposalId: "prop1", Voter: "voter1", YesWeight: sdkmath.NewUint(100)}

	msg, err := checker.Check(mockContext(), approval, col, "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, msg)
}

func TestVotingChallenges_QuorumNotMet(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(80),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters: []*types.Voter{
				{Address: "voter1", Weight: sdkmath.NewUint(50)},
				{Address: "voter2", Weight: sdkmath.NewUint(50)},
			},
		},
	}

	col := baseCollection()
	// Only voter1 votes yes at 50% -> (50*50/100)=25, pct=25%, need 80%
	voteKey := svc.ConstructVoteKey(col.CollectionId, "approver", "collection", "test", "prop1", "voter1")
	svc.votes[voteKey] = &types.VoteProof{ProposalId: "prop1", Voter: "voter1", YesWeight: sdkmath.NewUint(50)}

	msg, err := checker.Check(mockContext(), approval, col, "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "threshold not met")
}

func TestVotingChallenges_NilVoter(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(50),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters:           []*types.Voter{nil},
		},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "voter is nil")
}

func TestVotingChallenges_EmptyVoterAddress(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(50),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters:           []*types.Voter{{Address: "", Weight: sdkmath.NewUint(10)}},
		},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "voter address is empty")
}

func TestVotingChallenges_ZeroWeightVoter(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(50),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters:           []*types.Voter{{Address: "voter1", Weight: sdkmath.NewUint(0)}},
		},
	}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "zero weight")
}

func TestVotingChallenges_NilChallenge(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{nil}

	msg, err := checker.Check(mockContext(), approval, baseCollection(), "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "voting challenge is nil")
}

func TestVotingChallenges_YesWeightExceeds100(t *testing.T) {
	svc := newMockVotingServiceFull()
	checker := NewVotingChallengesChecker(svc)

	approval := baseApproval()
	approval.ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{
		{
			ProposalId:       "prop1",
			QuorumThreshold:  sdkmath.NewUint(50),
			DelayAfterQuorum: sdkmath.NewUint(0),
			Voters:           []*types.Voter{{Address: "voter1", Weight: sdkmath.NewUint(100)}},
		},
	}

	col := baseCollection()
	voteKey := svc.ConstructVoteKey(col.CollectionId, "approver", "collection", "test", "prop1", "voter1")
	svc.votes[voteKey] = &types.VoteProof{ProposalId: "prop1", Voter: "voter1", YesWeight: sdkmath.NewUint(101)}

	msg, err := checker.Check(mockContext(), approval, col, "to", "from", "init", "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, msg, "yesWeight exceeds 100")
}
