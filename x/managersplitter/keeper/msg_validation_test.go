package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestMsgServerRejectsNilAndMalformedMessagesWithoutPanic(t *testing.T) {
	ms, goCtx, _, _ := setupWithKeeper(t)
	admin := validAdmin(t)

	resp, err := ms.CreateManagerSplitter(goCtx, &types.MsgCreateManagerSplitter{Admin: admin})
	require.NoError(t, err)

	// Nil message pointers return errors cleanly
	require.NotPanics(t, func() {
		_, err = ms.CreateManagerSplitter(goCtx, nil)
		require.Error(t, err)
		_, err = ms.UpdateManagerSplitter(goCtx, nil)
		require.Error(t, err)
		_, err = ms.DeleteManagerSplitter(goCtx, nil)
		require.Error(t, err)
		_, err = ms.ExecuteUniversalUpdateCollection(goCtx, nil)
		require.Error(t, err)
	})

	// Nil nested UniversalUpdateCollectionMsg returns error without panicking
	require.NotPanics(t, func() {
		_, err = ms.ExecuteUniversalUpdateCollection(goCtx, &types.MsgExecuteUniversalUpdateCollection{
			Executor:                     admin,
			ManagerSplitterAddress:       resp.Address,
			UniversalUpdateCollectionMsg: nil,
		})
		require.Error(t, err)
	})

	// Uninitialized CollectionId on nested message returns error without panicking
	require.NotPanics(t, func() {
		_, err = ms.ExecuteUniversalUpdateCollection(goCtx, &types.MsgExecuteUniversalUpdateCollection{
			Executor:               admin,
			ManagerSplitterAddress: resp.Address,
			UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
				Creator:                  admin,
				UpdateCollectionMetadata: true,
				CollectionMetadata:       &tokenizationtypes.CollectionMetadata{Uri: "https://example.com"},
			},
		})
		require.Error(t, err)
	})

	// Invalid approved address in optional permission criteria is rejected
	_, err = ms.UpdateManagerSplitter(goCtx, &types.MsgUpdateManagerSplitter{
		Admin:   admin,
		Address: resp.Address,
		Permissions: &types.ManagerSplitterPermissions{
			CanAddMoreAliasPaths: &types.PermissionCriteria{
				ApprovedAddresses: []string{"not-bech32"},
			},
		},
	})
	require.Error(t, err)

	// Equal admin and splitter address rejected on update and delete
	_, err = ms.UpdateManagerSplitter(goCtx, &types.MsgUpdateManagerSplitter{
		Admin:   admin,
		Address: admin,
	})
	require.Error(t, err)
	_, err = ms.DeleteManagerSplitter(goCtx, &types.MsgDeleteManagerSplitter{
		Admin:   admin,
		Address: admin,
	})
	require.Error(t, err)

	// Sanity check valid delete still succeeds
	_, err = ms.DeleteManagerSplitter(goCtx, &types.MsgDeleteManagerSplitter{
		Admin:   admin,
		Address: resp.Address,
	})
	require.NoError(t, err)
	_ = sdk.AccAddress{}
}
