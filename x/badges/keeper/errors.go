package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/collections/types"
)

var (
	ErrInvalidNFT                        = sdkerrors.Register(types.ModuleName, 2, "invalid nft")
	ErrBadgeExists                       = sdkerrors.Register(types.ModuleName, 3, "badge already exists")
	ErrBadgeNotExists                    = sdkerrors.Register(types.ModuleName, 4, "badge does not exist")
	ErrSubBadgeExists                    = sdkerrors.Register(types.ModuleName, 5, "subbadge already exists")
	ErrSubBadgeNotExists                 = sdkerrors.Register(types.ModuleName, 6, "subbadge does not exist")
	ErrBadgeBalanceExists                = sdkerrors.Register(types.ModuleName, 7, "BadgeBalance already exists")
	ErrBadgeBalanceNotExists             = sdkerrors.Register(types.ModuleName, 8, "BadgeBalance does not exist")
	ErrInvalidBadgeID                    = sdkerrors.Register(types.ModuleName, 9, "invalid format of badge id")
	ErrInvalidPermissionsLeadingZeroes   = sdkerrors.Register(types.ModuleName, 10, "permissions does not start with correct amount of leading zeroes")
	ErrInvalidPermissionsUpdateLocked    = sdkerrors.Register(types.ModuleName, 11, "permission has previously been locked so cannot be updated")
	ErrInvalidPermissionsUpdatePermanent = sdkerrors.Register(types.ModuleName, 12, "permission is permanent and cannot be updated")
	ErrSupplyEqualsZero                  = sdkerrors.Register(types.ModuleName, 13, "supply can't equal zero")
	ErrSenderIsNotManager                = sdkerrors.Register(types.ModuleName, 14, "sender is not the manager of the badge. only the manager potentially has access to this privilege")
	ErrInvalidPermissions                = sdkerrors.Register(types.ModuleName, 15, "the badge permissions that are set do not allow this action")
	ErrBalanceIsZero                     = sdkerrors.Register(types.ModuleName, 16, "the balance to add can't be zero")
	ErrInvalidUri                        = sdkerrors.Register(types.ModuleName, 17, "invalid format of uri")
	ErrBadgeBalanceTooLow                = sdkerrors.Register(types.ModuleName, 18, "the badge balance is too low to perform this action")
	ErrSenderAndReceiverSame             = sdkerrors.Register(types.ModuleName, 19, "sender and receiver cannot be the same")
	ErrCantAcceptOwnTransferRequest      = sdkerrors.Register(types.ModuleName, 20, "cannot accept own transfer request")
	ErrInsufficientApproval              = sdkerrors.Register(types.ModuleName, 21, "insufficient approval")
	ErrAccountsAreNotRegistered          = sdkerrors.Register(types.ModuleName, 22, "accounts are not registered")
	ErrNoPendingTransferFound            = sdkerrors.Register(types.ModuleName, 23, "no pending transfer found")
	ErrPendingNotFound                   = sdkerrors.Register(types.ModuleName, 24, "pending transfer not found")
	ErrOverflow                          = sdkerrors.Register(types.ModuleName, 25, "overflow")
	ErrAddressFrozen                     = sdkerrors.Register(types.ModuleName, 26, "address is frozen")
	ErrAddressAlreadyFrozen			  	= sdkerrors.Register(types.ModuleName, 27, "address is already frozen")
	ErrAddressNotFrozen				  	= sdkerrors.Register(types.ModuleName, 28, "address is not frozen")
)
