package keeper

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgCreateManagerSplitter(t *testing.T) {
	// This is a basic test structure - expand with actual test setup
	// when you have the full test infrastructure
	msg := &types.MsgCreateManagerSplitter{
		Admin: "bb1test",
		Permissions: &types.ManagerSplitterPermissions{
			CanDeleteCollection: &types.PermissionCriteria{
				ApprovedAddresses: []string{"bb1approved1", "bb1approved2"},
			},
		},
	}

	require.NotNil(t, msg)
	require.Equal(t, "bb1test", msg.Admin)
	require.NotNil(t, msg.Permissions)
}

func TestDeriveManagerSplitterAddress(t *testing.T) {
	id := sdkmath.NewUint(1)
	address := types.DeriveManagerSplitterAddress(id)
	
	require.NotEmpty(t, address)
	// Address should be a valid bech32 address
	_, err := sdk.AccAddressFromBech32(address)
	require.NoError(t, err)
}

