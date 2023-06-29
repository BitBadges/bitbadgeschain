package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateManager = "transfer_manager"

var _ sdk.Msg = &MsgUpdateManager{}

func NewMsgUpdateManager(creator string, collectionId sdkmath.Uint, managerTimeline []*ManagerTimeline) *MsgUpdateManager {
	return &MsgUpdateManager{
		Creator:      creator,
		CollectionId: collectionId,
		ManagerTimeline: managerTimeline,	
	}
}

func (msg *MsgUpdateManager) Route() string {
	return RouterKey
}

func (msg *MsgUpdateManager) Type() string {
	return TypeMsgUpdateManager
}

func (msg *MsgUpdateManager) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateManager) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateManager) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	for _, timelineVal := range msg.ManagerTimeline {	
		_, err = sdk.AccAddressFromBech32(timelineVal.Manager)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid provided address (%s)", err)
		}
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	return nil
}
