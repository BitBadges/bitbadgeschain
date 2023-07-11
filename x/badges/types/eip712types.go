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

	// permissionType := []apitypes.Type{
	// 	{Name: "isFrozen", Type: "bool"},
	// 	{Name: "timeIntervals", Type: "UintRange[]"},
	// }

	// permissionsTypes := []apitypes.Type{
	// 	{Name: "canArchive", Type: "Permission"},
	// 	{Name: "canUpdateContractAddress", Type: "Permission"},
	// 	{Name: "canUpdateOffChainBalancesMetadata", Type: "Permission"},
	// 	{Name: "canDeleteCollection", Type: "Permission"},
	// 	{Name: "canUpdateCustomData", Type: "Permission"},
	// 	{Name: "canUpdateManager", Type: "Permission"},
	// 	{Name: "canUpdateCollectionMetadata", Type: "Permission"},
	// 	{Name: "canUpdateBadgeMetadata", Type: "Permission"},
	// 	{Name: "canCreateMoreBadges", Type: "Permission"},
	// 	{Name: "canUpdateCollectionApprovedTransfers", Type: "Permission"},
	// }

	// collectionApprovedTransferTypes := []apitypes.Type{
	// 	{Name: "from", Type: "AddressMapping"},
	// 	{Name: "to", Type: "AddressMapping"},
	// 	{Name: "initiatedBy", Type: "AddressMapping"},
	// 	{Name: "includeMints", Type: "bool"},
	// 	{Name: "isAllowed", Type: "bool"},
	// 	{Name: "isFrozen", Type: "bool"},
	// 	{Name: "noApprovalRequired", Type: "bool"},
	// 	{Name: "badgeIds", Type: "UintRange[]"},
	// 	{Name: "timeIntervals", Type: "UintRange[]"},
	// }

	// addressesTypes := []apitypes.Type{
	// 	{Name: "addresses", Type: "string[]"},
	// 	{Name: "includeAddresses", Type: "bool"},
	// }

	// uintRangeTypes := []apitypes.Type{
	// 	{Name: "start", Type: "uint64"},
	// 	{Name: "end", Type: "uint64"},
	// }

	// balanceTypes := []apitypes.Type{
	// 	{Name: "amount", Type: "uint64"},
	// 	{Name: "badgeIds", Type: "UintRange[]"},
	// }

	// badgeSupplyAndAmountTypes := []apitypes.Type{
	// 	{Name: "supply", Type: "uint64"},
	// 	{Name: "amount", Type: "uint64"},
	// }

	// transfersTypes := []apitypes.Type{
	// 	{Name: "toAddresses", Type: "string[]"},
	// 	{Name: "balances", Type: "Balance[]"},
	// }

	// claimsTypes := []apitypes.Type{
	// 	{Name: "balances", Type: "Balance[]"},
	// 	{Name: "timeIntervals", Type: "UintRange[]"},
	// 	{Name: "uri", Type: "string"},
	// 	{Name: "numClaimsPerAddress", Type: "uint64"},
	// 	{Name: "IncrementBadgeIdsBy", Type: "uint64"},
	// 	{Name: "startingClaimAmounts", Type: "Balance[]"},
	// 	{Name: "challenges", Type: "Challenge[]"},
	// 	{Name: "totalClaimsProcessed", Type: "uint64"},
	// 	{Name: "isAssignable", Type: "bool"},
	// 	{Name: "createdBy", Type: "string"},
	// 	{Name: "isReturnable", Type: "bool"},
	// }

	// proofItemTypes := []apitypes.Type{
	// 	{Name: "aunt", Type: "string"},
	// 	{Name: "onRight", Type: "bool"},
	// }

	// proofTypes := []apitypes.Type{
	// 	{Name: "aunts", Type: "MerklePathItem[]"},
	// 	{Name: "leaf", Type: "string"},
	// }

	// badgeMetadataType := []apitypes.Type{
	// 	{Name: "uri", Type: "string"},
	// 	{Name: "badgeIds", Type: "UintRange[]"},
	// 	{Name: "customData", Type: "string"},
	// 	{Name: "isFrozen", Type: "bool"},
	// }

	// collectionMetadataType := []apitypes.Type{
	// 	{Name: "uri", Type: "string"},
	// 	{Name: "customData", Type: "string"},
	// }

	// offChainBalancesMetadataType := []apitypes.Type{
	// 	{Name: "uri", Type: "string"},
	// 	{Name: "customData", Type: "string"},
	// }

	// challengeType := []apitypes.Type{
	// 	{Name: "root", Type: "string"},
	// 	{Name: "expectedProofLength", Type: "uint64"},
	// 	{Name: "isWhitelistTree", Type: "bool"},
	// 	{Name: "maxOneUsePerLeaf", Type: "bool"},
	// 	{Name: "useLeafIndexForBadgeIds", Type: "bool"},
	// }

	// challengeSolutionType := []apitypes.Type{
	// 	{Name: "proof", Type: "ClaimProof"},
	// }

	switch route {

	// case TypeMsgDeleteCollection:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 		},
	// 	}

	// case TypeMsgNewCollection:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionMetadata", Type: "CollectionMetadata"},
	// 			{Name: "badgeMetadata", Type: "BadgeMetadata[]"},
	// 			{Name: "offChainBalancesMetadata", Type: "OffChainBalancesMetadata"},
	// 			{Name: "customData", Type: "string"},
	// 			{Name: "permissions", Type: "Permissions"},
	// 			{Name: "approvedTransfers", Type: "CollectionApprovedTransfer[]"},
	// 			{Name: "standard", Type: "uint64"},
	// 			{Name: "badgesToCreate", Type: "BadgeSupplyAndAmount[]"},
	// 			{Name: "transfers", Type: "Transfer[]"},
	// 			{Name: "claims", Type: "Claim[]"},
	// 			{Name: "contractAddress", Type: "string"},
	// 		},
	// 		"CollectionApprovedTransfer": collectionApprovedTransferTypes,
	// 		"AddressMapping":           addressesTypes,
	// 		"UintRange":                    uintRangeTypes,
	// 		"BadgeSupplyAndAmount":       badgeSupplyAndAmountTypes,
	// 		"Transfer":                   transfersTypes,
	// 		"Claim":                      claimsTypes,
	// 		"Balance":                    balanceTypes,
	// 		"BadgeMetadata":              badgeMetadataType,
	// 		"CollectionMetadata":         collectionMetadataType,
	// 		"OffChainBalancesMetadata":           offChainBalancesMetadataType,
	// 		"Challenge":                  challengeType,
	// 		"Permissions":                permissionsTypes,
	// 		"Permission":                 permissionType,
	// 	}
	// case TypeMsgMintAndDistributeBadges:

	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "badgesToCreate", Type: "BadgeSupplyAndAmount[]"},
	// 			{Name: "transfers", Type: "Transfer[]"},
	// 			{Name: "claims", Type: "Claim[]"},
	// 			{Name: "collectionMetadata", Type: "CollectionMetadata"},
	// 			{Name: "badgeMetadata", Type: "BadgeMetadata[]"},
	// 			{Name: "offChainBalancesMetadata", Type: "OffChainBalancesMetadata"},
	// 			{Name: "approvedTransfers", Type: "CollectionApprovedTransfer[]"},
	// 		},
	// 		"BadgeSupplyAndAmount":       badgeSupplyAndAmountTypes,
	// 		"Transfer":                   transfersTypes,
	// 		"Claim":                      claimsTypes,
	// 		"Balance":                    balanceTypes,
	// 		"UintRange":                    uintRangeTypes,
	// 		"BadgeMetadata":              badgeMetadataType,
	// 		"CollectionMetadata":         collectionMetadataType,
	// 		"OffChainBalancesMetadata":           offChainBalancesMetadataType,
	// 		"Challenge":                  challengeType,
	// 		"CollectionApprovedTransfer": collectionApprovedTransferTypes,
	// 	}

	// case TypeMsgTransferBadges:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "from", Type: "string"},
	// 			{Name: "transfers", Type: "Transfer[]"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 		},
	// 		"Transfer": transfersTypes,
	// 		"Balance":  balanceTypes,
	// 		"UintRange":  uintRangeTypes,
	// 	}
	// case TypeMsgUpdateUserApprovedTransfers:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "approvedTransfers", Type: "UserApprovedOutgoingTransfer[]"},
	// 		},
	// 		"UserApprovedOutgoingTransfer": 			UserApprovedOutgoingTransferTypes,
	// 		"AddressMapping":           addressesTypes,
	// 		"UintRange":                    uintRangeTypes,
	// 	}
	// case TypeMsgUpdateCollectionApprovedTransfers:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "approvedTransfers", Type: "CollectionApprovedTransfer[]"},
	// 		},
	// 		"CollectionApprovedTransfer": collectionApprovedTransferTypes,
	// 		"AddressMapping":           addressesTypes,
	// 		"UintRange":                    uintRangeTypes,
	// 	}
	// case TypeMsgUpdateMetadata:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "collectionMetadata", Type: "CollectionMetadata"},
	// 			{Name: "badgeMetadata", Type: "BadgeMetadata[]"},
	// 			{Name: "offChainBalancesMetadata", Type: "OffChainBalancesMetadata"},
	// 			{Name: "customData", Type: "string"},
	// 			{Name: "contractAddress", Type: "string"},
	// 		},
	// 		"BadgeMetadata":      badgeMetadataType,
	// 		"CollectionMetadata": collectionMetadataType,
	// 		"OffChainBalancesMetadata":   offChainBalancesMetadataType,
	// 		"UintRange":            uintRangeTypes,
	// 	}
	// case TypeMsgUpdateCollectionPermissions:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "permissions", Type: "Permissions"},
	// 		},
	// 		"Permissions": permissionsTypes,
	// 		"Permission":  permissionType,
	// 	}
	// case TypeMsgUpdateManager:
	// 	return map[string][]apitypes.Type{
	// 		"MsgValue": {
	// 			{Name: "creator", Type: "string"},
	// 			{Name: "collectionId", Type: "uint64"},
	// 			{Name: "address", Type: "string"},
	// 		},
	// 	}
	default:
		return map[string][]apitypes.Type{}
	}
}
