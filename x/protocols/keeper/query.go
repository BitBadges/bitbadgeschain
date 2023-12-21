package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
)

var _ types.QueryServer = Keeper{}
