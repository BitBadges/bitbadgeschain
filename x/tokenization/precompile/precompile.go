// Package tokenization implements a precompiled contract for the BitBadges tokenization module.
// This precompile enables Solidity smart contracts to interact with the tokenization system
// through a standardized EVM interface.
//
// The precompile is available at address 0x0000000000000000000000000000000000001001 and provides
// both transaction methods (state-changing operations) and query methods (read-only operations).
//
// Transaction Methods:
//   - transferTokens: Transfer tokens from the caller to one or more recipients
//   - setIncomingApproval: Set an incoming approval for the caller
//   - setOutgoingApproval: Set an outgoing approval for the caller
//
// Query Methods:
//   - getCollection: Query collection data by ID
//   - getBalance: Query balance data for a user address
//   - getBalanceAmount: Query balance amount for specific token IDs and ownership times
//   - getTotalSupply: Query total supply for specific token IDs and ownership times
//   - And many more query methods for approvals, trackers, votes, etc.
//
// All methods use structured error handling with error codes for consistent error reporting.
// Input validation is performed on all parameters to ensure security and correctness.
package tokenization

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"

	"github.com/cosmos/gogoproto/proto"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

const (
	// Base gas costs for transactions
	// IMPORTANT: These values are DEDUCTED from the transaction gas before the precompile runs.
	// The actual execution gas comes from the remaining gas (contract.Gas after deduction).
	// Setting these too high causes "out of gas" errors because there's not enough remaining
	// gas for the Cosmos SDK operations. Keep these as minimal entry fees.
	GasTransferTokensBase            = 30_000
	GasSetIncomingApprovalBase       = 20_000
	GasSetOutgoingApprovalBase       = 20_000
	GasCreateCollectionBase          = 50_000
	GasUpdateCollectionBase          = 40_000
	GasDeleteCollectionBase          = 20_000
	GasCreateAddressListsBase        = 30_000
	GasUpdateUserApprovalsBase       = 30_000
	GasDeleteIncomingApprovalBase    = 15_000
	GasDeleteOutgoingApprovalBase    = 15_000
	GasPurgeApprovalsBase            = 25_000
	GasCreateDynamicStoreBase        = 20_000
	GasUpdateDynamicStoreBase        = 20_000
	GasDeleteDynamicStoreBase        = 15_000
	GasSetDynamicStoreValueBase      = 15_000
	GasSetValidTokenIdsBase          = 20_000
	GasSetManagerBase                = 15_000
	GasSetCollectionMetadataBase     = 15_000
	GasSetTokenMetadataBase          = 20_000
	GasSetCustomDataBase             = 15_000
	GasSetStandardsBase              = 15_000
	GasSetCollectionApprovalsBase    = 30_000
	GasSetIsArchivedBase             = 15_000
	GasCastVoteBase                  = 15_000
	GasUniversalUpdateCollectionBase = 50_000

	// Gas costs per element for dynamic calculations
	GasPerRecipient          = 5_000
	GasPerTokenIdRange       = 1_000
	GasPerOwnershipTimeRange = 1_000
	GasPerApprovalField      = 500

	// Gas costs for queries (lower since they're read-only)
	GasGetCollectionBase         = 3_000
	GasGetBalanceBase            = 3_000
	GasGetAddressList            = 5_000
	GasGetApprovalTracker        = 5_000
	GasGetChallengeTracker       = 5_000
	GasGetETHSignatureTracker    = 5_000
	GasGetDynamicStore           = 5_000
	GasGetDynamicStoreValue      = 5_000
	GasGetWrappableBalances      = 5_000
	GasIsAddressReservedProtocol = 2_000
	GasGetAllReservedProtocol    = 5_000
	GasGetVote                   = 5_000
	GasGetVotes                  = 5_000
	GasParams                    = 2_000
	GasGetBalanceAmountBase      = 3_000
	GasGetTotalSupplyBase        = 3_000
	GasPerQueryRange             = 500
	GasExecuteMultipleBase       = 10_000
	GasPerMessageInBatch         = 1_000
	MaxMessagesPerBatch          = 50   // Maximum number of messages allowed in executeMultiple batch
	MaxQueryArraySize            = 1000 // Maximum size for tokenIds and ownershipTimes arrays in queries
)

var _ vm.PrecompiledContract = &Precompile{}

var (
	// Embed abi json file to the executable binary. Needed when importing as dependency.
	//
	//go:embed abi.json
	f   embed.FS
	ABI abi.ABI
	// abiLoadError stores any error from ABI loading for lazy error reporting
	abiLoadError error
)

func init() {
	ABI, abiLoadError = cmn.LoadABI(f, "abi.json")
	if abiLoadError != nil {
		// Log the error but don't panic - the error will be returned when the precompile is used
		// This allows the chain to start even if the ABI is malformed, but the precompile will be disabled
		fmt.Printf("WARNING: Failed to load tokenization precompile ABI: %v\n", abiLoadError)
	}
}

// GetABILoadError returns any error that occurred during ABI loading
// This can be checked by callers to verify the precompile is properly initialized
func GetABILoadError() error {
	return abiLoadError
}

// Precompile defines the tokenization precompile
type Precompile struct {
	cmn.Precompile

	abi.ABI
	tokenizationKeeper tokenizationkeeper.Keeper
}

// NewPrecompile creates a new tokenization Precompile instance implementing the
// PrecompiledContract interface.
func NewPrecompile(
	tokenizationKeeper tokenizationkeeper.Keeper,
) *Precompile {
	return &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:          storetypes.GasConfig{},
			TransientKVGasConfig: storetypes.GasConfig{},
			ContractAddress:      common.HexToAddress(TokenizationPrecompileAddress),
		},
		ABI:                ABI,
		tokenizationKeeper: tokenizationKeeper,
	}
}

// TokenizationPrecompileAddress is the address of the tokenization precompile
// Using standard precompile address range: 0x0000000000000000000000000000000000001001
const TokenizationPrecompileAddress = "0x0000000000000000000000000000000000001001"

// GetCallerAddress gets the caller address and converts it to Cosmos format
// This should be used for ALL transaction methods to set the Creator field
// SECURITY: This ensures the creator is always the actual caller, preventing impersonation
// The caller is obtained from contract.Caller() which returns the EVM msg.sender
// and cannot be spoofed by malicious contracts
func (p Precompile) GetCallerAddress(contract *vm.Contract) (string, error) {
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return "", err
	}
	return sdk.AccAddress(caller.Bytes()).String(), nil
}

