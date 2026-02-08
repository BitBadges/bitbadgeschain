package helpers

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/app/params"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

// init ensures SDK config is initialized with "bb" prefix before any tests run
// This must be called before any address operations to ensure correct Bech32 prefix
func init() {
	// Initialize SDK config with "bb" prefix if not already set
	// This is safe to call multiple times - it will only set if not already "bb"
	params.InitSDKConfigWithoutSeal()
}

// TestSuite provides common test utilities and fixtures
type TestSuite struct {
	Keeper      tokenizationkeeper.Keeper
	Ctx         sdk.Context
	Precompile  *tokenization.Precompile
	MsgServer   tokenizationtypes.MsgServer
	QueryClient tokenizationtypes.QueryClient

	// Test addresses (EVM format)
	AliceEVM   common.Address
	BobEVM     common.Address
	CharlieEVM common.Address
	ManagerEVM common.Address

	// Test addresses (Cosmos format)
	Alice   sdk.AccAddress
	Bob     sdk.AccAddress
	Charlie sdk.AccAddress
	Manager sdk.AccAddress

	// Test data
	CollectionId sdkmath.Uint
}

// NewTestSuite creates a new test suite with initialized keeper and context
func NewTestSuite() *TestSuite {
	// Ensure SDK config is initialized with "bb" prefix before any address operations
	// This must be called before creating addresses or calling keeper functions
	params.InitSDKConfigWithoutSeal()

	keeper, ctx := keepertest.TokenizationKeeper(nil)
	precompile := tokenization.NewPrecompile(keeper)
	msgServer := tokenizationkeeper.NewMsgServerImpl(keeper)

	// Create test EVM addresses
	aliceEVM := common.HexToAddress("0x1111111111111111111111111111111111111111")
	bobEVM := common.HexToAddress("0x2222222222222222222222222222222222222222")
	charlieEVM := common.HexToAddress("0x3333333333333333333333333333333333333333")
	managerEVM := common.HexToAddress("0x4444444444444444444444444444444444444444")

	// Convert to Cosmos addresses
	alice := sdk.AccAddress(aliceEVM.Bytes())
	bob := sdk.AccAddress(bobEVM.Bytes())
	charlie := sdk.AccAddress(charlieEVM.Bytes())
	manager := sdk.AccAddress(managerEVM.Bytes())

	return &TestSuite{
		Keeper:       keeper,
		Ctx:          ctx,
		Precompile:   precompile,
		MsgServer:    msgServer,
		AliceEVM:     aliceEVM,
		BobEVM:       bobEVM,
		CharlieEVM:   charlieEVM,
		ManagerEVM:   managerEVM,
		Alice:        alice,
		Bob:          bob,
		Charlie:      charlie,
		Manager:      manager,
		CollectionId: sdkmath.NewUint(0),
	}
}

// CreateMockContract creates a mock EVM contract for testing
func (ts *TestSuite) CreateMockContract(caller common.Address, input []byte) *vm.Contract {
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	if len(input) > 0 {
		contract.SetCallCode(common.Hash{}, input)
	}
	return contract
}

// CreateTestCollection creates a basic test collection with transfer approvals
func (ts *TestSuite) CreateTestCollection(creator string) (sdkmath.Uint, error) {
	validTokenIds := []*tokenizationtypes.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(100),
		},
	}

	msg := &tokenizationtypes.MsgCreateCollection{
		Creator:       creator,
		ValidTokenIds: validTokenIds,
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "test data",
		},
		Manager:               creator,
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
		IsArchived:            false,
	}

	resp, err := ts.MsgServer.CreateCollection(ts.Ctx, msg)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	collectionId := resp.CollectionId

	// Set up collection approvals to allow transfers
	// This is needed for transfers to work in tests
	getFullUintRanges := func() []*tokenizationtypes.UintRange {
		return []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		}
	}

	// Approval for regular transfers (AllWithoutMint to All)
	transferApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "transfer_approval",
		FromListId:        "AllWithoutMint",
		ToListId:          "All",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "transfer-tracker",
			},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "transfer-tracker",
			},
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   creator,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{transferApproval},
	}
	_, err = ts.MsgServer.UniversalUpdateCollection(ts.Ctx, updateMsg)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	// resp.CollectionId is already a Uint type
	ts.CollectionId = collectionId
	return collectionId, nil
}

