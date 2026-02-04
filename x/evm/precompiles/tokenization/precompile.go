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
	// Gas costs for transactions
	GasTransferTokens      = 50_000
	GasSetIncomingApproval = 30_000
	GasSetOutgoingApproval = 30_000

	// Gas costs for queries (lower since they're read-only)
	GasGetCollection             = 5_000
	GasGetBalance                = 5_000
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
)

var _ vm.PrecompiledContract = &Precompile{}

var (
	// Embed abi json file to the executable binary. Needed when importing as dependency.
	//
	//go:embed abi.json
	f   embed.FS
	ABI abi.ABI
)

func init() {
	var err error
	ABI, err = cmn.LoadABI(f, "abi.json")
	if err != nil {
		panic(err)
	}
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

// TokenizationPrecompileAddress is the address of the badges precompile
// Using standard precompile address range: 0x0000000000000000000000000000000000001001
const TokenizationPrecompileAddress = "0x0000000000000000000000000000000000001001"

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

	switch method.Name {
	case TransferTokensMethod:
		return GasTransferTokens
	case SetIncomingApprovalMethod:
		return GasSetIncomingApproval
	case SetOutgoingApprovalMethod:
		return GasSetOutgoingApproval
	case GetCollectionMethod:
		return GasGetCollection
	case GetBalanceMethod:
		return GasGetBalance
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
	}

	return 0
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		return p.Execute(ctx, contract, readonly)
	})
}

