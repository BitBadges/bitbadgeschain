package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateUris = "update_uris"

var _ sdk.Msg = &MsgUpdateUris{}

func NewMsgUpdateUris(creator string, collectionId uint64, collectionUri string, badgeUris []*BadgeUri) *MsgUpdateUris {
	return &MsgUpdateUris{
		Creator:       creator,
		CollectionId:  collectionId,
		CollectionUri: collectionUri,
		BadgeUris:     badgeUris,
	}
}

func (msg *MsgUpdateUris) Route() string {
	return RouterKey
}

func (msg *MsgUpdateUris) Type() string {
	return TypeMsgUpdateUris
}

func (msg *MsgUpdateUris) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateUris) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateUris) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.BadgeUris == nil || len(msg.BadgeUris) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "badgeUris cannot be nil")
	}

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

	if err := ValidateURI(msg.CollectionUri); err != nil {
		return err
	}

	// hasId := strings.Contains(msg.BadgeUri, "{id}")
	// if !hasId {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "badgeUri must contain \"{id}\"")
	// }

	return nil
}
