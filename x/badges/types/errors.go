package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/badges module sentinel errors
var (
	ErrSample                            = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout              = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion                    = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidBadgeID                    = sdkerrors.Register(ModuleName, 1502, "invalid badge ID")
	ErrInvalidBadgeURI                   = sdkerrors.Register(ModuleName, 1503, "invalid badge URI")
	ErrInvalidPermissionsLeadingZeroes   = sdkerrors.Register(ModuleName, 1504, "invalid permissions leading zeroes")
	ErrInvalidPermissions                = sdkerrors.Register(ModuleName, 1505, "invalid permissions")
	ErrInvalidPermissionsUpdateLocked    = sdkerrors.Register(ModuleName, 1506, "invalid permissions update locked")
	ErrInvalidPermissionsUpdatePermanent = sdkerrors.Register(ModuleName, 1507, "invalid permissions update permanent")
	ErrSupplyEqualsZero                  = sdkerrors.Register(ModuleName, 1508, "supply equals zero")
	ErrSenderAndReceiverSame             = sdkerrors.Register(ModuleName, 1509, "sender and receiver same")
	ErrInvalidSupplyAndAmounts           = sdkerrors.Register(ModuleName, 1510, "invalid supply and amounts")
	ErrAmountEqualsZero                  = sdkerrors.Register(ModuleName, 1511, "amount to create equals zero")
	ErrInvalidAmountsAndAddressesLength  = sdkerrors.Register(ModuleName, 1512, "invalid amounts and addresses length")
	ErrInvalidBadgeHash				  	 = sdkerrors.Register(ModuleName, 1513, "invalid badge hash")
)
