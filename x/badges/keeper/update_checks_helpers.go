package keeper

import (
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	proto "github.com/gogo/protobuf/proto"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file is responsible for verifying that if we go from Value A to Value B for a timeline, that the update is valid
// So here, we check the following:
//-Assert the collection has the correct balances type, if we are updating a balance type-specific field
//-For all updates, check that we are able to update according to the permissions.
// This means we have to verify that the current permissions do not forbid the update.
//
//To do this, we have to do the following:
//-1) Get the combination of values which are "updated". Note this depends on the field, so this is kept generic through a function (GetUpdateCombinationsToCheck)
//		We are also dealing with timelines, so "updated" depends on the respective times, as well as the field value.
//		We do this in a multi step-process:
//		-First, we cast the timeline to UniversalPermission using only the TimelineTimes. We store the timeline VALUE in the ArbitraryValue field.
//		-Second, we get the overlaps and non-overlaps between the (old, new) timeline times.
//		 Note we also have to handle edge cases (in one but not the other). We add empty values where needed.
//		 This then leaves us with a list of all the (timeA-timeB, valueX) - (timeA-timeB, valueY) pairs we need to check.
//		-Third, we compare all valueX and valueY values to see if the actual value was updated.
//		 If the value was not updated, then for timeA-timeB, we do not need to check the permissions.
//		 If it was updated, we need to check the permissions for timeA-timeB.
//		-Lastly, if it was updated, in addition to just simply checking timeA-timeB, we may also have to be more specific with what that we need to check.
//		 Ex: If we go from [tokenIDs 1 to 10 -> www.example.com] to [tokenIDs 1 to 2 -> www.example2.com, tokenIDs 3 to 10 -> www.example.com],
//				 we only need to check tokenIDs 1 to 2 from timeA-timeB
//		 We eventually end with a (timeA-timeB, tokenIds, transferTimes, toList, fromList, initiatedByList) tuples array[] that we need to check, adding dummy values where needed.
//		 This step and the third step is field-specific, so that is why we do it via a generic custom function (GetUpdatedStringCombinations, GetUpdatedBoolCombinations, etc...)
//-2) For all the values that are considered "updated", we check if we are allowed to update them, according to the permissions.
//		This is done by fetching wherever the returned tuples from above overlaps any of the permission's (timelineTime, tokenIds, transferTimes, toList, fromList, initiatedByList) tuples, again adding dummy values where needed.
//		For all overlaps, we then assert that the current block time is NOT forbidden (permitted or undefined both correspond to allowed)
//		If all are not forbidden, it is a valid update.

// To make it easier, we first
func GetPotentialUpdatesForTimelineValues(ctx sdk.Context,
	times [][]*types.UintRange, values []interface{},
) []*types.UniversalPermissionDetails {
	castedPermissions := []*types.UniversalPermission{}
	for idx, time := range times {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TimelineTimes:     time,
			ArbitraryValue:    values[idx],
			UsesTimelineTimes: true,
		})
	}

	// I think this is unnecessary because we already disallow duplicate timeline times in ValidateBasic but I will keep it here for now
	return types.GetFirstMatchOnly(ctx, castedPermissions)
}

// Make a struct with a bool flag isApproved and an approval details arr
type ApprovalCriteriaWithIsApproved struct {
	ApprovalCriteria *types.ApprovalCriteria
	Uri              string
	CustomData       string
}

