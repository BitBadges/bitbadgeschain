package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// We can clean this up in the future. A lot of errors are deprecated and can be removed.
var (
	ErrInvalidNFT                                   = sdkerrors.Register(types.ModuleName, 2, "invalid nft")
	ErrCollectionExists                             = sdkerrors.Register(types.ModuleName, 3, "collection already exists")
	ErrCollectionNotExists                          = sdkerrors.Register(types.ModuleName, 4, "collection does not exist")
	ErrBadgeExists                                  = sdkerrors.Register(types.ModuleName, 5, "badge already exists")
	ErrBadgeNotExists                               = sdkerrors.Register(types.ModuleName, 6, "badge does not exist")
	ErrUserBalanceExists                            = sdkerrors.Register(types.ModuleName, 7, "user balance already exists")
	ErrUserBalanceNotExists                         = sdkerrors.Register(types.ModuleName, 8, "user balance does not exist")
	ErrInvalidBadgeID                               = sdkerrors.Register(types.ModuleName, 9, "invalid badge id")
	ErrInvalidPermissionsLeadingZeroes              = sdkerrors.Register(types.ModuleName, 10, "permissions does not start with correct amount of leading zeroes")
	ErrInvalidPermissionsUpdateLocked               = sdkerrors.Register(types.ModuleName, 11, "permission has previously been locked so cannot be updated")
	ErrInvalidPermissionsUpdatePermanent            = sdkerrors.Register(types.ModuleName, 12, "permission is permanent and cannot be updated")
	ErrSupplyEqualsZero                             = sdkerrors.Register(types.ModuleName, 13, "supply can't equal zero")
	ErrSenderIsNotManager                           = sdkerrors.Register(types.ModuleName, 14, "sender is not the manager of the badge. only the manager has access to this privilege")
	ErrInvalidPermissions                           = sdkerrors.Register(types.ModuleName, 15, "the badge permissions that are set do not allow this action")
	ErrBalanceIsZero                                = sdkerrors.Register(types.ModuleName, 16, "the balance to add can't be zero")
	ErrInvalidUri                                   = sdkerrors.Register(types.ModuleName, 17, "invalid format of uri")
	ErrUserBalanceTooLow                            = sdkerrors.Register(types.ModuleName, 18, "the badge balance is too low to perform this action")
	ErrSenderAndReceiverSame                        = sdkerrors.Register(types.ModuleName, 19, "sender and receiver cannot be the same")
	ErrCantAcceptOwnTransferRequest                 = sdkerrors.Register(types.ModuleName, 20, "cannot accept own transfer request")
	ErrInsufficientApproval                         = sdkerrors.Register(types.ModuleName, 21, "insufficient approval")
	ErrAccountsAreNotRegistered                     = sdkerrors.Register(types.ModuleName, 22, "accounts are not registered")
	ErrNoPendingTransferFound                       = sdkerrors.Register(types.ModuleName, 23, "no pending transfer found")
	ErrPendingNotFound                              = sdkerrors.Register(types.ModuleName, 24, "pending transfer not found")
	ErrOverflow                                     = sdkerrors.Register(types.ModuleName, 25, "overflow")
	ErrAddressFrozen                                = sdkerrors.Register(types.ModuleName, 26, "address is frozen")
	ErrAddressAlreadyFrozen                         = sdkerrors.Register(types.ModuleName, 27, "address is already frozen")
	ErrAddressNotFrozen                             = sdkerrors.Register(types.ModuleName, 28, "address is not frozen")
	ErrAddressNeedsToOptInAndRequestManagerTransfer = sdkerrors.Register(types.ModuleName, 29, "address needs to opt in and request manager transfer")
	ErrMustOwnTotalSupplyToSelfDestruct             = sdkerrors.Register(types.ModuleName, 30, "must own total supply to self destruct")
	ErrBadgeCanNotBeSelfDestructed                  = sdkerrors.Register(types.ModuleName, 31, "badge can not be self destructed")
	ErrNotApproved                                  = sdkerrors.Register(types.ModuleName, 32, "not approved")
	ErrPendingTransferExpired                       = sdkerrors.Register(types.ModuleName, 33, "pending transfer expired")
	ErrInvalidBadgeRange                            = sdkerrors.Register(types.ModuleName, 34, "invalid badge range")
	ErrAccountCanNotEqualCreator                    = sdkerrors.Register(types.ModuleName, 35, "account can not equal creator")
	ErrCantPruneBalanceYet                          = sdkerrors.Register(types.ModuleName, 36, "cant prune balance yet")
	ErrCantCancelYet                                = sdkerrors.Register(types.ModuleName, 37, "cant cancel yet")
	ErrCancelTimeIsGreaterThanExpirationTime        = sdkerrors.Register(types.ModuleName, 38, "cancel time is greater than expiration time")
	ErrApprovalForAddressDoesntExist                = sdkerrors.Register(types.ModuleName, 39, "approval for address doesn't exist")
	ErrUnderflow                                    = sdkerrors.Register(types.ModuleName, 40, "error underflow")
	ErrAccountNotRegistered                         = sdkerrors.Register(types.ModuleName, 41, "account not registered")
	ErrClaimNotExists                               = sdkerrors.Register(types.ModuleName, 42, "claim does not exist")
	ErrClaimAlreadyUsed                             = sdkerrors.Register(types.ModuleName, 43, "claim already used")
	ErrClaimNotFound                                = sdkerrors.Register(types.ModuleName, 44, "claim not found")
	ErrClaimDataInvalid                             = sdkerrors.Register(types.ModuleName, 45, "claim data invalid")
	ErrClaimTimeInvalid                             = sdkerrors.Register(types.ModuleName, 46, "claim time invalid")
	ErrIdAlreadyInRanges                            = sdkerrors.Register(types.ModuleName, 47, "id already in ranges")
	ErrIdInRange                                    = sdkerrors.Register(types.ModuleName, 48, "id in ranges")
	ErrRootHashInvalid                              = sdkerrors.Register(types.ModuleName, 49, "root hash invalid")
	ErrInvalidAddress                               = sdkerrors.Register(types.ModuleName, 50, "invalid address")
	ErrLeafIsEmpty                                  = sdkerrors.Register(types.ModuleName, 51, "leaf is empty")
	ErrAuntsIsEmpty                                 = sdkerrors.Register(types.ModuleName, 52, "aunts is empty")
	ErrDecodingHexString                            = sdkerrors.Register(types.ModuleName, 53, "decoding hex string")
	ErrMustBeClaimee                                = sdkerrors.Register(types.ModuleName, 54, "must be claimee")
	ErrCodeLeafInvalid                              = sdkerrors.Register(types.ModuleName, 55, "code leaf invalid")
	ErrCodeMaxUsesExceeded                          = sdkerrors.Register(types.ModuleName, 56, "code max uses exceeded")
	ErrAddressAlreadyUsed                           = sdkerrors.Register(types.ModuleName, 57, "address already used")
	ErrProofLengthInvalid                           = sdkerrors.Register(types.ModuleName, 58, "proof length invalid")
	ErrAddressMaxUsesExceeded                       = sdkerrors.Register(types.ModuleName, 59, "address max uses exceeded")
	ErrSolutionsLengthInvalid                       = sdkerrors.Register(types.ModuleName, 60, "solutions length invalid. must match number of challenges for claim")
	ErrChallengeMaxUsesExceeded                     = sdkerrors.Register(types.ModuleName, 61, "challenge max uses exceeded")
	ErrMintNotAllowed                               = sdkerrors.Register(types.ModuleName, 62, "minting to this address is not allowed according to the collection's allowed transfers")
	ErrInvalidTime                                  = sdkerrors.Register(types.ModuleName, 63, "claim is not valid at this time")
	ErrOffChainBalances                             = sdkerrors.Register(types.ModuleName, 64, "this action is not supported for this balances type (off-chain balances or inherited balances)")
	ErrClaimNotAssignable                           = sdkerrors.Register(types.ModuleName, 65, "claim is not assignable")
	ErrBadgeMetadataMustBeFrozen                    = sdkerrors.Register(types.ModuleName, 66, "badge metadata must be frozen")
	ErrCollectionIsArchived                         = sdkerrors.Register(types.ModuleName, 67, "collection is archived")
	ErrApprovedTransfersMustBeFrozen                = sdkerrors.Register(types.ModuleName, 68, "allowed transfers must be frozen")
	ErrNotImplemented															 = sdkerrors.Register(types.ModuleName, 69, "not implemented")
	ErrBadgeIdTooHigh															 = sdkerrors.Register(types.ModuleName, 70, "badge id too high. for inherited balances, all badge ids specified must be less than the collection's next badge ID - 1.")
	ErrInadequateApprovals													 = sdkerrors.Register(types.ModuleName, 71, "inadequate approvals")
)
