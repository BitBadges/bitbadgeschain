package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgDeleteManagerSplitter) ValidateBasic() error {
	// Validate admin address format
	adminAddr, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAdmin, "invalid admin address (%s)", err)
	}
	if adminAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidAdmin, "admin address cannot be empty")
	}

	// Validate manager splitter address format
	if msg.Address == "" {
		return sdkerrors.Wrap(ErrInvalidAddress, "manager splitter address cannot be empty")
	}
	managerSplitterAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid manager splitter address format (%s)", err)
	}
	if managerSplitterAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidAddress, "manager splitter address cannot be empty")
	}

	// Validate that admin and manager splitter are different addresses
	if msg.Admin == msg.Address {
		return sdkerrors.Wrap(ErrInvalidRequest, "admin and manager splitter address cannot be the same")
	}

	return nil
}

