package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewCollection = "new_collection"

var _ sdk.Msg = &MsgNewCollection{}

func NewMsgNewCollection(creator string, standard uint64, collectionsToCreate []*BadgeSupplyAndAmount, collectionUri string, badgeUri string, permissions uint64, disallowedTransfers []*TransferMapping, managerApprovedTransfers []*TransferMapping, bytesToStore string, transfers []*Transfers, claims []*Claim) *MsgNewCollection {
	return &MsgNewCollection{
		Creator:                  creator,
		CollectionUri:            collectionUri,
		BadgeUri:                 badgeUri,
		BadgeSupplys:             collectionsToCreate,
		DisallowedTransfers:      disallowedTransfers,
		ManagerApprovedTransfers: managerApprovedTransfers,
		Bytes:                    bytesToStore,
		Permissions:              permissions,
		Standard:                 standard,
		Transfers:                transfers,
		Claims:                   claims,
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

	if err := ValidateURI(*&msg.BadgeUri); err != nil {
		return err
	}

	hasId := strings.Contains(msg.BadgeUri, "{id}")
	if !hasId {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "badge uri must contain {id}")
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	if err := ValidateBytes(msg.Bytes); err != nil {
		return err
	}

	amounts := make([]uint64, len(msg.BadgeSupplys))
	supplys := make([]uint64, len(msg.BadgeSupplys))
	for i, subasset := range msg.BadgeSupplys {
		amounts[i] = subasset.Amount
		supplys[i] = subasset.Supply
	}

	err = ValidateNoElementIsX(amounts, 0)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(supplys, 0)
	if err != nil {
		return err
	}

	// err = ValidateRangesAreValid(msg.FreezeAddressRanges)
	// if err != nil {
	// 	return err
	// }

	for _, claim := range msg.Claims {
		if claim.Uri != "" {
			err = ValidateURI(claim.Uri)
			if err != nil {
				return err
			}
		}

		
	
		if claim.TimeRange == nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
		}
	
		err = ValidateRangesAreValid([]*IdRange{claim.TimeRange})
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
		}

		if claim.Type == uint64(ClaimType_AccountNum) {
			if claim.Data == "" || len(claim.Data) == 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid claim data")
			}
		}

		if claim.Type == uint64(ClaimType_Code) {
			if claim.BadgeIds == nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
			}
		
			err = ValidateRangesAreValid([]*IdRange{claim.BadgeIds})
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
			}
		}
	}

	return nil
}
