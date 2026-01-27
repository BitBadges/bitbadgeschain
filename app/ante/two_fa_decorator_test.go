//go:build proto
// +build proto

// NOTE: These tests require proto generation to be run first:
//   ignite generate proto-go --yes
//   OR
//   make proto-gen
//
// The tests will not compile until types.User2FARequirements is generated from proto.

package ante

import (
	"math"
	"testing"
	"time"

	badgesmodulekeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	keepertest "github.com/bitbadges/bitbadgeschain/x/badges/testutil/keeper"
)

type TwoFADecoratorTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	badgesKeeper *badgesmodulekeeper.Keeper
	decorator    TwoFADecorator
	anteHandler  sdk.AnteHandler

	// Test addresses
	alice   sdk.AccAddress
	bob     sdk.AccAddress
	charlie sdk.AccAddress

	// Test collection IDs
	collectionId1 sdkmath.Uint
	collectionId2 sdkmath.Uint
}

func TestTwoFADecoratorTestSuite(t *testing.T) {
	suite.Run(t, new(TwoFADecoratorTestSuite))
}

func (suite *TwoFADecoratorTestSuite) SetupTest() {
	// Setup keeper
	k, ctx := keepertest.BadgesKeeper(suite.T())
	suite.badgesKeeper = &k
	suite.ctx = ctx.WithBlockTime(time.Now())

	// Setup test addresses
	privKey1 := secp256k1.GenPrivKey()
	privKey2 := secp256k1.GenPrivKey()
	privKey3 := secp256k1.GenPrivKey()
	suite.alice = sdk.AccAddress(privKey1.PubKey().Address())
	suite.bob = sdk.AccAddress(privKey2.PubKey().Address())
	suite.charlie = sdk.AccAddress(privKey3.PubKey().Address())

	// Setup decorator
	suite.decorator = NewTwoFADecorator(suite.badgesKeeper)

	// Setup ante handler chain (simplified for testing)
	// We only test the 2FA decorator, so we use a minimal chain
	// Note: We don't include SetUpContextDecorator as it requires GasTx interface
	// ChainAnteDecorators with a single decorator will have a no-op next handler
	suite.anteHandler = sdk.ChainAnteDecorators(
		suite.decorator, // Our 2FA decorator
	)

	// Create test collections
	suite.createTestCollections()
}

func (suite *TwoFADecoratorTestSuite) createTestCollections() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	msgServer := badgesmodulekeeper.NewMsgServerImpl(*suite.badgesKeeper)

	// Create collection 1
	msg1 := &types.MsgCreateCollection{
		Creator: suite.alice.String(),
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.alice.String(),
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "https://example.com/collection1",
		},
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}
	resp1, err := msgServer.CreateCollection(wctx, msg1)
	suite.Require().NoError(err)
	suite.collectionId1 = resp1.CollectionId

	// Create collection 2
	msg2 := &types.MsgCreateCollection{
		Creator: suite.bob.String(),
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.bob.String(),
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "https://example.com/collection2",
		},
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}
	resp2, err := msgServer.CreateCollection(wctx, msg2)
	suite.Require().NoError(err)
	suite.collectionId2 = resp2.CollectionId
}

func (suite *TwoFADecoratorTestSuite) mintBadges(collectionId sdkmath.Uint, to sdk.AccAddress, amount sdkmath.Uint, tokenIds []*types.UintRange) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	msgServer := badgesmodulekeeper.NewMsgServerImpl(*suite.badgesKeeper)

	// Setup mint approval using UniversalUpdateCollection
	collection, _ := suite.badgesKeeper.GetCollectionFromStore(suite.ctx, collectionId)

	// Check if mint approval already exists
	hasMintApproval := false
	for _, approval := range collection.CollectionApprovals {
		if approval.ApprovalId == "mint" && approval.FromListId == types.MintAddress {
			hasMintApproval = true
			break
		}
	}

	if !hasMintApproval {
		mintApproval := &types.CollectionApproval{
			ApprovalId:        "mint",
			FromListId:        types.MintAddress,
			ToListId:          "All",
			InitiatedByListId: "All",
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
			TransferTimes:  []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		}

		// Prepend mint approval to existing approvals
		newApprovals := []*types.CollectionApproval{mintApproval}
		newApprovals = append(newApprovals, collection.CollectionApprovals...)

		updateMsg := &types.MsgUniversalUpdateCollection{
			Creator:                   collection.Manager,
			CollectionId:              collectionId,
			UpdateCollectionApprovals: true,
			CollectionApprovals:       newApprovals,
		}
		_, err := msgServer.UniversalUpdateCollection(wctx, updateMsg)
		suite.Require().NoError(err, "failed to set up mint approval")
	}

	// Mint badges
	msg := &types.MsgTransferTokens{
		Creator:      suite.alice.String(),
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{to.String()},
				Balances: []*types.Balance{
					{
						Amount:         amount,
						TokenIds:       tokenIds,
						OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
					},
				},
			},
		},
	}
	_, err := msgServer.TransferTokens(wctx, msg)
	suite.Require().NoError(err)
}

