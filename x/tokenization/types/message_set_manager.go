package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetManager = "set_manager"

var _ sdk.Msg = &MsgSetManager{}

func NewMsgSetManager(creator string, collectionId Uint, manager string, canUpdateManager []*ActionPermission) *MsgSetManager {
	return &MsgSetManager{
		Creator:          creator,
		CollectionId:     collectionId,
		Manager:          manager,
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
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
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
		UpdateManager:               true,
		Manager:                     msg.Manager,
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
