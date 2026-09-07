package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PermissionFields(p *ManagerSplitterPermissions) []*PermissionCriteria {
	if p == nil {
		return nil
	}
	return []*PermissionCriteria{p.CanDeleteCollection, p.CanArchiveCollection, p.CanUpdateStandards, p.CanUpdateCustomData, p.CanUpdateManager, p.CanUpdateCollectionMetadata, p.CanUpdateValidTokenIds, p.CanUpdateTokenMetadata, p.CanUpdateCollectionApprovals, p.CanAddMoreAliasPaths, p.CanAddMoreCosmosCoinWrapperPaths}
}

func ValidateCanonicalAddresses(p *ManagerSplitterPermissions, addresses ...string) error {
	validate := func(value string) error {
		address, err := sdk.AccAddressFromBech32(value)
		if err != nil {
			return sdkerrors.Wrap(ErrInvalidAddress, err.Error())
		}
		if address.String() != value {
			return sdkerrors.Wrap(ErrInvalidAddress, "address must use canonical lowercase spelling")
		}
		return nil
	}
	for _, address := range addresses {
		if err := validate(address); err != nil {
			return err
		}
	}
	for _, field := range PermissionFields(p) {
		if field == nil {
			continue
		}
		for _, address := range field.ApprovedAddresses {
			if err := validate(address); err != nil {
				return err
			}
		}
	}
	return nil
}
