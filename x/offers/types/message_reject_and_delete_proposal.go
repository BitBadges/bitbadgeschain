package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRejectAndDeleteProposal{}

func NewMsgRejectAndDeleteProposal(creator string, id sdkmath.Uint) *MsgRejectAndDeleteProposal {
	return &MsgRejectAndDeleteProposal{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgRejectAndDeleteProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Id.IsNil() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid id: %s", msg.Id)
	}

	if msg.Id.IsZero() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid id: %s", msg.Id)
	}

	return nil
}
