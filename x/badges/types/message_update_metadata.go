package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateMetadata = "update_uris"

var _ sdk.Msg = &MsgUpdateMetadata{}

func NewMsgUpdateMetadata(creator string, collectionId sdk.Uint, collectionMetadata *CollectionMetadata, badgeMetadata []*BadgeMetadata, offChainBalancesMetadata *OffChainBalancesMetadata, customData string, contractAddress string) *MsgUpdateMetadata {
	for _, badgeMetadata := range badgeMetadata {
		badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)
	}

	return &MsgUpdateMetadata{
		Creator:            creator,
		CollectionId:       collectionId,
		CollectionMetadata: collectionMetadata,
		BadgeMetadata:      badgeMetadata,
		OffChainBalancesMetadata:   offChainBalancesMetadata,
		CustomData:         customData,
		ContractAddress:    contractAddress,
	}
}

func (msg *MsgUpdateMetadata) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMetadata) Type() string {
	return TypeMsgUpdateMetadata
}

func (msg *MsgUpdateMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMetadata) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.OffChainBalancesMetadata != nil {
		err = ValidateURI(msg.OffChainBalancesMetadata.Uri)
		if err != nil {
			return err
		}
	}

	if msg.BadgeMetadata != nil && len(msg.BadgeMetadata) > 0 {
		err = ValidateBadgeMetadata(msg.BadgeMetadata)
		if err != nil {
			return err
		}
	}

	if msg.CollectionMetadata != nil {
		err = ValidateURI(msg.CollectionMetadata.Uri)
		if err != nil {
			return err
		}
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "collectionId is 0 or nil")
	}

	return nil
}
