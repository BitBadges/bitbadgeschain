package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

var _ types.QueryServer = Keeper{}