// RequiredGas calculates the precompiled contract's base gas rate.
// Returns a conservative estimate that accounts for Cosmos SDK operations to help
// estimateGas converge on a working value.
func (p Precompile) RequiredGas(input []byte) uint64 {
	// NOTE: This check avoid panicking when trying to decode the method ID
	if len(input) < 4 {
		return 0
	}

	methodID := input[:4]

	method, err := p.MethodById(methodID)
	if err != nil {
		// This should never happen since this method is going to fail during Run
		return 0
	}

	// Get base gas for the method
	var baseGas uint64
	var isTransaction bool

	// For methods that require dynamic gas calculation, we return a base amount
	// The actual gas will be calculated in Execute based on input size
	switch method.Name {
	// Transaction methods
	case TransferTokensMethod:
		baseGas = GasTransferTokensBase
		isTransaction = true
	case SetIncomingApprovalMethod:
		baseGas = GasSetIncomingApprovalBase
		isTransaction = true
	case SetOutgoingApprovalMethod:
		baseGas = GasSetOutgoingApprovalBase
		isTransaction = true
	case CreateCollectionMethod:
		baseGas = GasCreateCollectionBase
		isTransaction = true
	case UpdateCollectionMethod:
		baseGas = GasUpdateCollectionBase
		isTransaction = true
	case DeleteCollectionMethod:
		baseGas = GasDeleteCollectionBase
		isTransaction = true
	case CreateAddressListsMethod:
		baseGas = GasCreateAddressListsBase
		isTransaction = true
	case UpdateUserApprovalsMethod:
		baseGas = GasUpdateUserApprovalsBase
		isTransaction = true
	case DeleteIncomingApprovalMethod:
		baseGas = GasDeleteIncomingApprovalBase
		isTransaction = true
	case DeleteOutgoingApprovalMethod:
		baseGas = GasDeleteOutgoingApprovalBase
		isTransaction = true
	case PurgeApprovalsMethod:
		baseGas = GasPurgeApprovalsBase
		isTransaction = true
	case CreateDynamicStoreMethod:
		baseGas = GasCreateDynamicStoreBase
		isTransaction = true
	case UpdateDynamicStoreMethod:
		baseGas = GasUpdateDynamicStoreBase
		isTransaction = true
	case DeleteDynamicStoreMethod:
		baseGas = GasDeleteDynamicStoreBase
		isTransaction = true
	case SetDynamicStoreValueMethod:
		baseGas = GasSetDynamicStoreValueBase
		isTransaction = true
	case SetValidTokenIdsMethod:
		baseGas = GasSetValidTokenIdsBase
		isTransaction = true
	case SetManagerMethod:
		baseGas = GasSetManagerBase
		isTransaction = true
	case SetCollectionMetadataMethod:
		baseGas = GasSetCollectionMetadataBase
		isTransaction = true
	case SetTokenMetadataMethod:
		baseGas = GasSetTokenMetadataBase
		isTransaction = true
	case SetCustomDataMethod:
		baseGas = GasSetCustomDataBase
		isTransaction = true
	case SetStandardsMethod:
		baseGas = GasSetStandardsBase
		isTransaction = true
	case SetCollectionApprovalsMethod:
		baseGas = GasSetCollectionApprovalsBase
		isTransaction = true
	case SetIsArchivedMethod:
		baseGas = GasSetIsArchivedBase
		isTransaction = true
	case CastVoteMethod:
		baseGas = GasCastVoteBase
		isTransaction = true
	case UniversalUpdateCollectionMethod:
		baseGas = GasUniversalUpdateCollectionBase
		isTransaction = true
	case ExecuteMultipleMethod:
		// For executeMultiple, we need to parse the input to count messages for dynamic gas calculation
		// The input format is: methodID (4 bytes) + ABI-encoded tuple array
		if len(input) >= 4 {
			// Try to unpack the messages array to count messages
			// This is a best-effort calculation - if parsing fails, we use base gas
			method, err := p.MethodById(input[:4])
			if err == nil && method.Name == ExecuteMultipleMethod {
				unpacked, err := method.Inputs.Unpack(input[4:])
				if err == nil && len(unpacked) == 1 {
					// Count messages in the array
					msgCount := uint64(0)
					if msgs, ok := unpacked[0].([]interface{}); ok {
						msgCount = uint64(len(msgs))
					} else {
						// Try reflection for struct slice
						val := reflect.ValueOf(unpacked[0])
						if val.Kind() == reflect.Slice {
							msgCount = uint64(val.Len())
						}
					}
					// Enforce max batch size in gas calculation
					if msgCount > MaxMessagesPerBatch {
						msgCount = MaxMessagesPerBatch
					}
					// Calculate dynamic gas: base + (message count * per-message gas)
					baseGas = GasExecuteMultipleBase + (msgCount * GasPerMessageInBatch)
				} else {
					// If parsing fails, use base gas (will be adjusted during execution)
					baseGas = GasExecuteMultipleBase
				}
			} else {
				baseGas = GasExecuteMultipleBase
			}
		} else {
			baseGas = GasExecuteMultipleBase
		}
		isTransaction = true
	// Query methods
	case GetCollectionMethod:
		baseGas = GasGetCollectionBase
	case GetBalanceMethod:
		baseGas = GasGetBalanceBase
	case GetAddressListMethod:
		baseGas = GasGetAddressList
	case GetApprovalTrackerMethod:
		baseGas = GasGetApprovalTracker
	case GetChallengeTrackerMethod:
		baseGas = GasGetChallengeTracker
	case GetETHSignatureTrackerMethod:
		baseGas = GasGetETHSignatureTracker
	case GetDynamicStoreMethod:
		baseGas = GasGetDynamicStore
	case GetDynamicStoreValueMethod:
		baseGas = GasGetDynamicStoreValue
	case GetWrappableBalancesMethod:
		baseGas = GasGetWrappableBalances
	case IsAddressReservedProtocolMethod:
		baseGas = GasIsAddressReservedProtocol
	case GetAllReservedProtocolAddressesMethod:
		baseGas = GasGetAllReservedProtocol
	case GetVoteMethod:
		baseGas = GasGetVote
	case GetVotesMethod:
		baseGas = GasGetVotes
	case ParamsMethod:
		baseGas = GasParams
	case GetBalanceAmountMethod:
		baseGas = GasGetBalanceAmountBase
	case GetTotalSupplyMethod:
		baseGas = GasGetTotalSupplyBase
	default:
		return 0
	}

	// Add buffer for Cosmos SDK operations to help estimateGas converge
	// Transactions need more buffer due to state writes, bank transfers, etc.
	if isTransaction {
		return baseGas + 200_000
	}

	// Queries need less buffer but still need some for state reads
	return baseGas + 50_000
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	// Check if ABI loaded successfully during init
	if abiLoadError != nil {
		return nil, fmt.Errorf("tokenization precompile unavailable: ABI failed to load: %w", abiLoadError)
	}

	// Verify contract.Input is set
	if len(contract.Input) == 0 {
		return nil, fmt.Errorf("contract.Input is empty - precompile cannot execute without input data")
	}

	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		// Add panic recovery to catch any unexpected panics
		defer func() {
			if r := recover(); r != nil {
				// Re-panic so it's properly handled by the EVM
				panic(r)
			}
		}()

		result, methodName, err := p.ExecuteWithMethodName(ctx, contract, readonly)

		// Gas is tracked by the EVM, we log the method for monitoring
		LogPrecompileUsage(ctx, methodName, err == nil, 0, err)

		return result, err
	})
}

