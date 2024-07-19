package keeper

import (
	"bitbadgeschain/x/anchor/types"
)

var _ types.QueryServer = Keeper{}