func GetFirstMatchOnlyWithApprovalCriteria(ctx sdk.Context, permissions []*types.UniversalPermission) ([]*types.UniversalPermissionDetails, error) {
	handled := []*types.UniversalPermissionDetails{}
	for _, permission := range permissions {
		tokenIds := types.GetUintRangesWithOptions(permission.TokenIds, permission.UsesTokenIds)
		timelineTimes := types.GetUintRangesWithOptions(permission.TimelineTimes, permission.UsesTimelineTimes)
		transferTimes := types.GetUintRangesWithOptions(permission.TransferTimes, permission.UsesTransferTimes)
		ownershipTimes := types.GetUintRangesWithOptions(permission.OwnershipTimes, permission.UsesOwnershipTimes)
		permanentlyPermittedTimes := types.GetUintRangesWithOptions(permission.PermanentlyPermittedTimes, true)
		permanentlyForbiddenTimes := types.GetUintRangesWithOptions(permission.PermanentlyForbiddenTimes, true)

		toList := types.GetListWithOptions(permission.ToList, permission.UsesToList)
		fromList := types.GetListWithOptions(permission.FromList, permission.UsesFromList)
		initiatedByList := types.GetListWithOptions(permission.InitiatedByList, permission.UsesInitiatedByList)
		approvalIdList := types.GetListWithOptions(permission.ApprovalIdList, permission.UsesApprovalId)

		for _, tokenId := range tokenIds {
			for _, timelineTime := range timelineTimes {
				for _, transferTime := range transferTimes {
					for _, ownershipTime := range ownershipTimes {
						approval := permission.ArbitraryValue.(*types.CollectionApproval)
						arbValue := []*ApprovalCriteriaWithIsApproved{
							{
								ApprovalCriteria: approval.ApprovalCriteria,
								Uri:              approval.Uri,
								CustomData:       approval.CustomData,
							},
						}

						brokenDown := []*types.UniversalPermissionDetails{
							{
								TokenId:         tokenId,
								TimelineTime:    timelineTime,
								TransferTime:    transferTime,
								OwnershipTime:   ownershipTime,
								ToList:          toList,
								FromList:        fromList,
								InitiatedByList: initiatedByList,
								ApprovalIdList:  approvalIdList,

								ArbitraryValue: arbValue,
							},
						}

						overlaps, inBrokenDownButNotHandled, inHandledButNotBrokenDown := types.GetOverlapsAndNonOverlaps(ctx, brokenDown, handled)
						handled = []*types.UniversalPermissionDetails{}
						// if no overlaps, we can just append all of them
						handled = append(handled, inHandledButNotBrokenDown...)
						handled = append(handled, inBrokenDownButNotHandled...)

						// for overlaps, we append approval details
						for _, overlap := range overlaps {
							secondCriteria, ok := overlap.SecondDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)
							if !ok {
								return nil, sdkerrors.Wrapf(types.ErrInvalidArbitraryValueType, "expected []*ApprovalCriteriaWithIsApproved for SecondDetails, got %T", overlap.SecondDetails.ArbitraryValue)
							}
							firstCriteria, ok := overlap.FirstDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)
							if !ok {
								return nil, sdkerrors.Wrapf(types.ErrInvalidArbitraryValueType, "expected []*ApprovalCriteriaWithIsApproved for FirstDetails, got %T", overlap.FirstDetails.ArbitraryValue)
							}
							mergedApprovalCriteria := secondCriteria
							mergedApprovalCriteria = append(mergedApprovalCriteria, firstCriteria...)

							newArbValue := mergedApprovalCriteria
							handled = append(handled, &types.UniversalPermissionDetails{
								TimelineTime:    overlap.Overlap.TimelineTime,
								TokenId:         overlap.Overlap.TokenId,
								TransferTime:    overlap.Overlap.TransferTime,
								OwnershipTime:   overlap.Overlap.OwnershipTime,
								ToList:          overlap.Overlap.ToList,
								FromList:        overlap.Overlap.FromList,
								InitiatedByList: overlap.Overlap.InitiatedByList,

								ApprovalIdList: overlap.Overlap.ApprovalIdList,

								// Appended for future lookups (not involved in overlap logic)
								PermanentlyPermittedTimes: permanentlyPermittedTimes,
								PermanentlyForbiddenTimes: permanentlyForbiddenTimes,
								ArbitraryValue:            newArbValue,
							})
						}
					}
				}
			}
		}
	}

	// It is first match only, so we can do this
	// To help with determinism in comparing later, we sort by token ID
	returnArr := []*types.UniversalPermissionDetails{}
	for _, handledItem := range handled {
		idxToInsert := 0
		for idxToInsert < len(returnArr) && handledItem.TokenId.Start.GT(returnArr[idxToInsert].TokenId.Start) {
			idxToInsert++
		}

		returnArr = append(returnArr, nil)
		copy(returnArr[idxToInsert+1:], returnArr[idxToInsert:])
		returnArr[idxToInsert] = handledItem
	}

	return handled, nil
}

