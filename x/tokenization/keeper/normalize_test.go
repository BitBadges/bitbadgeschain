package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
)

// Ensures NormalizeNilPointers is a no-op on nil / non-pointer / non-struct
// inputs (must not panic).
func TestNormalize_DegenerateInputs(t *testing.T) {
	require.NotPanics(t, func() { NormalizeNilPointers(nil) })

	var nilCollection *types.TokenCollection
	require.NotPanics(t, func() { NormalizeNilPointers(nilCollection) })

	notAPtr := types.TokenCollection{}
	require.NotPanics(t, func() { NormalizeNilPointers(notAPtr) })

	prim := 42
	require.NotPanics(t, func() { NormalizeNilPointers(&prim) })
}

// Top-level singular pointer-to-struct field: nil → fresh empty struct.
func TestNormalize_FillsNilSingularPointer(t *testing.T) {
	c := &types.TokenCollection{}
	require.Nil(t, c.CollectionPermissions)
	require.Nil(t, c.DefaultBalances)
	require.Nil(t, c.CollectionMetadata)
	require.Nil(t, c.Invariants)

	NormalizeNilPointers(c)

	require.NotNil(t, c.CollectionPermissions)
	require.NotNil(t, c.DefaultBalances)
	require.NotNil(t, c.CollectionMetadata)
	require.NotNil(t, c.Invariants)
}

// Nested fields: defaultBalances.userPermissions inside the collection
// should be filled after recursing through the parent pointer.
func TestNormalize_RecursesIntoFilledPointers(t *testing.T) {
	c := &types.TokenCollection{}

	NormalizeNilPointers(c)

	require.NotNil(t, c.DefaultBalances)
	require.NotNil(t, c.DefaultBalances.UserPermissions)
}

// CollectionApproval has ApprovalCriteria → multiple nested struct
// fields. Reaching the deepest one through the recursive walker is
// the load-bearing case for keeper code that does
// `approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId`.
func TestNormalize_DeepRecursionThroughApprovalCriteria(t *testing.T) {
	approval := &types.CollectionApproval{}

	NormalizeNilPointers(approval)

	require.NotNil(t, approval.ApprovalCriteria, "top-level criteria")
	require.NotNil(t, approval.ApprovalCriteria.PredeterminedBalances)
	require.NotNil(t, approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod)
	require.NotNil(t, approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances)
	require.NotNil(t, approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes)
	require.NotNil(t, approval.ApprovalCriteria.ApprovalAmounts)
	require.NotNil(t, approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals)
	require.NotNil(t, approval.ApprovalCriteria.MaxNumTransfers)
	require.NotNil(t, approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals)
	require.NotNil(t, approval.ApprovalCriteria.AutoDeletionOptions)
	require.NotNil(t, approval.ApprovalCriteria.SenderChecks)
	require.NotNil(t, approval.ApprovalCriteria.RecipientChecks)
	require.NotNil(t, approval.ApprovalCriteria.InitiatorChecks)
	require.NotNil(t, approval.ApprovalCriteria.AltTimeChecks)
	require.NotNil(t, approval.ApprovalCriteria.UserApprovalSettings)
	require.NotNil(t, approval.ApprovalCriteria.UserApprovalSettings.UserRoyalties)

	// The exact panic-site access from challenges.go:42 — must not panic
	// after normalize.
	require.NotPanics(t, func() {
		_ = approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId
		_ = approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex
	})
}

// Slices of pointers: each element should also be normalized.
// CollectionApprovals[].ApprovalCriteria etc.
func TestNormalize_RecursesIntoSliceElements(t *testing.T) {
	c := &types.TokenCollection{
		CollectionApprovals: []*types.CollectionApproval{
			{}, // empty approval — ApprovalCriteria is nil
			{}, // empty approval — ApprovalCriteria is nil
		},
	}

	NormalizeNilPointers(c)

	require.Len(t, c.CollectionApprovals, 2)
	for i, ap := range c.CollectionApprovals {
		require.NotNilf(t, ap.ApprovalCriteria, "approval %d: criteria", i)
		require.NotNilf(t, ap.ApprovalCriteria.PredeterminedBalances, "approval %d: predetermined", i)
		require.NotNilf(t, ap.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod, "approval %d: order calc", i)
	}
}

