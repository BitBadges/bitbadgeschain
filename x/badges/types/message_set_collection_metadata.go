package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetCollectionMetadata = "set_collection_metadata"

var _ sdk.Msg = &MsgSetCollectionMetadata{}

func NewMsgSetCollectionMetadata(creator string, collectionId Uint, collectionMetadataTimeline []*CollectionMetadataTimeline, canUpdateCollectionMetadata []*TimedUpdatePermission) *MsgSetCollectionMetadata {
	return &MsgSetCollectionMetadata{
		Creator:                     creator,
		CollectionId:                collectionId,
		CollectionMetadataTimeline:  collectionMetadataTimeline,
		CanUpdateCollectionMetadata: canUpdateCollectionMetadata,
	}
}

func (msg *MsgSetCollectionMetadata) Route() string {
	return RouterKey
}

func (msg *MsgSetCollectionMetadata) Type() string {
	return TypeMsgSetCollectionMetadata
}

func (msg *MsgSetCollectionMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCollectionMetadata) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetCollectionMetadata) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetCollectionMetadata) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                          msg.Creator,
		CollectionId:                     msg.CollectionId,
		UpdateCollectionMetadataTimeline: true,
		CollectionMetadataTimeline:       msg.CollectionMetadataTimeline,
		UpdateCollectionPermissions:      true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateCollectionMetadata: msg.CanUpdateCollectionMetadata,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
