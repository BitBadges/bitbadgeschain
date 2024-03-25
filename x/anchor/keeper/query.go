package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
)

var _ types.QueryServer = Keeper{}
