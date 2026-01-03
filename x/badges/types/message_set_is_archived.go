package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetIsArchived = "set_is_archived"

var _ sdk.Msg = &MsgSetIsArchived{}

func NewMsgSetIsArchived(creator string, collectionId Uint, isArchived bool, canArchiveCollection []*ActionPermission) *MsgSetIsArchived {
	return &MsgSetIsArchived{
		Creator:              creator,
		CollectionId:         collectionId,
		IsArchived:           isArchived,
		CanArchiveCollection: canArchiveCollection,
	}
}

func (msg *MsgSetIsArchived) Route() string {
	return RouterKey
}

func (msg *MsgSetIsArchived) Type() string {
	return TypeMsgSetIsArchived
}

func (msg *MsgSetIsArchived) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetIsArchived) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetIsArchived) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetIsArchived) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateIsArchived:            true,
		IsArchived:                  msg.IsArchived,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanArchiveCollection: msg.CanArchiveCollection,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
