package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrRateLimitExceeded = errorsmod.Register(ModuleName, 1, "rate limit exceeded: transfer would exceed maximum supply change percentage")
	ErrInvalidSigner     = errorsmod.Register(ModuleName, 2, "expected gov account as only signer for proposal message")
)
