package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetValidTokenIds = "set_valid_badge_ids"

var _ sdk.Msg = &MsgSetValidTokenIds{}

func NewMsgSetValidTokenIds(creator string, collectionId Uint, validTokenIds []*UintRange, canUpdateValidTokenIds []*TokenIdsActionPermission) *MsgSetValidTokenIds {
	return &MsgSetValidTokenIds{
		Creator:                creator,
		CollectionId:           collectionId,
		ValidTokenIds:          validTokenIds,
		CanUpdateValidTokenIds: canUpdateValidTokenIds,
	}
}

func (msg *MsgSetValidTokenIds) Route() string {
	return RouterKey
}

func (msg *MsgSetValidTokenIds) Type() string {
	return TypeMsgSetValidTokenIds
}

func (msg *MsgSetValidTokenIds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetValidTokenIds) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetValidTokenIds) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetValidTokenIds) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateValidTokenIds:         true,
		ValidTokenIds:               msg.ValidTokenIds,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateValidTokenIds: msg.CanUpdateValidTokenIds,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
