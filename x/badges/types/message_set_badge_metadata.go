package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetBadgeMetadata = "set_badge_metadata"

var _ sdk.Msg = &MsgSetBadgeMetadata{}

func NewMsgSetBadgeMetadata(creator string, collectionId Uint, badgeMetadataTimeline []*BadgeMetadataTimeline, canUpdateBadgeMetadata []*TimedUpdateWithBadgeIdsPermission) *MsgSetBadgeMetadata {
	return &MsgSetBadgeMetadata{
		Creator:                creator,
		CollectionId:           collectionId,
		BadgeMetadataTimeline:  badgeMetadataTimeline,
		CanUpdateBadgeMetadata: canUpdateBadgeMetadata,
	}
}

func (msg *MsgSetBadgeMetadata) Route() string {
	return RouterKey
}

func (msg *MsgSetBadgeMetadata) Type() string {
	return TypeMsgSetBadgeMetadata
}

func (msg *MsgSetBadgeMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetBadgeMetadata) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetBadgeMetadata) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetBadgeMetadata) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateBadgeMetadataTimeline: true,
		BadgeMetadataTimeline:       msg.BadgeMetadataTimeline,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateBadgeMetadata: msg.CanUpdateBadgeMetadata,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
