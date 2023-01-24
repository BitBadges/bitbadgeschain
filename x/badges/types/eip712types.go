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

//Certain fields are omitted (when uint64 is 0, bool is false, etc) when serialized with Proto and Amino. EIP712 doesn't support optional fields, so we add the omitted empty values back in here.
func NormalizeEIP712TypedData(typedData apitypes.TypedData, msgType string) (apitypes.TypedData, error) {
	actualMsgValueTypes, addUri, addIdRange, addWhitelistInfo, addSubassetSupply := GetMsgValueTypes(msgType)
	typedData.Types["MsgValue"] = actualMsgValueTypes
	if addUri {
		typedData.Types["UriObject"] = []apitypes.Type{
			{Name: "decodeScheme", Type: "uint64"},
			{Name: "scheme", Type: "uint64"},
			{Name: "uri", Type: "string"},
			{Name: "idxRangeToRemove", Type: "IdRange"},
			{Name: "insertSubassetBytesIdx", Type: "uint64"},
			{Name: "bytesToInsert", Type: "string"},
			{Name: "insertIdIdx", Type: "uint64"},
		}
	}

	if addIdRange {
		typedData.Types["IdRange"] = []apitypes.Type{
			{Name: "start", Type: "uint64"},
			{Name: "end", Type: "uint64"},
		}
	}

	if addWhitelistInfo {
		typedData.Types["WhitelistMintInfo"] = []apitypes.Type{
			{Name: "addresses", Type: "uint64[]"},
			{Name: "balanceAmounts", Type: "BalanceObject[]"},
		}

		typedData.Types["BalanceObject"] = []apitypes.Type{
			{Name: "balance", Type: "uint64"},
			{Name: "idRanges", Type: "IdRange[]"},
		}
	}

	if addSubassetSupply {
		typedData.Types["SubassetAmountAndSupply"] = []apitypes.Type{
			{Name: "supply", Type: "uint64"},
			{Name: "amount", Type: "uint64"},
		}
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

//first bool is if URI is needed, second is if ID Range is needed
func GetMsgValueTypes(route string) ([]apitypes.Type, bool, bool, bool, bool) {
	switch route {
	case TypeMsgNewCollection:
		//TODO:
		return []apitypes.Type{
			{Name: "creator", Type: "string"},
			{Name: "uri", Type: "UriObject"},
			{Name: "arbitraryBytes", Type: "string"},
			{Name: "permissions", Type: "uint64"},
			{Name: "defaultSubassetSupply", Type: "uint64"},
			{Name: "freezeAddressRanges", Type: "IdRange[]"},
			{Name: "standard", Type: "uint64"},
			{Name: "BadgeSupplys", Type: "SubassetAmountAndSupply[]"},
			{Name: "whitelistedRecipients", Type: "WhitelistMintInfo[]"},
		}, true, true, true, true
	case TypeMsgMintBadge:
		//TODO:
		return []apitypes.Type{
			{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "BadgeSupplys", Type: "SubassetAmountAndSupply[]"},
		}, false, false, false, true
	case TypeMsgTransferBadge:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "from", Type: "uint64"},
			{Name: "toAddresses", Type: "uint64[]"},
			{Name: "amounts", Type: "uint64[]"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "badgeIdRanges", Type: "IdRange[]"},
			{Name: "expirationTime", Type: "uint64"},
			{Name: "cantCancelBeforeTime", Type: "uint64"},
		}, false, true, false, false
	case TypeMsgSetApproval:
		return []apitypes.Type{
			{Name: "creator", Type: "string"},
			{Name: "amount", Type: "uint64"},
			{Name: "address", Type: "uint64"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "badgeIdRanges", Type: "IdRange[]"},
		}, false, true, false, false

	case TypeMsgUpdateDisallowedTransfers:
		return []apitypes.Type{
			//TODO:
			{Name: "creator", Type: "string"},
			{Name: "addressRanges", Type: "IdRange[]"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "add", Type: "bool"},
		}, false, true, false, false
	case TypeMsgUpdateUris:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "uri", Type: "UriObject"},
		}, true, true, false, false
	case TypeMsgUpdatePermissions:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "permissions", Type: "uint64"},
		}, false, false, false, false
	case TypeMsgUpdateBytes:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "newBytes", Type: "string"},
		}, false, false, false, false
	case TypeMsgTransferManager:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "address", Type: "uint64"},
		}, false, false, false, false
	case TypeMsgRequestTransferManager:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "badgeId", Type: "uint64"},
			{Name: "add", Type: "bool"},
		}, false, false, false, false
	case TypeMsgRegisterAddresses:
		return []apitypes.Type{{Name: "creator", Type: "string"},
			{Name: "addressesToRegister", Type: "string[]"},
		}, false, false, false, false
	default:
		return []apitypes.Type{}, false, false, false, false
	}
}
