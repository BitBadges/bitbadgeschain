package keeper_test

import (
	"context"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
)

// Like setupMsgServer, but also hands back the keeper and sdk.Context so a test
// can read the store back and assert the message actually persisted something.
func setupWithKeeper(t *testing.T) (types.MsgServer, context.Context, keeper.Keeper, sdk.Context) {
	t.Helper()
	app := bitbadgesapp.Setup(false)
	ctx := app.BaseApp.NewContext(false)
	return keeper.NewMsgServerImpl(app.ManagerSplitterKeeper), sdk.WrapSDKContext(ctx),
		app.ManagerSplitterKeeper, ctx
}

func validAdmin(t *testing.T) string {
	t.Helper()
	return sdk.AccAddress([]byte("managersplitter-admin")).String()
}

func TestCreateManagerSplitterStoresAndReturnsDerivedAddress(t *testing.T) {
	ms, goCtx, k, ctx := setupWithKeeper(t)
	admin := validAdmin(t)

	perms := &types.ManagerSplitterPermissions{
		CanDeleteCollection: &types.PermissionCriteria{
			ApprovedAddresses: []string{"bb1approved1", "bb1approved2"},
		},
	}

	resp, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{
		Admin:       admin,
		Permissions: perms,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// The address is derived from the id counter, not supplied by the caller.
	require.Equal(t, types.DeriveManagerSplitterAddress(sdkmath.NewUint(1)), resp.Address)

	stored, found := k.GetManagerSplitterFromStore(ctx, resp.Address)
	require.True(t, found, "CreateManagerSplitter must persist the splitter")
	require.Equal(t, admin, stored.Admin)
	require.Equal(t, perms.CanDeleteCollection.ApprovedAddresses,
		stored.Permissions.CanDeleteCollection.ApprovedAddresses)
}

func TestCreateManagerSplitterRejectsInvalidAdmin(t *testing.T) {
	ms, goCtx, _, _ := setupWithKeeper(t)

	// "bb1test" is not valid bech32. The placeholder test this replaced used it
	// as a fixture and asserted the field it had just set, so it never reached
	// the code that rejects it.
	_, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{Admin: "bb1test"})
	require.ErrorIs(t, err, types.ErrInvalidAdmin)
}

func TestCreateManagerSplitterIncrementsIdAndDerivesDistinctAddresses(t *testing.T) {
	ms, goCtx, _, _ := setupWithKeeper(t)
	admin := validAdmin(t)

	first, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{Admin: admin})
	require.NoError(t, err)
	second, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{Admin: admin})
	require.NoError(t, err)

	require.NotEqual(t, first.Address, second.Address,
		"a second splitter must not collide with the first")
	require.Equal(t, types.DeriveManagerSplitterAddress(sdkmath.NewUint(2)), second.Address)
}

func TestCreateManagerSplitterNormalisesNilPermissions(t *testing.T) {
	ms, goCtx, k, ctx := setupWithKeeper(t)

	resp, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{
		Admin:       validAdmin(t),
		Permissions: nil,
	})
	require.NoError(t, err)

	stored, found := k.GetManagerSplitterFromStore(ctx, resp.Address)
	require.True(t, found)
	require.NotNil(t, stored.Permissions,
		"nil permissions must be stored as an empty set, not left nil")
}
