package types_test

import (
	math "math"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestActionPermissionUpdate() {
	oldActionPermission := &types.ActionPermission{

		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	newActionPermission := &types.ActionPermission{

		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}
	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	err := keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.ActionPermission{

		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(120),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(suite.T(), err)
	newActionPermission = &types.ActionPermission{

		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(suite.T(), err)
}

func GetFullUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func (suite *TestSuite) TestActionPermissionUpdateWithTokenIds() {
	oldActionPermission := &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	newActionPermission := &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	err := keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(122),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(122),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),
		TokenIds: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(120),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.CollectionApprovalPermission{
		ApprovalId:        "All",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	//copy newActionPermission to newActionPermission2
	newActionPermission2 := &types.CollectionApprovalPermission{}
	*newActionPermission2 = *newActionPermission
	newActionPermission2.TokenIds = types.InvertUintRanges(newActionPermission2.TokenIds, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))
	//Everything else from newActionPermission

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(sdk.Context{}, []*types.CollectionApprovalPermission{oldActionPermission}, []*types.CollectionApprovalPermission{newActionPermission, newActionPermission2})
	require.NoError(suite.T(), err)
}

func (suite *TestSuite) TestTimedUpdatePermission() {
	oldActionPermission := &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	newActionPermission := &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	err := keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(120),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(suite.T(), err)

	newActionPermission = &types.ActionPermission{
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}

}

// COMMENTED OUT: This test depends on timeline time logic that was removed
/*
func (suite *TestSuite) TestValidateTimedUpdatePermissionWithTokenIds() {
	oldActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{},
		},
	}

	newActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	err := keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(122),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			},

			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(120),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	oldActionPermissions = []*types.TokenIdsActionPermission{{
		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {
		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}}

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(62),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(62),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}, {

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := []*types.TokenIdsActionPermission{
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	oldActionPermissions = []*types.TokenIdsActionPermission{{
		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}}

	newActionPermissions = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 = []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	// Note: Validation doesn't check for TokenIds shrinking, only PermanentlyPermittedTimes
	// Since PermanentlyPermittedTimes were expanded, this is valid
	err = keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

}
*/
// END COMMENTED OUT TEST

func (suite *TestSuite) TestValidateTokenIdsActionPermissionUpdate2() {
	oldActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}, {

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}}

	newActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			}},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			}},
	}

	newActionPermission2 := []*types.TokenIdsActionPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(9),
			},
			{
				Start: sdkmath.NewUint(50),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(9),
			},
			{
				Start: sdkmath.NewUint(50),
				End:   sdkmath.NewUint(100),
			},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
	}

	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)
}

func (suite *TestSuite) TestValidateTokenIdsActionPermissionUpdate3() {
	oldActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}, {

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermissions := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}, {

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := []*types.TokenIdsActionPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(9),
				},
				{
					Start: sdkmath.NewUint(50),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		}, {
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(9),
				},
				{
					Start: sdkmath.NewUint(50),
					End:   sdkmath.NewUint(100),
				},
			},
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keepertest.TokenizationKeeper(suite.T())
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateTokenIdsActionPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)
}

func (suite *TestSuite) TestValidateCollectionApprovalPermissionsUpdate() {

	oldActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
	}

	newActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(10),
				End:   sdkmath.NewUint(50),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		}},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			}},
	}

	newActionPermission2 := []*types.CollectionApprovalPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(9),
				},
				{
					Start: sdkmath.NewUint(50),
					End:   sdkmath.NewUint(100),
				},
			},
			TransferTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(9),
				},
				{
					Start: sdkmath.NewUint(50),
					End:   sdkmath.NewUint(100),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)

	newActionPermissions = []*types.CollectionApprovalPermission{
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			TransferTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "x",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "x",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 = []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(9),
			},
			{
				Start: sdkmath.NewUint(50),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(9),
				},
				{
					Start: sdkmath.NewUint(50),
					End:   sdkmath.NewUint(100),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err = keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.Error(suite.T(), err)
}

func (suite *TestSuite) TestValidateCollectionApprovalPermissionsUpdate2() {
	oldActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
		{

			TokenIds: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			TransferTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: alice,
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {
		TokenIds: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: alice,
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
	}

	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	err := keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.Error(suite.T(), err)
}

func (suite *TestSuite) TestValidateCollectionApprovalPermissionsUpdate3() {
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	err := keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "ABC",
		Addresses: []string{bob, alice, charlie},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "Alice",
		Addresses: []string{alice},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "BobCharlie",
		Addresses: []string{bob, charlie},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	oldActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "ABC",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "ABC",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
	}

	newActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "AllWithoutMint",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
		{

			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.NoError(suite.T(), err)
}

func (suite *TestSuite) TestValidateCollectionApprovalPermissionsUpdate4Invalid() {
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	err := keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "ABC",
		Addresses: []string{bob, alice, charlie},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "Alice",
		Addresses: []string{alice},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:    "BobCharlie",
		Addresses: []string{bob, charlie},
		Whitelist: true,
	})
	require.NoError(suite.T(), err)

	oldActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "ABC",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
		{
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			TransferTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			OwnershipTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
			},
			ToListId:          "ABC",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "All",
			PermanentlyPermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			PermanentlyForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermissions := []*types.CollectionApprovalPermission{{

		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "BobCharlie",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	}, {
		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		TransferTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
		},
		ToListId:          "BobCharlie",
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalId:        "All",
		PermanentlyPermittedTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(50),
			},
			{
				Start: sdkmath.NewUint(200),
				End:   sdkmath.NewUint(300),
			},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(101),
				End:   sdkmath.NewUint(199),
			},
		},
	},
	}

	err = keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.Error(suite.T(), err)
}
