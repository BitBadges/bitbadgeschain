package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/badges module sentinel errors
var (
	ErrSample                                      = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout                        = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion                              = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidCollectionID                         = sdkerrors.Register(ModuleName, 1502, "invalid collection ID")
	ErrInvalidURI                                  = sdkerrors.Register(ModuleName, 1503, "invalid URI. must be blank or a correctly formatted URI")
	ErrInvalidPermissions                          = sdkerrors.Register(ModuleName, 1504, "invalid permissions")
	ErrAmountEqualsZero                            = sdkerrors.Register(ModuleName, 1505, "amount cannot equal zero")
	ErrDuplicateAddresses                          = sdkerrors.Register(ModuleName, 1506, "duplicate addresses")
	ErrStartGreaterThanEnd                         = sdkerrors.Register(ModuleName, 1507, "start greater than end")
	ErrRangesIsNil                                 = sdkerrors.Register(ModuleName, 1508, "ranges is nil")
	ErrElementCantEqualThis                        = sdkerrors.Register(ModuleName, 1509, "element cant equal this")
	ErrInvalidUintRangeSpecified                   = sdkerrors.Register(ModuleName, 1510, "invalid id range specified")
	ErrInvalidTypedData                            = sdkerrors.Register(ModuleName, 1511, "invalid typed data")
	ErrNotImplemented                              = sdkerrors.Register(ModuleName, 1512, "not implemented")
	ErrRangesOverlap                               = sdkerrors.Register(ModuleName, 1513, "id ranges overlap. for example, { Start: sdkmath.NewUint(1), end: 5 } and { Start: sdkmath.NewUint(4), End: sdkmath.NewUint(10) } overlap")
	ErrUintUnititialized                           = sdkerrors.Register(ModuleName, 1514, "uint is uninitialized (nil)")
	ErrPrimaryChallengeMustBeOneUsePerLeaf         = sdkerrors.Register(ModuleName, 1515, "primary challenge must be one use per leaf")
	ErrCanOnlyUseMaxOneUsePerLeafWithWhitelistTree = sdkerrors.Register(ModuleName, 1516, "can only use non max one use per leaf with whitelist tree")
	ErrCanOnlyUseLeafIndexForBadgeIdsOnce          = sdkerrors.Register(ModuleName, 1517, "can only use leaf index for badge ids once")
	ErrRangeDoesNotOverlap                         = sdkerrors.Register(ModuleName, 1518, "range does not overlap with existing ranges")
	ErrAmountRestrictionsIsNil                     = sdkerrors.Register(ModuleName, 1519, "amount restrictions is nil")
	ErrPermissionsValueIsNil                       = sdkerrors.Register(ModuleName, 1520, "permissions is defined but default values is nil")
	ErrCombinationsIsNil                           = sdkerrors.Register(ModuleName, 1521, "permissions is defined but combinations is nil")
	ErrPermissionsIsNil                            = sdkerrors.Register(ModuleName, 1522, "permissions is nil")
	ErrInvalidCombinations                         = sdkerrors.Register(ModuleName, 1523, "invalid permission combinations. you have specified duplicate combinations and because of the first match policy, the second combination will never be used. please remove the duplicate combinations")
	ErrOverflow                                    = sdkerrors.Register(ModuleName, 1524, "overflow")
	ErrUnderflow                                   = sdkerrors.Register(ModuleName, 1525, "underflow")
	ErrInvalidAddress                              = sdkerrors.Register(ModuleName, 1526, "invalid address")
	ErrInvalidRequest                              = sdkerrors.Register(ModuleName, 1527, "invalid request")
	ErrUnknownRequest                              = sdkerrors.Register(ModuleName, 1528, "unknown request")
	ErrInvalidType                                 = sdkerrors.Register(ModuleName, 1529, "invalid type")
	ErrUnauthorized                                = sdkerrors.Register(ModuleName, 1530, "unauthorized")
	ErrInvalidPubKey                               = sdkerrors.Register(ModuleName, 1531, "invalid public key")
	ErrWrongSequence                               = sdkerrors.Register(ModuleName, 1532, "wrong sequence")
	ErrNotSupported                                = sdkerrors.Register(ModuleName, 1533, "not supported")
	ErrTooManySignatures                           = sdkerrors.Register(ModuleName, 1534, "too many signatures")
	ErrNoSignatures                                = sdkerrors.Register(ModuleName, 1535, "no signatures")
	ErrUnknownExtensionOptions                     = sdkerrors.Register(ModuleName, 1536, "unknown extension options")
	ErrInvalidChainID                              = sdkerrors.Register(ModuleName, 1537, "invalid chain id")
	ErrorInvalidSigner                             = sdkerrors.Register(ModuleName, 1538, "invalid signer")
	ErrLogic                                       = sdkerrors.Register(ModuleName, 1539, "logic")
	ErrNotFound                                    = sdkerrors.Register(ModuleName, 1540, "not found")
	ErrInvalidInheritedBadgeLength                 = sdkerrors.Register(ModuleName, 1541, "invalid inherited badge balances length. num parent badges must == 1 or equal num child badges")
	ErrNoTimelineTimeSpecified                     = sdkerrors.Register(ModuleName, 1542, "no timeline times specified (len 0 for Times)")
	ErrSenderAndReceiverSame                       = sdkerrors.Register(ModuleName, 1543, "sender and receiver cannot be the same")
	ErrInvalidTransfers                            = sdkerrors.Register(ModuleName, 1544, "invalid transfers")
	ErrUintGreaterThanMax                          = sdkerrors.Register(ModuleName, 1545, "uint greater than max uint")
	ErrExceedsThreshold														= sdkerrors.Register(ModuleName, 1546, "exceeds threshold")
	ErrApprovalTrackerIdIsNil 										= sdkerrors.Register(ModuleName, 1547, "approval tracker id is nil")
	ErrChallengeTrackerIdIsNil 										= sdkerrors.Register(ModuleName, 1548, "challenge tracker id is nil")
	ErrIdsContainsInvalidChars 										= sdkerrors.Register(ModuleName, 1549, "ids contains invalid chars")
)
