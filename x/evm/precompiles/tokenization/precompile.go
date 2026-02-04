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
	// GasTransferTokens defines the gas cost for transferTokens
	GasTransferTokens = 50_000
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
	case TransferTokensMethod:
		bz, err = p.TransferTokens(ctx, method, args, contract)
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, err
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
func (Precompile) IsTransaction(method *abi.Method) bool {
	switch method.Name {
	case TransferTokensMethod:
		return true
	default:
		return false
	}
}

// TransferTokensMethod is the method name for transferTokens
const TransferTokensMethod = "transferTokens"

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