func (suite *TwoFADecoratorTestSuite) set2FARequirements(address sdk.AccAddress, requirements []*types.MustOwnTokens) {
	// Directly set in store for testing
	// NOTE: This requires proto generation. Run: ignite generate proto-go --yes
	// After generation, types.User2FARequirements will be available from proto/badges/user_2fa_requirements.proto
	user2FARequirements := &types.User2FARequirements{
		MustOwnTokens:          requirements,
		DynamicStoreChallenges: []*types.DynamicStoreChallenge{},
	}
	err := suite.badgesKeeper.SetUser2FARequirementsInStore(suite.ctx, address.String(), user2FARequirements)
	suite.Require().NoError(err)
}

func (suite *TwoFADecoratorTestSuite) set2FARequirementsWithDynamicStores(address sdk.AccAddress, mustOwnTokens []*types.MustOwnTokens, dynamicStoreChallenges []*types.DynamicStoreChallenge) {
	user2FARequirements := &types.User2FARequirements{
		MustOwnTokens:          mustOwnTokens,
		DynamicStoreChallenges: dynamicStoreChallenges,
	}
	err := suite.badgesKeeper.SetUser2FARequirementsInStore(suite.ctx, address.String(), user2FARequirements)
	suite.Require().NoError(err)
}

// createDynamicStore creates a dynamic store and returns its storeId
func (suite *TwoFADecoratorTestSuite) createDynamicStore(creator sdk.AccAddress, defaultValue bool) sdkmath.Uint {
	wctx := sdk.WrapSDKContext(suite.ctx)
	msgServer := badgesmodulekeeper.NewMsgServerImpl(*suite.badgesKeeper)

	msg := &types.MsgCreateDynamicStore{
		Creator:      creator.String(),
		DefaultValue: defaultValue,
	}

	resp, err := msgServer.CreateDynamicStore(wctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	return resp.StoreId
}

// setDynamicStoreValue sets a value in a dynamic store for a specific address
func (suite *TwoFADecoratorTestSuite) setDynamicStoreValue(creator sdk.AccAddress, storeId sdkmath.Uint, address sdk.AccAddress, value bool) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	msgServer := badgesmodulekeeper.NewMsgServerImpl(*suite.badgesKeeper)

	msg := &types.MsgSetDynamicStoreValue{
		Creator: creator.String(),
		StoreId: storeId,
		Address: address.String(),
		Value:   value,
	}

	_, err := msgServer.SetDynamicStoreValue(wctx, msg)
	suite.Require().NoError(err)
}

// updateDynamicStoreGlobalEnabled updates the global enabled status of a dynamic store
func (suite *TwoFADecoratorTestSuite) updateDynamicStoreGlobalEnabled(creator sdk.AccAddress, storeId sdkmath.Uint, globalEnabled bool) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	msgServer := badgesmodulekeeper.NewMsgServerImpl(*suite.badgesKeeper)

	msg := &types.MsgUpdateDynamicStore{
		Creator:       creator.String(),
		StoreId:       storeId,
		GlobalEnabled: globalEnabled,
	}

	_, err := msgServer.UpdateDynamicStore(wctx, msg)
	suite.Require().NoError(err)
}

func (suite *TwoFADecoratorTestSuite) createBankMsgSend(from, to sdk.AccAddress, amount sdk.Coins) sdk.Msg {
	return &banktypes.MsgSend{
		FromAddress: from.String(),
		ToAddress:   to.String(),
		Amount:      amount,
	}
}

func (suite *TwoFADecoratorTestSuite) createBadgesMsgTransfer(from, to sdk.AccAddress, collectionId sdkmath.Uint) sdk.Msg {
	return &types.MsgTransferTokens{
		Creator:      from.String(),
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        from.String(),
				ToAddresses: []string{to.String()},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
					},
				},
			},
		},
	}
}

