package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetCustomData = "set_custom_data"

var _ sdk.Msg = &MsgSetCustomData{}

func NewMsgSetCustomData(creator string, collectionId Uint, customDataTimeline []*CustomDataTimeline, canUpdateCustomData []*TimedUpdatePermission) *MsgSetCustomData {
	return &MsgSetCustomData{
		Creator:             creator,
		CollectionId:        collectionId,
		CustomDataTimeline:  customDataTimeline,
		CanUpdateCustomData: canUpdateCustomData,
	}
}

func (msg *MsgSetCustomData) Route() string {
	return RouterKey
}

func (msg *MsgSetCustomData) Type() string {
	return TypeMsgSetCustomData
}

func (msg *MsgSetCustomData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCustomData) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetCustomData) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetCustomData) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateCustomDataTimeline:    true,
		CustomDataTimeline:          msg.CustomDataTimeline,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateCustomData: msg.CanUpdateCustomData,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
