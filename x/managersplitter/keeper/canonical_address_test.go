package keeper_test

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"
	splitterkeeper "github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperRejectsUppercaseSplitterIdentity(t *testing.T) {
	ms, ctx, _, _ := setupWithKeeper(t)
	admin := validAdmin(t)
	_, err := ms.CreateManagerSplitter(ctx, &types.MsgCreateManagerSplitter{Admin: strings.ToUpper(admin)})
	require.Error(t, err)
	delegate := sdk.AccAddress([]byte("splitter-delegate-aa")).String()
	_, err = ms.CreateManagerSplitter(ctx, &types.MsgCreateManagerSplitter{Admin: admin, Permissions: &types.ManagerSplitterPermissions{CanAddMoreAliasPaths: &types.PermissionCriteria{ApprovedAddresses: []string{strings.ToUpper(delegate)}}}})
	require.Error(t, err)
}

func TestV35CanonicalizesSplitterDelegates(t *testing.T) {
	_, _, k, ctx := setupWithKeeper(t)
	admin := validAdmin(t)
	address := types.DeriveManagerSplitterAddress(sdkmath.OneUint())
	require.NoError(t, k.SetManagerSplitterInStore(ctx, &types.ManagerSplitter{Address: address, Admin: strings.ToUpper(admin), Permissions: &types.ManagerSplitterPermissions{CanAddMoreAliasPaths: &types.PermissionCriteria{ApprovedAddresses: []string{strings.ToUpper(admin), admin}}}}))
	require.NoError(t, k.MigrateV35CanonicalAddresses(ctx))
	stored, found := k.GetManagerSplitterFromStore(ctx, address)
	require.True(t, found)
	require.Equal(t, admin, stored.Admin)
	require.Equal(t, []string{admin}, stored.Permissions.CanAddMoreAliasPaths.ApprovedAddresses)
	require.NoError(t, k.MigrateV35CanonicalAddresses(ctx))
}

func TestDeleteManagerSplitterRequiresManagerTransfer(t *testing.T) {
	app := bitbadgesapp.Setup(false)
	ctx := app.BaseApp.NewContext(false)
	ms := splitterkeeper.NewMsgServerImpl(app.ManagerSplitterKeeper)
	admin := validAdmin(t)
	created, err := ms.CreateManagerSplitter(ctx, &types.MsgCreateManagerSplitter{Admin: admin})
	require.NoError(t, err)
	collection := &tokenizationtypes.TokenCollection{CollectionId: sdkmath.OneUint(), Manager: created.Address}
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, collection, true))
	_, err = ms.DeleteManagerSplitter(ctx, &types.MsgDeleteManagerSplitter{Admin: admin, Address: created.Address})
	require.Error(t, err)
	_, found := app.ManagerSplitterKeeper.GetManagerSplitterFromStore(ctx, created.Address)
	require.True(t, found)
	collection.Manager = admin
	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, collection, true))
	_, err = ms.DeleteManagerSplitter(ctx, &types.MsgDeleteManagerSplitter{Admin: admin, Address: created.Address})
	require.NoError(t, err)
}
