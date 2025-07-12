package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrCollectionNotExists                      = sdkerrors.Register(types.ModuleName, 4, "collection does not exist")
	ErrForbiddenTime                            = sdkerrors.Register(types.ModuleName, 5, "this permission is forbidden at this time")
	ErrUserBalanceNotExists                     = sdkerrors.Register(types.ModuleName, 8, "user balance is empty or does not exist")
	ErrSupplyEqualsZero                         = sdkerrors.Register(types.ModuleName, 13, "can't create a badge with zero supply")
	ErrSenderIsNotManager                       = sdkerrors.Register(types.ModuleName, 14, "sender of tx is not the manager of the badge. must be manager to access this privilege")
	ErrAccountCanNotEqualCreator                = sdkerrors.Register(types.ModuleName, 35, "account can not equal creator")
	ErrRootHashInvalid                          = sdkerrors.Register(types.ModuleName, 49, "root hash invalid")
	ErrDecodingHexString                        = sdkerrors.Register(types.ModuleName, 53, "couldn't decode hex string")
	ErrWrongBalancesType                        = sdkerrors.Register(types.ModuleName, 64, "wrong balances type ")
	ErrCollectionIsArchived                     = sdkerrors.Register(types.ModuleName, 67, "collection is currently archived (read-only)")
	ErrNotImplemented                           = sdkerrors.Register(types.ModuleName, 69, "not implemented")
	ErrInadequateApprovals                      = sdkerrors.Register(types.ModuleName, 71, "inadequate approvals")
	ErrInvalidAddressListId                     = sdkerrors.Register(types.ModuleName, 72, "invalid address list id")
	ErrAddressListNotFound                      = sdkerrors.Register(types.ModuleName, 73, "address list not found")
	ErrCircularDependency                       = sdkerrors.Register(types.ModuleName, 75, "circular dependency")
	ErrDisallowedTransfer                       = sdkerrors.Register(types.ModuleName, 76, "disallowed transfer")
	ErrNoValidSolutionForChallenge              = sdkerrors.Register(types.ModuleName, 77, "challenge failed")
	ErrExceedsThreshold                         = sdkerrors.Register(types.ModuleName, 78, "exceeds threshold")
	ErrCircularInheritance                      = sdkerrors.Register(types.ModuleName, 79, "circular inheritance")
	ErrAddressListAlreadyExists                 = sdkerrors.Register(types.ModuleName, 80, "address list already exists")
	ErrGlobalArchive                            = sdkerrors.Register(types.ModuleName, 81, "global halt is active")
	ErrNoMatchingChallengeForChallengeTrackerId = sdkerrors.Register(types.ModuleName, 82, "no matching challenge for challenge tracker id")
	ErrInvalidDenom                             = sdkerrors.Register(types.ModuleName, 83, "invalid denom")
	ErrOverrideTimestampNotAllowed              = sdkerrors.Register(types.ModuleName, 84, "override timestamp not allowed")
	ErrAccountNotFound                          = sdkerrors.Register(types.ModuleName, 85, "account not found")
	ErrInvalidLeafSigner                        = sdkerrors.Register(types.ModuleName, 86, "invalid leaf signer")
	ErrInvalidConversion                        = sdkerrors.Register(types.ModuleName, 87, "invalid conversion")
)
