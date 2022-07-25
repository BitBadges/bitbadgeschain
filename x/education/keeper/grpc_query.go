package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/education/types"
)

var _ types.QueryServer = Keeper{}
