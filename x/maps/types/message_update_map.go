package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

const TypeMsgUpdateMap = "update_map"

var _ sdk.Msg = &MsgUpdateMap{}

func NewMsgUpdateMap(creator string, mapId string, updateManagerTimeline bool, managerTimeline []*ManagerTimeline, updateMetadataTimeline bool, metadataTimeline []*MapMetadataTimeline, updatePermissions bool, permissions *MapPermissions) *MsgUpdateMap {
	return &MsgUpdateMap{
		Creator:                creator,
		MapId:                  mapId,
		UpdateManagerTimeline:  updateManagerTimeline,
		ManagerTimeline:        managerTimeline,
		UpdateMetadataTimeline: updateMetadataTimeline,
		MetadataTimeline:       metadataTimeline,
		UpdatePermissions:      updatePermissions,
		Permissions:            permissions,
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

	err = badgestypes.ValidateManagerTimeline(CastManagerTimelineArray(msg.ManagerTimeline))
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "manager timeline cannot be invalid")
	}

	err = badgestypes.ValidateCollectionMetadataTimeline(CastMetadataTimelineArray(msg.MetadataTimeline))
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "metadata timeline cannot be invalid")
	}

	if ValidatePermissions(msg.Permissions, false) != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "permissions are invalid")
	}

	return nil
}
