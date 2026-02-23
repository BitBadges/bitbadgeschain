package precompile

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// VerifyCaller verifies that the caller address is valid
func VerifyCaller(caller common.Address) error {
	if caller == (common.Address{}) {
		return ErrInvalidInput("caller address is zero")
	}
	return nil
}

// unmarshalMsgFromJSON unmarshals a JSON string into the appropriate Msg type based on method name
// and sets the from_address field from the contract caller for security.
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
	case SendMethod:
		msg = &sendmanagertypes.MsgSendWithAliasRouting{}
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown method: %s", methodName))
	}

	// Unmarshal JSON into the message
	// Use standard json.Unmarshal since sendmanager types use standard protobuf JSON
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

	// Set from_address field from contract caller (security: override any value in JSON)
	switch m := msg.(type) {
	case *sendmanagertypes.MsgSendWithAliasRouting:
		m.FromAddress = senderCosmosAddr
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
