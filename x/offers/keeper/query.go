package keeper

import (
	"bitbadgeschain/x/offers/types"
)

var _ types.QueryServer = Keeper{}
