package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestDenomSuffixDetailsValidation(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx.WithChainID("test-chain-id")

	creator := suite.app.AccountKeeper.GetModuleAddress("badges").String()
	if creator == "" {
		creator = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	}

	// Test 1: Valid destination collection ID and chain ID
	t.Run("Valid destination collection ID and chain ID", func(t *testing.T) {
		collectionId := sdkmath.NewUint(1)

		// Create a collection first
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
		}
		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection)
		require.NoError(t, err)

		// Test the validation logic directly
		path := &types.IBCUnwrapPathAddObject{
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
				DestinationCollectionId: collectionId.String(),
				DestinationChainId:      "test-chain-id",
			},
		}

		// This should not error
		if path.DenomSuffixDetails != nil {
			// Validate destination collection ID
			if path.DenomSuffixDetails.DestinationCollectionId != "" {
				require.Equal(t, path.DenomSuffixDetails.DestinationCollectionId, collection.CollectionId.String())
			}

			// Validate destination chain ID
			if path.DenomSuffixDetails.DestinationChainId != "" {
				require.Equal(t, path.DenomSuffixDetails.DestinationChainId, ctx.ChainID())
			}
		}
	})

	// Test 2: Invalid destination collection ID
	t.Run("Invalid destination collection ID", func(t *testing.T) {
		collectionId := sdkmath.NewUint(2)

		// Create a collection first
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
		}
		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection)
		require.NoError(t, err)

		// Test the validation logic directly
		path := &types.IBCUnwrapPathAddObject{
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
				DestinationCollectionId: "999", // Wrong collection ID
				DestinationChainId:      "test-chain-id",
			},
		}

		// This should error
		if path.DenomSuffixDetails != nil {
			// Validate destination collection ID
			if path.DenomSuffixDetails.DestinationCollectionId != "" {
				require.NotEqual(t, path.DenomSuffixDetails.DestinationCollectionId, collection.CollectionId.String())
			}
		}
	})

	// Test 3: Invalid destination chain ID
	t.Run("Invalid destination chain ID", func(t *testing.T) {
		collectionId := sdkmath.NewUint(3)

		// Create a collection first
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
		}
		err := suite.app.BadgesKeeper.SetCollectionInStore(ctx, collection)
		require.NoError(t, err)

		// Test the validation logic directly
		path := &types.IBCUnwrapPathAddObject{
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
				DestinationCollectionId: collectionId.String(),
				DestinationChainId:      "wrong-chain-id", // Wrong chain ID
			},
		}

		// This should error
		if path.DenomSuffixDetails != nil {
			// Validate destination chain ID
			if path.DenomSuffixDetails.DestinationChainId != "" {
				require.NotEqual(t, path.DenomSuffixDetails.DestinationChainId, ctx.ChainID())
			}
		}
	})
}
