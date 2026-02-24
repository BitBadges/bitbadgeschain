package keeper_test

import (
	"encoding/hex"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// mockEVMKeeperForQuery is a mock EVM keeper that returns configurable results
type mockEVMKeeperForQuery struct {
	contracts    map[string]bool
	returnValues map[string][]byte // contract address -> return value
	returnError  error
}

func (m *mockEVMKeeperForQuery) IsContract(ctx sdk.Context, addr common.Address) bool {
	accAddr := sdk.AccAddress(addr.Bytes())
	return m.contracts[accAddr.String()]
}

func (m *mockEVMKeeperForQuery) CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool, gasCap *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}

	accAddr := sdk.AccAddress(contract.Bytes())
	retVal := m.returnValues[accAddr.String()]
	return &evmtypes.MsgEthereumTxResponse{
		Ret: retVal,
	}, nil
}

// TestEVMQueryChallenges_NilEVMKeeper tests that EVM query fails gracefully when EVM keeper is nil
func (suite *TestSuite) TestEVMQueryChallenges_NilEVMKeeper() {
	// Set EVM keeper to nil to test error handling
	suite.app.TokenizationKeeper.SetEVMKeeper(nil)

	// Test ExecuteEVMQuery with nil keeper - should fail with "EVM keeper not available"
	calldata, _ := hex.DecodeString("70a08231")
	_, err := suite.app.TokenizationKeeper.ExecuteEVMQuery(
		suite.ctx,
		"0x0399ac65c88a0dcdd54bf5c0e8fc1e11bccccd39",
		calldata,
		100000,
	)
	suite.Require().NotNil(err, "EVM query should fail when EVM keeper is nil")
	suite.Require().Contains(err.Error(), "EVM keeper not available", "error should mention EVM keeper not available")
}

// NOTE: Full EVM query execution tests with mocked contracts are in
// x/tokenization/approval_criteria/evm_query_challenges_test.go
// The keeper integration tests focus on validation and edge cases.
// Full mock testing is challenging due to value receiver semantics in the keeper.

// TestEVMQueryChallenges_ValidationErrors tests validation of EVM query challenges during collection creation
func (suite *TestSuite) TestEVMQueryChallenges_ValidationErrors() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list
	err := suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
		ListId:    "testList",
		Addresses: []string{bob},
		Whitelist: true,
	})
	suite.Require().Nil(err, "error creating address list")

	tests := []struct {
		name        string
		challenges  []*types.EVMQueryChallenge
		expectError bool
		errorMsg    string
	}{
		{
			name: "missing contract address",
			challenges: []*types.EVMQueryChallenge{
				{
					ContractAddress: "",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectError: true,
			errorMsg:    "contract address required",
		},
		{
			name: "missing calldata",
			challenges: []*types.EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectError: true,
			errorMsg:    "calldata required",
		},
		{
			name: "invalid comparison operator",
			challenges: []*types.EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: "invalid",
					GasLimit:           sdkmath.NewUint(100000),
				},
			},
			expectError: true,
			errorMsg:    "invalid comparison operator",
		},
		{
			name: "gas limit exceeds maximum",
			challenges: []*types.EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(500001),
				},
			},
			expectError: true,
			errorMsg:    "gas limit exceeds maximum",
		},
	}

	for i, tt := range tests {
		suite.Run(tt.name, func() {
			collectionsToCreate := []*types.MsgNewCollection{
				{
					Creator: bob,
					CollectionApprovals: []*types.CollectionApproval{
						{
							ToListId:          "testList",
							FromListId:        "testList",
							InitiatedByListId: "testList",
							TransferTimes:     GetFullUintRanges(),
							OwnershipTimes:    GetFullUintRanges(),
							TokenIds:          GetFullUintRanges(),
							ApprovalId:        "test-validation-" + string(rune(i)),
							ApprovalCriteria: &types.ApprovalCriteria{
								EvmQueryChallenges: tt.challenges,
								MaxNumTransfers: &types.MaxNumTransfers{
									OverallMaxNumTransfers: sdkmath.NewUint(1000),
									AmountTrackerId:        "tracker-" + string(rune(i)),
								},
								ApprovalAmounts: &types.ApprovalAmounts{
									PerFromAddressApprovalAmount: sdkmath.NewUint(1000000),
									AmountTrackerId:              "tracker-" + string(rune(i)),
								},
							},
						},
					},
					TokensToCreate: []*types.Balance{
						{
							Amount:         sdkmath.NewUint(10),
							TokenIds:       GetFullUintRanges(),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					Permissions: &types.CollectionPermissions{
						CanArchiveCollection:         []*types.ActionPermission{},
						CanUpdateStandards:           []*types.ActionPermission{},
						CanUpdateCustomData:          []*types.ActionPermission{},
						CanDeleteCollection:          []*types.ActionPermission{},
						CanUpdateManager:             []*types.ActionPermission{},
						CanUpdateCollectionMetadata:  []*types.ActionPermission{},
						CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
						CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
						CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
							{
								PermanentlyPermittedTimes: GetFullUintRanges(),
							},
						},
					},
				},
			}

			err = CreateCollections(suite, wctx, collectionsToCreate)

			if tt.expectError {
				suite.Require().NotNil(err, "expected error for: %s", tt.name)
				suite.Require().Contains(err.Error(), tt.errorMsg, "error message should contain: %s", tt.errorMsg)
			} else {
				suite.Require().Nil(err, "expected no error for: %s", tt.name)
			}
		})
	}
}
