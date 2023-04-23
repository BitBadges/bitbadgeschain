package types

import (
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		} else if typeStr == "AddressOptions" && value == nil {
			mapObject[typeObj.Name] = "0"
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

	transferMappingTypes := []apitypes.Type{
		{Name: "from", Type: "Addresses"},
		{Name: "to", Type: "Addresses"},
	}

	addressesTypes := []apitypes.Type{
		{Name: "accountIds", Type: "IdRange[]"},
		{Name: "options", Type: "uint64"},
	}

	idRangeTypes := []apitypes.Type{
		{Name: "start", Type: "uint64"},
		{Name: "end", Type: "uint64"},
	}

	balanceTypes := []apitypes.Type{
		{Name: "balance", Type: "uint64"},
		{Name: "badgeIds", Type: "IdRange[]"},
	}

	badgeSupplyAndAmountTypes := []apitypes.Type{
		{Name: "supply", Type: "uint64"},
		{Name: "amount", Type: "uint64"},
	}

	transfersTypes := []apitypes.Type{
		{Name: "toAddresses", Type: "uint64[]"},
		{Name: "balances", Type: "Balance[]"},
	}

	claimsTypes := []apitypes.Type{
		{Name: "balances", Type: "Balance[]"},
		{Name: "codeRoot", Type: "string"},
		{Name: "whitelistRoot", Type: "string"},
		{Name: "incrementIdsBy", Type: "uint64"},
		{Name: "amount", Type: "uint64"},
		{Name: "badgeIds", Type: "IdRange[]"},
		{Name: "restrictOptions", Type: "uint64"},
		{Name: "uri", Type: "string"},
		{Name: "timeRange", Type: "IdRange"},
		{Name: "expectedMerkleProofLength", Type: "uint64"},
	}

	proofItemTypes := []apitypes.Type{
		{Name: "aunt", Type: "string"},
		{Name: "onRight", Type: "bool"},
	}

	proofTypes := []apitypes.Type{
		{Name: "aunts", Type: "ClaimProofItem[]"},
		{Name: "leaf", Type: "string"},
	}

	badgeUrisType := []apitypes.Type{
		{Name: "uri", Type: "string"},
		{Name: "badgeIds", Type: "IdRange[]"},
	}

	switch route {

	case TypeMsgDeleteCollection:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
			},
		}

	case TypeMsgNewCollection:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionUri", Type: "string"},
				{Name: "badgeUris", Type: "BadgeUri[]"},
				{Name: "bytes", Type: "string"},
				{Name: "permissions", Type: "uint64"},
				{Name: "disallowedTransfers", Type: "TransferMapping[]"},
				{Name: "managerApprovedTransfers", Type: "TransferMapping[]"},
				{Name: "standard", Type: "uint64"},
				{Name: "badgeSupplys", Type: "BadgeSupplyAndAmount[]"},
				{Name: "transfers", Type: "Transfers[]"},
				{Name: "claims", Type: "Claim[]"},
			},
			"TransferMapping":      transferMappingTypes,
			"Addresses":            addressesTypes,
			"IdRange":              idRangeTypes,
			"BadgeSupplyAndAmount": badgeSupplyAndAmountTypes,
			"Transfers":            transfersTypes,
			"Claim":                claimsTypes,
			"Balance":              balanceTypes,
			"BadgeUri":             badgeUrisType,
		}
	case TypeMsgMintBadge:

		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "badgeSupplys", Type: "BadgeSupplyAndAmount[]"},
				{Name: "transfers", Type: "Transfers[]"},
				{Name: "claims", Type: "Claim[]"},
				{Name: "collectionUri", Type: "string"},
				{Name: "badgeUris", Type: "BadgeUri[]"},
			},
			"BadgeSupplyAndAmount": badgeSupplyAndAmountTypes,
			"Transfers":            transfersTypes,
			"Claim":                claimsTypes,
			"Balance":              balanceTypes,
			"IdRange":              idRangeTypes,
			"BadgeUri":             badgeUrisType,
		}

	case TypeMsgTransferBadge:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "from", Type: "uint64"},
				{Name: "transfers", Type: "Transfers[]"},
				{Name: "collectionId", Type: "uint64"},
			},
			"Transfers": transfersTypes,
			"Balance":   balanceTypes,
			"IdRange":   idRangeTypes,
		}
	case TypeMsgSetApproval:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "address", Type: "uint64"},
				{Name: "balances", Type: "Balance[]"},
			},
			"Balance": balanceTypes,
			"IdRange": idRangeTypes,
		}
	case TypeMsgUpdateDisallowedTransfers:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "disallowedTransfers", Type: "TransferMapping[]"},
			},
			"TransferMapping": transferMappingTypes,
			"Addresses":       addressesTypes,
			"IdRange":         idRangeTypes,
		}
	case TypeMsgUpdateUris:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "collectionUri", Type: "string"},
				{Name: "badgeUris", Type: "BadgeUri[]"},
			},
			"BadgeUri": badgeUrisType,
			"IdRange":  idRangeTypes,
		}
	case TypeMsgUpdatePermissions:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "permissions", Type: "uint64"},
			},
		}
	case TypeMsgUpdateBytes:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "newBytes", Type: "string"},
			},
		}
	case TypeMsgTransferManager:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "address", Type: "uint64"},
			},
		}
	case TypeMsgRequestTransferManager:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "addRequest", Type: "bool"},
			},
		}
	case TypeMsgRegisterAddresses:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "addressesToRegister", Type: "string[]"},
			},
		}
	case TypeMsgClaimBadge:
		return map[string][]apitypes.Type{
			"MsgValue": {
				{Name: "creator", Type: "string"},
				{Name: "claimId", Type: "uint64"},
				{Name: "collectionId", Type: "uint64"},
				{Name: "whitelistProof", Type: "ClaimProof"},
				{Name: "codeProof", Type: "ClaimProof"},
			},
			"ClaimProof":     proofTypes,
			"ClaimProofItem": proofItemTypes,
		}
	default:
		return map[string][]apitypes.Type{}
	}
}
