package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/council module sentinel errors
var (
	ErrInvalidAddress       = sdkerrors.Register(ModuleName, 1100, "invalid address")
	ErrInvalidThreshold     = sdkerrors.Register(ModuleName, 1101, "voting threshold must be between 1 and 100")
	ErrCouncilNotFound      = sdkerrors.Register(ModuleName, 1102, "council not found")
	ErrProposalNotFound     = sdkerrors.Register(ModuleName, 1103, "proposal not found")
	ErrNoCredential         = sdkerrors.Register(ModuleName, 1104, "sender does not hold the required credential token")
	ErrDisallowedMsgType   = sdkerrors.Register(ModuleName, 1105, "message type not allowed by this council")
	ErrProposalNotPassed    = sdkerrors.Register(ModuleName, 1106, "proposal has not passed")
	ErrExecutionDelayNotMet = sdkerrors.Register(ModuleName, 1107, "execution delay has not elapsed")
	ErrProposalExpired      = sdkerrors.Register(ModuleName, 1108, "proposal voting deadline has passed")
	ErrInvalidMsgCount      = sdkerrors.Register(ModuleName, 1109, "msg type URLs and msg bytes must have the same length")
	ErrMsgDispatchFailed    = sdkerrors.Register(ModuleName, 1110, "message dispatch failed")
	ErrInvalidExecutionDelay = sdkerrors.Register(ModuleName, 1111, "execution delay must be non-negative")
	ErrInvalidDeadline      = sdkerrors.Register(ModuleName, 1112, "deadline must be in the future")
	ErrAlreadyExecuted      = sdkerrors.Register(ModuleName, 1113, "proposal already executed")
	ErrNoMessages           = sdkerrors.Register(ModuleName, 1114, "proposal must contain at least one message")
)
