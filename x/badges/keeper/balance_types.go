package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func IsStandardBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Standard"
}

func IsOffChainBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Off-Chain"
}

func IsInheritedBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Inherited"
}
