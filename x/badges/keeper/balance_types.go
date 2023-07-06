package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func IsStandardBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType.Equal(sdkmath.NewUint(0))
}

func IsOffChainBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType.Equal(sdkmath.NewUint(1))
}

func IsInheritedBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType.Equal(sdkmath.NewUint(2))
}
