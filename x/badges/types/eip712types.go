package types

import (
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type TypedDataMessage struct {
	msgs []map[string]interface{}
}

//Certain fields are omitted when serialized with Proto and Amino. EIP712 doesn't support optional fields, so we add the omitted empty values back in here.
func NormalizeEIP712TypedData(typedData apitypes.TypedData, msgType string) (apitypes.TypedData, error) {
	actualMsgValueTypes, addUri, addIdRange, addWhitelistInfo := GetMsgValueTypes(msgType)
	typedData.Types["MsgValue"] = actualMsgValueTypes
	if addUri {
		typedData.Types["UriObject"] = []apitypes.Type{
			{ Name: "decodeScheme", Type: "uint64" },
			{ Name: "scheme", Type: "uint64" },
			{ Name: "uri", Type: "string" },
			{ Name: "idxRangeToRemove", Type: "IdRange" },
			{ Name: "insertSubassetBytesIdx", Type: "uint64" },
			{ Name: "bytesToInsert", Type: "string" },
			{ Name: "insertIdIdx", Type: "uint64" },
		}
	}

	if addIdRange {
		typedData.Types["IdRange"] = []apitypes.Type{ 
			{ Name: "start", Type: "uint64" },
			{ Name: "end", Type: "uint64" },
		}
	}

	if addWhitelistInfo {
		typedData.Types["WhitelistMintInfo"] = []apitypes.Type{
			{ Name: "addresses", Type: "uint64[]" },
			{ Name: "balanceAmounts", Type: "BalanceObject[]" },
		}

		typedData.Types["BalanceObject"] = []apitypes.Type{
			{ Name: "balance", Type: "uint64" },
			{ Name: "idRanges", Type: "IdRange[]" },
		}
	}



	//Add empty values for fields if omitted. EIP712 was signed with the fields, so we have to read them back here
	for _, typeObj := range typedData.Types["MsgValue"] {
		firstMsg, ok := typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"].(map[string]interface{}) //TODO: multuple messages
		if !ok {
			return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "message is not a map[string]interface{}")
		}
		msgValueForType := firstMsg[typeObj.Name]

		if strings.Contains(typeObj.Type, "[]") && msgValueForType == nil {
			firstMsg[typeObj.Name] = []interface{}{}
		} else if typeObj.Type == "string" && msgValueForType == nil {
			firstMsg[typeObj.Name] = ""
		} else if typeObj.Type == "uint64" && msgValueForType == nil {
			firstMsg[typeObj.Name] = "0"
		} else if typeObj.Type == "bool" && msgValueForType == nil {
			firstMsg[typeObj.Name] = false
		} else if typeObj.Type == "IdRange[]" {
			idRanges, ok := msgValueForType.([]interface{})
			if !ok {
				return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "3. message is not a map[string]interface{}")
			}

			for _, idRangeObj := range idRanges {
				idRange, ok := idRangeObj.(map[string]interface{})
				if !ok {
					return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "4. message is not a map[string]interface{}")
				}
				if idRange["start"] == nil {
					idRange["start"] = "0"
				} 
				if idRange["end"] == nil {
					idRange["end"] = "0"
				}
			}
			// if msgValueForType == nil {
			// 	typedData.Message[typeObj.Name] = map[string]interface{}{
			// 		"start": uint64(0),
			// 		"end": uint64(0),
			// 	}
			// }
			// currStart := msgValueForType["start"]
			// currEnd := msgValueForType["start"]
			// currIdRange, ok := 
			// if !ok {
			// 	return typedData, ErrInvalidIdRangeSpecified
			// }
		}
		// 	typedData.Message[typeObj.Name] = apitypes.IdRange{
		// 		Start: uint64(0),
		// 		End: uint64(0),
		// 	}
		// } else if typeObj.Type == "IdRange[]" && msgValueForType == nil {
		// 	typedData.Message[typeObj.Name] = apitypes.IdRange{
		// 		Start: uint64(0),
		// 		End: uint64(0),
		// 	}
		// } else if typeObj.Type == "UriObject" && msgValueForType == nil {
		// 	typedData.Message[typeObj.Name] = apitypes.IdRange{
		// 		Start: uint64(0),
		// 		End: uint64(0),
		// 	}
		// }



		
	}

	//TODO: this is awful code. I'm sorry. I'll fix it later
	//Add empty values for fields if omitted. EIP712 was signed with the fields, so we have to read them back here
	for _, typeObj := range typedData.Types["UriObject"] {
		uriObject, ok := typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"].(map[string]interface{})["uri"].(map[string]interface{})//TODO: multuple messages
		if !ok {
			return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "message is not a map[string]interface{}")
		}
		uriObjValueForType := uriObject[typeObj.Name]

		if strings.Contains(typeObj.Type, "[]") && uriObjValueForType == nil {
			uriObject[typeObj.Name] = []interface{}{}
		} else if typeObj.Type == "string" && uriObjValueForType == nil {
			uriObject[typeObj.Name] = ""
		} else if typeObj.Type == "uint64" && uriObjValueForType == nil {
			// b, err := parseInteger("uint64", 0)
			// if err != nil {
			// 	return typedData, err
			// }
			uriObject[typeObj.Name] = "0"
		}
	}

	if addUri {
		uriObject, ok := typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"].(map[string]interface{})["uri"].(map[string]interface{})//TODO: multuple messages
		if !ok {
			return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "message is not a map[string]interface{}")
		}
		uriObjValueForType := uriObject["idxRangeToRemove"].(map[string]interface{})

		if uriObjValueForType == nil {
			uriObject["idxRangeToRemove"] = map[string]interface{}{
				"start": "0",
				"end": "0",
			}
		} 
		if uriObjValueForType["start"] == nil {
			uriObjValueForType["start"] = "0"
		} 
		if uriObjValueForType["end"] == nil {
			uriObjValueForType["end"] = "0"
		}
	}

	if addWhitelistInfo {
		whitelistIdRangesArr, ok := typedData.Message["msgs"].([]interface{})[0].(map[string]interface{})["value"].(map[string]interface{})["whitelistedRecipients"] //TODO: multuple messages
		if !ok {
			return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "1. message is not a map[string]interface{}")
		}
		whitelistValuesArr := whitelistIdRangesArr.([]interface{})
		if !ok {
			return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "2. message is not a []interface{}")
		}

		for _, whitelistValue := range whitelistValuesArr {
			balanceAmounts := whitelistValue.(map[string]interface{})["balanceAmounts"].([]interface{})
			if !ok {
				return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "2. message is not a []interface{}")
			}
			for _, balanceAmount := range balanceAmounts {
				idRanges, ok := balanceAmount.(map[string]interface{})["idRanges"].([]interface{})
				if !ok {
					return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "3. message is not a map[string]interface{}")
				}
				for _, idRangeObj := range idRanges {
					idRange, ok := idRangeObj.(map[string]interface{})
					if !ok {
						return typedData, sdkerrors.Wrap(ErrInvalidTypedData, "4. message is not a map[string]interface{}")
					}
					if idRange["start"] == nil {
						idRange["start"] = "0"
					} 
					if idRange["end"] == nil {
						idRange["end"] = "0"
					}
				}
			}
		}
	}

	return typedData, nil
}

