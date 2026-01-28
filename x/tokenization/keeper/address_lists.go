package keeper

import (
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) CreateAddressList(ctx sdk.Context, addressList *types.AddressList) error {
	id := addressList.ListId

	// Check if ID is empty
	if id == "" {
		return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id cannot be empty")
	}

	// Check if ID starts with !
	if id[0] == '!' {
		return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id cannot start with !")
	}

	// Check if all characters are alphanumeric
	for _, char := range id {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
			return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list id can only contain alphanumeric characters")
		}
	}

	// Validate addresses in the list
	for _, address := range addressList.Addresses {
		if address == "" {
			return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list cannot contain empty addresses")
		}

		// Allow Mint address to be included in address lists
		if err := types.ValidateAddress(address, true); err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddressListId, "address list contains invalid address: %s", err)
		}
	}

	// Check for duplicate addresses
	seenAddresses := make(map[string]bool, len(addressList.Addresses))
	for _, address := range addressList.Addresses {
		if seenAddresses[address] {
			return sdkerrors.Wrapf(types.ErrDuplicateAddresses, "address list cannot contain duplicate addresses")
		}
		seenAddresses[address] = true
	}

	_, err := k.GetAddressListById(ctx, id)
	if err == nil {
		return sdkerrors.Wrapf(ErrAddressListAlreadyExists, "address list with id %s already exists or is reserved", id)
	}

	err = k.SetAddressListInStore(ctx, *addressList)
	if err != nil {
		return err
	}

	if err := k.IncrementNextAddressListCounter(ctx); err != nil {
		return err
	}

	return nil
}

// parseInversionPattern extracts inversion patterns from a list ID.
// If allowParentheses is true, supports both "!" and "!(...)" patterns.
// If false, only supports "!" pattern.
// Returns: (isInverted, cleanedId, originalId)
func parseInversionPattern(listId string, allowParentheses bool) (bool, string, string) {
	originalId := listId
	inverted := false

	if allowParentheses {
		// Support both ! and !(...) patterns
		if len(listId) > 0 && listId[0] == '!' && len(listId) > 1 && listId[len(listId)-1] != ')' {
			inverted = true
			listId = listId[1:]
		} else if strings.HasPrefix(listId, "!(") && strings.HasSuffix(listId, ")") {
			inverted = true
			listId = listId[2 : len(listId)-1]
		}
	} else {
		// Only support ! pattern
		if len(listId) > 0 && listId[0] == '!' {
			inverted = true
			listId = listId[1:]
		}
	}

	return inverted, listId, originalId
}

func getReservedListById(addressListId string, allowAliases bool) (*types.AddressList, error) {
	// Handle special reserved IDs
	switch {
	case addressListId == types.MintAddress:
		return &types.AddressList{
			ListId:    types.MintAddress,
			Addresses: []string{types.MintAddress},
			Whitelist: true,
		}, nil

	case strings.HasPrefix(addressListId, "AllWithout"):
		const allWithoutPrefix = "AllWithout"
		addresses := strings.Split(addressListId[len(allWithoutPrefix):], ":")
		for _, address := range addresses {
			if err := types.ValidateAddress(address, true); err != nil {
				return nil, sdkerrors.Wrapf(ErrInvalidAddressListId, "address list cannot contain invalid addresses")
			}
		}
		return &types.AddressList{
			ListId:    addressListId,
			Addresses: addresses,
			Whitelist: false,
		}, nil

	case addressListId == "All", addressListId == "AllWithMint":
		return &types.AddressList{
			ListId:    addressListId,
			Addresses: []string{},
			Whitelist: false,
		}, nil

	case addressListId == "None":
		return &types.AddressList{
			ListId:    addressListId,
			Addresses: []string{},
			Whitelist: true,
		}, nil
	}

	// Handle colon-separated addresses
	addresses := strings.Split(addressListId, ":")
	if !allowAliases {
		for _, address := range addresses {
			if err := types.ValidateAddress(address, true); err != nil {
				return nil, sdkerrors.Wrapf(ErrInvalidAddressListId, "address list cannot contain invalid addresses")
			}
		}
	}

	return &types.AddressList{
		ListId:    addressListId,
		Addresses: addresses,
		Whitelist: true,
	}, nil
}

func (k Keeper) GetTrackerListById(ctx sdk.Context, trackerListId string) (*types.AddressList, error) {
	inverted, cleanedId, originalId := parseInversionPattern(trackerListId, true)

	addressList, err := getReservedListById(cleanedId, true)
	if err != nil {
		return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "tracker list with id %s not a reserved ID", cleanedId)
	}

	if inverted {
		addressList.Whitelist = !addressList.Whitelist
	}
	addressList.ListId = originalId
	return addressList, nil
}

func (k Keeper) GetAddressListById(ctx sdk.Context, addressListId string) (*types.AddressList, error) {
	inverted, cleanedId, originalId := parseInversionPattern(addressListId, true)

	addressList, err := getReservedListById(cleanedId, false)
	if err != nil {
		addressListFetched, found := k.GetAddressListFromStore(ctx, cleanedId)
		if !found {
			return nil, sdkerrors.Wrapf(ErrAddressListNotFound, "address list with id %s not found", cleanedId)
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
func (k Keeper) CheckIfAddressesMatchCollectionListIds(ctx sdk.Context, approval *types.CollectionApproval, from string, to string, initiatedBy string) bool {
	if approval == nil {
		panic("approval cannot be nil")
	}

	fromFound, err := k.CheckAddresses(ctx, approval.FromListId, from)
	if err != nil {
		return false
	}

	toFound, err := k.CheckAddresses(ctx, approval.ToListId, to)
	if err != nil {
		return false
	}

	initiatedByFound, err := k.CheckAddresses(ctx, approval.InitiatedByListId, initiatedBy)
	if err != nil {
		return false
	}

	return fromFound && toFound && initiatedByFound
}
