package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ types.QueryServer = Keeper{}
