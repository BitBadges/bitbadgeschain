package keeper

import (
	"encoding/binary"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressMapping(ctx sdk.Context, addressMapping *types.AddressMapping) error {
	id := addressMapping.MappingId

	//if starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot start with !")
	}

	//if any char is a :
	for _, char := range id {
		if char == ':' || char == '_' {
			return sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping id cannot contain : or _")
		}
	}

	_, err := k.GetAddressMappingById(ctx, id)
	if err == nil {
		return sdkerrors.Wrapf(ErrAddressMappingAlreadyExists, "address mapping with id %s already exists or is reserved", id)
	}


	// From cosmos SDK x/group module
	// Generate account address for mapping
	var accountAddr sdk.AccAddress
	// loop here in the rare case where a ADR-028-derived address creates a
	// collision with an existing address.
	for {
		derivationKey := make([]byte, 8)
		nextId := k.GetNextAddressMappingCounter(ctx)
		binary.BigEndian.PutUint64(derivationKey, nextId.Uint64())

		ac, err := authtypes.NewModuleCredential(types.ModuleName, AddressGenerationPrefix, derivationKey)
		if err != nil {
			return err
		}
		//generate the address from the credential
		accountAddr = sdk.AccAddress(ac.Address())
		
		break
	}

	addressMapping.AliasAddress = accountAddr.String()

	err = k.SetAddressMappingInStore(ctx, *addressMapping)
	if err != nil {
		return err
	}

	k.IncrementNextAddressMappingCounter(ctx)

	return nil
}

func getReservedMappingById(addressMappingId string, allowAliases bool) (*types.AddressMapping, bool, error) {
	handled := false
	addressMapping := &types.AddressMapping{}

	if addressMappingId == "Mint" {
		addressMapping = &types.AddressMapping{
			MappingId:        "Mint",
			Addresses:        []string{"Mint"},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if len(addressMappingId) > 10 && addressMappingId[0:10] == "AllWithout" {
		//If starts with AllWithout, we create a mapping with all addresses except the ones specified delimited by :
		addresses := addressMappingId[10:]
		addressMapping = &types.AddressMapping{
			MappingId:        addressMappingId,
			Addresses:        []string{},
			IncludeAddresses: false,
			Uri:              "",
			CustomData:       "",
		}

		//split by :
		splitAdresses := strings.Split(addresses, ":")
		for _, address := range splitAdresses {
			addressMapping.Addresses = append(addressMapping.Addresses, address)

			if err := types.ValidateAddress(address, true); err != nil {
				return nil, false, sdkerrors.Wrapf(ErrInvalidAddressMappingId, "address mapping cannot contain invalid addresses")
			}
		}

		handled = true
	} else if addressMappingId == "All" {
		addressMapping = &types.AddressMapping{
			MappingId:        "All",
			Addresses:        []string{},
			IncludeAddresses: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if addressMappingId == "AllWithMint" {
		addressMapping = &types.AddressMapping{
			MappingId:        "AllWithMint",
			Addresses:        []string{},
			IncludeAddresses: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if addressMappingId == "None" {
		addressMapping = &types.AddressMapping{
			MappingId:        "None",
			Addresses:        []string{},
			IncludeAddresses: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	//Split by :
	if !handled {
		addresses := strings.Split(addressMappingId, ":")
		allAreValid := true
		if !allowAliases {
			for _, address := range addresses {
				if err := types.ValidateAddress(address, true); err != nil {
					allAreValid = false
				}
			}
		}

		if allAreValid {
			addressMapping = &types.AddressMapping{
				MappingId:        addressMappingId,
				Addresses:        addresses,
				IncludeAddresses: true,
				Uri:              "",
				CustomData:       "",
			}
			handled = true
		}
	}

	return addressMapping, handled, nil
}

func (k Keeper) GetTrackerMappingById(ctx sdk.Context, trackerMappingId string) (*types.AddressMapping, error) {
	inverted := false
	addressMapping := &types.AddressMapping{}
	if trackerMappingId[0] == '!' {
		inverted = true
		trackerMappingId = trackerMappingId[1:]
	}

	//Tracker mappings do not allow aliases and are only reserved IDs
	addressMapping, handled, err := getReservedMappingById(trackerMappingId, true)
	if err != nil {
		return nil, err
	}

	if !handled {
		return nil, sdkerrors.Wrapf(ErrAddressMappingNotFound, "tracker mapping with id %s not a reserved ID", trackerMappingId)
	}

	if inverted {
		addressMapping.IncludeAddresses = !addressMapping.IncludeAddresses
	}

	return addressMapping, nil
}

func (k Keeper) GetAddressMappingById(ctx sdk.Context, addressMappingId string) (*types.AddressMapping, error) {
	inverted := false
	addressMapping := &types.AddressMapping{}
	if addressMappingId[0] == '!' {
		inverted = true
		addressMappingId = addressMappingId[1:]
	}

	addressMapping, handled, err := getReservedMappingById(addressMappingId, false)
	if err != nil {
		return nil, err
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

func (k Keeper) CheckMappingAddresses(ctx sdk.Context, addressMappingId string, addressToCheck string) (bool, error) {
	addressMapping, err := k.GetAddressMappingById(ctx, addressMappingId)
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
func (k Keeper) CheckIfAddressesMatchCollectionMappingIds(ctx sdk.Context, collectionApproval *types.CollectionApproval, from string, to string, initiatedBy string) bool {
	fromFound, err := k.CheckMappingAddresses(ctx, collectionApproval.FromMappingId, from)
	if err != nil {
		return false
	}

	toFound, err := k.CheckMappingAddresses(ctx, collectionApproval.ToMappingId, to)
	if err != nil {
		return false
	}

	initiatedByFound, err := k.CheckMappingAddresses(ctx, collectionApproval.InitiatedByMappingId, initiatedBy)
	if err != nil {
		return false
	}

	return fromFound && toFound && initiatedByFound
}
