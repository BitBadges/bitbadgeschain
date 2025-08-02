package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetStandards = "set_standards"

var _ sdk.Msg = &MsgSetStandards{}

func NewMsgSetStandards(creator string, collectionId Uint, standardsTimeline []*StandardsTimeline, canUpdateStandards []*TimedUpdatePermission) *MsgSetStandards {
	return &MsgSetStandards{
		Creator:            creator,
		CollectionId:       collectionId,
		StandardsTimeline:  standardsTimeline,
		CanUpdateStandards: canUpdateStandards,
	}
}

func (msg *MsgSetStandards) Route() string {
	return RouterKey
}

func (msg *MsgSetStandards) Type() string {
	return TypeMsgSetStandards
}

func (msg *MsgSetStandards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetStandards) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func BlankUniversalMsg() *MsgUniversalUpdateCollection {
	return &MsgUniversalUpdateCollection{
		Creator:      "",
		CollectionId: sdkmath.NewUint(0), //We use 0 to indicate a new collection

		//Exclusive to collection creations
		BalancesType:    "Standard",
		DefaultBalances: &UserBalanceStore{},

		//Applicable to creations and updates
		ValidBadgeIds:                          []*UintRange{},
		UpdateCollectionPermissions:            false,
		CollectionPermissions:                  &CollectionPermissions{},
		UpdateManagerTimeline:                  false,
		ManagerTimeline:                        []*ManagerTimeline{},
		UpdateCollectionMetadataTimeline:       false,
		CollectionMetadataTimeline:             []*CollectionMetadataTimeline{},
		UpdateBadgeMetadataTimeline:            false,
		BadgeMetadataTimeline:                  []*BadgeMetadataTimeline{},
		UpdateOffChainBalancesMetadataTimeline: false,
		OffChainBalancesMetadataTimeline:       []*OffChainBalancesMetadataTimeline{},
		UpdateCustomDataTimeline:               false,
		CustomDataTimeline:                     []*CustomDataTimeline{},
		UpdateCollectionApprovals:              false,
		CollectionApprovals:                    []*CollectionApproval{},
		UpdateStandardsTimeline:                false,
		StandardsTimeline:                      []*StandardsTimeline{},
		UpdateIsArchivedTimeline:               false,
		IsArchivedTimeline:                     []*IsArchivedTimeline{},

		MintEscrowCoinsToTransfer:   []*sdk.Coin{},
		CosmosCoinWrapperPathsToAdd: []*CosmosCoinWrapperPathAddObject{},
		Invariants:                  &CollectionInvariants{},
	}
}

func (msg *MsgSetStandards) ValidateBasic() error {
	uni, err := msg.ToUniversalUpdateCollection()
	if err != nil {
		return err
	}
	return uni.ValidateBasic()
}

func (msg *MsgSetStandards) ToUniversalUpdateCollection() (*MsgUniversalUpdateCollection, error) {
	ms := BlankUniversalMsg()
	ms.Creator = msg.Creator
	ms.CollectionId = msg.CollectionId
	ms.UpdateStandardsTimeline = true
	ms.StandardsTimeline = msg.StandardsTimeline
	ms.UpdateCollectionPermissions = true
	ms.CollectionPermissions = &CollectionPermissions{
		CanUpdateStandards: msg.CanUpdateStandards,
	}
	err := ms.CheckAndCleanMsg(sdk.Context{}, true)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
