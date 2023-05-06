package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintAndDistributeBadges = "mint_and_distribute_badge"

var _ sdk.Msg = &MsgMintAndDistributeBadges{}

func NewMsgMintAndDistributeBadges(creator string, collectionId uint64, supplysAndAmounts []*BadgeSupplyAndAmount, transfers []*Transfers, claims []*Claim, collectionUri string, badgeUris []*BadgeUri, balancesUri string) *MsgMintAndDistributeBadges {
	return &MsgMintAndDistributeBadges{
		Creator:       creator,
		CollectionId:  collectionId,
		BadgeSupplys:  supplysAndAmounts,
		Transfers:     transfers,
		Claims:        claims,
		CollectionUri: collectionUri,
		BadgeUris:     badgeUris,
		BalancesUri:   balancesUri,
	}
}

func (msg *MsgMintAndDistributeBadges) Route() string {
	return RouterKey
}

func (msg *MsgMintAndDistributeBadges) Type() string {
	return TypeMsgMintAndDistributeBadges
}

func (msg *MsgMintAndDistributeBadges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintAndDistributeBadges) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintAndDistributeBadges) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid collection id")
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

	if msg.BalancesUri != "" {
		err = ValidateURI(msg.BalancesUri)
		if err != nil {
			return err
		}
	}

	if msg.BadgeUris != nil && len(msg.BadgeUris) > 0 {
		for _, badgeUri := range msg.BadgeUris {
			//Validate well-formedness of the message entries
			if err := ValidateURI(badgeUri.Uri); err != nil {
				return err
			}

			err = ValidateRangesAreValid(badgeUri.BadgeIds)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid badgeIds")
			}
		}
	}

	if msg.CollectionUri != "" {
		err = ValidateURI(msg.CollectionUri)
		if err != nil {
			return err
		}
	}

	//TODO: Validate transfers


	for _, claim := range msg.Claims {
		err = ValidateClaim(claim)
		if err != nil {
			return err
		}
	}

	return nil
}
