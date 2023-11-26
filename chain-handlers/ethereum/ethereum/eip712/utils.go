package eip712

import (
	"strings"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/tidwall/gjson"
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
			mapObject[typeObj.Name] = "0" //TODO: Does this really work / resolve correctly? We don't use it in any x/badges txs but it should be tested
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

// Recursively iterate through message and populate empty fields with default falsy values
func NormalizeEIP712TypedData(typedData apitypes.TypedData) (apitypes.TypedData, error) {
	msgValue := typedData.Message
	normalizedMsgValue, err := NormalizeEmptyTypes(typedData, typedData.Types["Tx"], msgValue)
	if err != nil {
		return typedData, err
	}

	typedData.Message = normalizedMsgValue

	return typedData, nil
}

// GetPopulatedSchemaForMsg returns a sample JSON object for a given message type with all potentially optional fields populated
func GetPopulatedSchemaForMsg(msg gjson.Result) (gjson.Result, error) {
	//1. Iterate through the ./example-jsons directory and find the file that matches the msg type
	//2. If a match is found, return the populated JSON object
	//3. If no match is found, return the original msg

	msgType := msg.Get("type").String()
	schemas := GetSchemas()

	for _, schema := range schemas {
		fileContents := schema

		//Check if the schema corresponds to the msg type
		if strings.Contains(fileContents, msgType) {
			return gjson.Parse(fileContents), nil
		}
	}

	return msg, nil
}
