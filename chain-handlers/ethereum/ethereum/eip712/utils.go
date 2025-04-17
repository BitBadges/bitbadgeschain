package eip712

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/tidwall/gjson"
)

// NormalizeEmptyTypes is a recursive function that adds empty values to fields that are omitted when serialized with Proto and Amino.
// EIP712 doesn't support optional fields, so we add the omitted empty values back in here.
// This includes empty strings, when a uint64 is 0 (for cosmos.Uint we do empty string since it is a custom type ""), and when a bool is false.
func NormalizeEmptyTypes(typedData apitypes.TypedData, typeObjArr []apitypes.Type, mapObject map[string]interface{}) (map[string]interface{}, error) {
	for _, typeObj := range typeObjArr {
		typeStr := typeObj.Type
		value := mapObject[typeObj.Name]

		if strings.Contains(typeStr, "[]") && value == nil {
			mapObject[typeObj.Name] = []interface{}{}
		} else if strings.Contains(typeStr, "[]") && value != nil {
			valueArr := value.([]interface{})
			// Get typeStr without the [] brackets at the end
			typeStr = typeStr[:len(typeStr)-2]

			//For multi-type arrays, we add Any[] to the end of the typeStr
			//And the individual element types are going to be the stripped typeStr + "0" or "1" etc
			isMultiTypeArray := strings.Contains(typeStr, "Any[]")
			if isMultiTypeArray {
				typeStr = typeStr[:len(typeStr)-3]
			}

			for i, value := range valueArr {

				elementTypeStr := typeStr
				if isMultiTypeArray {
					elementTypeStr += strconv.Itoa(i)
				}

				innerMap, ok := value.(map[string]interface{})
				if ok {
					newMap, err := NormalizeEmptyTypes(typedData, typedData.Types[elementTypeStr], innerMap)
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

	// Check if this is a MsgExec authz message
	if strings.Contains(msgType, "cosmos-sdk/MsgExec") && msg.Get("value.msgs").Exists() {
		// Get the inner messages array
		innerMsgs := msg.Get("value.msgs").Array()

		// Create a new array to hold populated messages
		populatedMsgs := make([]interface{}, len(innerMsgs))

		// Process each inner message
		for i, innerMsg := range innerMsgs {
			populatedInnerMsg, err := GetPopulatedSchemaForMsg(innerMsg)
			if err != nil {
				return msg, err
			}
			// Store the populated message
			populatedMsgs[i] = populatedInnerMsg.Value()
		}

		// Construct new message with populated msgs
		newMsg := map[string]interface{}{
			"type": msgType,
			"value": map[string]interface{}{
				"grantee": msg.Get("value.grantee").String(),
				"msgs":    populatedMsgs,
			},
		}

		jsonBytes, err := json.Marshal(newMsg)
		if err != nil {
			return msg, err
		}

		return gjson.ParseBytes(jsonBytes), nil
	}

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
