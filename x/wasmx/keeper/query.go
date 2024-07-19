package keeper

import (
	"bitbadgeschain/x/wasmx/types"
)

var _ types.QueryServer = Keeper{}
