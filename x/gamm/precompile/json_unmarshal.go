package gamm

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

// unmarshalMsgFromJSON unmarshals a JSON string into the appropriate Msg type based on method name
// and sets the Sender field from the contract caller for security.
func (p Precompile) unmarshalMsgFromJSON(methodName string, jsonStr string, contract *vm.Contract) (sdk.Msg, error) {
	// Get caller address
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	senderCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Create the appropriate Msg type based on method name
	var msg sdk.Msg
	switch methodName {
	case JoinPoolMethod:
		msg = &gammtypes.MsgJoinPool{}
	case ExitPoolMethod:
		msg = &gammtypes.MsgExitPool{}
	case SwapExactAmountInMethod:
		msg = &gammtypes.MsgSwapExactAmountIn{}
	case SwapExactAmountInWithIBCTransferMethod:
		msg = &gammtypes.MsgSwapExactAmountInWithIBCTransfer{}
	case CreatePoolMethod:
		msg = &balancer.MsgCreateBalancerPool{}
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown method: %s", methodName))
	}

	// Unmarshal JSON into the message
	// Use standard json.Unmarshal since gamm types use standard protobuf JSON
	if err := json.Unmarshal([]byte(jsonStr), msg); err != nil {
		// Validate JSON syntax separately for better error messages
		var jsonMap map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(jsonStr), &jsonMap); jsonErr != nil {
			// JSON syntax is invalid
			return nil, ErrInvalidInput(fmt.Sprintf("invalid JSON syntax: %s", jsonErr))
		}
		// JSON syntax is valid but protobuf unmarshaling failed
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON into %T: %s. JSON was: %s", msg, err, jsonStr))
	}

	// Set Sender field from contract caller (security: override any value in JSON)
	switch m := msg.(type) {
	case *gammtypes.MsgJoinPool:
		m.Sender = senderCosmosAddr
	case *gammtypes.MsgExitPool:
		m.Sender = senderCosmosAddr
	case *gammtypes.MsgSwapExactAmountIn:
		m.Sender = senderCosmosAddr
	case *gammtypes.MsgSwapExactAmountInWithIBCTransfer:
		m.Sender = senderCosmosAddr
	case *balancer.MsgCreateBalancerPool:
		m.Sender = senderCosmosAddr
	}

	// Validate message using ValidateBasic
	// Use panic recovery to handle cases where ValidateBasic might panic on nil/uninitialized fields
	if validator, ok := msg.(interface{ ValidateBasic() error }); ok {
		var validationErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error - this handles cases where ValidateBasic panics on nil fields
					// This is a safety measure for production readiness
					validationErr = fmt.Errorf("validation panic: %s", r)
				}
			}()
			validationErr = validator.ValidateBasic()
		}()
		if validationErr != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("message validation failed: %s", validationErr))
		}
	}

	return msg, nil
}

// unmarshalQueryFromJSON unmarshals a JSON string into the appropriate QueryRequest type
func (p Precompile) unmarshalQueryFromJSON(methodName string, jsonStr string) (interface{}, error) {
	var queryReq interface{}

	switch methodName {
	case GetPoolMethod:
		queryReq = &gammtypes.QueryPoolRequest{}
	case GetPoolsMethod:
		queryReq = &gammtypes.QueryPoolsRequest{}
	case GetPoolTypeMethod:
		queryReq = &gammtypes.QueryPoolTypeRequest{}
	case CalcJoinPoolNoSwapSharesMethod:
		queryReq = &gammtypes.QueryCalcJoinPoolNoSwapSharesRequest{}
	case CalcExitPoolCoinsFromSharesMethod:
		queryReq = &gammtypes.QueryCalcExitPoolCoinsFromSharesRequest{}
	case CalcJoinPoolSharesMethod:
		queryReq = &gammtypes.QueryCalcJoinPoolSharesRequest{}
	case GetPoolParamsMethod:
		queryReq = &gammtypes.QueryPoolParamsRequest{}
	case GetTotalSharesMethod:
		queryReq = &gammtypes.QueryTotalSharesRequest{}
	case GetTotalLiquidityMethod:
		queryReq = &gammtypes.QueryTotalLiquidityRequest{}
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown query method: %s", methodName))
	}

	// Unmarshal JSON into the query request
	// Use standard json.Unmarshal since gamm types use standard protobuf JSON
	if err := json.Unmarshal([]byte(jsonStr), queryReq); err != nil {
		// Validate JSON syntax separately for better error messages
		var jsonMap map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(jsonStr), &jsonMap); jsonErr != nil {
			// JSON syntax is invalid
			return nil, ErrInvalidInput(fmt.Sprintf("invalid JSON syntax: %s", jsonErr))
		}
		// JSON syntax is valid but protobuf unmarshaling failed
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal query JSON into %T: %s. JSON was: %s", queryReq, err, jsonStr))
	}

	// Validate query request using ValidateBasic if available
	// Use panic recovery to handle cases where ValidateBasic might panic on nil/uninitialized fields
	if validator, ok := queryReq.(interface{ ValidateBasic() error }); ok {
		var validationErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error - this handles cases where ValidateBasic panics on nil fields
					validationErr = fmt.Errorf("validation panic: %s", r)
				}
			}()
			validationErr = validator.ValidateBasic()
		}()
		if validationErr != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("query validation failed: %v", validationErr))
		}
	}

	return queryReq, nil
}

