package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// NormalizeEmptyTypes is a recursive function that adds empty values to fields that are omitted when serialized with Proto and Amino.
// EIP712 doesn't support optional fields, so we add the omitted empty values back in here.
// This includes empty strings, when a uint64 is 0, and when a bool is false.
func NormalizeEmptyTypes(typedData apitypes.TypedData, typeObjArr []apitypes.Type, mapObject map[string]interface{}) (map[string]interface{}, error) {
	for _, typeObj := range typeObjArr {
		typeStr := typeObj.Type
		value := mapObject[typeObj.Name]

		if strings.Contains(typeStr, "[]") && value == nil {
			mapObject[typeObj.Name] = []interface{}{}
		} else if strings.Contains(typeStr, "[]") && value != nil {
			valueArr := value.([]interface{})
			// Get typeStr without the brackets at the end
			typeStr = typeStr[:len(typeStr)-2]
			for i, value := range valueArr {
				innerMap, ok := value.(map[string]interface{})
				if ok {
					newMap, err := NormalizeEmptyTypes(typedData, typedData.Types[typeStr], innerMap)
					if err != nil {
						return mapObject, err
					}
					valueArr[i] = newMap
				}
			}
			mapObject[typeObj.Name] = valueArr
		} else if typeStr == "string" && value == nil {
			mapObject[typeObj.Name] = ""
		} else if typeStr == "uint64" && value == nil {
			mapObject[typeObj.Name] = "0"
		} else if typeStr == "bool" && value == nil {
			mapObject[typeObj.Name] = false
		} else {
			innerMap, ok := mapObject[typeObj.Name].(map[string]interface{})
			if ok {
				newMap, err := NormalizeEmptyTypes(typedData, typedData.Types[typeStr], innerMap)
				if err != nil {
					return mapObject, err
				}
				mapObject[typeObj.Name] = newMap
			}
		}
	}
	return mapObject, nil
}

// Certain fields are omitted (when uint64 is 0, bool is false, etc) when serialized with Proto and Amino. EIP712 doesn't support optional fields, so we add the omitted empty values back in here.
func NormalizeEIP712TypedData(typedData apitypes.TypedData, msgType string) (apitypes.TypedData, error) {
	typesMap := GetMsgValueTypes(msgType)
	for key, value := range typesMap {
		typedData.Types[key] = value
	}

	//Remove the types in typedData.Types that begin with the prefix Type
	//This is to get around the logic in the ethermint eip712 logic which assigns
	//a new type to each unknown field by Type + FieldName.
	//We don't want to use this because it is wrong and we add the types in this file.
	for key := range typedData.Types {
		if strings.HasPrefix(key, "Type") {
			delete(typedData.Types, key)
		}
	}

	msgValue, ok := typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"].(map[string]interface{})
	if !ok {
		return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "message is not a map[string]interface{}")
	}

	normalizedMsgValue, err := NormalizeEmptyTypes(typedData, typedData.Types["MsgValue"], msgValue)
	if err != nil {
		return typedData, err
	}

	typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"] = normalizedMsgValue

	return typedData, nil
}

