package keeper

import (
	"bitbadgeschain/x/badges/types"
)

var _ types.QueryServer = Keeper{}