//first bool is if URI is needed, second is if ID Range is needed
func GetMsgValueTypes(route string) ([]apitypes.Type, bool, bool, bool) {
	switch route {
	case TypeMsgNewBadge:
		return []apitypes.Type{ 
			{ Name: "creator", Type: "string" },
			{ Name: "uri", Type: "UriObject" },
			{ Name: "arbitraryBytes", Type: "string" },
			{ Name: "permissions", Type: "uint64" },
			{ Name: "defaultSubassetSupply", Type: "uint64" },
			{ Name: "freezeAddressRanges", Type: "IdRange[]" },
			{ Name: "standard", Type: "uint64" },
			{ Name: "subassetSupplys", Type: "uint64[]" },
			{ Name: "subassetAmountsToCreate", Type: "uint64[]" },
			{ Name: "whitelistedRecipients", Type: "WhitelistMintInfo[]" },
			
		}, true, true, true
	case TypeMsgNewSubBadge:
		return []apitypes.Type{ 
			{ Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "supplys", Type: "uint64[]" },
			{ Name: "amountsToCreate", Type: "uint64[]" },
		}, false, false, false
	case TypeMsgTransferBadge:
		return []apitypes.Type{	  { Name: "creator", Type: "string" },
			{ Name: "from", Type: "uint64" },
			{ Name: "toAddresses", Type: "uint64[]" },
			{ Name: "amounts", Type: "uint64[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
			{ Name: "expirationTime", Type: "uint64" },
			{ Name: "cantCancelBeforeTime", Type: "uint64" },
		}, false, true, false
	case TypeMsgRequestTransferBadge:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "from", Type: "uint64" },
			{ Name: "amount", Type: "uint64" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
			{ Name: "expirationTime", Type: "uint64" },
			{ Name: "cantCancelBeforeTime", Type: "uint64" },
		}, false, true, false
	case TypeMsgHandlePendingTransfer:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "accept", Type: "bool" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "nonceRanges", Type: "IdRange[]" },
			{ Name: "forcefulAccept", Type: "bool" },
		}, false, true, false
	case TypeMsgSetApproval:
		return []apitypes.Type{
			{ Name: "creator", Type: "string" },
			{ Name: "amount", Type: "uint64" },
			{ Name: "address", Type: "uint64" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
		}, false, true, false
	case TypeMsgRevokeBadge:
		return []apitypes.Type{  { Name: "creator", Type: "string" },
			{ Name: "addresses", Type: "uint64[]" },
			{ Name: "amounts", Type: "uint64[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
		}, false, true, false
	case TypeMsgFreezeAddress:
		return []apitypes.Type{
			{ Name: "creator", Type: "string" },
			{ Name: "addressRanges", Type: "IdRange[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "add", Type: "bool" },
		}, false, true, false
	case TypeMsgUpdateUris:
		return []apitypes.Type{	 { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "uri", Type: "UriObject" },
		}, true, true, false
	case TypeMsgUpdatePermissions:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "permissions", Type: "uint64" },
		}, false, false, false
	case TypeMsgUpdateBytes:
		return []apitypes.Type{  { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "newBytes", Type: "string" },
		}, false, false, false
	case TypeMsgTransferManager:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "address", Type: "uint64" },
		}, false, false, false
	case TypeMsgRequestTransferManager:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "add", Type: "bool" },
		}, false, false, false
	case TypeMsgSelfDestructBadge:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
		}, false, false, false
	case TypeMsgPruneBalances:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeIds", Type: "uint64[]" },
			{ Name: "addresses", Type: "uint64[]" },
		}, false, false, false
	case TypeMsgRegisterAddresses:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "addressesToRegister", Type: "string[]" },
		}, false, false, false
	default:
		return []apitypes.Type{}, false, false, false
	};
}
