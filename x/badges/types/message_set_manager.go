package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetManager = "set_manager"

var _ sdk.Msg = &MsgSetManager{}

func NewMsgSetManager(creator string, collectionId Uint, managerTimeline []*ManagerTimeline, canUpdateManager []*TimedUpdatePermission) *MsgSetManager {
	return &MsgSetManager{
		Creator:          creator,
		CollectionId:     collectionId,
		ManagerTimeline:  managerTimeline,
		CanUpdateManager: canUpdateManager,
	}
}

func (msg *MsgSetManager) Route() string {
	return RouterKey
}

func (msg *MsgSetManager) Type() string {
	return TypeMsgSetManager
}

func (msg *MsgSetManager) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetManager) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetManager) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetManager) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateManagerTimeline:       true,
		ManagerTimeline:             msg.ManagerTimeline,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateManager: msg.CanUpdateManager,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
