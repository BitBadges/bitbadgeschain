package keeper

import (
	"strings"

	"bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressList(ctx sdk.Context, addressList *types.AddressList) error {
	id := addressList.ListId

	//if starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id cannot start with !")
	}

	//Check if all characters are alphanumeric
	for _, char := range id {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id can only contain alphanumeric characters")
		}
	}

	_, err := k.GetAddressListById(ctx, id)
	if err == nil {
		return sdkerrors.Wrapf(ErrAddressListAlreadyExists, "address list with id %s already exists or is reserved", id)
	}

	err = k.SetAddressListInStore(ctx, *addressList)
	if err != nil {
		return err
	}

	k.IncrementNextAddressListCounter(ctx)

	return nil
}

func getReservedListById(addressListId string, allowAliases bool) (*types.AddressList, bool, error) {

	// Handle special reserved IDs
	switch {
	case addressListId == "Mint":
		return &types.AddressList{
			ListId:     "Mint",
			Addresses:  []string{"Mint"},
			Whitelist:  true,
			Uri:        "",
			CustomData: "",
		}, true, nil

	case strings.HasPrefix(addressListId, "AllWithout"):
		addresses := strings.Split(addressListId[10:], ":")
		for _, address := range addresses {
			if err := types.ValidateAddress(address, true); err != nil {
				return nil, false, sdkerrors.Wrapf(ErrInvalidAddressListId, "address list cannot contain invalid addresses")
			}
		}
		return &types.AddressList{
			ListId:     addressListId,
			Addresses:  addresses,
			Whitelist:  false,
			Uri:        "",
			CustomData: "",
		}, true, nil

	case addressListId == "All", addressListId == "AllWithMint":
		return &types.AddressList{
			ListId:     addressListId,
			Addresses:  []string{},
			Whitelist:  false,
			Uri:        "",
			CustomData: "",
		}, true, nil

	case addressListId == "None":
		return &types.AddressList{
			ListId:     addressListId,
			Addresses:  []string{},
			Whitelist:  true,
			Uri:        "",
			CustomData: "",
		}, true, nil
	}

	// Handle colon-separated addresses
	addresses := strings.Split(addressListId, ":")
	if !allowAliases {
		for _, address := range addresses {
			if err := types.ValidateAddress(address, true); err != nil {
				return nil, false, nil
			}
		}
	}

	return &types.AddressList{
		ListId:     addressListId,
		Addresses:  addresses,
		Whitelist:  true,
		Uri:        "",
		CustomData: "",
	}, true, nil
}

func (k Keeper) GetTrackerListById(ctx sdk.Context, trackerListId string) (*types.AddressList, error) {
	inverted := false
	originalId := trackerListId

	if trackerListId[0] == '!' {
		inverted = true
		trackerListId = trackerListId[1:]
	}

	addressList, handled, err := getReservedListById(trackerListId, true)
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "tracker list with id %s not a reserved ID", trackerListId)
	}

	if inverted {
		addressList.Whitelist = !addressList.Whitelist
	}
	addressList.ListId = originalId
	return addressList, nil
}

func (k Keeper) GetAddressListById(ctx sdk.Context, addressListId string) (*types.AddressList, error) {
	inverted := false
	originalId := addressListId

	// Handle inversion patterns
	if addressListId[0] == '!' && len(addressListId) > 1 && addressListId[len(addressListId)-1] != ')' {
		inverted = true
		addressListId = addressListId[1:]
	} else if strings.HasPrefix(addressListId, "!(") && strings.HasSuffix(addressListId, ")") {
		inverted = true
		addressListId = addressListId[2 : len(addressListId)-1]
	}

	addressList, handled, err := getReservedListById(addressListId, false)
	if err != nil {
		return nil, err
	}

	if !handled {
		addressListFetched, found := k.GetAddressListFromStore(ctx, addressListId)
		if !found {
			return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "address list with id %s not found", addressListId)
		}
		addressList = &addressListFetched
	}

	if inverted {
		addressList.Whitelist = !addressList.Whitelist
	}
	addressList.ListId = originalId
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
			break
		}
	}

	if !addressList.Whitelist {
		found = !found
	}

	if !found {
		return false, nil
	}

	return true, nil
}

// Checks if the addresses in the (to, from, initiatedBy) are approved
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
