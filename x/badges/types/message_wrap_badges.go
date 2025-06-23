package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWrapBadges = "wrap_badge"

var _ sdk.Msg = &MsgWrapBadges{}

func NewMsgWrapBadges(creator string, collectionId sdkmath.Uint) *MsgWrapBadges {
	return &MsgWrapBadges{
		Creator:      creator,
		CollectionId: collectionId,
	}
}

func (msg *MsgWrapBadges) Route() string {
	return RouterKey
}

func (msg *MsgWrapBadges) Type() string {
	return TypeMsgWrapBadges
}

func (msg *MsgWrapBadges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWrapBadges) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWrapBadges) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
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

	msg.Balances, err = ValidateBalances(ctx, msg.Balances, canChangeValues)
	if err != nil {
		return err
	}

	return nil
}

func (msg *MsgWrapBadges) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}
