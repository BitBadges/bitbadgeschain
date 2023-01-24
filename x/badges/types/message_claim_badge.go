package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimBadge = "claim_badge"

var _ sdk.Msg = &MsgClaimBadge{}

func NewMsgClaimBadge(creator string, claimId uint64, collectionId uint64, leaf []byte, proof *Proof) *MsgClaimBadge {
	return &MsgClaimBadge{
		Creator: creator,
		ClaimId: claimId,
		CollectionId: collectionId,
		Leaf: leaf,
		Proof: proof,
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
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Leaf == nil || len(msg.Leaf) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid leaf")
	}

	if msg.Proof.LeafHash == nil || len(msg.Proof.LeafHash) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid leaf hash in proof")
	}

	if msg.Proof.Aunts == nil || len(msg.Proof.Aunts) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid aunts in proof")
	}

	for _, aunt := range msg.Proof.Aunts {
		if aunt == nil || len(aunt) == 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid aunt in proof")
		}
	}
	

	return nil
}
