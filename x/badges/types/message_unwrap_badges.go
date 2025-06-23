package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUnwrapBadges = "unwrap_badge"

var _ sdk.Msg = &MsgUnwrapBadges{}

func NewMsgUnwrapBadges(creator string, collectionId sdkmath.Uint) *MsgUnwrapBadges {
	return &MsgUnwrapBadges{
		Creator:      creator,
		CollectionId: collectionId,
	}
}

func (msg *MsgUnwrapBadges) Route() string {
	return RouterKey
}

func (msg *MsgUnwrapBadges) Type() string {
	return TypeMsgUnwrapBadges
}

func (msg *MsgUnwrapBadges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnwrapBadges) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnwrapBadges) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.Denom == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid denom")
	}

	if msg.Amount.IsNil() || msg.Amount.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid amount")
	}

	return nil
}

func (msg *MsgUnwrapBadges) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}
