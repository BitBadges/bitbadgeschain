package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestCollectionInvariants_MaxSupplyPerId(t *testing.T) {
	// Test that we can create CollectionInvariants with maxSupplyPerId
	invariants := &types.CollectionInvariants{
		NoCustomOwnershipTimes: true,
		MaxSupplyPerId:         sdkmath.NewUint(1000),
	}

	require.True(t, invariants.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(1000), invariants.MaxSupplyPerId)

	// Test that we can create CollectionInvariants without maxSupplyPerId (defaults to nil = unlimited)
	invariants2 := &types.CollectionInvariants{
		NoCustomOwnershipTimes: false,
	}

	require.False(t, invariants2.NoCustomOwnershipTimes)
	// When not set, MaxSupplyPerId is nil (which means unlimited)
	require.True(t, invariants2.MaxSupplyPerId.IsNil())
}

func TestCollectionInvariants_MaxSupplyPerId_GreaterThanZero(t *testing.T) {
	// Test that we can create CollectionInvariants with maxSupplyPerId
	invariants := &types.CollectionInvariants{
		NoCustomOwnershipTimes: true,
		MaxSupplyPerId:         sdkmath.NewUint(1000),
	}

	require.True(t, invariants.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(1000), invariants.MaxSupplyPerId)

	// Test that we can create CollectionInvariants with maxSupplyPerId = 2
	invariants2 := &types.CollectionInvariants{
		NoCustomOwnershipTimes: false,
		MaxSupplyPerId:         sdkmath.NewUint(2),
	}

	require.False(t, invariants2.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(2), invariants2.MaxSupplyPerId)
}