// MockTx is a simple mock transaction for testing
type MockTx struct {
	msgs []sdk.Msg
}

func (m *MockTx) GetMsgs() []sdk.Msg {
	return m.msgs
}

func (m *MockTx) GetMsgsV2() ([]proto.Message, error) {
	result := make([]proto.Message, len(m.msgs))
	for i, msg := range m.msgs {
		if protoMsg, ok := msg.(proto.Message); ok {
			result[i] = protoMsg
		}
	}
	return result, nil
}

func (suite *TwoFADecoratorTestSuite) createMockTx(msgs []sdk.Msg) sdk.Tx {
	return &MockTx{msgs: msgs}
}

// TestBasic2FACheck tests basic 2FA requirement checking
func (suite *TwoFADecoratorTestSuite) TestBasic2FACheck() {
	// Mint badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements for alice
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because alice owns the required badges
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when 2FA requirements are met")
}

// Test2FAFailure tests that transactions fail when 2FA requirements are not met
func (suite *TwoFADecoratorTestSuite) Test2FAFailure() {
	// Set 2FA requirements for alice (but don't mint badges)
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Verify requirements are stored
	requirements, found := suite.badgesKeeper.GetUser2FARequirementsFromStore(suite.ctx, suite.alice.String())
	suite.Require().True(found, "2FA requirements should be stored")
	suite.Require().Len(requirements.MustOwnTokens, 1, "should have 1 requirement")

	// Verify alice has no badges (check balance)
	collection, _ := suite.badgesKeeper.GetCollectionFromStore(suite.ctx, suite.collectionId1)
	balance, _, err := suite.badgesKeeper.GetBalanceOrApplyDefault(suite.ctx, collection, suite.alice.String())
	suite.Require().NoError(err)
	suite.Require().Len(balance.Balances, 0, "alice should have no badges before the test")

	// Create a bank transfer message from alice
	msg := suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())
	tx := suite.createMockTx([]sdk.Msg{msg})

	// Verify signers can be extracted from the message (for bank messages, extract from FromAddress)
	if bankMsg, ok := msg.(*banktypes.MsgSend); ok {
		suite.Require().Equal(suite.alice.String(), bankMsg.FromAddress, "FromAddress should be alice")
	} else {
		// For other message types, try GetSigners()
		type msgWithSigners interface {
			GetSigners() []sdk.AccAddress
		}
		if msgWithSigs, ok := msg.(msgWithSigners); ok {
			signers := msgWithSigs.GetSigners()
			suite.Require().Len(signers, 1, "message should have 1 signer")
			suite.Require().Equal(suite.alice.String(), signers[0].String(), "signer should be alice")
		} else {
			suite.T().Fatalf("message does not implement GetSigners() and is not a known type")
		}
	}

	// Should fail because alice doesn't own the required badges
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when 2FA requirements are not met")
	suite.Require().Contains(err.Error(), "2FA requirement failed", "error should mention 2FA requirement")
}

// TestNo2FARequirements tests that transactions pass when no 2FA requirements are set
func (suite *TwoFADecoratorTestSuite) TestNo2FARequirements() {
	// Don't set any 2FA requirements

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because no 2FA requirements are set
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when no 2FA requirements are set")
}

// TestEmpty2FARequirements tests that empty 2FA requirements are treated as no requirements
func (suite *TwoFADecoratorTestSuite) TestEmpty2FARequirements() {
	// Set empty 2FA requirements
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because empty requirements are treated as no requirements
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when 2FA requirements are empty")
}

