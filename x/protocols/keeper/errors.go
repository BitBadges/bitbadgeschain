package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
)

var (
	ErrProtocolExists       = sdkerrors.Register(types.ModuleName, 1, "protocol already exists")
	ErrProtocolDoesNotExist = sdkerrors.Register(types.ModuleName, 2, "protocol does not exist")
	ErrNotProtocolCreator   = sdkerrors.Register(types.ModuleName, 3, "not protocol creator")
)
