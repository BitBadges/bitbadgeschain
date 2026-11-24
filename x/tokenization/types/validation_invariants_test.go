package types_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestApprovalCriteriaCastPreservesAllSharedFields verifies that casting between
// CollectionApproval, UserOutgoingApproval, and UserIncomingApproval preserves
// every shared ApprovalCriteria field, including EvmQueryChallenges.
func TestApprovalCriteriaCastPreservesAllSharedFields(t *testing.T) {
	evmChallenge := &types.EVMQueryChallenge{
		ContractAddress:    "0x1111111111111111111111111111111111111111",
		Calldata:           "0xdeadbeef",
		ExpectedResult:     "0x0000000000000000000000000000000000000000000000000000000000000001",
		ComparisonOperator: "eq",
		GasLimit:           sdkmath.NewUint(100000),
		Uri:                "https://example.com/q",
		CustomData:         "{}",
	}

	outApproval := &types.UserOutgoingApproval{
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "out-1",
		Version:           sdkmath.ZeroUint(),
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			EvmQueryChallenges: []*types.EVMQueryChallenge{evmChallenge},
		},
	}

	colFromOut := types.CastOutgoingTransferToCollectionTransfer(outApproval, "bb1from")
	require.Equal(t, []*types.EVMQueryChallenge{evmChallenge}, colFromOut.ApprovalCriteria.EvmQueryChallenges)

	roundTripOut := types.CastFromCollectionTransferToOutgoingTransfer(colFromOut)
	require.Equal(t, []*types.EVMQueryChallenge{evmChallenge}, roundTripOut.ApprovalCriteria.EvmQueryChallenges)

	inApproval := &types.UserIncomingApproval{
		FromListId:        "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "in-1",
		Version:           sdkmath.ZeroUint(),
		ApprovalCriteria: &types.IncomingApprovalCriteria{
			EvmQueryChallenges: []*types.EVMQueryChallenge{evmChallenge},
		},
	}

	colFromIn := types.CastIncomingTransferToCollectionTransfer(inApproval, "bb1to")
	require.Equal(t, []*types.EVMQueryChallenge{evmChallenge}, colFromIn.ApprovalCriteria.EvmQueryChallenges)

	roundTripIn := types.CastFromCollectionTransferToIncomingTransfer(colFromIn)
	require.Equal(t, []*types.EVMQueryChallenge{evmChallenge}, roundTripIn.ApprovalCriteria.EvmQueryChallenges)

	// Reflect over shared fields between OutgoingApprovalCriteria / IncomingApprovalCriteria and ApprovalCriteria
	// so any newly added shared field is covered automatically.
	fullType := reflect.TypeOf(types.ApprovalCriteria{})
	for _, subType := range []reflect.Type{
		reflect.TypeOf(types.OutgoingApprovalCriteria{}),
		reflect.TypeOf(types.IncomingApprovalCriteria{}),
	} {
		for i := 0; i < subType.NumField(); i++ {
			f := subType.Field(i)
			if !f.IsExported() {
				continue
			}
			fullField, ok := fullType.FieldByName(f.Name)
			require.True(t, ok, "field %s on %s must exist on ApprovalCriteria", f.Name, subType.Name())
			require.Equal(t, f.Type, fullField.Type, "field %s type mismatch", f.Name)
		}
	}
}

