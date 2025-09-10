package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUnwrapIBCDenom(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx.WithChainID("test-chain-id")
	wctx := sdk.WrapSDKContext(ctx)

	// Create a test collection with IBC unwrap paths
	collectionId := sdkmath.NewUint(1)
	creator := suite.app.AccountKeeper.GetModuleAddress("badges").String()
	if creator == "" {
		creator = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q" // fallback to a valid address
	}

	// Create collection with IBC unwrap paths
	collection := &types.BadgeCollection{
		CollectionId: collectionId,
		BalancesType: "Standard",
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1000),
					BadgeIds: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
					},
					OwnershipTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
					},
				},
			},
		},
		ValidBadgeIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
		IbcUnwrapPaths: []*types.IBCUnwrapPath{
			{
				ChannelId:                      "channel-0",
				PortId:                         "transfer",
				SourceCollectionId:             sdkmath.NewUint(20),
				Denom:                          "test-denom",
				AllowOverrideWithAnyValidToken: false,
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(100),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
				DenomSuffixDetails: &types.DenomSuffixDetails{
					WithAddress: false,
				},
			},
		},
	}

	// Store the collection
	err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection)
	require.NoError(t, err)

	// Test 1: Basic unwrap without suffixes
	t.Run("Basic unwrap without suffixes", func(t *testing.T) {
		// Calculate expected IBC denom hash
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom", collection.IbcUnwrapPaths[0].SourceCollectionId.String())
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Mint the IBC tokens to the creator first (simulating IBC transfer)
		creatorAddr, err := sdk.AccAddressFromBech32(creator)
		require.NoError(t, err)

		ibcCoin := sdk.NewCoin(ibcDenom, sdkmath.NewInt(1000000))
		err = suite.app.BankKeeper.MintCoins(ctx, "badges", sdk.NewCoins(ibcCoin))
		require.NoError(t, err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "badges", creatorAddr, sdk.NewCoins(ibcCoin))
		require.NoError(t, err)

		// Create message
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.NoError(t, err)

		// Check that badges were minted
		userBalance, err := GetUserBalance(suite, wctx, collectionId, creator)
		require.NoError(t, err)
		require.Len(t, userBalance.Balances, 1)
		require.Equal(t, sdkmath.NewUint(100001000), userBalance.Balances[0].Amount) // 1000 (default) + 100 * 1000000
	})

	// Test 2: Unwrap with address suffix
	t.Run("Unwrap with address suffix", func(t *testing.T) {
		// Create a new collection with address suffix
		collectionId2 := sdkmath.NewUint(2)
		collection2 := &types.BadgeCollection{
			CollectionId: collectionId2,
			BalancesType: "Standard",
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
			},
			ValidBadgeIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			IbcUnwrapPaths: []*types.IBCUnwrapPath{
				{
					ChannelId:                      "channel-0",
					PortId:                         "transfer",
					SourceCollectionId:             sdkmath.NewUint(20),
					Denom:                          "test-denom",
					AllowOverrideWithAnyValidToken: false,
					Balances: []*types.Balance{
						{
							Amount: sdkmath.NewUint(100),
							BadgeIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
					DenomSuffixDetails: &types.DenomSuffixDetails{
						WithAddress: true,
					},
				},
			},
		}

		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection2)
		require.NoError(t, err)

		// Calculate expected IBC denom hash with address suffix
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom:%s", collection2.IbcUnwrapPaths[0].SourceCollectionId.String(), creator)
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Mint the IBC tokens to the creator first (simulating IBC transfer)
		creatorAddr, err := sdk.AccAddressFromBech32(creator)
		require.NoError(t, err)

		ibcCoin := sdk.NewCoin(ibcDenom, sdkmath.NewInt(1000000))
		err = suite.app.BankKeeper.MintCoins(ctx, "badges", sdk.NewCoins(ibcCoin))
		require.NoError(t, err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "badges", creatorAddr, sdk.NewCoins(ibcCoin))
		require.NoError(t, err)

		// Create message
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId2,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.NoError(t, err)

		// Check that badges were minted
		userBalance, err := GetUserBalance(suite, wctx, collectionId2, creator)
		require.NoError(t, err)
		require.Len(t, userBalance.Balances, 1)
		require.Equal(t, sdkmath.NewUint(100001000), userBalance.Balances[0].Amount) // 1000 (default) + 100 * 1000000
	})

	// Test 3: Unwrap with destination collection ID and chain ID
	t.Run("Unwrap with destination collection ID and chain ID", func(t *testing.T) {
		// Create a new collection with destination gating
		collectionId3 := sdkmath.NewUint(3)
		collection3 := &types.BadgeCollection{
			CollectionId: collectionId3,
			BalancesType: "Standard",
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
			},
			ValidBadgeIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			IbcUnwrapPaths: []*types.IBCUnwrapPath{
				{
					ChannelId:                      "channel-0",
					PortId:                         "transfer",
					SourceCollectionId:             sdkmath.NewUint(20),
					Denom:                          "test-denom",
					AllowOverrideWithAnyValidToken: false,
					Balances: []*types.Balance{
						{
							Amount: sdkmath.NewUint(100),
							BadgeIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
					DenomSuffixDetails: &types.DenomSuffixDetails{
						WithAddress:             true,
						DestinationCollectionId: collectionId3.String(),
						DestinationChainId:      "test-chain-id",
					},
				},
			},
		}

		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection3)
		require.NoError(t, err)

		// Calculate expected IBC denom hash with all suffixes
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom:%s:%s:%s",
			collection3.IbcUnwrapPaths[0].SourceCollectionId.String(),
			creator,
			collectionId3.String(),
			"test-chain-id")
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Mint the IBC tokens to the creator first (simulating IBC transfer)
		creatorAddr, err := sdk.AccAddressFromBech32(creator)
		require.NoError(t, err)

		ibcCoin := sdk.NewCoin(ibcDenom, sdkmath.NewInt(1000000))
		err = suite.app.BankKeeper.MintCoins(ctx, "badges", sdk.NewCoins(ibcCoin))
		require.NoError(t, err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "badges", creatorAddr, sdk.NewCoins(ibcCoin))
		require.NoError(t, err)

		// Create message
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId3,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.NoError(t, err)

		// Check that badges were minted
		userBalance, err := GetUserBalance(suite, wctx, collectionId3, creator)
		require.NoError(t, err)
		require.Len(t, userBalance.Balances, 1)
		require.Equal(t, sdkmath.NewUint(100001000), userBalance.Balances[0].Amount) // 1000 (default) + 100 * 1000000
	})

	// Test 4: Unwrap with override token ID
	t.Run("Unwrap with override token ID", func(t *testing.T) {
		// Create a new collection with override support
		collectionId4 := sdkmath.NewUint(4)
		collection4 := &types.BadgeCollection{
			CollectionId: collectionId4,
			BalancesType: "Standard",
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
			},
			ValidBadgeIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			IbcUnwrapPaths: []*types.IBCUnwrapPath{
				{
					ChannelId:                      "channel-0",
					PortId:                         "transfer",
					SourceCollectionId:             sdkmath.NewUint(20),
					Denom:                          "test-denom-{id}",
					AllowOverrideWithAnyValidToken: true,
					Balances: []*types.Balance{
						{
							Amount: sdkmath.NewUint(100),
							BadgeIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
					DenomSuffixDetails: &types.DenomSuffixDetails{
						WithAddress: false,
					},
				},
			},
		}

		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection4)
		require.NoError(t, err)

		// Calculate expected IBC denom hash with override token ID
		overrideTokenId := "5"
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom-%s",
			collection4.IbcUnwrapPaths[0].SourceCollectionId.String(),
			overrideTokenId)
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Mint the IBC tokens to the creator first (simulating IBC transfer)
		creatorAddr, err := sdk.AccAddressFromBech32(creator)
		require.NoError(t, err)

		ibcCoin := sdk.NewCoin(ibcDenom, sdkmath.NewInt(1000000))
		err = suite.app.BankKeeper.MintCoins(ctx, "badges", sdk.NewCoins(ibcCoin))
		require.NoError(t, err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "badges", creatorAddr, sdk.NewCoins(ibcCoin))
		require.NoError(t, err)

		// Create message with override token ID
		msg := &types.MsgUnwrapIBCDenom{
			Creator:         creator,
			CollectionId:    collectionId4,
			OverrideTokenId: overrideTokenId,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.NoError(t, err)

		// Check that badges were minted
		userBalance, err := GetUserBalance(suite, wctx, collectionId4, creator)
		require.NoError(t, err)
		require.Len(t, userBalance.Balances, 1)
		require.Equal(t, sdkmath.NewUint(100001000), userBalance.Balances[0].Amount) // 1000 (default) + 100 * 1000000
	})

	// Test 5: Destination collection ID mismatch
	t.Run("Destination collection ID mismatch", func(t *testing.T) {
		// Create a new collection with destination gating
		collectionId5 := sdkmath.NewUint(5)
		collection5 := &types.BadgeCollection{
			CollectionId: collectionId5,
			BalancesType: "Standard",
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
			},
			ValidBadgeIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			IbcUnwrapPaths: []*types.IBCUnwrapPath{
				{
					ChannelId:                      "channel-0",
					PortId:                         "transfer",
					SourceCollectionId:             sdkmath.NewUint(20),
					Denom:                          "test-denom",
					AllowOverrideWithAnyValidToken: false,
					Balances: []*types.Balance{
						{
							Amount: sdkmath.NewUint(100),
							BadgeIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
					DenomSuffixDetails: &types.DenomSuffixDetails{
						WithAddress:             true,
						DestinationCollectionId: "999", // Different collection ID
						DestinationChainId:      ctx.ChainID(),
					},
				},
			},
		}

		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection5)
		require.NoError(t, err)

		// Calculate expected IBC denom hash with wrong destination collection ID
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom:%s:%s:%s",
			collection5.IbcUnwrapPaths[0].SourceCollectionId.String(),
			creator,
			"999", // Wrong collection ID
			ctx.ChainID())
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Create message
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId5,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message - should fail due to destination collection ID mismatch
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "destination collection ID mismatch")
	})

	// Test 6: Destination chain ID mismatch
	t.Run("Destination chain ID mismatch", func(t *testing.T) {
		// Create a new collection with destination gating
		collectionId6 := sdkmath.NewUint(6)
		collection6 := &types.BadgeCollection{
			CollectionId: collectionId6,
			BalancesType: "Standard",
			DefaultBalances: &types.UserBalanceStore{
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						BadgeIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
					},
				},
			},
			ValidBadgeIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			IbcUnwrapPaths: []*types.IBCUnwrapPath{
				{
					ChannelId:                      "channel-0",
					PortId:                         "transfer",
					SourceCollectionId:             sdkmath.NewUint(20),
					Denom:                          "test-denom",
					AllowOverrideWithAnyValidToken: false,
					Balances: []*types.Balance{
						{
							Amount: sdkmath.NewUint(100),
							BadgeIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
					DenomSuffixDetails: &types.DenomSuffixDetails{
						WithAddress:             true,
						DestinationCollectionId: collectionId6.String(),
						DestinationChainId:      "wrong-chain-id", // Wrong chain ID
					},
				},
			},
		}

		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection6)
		require.NoError(t, err)

		// Calculate expected IBC denom hash with wrong destination chain ID
		coreDenom := fmt.Sprintf("transfer/channel-0/badges:%s:test-denom:%s:%s:%s",
			collection6.IbcUnwrapPaths[0].SourceCollectionId.String(),
			creator,
			collectionId6.String(),
			"wrong-chain-id") // Wrong chain ID
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])
		ibcDenom := "ibc/" + expectedHash

		// Create message
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId6,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		// Execute the message - should fail due to destination chain ID mismatch
		_, err = suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "destination chain ID mismatch")
	})

	// Test 7: Invalid IBC denom format
	t.Run("Invalid IBC denom format", func(t *testing.T) {
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId,
			Amount: &sdk.Coin{
				Denom:  "invalid-denom", // Not ibc/ format
				Amount: sdkmath.NewInt(1000000),
			},
		}

		_, err := suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "denom must start with 'ibc/'")
	})

	// Test 8: No matching path found
	t.Run("No matching path found", func(t *testing.T) {
		// Create a random IBC denom that won't match any path
		randomHash := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
		ibcDenom := "ibc/" + randomHash

		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId,
			Amount: &sdk.Coin{
				Denom:  ibcDenom,
				Amount: sdkmath.NewInt(1000000),
			},
		}

		_, err := suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no matching IBC unwrap path found")
	})

	// Test 9: Collection not found
	t.Run("Collection not found", func(t *testing.T) {
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: sdkmath.NewUint(999), // Non-existent collection
			Amount: &sdk.Coin{
				Denom:  "ibc/abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				Amount: sdkmath.NewInt(1000000),
			},
		}

		_, err := suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "collection does not exist")
	})

	// Test 10: Zero amount
	t.Run("Zero amount", func(t *testing.T) {
		msg := &types.MsgUnwrapIBCDenom{
			Creator:      creator,
			CollectionId: collectionId,
			Amount: &sdk.Coin{
				Denom:  "ibc/abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				Amount: sdkmath.NewInt(0), // Zero amount
			},
		}

		_, err := suite.msgServer.UnwrapIBCDenom(wctx, msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "amount cannot be zero")
	})
}

