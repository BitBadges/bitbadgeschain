package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewCollection = "new_collection"

var _ sdk.Msg = &MsgNewCollection{}

func NewMsgNewCollection(creator string, standard sdk.Uint, collectionsToCreate []*BadgeSupplyAndAmount, collectionUri string, badgeUris []*BadgeUri, permissions sdk.Uint, allowedTransfers []*TransferMapping, managerApprovedTransfers []*TransferMapping, bytesToStore string, transfers []*Transfer, claims []*Claim, balancesUri string) *MsgNewCollection {
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	for _, badgeUri := range badgeUris {
		badgeUri.BadgeIds = SortAndMergeOverlapping(badgeUri.BadgeIds)
	}

	for _, claim := range claims {
		for _, balance := range claim.UndistributedBalances {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
		
		for _, balance := range claim.CurrentClaimAmounts {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	
	return &MsgNewCollection{
		Creator:                  creator,
		CollectionUri:            collectionUri,
		BadgeUris:                badgeUris,
		BadgeSupplys:             collectionsToCreate,
		AllowedTransfers:         allowedTransfers,
		ManagerApprovedTransfers: managerApprovedTransfers,
		Bytes:                    bytesToStore,
		Permissions:              permissions,
		Standard:                 standard,
		Transfers:                transfers,
		Claims:                   claims,
		BalancesUri: 							balancesUri,
	}
}

func (msg *MsgNewCollection) Route() string {
	return RouterKey
}

func (msg *MsgNewCollection) Type() string {
	return TypeMsgNewCollection
}

func (msg *MsgNewCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateURI(*&msg.CollectionUri); err != nil {
		return err
	}

	if msg.BalancesUri != "" {
		if err := ValidateURI(*&msg.BalancesUri); err != nil {
			return err
		}
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	if msg.BadgeUris != nil {
		if err := ValidateBadgeUris(msg.BadgeUris); err != nil {
			return err
		}
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	if err := ValidateBytes(msg.Bytes); err != nil {
		return err
	}

	amounts := make([]sdk.Uint, len(msg.BadgeSupplys))
	supplys := make([]sdk.Uint, len(msg.BadgeSupplys))
	for i, subasset := range msg.BadgeSupplys {
		if subasset.Amount.IsNil() || subasset.Supply.IsNil() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid amount (%s)", subasset.Amount)
		}

		amounts[i] = subasset.Amount
		supplys[i] = subasset.Supply
	}

	err = ValidateNoElementIsX(amounts, sdk.NewUint(0))
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(supplys, sdk.NewUint(0))
	if err != nil {
		return err
	}

	for _, claim := range msg.Claims {
		err = ValidateClaim(claim)
		if err != nil {
			return err
		}
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	for _, transfer := range msg.AllowedTransfers {
		err = ValidateTransferMapping(*transfer)
		if err != nil {
			return err
		}
	}

	for _, transfer := range msg.ManagerApprovedTransfers {
		err = ValidateTransferMapping(*transfer)
		if err != nil {
			return err
		}
	}

	if msg.Standard.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid standard (%s)", msg.Standard)
	}

	return nil
}
