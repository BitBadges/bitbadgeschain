package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var (
	ErrCollectionNotExists                          = sdkerrors.Register(types.ModuleName, 4, "collection does not exist")
	ErrForbiddenTime 																= sdkerrors.Register(types.ModuleName, 5, "this permission is forbidden at this time")
	ErrUserBalanceNotExists                         = sdkerrors.Register(types.ModuleName, 8, "user balance does not exist")
	ErrSupplyEqualsZero                             = sdkerrors.Register(types.ModuleName, 13, "can't create a badge with zero supply")
	ErrSenderIsNotManager                           = sdkerrors.Register(types.ModuleName, 14, "sender of tx is not the manager of the badge. must be manager to access this privilege")
	ErrAccountCanNotEqualCreator                    = sdkerrors.Register(types.ModuleName, 35, "account can not equal creator")
	ErrRootHashInvalid                              = sdkerrors.Register(types.ModuleName, 49, "root hash invalid")
	ErrDecodingHexString                            = sdkerrors.Register(types.ModuleName, 53, "couldn't decode hex string")
	ErrWrongBalancesType                             = sdkerrors.Register(types.ModuleName, 64, "wrong balances type (off-chain balances vs inherited balances vs on-chain balances) ")
	ErrCollectionIsArchived                         = sdkerrors.Register(types.ModuleName, 67, "collection is currently archived (read-only)")
	ErrNotImplemented															  = sdkerrors.Register(types.ModuleName, 69, "not implemented")
	ErrInadequateApprovals													= sdkerrors.Register(types.ModuleName, 71, "inadequate approvals")
	ErrInvalidAddressMappingId											= sdkerrors.Register(types.ModuleName, 72, "invalid address mapping id")
	ErrAddressMappingNotFound											  = sdkerrors.Register(types.ModuleName, 73, "address mapping not found")
	ErrCircularDependency 												 = sdkerrors.Register(types.ModuleName, 75, "circular dependency")
	ErrDisallowedTransfer 												 = sdkerrors.Register(types.ModuleName, 76, "disallowed transfer")
	ErrNoValidSolutionForChallenge 								 = sdkerrors.Register(types.ModuleName, 77, "no valid solution for challenge")
	ErrInheritedBalances 													 = sdkerrors.Register(types.ModuleName, 79, "inherited balances")
)
