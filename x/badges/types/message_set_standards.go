package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetStandards = "set_standards"

var _ sdk.Msg = &MsgSetStandards{}

func NewMsgSetStandards(creator string, collectionId Uint, standards []string, canUpdateStandards []*ActionPermission) *MsgSetStandards {
	return &MsgSetStandards{
		Creator:            creator,
		CollectionId:       collectionId,
		Standards:          standards,
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
		DefaultBalances: &UserBalanceStore{},

		//Applicable to creations and updates
		ValidTokenIds:               []*UintRange{},
		UpdateCollectionPermissions: false,
		CollectionPermissions:       &CollectionPermissions{},
		UpdateManager:               false,
		Manager:                     "",
		UpdateCollectionMetadata:    false,
		CollectionMetadata:          nil,
		UpdateTokenMetadata:         false,
		TokenMetadata:               []*TokenMetadata{},
		UpdateCustomData:            false,
		CustomData:                  "",
		UpdateCollectionApprovals:    false,
		CollectionApprovals:         []*CollectionApproval{},
		UpdateStandards:             false,
		Standards:                   []string{},
		UpdateIsArchived:            false,
		IsArchived:                  false,

		MintEscrowCoinsToTransfer:   []*sdk.Coin{},
		CosmosCoinWrapperPathsToAdd: []*CosmosCoinWrapperPathAddObject{},
		AliasPathsToAdd:             []*AliasPathAddObject{},
		Invariants:                  &InvariantsAddObject{},
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
	ms.UpdateStandards = true
	ms.Standards = msg.Standards
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
