package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"
)

var _ types.QueryServer = Keeper{}
