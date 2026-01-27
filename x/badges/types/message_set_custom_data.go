package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetCustomData = "set_custom_data"

var _ sdk.Msg = &MsgSetCustomData{}

func NewMsgSetCustomData(creator string, collectionId Uint, customData string, canUpdateCustomData []*ActionPermission) *MsgSetCustomData {
	return &MsgSetCustomData{
		Creator:             creator,
		CollectionId:        collectionId,
		CustomData:          customData,
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
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
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
		UpdateCustomData:            true,
		CustomData:                  msg.CustomData,
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
