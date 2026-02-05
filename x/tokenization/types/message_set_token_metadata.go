package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetTokenMetadata = "set_token_metadata"

var _ sdk.Msg = &MsgSetTokenMetadata{}

func NewMsgSetTokenMetadata(creator string, collectionId Uint, tokenMetadata []*TokenMetadata, canUpdateTokenMetadata []*TokenIdsActionPermission) *MsgSetTokenMetadata {
	return &MsgSetTokenMetadata{
		Creator:                creator,
		CollectionId:           collectionId,
		TokenMetadata:          tokenMetadata,
		CanUpdateTokenMetadata: canUpdateTokenMetadata,
	}
}

func (msg *MsgSetTokenMetadata) Route() string {
	return RouterKey
}

func (msg *MsgSetTokenMetadata) Type() string {
	return TypeMsgSetTokenMetadata
}

func (msg *MsgSetTokenMetadata) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetTokenMetadata) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgSetTokenMetadata) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetTokenMetadata) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := &MsgUniversalUpdateCollection{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateTokenMetadata:         true,
		TokenMetadata:               msg.TokenMetadata,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &CollectionPermissions{
			CanUpdateTokenMetadata: msg.CanUpdateTokenMetadata,
		},
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