// TestMultiple2FARequirements tests multiple 2FA requirements
func (suite *TwoFADecoratorTestSuite) TestMultiple2FARequirements() {
	// Mint badges to alice in both collections
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
	suite.mintBadges(suite.collectionId2, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set multiple 2FA requirements
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
		{
			CollectionId: suite.collectionId2,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because alice owns badges in both collections
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when all 2FA requirements are met")
}

// TestMultiple2FARequirementsFailure tests that transactions fail if any requirement is not met
func (suite *TwoFADecoratorTestSuite) TestMultiple2FARequirementsFailure() {
	// Mint badges to alice in only one collection
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set multiple 2FA requirements
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
		{
			CollectionId: suite.collectionId2,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because alice doesn't own badges in collection 2
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when any 2FA requirement is not met")
}

// TestTimeDependent2FA tests time-dependent 2FA requirements
func (suite *TwoFADecoratorTestSuite) TestTimeDependent2FA() {
	// Set block time to a specific time
	blockTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	// Mint badges with specific ownership times
	ownershipStart := sdkmath.NewUint(uint64(blockTime.UnixMilli()))
	ownershipEnd := sdkmath.NewUint(uint64(blockTime.Add(24 * time.Hour).UnixMilli()))

	balance := &types.UserBalanceStore{
		Balances: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: []*types.UintRange{{Start: ownershipStart, End: ownershipEnd}},
			},
		},
	}
	balanceKey := badgesmodulekeeper.ConstructBalanceKey(suite.alice.String(), suite.collectionId1)
	suite.badgesKeeper.SetUserBalanceInStore(suite.ctx, balanceKey, balance, true)

	// Set 2FA requirements with matching ownership times
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: ownershipStart, End: ownershipEnd}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because ownership times match
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when ownership times match")
}

// TestTimeDependent2FAFailure tests that transactions fail when ownership times don't match
func (suite *TwoFADecoratorTestSuite) TestTimeDependent2FAFailure() {
	// Set block time to a specific time
	blockTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	// Mint badges with specific ownership times
	ownershipStart := sdkmath.NewUint(uint64(blockTime.Add(-48 * time.Hour).UnixMilli()))
	ownershipEnd := sdkmath.NewUint(uint64(blockTime.Add(-24 * time.Hour).UnixMilli()))

	balance := &types.UserBalanceStore{
		Balances: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: []*types.UintRange{{Start: ownershipStart, End: ownershipEnd}},
			},
		},
	}
	balanceKey := badgesmodulekeeper.ConstructBalanceKey(suite.alice.String(), suite.collectionId1)
	suite.badgesKeeper.SetUserBalanceInStore(suite.ctx, balanceKey, balance, true)

	// Set 2FA requirements with different ownership times (current time)
	currTime := sdkmath.NewUint(uint64(blockTime.UnixMilli()))
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: currTime, End: currTime}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because ownership times don't match
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when ownership times don't match")
}

// TestOverrideWithCurrentTime tests the OverrideWithCurrentTime feature
func (suite *TwoFADecoratorTestSuite) TestOverrideWithCurrentTime() {
	// Mint badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements with OverrideWithCurrentTime
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:            suite.collectionId1,
			AmountRange:             &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:                []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes:          []*types.UintRange{}, // Empty, will be overridden
			OverrideWithCurrentTime: true,
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because OverrideWithCurrentTime uses current block time
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when OverrideWithCurrentTime is used")
}

// TestMultipleSigners tests transactions with multiple signers
func (suite *TwoFADecoratorTestSuite) TestMultipleSigners() {
	// Mint badges to both alice and bob
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
	suite.mintBadges(suite.collectionId1, suite.bob, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements for both
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})
	suite.set2FARequirements(suite.bob, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a transaction with messages from both alice and bob
	tx := suite.createMockTx([]sdk.Msg{
		suite.createBankMsgSend(suite.alice, suite.charlie, sdk.NewCoins()),
		suite.createBankMsgSend(suite.bob, suite.charlie, sdk.NewCoins()),
	})

	// Should pass because both signers meet their 2FA requirements
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when all signers meet 2FA requirements")
}

// TestMultipleSignersFailure tests that transactions fail if any signer doesn't meet requirements
func (suite *TwoFADecoratorTestSuite) TestMultipleSignersFailure() {
	// Mint badges only to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements for both
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})
	suite.set2FARequirements(suite.bob, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a transaction with messages from both alice and bob
	tx := suite.createMockTx([]sdk.Msg{
		suite.createBankMsgSend(suite.alice, suite.charlie, sdk.NewCoins()),
		suite.createBankMsgSend(suite.bob, suite.charlie, sdk.NewCoins()),
	})

	// Should fail because bob doesn't meet 2FA requirements
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when any signer doesn't meet 2FA requirements")
}

// TestDifferentTransactionTypes tests that 2FA works with different transaction types
func (suite *TwoFADecoratorTestSuite) TestDifferentTransactionTypes() {
	// Mint badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	testCases := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			name: "Bank transfer",
			msg:  suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins()),
		},
		{
			name: "Badges transfer",
			msg:  suite.createBadgesMsgTransfer(suite.alice, suite.bob, suite.collectionId1),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tx := suite.createMockTx([]sdk.Msg{tc.msg})
			_, err := suite.anteHandler(suite.ctx, tx, false)
			suite.Require().NoError(err, "%s should pass when 2FA requirements are met", tc.name)
		})
	}
}

