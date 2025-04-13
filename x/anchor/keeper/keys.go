package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
)

var (
	NextLocationIdKey = []byte{0x0A}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

const StoreKey = types.ModuleName
