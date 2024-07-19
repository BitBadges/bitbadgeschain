package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/anchor module sentinel errors
var (
	ErrInvalidSigner        = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample               = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidAddress       = sdkerrors.Register(ModuleName, 1503, "invalid address")
	ErrInvalidRequest       = sdkerrors.Register(ModuleName, 1504, "invalid request")
	ErrUnknownRequest       = sdkerrors.Register(ModuleName, 1505, "unknown request")
)