// TestAmountRange tests amount range validation
func (suite *TwoFADecoratorTestSuite) TestAmountRange() {
	// Mint 5 badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(5), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)}})

	// Set 2FA requirements with amount range 1-10
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because amount (5) is within range (1-10)
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when amount is within range")
}

// TestAmountRangeFailure tests that transactions fail when amount is outside range
func (suite *TwoFADecoratorTestSuite) TestAmountRangeFailure() {
	// Mint 20 badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(20), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(20)}})

	// Set 2FA requirements with amount range 1-10
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(20)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because amount (20) is outside range (1-10)
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when amount is outside range")
}

// TestMustSatisfyForAllAssets tests the MustSatisfyForAllAssets flag
func (suite *TwoFADecoratorTestSuite) TestMustSatisfyForAllAssets() {
	// Mint badges to alice - some within range, some outside
	balance := &types.UserBalanceStore{
		Balances: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(5), // Within range
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			},
			{
				Amount:         sdkmath.NewUint(20), // Outside range
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			},
		},
	}
	balanceKey := badgesmodulekeeper.ConstructBalanceKey(suite.alice.String(), suite.collectionId1)
	suite.badgesKeeper.SetUserBalanceInStore(suite.ctx, balanceKey, balance, true)

	// Set 2FA requirements with MustSatisfyForAllAssets = false (any one passes)
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TokenIds:                []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
			OwnershipTimes:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			MustSatisfyForAllAssets: false,
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because at least one balance (token 1 with amount 5) is within range
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when MustSatisfyForAllAssets is false and at least one balance passes")
}

// TestMustSatisfyForAllAssetsFailure tests that transactions fail when MustSatisfyForAllAssets is true and not all pass
func (suite *TwoFADecoratorTestSuite) TestMustSatisfyForAllAssetsFailure() {
	// Mint badges to alice - some within range, some outside
	balance := &types.UserBalanceStore{
		Balances: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(5), // Within range
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			},
			{
				Amount:         sdkmath.NewUint(20), // Outside range
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
				OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			},
		},
	}
	balanceKey := badgesmodulekeeper.ConstructBalanceKey(suite.alice.String(), suite.collectionId1)
	suite.badgesKeeper.SetUserBalanceInStore(suite.ctx, balanceKey, balance, true)

	// Set 2FA requirements with MustSatisfyForAllAssets = true (all must pass)
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TokenIds:                []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
			OwnershipTimes:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			MustSatisfyForAllAssets: true,
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because not all balances are within range (token 2 has amount 20)
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when MustSatisfyForAllAssets is true and not all balances pass")
}