func (k Keeper) GetDetailsToCheck(ctx sdk.Context, collection *types.TokenCollection, oldApprovals []*types.CollectionApproval, newApprovals []*types.CollectionApproval) ([]*types.UniversalPermissionDetails, error) {
	x := [][]*types.UintRange{}
	x = append(x, []*types.UintRange{
		// Dummmy range since collection approvals dont use timeline times
		{
			Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64),
		},
	})

	y := [][]*types.UintRange{}
	y = append(y, []*types.UintRange{
		// Dummmy range since collection approvals dont use timeline times
		{
			Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64),
		},
	})

	// This is just to maintain consistency with the legacy features when we used to have timeline times
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, x, []interface{}{oldApprovals})
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, y, []interface{}{newApprovals})

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.CollectionApproval{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
		// This is a little different from the other functions because it is not first match only

		// Expand all collection approved transfers so that they are manipulated according to options and approvalCriteria / allowedCombinations are len 1
		oldApprovals := oldValue.([]*types.CollectionApproval)
		newApprovals := newValue.([]*types.CollectionApproval)

		// Step 1: Merge so we get approvalCriteria arrays of proper length such that it is first match and each (to, from, init, time, ids, ownershipTimes) is only seen once
		// Step 2: Compare as we had previously

		// Step 1:
		oldApprovalsCasted, err := k.CastCollectionApprovalToUniversalPermission(ctx, oldApprovals)
		if err != nil {
			return nil, err
		}
		firstMatchesForOld, err := GetFirstMatchOnlyWithApprovalCriteria(ctx, oldApprovalsCasted)
		if err != nil {
			return nil, err
		}

		newApprovalsCasted, err := k.CastCollectionApprovalToUniversalPermission(ctx, newApprovals)
		if err != nil {
			return nil, err
		}
		firstMatchesForNew, err := GetFirstMatchOnlyWithApprovalCriteria(ctx, newApprovalsCasted)
		if err != nil {
			return nil, err
		}

		// Step 2:
		// For every token, we need to check if the new provided value is different in any way from the old value for each token ID
		// The overlapObjects from GetOverlapsAndNonOverlaps will return which token IDs overlap
		// Note this okay since we already converted everything to first match only in the previous step
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(ctx, firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {
				oldArbVal := oldDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)
				newArbVal := newDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)

				oldVal := oldArbVal
				newVal := newArbVal

				// Go one by one comparing old to new as flat array (if 2d array is empty we still treat it as an empty element
				if len(oldVal) != len(newVal) {
					different = true
				} else {
					for i := 0; i < len(oldVal); i++ {
						oldApprovalCriteria := oldVal[i].ApprovalCriteria
						newApprovalCriteria := newVal[i].ApprovalCriteria
						// Check approval criteria changes using stringification
						if proto.MarshalTextString(oldApprovalCriteria) != proto.MarshalTextString(newApprovalCriteria) {
							different = true
						}
						// Check URI changes
						if oldVal[i].Uri != newVal[i].Uri {
							different = true
						}
						// Check custom data changes
						if oldVal[i].CustomData != newVal[i].CustomData {
							different = true
						}
					}
				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		// If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return nil, err
	}

	return detailsToCheck, nil
}

// ValidateCollectionApprovalsWithInvariants validates collection approvals with invariants using keeper context
// This allows proper checking of address lists (e.g., whether FromListId includes "Mint")
func (k Keeper) ValidateCollectionApprovalsWithInvariants(ctx sdk.Context, collectionApprovals []*types.CollectionApproval, canChangeValues bool, collection *types.TokenCollection) error {
	// First validate the basic collection approvals
	if err := types.ValidateCollectionApprovals(ctx, collectionApprovals, canChangeValues); err != nil {
		return err
	}

	// Check invariants if collection is provided
	if collection != nil && collection.Invariants != nil {
		if collection.Invariants.NoCustomOwnershipTimes {
			for _, collectionApproval := range collectionApprovals {
				if err := types.ValidateNoCustomOwnershipTimesInvariant(collectionApproval.OwnershipTimes, true); err != nil {
					return err
				}
			}
		}

		// Check noForcefulPostMintTransfers invariant
		// This only applies to approvals where FromListId is not only "Mint"
		if collection.Invariants.NoForcefulPostMintTransfers {
			for _, collectionApproval := range collectionApprovals {
				if collectionApproval.ApprovalCriteria != nil {
					// Get the address list object and check if it's only Mint
					fromList, err := k.GetAddressListById(ctx, collectionApproval.FromListId)
					if err == nil && isOnlyMint(fromList) {
						// FromListId is only Mint, so skip the check (allow forceful transfers from Mint)
						continue
					}

					// FromListId is not only Mint, so disallow forceful transfers
					if collectionApproval.ApprovalCriteria.OverridesFromOutgoingApprovals {
						return sdkerrors.Wrapf(types.ErrInvalidRequest, "collection approval %s has overridesFromOutgoingApprovals set to true, which is disallowed when noForcefulPostMintTransfers invariant is enabled (unless FromListId is only Mint)", collectionApproval.ApprovalId)
					}
					if collectionApproval.ApprovalCriteria.OverridesToIncomingApprovals {
						return sdkerrors.Wrapf(types.ErrInvalidRequest, "collection approval %s has overridesToIncomingApprovals set to true, which is disallowed when noForcefulPostMintTransfers invariant is enabled (unless FromListId is only Mint)", collectionApproval.ApprovalId)
					}
				}
			}
		}
	}

	return nil
}

