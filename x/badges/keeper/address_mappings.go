package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressMapping(ctx sdk.Context, addressMapping *types.AddressMapping) error {
	id := addressMapping.MappingId

	//Validate ID
	if id == "Mint" || id == "Manager" || id == "AllWithoutMint" || id == "None" || id == "AllWithMint" {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot be %s", id)
	}

	//if starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot start with !")
	}

	//if any char is a :
	for _, char := range id {
		if char == ':'  || char == '_' {
			return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot contain : or _")
		}
	}

	if types.ValidateAddress(addressMapping.MappingId, false) == nil {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot be a valid cosmos address")
	}

	err := k.SetAddressMappingInStore(ctx, *addressMapping)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAddressMappingById(ctx sdk.Context, addressMappingId string, managerAddress string) (*types.AddressMapping, error) {
	inverted := false
	handled := false
	addressMapping := &types.AddressMapping{}
	if addressMappingId[0] == '!' {
		inverted = true
		addressMappingId = addressMappingId[1:]
	}

	if addressMappingId == "Mint" {
		addressMapping = &types.AddressMapping{
			MappingId:        "Mint",
			Addresses:        []string{"Mint"},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if addressMappingId == "Manager" {
		addressMapping = &types.AddressMapping{
			MappingId:        "Manager",
			Addresses:        []string{managerAddress},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if addressMappingId == "AllWithoutMint" {
		addressMapping = &types.AddressMapping{
			MappingId:        "AllWithoutMint",
			Addresses:        []string{"Mint"},
			IncludeAddresses: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if addressMappingId == "AllWithMint" {
		addressMapping = &types.AddressMapping{
			MappingId:        "AllWithMint",
			Addresses:        []string{},
			IncludeAddresses: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if addressMappingId == "None" {
		addressMapping = &types.AddressMapping{
			MappingId:        "None",
			Addresses:        []string{},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if types.ValidateAddress(addressMappingId, false) == nil {
		addressMapping = &types.AddressMapping{
			MappingId:        addressMappingId,
			Addresses:        []string{addressMappingId},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	if !handled {
		addressMappingFetched, found := k.GetAddressMappingFromStore(ctx, addressMappingId)
		if found {
			addressMapping = &addressMappingFetched
		} else {
			return nil, sdkerrors.Wrapf(ErrAddressMappingNotFound, "address mapping with id %s not found", addressMappingId)
		}
	}

	if inverted {
		addressMapping.IncludeAddresses = !addressMapping.IncludeAddresses
	}

	return addressMapping, nil
}

func (k Keeper) CheckMappingAddresses(ctx sdk.Context, addressMappingId string, addressToCheck string, managerAddress string) (bool, error) {
	addressMapping, err := k.GetAddressMappingById(ctx, addressMappingId, managerAddress)
	if err != nil {
		return false, err
	}

	found := false
	for _, address := range addressMapping.Addresses {
		if address == addressToCheck {
			found = true
		}
	}

	if !addressMapping.IncludeAddresses {
		found = !found
	}

	if !found {
		return false, nil
	}

	return true, nil
}

// Checks if the addresses are in their respective mapping.
// If includeAddresses is true, then we check if the address is in the Addresses field.
// If includeAddresses is false, then we check if the address is NOT in the Addresses field.

// Note addresses matching does not mean the transfer is allowed. It just means the addresses match.
// All other criteria must also be met.
func (k Keeper) CheckIfAddressesMatchCollectionMappingIds(ctx sdk.Context, collectionApproval *types.CollectionApproval, from string, to string, initiatedBy string, managerAddress string) bool {
	fromFound, err := k.CheckMappingAddresses(ctx, collectionApproval.FromMappingId, from, managerAddress)
	if err != nil {
		return false
	}

	toFound, err := k.CheckMappingAddresses(ctx, collectionApproval.ToMappingId, to, managerAddress)
	if err != nil {
		return false
	}

	initiatedByFound, err := k.CheckMappingAddresses(ctx, collectionApproval.InitiatedByMappingId, initiatedBy, managerAddress)
	if err != nil {
		return false
	}

	return fromFound && toFound && initiatedByFound
}