// TestInvalidCollection tests that transactions fail when collection doesn't exist
func (suite *TwoFADecoratorTestSuite) TestInvalidCollection() {
	// Set 2FA requirements with non-existent collection
	invalidCollectionId := sdkmath.NewUint(99999)
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:   invalidCollectionId,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because collection doesn't exist
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when collection doesn't exist")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSimulateMode tests that simulate mode bypasses 2FA checks
func (suite *TwoFADecoratorTestSuite) TestSimulateMode() {
	// Set 2FA requirements but don't mint badges
	suite.set2FARequirements(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId:   suite.collectionId1,
			AmountRange:    &types.UintRange{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// In simulate mode, should still check 2FA (we don't bypass it)
	// But let's verify the behavior is consistent
	_, err := suite.anteHandler(suite.ctx, tx, true)
	// The decorator doesn't bypass simulate mode, so it should still fail
	suite.Require().Error(err, "transaction should fail in simulate mode if 2FA requirements are not met")
}

// ========== Dynamic Store 2FA Tests ==========

// TestDynamicStore2FABasicSuccess tests that transactions pass when dynamic store 2FA requirements are met
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FABasicSuccess() {
	// Create a dynamic store with defaultValue = false
	storeId := suite.createDynamicStore(suite.alice, false)

	// Set value to true for alice
	suite.setDynamicStoreValue(suite.alice, storeId, suite.alice, true)

	// Set 2FA requirements with dynamic store challenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator", // Will check the signer (alice)
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because alice has true value in the dynamic store
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when dynamic store 2FA requirements are met")
}

// TestDynamicStore2FABasicFailure tests that transactions fail when dynamic store 2FA requirements are not met
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FABasicFailure() {
	// Create a dynamic store with defaultValue = false
	storeId := suite.createDynamicStore(suite.alice, false)

	// Don't set value for alice (will use defaultValue = false)

	// Set 2FA requirements with dynamic store challenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because alice doesn't have true value in the dynamic store
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when dynamic store 2FA requirements are not met")
	suite.Require().Contains(err.Error(), "2FA dynamic store challenge failed", "error should mention dynamic store challenge")
}

// TestDynamicStore2FAWithDefaultValueTrue tests that defaultValue = true works correctly
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithDefaultValueTrue() {
	// Create a dynamic store with defaultValue = true
	storeId := suite.createDynamicStore(suite.alice, true)

	// Don't set explicit value for alice (will use defaultValue = true)

	// Set 2FA requirements with dynamic store challenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because defaultValue is true
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when dynamic store has defaultValue = true")
}

// TestDynamicStore2FAWithGlobalKillSwitch tests that globally disabled stores fail immediately
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithGlobalKillSwitch() {
	// Create a dynamic store with defaultValue = true
	storeId := suite.createDynamicStore(suite.alice, true)

	// Disable the store globally
	suite.updateDynamicStoreGlobalEnabled(suite.alice, storeId, false)

	// Set 2FA requirements with dynamic store challenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because the store is globally disabled
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when dynamic store is globally disabled")
	suite.Require().Contains(err.Error(), "globally disabled", "error should mention global disable")
}

// TestDynamicStore2FAWithInvalidStoreId tests that transactions fail when store doesn't exist
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithInvalidStoreId() {
	// Use a non-existent store ID
	invalidStoreId := sdkmath.NewUint(99999)

	// Set 2FA requirements with invalid dynamic store challenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             invalidStoreId,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should fail because store doesn't exist
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when dynamic store doesn't exist")
	suite.Require().Contains(err.Error(), "dynamic store", "error should mention dynamic store")
}

// TestDynamicStore2FAWithMultipleChallenges tests multiple dynamic store challenges
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithMultipleChallenges() {
	// Create two dynamic stores
	storeId1 := suite.createDynamicStore(suite.alice, false)
	storeId2 := suite.createDynamicStore(suite.alice, false)

	// Set value to true for alice in both stores
	suite.setDynamicStoreValue(suite.alice, storeId1, suite.alice, true)
	suite.setDynamicStoreValue(suite.alice, storeId2, suite.alice, true)

	// Set 2FA requirements with multiple dynamic store challenges
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId1,
			OwnershipCheckParty: "initiator",
		},
		{
			StoreId:             storeId2,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because alice has true value in both stores
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when all dynamic store challenges are met")

	// Now set one store to false
	suite.setDynamicStoreValue(suite.alice, storeId1, suite.alice, false)

	// Should fail because one store is false
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when any dynamic store challenge is not met")
}

// TestDynamicStore2FAWithOwnershipCheckParty tests different ownershipCheckParty values
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithOwnershipCheckParty() {
	// Create a dynamic store
	storeId := suite.createDynamicStore(suite.alice, false)

	// Set value to true for alice
	suite.setDynamicStoreValue(suite.alice, storeId, suite.alice, true)

	// Test with "initiator" (default)
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator",
		},
	})

	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass with 'initiator' ownershipCheckParty")

	// Test with "sender" (should also check signer for 2FA)
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "sender",
		},
	})

	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass with 'sender' ownershipCheckParty")

	// Test with empty string (should default to signer)
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "",
		},
	})

	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass with empty ownershipCheckParty")
}

// TestDynamicStore2FAWithCustomAddress tests checking a custom address
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithCustomAddress() {
	// Create a dynamic store
	storeId := suite.createDynamicStore(suite.alice, false)

	// Set value to true for bob (not alice)
	suite.setDynamicStoreValue(suite.alice, storeId, suite.bob, true)

	// Set 2FA requirements with custom address (bob)
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: suite.bob.String(), // Check bob's value
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because bob has true value (even though alice is the signer)
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when custom address has true value")

	// Now set bob's value to false
	suite.setDynamicStoreValue(suite.alice, storeId, suite.bob, false)

	// Should fail because bob's value is now false
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when custom address has false value")
}

