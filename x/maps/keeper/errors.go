package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrMapExists             = sdkerrors.Register(types.ModuleName, 1, "map already exists")
	ErrMapDoesNotExist       = sdkerrors.Register(types.ModuleName, 2, "map does not exist")
	ErrNotMapCreator         = sdkerrors.Register(types.ModuleName, 3, "not map creator")
	ErrMapIsFrozen           = sdkerrors.Register(types.ModuleName, 4, "map is frozen")
	ErrMapNotEditable        = sdkerrors.Register(types.ModuleName, 5, "map is not editable at the moment")
	ErrValueAlreadySet       = sdkerrors.Register(types.ModuleName, 6, "value already set")
	ErrCannotUpdateMapValues = sdkerrors.Register(types.ModuleName, 7, "cannot update map values")
	ErrCannotUpdateMapValue  = sdkerrors.Register(types.ModuleName, 8, "cannot update map value")
	ErrDuplicateValue        = sdkerrors.Register(types.ModuleName, 9, "duplicate value")
	ErrInvalidMapId          = sdkerrors.Register(types.ModuleName, 10, "invalid map id")
	ErrInvalidValue          = sdkerrors.Register(types.ModuleName, 11, "invalid value")
)
