package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetValidBadgeIds = "set_valid_badge_ids"

var _ sdk.Msg = &MsgSetValidBadgeIds{}

func NewMsgSetValidBadgeIds(creator string, collectionId Uint, validBadgeIds []*UintRange, canUpdateValidBadgeIds []*BadgeIdsActionPermission) *MsgSetValidBadgeIds {
	return &MsgSetValidBadgeIds{
		Creator:                creator,
		CollectionId:           collectionId,
		ValidBadgeIds:          validBadgeIds,
		CanUpdateValidBadgeIds: canUpdateValidBadgeIds,
	}
}

func (msg *MsgSetValidBadgeIds) Route() string {
	return RouterKey
}

func (msg *MsgSetValidBadgeIds) Type() string {
	return TypeMsgSetValidBadgeIds
}

func (msg *MsgSetValidBadgeIds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetValidBadgeIds) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetValidBadgeIds) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetValidBadgeIds) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateValidBadgeIds:         true,
		ValidBadgeIds:               msg.ValidBadgeIds,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateValidBadgeIds: msg.CanUpdateValidBadgeIds,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
