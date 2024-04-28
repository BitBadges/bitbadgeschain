package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/maps module sentinel errors
var (
	ErrSample               = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrPermissionsIsNil     = sdkerrors.Register(ModuleName, 1502, "permissions is nil")
	ErrInvalidAddress       = sdkerrors.Register(ModuleName, 1503, "invalid address")
	ErrInvalidRequest       = sdkerrors.Register(ModuleName, 1504, "invalid request")
	ErrUnknownRequest       = sdkerrors.Register(ModuleName, 1505, "unknown request")
)
