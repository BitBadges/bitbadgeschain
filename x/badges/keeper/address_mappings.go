package keeper

import (
	"math"
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

func (k Keeper) CreateAddressMapping(ctx sdk.Context, addressMapping *types.AddressMapping) error {
	id := addressMapping.MappingId

	//Validate ID 
	if id == "Mint" || id == "Manager" || id == "All" || id == "None"	{
		return ErrInvalidAddressMappingId
	}
	
	//if starts with !
	if id[0] == '!' {
		return ErrInvalidAddressMappingId
	}

	//if any char is a :
	for _, char := range id {
		if char == ':' {
			return ErrInvalidAddressMappingId
		}
	}

	if types.ValidateAddress(addressMapping.MappingId, false) == nil {
		return ErrInvalidAddressMappingId
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
			Filters: nil,
		}, nil
	}

	if addressMappingId == "Manager" {
		return &types.AddressMapping{
			MappingId: "Manager",
			Addresses: []string{"Manager"},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
			Filters: nil,
		}, nil
	}

	if addressMappingId == "All" {
		return &types.AddressMapping{
			MappingId: "All",
			Addresses: []string{},
			IncludeOnlySpecified: false,
			Uri: "",
			CustomData: "",
			Filters: nil,
		}, nil
	}

	if addressMappingId == "None" {
		return &types.AddressMapping{
			MappingId: "None",
			Addresses: []string{},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
			Filters: nil,
		}, nil
	}

	if types.ValidateAddress(addressMappingId, false) == nil {
		return &types.AddressMapping{
			MappingId: addressMappingId,
			Addresses: []string{addressMappingId},
			IncludeOnlySpecified: true,
			Uri: "",
			CustomData: "",
			Filters: nil,
		}, nil
	}


	addressMappingIdSplit := strings.Split(addressMappingId, ":")
	if len(addressMappingIdSplit) == 2 {
		parsedCollectionId, err := strconv.ParseUint(addressMappingIdSplit[0], 10, 64)
		if err != nil {
			return nil, ErrInvalidAddressMappingId
		}

		collectionId := sdkmath.NewUint(parsedCollectionId)

		parsedBadgeId, err := strconv.ParseUint(addressMappingIdSplit[1], 10, 64)
		if err != nil {
			return nil, ErrInvalidAddressMappingId
		}

		badgeId := sdkmath.NewUint(parsedBadgeId)

		return &types.AddressMapping{
			MappingId: addressMappingId,
			Addresses: []string{},
			IncludeOnlySpecified: false,
			Uri: "",
			CustomData: "",
			Filters: []*types.AddressMappingFilters{
				{
					MustSatisfyMinX: sdkmath.NewUint(1),
					Conditions: []*types.AddressMappingConditions{
						{
							MustOwnBadges: []*types.MinMaxBalance{
								{
									CollectionId: collectionId,
									BadgeIds: []*types.IdRange{
										{
											Start: badgeId,
											End: badgeId,
										},
									},
									Amount: &types.IdRange{
										Start: sdkmath.NewUint(1),
										End: sdkmath.NewUint(math.MaxUint64),
									},
								},
							},
						},
					},
				},
			},
		}, nil
	}

	addressMapping, found := k.GetAddressMappingFromStore(ctx, addressMappingId)
	if found {
		return &addressMapping, nil
	}

	return nil, ErrAddressMappingNotFound
}


func (k Keeper) CheckIfManager(ctx sdk.Context, collectionId sdkmath.Uint, address string) (bool, error) {
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return false, ErrCollectionNotFound
	}

	return types.GetCurrentManager(ctx, collection) == address, nil
}