// Helper function to simulate IBC denom hashing
func simulateIBCDenomHash(portId, channelId, sourceCollectionId, denom string, suffixes ...string) string {
	coreDenom := fmt.Sprintf("%s/%s/badges:%s:%s", portId, channelId, sourceCollectionId, denom)

	// Add suffixes
	for _, suffix := range suffixes {
		if suffix != "" {
			coreDenom += ":" + suffix
		}
	}

	hash := sha256.Sum256([]byte(coreDenom))
	return "ibc/" + hex.EncodeToString(hash[:])
}

func TestIBCDenomHashing(t *testing.T) {
	t.Run("Test IBC denom hashing simulation", func(t *testing.T) {
		// Test basic denom
		hash1 := simulateIBCDenomHash("transfer", "channel-0", "20", "test-denom")
		require.True(t, strings.HasPrefix(hash1, "ibc/"))
		require.Len(t, hash1, 68) // "ibc/" + 64 char hash

		// Test denom with address suffix
		hash2 := simulateIBCDenomHash("transfer", "channel-0", "20", "test-denom", "bb1abc123")
		require.True(t, strings.HasPrefix(hash2, "ibc/"))
		require.Len(t, hash2, 68)

		// Test denom with multiple suffixes
		hash3 := simulateIBCDenomHash("transfer", "channel-0", "20", "test-denom", "bb1abc123", "123", "bitbadges-1")
		require.True(t, strings.HasPrefix(hash3, "ibc/"))
		require.Len(t, hash3, 68)

		// Ensure different inputs produce different hashes
		require.NotEqual(t, hash1, hash2)
		require.NotEqual(t, hash2, hash3)
		require.NotEqual(t, hash1, hash3)
	})
}
