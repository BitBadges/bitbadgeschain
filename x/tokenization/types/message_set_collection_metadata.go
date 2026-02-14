package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetCollectionMetadata = "set_collection_metadata"

var _ sdk.Msg = &MsgSetCollectionMetadata{}

func NewMsgSetCollectionMetadata(creator string, collectionId Uint, collectionMetadata *CollectionMetadata, canUpdateCollectionMetadata []*ActionPermission) *MsgSetCollectionMetadata {
	return &MsgSetCollectionMetadata{
		Creator:                     creator,
		CollectionId:                collectionId,
		CollectionMetadata:          collectionMetadata,
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
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
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
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateCollectionMetadata:    true,
		CollectionMetadata:          msg.CollectionMetadata,
		UpdateCollectionPermissions: true,
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
