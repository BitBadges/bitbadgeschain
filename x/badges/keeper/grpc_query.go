package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var _ types.QueryServer = Keeper{}
