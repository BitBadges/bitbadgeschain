package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimBadge = "claim_badge"

var _ sdk.Msg = &MsgClaimBadge{}

func NewMsgClaimBadge(creator string, claimId uint64, collectionId uint64, solutions []*ChallengeSolution) *MsgClaimBadge {
	return &MsgClaimBadge{
		Creator:        creator,
		ClaimId:        claimId,
		CollectionId:   collectionId,
		Solutions:      solutions,
	}
}

func (msg *MsgClaimBadge) Route() string {
	return RouterKey
}

func (msg *MsgClaimBadge) Type() string {
	return TypeMsgClaimBadge
}

func (msg *MsgClaimBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimBadge) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.ClaimId == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid claim id")
	}

	if msg.CollectionId == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid collection id")
	}

	for _, solution := range msg.Solutions {
		if solution.Proof != nil {
			for _, aunt := range solution.Proof.Aunts {
				if aunt.Aunt == "" || len(aunt.Aunt) == 0 {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid aunt in proof")
				}
			}
		}
	}

	return nil
}