// Already-filled fields are left intact (we don't overwrite caller data).
func TestNormalize_PreservesExistingValues(t *testing.T) {
	c := &types.TokenCollection{
		Manager: "bb1existing",
		CollectionPermissions: &types.CollectionPermissions{
			CanArchiveCollection: []*types.ActionPermission{
				{},
			},
		},
		Invariants: &types.CollectionInvariants{
			DisablePoolCreation: true,
		},
	}

	NormalizeNilPointers(c)

	require.Equal(t, "bb1existing", c.Manager)
	require.NotNil(t, c.CollectionPermissions)
	require.Len(t, c.CollectionPermissions.CanArchiveCollection, 1)
	require.NotNil(t, c.Invariants)
	require.True(t, c.Invariants.DisablePoolCreation)
}

// Empty slice fields remain empty (we don't materialize bogus elements).
func TestNormalize_EmptySlicesStayEmpty(t *testing.T) {
	c := &types.TokenCollection{}
	NormalizeNilPointers(c)

	require.NotNil(t, c.CollectionPermissions)
	// CollectionPermissions has many `[]*TimedUpdatePermission` slice fields.
	// Walker should NOT add elements — slice is nil/empty, so we just leave it.
	require.Empty(t, c.CollectionPermissions.CanArchiveCollection)
	require.Empty(t, c.CollectionPermissions.CanDeleteCollection)
}

// Deeply nested AliasPath: TokenCollection → AliasPaths[] → Conversion → SideA, etc.
func TestNormalize_AliasPathChain(t *testing.T) {
	c := &types.TokenCollection{
		AliasPaths: []*types.AliasPath{
			{Denom: "ucredit"}, // Conversion is nil
		},
	}

	NormalizeNilPointers(c)

	require.Len(t, c.AliasPaths, 1)
	ap := c.AliasPaths[0]
	require.Equal(t, "ucredit", ap.Denom)
	require.NotNil(t, ap.Conversion)
	require.NotNil(t, ap.Conversion.SideA)
	require.NotNil(t, ap.Metadata)
}

// UserBalanceStore: defaultBalances analogue used for user storage.
func TestNormalize_UserBalanceStore(t *testing.T) {
	ub := &types.UserBalanceStore{}
	NormalizeNilPointers(ub)
	require.NotNil(t, ub.UserPermissions)
}

// `sdkmath.Uint{}` (the Go zero value) wraps a nil `*big.Int`; any math
// op panics on it. Proto fields tagged `customtype = "Uint", nullable
// = false` deserialize to the zero value when the wire bytes don't
// include the field — common after our SDK strips empty Uint strings.
// Normalize must replace with `NewUint(0)` so downstream math is safe.
func TestNormalize_RecoversNilUintInternals(t *testing.T) {
	// AltTimeChecks.TimezoneOffsetMinutes is exactly such a field.
	atc := &types.AltTimeChecks{}
	require.True(t, atc.TimezoneOffsetMinutes.IsNil(), "before normalize: nil internal big.Int")

	NormalizeNilPointers(atc)

	require.False(t, atc.TimezoneOffsetMinutes.IsNil(), "after normalize: must be initialized to 0")

	// Math operations must no longer panic.
	require.NotPanics(t, func() {
		_ = atc.TimezoneOffsetMinutes.GT(sdkmath.NewUint(0))
	})
}

// Same fix should apply through nested chains: Approval →
// ApprovalCriteria → ApprovalAmounts → resetTimeIntervals → uint
// fields, etc.
func TestNormalize_RecoversNilUintsThroughNesting(t *testing.T) {
	approval := &types.CollectionApproval{}
	NormalizeNilPointers(approval)

	require.NotPanics(t, func() {
		ac := approval.ApprovalCriteria
		require.NotNil(t, ac)
		require.NotNil(t, ac.AltTimeChecks)
		_ = ac.AltTimeChecks.TimezoneOffsetMinutes.GT(sdkmath.NewUint(0))

		require.NotNil(t, ac.ApprovalAmounts)
		require.NotNil(t, ac.ApprovalAmounts.ResetTimeIntervals)
		_ = ac.ApprovalAmounts.ResetTimeIntervals.StartTime.GT(sdkmath.NewUint(0))
		_ = ac.ApprovalAmounts.ResetTimeIntervals.IntervalLength.GT(sdkmath.NewUint(0))

		require.NotNil(t, ac.MaxNumTransfers)
		_ = ac.MaxNumTransfers.OverallMaxNumTransfers.GT(sdkmath.NewUint(0))
	})
}
