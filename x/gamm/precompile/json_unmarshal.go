package gamm

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

// convertEVMAddressToBech32 converts an EVM address (0x...) to bech32 format if needed
// If the address is already in bech32 format, it returns it unchanged
// If it's an EVM address, it converts it to bech32
// Note: Invalid addresses are passed through for validation by ValidateBasic() which
// provides more comprehensive error messages. This function only handles format conversion.
func convertEVMAddressToBech32(addr string) string {
	if addr == "" {
		return addr // Empty is allowed (optional field)
	}
	// Check if it's already a bech32 address
	if _, err := sdk.AccAddressFromBech32(addr); err == nil {
		// Already valid bech32, return as-is
		return addr
	}
	// Check if it's an EVM address (0x...)
	if common.IsHexAddress(addr) {
		evmAddr := common.HexToAddress(addr)
		return sdk.AccAddress(evmAddr.Bytes()).String()
	}
	// Invalid format - return as-is for ValidateBasic() to catch with proper error context
	return addr
}

func meterGammJSONInput(ctx sdk.Context, input string) error {
	opts := evmcompat.DefaultJSONBudgetOptions()
	opts.MaxStringLength = MaxStringLength
	opts.GasPerValue = 100
	return evmcompat.ValidateJSONInput(ctx, input, opts)
}

// unmarshalMsgFromJSON unmarshals a JSON string into the appropriate Msg type based on method name
// and sets the Sender field from the contract caller.
func (p Precompile) unmarshalMsgFromJSON(ctx sdk.Context, methodName string, jsonStr string, contract *vm.Contract) (sdk.Msg, error) {
	// Get caller address
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	senderCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	if err := meterGammJSONInput(ctx, jsonStr); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid JSON input: %s", err))
	}

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
	if err := json.Unmarshal([]byte(jsonStr), msg); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON into %T: %s", msg, err))
	}

	// Set Sender field from contract caller (security: override any value in JSON)
	// Also convert affiliate addresses from EVM to bech32 format if needed
	switch m := msg.(type) {
	case *gammtypes.MsgJoinPool:
		m.Sender = senderCosmosAddr
	case *gammtypes.MsgExitPool:
		m.Sender = senderCosmosAddr
	case *gammtypes.MsgSwapExactAmountIn:
		m.Sender = senderCosmosAddr
		// Convert affiliate addresses from EVM to bech32 format
		for i := range m.Affiliates {
			if m.Affiliates[i].Address != "" {
				m.Affiliates[i].Address = convertEVMAddressToBech32(m.Affiliates[i].Address)
			}
		}
	case *gammtypes.MsgSwapExactAmountInWithIBCTransfer:
		m.Sender = senderCosmosAddr
		// Convert affiliate addresses from EVM to bech32 format
		// NOTE: Do NOT convert IbcTransferInfo.Receiver - it's an IBC receiver on another chain
		for i := range m.Affiliates {
			if m.Affiliates[i].Address != "" {
				m.Affiliates[i].Address = convertEVMAddressToBech32(m.Affiliates[i].Address)
			}
		}
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
func (p Precompile) unmarshalQueryFromJSON(ctx sdk.Context, methodName string, jsonStr string) (interface{}, error) {
	if err := meterGammJSONInput(ctx, jsonStr); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid query JSON input: %s", err))
	}

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
	if err := json.Unmarshal([]byte(jsonStr), queryReq); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal query JSON into %T: %s", queryReq, err))
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