// Execute executes the precompiled contract tokenization methods defined in the ABI.
// Deprecated: Use ExecuteWithMethodName instead for better performance (avoids double method lookup).
func (p Precompile) Execute(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, error) {
	bz, _, err := p.ExecuteWithMethodName(ctx, contract, readOnly)
	return bz, err
}

// ExecuteWithMethodName executes the precompiled contract and returns the method name for logging.
// This avoids the double MethodById() lookup that occurs when logging separately.
func (p Precompile) ExecuteWithMethodName(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, string, error) {
	method, args, err := cmn.SetupABI(p.ABI, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, "unknown", fmt.Errorf("SetupABI failed: %w", err)
	}

	// Handle executeMultiple specially as it has different argument structure
	if method.Name == ExecuteMultipleMethod {
		// Manually unpack to handle tuple arrays correctly
		// SetupABI may not handle tuple arrays the way we need
		if len(contract.Input) < 4 {
			return nil, method.Name, ErrInvalidInput("contract input too short")
		}

		// Unpack arguments using method's input definition (skip 4-byte method selector)
		unpacked, err := method.Inputs.Unpack(contract.Input[4:])
		if err != nil {
			return nil, method.Name, WrapError(err, ErrorCodeInvalidInput, "failed to unpack method arguments")
		}

		if len(unpacked) != 1 {
			return nil, method.Name, ErrInvalidInput(fmt.Sprintf("expected 1 unpacked argument, got %d", len(unpacked)))
		}

		// ABI unpacks tuple arrays - handle different formats
		var messages []interface{}

		// Case 1: []interface{} (most common)
		if msgs, ok := unpacked[0].([]interface{}); ok {
			messages = msgs
		} else {
			// Case 2: Struct slice with JSON tags (what go-ethereum ABI returns for tuples)
			// Use reflection to handle the struct slice
			val := reflect.ValueOf(unpacked[0])
			if val.Kind() == reflect.Slice {
				messages = make([]interface{}, val.Len())
				for i := 0; i < val.Len(); i++ {
					elem := val.Index(i)
					if elem.Kind() == reflect.Struct {
						// Extract MessageType and MsgJson fields
						msgTypeField := elem.FieldByName("MessageType")
						msgJsonField := elem.FieldByName("MsgJson")

						if msgTypeField.IsValid() && msgJsonField.IsValid() {
							// Convert to []interface{} format [messageType, msgJson]
							messages[i] = []interface{}{
								msgTypeField.String(),
								msgJsonField.String(),
							}
						} else {
							return nil, method.Name, ErrInvalidInput(fmt.Sprintf(
								"struct at index %d missing MessageType or MsgJson fields", i))
						}
					} else {
						// If not a struct, try to use as-is
						messages[i] = elem.Interface()
					}
				}
			} else {
				return nil, method.Name, ErrInvalidInput(fmt.Sprintf(
					"expected messages array (slice), got type %T. "+
						"Value: %+v. "+
						"This indicates the ABI tuple array format is not as expected.",
					unpacked[0], unpacked[0]))
			}
		}

		bz, err := p.HandleExecuteMultiple(ctx, method, messages, contract)
		return bz, method.Name, err
	}

	// Extract JSON string from args for other methods
	if len(args) != 1 {
		return nil, method.Name, ErrInvalidInput(fmt.Sprintf("expected 1 argument (JSON string), got %d", len(args)))
	}

	jsonStr, ok := args[0].(string)
	if !ok {
		return nil, method.Name, ErrInvalidInput("expected JSON string as first argument")
	}

	// Route to transaction or query handler
	// Note: GetBalanceAmount, GetTotalSupply, and GetAllReservedProtocolAddresses need special handling
	var bz []byte
	if p.IsTransaction(method) {
		bz, err = p.HandleTransaction(ctx, method, jsonStr, contract)
	} else if method.Name == GetBalanceAmountMethod {
		bz, err = p.HandleGetBalanceAmount(ctx, method, jsonStr)
	} else if method.Name == GetTotalSupplyMethod {
		bz, err = p.HandleGetTotalSupply(ctx, method, jsonStr)
	} else if method.Name == GetAllReservedProtocolAddressesMethod {
		bz, err = p.HandleGetAllReservedProtocolAddresses(ctx, method, jsonStr)
	} else if method.Name == IsAddressReservedProtocolMethod {
		bz, err = p.HandleIsAddressReservedProtocol(ctx, method, jsonStr)
	} else {
		bz, err = p.HandleQuery(ctx, method, jsonStr)
	}

	return bz, method.Name, err
}

// HandleTransaction handles a transaction by unmarshaling JSON and executing via keeper
func (p Precompile) HandleTransaction(ctx sdk.Context, method *abi.Method, jsonStr string, contract *vm.Contract) ([]byte, error) {
	// Unmarshal JSON to Msg
	msg, err := p.unmarshalMsgFromJSON(method.Name, jsonStr, contract)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON for method %s: %w", method.Name, err)
	}

	// Execute message via keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)

	// Route to appropriate handler based on message type
	// Extract return values according to ABI outputs
	var result interface{}
	switch m := msg.(type) {
	case *tokenizationtypes.MsgTransferTokens:
		_, err = msgServer.TransferTokens(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgSetIncomingApproval:
		_, err = msgServer.SetIncomingApproval(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgSetOutgoingApproval:
		_, err = msgServer.SetOutgoingApproval(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgCreateCollection:
		resp, createErr := msgServer.CreateCollection(ctx, m)
		err = createErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "CreateCollection returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgUpdateCollection:
		resp, updateErr := msgServer.UpdateCollection(ctx, m)
		err = updateErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "UpdateCollection returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgDeleteCollection:
		_, err = msgServer.DeleteCollection(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgCreateAddressLists:
		_, err = msgServer.CreateAddressLists(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgUpdateUserApprovals:
		_, err = msgServer.UpdateUserApprovals(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgDeleteIncomingApproval:
		_, err = msgServer.DeleteIncomingApproval(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgDeleteOutgoingApproval:
		_, err = msgServer.DeleteOutgoingApproval(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgPurgeApprovals:
		resp, purgeErr := msgServer.PurgeApprovals(ctx, m)
		err = purgeErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "PurgeApprovals returned nil response")
		} else {
			result = resp.NumPurged.BigInt() // ABI: uint256 numPurged
		}
	case *tokenizationtypes.MsgCreateDynamicStore:
		resp, createErr := msgServer.CreateDynamicStore(ctx, m)
		err = createErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "CreateDynamicStore returned nil response")
		} else {
			result = resp.StoreId.BigInt() // ABI: uint256 storeId
		}
	case *tokenizationtypes.MsgUpdateDynamicStore:
		_, err = msgServer.UpdateDynamicStore(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgDeleteDynamicStore:
		_, err = msgServer.DeleteDynamicStore(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgSetDynamicStoreValue:
		_, err = msgServer.SetDynamicStoreValue(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgSetValidTokenIds:
		resp, setErr := msgServer.SetValidTokenIds(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetValidTokenIds returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetManager:
		resp, setErr := msgServer.SetManager(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetManager returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetCollectionMetadata:
		resp, setErr := msgServer.SetCollectionMetadata(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetCollectionMetadata returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetTokenMetadata:
		resp, setErr := msgServer.SetTokenMetadata(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetTokenMetadata returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetCustomData:
		resp, setErr := msgServer.SetCustomData(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetCustomData returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetStandards:
		resp, setErr := msgServer.SetStandards(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetStandards returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetCollectionApprovals:
		resp, setErr := msgServer.SetCollectionApprovals(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetCollectionApprovals returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgSetIsArchived:
		resp, setErr := msgServer.SetIsArchived(ctx, m)
		err = setErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "SetIsArchived returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	case *tokenizationtypes.MsgCastVote:
		_, err = msgServer.CastVote(ctx, m)
		result = true // ABI: bool success
	case *tokenizationtypes.MsgUniversalUpdateCollection:
		resp, updateErr := msgServer.UniversalUpdateCollection(ctx, m)
		err = updateErr
		if err != nil {
			// Error will be handled below
		} else if resp == nil {
			return nil, WrapError(fmt.Errorf("response is nil"), ErrorCodeInternalError, "UniversalUpdateCollection returned nil response")
		} else {
			result = resp.CollectionId.BigInt() // ABI: uint256 collectionId
		}
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unsupported message type for method: %s", method.Name))
	}

	if err != nil {
		return nil, WrapError(err, ErrorCodeTransferFailed, "transaction failed")
	}

	// Pack response - methods return specific values (collectionId, storeId) or bool success
	// Check if method has outputs defined
	if len(method.Outputs) == 0 {
		return []byte{}, nil
	}

	packed, packErr := method.Outputs.Pack(result)
	if packErr != nil {
		return nil, WrapError(packErr, ErrorCodeInternalError, "failed to pack result")
	}
	if len(packed) == 0 {
		return nil, WrapError(fmt.Errorf("packing returned empty result"), ErrorCodeInternalError, "packing returned empty byte slice")
	}
	return packed, nil
}

// HandleExecuteMultiple handles multiple messages in a single atomic transaction
func (p Precompile) HandleExecuteMultiple(ctx sdk.Context, method *abi.Method, messages []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(messages) == 0 {
		return nil, ErrInvalidInput("messages array cannot be empty")
	}

	// Enforce maximum batch size to prevent gas exhaustion attacks
	if len(messages) > MaxMessagesPerBatch {
		return nil, ErrInvalidInput(fmt.Sprintf("messages array exceeds maximum batch size: %d > %d", len(messages), MaxMessagesPerBatch))
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	results := make([][]byte, 0, len(messages))

	// Process each message sequentially
	for i, msgInterface := range messages {
		var messageType, msgJson string

		// Parse message tuple - can be either []interface{} or map[string]interface{}
		if msgTuple, ok := msgInterface.([]interface{}); ok && len(msgTuple) == 2 {
			// Tuple as slice: [messageType, msgJson]
			var ok1, ok2 bool
			messageType, ok1 = msgTuple[0].(string)
			msgJson, ok2 = msgTuple[1].(string)
			if !ok1 || !ok2 {
				return nil, WrapErrorWithContext(
					fmt.Errorf("invalid message tuple types at index %d", i),
					ErrorCodeInvalidInput,
					"messageType and msgJson must be strings",
					fmt.Sprintf("message index %d", i),
				)
			}
		} else if msgMap, ok := msgInterface.(map[string]interface{}); ok {
			// Tuple as map: {"messageType": "...", "msgJson": "..."}
			var ok1, ok2 bool
			messageType, ok1 = msgMap["messageType"].(string)
			msgJson, ok2 = msgMap["msgJson"].(string)
			if !ok1 || !ok2 {
				return nil, WrapErrorWithContext(
					fmt.Errorf("invalid message map format at index %d", i),
					ErrorCodeInvalidInput,
					"message map must contain 'messageType' and 'msgJson' string fields",
					fmt.Sprintf("message index %d", i),
				)
			}
		} else {
			return nil, WrapErrorWithContext(
				fmt.Errorf("invalid message format at index %d: expected tuple or map", i),
				ErrorCodeInvalidInput,
				"each message must be a tuple [messageType, msgJson] or map {messageType, msgJson}",
				fmt.Sprintf("message index %d", i),
			)
		}

		// Route and unmarshal message
		msg, err := p.routeMessageByType(messageType, msgJson, contract)
		if err != nil {
			return nil, WrapErrorWithContext(
				err,
				ErrorCodeInvalidInput,
				"failed to route message",
				fmt.Sprintf("message index %d, type: %s", i, messageType),
			)
		}

		// Execute message via keeper
		var result interface{}
		switch m := msg.(type) {
		case *tokenizationtypes.MsgTransferTokens:
			_, err = msgServer.TransferTokens(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgSetIncomingApproval:
			_, err = msgServer.SetIncomingApproval(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgSetOutgoingApproval:
			_, err = msgServer.SetOutgoingApproval(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgCreateCollection:
			resp, createErr := msgServer.CreateCollection(ctx, m)
			err = createErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"CreateCollection returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgUpdateCollection:
			resp, updateErr := msgServer.UpdateCollection(ctx, m)
			err = updateErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"UpdateCollection returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgDeleteCollection:
			_, err = msgServer.DeleteCollection(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgCreateAddressLists:
			_, err = msgServer.CreateAddressLists(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgUpdateUserApprovals:
			_, err = msgServer.UpdateUserApprovals(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgDeleteIncomingApproval:
			_, err = msgServer.DeleteIncomingApproval(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgDeleteOutgoingApproval:
			_, err = msgServer.DeleteOutgoingApproval(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgPurgeApprovals:
			resp, purgeErr := msgServer.PurgeApprovals(ctx, m)
			err = purgeErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"PurgeApprovals returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.NumPurged.BigInt()
			}
		case *tokenizationtypes.MsgCreateDynamicStore:
			resp, createErr := msgServer.CreateDynamicStore(ctx, m)
			err = createErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"CreateDynamicStore returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.StoreId.BigInt()
			}
		case *tokenizationtypes.MsgUpdateDynamicStore:
			_, err = msgServer.UpdateDynamicStore(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgDeleteDynamicStore:
			_, err = msgServer.DeleteDynamicStore(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgSetDynamicStoreValue:
			_, err = msgServer.SetDynamicStoreValue(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgSetValidTokenIds:
			resp, setErr := msgServer.SetValidTokenIds(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetValidTokenIds returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetManager:
			resp, setErr := msgServer.SetManager(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetManager returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetCollectionMetadata:
			resp, setErr := msgServer.SetCollectionMetadata(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetCollectionMetadata returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetTokenMetadata:
			resp, setErr := msgServer.SetTokenMetadata(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetTokenMetadata returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetCustomData:
			resp, setErr := msgServer.SetCustomData(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetCustomData returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetStandards:
			resp, setErr := msgServer.SetStandards(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetStandards returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetCollectionApprovals:
			resp, setErr := msgServer.SetCollectionApprovals(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetCollectionApprovals returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgSetIsArchived:
			resp, setErr := msgServer.SetIsArchived(ctx, m)
			err = setErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"SetIsArchived returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		case *tokenizationtypes.MsgCastVote:
			_, err = msgServer.CastVote(ctx, m)
			if err == nil {
				result = true
			}
		case *tokenizationtypes.MsgUniversalUpdateCollection:
			resp, updateErr := msgServer.UniversalUpdateCollection(ctx, m)
			err = updateErr
			if err == nil {
				if resp == nil {
					return nil, WrapErrorWithContext(
						fmt.Errorf("response is nil"),
						ErrorCodeInternalError,
						"UniversalUpdateCollection returned nil response",
						fmt.Sprintf("message index %d", i),
					)
				}
				result = resp.CollectionId.BigInt()
			}
		default:
			return nil, WrapErrorWithContext(
				fmt.Errorf("unsupported message type: %T", msg),
				ErrorCodeInvalidInput,
				"unsupported message type",
				fmt.Sprintf("message index %d, type: %s", i, messageType),
			)
		}

		// If execution failed, return error (atomic rollback handled by transaction context)
		if err != nil {
			return nil, WrapErrorWithContext(
				err,
				ErrorCodeTransferFailed,
				"message execution failed",
				fmt.Sprintf("message index %d, type: %s", i, messageType),
			)
		}

		// Pack result as bytes
		var resultBytes []byte
		if result != nil {
			// Determine result type and pack accordingly
			switch r := result.(type) {
			case bool:
				// Pack bool as ABI-encoded uint256 (32 bytes: 0x00...00 or 0x00...01)
				// ABI encoding requires booleans to be 32 bytes
				resultBytes = make([]byte, 32)
				if r {
					resultBytes[31] = 0x01
				}
				// else: all bytes remain 0x00
			case *big.Int:
				// Pack uint256 as bytes (32 bytes)
				resultBytes = r.Bytes()
				// Pad to 32 bytes if needed
				if len(resultBytes) < 32 {
					padded := make([]byte, 32)
					copy(padded[32-len(resultBytes):], resultBytes)
					resultBytes = padded
				} else if len(resultBytes) > 32 {
					// Truncate if too long (shouldn't happen for uint256)
					resultBytes = resultBytes[len(resultBytes)-32:]
				}
			default:
				// For other types, try to marshal as JSON then convert to bytes
				jsonBytes, jsonErr := json.Marshal(result)
				if jsonErr != nil {
					return nil, WrapErrorWithContext(
						jsonErr,
						ErrorCodeInternalError,
						"failed to marshal result",
						fmt.Sprintf("message index %d", i),
					)
				}
				resultBytes = jsonBytes
			}
		}

		results = append(results, resultBytes)
	}

	// Pack results: (bool success, bytes[] results)
	success := true
	packed, packErr := method.Outputs.Pack(success, results)
	if packErr != nil {
		return nil, WrapError(packErr, ErrorCodeInternalError, "failed to pack results")
	}

	return packed, nil
}

// HandleQuery handles a query by unmarshaling JSON and executing via keeper
func (p Precompile) HandleQuery(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Unmarshal JSON to QueryRequest
	queryReq, err := p.unmarshalQueryFromJSON(method.Name, jsonStr)
	if err != nil {
		return nil, err
	}

	// Validate query request before processing
	if err := p.validateQueryRequest(queryReq); err != nil {
		return nil, err
	}

	// Execute query via keeper
	var resp interface{}
	switch req := queryReq.(type) {
	case *tokenizationtypes.QueryGetCollectionRequest:
		resp, err = p.tokenizationKeeper.GetCollection(ctx, req)
	case *tokenizationtypes.QueryGetBalanceRequest:
		resp, err = p.tokenizationKeeper.GetBalance(ctx, req)
	case *tokenizationtypes.QueryGetAddressListRequest:
		resp, err = p.tokenizationKeeper.GetAddressList(ctx, req)
	case *tokenizationtypes.QueryGetApprovalTrackerRequest:
		resp, err = p.tokenizationKeeper.GetApprovalTracker(ctx, req)
	case *tokenizationtypes.QueryGetChallengeTrackerRequest:
		resp, err = p.tokenizationKeeper.GetChallengeTracker(ctx, req)
	case *tokenizationtypes.QueryGetETHSignatureTrackerRequest:
		resp, err = p.tokenizationKeeper.GetETHSignatureTracker(ctx, req)
	case *tokenizationtypes.QueryGetDynamicStoreRequest:
		resp, err = p.tokenizationKeeper.GetDynamicStore(ctx, req)
	case *tokenizationtypes.QueryGetDynamicStoreValueRequest:
		resp, err = p.tokenizationKeeper.GetDynamicStoreValue(ctx, req)
	case *tokenizationtypes.QueryGetWrappableBalancesRequest:
		resp, err = p.tokenizationKeeper.GetWrappableBalances(ctx, req)
	case *tokenizationtypes.QueryIsAddressReservedProtocolRequest:
		resp, err = p.tokenizationKeeper.IsAddressReservedProtocol(ctx, req)
	case *tokenizationtypes.QueryGetAllReservedProtocolAddressesRequest:
		resp, err = p.tokenizationKeeper.GetAllReservedProtocolAddresses(ctx, req)
	case *tokenizationtypes.QueryGetVoteRequest:
		resp, err = p.tokenizationKeeper.GetVote(ctx, req)
	case *tokenizationtypes.QueryGetVotesRequest:
		resp, err = p.tokenizationKeeper.GetVotes(ctx, req)
	case *tokenizationtypes.QueryParamsRequest:
		resp, err = p.tokenizationKeeper.Params(ctx, req)
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unsupported query type for method: %s", method.Name))
	}

	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "query failed")
	}

	// Handle special query methods that return uint256 instead of bytes
	switch method.Name {
	case GetChallengeTrackerMethod:
		if challengeResp, ok := resp.(*tokenizationtypes.QueryGetChallengeTrackerResponse); ok {
			// Parse numUsed string to uint256
			numUsed, err := sdkmath.ParseUint(challengeResp.NumUsed)
			if err != nil {
				return nil, WrapError(err, ErrorCodeInternalError, "failed to parse numUsed")
			}
			return method.Outputs.Pack(numUsed.BigInt())
		}
		return nil, WrapError(fmt.Errorf("invalid response type for getChallengeTracker"), ErrorCodeInternalError, "expected QueryGetChallengeTrackerResponse")
	case GetETHSignatureTrackerMethod:
		if ethResp, ok := resp.(*tokenizationtypes.QueryGetETHSignatureTrackerResponse); ok {
			// Parse numUsed string to uint256
			numUsed, err := sdkmath.ParseUint(ethResp.NumUsed)
			if err != nil {
				return nil, WrapError(err, ErrorCodeInternalError, "failed to parse numUsed")
			}
			return method.Outputs.Pack(numUsed.BigInt())
		}
		return nil, WrapError(fmt.Errorf("invalid response type for getETHSignatureTracker"), ErrorCodeInternalError, "expected QueryGetETHSignatureTrackerResponse")
	case GetWrappableBalancesMethod:
		if wrappableResp, ok := resp.(*tokenizationtypes.QueryGetWrappableBalancesResponse); ok {
			// Convert Uint to uint256
			return method.Outputs.Pack(wrappableResp.Amount.BigInt())
		}
		return nil, WrapError(fmt.Errorf("invalid response type for getWrappableBalances"), ErrorCodeInternalError, "expected QueryGetWrappableBalancesResponse")
	}

	// Marshal response to bytes (protobuf) for all other queries
	// Type assert to proto.Message
	if protoMsg, ok := resp.(proto.Message); ok {
		bz, err := tokenizationtypes.ModuleCdc.Marshal(protoMsg)
		if err != nil {
			return nil, WrapError(err, ErrorCodeInternalError, "failed to marshal query response")
		}
		return method.Outputs.Pack(bz)
	}

	return nil, WrapError(fmt.Errorf("response is not a proto.Message"), ErrorCodeInternalError, "invalid response type")
}

// validateQueryRequest validates query request inputs
func (p Precompile) validateQueryRequest(queryReq interface{}) error {
	switch req := queryReq.(type) {
	case *tokenizationtypes.QueryGetCollectionRequest:
		// Validate collection ID
		if req.CollectionId == "" {
			return ErrInvalidInput("collection ID cannot be empty")
		}
		collectionId, err := sdkmath.ParseUint(req.CollectionId)
		if err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid collection ID: %v", err))
		}
		if collectionId.IsZero() {
			return ErrInvalidInput("collection ID cannot be zero")
		}
	case *tokenizationtypes.QueryGetBalanceRequest:
		// Validate collection ID
		if req.CollectionId == "" {
			return ErrInvalidInput("collection ID cannot be empty")
		}
		collectionId, err := sdkmath.ParseUint(req.CollectionId)
		if err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid collection ID: %v", err))
		}
		if collectionId.IsZero() {
			return ErrInvalidInput("collection ID cannot be zero")
		}
		// Validate address
		if req.Address == "" {
			return ErrInvalidInput("address cannot be empty")
		}
		// Check if address is zero (all zeros)
		accAddr, err := sdk.AccAddressFromBech32(req.Address)
		if err == nil {
			// Check if it's a zero address (all bytes are zero)
			if accAddr.Empty() {
				return ErrInvalidInput("address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("address cannot be zero address")
			}
		}
		// Also check if it's a hex zero address
		if req.Address == "0x0000000000000000000000000000000000000000" || req.Address == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("address cannot be zero address")
		}
	case *tokenizationtypes.QueryGetDynamicStoreRequest:
		// Validate store ID
		if req.StoreId == "" {
			return ErrInvalidInput("store ID cannot be empty")
		}
		storeId, err := sdkmath.ParseUint(req.StoreId)
		if err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid store ID: %v", err))
		}
		if storeId.IsZero() {
			return ErrInvalidInput("store ID cannot be zero")
		}
	case *tokenizationtypes.QueryGetDynamicStoreValueRequest:
		// Validate store ID
		if req.StoreId == "" {
			return ErrInvalidInput("store ID cannot be empty")
		}
		storeId, err := sdkmath.ParseUint(req.StoreId)
		if err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid store ID: %v", err))
		}
		if storeId.IsZero() {
			return ErrInvalidInput("store ID cannot be zero")
		}
		// Validate address
		if req.Address == "" {
			return ErrInvalidInput("address cannot be empty")
		}
		// Check if address is zero
		accAddr, err := sdk.AccAddressFromBech32(req.Address)
		if err == nil {
			if accAddr.Empty() {
				return ErrInvalidInput("address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("address cannot be zero address")
			}
		}
		if req.Address == "0x0000000000000000000000000000000000000000" || req.Address == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("address cannot be zero address")
		}
	case *tokenizationtypes.QueryGetApprovalTrackerRequest:
		// Validate collection ID
		if req.CollectionId == "" {
			return ErrInvalidInput("collection ID cannot be empty")
		}
		collectionId, err := sdkmath.ParseUint(req.CollectionId)
		if err != nil {
			return ErrInvalidInput(fmt.Sprintf("invalid collection ID: %v", err))
		}
		if collectionId.IsZero() {
			return ErrInvalidInput("collection ID cannot be zero")
		}
		// Validate addresses - empty string means zero address, which is invalid
		if req.ApproverAddress == "" {
			return ErrInvalidInput("approver address cannot be empty or zero address")
		}
		accAddr, err := sdk.AccAddressFromBech32(req.ApproverAddress)
		if err == nil {
			if accAddr.Empty() {
				return ErrInvalidInput("approver address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("approver address cannot be zero address")
			}
		}
		if req.ApproverAddress == "0x0000000000000000000000000000000000000000" || req.ApproverAddress == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("approver address cannot be zero address")
		}
		if req.ApprovedAddress == "" {
			return ErrInvalidInput("approved address cannot be empty or zero address")
		}
		accAddr2, err := sdk.AccAddressFromBech32(req.ApprovedAddress)
		if err == nil {
			if accAddr2.Empty() {
				return ErrInvalidInput("approved address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr2 {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("approved address cannot be zero address")
			}
		}
		if req.ApprovedAddress == "0x0000000000000000000000000000000000000000" || req.ApprovedAddress == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("approved address cannot be zero address")
		}
	case *tokenizationtypes.QueryGetVoteRequest:
		// Validate addresses - both are required and cannot be zero
		if req.VoterAddress == "" {
			return ErrInvalidInput("voter address cannot be empty or zero address")
		}
		accAddr, err := sdk.AccAddressFromBech32(req.VoterAddress)
		if err == nil {
			if accAddr.Empty() {
				return ErrInvalidInput("voter address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("voter address cannot be zero address")
			}
		}
		if req.VoterAddress == "0x0000000000000000000000000000000000000000" || req.VoterAddress == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("voter address cannot be zero address")
		}
		if req.ApproverAddress == "" {
			return ErrInvalidInput("approver address cannot be empty or zero address")
		}
		accAddr2, err := sdk.AccAddressFromBech32(req.ApproverAddress)
		if err == nil {
			if accAddr2.Empty() {
				return ErrInvalidInput("approver address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr2 {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("approver address cannot be zero address")
			}
		}
		if req.ApproverAddress == "0x0000000000000000000000000000000000000000" || req.ApproverAddress == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("approver address cannot be zero address")
		}
	case *tokenizationtypes.QueryGetVotesRequest:
		// Validate approver address - required, cannot be empty or zero
		if req.ApproverAddress == "" {
			return ErrInvalidInput("approver address cannot be empty or zero address")
		}
		accAddr, err := sdk.AccAddressFromBech32(req.ApproverAddress)
		if err == nil {
			if accAddr.Empty() {
				return ErrInvalidInput("approver address cannot be zero address")
			}
			// Check if all bytes are zero
			allZero := true
			for _, b := range accAddr {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				return ErrInvalidInput("approver address cannot be zero address")
			}
		}
		if req.ApproverAddress == "0x0000000000000000000000000000000000000000" || req.ApproverAddress == "0000000000000000000000000000000000000000" {
			return ErrInvalidInput("approver address cannot be zero address")
		}
	case *tokenizationtypes.QueryGetChallengeTrackerRequest:
		// Validate leaf index (should be non-negative if provided)
		if req.LeafIndex != "" {
			leafIndex, err := sdkmath.ParseUint(req.LeafIndex)
			if err != nil {
				return ErrInvalidInput(fmt.Sprintf("invalid leaf index: %v", err))
			}
			// Leaf index can be 0, so we just validate it's a valid uint
			_ = leafIndex
		}
	case *tokenizationtypes.QueryGetAddressListRequest:
		// Validate list ID
		if req.ListId == "" {
			return ErrInvalidInput("list ID cannot be empty")
		}
	}
	return nil
}

// HandleGetBalanceAmount handles getBalanceAmount query with custom logic
// This returns the exact balance amount for a single (tokenId, ownershipTime) combination.
// Uses the same logic as the SDK's getBalanceForIdAndTime function.
func (p Precompile) HandleGetBalanceAmount(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Parse JSON request - single token ID and ownership time
	var req struct {
		CollectionId  string `json:"collectionId"`
		Address       string `json:"address"`
		TokenId       string `json:"tokenId"`
		OwnershipTime string `json:"ownershipTime"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON: %v", err))
	}

	// Convert collectionId
	collectionId, err := sdkmath.ParseUint(req.CollectionId)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid collectionId: %v", err))
	}

	// Parse tokenId
	tokenId, err := sdkmath.ParseUint(req.TokenId)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid tokenId: %v", err))
	}

	// Parse ownershipTime
	ownershipTime, err := sdkmath.ParseUint(req.OwnershipTime)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid ownershipTime: %v", err))
	}

	// Get collection
	collection, found := p.tokenizationKeeper.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotFound(req.CollectionId)
	}

	// Get user balance store
	userBalanceStore, _, err := p.tokenizationKeeper.GetBalanceOrApplyDefault(ctx, collection, req.Address)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "failed to get balance")
	}

	// Get the balance amount for this specific token ID and ownership time
	// This mirrors the SDK's getBalanceForIdAndTime function
	totalAmount := getBalanceForIdAndTime(userBalanceStore.Balances, tokenId, ownershipTime)

	// Emit event
	EmitGetBalanceAmountEventSingle(ctx, collectionId, req.Address, tokenId, ownershipTime, totalAmount)

	// Return uint256
	return method.Outputs.Pack(totalAmount.BigInt())
}

// getBalanceForIdAndTime returns the balance amount for a specific token ID and ownership time.
// This mirrors the SDK's getBalanceForIdAndTime function - it sums all matching balance amounts
// since the same (id, time) could appear in multiple overlapping Balance objects.
func getBalanceForIdAndTime(balances []*tokenizationtypes.Balance, tokenId sdkmath.Uint, ownershipTime sdkmath.Uint) sdkmath.Uint {
	amount := sdkmath.ZeroUint()

	for _, balance := range balances {
		if balance == nil {
			continue
		}

		// Check if tokenId is in this balance's token ID ranges
		foundTokenId := false
		for _, tokenRange := range balance.TokenIds {
			if tokenRange != nil && tokenId.GTE(tokenRange.Start) && tokenId.LTE(tokenRange.End) {
				foundTokenId = true
				break
			}
		}

		if !foundTokenId {
			continue
		}

		// Check if ownershipTime is in this balance's ownership time ranges
		foundTime := false
		for _, timeRange := range balance.OwnershipTimes {
			if timeRange != nil && ownershipTime.GTE(timeRange.Start) && ownershipTime.LTE(timeRange.End) {
				foundTime = true
				break
			}
		}

		if foundTime {
			amount = amount.Add(balance.Amount)
		}
	}

	return amount
}

// HandleGetTotalSupply handles getTotalSupply query with custom logic
// This returns the exact total supply for a single (tokenId, ownershipTime) combination.
func (p Precompile) HandleGetTotalSupply(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Parse JSON request - single token ID and ownership time
	var req struct {
		CollectionId  string `json:"collectionId"`
		TokenId       string `json:"tokenId"`
		OwnershipTime string `json:"ownershipTime"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON: %v", err))
	}

	// Convert collectionId
	collectionId, err := sdkmath.ParseUint(req.CollectionId)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid collectionId: %v", err))
	}

	// Parse tokenId
	tokenId, err := sdkmath.ParseUint(req.TokenId)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid tokenId: %v", err))
	}

	// Parse ownershipTime
	ownershipTime, err := sdkmath.ParseUint(req.OwnershipTime)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid ownershipTime: %v", err))
	}

	// Get collection (validate it exists)
	_, found := p.tokenizationKeeper.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotFound(req.CollectionId)
	}

	// Get all balances for the collection (from Total address)
	// Note: We need to get the Total address balance directly from the store,
	// as GetBalanceOrApplyDefault returns an empty store for Total address
	totalAddress := tokenizationtypes.TotalAddress
	balanceKey := tokenizationkeeper.ConstructBalanceKey(totalAddress, collectionId)
	userBalanceStore, found := p.tokenizationKeeper.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		// If Total address balance doesn't exist, return 0
		userBalanceStore = &tokenizationtypes.UserBalanceStore{Balances: []*tokenizationtypes.Balance{}}
	}

	// Get the total supply for this specific token ID and ownership time
	totalAmount := getBalanceForIdAndTime(userBalanceStore.Balances, tokenId, ownershipTime)

	// Emit event
	EmitGetTotalSupplyEventSingle(ctx, collectionId, tokenId, ownershipTime, totalAmount)

	// Return uint256
	return method.Outputs.Pack(totalAmount.BigInt())
}

// HandleGetAllReservedProtocolAddresses handles getAllReservedProtocolAddresses query with custom logic
func (p Precompile) HandleGetAllReservedProtocolAddresses(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Parse JSON request (empty request is fine)
	var req struct{}
	if jsonStr != "" && jsonStr != "{}" {
		if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON: %v", err))
		}
	}

	// Call keeper query
	queryReq := &tokenizationtypes.QueryGetAllReservedProtocolAddressesRequest{}
	resp, err := p.tokenizationKeeper.GetAllReservedProtocolAddresses(ctx, queryReq)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "query failed")
	}

	// Convert Cosmos addresses (strings) to EVM addresses
	evmAddresses := make([]common.Address, len(resp.Addresses))
	for i, cosmosAddr := range resp.Addresses {
		// Try to parse as Cosmos address and convert to EVM
		accAddr, parseErr := sdk.AccAddressFromBech32(cosmosAddr)
		if parseErr != nil {
			// If not a valid Cosmos address, try to convert directly from hex
			evmAddr := common.HexToAddress(cosmosAddr)
			evmAddresses[i] = evmAddr
		} else {
			// Convert Cosmos address bytes to EVM address
			evmAddresses[i] = common.BytesToAddress(accAddr.Bytes())
		}
	}

	// Pack as address[] according to ABI
	return method.Outputs.Pack(evmAddresses)
}

// HandleIsAddressReservedProtocol handles isAddressReservedProtocol query with custom logic
func (p Precompile) HandleIsAddressReservedProtocol(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Parse JSON request
	var req tokenizationtypes.QueryIsAddressReservedProtocolRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON: %v", err))
	}

	// Validate address - cannot be empty or zero
	if req.Address == "" {
		return nil, ErrInvalidInput("address cannot be empty or zero address")
	}
	accAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err == nil {
		if accAddr.Empty() {
			return nil, ErrInvalidInput("address cannot be zero address")
		}
		// Check if all bytes are zero
		allZero := true
		for _, b := range accAddr {
			if b != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return nil, ErrInvalidInput("address cannot be zero address")
		}
	}
	if req.Address == "0x0000000000000000000000000000000000000000" || req.Address == "0000000000000000000000000000000000000000" {
		return nil, ErrInvalidInput("address cannot be zero address")
	}

	// Call keeper query
	resp, err := p.tokenizationKeeper.IsAddressReservedProtocol(ctx, &req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "query failed")
	}

	// Return bool according to ABI
	return method.Outputs.Pack(resp.IsReservedProtocol)
}

// transactionMethods is a map of method names that are transactions (state-changing).
// Using a map provides O(1) lookup instead of O(n) switch statement.
var transactionMethods = map[string]bool{
	TransferTokensMethod:            true,
	SetIncomingApprovalMethod:       true,
	SetOutgoingApprovalMethod:       true,
	CreateCollectionMethod:          true,
	UpdateCollectionMethod:          true,
	DeleteCollectionMethod:          true,
	CreateAddressListsMethod:        true,
	UpdateUserApprovalsMethod:       true,
	DeleteIncomingApprovalMethod:    true,
	DeleteOutgoingApprovalMethod:    true,
	PurgeApprovalsMethod:            true,
	CreateDynamicStoreMethod:        true,
	UpdateDynamicStoreMethod:        true,
	DeleteDynamicStoreMethod:        true,
	SetDynamicStoreValueMethod:      true,
	SetValidTokenIdsMethod:          true,
	SetManagerMethod:                true,
	SetCollectionMetadataMethod:     true,
	SetTokenMetadataMethod:          true,
	SetCustomDataMethod:             true,
	SetStandardsMethod:              true,
	SetCollectionApprovalsMethod:    true,
	SetIsArchivedMethod:             true,
	CastVoteMethod:                  true,
	UniversalUpdateCollectionMethod: true,
	ExecuteMultipleMethod:           true,
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
// Uses O(1) map lookup for better performance.
func (Precompile) IsTransaction(method *abi.Method) bool {
	return transactionMethods[method.Name]
}

// Method name constants
const (
	// Transaction methods
	TransferTokensMethod            = "transferTokens"
	SetIncomingApprovalMethod       = "setIncomingApproval"
	SetOutgoingApprovalMethod       = "setOutgoingApproval"
	CreateCollectionMethod          = "createCollection"
	UpdateCollectionMethod          = "updateCollection"
	DeleteCollectionMethod          = "deleteCollection"
	CreateAddressListsMethod        = "createAddressLists"
	UpdateUserApprovalsMethod       = "updateUserApprovals"
	DeleteIncomingApprovalMethod    = "deleteIncomingApproval"
	DeleteOutgoingApprovalMethod    = "deleteOutgoingApproval"
	PurgeApprovalsMethod            = "purgeApprovals"
	CreateDynamicStoreMethod        = "createDynamicStore"
	UpdateDynamicStoreMethod        = "updateDynamicStore"
	DeleteDynamicStoreMethod        = "deleteDynamicStore"
	SetDynamicStoreValueMethod      = "setDynamicStoreValue"
	SetValidTokenIdsMethod          = "setValidTokenIds"
	SetManagerMethod                = "setManager"
	SetCollectionMetadataMethod     = "setCollectionMetadata"
	SetTokenMetadataMethod          = "setTokenMetadata"
	SetCustomDataMethod             = "setCustomData"
	SetStandardsMethod              = "setStandards"
	SetCollectionApprovalsMethod    = "setCollectionApprovals"
	SetIsArchivedMethod             = "setIsArchived"
	CastVoteMethod                  = "castVote"
	UniversalUpdateCollectionMethod = "universalUpdateCollection"
	ExecuteMultipleMethod           = "executeMultiple"

	// Query methods
	GetCollectionMethod                   = "getCollection"
	GetBalanceMethod                      = "getBalance"
	GetAddressListMethod                  = "getAddressList"
	GetApprovalTrackerMethod              = "getApprovalTracker"
	GetChallengeTrackerMethod             = "getChallengeTracker"
	GetETHSignatureTrackerMethod          = "getETHSignatureTracker"
	GetDynamicStoreMethod                 = "getDynamicStore"
	GetDynamicStoreValueMethod            = "getDynamicStoreValue"
	GetWrappableBalancesMethod            = "getWrappableBalances"
	IsAddressReservedProtocolMethod       = "isAddressReservedProtocol"
	GetAllReservedProtocolAddressesMethod = "getAllReservedProtocolAddresses"
	GetVoteMethod                         = "getVote"
	GetVotesMethod                        = "getVotes"
	ParamsMethod                          = "params"
	GetBalanceAmountMethod                = "getBalanceAmount"
	GetTotalSupplyMethod                  = "getTotalSupply"
)