// CreateTestBalance creates a test balance for a user
func (ts *TestSuite) CreateTestBalance(collectionId sdkmath.Uint, user string, amount sdkmath.Uint, tokenIds, ownershipTimes []*tokenizationtypes.UintRange) error {
	balance := &tokenizationtypes.Balance{
		Amount:         amount,
		TokenIds:       tokenIds,
		OwnershipTimes: ownershipTimes,
	}

	balances := []*tokenizationtypes.Balance{balance}

	store := &tokenizationtypes.UserBalanceStore{
		Balances: balances,
		AutoApproveSelfInitiatedOutgoingTransfers: true, // Allow self-initiated transfers for testing
		AutoApproveSelfInitiatedIncomingTransfers: true, // Allow self-initiated incoming transfers for testing
		AutoApproveAllIncomingTransfers:           true, // Allow all incoming transfers for testing
	}

	// Get collection
	collection, found := ts.Keeper.GetCollectionFromStore(ts.Ctx, collectionId)
	if !found {
		return sdkerrors.Wrapf(tokenizationtypes.ErrInvalidRequest, "collection %s not found", collectionId.String())
	}

	// Set balance using SetBalanceForAddress
	return ts.Keeper.SetBalanceForAddress(ts.Ctx, collection, user, store)
}

// CreateTestUintRange creates a UintRange for testing
func CreateTestUintRange(start, end uint64) *tokenizationtypes.UintRange {
	return &tokenizationtypes.UintRange{
		Start: sdkmath.NewUint(start),
		End:   sdkmath.NewUint(end),
	}
}

// CreateTestUintRangeArray creates an array of UintRanges
func CreateTestUintRangeArray(ranges [][2]uint64) []*tokenizationtypes.UintRange {
	result := make([]*tokenizationtypes.UintRange, len(ranges))
	for i, r := range ranges {
		result[i] = CreateTestUintRange(r[0], r[1])
	}
	return result
}

// CreateTestBalanceStruct creates a Balance struct for testing
func CreateTestBalanceStruct(amount uint64, tokenIds, ownershipTimes [][2]uint64) *tokenizationtypes.Balance {
	return &tokenizationtypes.Balance{
		Amount:         sdkmath.NewUint(amount),
		TokenIds:       CreateTestUintRangeArray(tokenIds),
		OwnershipTimes: CreateTestUintRangeArray(ownershipTimes),
	}
}

// BigIntToUint converts *big.Int to sdkmath.Uint
func BigIntToUint(bi *big.Int) sdkmath.Uint {
	return sdkmath.NewUintFromBigInt(bi)
}

// UintToBigInt converts sdkmath.Uint to *big.Int
func UintToBigInt(u sdkmath.Uint) *big.Int {
	return u.BigInt()
}

// CreateMockMethod creates a mock abi.Method for testing when the method doesn't exist in ABI
// This is a workaround for methods that haven't been added to the ABI JSON yet
func CreateMockMethod(name string, inputs, outputs abi.Arguments) abi.Method {
	// Create a mock method ID (first 4 bytes of keccak256 hash)
	// For testing, we use a simple hash of the name
	methodSig := fmt.Sprintf("%s(%s)", name, "")
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(methodSig))
	methodID := hash.Sum(nil)[:4]

	return abi.Method{
		Name:    name,
		RawName: name,
		Type:    abi.Function,
		Inputs:  inputs,
		Outputs: outputs,
		ID:      methodID,
	}
}

// CreateMockBoolOutput creates a simple abi.Arguments with a single bool output
func CreateMockBoolOutput() abi.Arguments {
	boolType, _ := abi.NewType("bool", "", nil)
	return abi.Arguments{
		{Type: boolType, Name: "success"},
	}
}
