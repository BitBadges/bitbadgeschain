package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/badges module sentinel errors
var (
	ErrSample                                = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout                  = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion                        = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidBadgeID                        = sdkerrors.Register(ModuleName, 1502, "invalid badge ID")
	ErrInvalidBadgeURI                       = sdkerrors.Register(ModuleName, 1503, "invalid badge URI")
	ErrInvalidPermissionsLeadingZeroes       = sdkerrors.Register(ModuleName, 1504, "invalid permissions leading zeroes")
	ErrInvalidPermissions                    = sdkerrors.Register(ModuleName, 1505, "invalid permissions")
	ErrInvalidPermissionsUpdateLocked        = sdkerrors.Register(ModuleName, 1506, "invalid permissions update locked")
	ErrInvalidPermissionsUpdatePermanent     = sdkerrors.Register(ModuleName, 1507, "invalid permissions update permanent")
	ErrSupplyEqualsZero                      = sdkerrors.Register(ModuleName, 1508, "supply equals zero")
	ErrSenderAndReceiverSame                 = sdkerrors.Register(ModuleName, 1509, "sender and receiver same")
	ErrInvalidSupplyAndAmounts               = sdkerrors.Register(ModuleName, 1510, "invalid supply and amounts")
	ErrAmountEqualsZero                      = sdkerrors.Register(ModuleName, 1511, "amount to create equals zero")
	ErrInvalidAmountsAndAddressesLength      = sdkerrors.Register(ModuleName, 1512, "invalid amounts and addresses length")
	ErrInvalidBadgeHash                      = sdkerrors.Register(ModuleName, 1513, "invalid badge hash")
	ErrDuplicateAddresses                    = sdkerrors.Register(ModuleName, 1514, "duplicate addresses")
	ErrStartGreaterThanEnd                   = sdkerrors.Register(ModuleName, 1515, "start greater than end")
	ErrDefaultSupplyEqualsZero               = sdkerrors.Register(ModuleName, 1516, "default supply equals zero")
	ErrInvalidArgumentLengths                = sdkerrors.Register(ModuleName, 1517, "invalid argument lengths")
	ErrRangesIsNil                           = sdkerrors.Register(ModuleName, 1518, "ranges is nil")
	ErrBytesGreaterThan256                   = sdkerrors.Register(ModuleName, 1519, "bytes greater than 256")
	ErrInvalidUriScheme                      = sdkerrors.Register(ModuleName, 1520, "invalid uri scheme")
	ErrCancelTimeIsGreaterThanExpirationTime = sdkerrors.Register(ModuleName, 1521, "cancel time is greater than expiration time")
	ErrDuplicateAmounts                      = sdkerrors.Register(ModuleName, 1522, "duplicate amounts")
	ErrElementCantEqualThis                  = sdkerrors.Register(ModuleName, 1523, "element cant equal this")
	ErrInvalidIdRangeSpecified			   = sdkerrors.Register(ModuleName, 1524, "invalid id range specified")
	ErrInvalidTypedData					  = sdkerrors.Register(ModuleName, 1525, "invalid typed data")
	ErrActionOutOfRange					  = sdkerrors.Register(ModuleName, 1526, "action out of range")
	ErrActionsEmpty					  = sdkerrors.Register(ModuleName, 1527, "actions empty")
	ErrActionsLengthNotEqualToRangesLength	  = sdkerrors.Register(ModuleName, 1528, "actions length not equal to ranges length")

)
