package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/offers module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample                     = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidPacketTimeout       = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion             = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrProposalNotFound           = sdkerrors.Register(ModuleName, 1502, "proposal not found")
	ErrUnauthorized               = sdkerrors.Register(ModuleName, 1503, "unauthorized party")
	ErrProposalAlreadyAccepted    = sdkerrors.Register(ModuleName, 1504, "proposal already accepted")
	ErrProposalAlreadyRejected    = sdkerrors.Register(ModuleName, 1505, "proposal already rejected")
	ErrProposalNotAccepted        = sdkerrors.Register(ModuleName, 1506, "proposal not accepted")
	ErrProposalNotValidAtThisTime = sdkerrors.Register(ModuleName, 1507, "proposal not valid at this time")
)
