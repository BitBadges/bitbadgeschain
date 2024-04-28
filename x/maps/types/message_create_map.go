package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreateMap = "create_map"

var _ sdk.Msg = &MsgCreateMap{}

func NewMsgCreateMap(creator string, mapId string, updateCriteria *MapUpdateCriteria, valueOptions *ValueOptions, defaultValue string, managerTimeline []*ManagerTimeline, metadataTimeline []*MapMetadataTimeline, permissions *MapPermissions, inheritManagerTimelineFrom sdkmath.Uint) *MsgCreateMap {
	return &MsgCreateMap{
		Creator:                    creator,
		MapId:                      mapId,
		UpdateCriteria:             updateCriteria,
		ValueOptions:               valueOptions,
		DefaultValue:               defaultValue,
		ManagerTimeline:            managerTimeline,
		MetadataTimeline:           metadataTimeline,
		Permissions:                permissions,
		InheritManagerTimelineFrom: inheritManagerTimelineFrom,
	}
}

func (msg *MsgCreateMap) Route() string {
	return RouterKey
}

func (msg *MsgCreateMap) Type() string {
	return TypeMsgCreateMap
}

func (msg *MsgCreateMap) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateMap) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateMap) ValidateBasic() error {
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

	//Validate update criteria
	if msg.UpdateCriteria == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "update criteria cannot be nil")
	}

	numDefined := 0
	if msg.UpdateCriteria.ManagerOnly {
		numDefined++
	}
	if !msg.UpdateCriteria.CollectionId.IsNil() && !msg.UpdateCriteria.CollectionId.IsZero() {
		numDefined++
	}
	if msg.UpdateCriteria.CreatorOnly {
		numDefined++
	}
	if msg.UpdateCriteria.FirstComeFirstServe {
		numDefined++
	}

	if numDefined != 1 {
		return sdkerrors.Wrap(ErrInvalidRequest, "update criteria must have exactly one field defined")
	}

	if ValidatePermissions(msg.Permissions, false) != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "permissions are invalid")
	}

	return nil
}
