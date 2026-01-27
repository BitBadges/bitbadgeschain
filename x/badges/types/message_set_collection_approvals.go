package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetCollectionApprovals = "set_collection_approvals"

var _ sdk.Msg = &MsgSetCollectionApprovals{}

func NewMsgSetCollectionApprovals(creator string, collectionId Uint, collectionApprovals []*CollectionApproval, canUpdateCollectionApprovals []*CollectionApprovalPermission) *MsgSetCollectionApprovals {
	return &MsgSetCollectionApprovals{
		Creator:                      creator,
		CollectionId:                 collectionId,
		CollectionApprovals:          collectionApprovals,
		CanUpdateCollectionApprovals: canUpdateCollectionApprovals,
	}
}

func (msg *MsgSetCollectionApprovals) Route() string {
	return RouterKey
}

func (msg *MsgSetCollectionApprovals) Type() string {
	return TypeMsgSetCollectionApprovals
}

func (msg *MsgSetCollectionApprovals) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCollectionApprovals) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetCollectionApprovals) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetCollectionApprovals) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateCollectionApprovals:   true,
		CollectionApprovals:         msg.CollectionApprovals,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateCollectionApprovals: msg.CanUpdateCollectionApprovals,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