func (k Keeper) CheckIfOwnsBadges(ctx sdk.Context, collectionId sdkmath.Uint, address string, badgeIds []*types.IdRange, minMaxAmount *types.IdRange) (bool, error) {
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return false, ErrCollectionNotFound
	}

	unhandledBadgeIds := make([]*types.IdRange, len(badgeIds))
	copy(unhandledBadgeIds, badgeIds)


	balanceKey := ConstructBalanceKey(address, collection.CollectionId)
	userBalance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		userBalance = &types.UserBalanceStore{
			Balances : []*types.Balance{},
		}
	}

	balances, err := types.GetBalancesForIds(badgeIds, []*types.IdRange{{
			Start: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
			End: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
	}}, userBalance.Balances)
	if err != nil {
		return false, err
	}

	for _, balance := range balances {
		if balance.Amount.LT(minMaxAmount.Start) || balance.Amount.GT(minMaxAmount.End) {
			return false, nil
		}

		unhandledBadgeIds, _ = types.RemoveIdRangeFromIdRange(balance.BadgeIds, unhandledBadgeIds)
	}

	//Sanity check
	if len(unhandledBadgeIds) != 0 {
		return false, ErrUnhandledBadgeIds
	}
	
	return true, nil
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

	//Must check filters
	for _, filter := range addressMapping.Filters {
		mustSatisfyMinX := filter.MustSatisfyMinX
		conditions := filter.Conditions
		satisfyCount := sdkmath.NewUint(0)

		//TODO: support inherited balances with no circular references
		for _, condition := range conditions {
			satisfied := true
			if condition.MustOwnBadges != nil && len(condition.MustOwnBadges) > 0 {
				for _, minMaxBalance := range condition.MustOwnBadges {
					ownsBadges, err := k.CheckIfOwnsBadges(ctx, minMaxBalance.CollectionId, addressToCheck, minMaxBalance.BadgeIds, minMaxBalance.Amount)
					if err != nil {
						return false, err
					}

					if !ownsBadges {
						satisfied = false
						continue
					}
				}
			}

			if condition.MustNotOwnBadges != nil && len(condition.MustNotOwnBadges) > 0 {
				for _, minMaxBalance := range condition.MustNotOwnBadges {
					ownsBadges, err := k.CheckIfOwnsBadges(ctx, minMaxBalance.CollectionId, addressToCheck, minMaxBalance.BadgeIds, minMaxBalance.Amount)
					if err != nil {
						return false, err
					}

					if ownsBadges {
						satisfied = false
						continue
					}
				}
			}

			if condition.MustBeManager != nil && len(condition.MustBeManager) > 0 {
				for _, collectionId := range condition.MustBeManager {
					isManager, err := k.CheckIfManager(ctx, collectionId, addressToCheck)
					if err != nil {
						return false, err
					}

					if !isManager {
						satisfied = false
					}
				}
			}

			if condition.MustNotBeManager != nil && len(condition.MustNotBeManager) > 0 {
				for _, collectionId := range condition.MustNotBeManager {
					isManager, err := k.CheckIfManager(ctx, collectionId, addressToCheck)
					if err != nil {
						return false, err
					}

					if isManager {
						satisfied = false
					}
				}
			}

			if condition.MustBeInMapping != nil && len(condition.MustBeInMapping) > 0 {
				for _, mapping := range condition.MustBeInMapping {
					matches, err := k.CheckMappingAddresses(ctx, mapping, addressToCheck, managerAddress, newCheckedMappingIds)
					if err != nil {
						return false, err
					}

					if !matches {
						satisfied = false
					}
				}
			}

			if condition.MustNotBeInMapping != nil && len(condition.MustNotBeInMapping) > 0 {
				for _, mapping := range condition.MustNotBeInMapping {
					matches, err := k.CheckMappingAddresses(ctx, mapping, addressToCheck, managerAddress, newCheckedMappingIds)
					if err != nil {
						return false, err
					}

					if matches {
						satisfied = false
					}
				}
			}

			if satisfied {
				satisfyCount = satisfyCount.Add(sdkmath.NewUint(1))
			}
		}
		

		if satisfyCount.LT(mustSatisfyMinX) {
			return false, nil
		}
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