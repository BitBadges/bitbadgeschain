package keeper_test

import (
	"bytes"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/module"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestCollectionCreatedBlockAndTimestampImmutableAndGenesisRoundTrip(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()

	createdTime := time.Date(2026, 6, 15, 12, 30, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockHeight(4200).WithBlockTime(createdTime)

	creator := sdk.AccAddress(bytes.Repeat([]byte{0x11}, 20)).String()
	createMsg := &types.MsgCreateCollection{
		Creator: creator,
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/col",
			CustomData: "v1",
		},
		Manager:               creator,
		CollectionPermissions: &types.CollectionPermissions{},
	}

	res, err := suite.msgServer.CreateCollection(suite.ctx, createMsg)
	require.NoError(t, err)

	col, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, res.CollectionId)
	require.True(t, found)
	require.Equal(t, uint64(4200), col.CreatedBlock)
	require.Equal(t, uint64(createdTime.UnixMilli()), col.CreatedTimestamp)

	// Advance block height & time and update the collection; createdBlock and createdTimestamp must remain immutable.
	laterTime := createdTime.Add(24 * time.Hour)
	suite.ctx = suite.ctx.WithBlockHeight(9999).WithBlockTime(laterTime)

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:          creator,
		CollectionId:     res.CollectionId,
		UpdateCustomData: true,
		CustomData:       "v2",
	}
	_, err = suite.msgServer.UniversalUpdateCollection(suite.ctx, updateMsg)
	require.NoError(t, err)

	updated, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, res.CollectionId)
	require.True(t, found)
	require.Equal(t, "v2", updated.CustomData)
	require.Equal(t, uint64(4200), updated.CreatedBlock)
	require.Equal(t, uint64(createdTime.UnixMilli()), updated.CreatedTimestamp)

	// Export and re-import genesis: createdBlock and createdTimestamp must survive intact.
	genState := tokenization.ExportGenesis(suite.ctx, *suite.app.TokenizationKeeper)
	var exportedCol *types.TokenCollection
	for _, c := range genState.Collections {
		if c.CollectionId.Equal(res.CollectionId) {
			exportedCol = c
			break
		}
	}
	require.NotNil(t, exportedCol)
	require.Equal(t, uint64(4200), exportedCol.CreatedBlock)
	require.Equal(t, uint64(createdTime.UnixMilli()), exportedCol.CreatedTimestamp)

	suite2 := new(TestSuite)
	suite2.SetT(t)
	suite2.SetupTest()
	tokenization.InitGenesis(suite2.ctx, *suite2.app.TokenizationKeeper, *genState)
	reimported, found := suite2.app.TokenizationKeeper.GetCollectionFromStore(suite2.ctx, res.CollectionId)
	require.True(t, found)
	require.Equal(t, uint64(4200), reimported.CreatedBlock)
	require.Equal(t, uint64(createdTime.UnixMilli()), reimported.CreatedTimestamp)
}

func TestQueryResolveDenom(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()

	creator := sdk.AccAddress(bytes.Repeat([]byte{0x11}, 20)).String()
	fullTimes := GetFullUintRanges()
	createMsg := &types.MsgCreateCollection{
		Creator: creator,
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		},
		Manager:               creator,
		CollectionPermissions: &types.CollectionPermissions{},
		CosmosCoinWrapperPathsToAdd: []*types.CosmosCoinWrapperPathAddObject{
			{
				Denom:  "mycoin",
				Symbol: "MYC",
				Conversion: &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{Amount: sdkmath.NewUint(1)},
					SideB: []*types.Balance{
						{
							Amount:         sdkmath.NewUint(1),
							TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
							OwnershipTimes: fullTimes,
						},
					},
				},
			},
			{
				Denom:                          "nft-{id}",
				Symbol:                         "NFT",
				AllowOverrideWithAnyValidToken: true,
				Conversion: &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{Amount: sdkmath.NewUint(1)},
					SideB: []*types.Balance{
						{
							Amount:         sdkmath.NewUint(1),
							TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
							OwnershipTimes: fullTimes,
						},
					},
				},
			},
		},
		AliasPathsToAdd: []*types.AliasPathAddObject{
			{
				Denom:  "myalias",
				Symbol: "MYA",
				Conversion: &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{Amount: sdkmath.NewUint(1)},
					SideB: []*types.Balance{
						{
							Amount:         sdkmath.NewUint(1),
							TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
							OwnershipTimes: fullTimes,
						},
					},
				},
			},
		},
	}

	res, err := suite.msgServer.CreateCollection(suite.ctx, createMsg)
	require.NoError(t, err)
	cid := res.CollectionId

	k := suite.app.TokenizationKeeper

	// 1. Exact wrapper match
	wrapResp, err := k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{
		Denom: "badges:" + cid.String() + ":mycoin",
	})
	require.NoError(t, err)
	require.Equal(t, "wrapper", wrapResp.PathType)
	require.True(t, wrapResp.CollectionId.Equal(cid))
	require.Equal(t, "mycoin", wrapResp.BaseDenom)
	require.Empty(t, wrapResp.MatchedTokenId)
	require.NotNil(t, wrapResp.WrapperPath)
	require.Equal(t, "MYC", wrapResp.WrapperPath.Symbol)
	require.Nil(t, wrapResp.AliasPath)

	// 2. Placeholder {id} wrapper match (valid tokenId)
	placeholderResp, err := k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{
		Denom: "badges:" + cid.String() + ":nft-42",
	})
	require.NoError(t, err)
	require.Equal(t, "wrapper", placeholderResp.PathType)
	require.True(t, placeholderResp.CollectionId.Equal(cid))
	require.Equal(t, "nft-42", placeholderResp.BaseDenom)
	require.Equal(t, "42", placeholderResp.MatchedTokenId)
	require.NotNil(t, placeholderResp.WrapperPath)
	require.Equal(t, "nft-{id}", placeholderResp.WrapperPath.Denom)

	// 3. Placeholder {id} out of ValidTokenIds range (51 > 50) or non-canonical ("042") -> error
	_, err = k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{
		Denom: "badges:" + cid.String() + ":nft-51",
	})
	require.Error(t, err)
	_, err = k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{
		Denom: "badges:" + cid.String() + ":nft-042",
	})
	require.Error(t, err)

	// 4. Alias match
	aliasResp, err := k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{
		Denom: "badgeslp:" + cid.String() + ":myalias",
	})
	require.NoError(t, err)
	require.Equal(t, "alias", aliasResp.PathType)
	require.True(t, aliasResp.CollectionId.Equal(cid))
	require.Equal(t, "myalias", aliasResp.BaseDenom)
	require.NotNil(t, aliasResp.AliasPath)
	require.Equal(t, "MYA", aliasResp.AliasPath.Symbol)
	require.Nil(t, aliasResp.WrapperPath)

	// 5. Invalid / unknown denoms
	for _, bad := range []string{"", "ubadge", "badges:9999:mycoin", "badges:" + cid.String() + ":unknown", "badgeslp:" + cid.String() + ":unknown"} {
		_, err = k.ResolveDenom(suite.ctx, &types.QueryResolveDenomRequest{Denom: bad})
		require.Error(t, err, "expected error for denom %q", bad)
	}
}
