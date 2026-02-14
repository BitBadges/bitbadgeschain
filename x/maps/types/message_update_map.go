package types

import (
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateMap = "update_map"

var _ sdk.Msg = &MsgUpdateMap{}

func NewMsgUpdateMap(creator string, mapId string, updateManager bool, manager string, updateMetadata bool, metadata *Metadata, updatePermissions bool, permissions *MapPermissions) *MsgUpdateMap {
	return &MsgUpdateMap{
		Creator:           creator,
		MapId:             mapId,
		UpdateManager:     updateManager,
		Manager:           manager,
		UpdateMetadata:    updateMetadata,
		Metadata:          metadata,
		UpdatePermissions: updatePermissions,
		Permissions:       permissions,
	}
}

func (msg *MsgUpdateMap) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMap) Type() string {
	return TypeMsgUpdateMap
}

func (msg *MsgUpdateMap) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMap) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.MapId) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "map ID cannot be empty")
	}

	err = tokenizationtypes.ValidateManager(msg.Manager)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "manager timeline cannot be invalid")
	}

	err = tokenizationtypes.ValidateCollectionMetadata(CastMapMetadataToCollectionMetadata(msg.Metadata))
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "metadata timeline cannot be invalid")
	}

	if ValidatePermissions(msg.Permissions, false) != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "permissions are invalid")
	}

	return nil
}
