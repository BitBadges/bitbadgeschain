package keeper

import (
	"encoding/binary"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressList(ctx sdk.Context, addressList *types.AddressList) error {
	id := addressList.ListId

	//if starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id cannot start with !")
	}

	//if any char is a :
	for _, char := range id {
		if char == ':' || char == '_' {
			return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id cannot contain : or _")
		}
	}

	_, err := k.GetAddressListById(ctx, id)
	if err == nil {
		return sdkerrors.Wrapf(ErrAddressListAlreadyExists, "address list with id %s already exists or is reserved", id)
	}


	// From cosmos SDK x/group module
	// Generate account address for list
	var accountAddr sdk.AccAddress
	// loop here in the rare case where a ADR-028-derived address creates a
	// collision with an existing address.
	for {
		derivationKey := make([]byte, 8)
		nextId := k.GetNextAddressListCounter(ctx)
		binary.BigEndian.PutUint64(derivationKey, nextId.Uint64())

		ac, err := authtypes.NewModuleCredential(types.ModuleName, AddressGenerationPrefix, derivationKey)
		if err != nil {
			return err
		}
		//generate the address from the credential
		accountAddr = sdk.AccAddress(ac.Address())
		
		break
	}

	addressList.AliasAddress = accountAddr.String()

	err = k.SetAddressListInStore(ctx, *addressList)
	if err != nil {
		return err
	}

	k.IncrementNextAddressListCounter(ctx)

	return nil
}

func getReservedListById(addressListId string, allowAliases bool) (*types.AddressList, bool, error) {
	handled := false
	addressList := &types.AddressList{}

	if addressListId == "Mint" {
		addressList = &types.AddressList{
			ListId:        "Mint",
			Addresses:        []string{"Mint"},
			Allowlist: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if len(addressListId) > 10 && addressListId[0:10] == "AllWithout" {
		//If starts with AllWithout, we create a list with all addresses except the ones specified delimited by :
		addresses := addressListId[10:]
		addressList = &types.AddressList{
			ListId:        addressListId,
			Addresses:        []string{},
			Allowlist: false,
			Uri:              "",
			CustomData:       "",
		}

		//split by :
		splitAdresses := strings.Split(addresses, ":")
		for _, address := range splitAdresses {
			addressList.Addresses = append(addressList.Addresses, address)

			if err := types.ValidateAddress(address, true); err != nil {
				return nil, false, sdkerrors.Wrapf(ErrInvalidAddressListId, "address list cannot contain invalid addresses")
			}
		}

		handled = true
	} else if addressListId == "All" {
		addressList = &types.AddressList{
			ListId:        "All",
			Addresses:        []string{},
			Allowlist: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if addressListId == "AllWithMint" {
		addressList = &types.AddressList{
			ListId:        "AllWithMint",
			Addresses:        []string{},
			Allowlist: false,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	} else if addressListId == "None" {
		addressList = &types.AddressList{
			ListId:        "None",
			Addresses:        []string{},
			Allowlist: true,
			Uri:              "",
			CustomData:       "",
		}
		handled = true
	}

	//Split by :
	if !handled {
		addresses := strings.Split(addressListId, ":")
		allAreValid := true
		if !allowAliases {
			for _, address := range addresses {
				if err := types.ValidateAddress(address, true); err != nil {
					allAreValid = false
				}
			}
		}

		if allAreValid {
			addressList = &types.AddressList{
				ListId:        addressListId,
				Addresses:        addresses,
				Allowlist: true,
				Uri:              "",
				CustomData:       "",
			}
			handled = true
		}
	}

	return addressList, handled, nil
}

func (k Keeper) GetTrackerListById(ctx sdk.Context, trackerListId string) (*types.AddressList, error) {
	inverted := false
	addressList := &types.AddressList{}
	if trackerListId[0] == '!' {
		inverted = true
		trackerListId = trackerListId[1:]
	}

	//Tracker lists do not allow aliases and are only reserved IDs
	addressList, handled, err := getReservedListById(trackerListId, true)
	if err != nil {
		return nil, err
	}

	if !handled {
		return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "tracker list with id %s not a reserved ID", trackerListId)
	}

	if inverted {
		addressList.Allowlist = !addressList.Allowlist
	}

	return addressList, nil
}

func (k Keeper) GetAddressListById(ctx sdk.Context, addressListId string) (*types.AddressList, error) {
	inverted := false
	addressList := &types.AddressList{}
	if addressListId[0] == '!' && len(addressListId) > 1 && addressListId[len(addressListId)-1] != ')' {
		inverted = true
		addressListId = addressListId[1:]
	} else if  len(addressListId) > 3 && addressListId[0:2] == "!(" && addressListId[len(addressListId)-1] == ')' {
		inverted = true
		addressListId = addressListId[2:len(addressListId)-1]
	}


	addressList, handled, err := getReservedListById(addressListId, false)
	if err != nil {
		return nil, err
	}


	if !handled {
		addressListFetched, found := k.GetAddressListFromStore(ctx, addressListId)
		if found {
			addressList = &addressListFetched
		} else {
			return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "address list with id %s not found", addressListId)
		}
	}

	if inverted {
		addressList.Allowlist = !addressList.Allowlist
	}

	return addressList, nil
}

func (k Keeper) CheckAddresses(ctx sdk.Context, addressListId string, addressToCheck string) (bool, error) {
	addressList, err := k.GetAddressListById(ctx, addressListId)
	if err != nil {
		return false, err
	}

	found := false
	for _, address := range addressList.Addresses {
		if address == addressToCheck {
			found = true
		}
	}

	if !addressList.Allowlist {
		found = !found
	}

	if !found {
		return false, nil
	}

	return true, nil
}

// Checks if the addresses are in their respective list.
// If allowlist is true, then we check if the address is in the Addresses field.
// If allowlist is false, then we check if the address is NOT in the Addresses field.

// Note addresses matching does not mean the transfer is allowed. It just means the addresses match.
// All other criteria must also be met.
func (k Keeper) CheckIfAddressesMatchCollectionListIds(ctx sdk.Context, collectionApproval *types.CollectionApproval, from string, to string, initiatedBy string) bool {
	fromFound, err := k.CheckAddresses(ctx, collectionApproval.FromListId, from)
	if err != nil {
		return false
	}

	toFound, err := k.CheckAddresses(ctx, collectionApproval.ToListId, to)
	if err != nil {
		return false
	}

	initiatedByFound, err := k.CheckAddresses(ctx, collectionApproval.InitiatedByListId, initiatedBy)
	if err != nil {
		return false
	}

	return fromFound && toFound && initiatedByFound
}