// isOnlyMint checks if an address list contains only Mint
// It must be: whitelist = true, len(addresses) == 1, and addresses[0] == "Mint"
func isOnlyMint(addressList *types.AddressList) bool {
	return addressList != nil &&
		addressList.Whitelist &&
		len(addressList.Addresses) == 1 &&
		addressList.Addresses[0] == types.MintAddress
}

func (k Keeper) ValidateCollectionApprovalsUpdate(ctx sdk.Context, collection *types.TokenCollection, oldApprovals []*types.CollectionApproval, newApprovals []*types.CollectionApproval, CanUpdateCollectionApprovals []*types.CollectionApprovalPermission) error {
	// Validate new approvals with invariants (with keeper context for proper address list checking)
	if err := k.ValidateCollectionApprovalsWithInvariants(ctx, newApprovals, true, collection); err != nil {
		return err
	}

	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, oldApprovals, newApprovals)
	if err != nil {
		return err
	}

	err = k.CheckIfCollectionApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserOutgoingApprovalsUpdate(ctx sdk.Context, collection *types.TokenCollection, oldApprovals []*types.UserOutgoingApproval, newApprovals []*types.UserOutgoingApproval, CanUpdateCollectionApprovals []*types.UserOutgoingApprovalPermission, fromAddress string) error {
	old := types.CastOutgoingTransfersToCollectionTransfers(oldApprovals, fromAddress)
	new := types.CastOutgoingTransfersToCollectionTransfers(newApprovals, fromAddress)

	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, old, new)
	if err != nil {
		return err
	}

	err = k.CheckIfUserOutgoingApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserIncomingApprovalsUpdate(ctx sdk.Context, collection *types.TokenCollection, oldApprovals []*types.UserIncomingApproval, newApprovals []*types.UserIncomingApproval, CanUpdateCollectionApprovals []*types.UserIncomingApprovalPermission, toAddress string) error {
	old := types.CastIncomingTransfersToCollectionTransfers(oldApprovals, toAddress)
	new := types.CastIncomingTransfersToCollectionTransfers(newApprovals, toAddress)

	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, old, new)
	if err != nil {
		return err
	}

	err = k.CheckIfUserIncomingApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateTokenMetadataUpdate(ctx sdk.Context, oldTokenMetadata []*types.TokenMetadata, newTokenMetadata []*types.TokenMetadata, canUpdateTokenMetadata []*types.TokenIdsActionPermission) error {
	// Cast to UniversalPermissionDetails for compatibility with overlap functions and get first matches only (i.e. first match for each token ID)
	oldCasted := k.CastTokenMetadataToUniversalPermission(oldTokenMetadata)
	firstMatchesForOld := types.GetFirstMatchOnly(ctx, oldCasted)

	newCasted := k.CastTokenMetadataToUniversalPermission(newTokenMetadata)
	firstMatchesForNew := types.GetFirstMatchOnly(ctx, newCasted)

	// For every token, we need to check if the new provided value is different in any way from the old value for each specific token ID
	// The overlapObjects from GetOverlapsAndNonOverlaps will return which token IDs overlap
	detailsToCheck := []*types.UniversalPermissionDetails{}
	overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(ctx, firstMatchesForOld, firstMatchesForNew)

	for _, overlapObject := range overlapObjects {
		overlap := overlapObject.Overlap
		oldDetails := overlapObject.FirstDetails
		newDetails := overlapObject.SecondDetails

		// ArbitraryValue contains the JSON-marshaled metadata string
		if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
			detailsToCheck = append(detailsToCheck, overlap)
		} else if oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue != nil {
			oldVal, ok1 := oldDetails.ArbitraryValue.(string)
			newVal, ok2 := newDetails.ArbitraryValue.(string)

			if !ok1 || !ok2 {
				// Type assertion failed - treat as changed to be safe
				detailsToCheck = append(detailsToCheck, overlap)
			} else if newVal != oldVal {
				detailsToCheck = append(detailsToCheck, overlap)
			}
		}
		// If both are nil, no change, so we don't add to detailsToCheck
	}

	// If metadata is in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
	detailsToCheck = append(detailsToCheck, inOldButNotNew...)
	detailsToCheck = append(detailsToCheck, inNewButNotOld...)

	// Only check permissions if there are actually changes
	if len(detailsToCheck) > 0 {
		if err := k.CheckIfTokenIdsActionPermissionPermits(ctx, detailsToCheck, canUpdateTokenMetadata, "update token metadata"); err != nil {
			return err
		}
	}

	return nil
}

