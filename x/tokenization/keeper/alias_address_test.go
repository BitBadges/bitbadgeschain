package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func makeTestCollection() *types.TokenCollection {
	return &types.TokenCollection{
		CollectionId:      sdkmath.NewUint(1),
		MintEscrowAddress: "bb1mintescrow",
		CosmosCoinWrapperPaths: []*types.CosmosCoinWrapperPath{
			{Address: "bb1wrapper0", Denom: "uatom"},
			{Address: "bb1wrapper1", Denom: "uosmo"},
		},
		Invariants: &types.CollectionInvariants{
			CosmosCoinBackedPath: &types.CosmosCoinBackedPath{
				Address: "bb1ibcbacking",
			},
		},
	}
}

func TestResolveAddressAlias(t *testing.T) {
	collection := makeTestCollection()

	tests := []struct {
		name    string
		alias   string
		want    string
		wantErr bool
	}{
		{"MintEscrow", "MintEscrow", "bb1mintescrow", false},
		{"CosmosWrapper/0", "CosmosWrapper/0", "bb1wrapper0", false},
		{"CosmosWrapper/1", "CosmosWrapper/1", "bb1wrapper1", false},
		{"IBCBacking", "IBCBacking", "bb1ibcbacking", false},
		{"CosmosWrapper out of range", "CosmosWrapper/5", "", true},
		{"unknown alias", "UnknownAlias", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keeper.ResolveAddressAlias(collection, tt.alias)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestResolveAddressIfAlias(t *testing.T) {
	collection := makeTestCollection()

	// Alias should resolve
	got, err := keeper.ResolveAddressIfAlias(collection, "MintEscrow")
	require.NoError(t, err)
	require.Equal(t, "bb1mintescrow", got)

	// Regular address should pass through
	got, err = keeper.ResolveAddressIfAlias(collection, "bb1someregularaddr")
	require.NoError(t, err)
	require.Equal(t, "bb1someregularaddr", got)

	// Empty should pass through
	got, err = keeper.ResolveAddressIfAlias(collection, "")
	require.NoError(t, err)
	require.Equal(t, "", got)
}

func TestResolveListIdAliases(t *testing.T) {
	collection := makeTestCollection()

	tests := []struct {
		name    string
		listId  string
		want    string
		wantErr bool
	}{
		// Reserved keywords pass through
		{"All", "All", "All", false},
		{"None", "None", "None", false},
		{"Mint", "Mint", "Mint", false},
		{"AllWithMint", "AllWithMint", "AllWithMint", false},

		// Single alias
		{"single alias", "MintEscrow", "bb1mintescrow", false},
		{"single IBCBacking", "IBCBacking", "bb1ibcbacking", false},

		// Colon-separated with aliases
		{"alias:addr", "MintEscrow:bb1other", "bb1mintescrow:bb1other", false},
		{"alias:alias", "MintEscrow:CosmosWrapper/0", "bb1mintescrow:bb1wrapper0", false},

		// AllWithout
		{"AllWithoutAlias", "AllWithoutMintEscrow", "AllWithoutbb1mintescrow", false},
		{"AllWithoutMultiple", "AllWithoutMintEscrow:CosmosWrapper/0", "AllWithoutbb1mintescrow:bb1wrapper0", false},

		// Inversion
		{"!alias", "!MintEscrow", "!bb1mintescrow", false},
		{"!(alias:addr)", "!(MintEscrow:bb1other)", "!(bb1mintescrow:bb1other)", false},

		// Regular addresses pass through
		{"regular addr", "bb1someaddr", "bb1someaddr", false},

		// Empty
		{"empty", "", "", false},

		// Error case
		{"bad alias", "CosmosWrapper/99", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keeper.ResolveListIdAliases(collection, tt.listId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestResolveAliasNoMintEscrow(t *testing.T) {
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
	}
	_, err := keeper.ResolveAddressAlias(collection, "MintEscrow")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no mint escrow address")
}

func TestResolveAliasNoIBCBacking(t *testing.T) {
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
	}
	_, err := keeper.ResolveAddressAlias(collection, "IBCBacking")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no IBC backing path")
}

func TestIsAddressAlias(t *testing.T) {
	require.True(t, types.IsAddressAlias("MintEscrow"))
	require.True(t, types.IsAddressAlias("IBCBacking"))
	require.True(t, types.IsAddressAlias("CosmosWrapper/0"))
	require.True(t, types.IsAddressAlias("CosmosWrapper/123"))

	require.False(t, types.IsAddressAlias("Mint"))
	require.False(t, types.IsAddressAlias("bb1someaddr"))
	require.False(t, types.IsAddressAlias(""))
	require.False(t, types.IsAddressAlias("CosmosWrapper/"))
	require.False(t, types.IsAddressAlias("CosmosWrapper/abc"))
	require.False(t, types.IsAddressAlias("CosmosWrapper"))
}

func TestIsReservedAliasListId(t *testing.T) {
	require.True(t, types.IsReservedAliasListId("MintEscrow"))
	require.True(t, types.IsReservedAliasListId("IBCBacking"))
	require.True(t, types.IsReservedAliasListId("CosmosWrapper/0"))
	require.True(t, types.IsReservedAliasListId("CosmosWrapper/anything"))

	require.False(t, types.IsReservedAliasListId("Mint"))
	require.False(t, types.IsReservedAliasListId("my-custom-list"))
	require.False(t, types.IsReservedAliasListId(""))
}
