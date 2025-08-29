package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestCollectionInvariants_MaxSupplyPerId(t *testing.T) {
	// Test that we can create CollectionInvariants with maxSupplyPerId
	invariants := &types.CollectionInvariants{
		NoCustomOwnershipTimes: true,
		MaxSupplyPerId:         sdkmath.NewUint(1000),
	}

	require.True(t, invariants.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(1000), invariants.MaxSupplyPerId)

	// Test that we can create CollectionInvariants without maxSupplyPerId
	invariants2 := &types.CollectionInvariants{
		NoCustomOwnershipTimes: false,
		MaxSupplyPerId:         sdkmath.NewUint(0),
	}

	require.False(t, invariants2.NoCustomOwnershipTimes)
	require.True(t, invariants2.MaxSupplyPerId.IsZero())
}

func TestCollectionInvariants_MaxSupplyPerId_GreaterThanZero(t *testing.T) {
	// Test that we can create CollectionInvariants with maxSupplyPerId
	invariants := &types.CollectionInvariants{
		NoCustomOwnershipTimes: true,
		MaxSupplyPerId:         sdkmath.NewUint(1000),
	}

	require.True(t, invariants.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(1000), invariants.MaxSupplyPerId)

	// Test that we can create CollectionInvariants without maxSupplyPerId
	invariants2 := &types.CollectionInvariants{
		NoCustomOwnershipTimes: false,
		MaxSupplyPerId:         sdkmath.NewUint(2),
	}

	require.False(t, invariants2.NoCustomOwnershipTimes)
	require.Equal(t, sdkmath.NewUint(2), invariants2.MaxSupplyPerId)
}

