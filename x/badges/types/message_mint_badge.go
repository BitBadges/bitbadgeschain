package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintBadge = "new_sub_badge"

var _ sdk.Msg = &MsgMintBadge{}

func NewMsgMintBadge(creator string, collectionId uint64, supplysAndAmounts []*BadgeSupplyAndAmount, transfers []*Transfers, claims []*Claim) *MsgMintBadge {
	return &MsgMintBadge{
		Creator:      creator,
		CollectionId: collectionId,
		BadgeSupplys: supplysAndAmounts,
		Transfers:    transfers,
		Claims:       claims,
	}
}

func (msg *MsgMintBadge) Route() string {
	return RouterKey
}

func (msg *MsgMintBadge) Type() string {
	return TypeMsgMintBadge
}

func (msg *MsgMintBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintBadge) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	amounts := make([]uint64, len(msg.BadgeSupplys))
	supplys := make([]uint64, len(msg.BadgeSupplys))
	for i, subasset := range msg.BadgeSupplys {
		amounts[i] = subasset.Amount
		supplys[i] = subasset.Supply
	}

	err = ValidateNoElementIsX(amounts, 0)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(supplys, 0)
	if err != nil {
		return err
	}

	for _, claim := range msg.Claims {
		if claim.Uri != "" {
			err = ValidateURI(claim.Uri)
			if err != nil {
				return err
			}
		}
	
		if claim.TimeRange == nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
		}
	
		err = ValidateRangesAreValid([]*IdRange{claim.TimeRange})
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
		}

		if claim.Type == uint64(ClaimType_MerkleTree) {
			if claim.Data == "" || len(claim.Data) == 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid claim data")
			}
		}

		if claim.Type == uint64(ClaimType_FirstCome) {
			if claim.BadgeIds == nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
			}
		
			err = ValidateRangesAreValid(claim.BadgeIds)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
			}
		}
	}

	return nil
}
