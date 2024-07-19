package keeper

import (
	"bitbadgeschain/x/badges/types"
)

func IsStandardBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Standard"
}

func IsOffChainBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Off-Chain - Indexed"
}

func IsNonPublicBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Non-Public"
}

func IsNonIndexedBalances(collection *types.BadgeCollection) bool {
	return collection.BalancesType == "Off-Chain - Non-Indexed"
}