// TestDynamicStore2FACombinedWithMustOwnTokens tests combining dynamic store and badge ownership 2FA
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FACombinedWithMustOwnTokens() {
	// Create a dynamic store
	storeId := suite.createDynamicStore(suite.alice, false)
	suite.setDynamicStoreValue(suite.alice, storeId, suite.alice, true)

	// Mint badges to alice
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})

	// Set 2FA requirements with both MustOwnTokens and DynamicStoreChallenge
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{
		{
			CollectionId: suite.collectionId1,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
		},
	}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because both requirements are met
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when both badge and dynamic store requirements are met")

	// Remove alice's badge
	balanceKey := badgesmodulekeeper.ConstructBalanceKey(suite.alice.String(), suite.collectionId1)
	suite.badgesKeeper.DeleteUserBalanceFromStore(suite.ctx, balanceKey)

	// Should fail because badge requirement is not met
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when badge requirement is not met")

	// Restore badge but set dynamic store to false
	suite.mintBadges(suite.collectionId1, suite.alice, sdkmath.NewUint(1), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
	suite.setDynamicStoreValue(suite.alice, storeId, suite.alice, false)

	// Should fail because dynamic store requirement is not met
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when dynamic store requirement is not met")
}

// TestDynamicStore2FAWithMultipleSigners tests dynamic store 2FA with multiple signers
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAWithMultipleSigners() {
	// Create dynamic stores for alice and bob
	storeId1 := suite.createDynamicStore(suite.alice, false)
	storeId2 := suite.createDynamicStore(suite.bob, false)

	// Verify stores were created correctly and have different IDs
	suite.Require().NotEqual(storeId1, storeId2, "storeId1 and storeId2 should be different")
	
	store1, found1 := suite.badgesKeeper.GetDynamicStoreFromStore(suite.ctx, storeId1)
	suite.Require().True(found1, "storeId1 should exist")
	suite.Require().Equal(suite.alice.String(), store1.CreatedBy, "storeId1 should be created by alice, got %s", store1.CreatedBy)

	store2, found2 := suite.badgesKeeper.GetDynamicStoreFromStore(suite.ctx, storeId2)
	suite.Require().True(found2, "storeId2 should exist")
	suite.Require().Equal(suite.bob.String(), store2.CreatedBy, "storeId2 should be created by bob, got %s", store2.CreatedBy)

	// Set values to true
	suite.setDynamicStoreValue(suite.alice, storeId1, suite.alice, true)
	suite.setDynamicStoreValue(suite.bob, storeId2, suite.bob, true)

	// Set 2FA requirements for alice
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId1,
			OwnershipCheckParty: "initiator",
		},
	})

	// Set 2FA requirements for bob
	suite.set2FARequirementsWithDynamicStores(suite.bob, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{
		{
			StoreId:             storeId2,
			OwnershipCheckParty: "initiator",
		},
	})

	// Create a transaction with both alice and bob as signers
	tx := suite.createMockTx([]sdk.Msg{
		suite.createBankMsgSend(suite.alice, suite.charlie, sdk.NewCoins()),
		suite.createBankMsgSend(suite.bob, suite.charlie, sdk.NewCoins()),
	})

	// Should pass because both signers meet their 2FA requirements
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when all signers meet their 2FA requirements")

	// Set bob's value to false
	suite.setDynamicStoreValue(suite.bob, storeId2, suite.bob, false)

	// Should fail because bob doesn't meet his 2FA requirement
	_, err = suite.anteHandler(suite.ctx, tx, false)
	suite.Require().Error(err, "transaction should fail when any signer doesn't meet 2FA requirements")
	suite.Require().Contains(err.Error(), suite.bob.String(), "error should mention the failing signer")
}

// TestDynamicStore2FAEmptyChallenges tests that empty dynamic store challenges are treated as no requirements
func (suite *TwoFADecoratorTestSuite) TestDynamicStore2FAEmptyChallenges() {
	// Set 2FA requirements with empty dynamic store challenges
	suite.set2FARequirementsWithDynamicStores(suite.alice, []*types.MustOwnTokens{}, []*types.DynamicStoreChallenge{})

	// Create a bank transfer message from alice
	tx := suite.createMockTx([]sdk.Msg{suite.createBankMsgSend(suite.alice, suite.bob, sdk.NewCoins())})

	// Should pass because no dynamic store requirements are set
	_, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "transaction should pass when dynamic store challenges are empty")
}
