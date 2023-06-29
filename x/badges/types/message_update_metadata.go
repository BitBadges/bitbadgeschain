package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateMetadata = "update_uris"

var _ sdk.Msg = &MsgUpdateMetadata{}

func NewMsgUpdateMetadata(creator string, collectionId sdk.Uint, collectionMetadataTimeline []*CollectionMetadataTimeline, badgeMetadataTimeline []*BadgeMetadataTimeline, offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline, customDataTimeline []*CustomDataTimeline, contractAddressTimeline []*ContractAddressTimeline) *MsgUpdateMetadata {
	// for _, badgeMetadata := range badgeMetadata {
	// 	badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)
	// }

	return &MsgUpdateMetadata{
		Creator:            creator,
		CollectionId:       collectionId,
		CollectionMetadataTimeline: collectionMetadataTimeline,
		BadgeMetadataTimeline:      badgeMetadataTimeline,
		OffChainBalancesMetadataTimeline:   offChainBalancesMetadataTimeline,
		CustomDataTimeline:         customDataTimeline,
		ContractAddressTimeline:    contractAddressTimeline,
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.OffChainBalancesMetadataTimeline != nil {
		for _, timelineVal := range msg.OffChainBalancesMetadataTimeline {	
			err = ValidateURI(timelineVal.OffChainBalancesMetadata.Uri)
			if err != nil {
				return err
			}
		}
	}

	if msg.BadgeMetadataTimeline != nil && len(msg.BadgeMetadataTimeline) > 0 {
		for _, timelineVal := range msg.BadgeMetadataTimeline {
			err = ValidateBadgeMetadata(timelineVal.BadgeMetadata)
			if err != nil {
				return err
			}
		}
	}

	if msg.CollectionMetadataTimeline != nil {
		for _, timelineVal := range msg.CollectionMetadataTimeline {
			err = ValidateURI(timelineVal.CollectionMetadata.Uri)
			if err != nil {
				return err
			}
		}
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "collectionId is 0 or nil")
	}

	return nil
}