// TestIncrementBalancesRecurringPeriodContinuesAllEntries verifies that when the
// first recurring charge period matches a balance entry, subsequent entries in
// the slice are also updated, and arithmetic overflow returns an error.
func TestIncrementBalancesRecurringPeriodContinuesAllEntries(t *testing.T) {
	// Block time 950ms sits in first charge period [900, 1000) for startTime=1000, chargePeriodLength=100.
	ctx := sdk.Context{}.WithBlockTime(time.UnixMilli(950))
	inc := &types.IncrementedBalances{
		StartBalances: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(10),
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
			},
			{
				Amount:         sdkmath.NewUint(20),
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
			},
		},
		IncrementTokenIdsBy:       sdkmath.NewUint(5),
		IncrementOwnershipTimesBy: sdkmath.ZeroUint(),
		DurationFromTimestamp:     sdkmath.ZeroUint(),
		RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
			StartTime:          sdkmath.NewUint(1000),
			IntervalLength:     sdkmath.NewUint(500),
			ChargePeriodLength: sdkmath.NewUint(100),
		},
	}

	got, err := types.IncrementBalances(ctx, inc, sdkmath.NewUint(1), nil, nil)
	require.NoError(t, err)
	require.Len(t, got, 2)
	// Both balance entries must have recurring ownership times applied AND token IDs incremented by 5.
	require.Equal(t, sdkmath.NewUint(1000), got[0].OwnershipTimes[0].Start)
	require.Equal(t, sdkmath.NewUint(1499), got[0].OwnershipTimes[0].End)
	require.Equal(t, sdkmath.NewUint(6), got[0].TokenIds[0].Start)
	require.Equal(t, sdkmath.NewUint(1000), got[1].OwnershipTimes[0].Start)
	require.Equal(t, sdkmath.NewUint(1499), got[1].OwnershipTimes[0].End)
	require.Equal(t, sdkmath.NewUint(7), got[1].TokenIds[0].Start, "second balance entry must also be processed on first charge period")

	// Overflow check: uint256 max + intervalLength returns error instead of panicking.
	maxUint256 := sdkmath.NewUintFromBigInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)))
	overflowInc := &types.IncrementedBalances{
		StartBalances: []*types.Balance{{
			Amount:         sdkmath.OneUint(),
			TokenIds:       []*types.UintRange{{Start: maxUint256, End: maxUint256}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.NewUint(10)}},
		}},
		IncrementTokenIdsBy:       sdkmath.OneUint(),
		IncrementOwnershipTimesBy: sdkmath.ZeroUint(),
		DurationFromTimestamp:     sdkmath.ZeroUint(),
	}
	require.NotPanics(t, func() {
		_, err = types.IncrementBalances(ctx, overflowInc, sdkmath.OneUint(), nil, nil)
		require.Error(t, err)
	})
}

// TestValidateUserPermissionsChecksActionPermissionsAndIdentifiers verifies
// UserPermissions validates all ActionPermission fields and non-empty identifiers.
func TestValidateUserPermissionsChecksActionPermissionsAndIdentifiers(t *testing.T) {
	overlap := []*types.ActionPermission{{
		PermanentlyPermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		PermanentlyForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(5), End: sdkmath.NewUint(15)}},
	}}

	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateAutoApproveSelfInitiatedIncomingTransfers: overlap,
	}, false))
	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: overlap,
	}, false))
	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateAutoApproveAllIncomingTransfers: overlap,
	}, false))

	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateAutoApproveAllIncomingTransfers: []*types.ActionPermission{{}, {}},
	}, false))

	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{{
			ApprovalId: "", ToListId: "All", InitiatedByListId: "All",
		}},
	}, false))
	require.Error(t, types.ValidateUserPermissions(&types.UserPermissions{
		CanUpdateIncomingApprovals: []*types.UserIncomingApprovalPermission{{
			ApprovalId: "id", FromListId: "", InitiatedByListId: "All",
		}},
	}, false))
}

