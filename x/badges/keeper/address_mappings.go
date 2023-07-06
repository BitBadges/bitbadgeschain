package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressMapping(ctx sdk.Context, addressMapping *types.AddressMapping) error {
	id := addressMapping.MappingId

	//Validate ID 
	if id == "Mint" || id == "Manager" || id == "All" || id == "None"	{
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot be %s", id)
	}
	
	//if starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot start with !")
	}

	//if any char is a :
	for _, char := range id {
		if char == ':' {
			return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot contain :")
		}
	}

	if types.ValidateAddress(addressMapping.MappingId, false) == nil {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot be a valid address")
	}

	err := k.SetAddressMappingInStore(ctx, *addressMapping)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAddressMapping(ctx sdk.Context, addressMappingId string, managerAddress string) (*types.AddressMapping, error) {
	if addressMappingId[0] == '!' {
		return nil, ErrInvalidAddressMappingId
	}
	
	if addressMappingId == "Mint" {
		return &types.AddressMapping{
			MappingId: "Mint",
			Addresses: []string{"Mint"},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
		}, nil
	}

	if addressMappingId == "Manager" {
		return &types.AddressMapping{
			MappingId: "Manager",
			Addresses: []string{"Manager"},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
		}, nil
	}

	if addressMappingId == "All" {
		return &types.AddressMapping{
			MappingId: "All",
			Addresses: []string{},
			IncludeOnlySpecified: false,
			Uri: "",
			CustomData: "",
		}, nil
	}

	if addressMappingId == "None" {
		return &types.AddressMapping{
			MappingId: "None",
			Addresses: []string{},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
		}, nil
	}

	if types.ValidateAddress(addressMappingId, false) == nil {
		return &types.AddressMapping{
			MappingId: addressMappingId,
			Addresses: []string{addressMappingId},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
		}, nil
	}

	addressMapping, found := k.GetAddressMappingFromStore(ctx, addressMappingId)
	if found {
		return &addressMapping, nil
	}

	return nil, ErrAddressMappingNotFound
}


//Avoid circular dependencies through the checkedMappingIds
func (k Keeper) CheckMappingAddresses(ctx sdk.Context, addressMappingId string, addressToCheck string, managerAddress string, checkedMappingIds []string) (bool, error) {
	addressMapping, err := k.GetAddressMapping(ctx, addressMappingId, managerAddress)
	if err != nil {
		return false, err
	}
	
	newCheckedMappingIds := make([]string, len(checkedMappingIds) + 1)
	copy(newCheckedMappingIds, checkedMappingIds)
	newCheckedMappingIds[len(checkedMappingIds)] = addressMappingId

	for _, checkedMappingId := range checkedMappingIds {
		if checkedMappingId == addressMappingId {
			return false, ErrCircularDependency
		}
	}

	found := false
	for _, address := range addressMapping.Addresses {
		if address == addressToCheck {
			found = true
		}
		

		//Support the manager alias
		if address == "Manager" && (addressToCheck == managerAddress || addressToCheck == "Manager") {
			found = true
		}
	}

	if !addressMapping.IncludeOnlySpecified {
		found = !found
	}

	if found && addressMapping.MappingId == "All" && addressToCheck == "Mint" {
		return false, nil
	}

	if !found {
		return false, nil
	}
	
	return true, nil
}



// Checks if the from and to addresses are in the transfer approvedTransfer.
// Handles the manager options for the from and to addresses.
// If includeOnlySpecified is true, then we check if the address is in the Addresses field.
// If includeOnlySpecified is false, then we check if the address is NOT in the Addresses field.

// Note addresses matching does not mean the transfer is allowed. It just means the addresses match.
// All other criteria must also be met.
func (k Keeper) CheckIfAddressesMatchCollectionMappingIds(ctx sdk.Context, collectionApprovedTransfer *types.CollectionApprovedTransfer, from string, to string, initiatedBy string, managerAddress string) bool {
	fromFound, err := k.CheckMappingAddresses(ctx, collectionApprovedTransfer.FromMappingId, from, managerAddress, []string{})
	if err != nil {
		return false
	}

	toFound, err := k.CheckMappingAddresses(ctx, collectionApprovedTransfer.ToMappingId, to, managerAddress, []string{})
	if err != nil {
		return false
	}
	
	initiatedByFound, err := k.CheckMappingAddresses(ctx, collectionApprovedTransfer.InitiatedByMappingId, initiatedBy, managerAddress, []string{})
	if err != nil {
		return false
	}

	return fromFound && toFound && initiatedByFound
}