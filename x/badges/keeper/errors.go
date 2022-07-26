package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/collections/types"
)

// x/nft module sentinel errors
var (
	// ErrInvalidNFT     = sdkerrors.Register(ModuleName, 2, "invalid nft")
	ErrBadgeExists    = sdkerrors.Register(types.ModuleName, 3, "badge already exists")
	ErrBadgeNotExists = sdkerrors.Register(types.ModuleName, 4, "badge does not exist")
	// ErrNFTExists      = sdkerrors.Register(ModuleName, 5, "nft already exist")
	// ErrNFTNotExists   = sdkerrors.Register(ModuleName, 6, "nft does not exist")
	// ErrInvalidID      = sdkerrors.Register(ModuleName, 7, "invalid id")
	ErrInvalidBadgeID = sdkerrors.Register(types.ModuleName, 8, "invalid badge id")
)