// Execute executes the precompiled contract tokenization methods defined in the ABI.
func (p Precompile) Execute(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, error) {
	method, args, err := cmn.SetupABI(p.ABI, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
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
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, err
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
func (Precompile) IsTransaction(method *abi.Method) bool {
	switch method.Name {
	case TransferTokensMethod, SetIncomingApprovalMethod, SetOutgoingApprovalMethod:
		return true
	default:
		// All other methods are queries (view functions)
		return false
	}
}

// Method name constants
const (
	TransferTokensMethod                  = "transferTokens"
	SetIncomingApprovalMethod             = "setIncomingApproval"
	SetOutgoingApprovalMethod             = "setOutgoingApproval"
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
)

// TransferTokens executes a token transfer via the tokenization module
func (p Precompile) TransferTokens(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("invalid number of arguments, expected 5, got %d", len(args))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	toAddresses, ok := args[1].([]common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid toAddresses type")
	}
	amountBig, ok := args[2].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid amount type")
	}
	tokenIdsRanges, ok := args[3].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, fmt.Errorf("invalid tokenIds type")
	}
	ownershipTimesRanges, ok := args[4].([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	})
	if !ok {
		return nil, fmt.Errorf("invalid ownershipTimes type")
	}

	// Get the caller (msg.sender) from the contract
	caller := contract.Caller()
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

	// Convert tokenIds ranges
	tokenIds := make([]*tokenizationtypes.UintRange, len(tokenIdsRanges))
	for i, r := range tokenIdsRanges {
		tokenIds[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Convert ownershipTimes ranges
	ownershipTimes := make([]*tokenizationtypes.UintRange, len(ownershipTimesRanges))
	for i, r := range ownershipTimesRanges {
		ownershipTimes[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

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
	_, err := msgServer.TransferTokens(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("transfer failed: %w", err)
	}

	// Return success
	return method.Outputs.Pack(true)
}

// SetIncomingApproval executes setting an incoming approval via the tokenization module
func (p Precompile) SetIncomingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments, expected 2, got %d", len(args))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}

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
		return nil, fmt.Errorf("invalid approval type")
	}

	// Get the caller (msg.sender) from the contract
	caller := contract.Caller()
	fromCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Convert collectionId
	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Convert approval struct to UserIncomingApproval
	approval := &tokenizationtypes.UserIncomingApproval{
		ApprovalId:        approvalStruct.ApprovalId,
		FromListId:        approvalStruct.FromListId,
		InitiatedByListId: approvalStruct.InitiatedByListId,
		Uri:               approvalStruct.Uri,
		CustomData:        approvalStruct.CustomData,
		Version:           sdkmath.NewUint(0),
	}

	// Convert transferTimes
	approval.TransferTimes = make([]*tokenizationtypes.UintRange, len(approvalStruct.TransferTimes))
	for i, r := range approvalStruct.TransferTimes {
		approval.TransferTimes[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Convert tokenIds
	approval.TokenIds = make([]*tokenizationtypes.UintRange, len(approvalStruct.TokenIds))
	for i, r := range approvalStruct.TokenIds {
		approval.TokenIds[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Convert ownershipTimes
	approval.OwnershipTimes = make([]*tokenizationtypes.UintRange, len(approvalStruct.OwnershipTimes))
	for i, r := range approvalStruct.OwnershipTimes {
		approval.OwnershipTimes[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Initialize empty approval criteria
	approval.ApprovalCriteria = &tokenizationtypes.IncomingApprovalCriteria{}

	// Create the message
	msg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      fromCosmosAddr,
		CollectionId: collectionId,
		Approval:     approval,
	}

	// Execute via the keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err := msgServer.SetIncomingApproval(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("set incoming approval failed: %w", err)
	}

	// Return success
	return method.Outputs.Pack(true)
}

// SetOutgoingApproval executes setting an outgoing approval via the tokenization module
func (p Precompile) SetOutgoingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments, expected 2, got %d", len(args))
	}

	// Extract arguments
	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}

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
		return nil, fmt.Errorf("invalid approval type")
	}

	// Get the caller (msg.sender) from the contract
	caller := contract.Caller()
	fromCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Convert collectionId
	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Convert approval struct to UserOutgoingApproval
	approval := &tokenizationtypes.UserOutgoingApproval{
		ApprovalId:        approvalStruct.ApprovalId,
		ToListId:          approvalStruct.ToListId,
		InitiatedByListId: approvalStruct.InitiatedByListId,
		Uri:               approvalStruct.Uri,
		CustomData:        approvalStruct.CustomData,
		Version:           sdkmath.NewUint(0),
	}

	// Convert transferTimes
	approval.TransferTimes = make([]*tokenizationtypes.UintRange, len(approvalStruct.TransferTimes))
	for i, r := range approvalStruct.TransferTimes {
		approval.TransferTimes[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Convert tokenIds
	approval.TokenIds = make([]*tokenizationtypes.UintRange, len(approvalStruct.TokenIds))
	for i, r := range approvalStruct.TokenIds {
		approval.TokenIds[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Convert ownershipTimes
	approval.OwnershipTimes = make([]*tokenizationtypes.UintRange, len(approvalStruct.OwnershipTimes))
	for i, r := range approvalStruct.OwnershipTimes {
		approval.OwnershipTimes[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	// Initialize empty approval criteria
	approval.ApprovalCriteria = &tokenizationtypes.OutgoingApprovalCriteria{}

	// Create the message
	msg := &tokenizationtypes.MsgSetOutgoingApproval{
		Creator:      fromCosmosAddr,
		CollectionId: collectionId,
		Approval:     approval,
	}

	// Execute via the keeper
	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err := msgServer.SetOutgoingApproval(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("set outgoing approval failed: %w", err)
	}

	// Return success
	return method.Outputs.Pack(true)
}

// GetCollection queries a collection by ID
func (p Precompile) GetCollection(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments, expected 1, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Query the collection
	req := &tokenizationtypes.QueryGetCollectionRequest{
		CollectionId: collectionId.String(),
	}
	resp, err := p.tokenizationKeeper.GetCollection(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get collection failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal collection failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetBalance queries a balance for a user address
func (p Precompile) GetBalance(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments, expected 2, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}

	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid userAddress type")
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
		return nil, fmt.Errorf("get balance failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal balance failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetAddressList queries an address list by ID
func (p Precompile) GetAddressList(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments, expected 1, got %d", len(args))
	}

	listId, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid listId type")
	}

	// Query the address list
	req := &tokenizationtypes.QueryGetAddressListRequest{
		ListId: listId,
	}
	resp, err := p.tokenizationKeeper.GetAddressList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get address list failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal address list failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetApprovalTracker queries an approval tracker
func (p Precompile) GetApprovalTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 7 {
		return nil, fmt.Errorf("invalid number of arguments, expected 7, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approverAddress type")
	}
	amountTrackerId, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf("invalid amountTrackerId type")
	}
	trackerType, ok := args[4].(string)
	if !ok {
		return nil, fmt.Errorf("invalid trackerType type")
	}
	approvedAddress, ok := args[5].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approvedAddress type")
	}
	approvalId, ok := args[6].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalId type")
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
		return nil, fmt.Errorf("get approval tracker failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal approval tracker failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetChallengeTracker queries a challenge tracker
func (p Precompile) GetChallengeTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf("invalid number of arguments, expected 6, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approverAddress type")
	}
	challengeTrackerId, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf("invalid challengeTrackerId type")
	}
	leafIndexBig, ok := args[4].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid leafIndex type")
	}
	approvalId, ok := args[5].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalId type")
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
		return nil, fmt.Errorf("get challenge tracker failed: %w", err)
	}

	// Convert numUsed string to uint256
	numUsed := sdkmath.NewUintFromString(resp.NumUsed)

	return method.Outputs.Pack(numUsed.BigInt())
}

// GetETHSignatureTracker queries an ETH signature tracker
func (p Precompile) GetETHSignatureTracker(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf("invalid number of arguments, expected 6, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalId type")
	}
	challengeTrackerId, ok := args[4].(string)
	if !ok {
		return nil, fmt.Errorf("invalid challengeTrackerId type")
	}
	signature, ok := args[5].(string)
	if !ok {
		return nil, fmt.Errorf("invalid signature type")
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
		return nil, fmt.Errorf("get ETH signature tracker failed: %w", err)
	}

	// Convert numUsed string to uint256
	numUsed := sdkmath.NewUintFromString(resp.NumUsed)

	return method.Outputs.Pack(numUsed.BigInt())
}

// GetDynamicStore queries a dynamic store
func (p Precompile) GetDynamicStore(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments, expected 1, got %d", len(args))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid storeId type")
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)

	// Query the dynamic store
	req := &tokenizationtypes.QueryGetDynamicStoreRequest{
		StoreId: storeId.String(),
	}
	resp, err := p.tokenizationKeeper.GetDynamicStore(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get dynamic store failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal dynamic store failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetDynamicStoreValue queries a dynamic store value
func (p Precompile) GetDynamicStoreValue(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments, expected 2, got %d", len(args))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid storeId type")
	}
	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid userAddress type")
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
		return nil, fmt.Errorf("get dynamic store value failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal dynamic store value failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetWrappableBalances queries wrappable balances
func (p Precompile) GetWrappableBalances(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments, expected 2, got %d", len(args))
	}

	denom, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid denom type")
	}
	userAddress, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid userAddress type")
	}

	userCosmosAddr := sdk.AccAddress(userAddress.Bytes()).String()

	// Query wrappable balances
	req := &tokenizationtypes.QueryGetWrappableBalancesRequest{
		Denom:   denom,
		Address: userCosmosAddr,
	}
	resp, err := p.tokenizationKeeper.GetWrappableBalances(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get wrappable balances failed: %w", err)
	}

	return method.Outputs.Pack(resp.Amount.BigInt())
}

