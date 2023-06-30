package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateMetadata = "update_uris"

var _ sdk.Msg = &MsgUpdateMetadata{}

func NewMsgUpdateMetadata(creator string, collectionId sdkmath.Uint, collectionMetadataTimeline []*CollectionMetadataTimeline, badgeMetadataTimeline []*BadgeMetadataTimeline, offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline, customDataTimeline []*CustomDataTimeline, contractAddressTimeline []*ContractAddressTimeline, standardsTimeline []*StandardTimeline) *MsgUpdateMetadata {
	return &MsgUpdateMetadata{
		Creator:            creator,
		CollectionId:       collectionId,
		CollectionMetadataTimeline: collectionMetadataTimeline,
		BadgeMetadataTimeline:      badgeMetadataTimeline,
		OffChainBalancesMetadataTimeline:   offChainBalancesMetadataTimeline,
		CustomDataTimeline:         customDataTimeline,
		ContractAddressTimeline:    contractAddressTimeline,
		StandardsTimeline: 			standardsTimeline,
	}
}

func (msg *MsgUpdateMetadata) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMetadata) Type() string {
	return TypeMsgUpdateMetadata
}

func (msg *MsgUpdateMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMetadata) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateOffChainBalancesMetadataTimeline(msg.OffChainBalancesMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateContractAddressTimeline(msg.ContractAddressTimeline); err != nil {
		return err
	}

	if err := ValidateCustomDataTimeline(msg.CustomDataTimeline); err != nil {
		return err
	}

	if err := ValidateStandardsTimeline(msg.StandardsTimeline); err != nil {
		return err
	}

	if msg.CollectionId.IsNil() || msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "collectionId is 0 or nil")
	}

	return nil
}
