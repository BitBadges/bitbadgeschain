package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPruneBalances = "prune_balances"

var _ sdk.Msg = &MsgPruneBalances{}

func NewMsgPruneBalances(creator string, badgeIds []uint64, addresses []uint64) *MsgPruneBalances {
	return &MsgPruneBalances{
		Creator:   creator,
		BadgeIds:  badgeIds,
		Addresses: addresses,
	}
}

func (msg *MsgPruneBalances) Route() string {
	return RouterKey
}

func (msg *MsgPruneBalances) Type() string {
	return TypeMsgPruneBalances
}

func (msg *MsgPruneBalances) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPruneBalances) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPruneBalances) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
