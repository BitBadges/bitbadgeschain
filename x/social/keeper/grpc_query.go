package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/social/types"
)

var _ types.QueryServer = Keeper{}
