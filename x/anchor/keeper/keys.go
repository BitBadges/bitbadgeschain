package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var (
	NextLocationIdKey = []byte{0x0A}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName