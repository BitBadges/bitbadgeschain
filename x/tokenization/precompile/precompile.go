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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

const (
	// Base gas costs for transactions
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

	// For methods that require dynamic gas calculation, we return a base amount
	// The actual gas will be calculated in Execute based on input size
	switch method.Name {
	// Transaction methods
	case TransferTokensMethod:
		return GasTransferTokensBase
	case SetIncomingApprovalMethod:
		return GasSetIncomingApprovalBase
	case SetOutgoingApprovalMethod:
		return GasSetOutgoingApprovalBase
	case CreateCollectionMethod:
		return GasCreateCollectionBase
	case UpdateCollectionMethod:
		return GasUpdateCollectionBase
	case DeleteCollectionMethod:
		return GasDeleteCollectionBase
	case CreateAddressListsMethod:
		return GasCreateAddressListsBase
	case UpdateUserApprovalsMethod:
		return GasUpdateUserApprovalsBase
	case DeleteIncomingApprovalMethod:
		return GasDeleteIncomingApprovalBase
	case DeleteOutgoingApprovalMethod:
		return GasDeleteOutgoingApprovalBase
	case PurgeApprovalsMethod:
		return GasPurgeApprovalsBase
	case CreateDynamicStoreMethod:
		return GasCreateDynamicStoreBase
	case UpdateDynamicStoreMethod:
		return GasUpdateDynamicStoreBase
	case DeleteDynamicStoreMethod:
		return GasDeleteDynamicStoreBase
	case SetDynamicStoreValueMethod:
		return GasSetDynamicStoreValueBase
	case SetValidTokenIdsMethod:
		return GasSetValidTokenIdsBase
	case SetManagerMethod:
		return GasSetManagerBase
	case SetCollectionMetadataMethod:
		return GasSetCollectionMetadataBase
	case SetTokenMetadataMethod:
		return GasSetTokenMetadataBase
	case SetCustomDataMethod:
		return GasSetCustomDataBase
	case SetStandardsMethod:
		return GasSetStandardsBase
	case SetCollectionApprovalsMethod:
		return GasSetCollectionApprovalsBase
	case SetIsArchivedMethod:
		return GasSetIsArchivedBase
	case CastVoteMethod:
		return GasCastVoteBase
	case UniversalUpdateCollectionMethod:
		return GasUniversalUpdateCollectionBase
	// Query methods
	case GetCollectionMethod:
		return GasGetCollectionBase
	case GetBalanceMethod:
		return GasGetBalanceBase
	case GetAddressListMethod:
		return GasGetAddressList
	case GetApprovalTrackerMethod:
		return GasGetApprovalTracker
	case GetChallengeTrackerMethod:
		return GasGetChallengeTracker
	case GetETHSignatureTrackerMethod:
		return GasGetETHSignatureTracker
	case GetDynamicStoreMethod:
		return GasGetDynamicStore
	case GetDynamicStoreValueMethod:
		return GasGetDynamicStoreValue
	case GetWrappableBalancesMethod:
		return GasGetWrappableBalances
	case IsAddressReservedProtocolMethod:
		return GasIsAddressReservedProtocol
	case GetAllReservedProtocolAddressesMethod:
		return GasGetAllReservedProtocol
	case GetVoteMethod:
		return GasGetVote
	case GetVotesMethod:
		return GasGetVotes
	case ParamsMethod:
		return GasParams
	case GetBalanceAmountMethod:
		return GasGetBalanceAmountBase
	case GetTotalSupplyMethod:
		return GasGetTotalSupplyBase
	}

	return 0
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	// Check if ABI loaded successfully during init
	if abiLoadError != nil {
		return nil, fmt.Errorf("tokenization precompile unavailable: ABI failed to load: %w", abiLoadError)
	}

	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
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
		return nil, "unknown", err
	}

	var bz []byte
	switch method.Name {
	// Transactions
	case TransferTokensMethod:
		bz, err = p.TransferTokens(ctx, method, args, contract)
	case SetIncomingApprovalMethod:
		bz, err = p.SetIncomingApproval(ctx, method, args, contract)
	case SetOutgoingApprovalMethod:
		bz, err = p.SetOutgoingApproval(ctx, method, args, contract)
	case CreateCollectionMethod:
		bz, err = p.CreateCollection(ctx, method, args, contract)
	case UpdateCollectionMethod:
		bz, err = p.UpdateCollection(ctx, method, args, contract)
	case DeleteCollectionMethod:
		bz, err = p.DeleteCollection(ctx, method, args, contract)
	case CreateAddressListsMethod:
		bz, err = p.CreateAddressLists(ctx, method, args, contract)
	case UpdateUserApprovalsMethod:
		bz, err = p.UpdateUserApprovals(ctx, method, args, contract)
	case DeleteIncomingApprovalMethod:
		bz, err = p.DeleteIncomingApproval(ctx, method, args, contract)
	case DeleteOutgoingApprovalMethod:
		bz, err = p.DeleteOutgoingApproval(ctx, method, args, contract)
	case PurgeApprovalsMethod:
		bz, err = p.PurgeApprovals(ctx, method, args, contract)
	case CreateDynamicStoreMethod:
		bz, err = p.CreateDynamicStore(ctx, method, args, contract)
	case UpdateDynamicStoreMethod:
		bz, err = p.UpdateDynamicStore(ctx, method, args, contract)
	case DeleteDynamicStoreMethod:
		bz, err = p.DeleteDynamicStore(ctx, method, args, contract)
	case SetDynamicStoreValueMethod:
		bz, err = p.SetDynamicStoreValue(ctx, method, args, contract)
	case SetValidTokenIdsMethod:
		bz, err = p.SetValidTokenIds(ctx, method, args, contract)
	case SetManagerMethod:
		bz, err = p.SetManager(ctx, method, args, contract)
	case SetCollectionMetadataMethod:
		bz, err = p.SetCollectionMetadata(ctx, method, args, contract)
	case SetTokenMetadataMethod:
		bz, err = p.SetTokenMetadata(ctx, method, args, contract)
	case SetCustomDataMethod:
		bz, err = p.SetCustomData(ctx, method, args, contract)
	case SetStandardsMethod:
		bz, err = p.SetStandards(ctx, method, args, contract)
	case SetCollectionApprovalsMethod:
		bz, err = p.SetCollectionApprovals(ctx, method, args, contract)
	case SetIsArchivedMethod:
		bz, err = p.SetIsArchived(ctx, method, args, contract)
	case CastVoteMethod:
		bz, err = p.CastVote(ctx, method, args, contract)
	case UniversalUpdateCollectionMethod:
		bz, err = p.UniversalUpdateCollection(ctx, method, args, contract)
	// Queries
	case GetCollectionMethod:
		bz, err = p.GetCollection(ctx, method, args)
	case GetBalanceMethod:
		bz, err = p.GetBalance(ctx, method, args)
	case GetAddressListMethod:
		bz, err = p.GetAddressList(ctx, method, args)
	case GetApprovalTrackerMethod:
		bz, err = p.GetApprovalTracker(ctx, method, args)
	case GetChallengeTrackerMethod:
		bz, err = p.GetChallengeTracker(ctx, method, args)
	case GetETHSignatureTrackerMethod:
		bz, err = p.GetETHSignatureTracker(ctx, method, args)
	case GetDynamicStoreMethod:
		bz, err = p.GetDynamicStore(ctx, method, args)
	case GetDynamicStoreValueMethod:
		bz, err = p.GetDynamicStoreValue(ctx, method, args)
	case GetWrappableBalancesMethod:
		bz, err = p.GetWrappableBalances(ctx, method, args)
	case IsAddressReservedProtocolMethod:
		bz, err = p.IsAddressReservedProtocol(ctx, method, args)
	case GetAllReservedProtocolAddressesMethod:
		bz, err = p.GetAllReservedProtocolAddresses(ctx, method, args)
	case GetVoteMethod:
		bz, err = p.GetVote(ctx, method, args)
	case GetVotesMethod:
		bz, err = p.GetVotes(ctx, method, args)
	case ParamsMethod:
		bz, err = p.Params(ctx, method, args)
	case GetBalanceAmountMethod:
		bz, err = p.GetBalanceAmount(ctx, method, args)
	case GetTotalSupplyMethod:
		bz, err = p.GetTotalSupply(ctx, method, args)
	default:
		return nil, method.Name, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, method.Name, err
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

// TransferTokens executes a token transfer via the tokenization module.
// Transfers tokens from the caller (msg.sender) to one or more recipient addresses.
//
// Parameters:
//   - collectionId: The collection ID to transfer from (uint256)
//   - toAddresses: Array of recipient EVM addresses (address[])
//   - amount: Amount to transfer to each recipient (uint256)
//   - tokenIds: Array of token ID ranges to transfer (UintRange[])
//   - ownershipTimes: Array of ownership time ranges to transfer (UintRange[])
//
// Returns:
//   - bool: true if transfer succeeded
//
// Errors:
//   - ErrorCodeInvalidInput: Invalid input parameters (zero addresses, invalid ranges, etc.)
//   - ErrorCodeCollectionNotFound: Collection does not exist
//   - ErrorCodeTransferFailed: Transfer operation failed (insufficient balance, approval issues, etc.)
func (p Precompile) TransferTokens(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 5 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 5, got %d", len(args)))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}
	toAddresses, ok := args[1].([]common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid toAddresses type, expected []common.Address")
	}
	amountBig, ok := args[2].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid amount type, expected *big.Int")
	}
	tokenIdsRanges, ok := args[3].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid tokenIds type, expected []struct{Start, End *big.Int}")
	}
	ownershipTimesRanges, ok := args[4].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid ownershipTimes type, expected []struct{Start, End *big.Int}")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateAddresses(toAddresses, "toAddresses"); err != nil {
		return nil, err
	}
	if err := ValidateAmount(amountBig, "amount"); err != nil {
		return nil, err
	}
	if err := ValidateBigIntRanges(tokenIdsRanges, "tokenIds"); err != nil {
		return nil, err
	}
	if err := ValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes"); err != nil {
		return nil, err
	}
	if err := ValidateTransferInputs(toAddresses, tokenIdsRanges, ownershipTimesRanges); err != nil {
		return nil, err
	}

	// Security: Verify caller and check for overflow
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	if err := CheckOverflow(amountBig, "amount"); err != nil {
		return nil, err
	}
	fromCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Convert toAddresses from EVM addresses to Cosmos addresses
	toCosmosAddrs := make([]string, len(toAddresses))
	for i, addr := range toAddresses {
		toCosmosAddrs[i] = sdk.AccAddress(addr.Bytes()).String()
	}

	// Convert collectionId
	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Convert amount
	amount := sdkmath.NewUintFromBigInt(amountBig)

	// Convert and validate tokenIds ranges (overlaps allowed for token IDs)
	tokenIds, err := ConvertAndValidateBigIntRanges(tokenIdsRanges, "tokenIds")
	if err != nil {
		return nil, err
	}

	// Convert and validate ownershipTimes ranges (overlaps allowed for ownership times)
	ownershipTimes, err := ConvertAndValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes")
	if err != nil {
		return nil, err
	}

	// Note: Overlapping ranges are allowed for tokenIds and ownershipTimes
	// as they represent valid ranges that may overlap in the collection's state

	// Create the transfer message
	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      fromCosmosAddr,
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        fromCosmosAddr,
				ToAddresses: toCosmosAddrs,
				Balances: []*tokenizationtypes.Balance{
					{
						Amount:         amount,
						TokenIds:       tokenIds,
						OwnershipTimes: ownershipTimes,
					},
				},
			},
		},
	}

	// Execute the transfer via the keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.TransferTokens(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeTransferFailed, "transfer operation failed")
	}

	// Emit event
	EmitTransferEvent(ctx, collectionId, caller, toAddresses, amount, tokenIds, ownershipTimes)

	// Return success
	return method.Outputs.Pack(true)
}

// SetIncomingApproval executes setting an incoming approval via the tokenization module.
// Sets an approval that allows others to transfer tokens to the caller.
//
// Parameters:
//   - collectionId: The collection ID (uint256)
//   - approval: UserIncomingApproval struct containing:
//   - approvalId: Unique identifier for the approval (string)
//   - fromListId: List ID of addresses that can transfer to the caller (string)
//   - initiatedByListId: List ID of addresses that can initiate the transfer (string)
//   - transferTimes: Array of transfer time ranges (UintRange[])
//   - tokenIds: Array of token ID ranges (UintRange[])
//   - ownershipTimes: Array of ownership time ranges (UintRange[])
//   - uri: Optional URI for the approval (string)
//   - customData: Optional custom data (string)
//
// Returns:
//   - bool: true if approval was set successfully
//
// Errors:
//   - ErrorCodeInvalidInput: Invalid input parameters
//   - ErrorCodeCollectionNotFound: Collection does not exist
//   - ErrorCodeApprovalFailed: Approval operation failed
func (p Precompile) SetIncomingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	// Try to extract approval - handle both struct (from ABI) and map (from tests)
	var approval *tokenizationtypes.UserIncomingApproval
	var err error

	if approvalMap, ok := args[1].(map[string]interface{}); ok {
		// Handle map format (for testing)
		approval, err = ConvertUserIncomingApproval(approvalMap)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("approval conversion failed: %v", err))
		}
	} else {
		// Try struct format (from ABI unpacking)
		approvalStruct, ok := args[1].(struct {
			ApprovalId        string `json:"approvalId"`
			FromListId        string `json:"fromListId"`
			InitiatedByListId string `json:"initiatedByListId"`
			TransferTimes     []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"transferTimes"`
			TokenIds []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"tokenIds"`
			OwnershipTimes []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"ownershipTimes"`
			Uri        string `json:"uri"`
			CustomData string `json:"customData"`
		})
		if !ok {
			return nil, ErrInvalidInput("invalid approval type, expected UserIncomingApproval struct or map")
		}

		// Convert struct to UserIncomingApproval
		approval = &tokenizationtypes.UserIncomingApproval{
			ApprovalId:        approvalStruct.ApprovalId,
			FromListId:        approvalStruct.FromListId,
			InitiatedByListId: approvalStruct.InitiatedByListId,
			Uri:               approvalStruct.Uri,
			CustomData:        approvalStruct.CustomData,
			Version:           sdkmath.NewUint(0),
		}

		// Convert and validate transferTimes
		approval.TransferTimes, err = ConvertAndValidateBigIntRanges(approvalStruct.TransferTimes, "transferTimes")
		if err != nil {
			return nil, err
		}

		// Convert and validate tokenIds
		approval.TokenIds, err = ConvertAndValidateBigIntRanges(approvalStruct.TokenIds, "tokenIds")
		if err != nil {
			return nil, err
		}

		// Convert and validate ownershipTimes
		approval.OwnershipTimes, err = ConvertAndValidateBigIntRanges(approvalStruct.OwnershipTimes, "ownershipTimes")
		if err != nil {
			return nil, err
		}

		// Initialize empty approval criteria
		approval.ApprovalCriteria = &tokenizationtypes.IncomingApprovalCriteria{}
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateString(approval.ApprovalId, "approvalId"); err != nil {
		return nil, err
	}

	// Security: Verify caller
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	fromCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Convert collectionId
	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Create the message
	msg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      fromCosmosAddr,
		CollectionId: collectionId,
		Approval:     approval,
	}

	// Execute via the keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetIncomingApproval(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "set incoming approval failed")
	}

	// Emit event
	EmitIncomingApprovalEvent(ctx, collectionId, caller, approval.ApprovalId)

	// Return success
	return method.Outputs.Pack(true)
}

// SetOutgoingApproval executes setting an outgoing approval via the tokenization module.
// Sets an approval that allows the caller to transfer tokens to others.
//
// Parameters:
//   - collectionId: The collection ID (uint256)
//   - approval: UserOutgoingApproval struct containing:
//   - approvalId: Unique identifier for the approval (string)
//   - toListId: List ID of addresses that can receive tokens (string)
//   - initiatedByListId: List ID of addresses that can initiate the transfer (string)
//   - transferTimes: Array of transfer time ranges (UintRange[])
//   - tokenIds: Array of token ID ranges (UintRange[])
//   - ownershipTimes: Array of ownership time ranges (UintRange[])
//   - uri: Optional URI for the approval (string)
//   - customData: Optional custom data (string)
//
// Returns:
//   - bool: true if approval was set successfully
//
// Errors:
//   - ErrorCodeInvalidInput: Invalid input parameters
//   - ErrorCodeCollectionNotFound: Collection does not exist
//   - ErrorCodeApprovalFailed: Approval operation failed
func (p Precompile) SetOutgoingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	// Try to extract approval - handle both struct (from ABI) and map (from tests)
	var approval *tokenizationtypes.UserOutgoingApproval
	var err error

	if approvalMap, ok := args[1].(map[string]interface{}); ok {
		// Handle map format (for testing)
		approval, err = ConvertUserOutgoingApproval(approvalMap)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("approval conversion failed: %v", err))
		}
	} else {
		// Try struct format (from ABI unpacking)
		approvalStruct, ok := args[1].(struct {
			ApprovalId        string `json:"approvalId"`
			ToListId          string `json:"toListId"`
			InitiatedByListId string `json:"initiatedByListId"`
			TransferTimes     []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"transferTimes"`
			TokenIds []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"tokenIds"`
			OwnershipTimes []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			} `json:"ownershipTimes"`
			Uri        string `json:"uri"`
			CustomData string `json:"customData"`
		})
		if !ok {
			return nil, ErrInvalidInput("invalid approval type, expected UserOutgoingApproval struct or map")
		}

		// Convert struct to UserOutgoingApproval
		approval = &tokenizationtypes.UserOutgoingApproval{
			ApprovalId:        approvalStruct.ApprovalId,
			ToListId:          approvalStruct.ToListId,
			InitiatedByListId: approvalStruct.InitiatedByListId,
			Uri:               approvalStruct.Uri,
			CustomData:        approvalStruct.CustomData,
			Version:           sdkmath.NewUint(0),
		}

		// Convert and validate transferTimes
		approval.TransferTimes, err = ConvertAndValidateBigIntRanges(approvalStruct.TransferTimes, "transferTimes")
		if err != nil {
			return nil, err
		}

		// Convert and validate tokenIds
		approval.TokenIds, err = ConvertAndValidateBigIntRanges(approvalStruct.TokenIds, "tokenIds")
		if err != nil {
			return nil, err
		}

		// Convert and validate ownershipTimes
		approval.OwnershipTimes, err = ConvertAndValidateBigIntRanges(approvalStruct.OwnershipTimes, "ownershipTimes")
		if err != nil {
			return nil, err
		}

		// Initialize empty approval criteria
		approval.ApprovalCriteria = &tokenizationtypes.OutgoingApprovalCriteria{}
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateString(approval.ApprovalId, "approvalId"); err != nil {
		return nil, err
	}

	// Security: Verify caller
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	fromCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Convert collectionId
	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Create the message
	msg := &tokenizationtypes.MsgSetOutgoingApproval{
		Creator:      fromCosmosAddr,
		CollectionId: collectionId,
		Approval:     approval,
	}

	// Execute via the keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetOutgoingApproval(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "set outgoing approval failed")
	}

	// Emit event
	EmitOutgoingApprovalEvent(ctx, collectionId, caller, approval.ApprovalId)

	// Return success
	return method.Outputs.Pack(true)
}

// GetCollection queries a collection by ID
// Returns the collection data as a structured Solidity tuple (TokenCollection struct)
// Errors: ErrorCodeInvalidInput, ErrorCodeCollectionNotFound, ErrorCodeQueryFailed
func (p Precompile) GetCollection(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}

	// Validate collection ID
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid collectionId")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Query the collection
	req := &tokenizationtypes.QueryGetCollectionRequest{
		CollectionId: collectionId.String(),
	}
	resp, err := p.tokenizationKeeper.GetCollection(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get collection failed for collectionId: %s", collectionId.String()))
	}

	// Convert to Solidity struct and pack
	return PackCollectionAsStruct(method, resp.Collection)
}

// GetBalance queries a balance for a user address
// Returns the balance data as a structured Solidity tuple (UserBalanceStore struct)
// Errors: ErrorCodeInvalidInput, ErrorCodeCollectionNotFound, ErrorCodeBalanceNotFound, ErrorCodeQueryFailed
func (p Precompile) GetBalance(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}

	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid userAddress type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(userAddress, "userAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid userAddress")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	userCosmosAddr := sdk.AccAddress(userAddress.Bytes()).String()

	// Query the balance
	req := &tokenizationtypes.QueryGetBalanceRequest{
		CollectionId: collectionId.String(),
		Address:      userCosmosAddr,
	}
	resp, err := p.tokenizationKeeper.GetBalance(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get balance failed for collectionId: %s, address: %s", collectionId.String(), userCosmosAddr))
	}

	// Convert to Solidity struct and pack
	return PackUserBalanceStoreAsStruct(method, resp.Balance)
}

// GetAddressList queries an address list by ID
// Returns the address list data as a structured Solidity tuple (AddressList struct)
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetAddressList(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	listId, ok := args[0].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid listId type")
	}

	// Validate listId
	if err := ValidateString(listId, "listId"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid listId")
	}

	// Query the address list
	req := &tokenizationtypes.QueryGetAddressListRequest{
		ListId: listId,
	}
	resp, err := p.tokenizationKeeper.GetAddressList(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get address list failed for listId: %s", listId))
	}

	// Convert to Solidity struct and pack
	return PackAddressListAsStruct(method, resp.List)
}

// GetApprovalTracker queries an approval tracker
// Returns the approval tracker data as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetApprovalTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 7 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 7, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type")
	}
	amountTrackerId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid amountTrackerId type")
	}
	trackerType, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid trackerType type")
	}
	approvedAddress, ok := args[5].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approvedAddress type")
	}
	approvalId, ok := args[6].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid collectionId")
	}
	if err := ValidateAddress(approverAddress, "approverAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approverAddress")
	}
	if err := ValidateAddress(approvedAddress, "approvedAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approvedAddress")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	approverCosmosAddr := sdk.AccAddress(approverAddress.Bytes()).String()
	approvedCosmosAddr := sdk.AccAddress(approvedAddress.Bytes()).String()

	// Query the approval tracker
	req := &tokenizationtypes.QueryGetApprovalTrackerRequest{
		CollectionId:    collectionId.String(),
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverCosmosAddr,
		AmountTrackerId: amountTrackerId,
		TrackerType:     trackerType,
		ApprovedAddress: approvedCosmosAddr,
		ApprovalId:      approvalId,
	}
	resp, err := p.tokenizationKeeper.GetApprovalTracker(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get approval tracker failed for collectionId: %s", collectionId.String()))
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal approval tracker failed")
	}

	return method.Outputs.Pack(bz)
}

// GetChallengeTracker queries a challenge tracker
// Returns the number of times the challenge has been used as uint256
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetChallengeTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 6, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type")
	}
	challengeTrackerId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid challengeTrackerId type")
	}
	leafIndexBig, ok := args[4].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid leafIndex type")
	}
	approvalId, ok := args[5].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid collectionId")
	}
	if err := ValidateAddress(approverAddress, "approverAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approverAddress")
	}
	if err := CheckOverflow(leafIndexBig, "leafIndex"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid leafIndex")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	approverCosmosAddr := sdk.AccAddress(approverAddress.Bytes()).String()
	leafIndex := sdkmath.NewUintFromBigInt(leafIndexBig)

	// Query the challenge tracker
	req := &tokenizationtypes.QueryGetChallengeTrackerRequest{
		CollectionId:       collectionId.String(),
		ApprovalLevel:      approvalLevel,
		ApproverAddress:    approverCosmosAddr,
		ChallengeTrackerId: challengeTrackerId,
		LeafIndex:          leafIndex.String(),
		ApprovalId:         approvalId,
	}
	resp, err := p.tokenizationKeeper.GetChallengeTracker(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get challenge tracker failed for collectionId: %s", collectionId.String()))
	}

	// Convert numUsed string to uint256
	numUsed := sdkmath.NewUintFromString(resp.NumUsed)

	return method.Outputs.Pack(numUsed.BigInt())
}

// GetETHSignatureTracker queries an ETH signature tracker
// Returns the number of times the signature has been used as uint256
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetETHSignatureTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 6, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type")
	}
	challengeTrackerId, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid challengeTrackerId type")
	}
	signature, ok := args[5].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid signature type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(approverAddress, "approverAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approverAddress")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	approverCosmosAddr := sdk.AccAddress(approverAddress.Bytes()).String()

	// Query the ETH signature tracker
	req := &tokenizationtypes.QueryGetETHSignatureTrackerRequest{
		CollectionId:       collectionId.String(),
		ApprovalLevel:      approvalLevel,
		ApproverAddress:    approverCosmosAddr,
		ApprovalId:         approvalId,
		ChallengeTrackerId: challengeTrackerId,
		Signature:          signature,
	}
	resp, err := p.tokenizationKeeper.GetETHSignatureTracker(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get ETH signature tracker failed for collectionId: %s", collectionId.String()))
	}

	// Convert numUsed string to uint256
	numUsed := sdkmath.NewUintFromString(resp.NumUsed)

	return method.Outputs.Pack(numUsed.BigInt())
}

// GetDynamicStore queries a dynamic store
// Returns the dynamic store data as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetDynamicStore(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid storeId type")
	}

	// Validate storeId (treating it like a collectionId for validation)
	if err := ValidateCollectionId(storeIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)

	// Query the dynamic store
	req := &tokenizationtypes.QueryGetDynamicStoreRequest{
		StoreId: storeId.String(),
	}
	resp, err := p.tokenizationKeeper.GetDynamicStore(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get dynamic store failed for storeId: %s", storeId.String()))
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal dynamic store failed")
	}

	return method.Outputs.Pack(bz)
}

// GetDynamicStoreValue queries a dynamic store value for a user address
// Returns the dynamic store value data as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetDynamicStoreValue(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid storeId type")
	}
	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid userAddress type")
	}

	// Validate inputs
	if err := ValidateCollectionId(storeIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(userAddress, "userAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid userAddress")
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)
	userCosmosAddr := sdk.AccAddress(userAddress.Bytes()).String()

	// Query the dynamic store value
	req := &tokenizationtypes.QueryGetDynamicStoreValueRequest{
		StoreId: storeId.String(),
		Address: userCosmosAddr,
	}
	resp, err := p.tokenizationKeeper.GetDynamicStoreValue(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get dynamic store value failed for storeId: %s, address: %s", storeId.String(), userCosmosAddr))
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal dynamic store value failed")
	}

	return method.Outputs.Pack(bz)
}

// GetWrappableBalances queries wrappable balances for a user address
// Returns the wrappable balance amount as uint256
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetWrappableBalances(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	denom, ok := args[0].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid denom type")
	}
	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid userAddress type")
	}

	// Validate inputs
	if err := ValidateAddress(userAddress, "userAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid userAddress")
	}

	userCosmosAddr := sdk.AccAddress(userAddress.Bytes()).String()

	// Query wrappable balances
	req := &tokenizationtypes.QueryGetWrappableBalancesRequest{
		Denom:   denom,
		Address: userCosmosAddr,
	}
	resp, err := p.tokenizationKeeper.GetWrappableBalances(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get wrappable balances failed for denom: %s, address: %s", denom, userCosmosAddr))
	}

	return method.Outputs.Pack(resp.Amount.BigInt())
}

// IsAddressReservedProtocol checks if an address is reserved protocol
// Returns true if the address is reserved protocol, false otherwise
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) IsAddressReservedProtocol(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	addr, ok := args[0].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid addr type")
	}

	// Validate address
	if err := ValidateAddress(addr, "addr"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid addr")
	}

	cosmosAddr := sdk.AccAddress(addr.Bytes()).String()

	// Query if address is reserved protocol
	req := &tokenizationtypes.QueryIsAddressReservedProtocolRequest{
		Address: cosmosAddr,
	}
	resp, err := p.tokenizationKeeper.IsAddressReservedProtocol(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("is address reserved protocol failed for address: %s", cosmosAddr))
	}

	return method.Outputs.Pack(resp.IsReservedProtocol)
}

// GetAllReservedProtocolAddresses gets all reserved protocol addresses
// Returns an array of EVM addresses that are reserved protocol addresses
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed, ErrorCodeInternalError
func (p Precompile) GetAllReservedProtocolAddresses(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 0, got %d", len(args)))
	}

	// Query all reserved protocol addresses
	req := &tokenizationtypes.QueryGetAllReservedProtocolAddressesRequest{}
	resp, err := p.tokenizationKeeper.GetAllReservedProtocolAddresses(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "get all reserved protocol addresses failed")
	}

	// Convert Cosmos addresses to EVM addresses
	evmAddresses := make([]common.Address, len(resp.Addresses))
	for i, addrStr := range resp.Addresses {
		cosmosAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, WrapError(err, ErrorCodeInternalError, fmt.Sprintf("invalid address %s", addrStr))
		}
		evmAddresses[i] = common.BytesToAddress(cosmosAddr.Bytes())
	}

	return method.Outputs.Pack(evmAddresses)
}

// GetVote queries a vote for a specific proposal
// Returns the vote data as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetVote(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 6, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type")
	}
	proposalId, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid proposalId type")
	}
	voterAddress, ok := args[5].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid voterAddress type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(approverAddress, "approverAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approverAddress")
	}
	if err := ValidateAddress(voterAddress, "voterAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid voterAddress")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	approverCosmosAddr := sdk.AccAddress(approverAddress.Bytes()).String()
	voterCosmosAddr := sdk.AccAddress(voterAddress.Bytes()).String()

	// Query the vote
	req := &tokenizationtypes.QueryGetVoteRequest{
		CollectionId:    collectionId.String(),
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverCosmosAddr,
		ApprovalId:      approvalId,
		ProposalId:      proposalId,
		VoterAddress:    voterCosmosAddr,
	}
	resp, err := p.tokenizationKeeper.GetVote(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get vote failed for collectionId: %s", collectionId.String()))
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal vote failed")
	}

	return method.Outputs.Pack(bz)
}

// GetVotes queries all votes for a proposal
// Returns all votes data as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) GetVotes(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 5 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 5, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type")
	}
	proposalId, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid proposalId type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(approverAddress, "approverAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid approverAddress")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	approverCosmosAddr := sdk.AccAddress(approverAddress.Bytes()).String()

	// Query the votes
	req := &tokenizationtypes.QueryGetVotesRequest{
		CollectionId:    collectionId.String(),
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverCosmosAddr,
		ApprovalId:      approvalId,
		ProposalId:      proposalId,
	}
	resp, err := p.tokenizationKeeper.GetVotes(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("get votes failed for collectionId: %s", collectionId.String()))
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal votes failed")
	}

	return method.Outputs.Pack(bz)
}

// Params queries the module parameters
// Returns the module parameters as protobuf-encoded bytes
// Errors: ErrorCodeInvalidInput, ErrorCodeQueryFailed
func (p Precompile) Params(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 0, got %d", len(args)))
	}

	// Query params
	req := &tokenizationtypes.QueryParamsRequest{}
	resp, err := p.tokenizationKeeper.Params(ctx, req)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "get params failed")
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "marshal params failed")
	}

	return method.Outputs.Pack(bz)
}

// GetBalanceAmount queries a balance amount for a user address with specific token IDs and ownership times
// Returns the total balance amount as uint256
// Errors: ErrorCodeInvalidInput, ErrorCodeCollectionNotFound, ErrorCodeQueryFailed
func (p Precompile) GetBalanceAmount(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 4 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 4, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}

	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid userAddress type")
	}

	tokenIdsRanges, ok := args[2].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid tokenIds type")
	}

	ownershipTimesRanges, ok := args[3].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid ownershipTimes type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateAddress(userAddress, "userAddress"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid userAddress")
	}
	if err := ValidateBigIntRanges(tokenIdsRanges, "tokenIds"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	userCosmosAddr := sdk.AccAddress(userAddress.Bytes()).String()

	// Get the collection
	collection, found := p.tokenizationKeeper.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotFound(collectionId.String())
	}

	// Get user balance
	userBalance, _, err := p.tokenizationKeeper.GetBalanceOrApplyDefault(ctx, collection, userCosmosAddr)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "get balance failed")
	}

	// Convert tokenIds ranges
	tokenIds, err := ConvertAndValidateBigIntRanges(tokenIdsRanges, "tokenIds")
	if err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid tokenIds ranges")
	}

	// Convert ownershipTimes ranges
	ownershipTimes, err := ConvertAndValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes")
	if err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "invalid ownershipTimes ranges")
	}

	// Get balances for the specified token IDs and ownership times
	fetchedBalances, err := tokenizationtypes.GetBalancesForIds(ctx, tokenIds, ownershipTimes, userBalance.Balances)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "get balances for ids failed")
	}

	// Sum up all amounts
	totalAmount := sdkmath.NewUint(0)
	for _, balance := range fetchedBalances {
		totalAmount = totalAmount.Add(balance.Amount)
	}

	// Emit event
	EmitGetBalanceAmountEvent(ctx, collectionId, userCosmosAddr, tokenIds, ownershipTimes, totalAmount)

	return method.Outputs.Pack(totalAmount.BigInt())
}

// GetTotalSupply queries the total supply for a collection with specific token IDs and ownership times
// Returns the total supply amount as uint256
// Errors: ErrorCodeInvalidInput, ErrorCodeCollectionNotFound, ErrorCodeQueryFailed
func (p Precompile) GetTotalSupply(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 3 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 3, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type")
	}

	tokenIdsRanges, ok := args[1].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid tokenIds type")
	}

	ownershipTimesRanges, ok := args[2].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, ErrInvalidInput("invalid ownershipTimes type")
	}

	// Validate inputs
	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateBigIntRanges(tokenIdsRanges, "tokenIds"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}
	if err := ValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes"); err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Verify collection exists
	_, found := p.tokenizationKeeper.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotFound(collectionId.String())
	}

	// Get total balance (Total address)
	// Note: GetBalanceOrApplyDefault returns empty for Total address, so we need to get it directly from store
	balanceKey := tokenizationkeeper.ConstructBalanceKey(tokenizationtypes.TotalAddress, collectionId)
	totalBalance, found := p.tokenizationKeeper.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		// If Total balance doesn't exist yet, return 0
		totalBalance = &tokenizationtypes.UserBalanceStore{
			Balances: []*tokenizationtypes.Balance{},
		}
	}

	// Convert and validate tokenIds ranges
	tokenIds, err := ConvertAndValidateBigIntRanges(tokenIdsRanges, "tokenIds")
	if err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}

	// Convert and validate ownershipTimes ranges
	ownershipTimes, err := ConvertAndValidateBigIntRanges(ownershipTimesRanges, "ownershipTimes")
	if err != nil {
		return nil, WrapError(err, ErrorCodeInvalidInput, "validation failed")
	}

	// Get balances for the specified token IDs and ownership times
	fetchedBalances, err := tokenizationtypes.GetBalancesForIds(ctx, tokenIds, ownershipTimes, totalBalance.Balances)
	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, "get balances for ids failed")
	}

	// Sum up all amounts
	totalAmount := sdkmath.NewUint(0)
	for _, balance := range fetchedBalances {
		totalAmount = totalAmount.Add(balance.Amount)
	}

	// Emit event
	EmitGetTotalSupplyEvent(ctx, collectionId, tokenIds, ownershipTimes, totalAmount)

	return method.Outputs.Pack(totalAmount.BigInt())
}
