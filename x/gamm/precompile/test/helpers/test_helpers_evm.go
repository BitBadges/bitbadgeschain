package helpers

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// CreateEVMAccount creates an EVM account with private key
// Exported for use in test packages
func CreateEVMAccount() (*ecdsa.PrivateKey, common.Address, sdk.AccAddress) {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	return key, addr, cosmosAddr
}

// FundEVMAccount funds an EVM account with native tokens for gas
// Exported for use in test packages
func FundEVMAccount(ctx sdk.Context, bankKeeper bankkeeper.Keeper, addr sdk.AccAddress, amount sdk.Coins) error {
	// Mint coins to the account
	err := bankKeeper.MintCoins(ctx, "mint", amount)
	if err != nil {
		return err
	}
	// Send coins from mint module to the account
	return bankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, amount)
}

// BuildEVMTransaction builds an EVM transaction
// Exported for use in test packages
// For contract creation, pass nil as the 'to' parameter
func BuildEVMTransaction(
	fromKey *ecdsa.PrivateKey,
	to *common.Address, // Use pointer - nil means contract creation
	data []byte,
	value *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	nonce uint64,
	chainID *big.Int,
) (*types.Transaction, error) {
	var tx *types.Transaction
	if to == nil {
		// Contract creation - use NewContractCreation which sets To to nil
		tx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, data)
	} else {
		// Regular transaction - call to existing address
		tx = types.NewTransaction(nonce, *to, value, gasLimit, gasPrice, data)
	}
	return types.SignTx(tx, types.NewEIP155Signer(chainID), fromKey)
}

// ExecuteEVMTransaction executes a transaction through EVM keeper
// Exported for use in test packages
// NOTE: This function includes error handling for known snapshot issues in cosmos/evm
// When precompiles return errors, the EVM tries to revert, but the snapshot stack
// can be empty, causing "snapshot index 0 out of bound [0..0)" panics.
// This is a workaround that catches and handles these errors gracefully.
func ExecuteEVMTransaction(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	tx *types.Transaction,
) (*evmtypes.MsgEthereumTxResponse, error) {
	// Convert transaction to MsgEthereumTx
	msg := &evmtypes.MsgEthereumTx{}
	msg.FromEthereumTx(tx)

	// Use recover to catch panics from snapshot errors
	// ROOT CAUSE: The upstream cosmos/evm module's RunNativeAction doesn't create a snapshot
	// before calling precompiles. When a precompile returns an error, the EVM tries to revert
	// to snapshot 0, but no snapshots exist, causing a panic.
	var response *evmtypes.MsgEthereumTxResponse
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if this is a snapshot error panic
				errStr := fmt.Sprintf("%v", r)
				if strings.Contains(errStr, "snapshot index") && strings.Contains(errStr, "out of bound") {
					// This is the known snapshot bug - the precompile likely executed but returned an error
					// and the EVM tried to revert with an empty snapshot stack
					// Create a response indicating the error occurred during revert
					response = &evmtypes.MsgEthereumTxResponse{
						VmError: "snapshot revert error: " + errStr,
						Ret:     []byte{},
						GasUsed: 0,
					}
					err = nil // Don't propagate the panic as an error
				} else {
					// Re-panic for non-snapshot errors
					panic(r)
				}
			}
		}()

		// Execute through keeper
		response, err = evmKeeper.EthereumTx(ctx, msg)
	}()

	// If we caught a snapshot error panic, return the response with VmError set
	if response != nil && strings.Contains(response.VmError, "snapshot revert error") {
		return response, nil
	}

	// Check for snapshot-related errors in the response (non-panic case)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "snapshot index") && strings.Contains(errStr, "out of bound") {
			// This is the known snapshot bug
			if response == nil {
				response = &evmtypes.MsgEthereumTxResponse{
					VmError: "snapshot revert error: " + errStr,
					Ret:     []byte{},
					GasUsed: 0,
				}
			} else {
				response.VmError = "snapshot revert error: " + errStr
			}
			return response, nil
		}
		// For other errors, return as-is
		return response, err
	}

	return response, nil
}

// GetEventsFromContext extracts events from a context's event manager
// This is useful for verifying events emitted by precompiles
func GetEventsFromContext(ctx sdk.Context) []sdk.Event {
	return ctx.EventManager().Events()
}

// FindEventByName finds an event by type name in a slice of events
func FindEventByName(events []sdk.Event, eventType string) *sdk.Event {
	for i := range events {
		if events[i].Type == eventType {
			return &events[i]
		}
	}
	return nil
}

// FindEventsByModule finds all events with a specific module attribute
func FindEventsByModule(events []sdk.Event, moduleName string) []sdk.Event {
	result := []sdk.Event{}
	for _, event := range events {
		for _, attr := range event.Attributes {
			if attr.Key == sdk.AttributeKeyModule && attr.Value == moduleName {
				result = append(result, event)
				break
			}
		}
	}
	return result
}

// VerifyEventAttributes verifies that an event has the expected attributes
func VerifyEventAttributes(event *sdk.Event, expectedAttrs map[string]string) bool {
	if event == nil {
		return false
	}

	attrMap := make(map[string]string)
	for _, attr := range event.Attributes {
		attrMap[attr.Key] = attr.Value
	}

	for key, expectedValue := range expectedAttrs {
		if actualValue, found := attrMap[key]; !found || actualValue != expectedValue {
			return false
		}
	}

	return true
}