// TestApprovalValidationChecksWithoutMutation verifies that MustOwnTokens,
// DynamicStoreChallenges, VotingChallenges, reserved approval IDs, auto-scan
// eligibility, nil EVMQueryChallenge.GasLimit, and MerkleChallenges are all
// checked even when canChangeValues is false.
func TestApprovalValidationChecksWithoutMutation(t *testing.T) {
	overMaxUint64 := sdkmath.NewUintFromBigInt(new(big.Int).Add(new(big.Int).SetUint64(math.MaxUint64), big.NewInt(1)))
	wrapApproval := func(criteria *types.ApprovalCriteria) []*types.CollectionApproval {
		return []*types.CollectionApproval{{
			ApprovalId:        "test-approval",
			FromListId:        "All",
			ToListId:          "All",
			InitiatedByListId: "All",
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			Version:           sdkmath.ZeroUint(),
			ApprovalCriteria:  criteria,
		}}
	}

	// MustOwnTokens CollectionId > MaxUint64 rejected when canChangeValues == false
	err := types.ValidateCollectionApprovals(sdk.Context{}, wrapApproval(&types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{{
			CollectionId:   overMaxUint64,
			AmountRange:    &types.UintRange{Start: sdkmath.OneUint(), End: sdkmath.NewUint(10)},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		}},
	}), false)
	require.Error(t, err)

	// DynamicStoreChallenges StoreId > MaxUint64 rejected when canChangeValues == false
	err = types.ValidateCollectionApprovals(sdk.Context{}, wrapApproval(&types.ApprovalCriteria{
		DynamicStoreChallenges: []*types.DynamicStoreChallenge{{
			StoreId: overMaxUint64,
		}},
	}), false)
	require.Error(t, err)

	// VotingChallenges invalid quorum > 100 rejected when canChangeValues == false
	err = types.ValidateCollectionApprovals(sdk.Context{}, wrapApproval(&types.ApprovalCriteria{
		VotingChallenges: []*types.VotingChallenge{{
			ProposalId:      "p1",
			QuorumThreshold: sdkmath.NewUint(101),
			Voters:          []*types.Voter{{Address: sample.AccAddress(), Weight: sdkmath.OneUint()}},
		}},
	}), false)
	require.Error(t, err)

	// EVMQueryChallenge with zero/uninitialized GasLimit does not panic
	require.NotPanics(t, func() {
		err = types.ValidateEVMQueryChallenges([]*types.EVMQueryChallenge{{
			ContractAddress:    "0x1111111111111111111111111111111111111111",
			Calldata:           "0xdeadbeef",
			ExpectedResult:     "0x0000000000000000000000000000000000000000000000000000000000000001",
			ComparisonOperator: "eq",
		}})
		require.NoError(t, err)
	})

	// Reserved approval IDs reject all-incoming and one-time-outgoing
	for _, reserved := range []string{"all-incoming", "one-time-outgoing"} {
		approvals := wrapApproval(&types.ApprovalCriteria{})
		approvals[0].ApprovalId = reserved
		require.Error(t, types.ValidateCollectionApprovals(sdk.Context{}, approvals, false))
	}

	// Approvals with EvmQueryChallenges or VotingChallenges are not auto-scannable
	require.False(t, types.CollectionApprovalIsAutoScannable(&types.ApprovalCriteria{
		EvmQueryChallenges: []*types.EVMQueryChallenge{{ContractAddress: "0x1111111111111111111111111111111111111111"}},
	}))
	require.False(t, types.CollectionApprovalIsAutoScannable(&types.ApprovalCriteria{
		VotingChallenges: []*types.VotingChallenge{{ProposalId: "p1"}},
	}))

	// MerkleChallenges rejects nil element or empty root
	require.Error(t, types.ValidateMerkleChallenges([]*types.MerkleChallenge{nil}, false, ""))
	require.Error(t, types.ValidateMerkleChallenges([]*types.MerkleChallenge{{Root: ""}}, false, ""))
}

// TestDefaultBalancesAndGenesisValidation checks NoCustomOwnershipTimes on
// DefaultBalances, MsgCastVote nil YesWeight, and TotalAddress in genesis balances.
func TestDefaultBalancesAndGenesisValidation(t *testing.T) {
	// DefaultBalances with custom ownership times fails when NoCustomOwnershipTimes is set
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      sample.AccAddress(),
		CollectionId: sdkmath.ZeroUint(),
		Invariants:   &types.InvariantsAddObject{NoCustomOwnershipTimes: true},
		DefaultBalances: &types.UserBalanceStore{
			OutgoingApprovals: []*types.UserOutgoingApproval{{
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes:     GetFullUintRanges(),
				TokenIds:          GetFullUintRanges(),
				OwnershipTimes:    []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
				ApprovalId:        "out-1",
				Version:           sdkmath.ZeroUint(),
			}},
		},
	}
	require.Error(t, msg.ValidateBasic())

	// MsgCastVote with uninitialized YesWeight returns error without panicking
	castMsg := &types.MsgCastVote{
		Creator:       sample.AccAddress(),
		CollectionId:  sdkmath.OneUint(),
		ApprovalLevel: "collection",
		ApprovalId:    "a1",
		ProposalId:    "p1",
	}
	require.NotPanics(t, func() {
		require.Error(t, castMsg.ValidateBasic())
	})

	// GenesisState allows TotalAddress in balance store keys
	gs := types.DefaultGenesis()
	gs.Collections = []*types.TokenCollection{{CollectionId: sdkmath.OneUint()}}
	gs.NextCollectionId = sdkmath.NewUint(2)
	gs.BalanceStoreKeys = []string{"1-" + types.TotalAddress}
	gs.Balances = []*types.UserBalanceStore{{
		Balances: []*types.Balance{},
	}}
	require.NoError(t, gs.Validate())
}