func GetUpdatedCollectionMetadataCombinations(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
	x := []*types.UniversalPermissionDetails{}
	if (oldValue == nil && newValue != nil) || (oldValue != nil && newValue == nil) {
		x = append(x, &types.UniversalPermissionDetails{})
	} else {
		oldVal := oldValue.(*types.CollectionMetadata)
		newVal := newValue.(*types.CollectionMetadata)

		if oldVal.Uri != newVal.Uri || oldVal.CustomData != newVal.CustomData {
			x = append(x, &types.UniversalPermissionDetails{})
		}
	}
	return x, nil
}

func (k Keeper) ValidateCollectionMetadataUpdate(ctx sdk.Context, oldCollectionMetadata *types.CollectionMetadata, newCollectionMetadata *types.CollectionMetadata, canUpdateCollectionMetadata []*types.ActionPermission) error {
	// Check if value changed
	if oldCollectionMetadata == nil && newCollectionMetadata == nil {
		return nil
	}
	if (oldCollectionMetadata == nil && newCollectionMetadata != nil) || (oldCollectionMetadata != nil && newCollectionMetadata == nil) {
		// Value changed, need to check permissions
	} else if oldCollectionMetadata.Uri == newCollectionMetadata.Uri && oldCollectionMetadata.CustomData == newCollectionMetadata.CustomData {
		return nil
	}

	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canUpdateCollectionMetadata, "update collection metadata"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateManagerUpdate(ctx sdk.Context, oldManager string, newManager string, canUpdateManager []*types.ActionPermission) error {
	// Check if value changed
	if oldManager == newManager {
		return nil
	}

	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canUpdateManager, "update manager"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCustomDataUpdate(ctx sdk.Context, oldCustomData string, newCustomData string, canUpdateCustomData []*types.ActionPermission) error {
	// Check if value changed
	if oldCustomData == newCustomData {
		return nil
	}

	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canUpdateCustomData, "update custom data"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateStandardsUpdate(ctx sdk.Context, oldStandards []string, newStandards []string, canUpdateStandards []*types.ActionPermission) error {
	// Check if value changed
	if len(oldStandards) != len(newStandards) {
		// Values changed, need to check permissions
	} else {
		allSame := true
		for i := 0; i < len(oldStandards); i++ {
			if oldStandards[i] != newStandards[i] {
				allSame = false
				break
			}
		}
		if allSame {
			return nil
		}
	}

	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canUpdateStandards, "update standards"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateIsArchivedUpdate(ctx sdk.Context, oldIsArchived bool, newIsArchived bool, canUpdateIsArchived []*types.ActionPermission) error {
	// Check if value changed
	if oldIsArchived == newIsArchived {
		return nil
	}

	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canUpdateIsArchived, "update is archived"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateAliasPathsAdd(ctx sdk.Context, canAddMoreAliasPaths []*types.ActionPermission) error {
	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canAddMoreAliasPaths, "add alias paths"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCosmosCoinWrapperPathsAdd(ctx sdk.Context, canAddMoreCosmosCoinWrapperPaths []*types.ActionPermission) error {
	// Check permissions
	if err := k.CheckIfActionPermissionPermits(ctx, canAddMoreCosmosCoinWrapperPaths, "add cosmos coin wrapper paths"); err != nil {
		return err
	}

	return nil
}