// IsAddressReservedProtocol checks if an address is reserved protocol
func (p Precompile) IsAddressReservedProtocol(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments, expected 1, got %d", len(args))
	}

	addr, ok := args[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid addr type")
	}

	cosmosAddr := sdk.AccAddress(addr.Bytes()).String()

	// Query if address is reserved protocol
	req := &tokenizationtypes.QueryIsAddressReservedProtocolRequest{
		Address: cosmosAddr,
	}
	resp, err := p.tokenizationKeeper.IsAddressReservedProtocol(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("is address reserved protocol failed: %w", err)
	}

	return method.Outputs.Pack(resp.IsReservedProtocol)
}

// GetAllReservedProtocolAddresses gets all reserved protocol addresses
func (p Precompile) GetAllReservedProtocolAddresses(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("invalid number of arguments, expected 0, got %d", len(args))
	}

	// Query all reserved protocol addresses
	req := &tokenizationtypes.QueryGetAllReservedProtocolAddressesRequest{}
	resp, err := p.tokenizationKeeper.GetAllReservedProtocolAddresses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get all reserved protocol addresses failed: %w", err)
	}

	// Convert Cosmos addresses to EVM addresses
	evmAddresses := make([]common.Address, len(resp.Addresses))
	for i, addrStr := range resp.Addresses {
		cosmosAddr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, fmt.Errorf("invalid address %s: %w", addrStr, err)
		}
		evmAddresses[i] = common.BytesToAddress(cosmosAddr.Bytes())
	}

	return method.Outputs.Pack(evmAddresses)
}

// GetVote queries a vote
func (p Precompile) GetVote(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf("invalid number of arguments, expected 6, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalId type")
	}
	proposalId, ok := args[4].(string)
	if !ok {
		return nil, fmt.Errorf("invalid proposalId type")
	}
	voterAddress, ok := args[5].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid voterAddress type")
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
		return nil, fmt.Errorf("get vote failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal vote failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// GetVotes queries all votes for a proposal
func (p Precompile) GetVotes(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("invalid number of arguments, expected 5, got %d", len(args))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid collectionId type")
	}
	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalLevel type")
	}
	approverAddress, ok := args[2].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid approverAddress type")
	}
	approvalId, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf("invalid approvalId type")
	}
	proposalId, ok := args[4].(string)
	if !ok {
		return nil, fmt.Errorf("invalid proposalId type")
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
		return nil, fmt.Errorf("get votes failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal votes failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}

// Params queries the module parameters
func (p Precompile) Params(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("invalid number of arguments, expected 0, got %d", len(args))
	}

	// Query params
	req := &tokenizationtypes.QueryParamsRequest{}
	resp, err := p.tokenizationKeeper.Params(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get params failed: %w", err)
	}

	// Marshal to bytes using types codec
	bz, err := tokenizationtypes.ModuleCdc.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal params failed: %w", err)
	}

	return method.Outputs.Pack(bz)
}
