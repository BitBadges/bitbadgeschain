package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/managersplitter module sentinel errors
var (
	ErrManagerSplitterNotFound    = sdkerrors.Register(ModuleName, 1100, "manager splitter not found")
	ErrUnauthorized              = sdkerrors.Register(ModuleName, 1101, "unauthorized: must be admin")
	ErrInvalidAddress             = sdkerrors.Register(ModuleName, 1102, "invalid address")
	ErrInvalidRequest             = sdkerrors.Register(ModuleName, 1103, "invalid request")
	ErrUnknownRequest             = sdkerrors.Register(ModuleName, 1104, "unknown request")
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1105, "invalid signer")
	ErrPermissionDenied           = sdkerrors.Register(ModuleName, 1106, "permission denied: address not approved for this action")
	ErrInvalidAdmin               = sdkerrors.Register(ModuleName, 1107, "invalid admin address")
	ErrManagerSplitterExists      = sdkerrors.Register(ModuleName, 1108, "manager splitter already exists")
)

