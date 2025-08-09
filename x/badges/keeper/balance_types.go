package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func IsStandardBalances(collection *types.TokenCollection) bool {
	return collection.BalancesType == "Standard"
}

func IsOffChainBalances(collection *types.TokenCollection) bool {
	return collection.BalancesType == "Off-Chain - Indexed"
}

func IsNonPublicBalances(collection *types.TokenCollection) bool {
	return collection.BalancesType == "Non-Public"
}

func IsNonIndexedBalances(collection *types.TokenCollection) bool {
	return collection.BalancesType == "Off-Chain - Non-Indexed"
}
