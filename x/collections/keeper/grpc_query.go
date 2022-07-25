package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/collections/types"
)

var _ types.QueryServer = Keeper{}