// first bool is if URI is needed, second is if ID Range is needed
func GetMsgValueTypes(route string) map[string][]apitypes.Type {
	valueOptionsTypes := []apitypes.Type{
		{Name: "invertDefault", Type: "bool"},
		{Name: "allValues", Type: "bool"}, //Override default values with all possible values
		{Name: "noValues", Type: "bool"}, //Override default values with no values
	}

	uintRangeTypes := []apitypes.Type{
		{Name: "start", Type: "string"},
		{Name: "end", Type: "string"},
	}

	balanceTypes := []apitypes.Type{
		{Name: "amount", Type: "string"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
	}

	proofItemTypes := []apitypes.Type{
		{Name: "aunt", Type: "string"},
		{Name: "onRight", Type: "bool"},
	}

	proofTypes := []apitypes.Type{
		{Name: "aunts", Type: "MerklePathItem[]"},
		{Name: "leaf", Type: "string"},
	}
	// UserApprovedOutgoingTransferTimelineTypes := []apitypes.Type{
	// 	{Name: "approvedOutgoingTransfers", Type: "UserApprovedOutgoingTransfer[]"},
	// 	{Name: "timelineTimes", Type: "UintRange[]"},
	// }
	// UserApprovedIncomingTransferTimelineTypes := []apitypes.Type{
	// 	{Name: "approvedIncomingTransfers", Type: "UserApprovedIncomingTransfer[]"},
	// 	{Name: "timelineTimes", Type: "UintRange[]"},
	// }
	UserPermissionsTypes := []apitypes.Type{
		{Name: "canUpdateApprovedOutgoingTransfers", Type: "UserApprovedOutgoingTransferPermission[]"},
		{Name: "canUpdateApprovedIncomingTransfers", Type: "UserApprovedIncomingTransferPermission[]"},
	}
	UserApprovedOutgoingTransferTypes := []apitypes.Type{
		{Name: "toMappingId", Type: "string"},
		{Name: "initiatedByMappingId", Type: "string"},
		{Name: "transferTimes", Type: "UintRange[]"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
		{Name: "approvalId", Type: "string"},
		{Name: "approvalTrackerId", Type: "string"},
		{Name: "challengeTrackerId", Type: "string"},
		{Name: "allowedCombinations", Type: "IsUserOutgoingTransferAllowed[]"},
		{Name: "approvalDetails", Type: "OutgoingApprovalDetails"},
	}
	UserApprovedIncomingTransferTypes := []apitypes.Type{
		{Name: "fromMappingId", Type: "string"},
		{Name: "initiatedByMappingId", Type: "string"},
		{Name: "transferTimes", Type: "UintRange[]"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
		{Name: "approvalId", Type: "string"},
		{Name: "approvalTrackerId", Type: "string"},
		{Name: "challengeTrackerId", Type: "string"},
		{Name: "allowedCombinations", Type: "IsUserIncomingTransferAllowed[]"},
		{Name: "approvalDetails", Type: "IncomingApprovalDetails"},
	}
	UserApprovedOutgoingTransferPermissionTypes := []apitypes.Type{
		{Name: "defaultValues", Type: "UserApprovedOutgoingTransferDefaultValues"},
		{Name: "combinations", Type: "UserApprovedOutgoingTransferCombination[]"},
	}
	UserApprovedIncomingTransferPermissionTypes := []apitypes.Type{
		{Name: "defaultValues", Type: "UserApprovedIncomingTransferDefaultValues"},
		{Name: "combinations", Type: "UserApprovedIncomingTransferCombination[]"},
	}
	UserApprovedOutgoingTransferDefaultValuesTypes := []apitypes.Type{
		{Name: "toMappingId", Type: "string"},
		{Name: "initiatedByMappingId", Type: "string"},
		{Name: "transferTimes", Type: "UintRange[]"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
		{Name: "approvalTrackerId", Type: "string"},
		{Name: "challengeTrackerId", Type: "string"},

		{Name: "permittedTimes", Type: "UintRange[]"},
		{Name: "forbiddenTimes", Type: "UintRange[]"},
	}
	UserApprovedOutgoingTransferCombinationTypes := []apitypes.Type{
		{Name: "toMappingOptions", Type: "ValueOptions"},
		{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
		{Name: "transferTimesOptions", Type: "ValueOptions"},
		{Name: "badgeIdsOptions", Type: "ValueOptions"},
		{Name: "ownershipTimesOptions", Type: "ValueOptions"},
		{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
		{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
		{Name: "permittedTimesOptions", Type: "ValueOptions"},
		{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
	}
	UserApprovedIncomingTransferDefaultValuesTypes := []apitypes.Type{
		{Name: "fromMappingId", Type: "string"},
		{Name: "initiatedByMappingId", Type: "string"},
		{Name: "transferTimes", Type: "UintRange[]"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
		{Name: "approvalTrackerId", Type: "string"},
		{Name: "challengeTrackerId", Type: "string"},

		{Name: "permittedTimes", Type: "UintRange[]"},
		{Name: "forbiddenTimes", Type: "UintRange[]"},
	}
	UserApprovedIncomingTransferCombinationTypes := []apitypes.Type{
		{Name: "fromMappingOptions", Type: "ValueOptions"},
		{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
		{Name: "transferTimesOptions", Type: "ValueOptions"},
		{Name: "badgeIdsOptions", Type: "ValueOptions"},
		{Name: "ownershipTimesOptions", Type: "ValueOptions"},
		{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
		{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
		{Name: "permittedTimesOptions", Type: "ValueOptions"},
		{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
	}
	IsUserOutgoingTransferAllowedTypes := []apitypes.Type{
		{Name: "toMappingOptions", Type: "ValueOptions"},
		{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
		{Name: "transferTimesOptions", Type: "ValueOptions"},
		{Name: "badgeIdsOptions", Type: "ValueOptions"},
		{Name: "ownershipTimesOptions", Type: "ValueOptions"},
		{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
		{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
		{Name: "isApproved", Type: "bool"},
	}
	OutgoingApprovalDetailsTypes := []apitypes.Type{
		{Name: "uri", Type: "string"},
		{Name: "customData", Type: "string"},

		{Name: "mustOwnBadges", Type: "MustOwnBadges[]"},
		{Name: "merkleChallenge", Type: "MerkleChallenge"},
		{Name: "predeterminedBalances", Type: "PredeterminedBalances"},
		{Name: "approvalAmounts", Type: "ApprovalAmounts"},
		{Name: "maxNumTransfers", Type: "MaxNumTransfers"},

		{Name: "requireToEqualsInitiatedBy", Type: "bool"},
		{Name: "requireToDoesNotEqualInitiatedBy", Type: "bool"},
	}
	IsUserIncomingTransferAllowedTypes := []apitypes.Type{
		{Name: "fromMappingOptions", Type: "ValueOptions"},
		{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
		{Name: "transferTimesOptions", Type: "ValueOptions"},
		{Name: "badgeIdsOptions", Type: "ValueOptions"},
		{Name: "ownershipTimesOptions", Type: "ValueOptions"},
		{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
		{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
		{Name: "isApproved", Type: "bool"},
	}
	IncomingApprovalDetailsTypes := []apitypes.Type{
		{Name: "uri", Type: "string"},
		{Name: "customData", Type: "string"},

		{Name: "mustOwnBadges", Type: "MustOwnBadges[]"},
		{Name: "merkleChallenge", Type: "MerkleChallenge"},
		{Name: "predeterminedBalances", Type: "PredeterminedBalances"},
		{Name: "approvalAmounts", Type: "ApprovalAmounts"},
		{Name: "maxNumTransfers", Type: "MaxNumTransfers"},

		{Name: "requireFromEqualsInitiatedBy", Type: "bool"},
		{Name: "requireFromDoesNotEqualInitiatedBy", Type: "bool"},
	}
	MustOwnBadgesTypes := []apitypes.Type{
		{Name: "collectionId", Type: "string"},
		{Name: "amountRange", Type: "UintRange"},
		{Name: "badgeIds", Type: "UintRange[]"},
		{Name: "ownershipTimes", Type: "UintRange[]"},
		{Name: "overrideWithCurrentTime", Type: "bool"},
		{Name: "mustOwnAll", Type: "bool"},
	}
	MerkleChallengeTypes := []apitypes.Type{
		{Name: "root", Type: "string"},
		{Name: "expectedProofLength", Type: "string"},
		{Name: "useCreatorAddressAsLeaf", Type: "bool"},
		{Name: "maxOneUsePerLeaf", Type: "bool"},
		{Name: "useLeafIndexForTransferOrder", Type: "bool"},
		{Name: "uri", Type: "string"},
		{Name: "customData", Type: "string"},
	}
	PredeterminedBalancesTypes := []apitypes.Type{
		{Name: "manualBalances", Type: "ManualBalances[]"},
		{Name: "incrementedBalances", Type: "IncrementedBalances"},
		{Name: "orderCalculationMethod", Type: "PredeterminedOrderCalculationMethod"},
	}
	ApprovalAmountsTypes := []apitypes.Type{
		{Name: "overallApprovalAmount", Type: "string"},
		{Name: "perToAddressApprovalAmount", Type: "string"},
		{Name: "perFromAddressApprovalAmount", Type: "string"},
		{Name: "perInitiatedByAddressApprovalAmount", Type: "string"},
	}
	MaxNumTransfersTypes := []apitypes.Type{
		{Name: "overallMaxNumTransfers", Type: "string"},
		{Name: "perToAddressMaxNumTransfers", Type: "string"},
		{Name: "perFromAddressMaxNumTransfers", Type: "string"},
		{Name: "perInitiatedByAddressMaxNumTransfers", Type: "string"},
	}
	ManualBalancesTypes := []apitypes.Type{
		{Name: "balances", Type: "Balance[]"},
	}
	IncrementedBalancesTypes := []apitypes.Type{
		{Name: "startBalances", Type: "Balance[]"},
		{Name: "incrementBadgeIdsBy", Type: "string"},
		{Name: "incrementOwnershipTimesBy", Type: "string"},
	}
	PredeterminedOrderCalculationMethodTypes := []apitypes.Type{
		{Name: "useOverallNumTransfers", Type: "bool"},
		{Name: "usePerToAddressNumTransfers", Type: "bool"},
		{Name: "usePerFromAddressNumTransfers", Type: "bool"},
		{Name: "usePerInitiatedByAddressNumTransfers", Type: "bool"},
		{Name: "useMerkleChallengeLeafIndex", Type: "bool"},
	}
	
	switch route {

	case TypeMsgDeleteCollection:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "string"},
			},
		}
		
	case TypeMsgCreateAddressMappings:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "addressMappings", Type: "AddressMapping[]"},
			},
			"AddressMapping": {
				{Name: "mappingId", Type: "string"},
				{Name: "addresses", Type: "string[]"},
				{Name: "includeAddresses", Type: "bool"},
				{Name: "uri", Type: "string"},
				{Name: "customData", Type: "string"},
			},
		}
	case TypeMsgTransferBadges:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "string"},
				{Name: "transfers", Type: "Transfer[]"},
			},
			"Transfer": {
				{Name: "from", Type: "string"},
				{Name: "toAddresses", Type: "string[]"},
				{Name: "balances", Type: "Balance[]"},
				{Name: "precalculationDetails", Type: "ApprovalIdentifierDetails"},
				{Name: "merkleProofs", Type: "MerkleProof[]"},
				{Name: "memo", Type: "string"},
				{Name: "prioritizedApprovals", Type: "ApprovalIdentifierDetails[]"},
				{Name: "onlyCheckPrioritizedApprovals", Type: "bool"},
			},

			"UintRange": uintRangeTypes,
			"Balance": balanceTypes,
			"ApprovalIdentifierDetails": {
				{Name: "approvalId", Type: "string"},
				{Name: "approvalLevel", Type: "string"},
				{Name: "approverAddress", Type: "string"},
			},
			"MerkleProof": proofTypes,
			"MerklePathItem": proofItemTypes,
		}

	case TypeMsgUpdateUserApprovedTransfers:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "string"},
				{Name: "updateApprovedOutgoingTransfers", Type: "bool"},
				{Name: "approvedOutgoingTransfers", Type: "UserApprovedOutgoingTransfer[]"},
				{Name: "updateApprovedIncomingTransfers", Type: "bool"},
				{Name: "approvedIncomingTransfers", Type: "UserApprovedIncomingTransfer[]"},
				{Name: "updateUserPermissions", Type: "bool"},
				{Name: "userPermissions", Type: "UserPermissions"},
			},
			"UserPermissions": UserPermissionsTypes,
			"UserApprovedOutgoingTransfer": UserApprovedOutgoingTransferTypes,
			"UserApprovedIncomingTransfer": UserApprovedIncomingTransferTypes,
			"UserApprovedOutgoingTransferPermission": UserApprovedOutgoingTransferPermissionTypes,
			"UserApprovedIncomingTransferPermission": UserApprovedIncomingTransferPermissionTypes,
			"UserApprovedOutgoingTransferDefaultValues": UserApprovedOutgoingTransferDefaultValuesTypes,
			"UserApprovedOutgoingTransferCombination": UserApprovedOutgoingTransferCombinationTypes,
			"UserApprovedIncomingTransferDefaultValues": UserApprovedIncomingTransferDefaultValuesTypes,
			"UserApprovedIncomingTransferCombination": UserApprovedIncomingTransferCombinationTypes,
			"Balance": balanceTypes,
			"UintRange": uintRangeTypes,
			"ValueOptions": valueOptionsTypes,
			"IsUserOutgoingTransferAllowed": IsUserOutgoingTransferAllowedTypes,
			"OutgoingApprovalDetails": OutgoingApprovalDetailsTypes,
			"IsUserIncomingTransferAllowed": IsUserIncomingTransferAllowedTypes,
			"IncomingApprovalDetails": IncomingApprovalDetailsTypes,
			"MustOwnBadges": MustOwnBadgesTypes,
			"MerkleChallenge": MerkleChallengeTypes,
			"PredeterminedBalances": PredeterminedBalancesTypes,
			"ApprovalAmounts": ApprovalAmountsTypes,
			"MaxNumTransfers":	MaxNumTransfersTypes,
			"ManualBalances":  ManualBalancesTypes,
			"IncrementedBalances": IncrementedBalancesTypes,
			"PredeterminedOrderCalculationMethod": PredeterminedOrderCalculationMethodTypes,
		}
	
	
	case TypeMsgUpdateCollection:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "string"},
				{Name: "balancesType", Type: "string"},
				{Name: "defaultApprovedOutgoingTransfers", Type: "UserApprovedOutgoingTransfer[]"},
				{Name: "defaultApprovedIncomingTransfers", Type: "UserApprovedIncomingTransfer[]"},
				{Name: "defaultUserPermissions", Type: "UserPermissions"},
				{Name: "badgesToCreate", Type: "Balance[]"},
				{Name: "updateCollectionPermissions", Type: "bool"},
				{Name: "collectionPermissions", Type: "CollectionPermissions"},
				{Name: "updateManagerTimeline", Type: "bool"},
				{Name: "managerTimeline", Type: "ManagerTimeline[]"},
				{Name: "updateCollectionMetadataTimeline", Type: "bool"},
				{Name: "collectionMetadataTimeline", Type: "CollectionMetadataTimeline[]"},
				{Name: "updateBadgeMetadataTimeline", Type: "bool"},
				{Name: "badgeMetadataTimeline", Type: "BadgeMetadataTimeline[]"},
				{Name: "updateOffChainBalancesMetadataTimeline", Type: "bool"},
				{Name: "offChainBalancesMetadataTimeline", Type: "OffChainBalancesMetadataTimeline[]"},
				{Name: "updateCustomDataTimeline", Type: "bool"},
				{Name: "customDataTimeline", Type: "CustomDataTimeline[]"},
				// {Name: "inheritedCollectionId", Type: "string"},
				{Name: "updateCollectionApprovedTransfers", Type: "bool"},
				{Name: "collectionApprovedTransfers", Type: "CollectionApprovedTransfer[]"},
				{Name: "updateStandardsTimeline", Type: "bool"},
				{Name: "standardsTimeline", Type: "StandardsTimeline[]"},
				{Name: "updateContractAddressTimeline", Type: "bool"},
				{Name: "contractAddressTimeline", Type: "ContractAddressTimeline[]"},
				{Name: "updateIsArchivedTimeline", Type: "bool"},
				{Name: "isArchivedTimeline", Type: "IsArchivedTimeline[]"},
			},
			
			"CollectionPermissions": {
				{Name: "canDeleteCollection", Type: "ActionPermission[]"},
				{Name: "canArchiveCollection", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateContractAddress", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateOffChainBalancesMetadata", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateStandards", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateCustomData", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateManager", Type: "TimedUpdatePermission[]"},
				{Name: "canUpdateCollectionMetadata", Type: "TimedUpdatePermission[]"},
				{Name: "canCreateMoreBadges", Type: "BalancesActionPermission[]"},
				{Name: "canUpdateBadgeMetadata", Type: "TimedUpdateWithBadgeIdsPermission[]"},
				{Name: "canUpdateCollectionApprovedTransfers", Type: "CollectionApprovedTransferPermission[]"},
			},
			"ManagerTimeline": {
				{Name: "manager", Type: "string"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"CollectionMetadataTimeline": {
				{Name: "collectionMetadata", Type: "CollectionMetadata"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"BadgeMetadataTimeline": {
				{Name: "badgeMetadata", Type: "BadgeMetadata[]"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"OffChainBalancesMetadataTimeline": {
				{Name: "offChainBalancesMetadata", Type: "OffChainBalancesMetadata"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"CustomDataTimeline": {
				{Name: "customData", Type: "string"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"StandardsTimeline": {
				{Name: "standards", Type: "string[]"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"ContractAddressTimeline": {
				{Name: "contractAddress", Type: "string"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"IsArchivedTimeline": {
				{Name: "isArchived", Type: "bool"},
				{Name: "timelineTimes", Type: "UintRange[]"},
			},
			"CollectionApprovedTransfer": {
				{Name: "fromMappingId", Type: "string"},
				{Name: "toMappingId", Type: "string"},
				{Name: "initiatedByMappingId", Type: "string"},
				{Name: "transferTimes", Type: "UintRange[]"},
				{Name: "badgeIds", Type: "UintRange[]"},
				{Name: "ownershipTimes", Type: "UintRange[]"},
				{Name: "approvalId", Type: "string"},
				{Name: "approvalTrackerId", Type: "string"},
				{Name: "challengeTrackerId", Type: "string"},
				{Name: "allowedCombinations", Type: "IsCollectionTransferAllowed[]"},
				{Name: "approvalDetails", Type: "ApprovalDetails"},
			},
			"UintRange": uintRangeTypes,
			"Balance": balanceTypes,
			"BadgeMetadata": {
				{Name: "uri", Type: "string"},
				{Name: "customData", Type: "string"},
				{Name: "badgeIds", Type: "UintRange[]"},
			},
			"CollectionMetadata": {
				{Name: "uri", Type: "string"},
				{Name: "customData", Type: "string"},
			},
			"OffChainBalancesMetadata": {
				{Name: "uri", Type: "string"},
				{Name: "customData", Type: "string"},
			},
			"IsCollectionTransferAllowed": {
				{Name: "fromMappingOptions", Type: "ValueOptions"},
				{Name: "toMappingOptions", Type: "ValueOptions"},
				{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
				{Name: "transferTimesOptions", Type: "ValueOptions"},
				{Name: "badgeIdsOptions", Type: "ValueOptions"},
				{Name: "ownershipTimesOptions", Type: "ValueOptions"},
				{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
				{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
				{Name: "isApproved", Type: "bool"},
			},
			"ActionPermission": {
				{Name: "defaultValues", Type: "ActionDefaultValues"},
				{Name: "combinations", Type: "ActionCombination[]"},
			},
			"TimedUpdatePermission": {
				{Name: "defaultValues", Type: "TimedUpdateDefaultValues"},
				{Name: "combinations", Type: "TimedUpdateCombination[]"},
			},
			"BalancesActionPermission": {
				{Name: "defaultValues", Type: "BalancesActionDefaultValues"},
				{Name: "combinations", Type: "BalancesActionCombination[]"},
			},
			"TimedUpdateWithBadgeIdsPermission": {
				{Name: "defaultValues", Type: "TimedUpdateWithBadgeIdsDefaultValues"},
				{Name: "combinations", Type: "TimedUpdateWithBadgeIdsCombination[]"},
			},
			"CollectionApprovedTransferPermission": {
				{Name: "defaultValues", Type: "CollectionApprovedTransferDefaultValues"},
				{Name: "combinations", Type: "CollectionApprovedTransferCombination[]"},
			},
			"CollectionApprovedTransferCombination": {
				{Name: "fromMappingOptions", Type: "ValueOptions"},
				{Name: "toMappingOptions", Type: "ValueOptions"},
				{Name: "initiatedByMappingOptions", Type: "ValueOptions"},
				{Name: "transferTimesOptions", Type: "ValueOptions"},
				{Name: "badgeIdsOptions", Type: "ValueOptions"},
				{Name: "ownershipTimesOptions", Type: "ValueOptions"},
				{Name: "approvalTrackerIdOptions", Type: "ValueOptions"},
				{Name: "challengeTrackerIdOptions", Type: "ValueOptions"},
				{Name: "permittedTimesOptions", Type: "ValueOptions"},
				{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
			},
			"CollectionApprovedTransferDefaultValues": {
				{Name: "fromMappingId", Type: "string"},
				{Name: "toMappingId", Type: "string"},
				{Name: "initiatedByMappingId", Type: "string"},
				{Name: "transferTimes", Type: "UintRange[]"},
				{Name: "badgeIds", Type: "UintRange[]"},
				{Name: "ownershipTimes", Type: "UintRange[]"},
				{Name: "approvalTrackerId", Type: "string"},
				{Name: "challengeTrackerId", Type: "string"},
				{Name: "permittedTimes", Type: "UintRange[]"},
				{Name: "forbiddenTimes", Type: "UintRange[]"},
			},
			"BalancesActionCombination": {
				{Name: "badgeIdsOptions", Type: "ValueOptions"},
				{Name: "ownershipTimesOptions", Type: "ValueOptions"},
				{Name: "permittedTimesOptions", Type: "ValueOptions"},
				{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
			},
			"BalancesActionDefaultValues": {
				{Name: "badgeIds", Type: "UintRange[]"},
				{Name: "ownershipTimes", Type: "UintRange[]"},
				{Name: "permittedTimes", Type: "UintRange[]"},
				{Name: "forbiddenTimes", Type: "UintRange[]"},
			},
			"ActionCombination": {
				{Name: "permittedTimesOptions", Type: "ValueOptions"},
				{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
			},
			"ActionDefaultValues": {
				{Name: "permittedTimes", Type: "UintRange[]"},
				{Name: "forbiddenTimes", Type: "UintRange[]"},
			},
			"TimedUpdateCombination": {
				{Name: "timelineTimesOptions", Type: "ValueOptions"},
				{Name: "permittedTimesOptions", Type: "ValueOptions"},
				{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
			},
			"TimedUpdateDefaultValues": {
				{Name: "timelineTimes", Type: "UintRange[]"},
				{Name: "permittedTimes", Type: "UintRange[]"},
				{Name: "forbiddenTimes", Type: "UintRange[]"},
			},
			"TimedUpdateWithBadgeIdsCombination": {
				{Name: "timelineTimesOptions", Type: "ValueOptions"},
				{Name: "badgeIdsOptions", Type: "ValueOptions"},
				{Name: "permittedTimesOptions", Type: "ValueOptions"},
				{Name: "forbiddenTimesOptions", Type: "ValueOptions"},
			},
			"TimedUpdateWithBadgeIdsDefaultValues": {
				{Name: "timelineTimes", Type: "UintRange[]"},
				{Name: "badgeIds", Type: "UintRange[]"},
				{Name: "permittedTimes", Type: "UintRange[]"},
				{Name: "forbiddenTimes", Type: "UintRange[]"},
			},
			"ApprovalDetails": {
				{Name: "uri", Type: "string"},
				{Name: "customData", Type: "string"},
				{Name: "mustOwnBadges", Type: "MustOwnBadges[]"},
				{Name: "merkleChallenge", Type: "MerkleChallenge"},
				{Name: "predeterminedBalances", Type: "PredeterminedBalances"},
				{Name: "approvalAmounts", Type: "ApprovalAmounts"},
				{Name: "maxNumTransfers", Type: "MaxNumTransfers"},
				{Name: "requireToEqualsInitiatedBy", Type: "bool"},
				{Name: "requireToDoesNotEqualInitiatedBy", Type: "bool"},
				{Name: "requireFromEqualsInitiatedBy", Type: "bool"},
				{Name: "requireFromDoesNotEqualInitiatedBy", Type: "bool"},
				{Name: "overridesFromApprovedOutgoingTransfers", Type: "bool"},
				{Name: "overridesToApprovedIncomingTransfers", Type: "bool"},
			},
			"UserPermissions": UserPermissionsTypes,
			"UserApprovedOutgoingTransfer": UserApprovedOutgoingTransferTypes,
			"UserApprovedIncomingTransfer": UserApprovedIncomingTransferTypes,
			"UserApprovedOutgoingTransferPermission": UserApprovedOutgoingTransferPermissionTypes,
			"UserApprovedIncomingTransferPermission": UserApprovedIncomingTransferPermissionTypes,
			"UserApprovedOutgoingTransferDefaultValues": UserApprovedOutgoingTransferDefaultValuesTypes,
			"UserApprovedOutgoingTransferCombination": UserApprovedOutgoingTransferCombinationTypes,
			"UserApprovedIncomingTransferDefaultValues": UserApprovedIncomingTransferDefaultValuesTypes,
			"UserApprovedIncomingTransferCombination": UserApprovedIncomingTransferCombinationTypes,
			"ValueOptions": valueOptionsTypes,
			"IsUserOutgoingTransferAllowed": IsUserOutgoingTransferAllowedTypes,
			"OutgoingApprovalDetails": OutgoingApprovalDetailsTypes,
			"IsUserIncomingTransferAllowed": IsUserIncomingTransferAllowedTypes,
			"IncomingApprovalDetails": IncomingApprovalDetailsTypes,
			"MustOwnBadges": MustOwnBadgesTypes,
			"MerkleChallenge": MerkleChallengeTypes,
			"PredeterminedBalances": PredeterminedBalancesTypes,
			"ApprovalAmounts": ApprovalAmountsTypes,
			"MaxNumTransfers":	MaxNumTransfersTypes,
			"ManualBalances":  ManualBalancesTypes,
			"IncrementedBalances": IncrementedBalancesTypes,
			"PredeterminedOrderCalculationMethod": PredeterminedOrderCalculationMethodTypes,
		}
	default:
		return map[string][]apitypes.Type{}
	}
}
