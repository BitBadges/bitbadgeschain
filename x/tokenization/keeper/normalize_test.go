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

// NormalizeNilPointers does NOT fill nil pointer-to-struct fields —
// that would break nil-vs-empty equality semantics. Confirm.
func TestNormalize_DoesNotFillNilPointers(t *testing.T) {
	c := &types.TokenCollection{}
	require.Nil(t, c.CollectionPermissions)
	require.Nil(t, c.DefaultBalances)

	NormalizeNilPointers(c)

	require.Nil(t, c.CollectionPermissions, "must NOT be filled — preserves 'unset' semantic")
	require.Nil(t, c.DefaultBalances)
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

// Walks into ALREADY-non-nil nested pointers (does not fill nils),
// so any uninitialized Uints inside user-set sub-messages get fixed.
func TestNormalize_RecoversNilUintsInNestedFilledPointers(t *testing.T) {
	approval := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			AltTimeChecks: &types.AltTimeChecks{},
		},
	}

	NormalizeNilPointers(approval)

	require.NotPanics(t, func() {
		_ = approval.ApprovalCriteria.AltTimeChecks.TimezoneOffsetMinutes.GT(sdkmath.NewUint(0))
	})
}
