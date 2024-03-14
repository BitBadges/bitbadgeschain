package types_test

import (
	math "math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestActionPermissionUpdate(t *testing.T) {
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
	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateActionPermissionUpdate(sdk.Context{}, []*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(t, err)

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
	require.Error(t, err)
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
	require.Error(t, err)
}

func TestActionPermissionUpdateWithBadgeIds(t *testing.T) {
	oldActionPermission := &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	newActionPermission := &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: types.InvertUintRanges([]*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		BadgeIds: types.InvertUintRanges([]*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{

		BadgeIds: []*types.UintRange{
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
	newActionPermission2 := &types.BalancesActionPermission{}
	*newActionPermission2 = *newActionPermission
	newActionPermission2.BadgeIds = types.InvertUintRanges(newActionPermission2.BadgeIds, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))
	//Everything else from newActionPermission

	err = keeper.ValidateBalancesActionPermissionUpdate(sdk.Context{}, []*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestTimedUpdatePermission(t *testing.T) {
	oldActionPermission := &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	newActionPermission := &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		TimelineTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(80),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		TimelineTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{

		TimelineTimes: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		TimelineTimes: []*types.UintRange{
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
	newActionPermission2 := &types.TimedUpdatePermission{}
	*newActionPermission2 = *newActionPermission
	newActionPermission2.TimelineTimes = types.InvertUintRanges(newActionPermission2.TimelineTimes, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))

	err = keeper.ValidateTimedUpdatePermissionUpdate(sdk.Context{}, []*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestValidateTimedUpdatePermissionWithBadgeIds(t *testing.T) {
	oldActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

	newActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(122),
				},
			},
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			},
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{
			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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
			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	oldActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{{
		TimelineTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		BadgeIds: []*types.UintRange{
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
		TimelineTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		BadgeIds: []*types.UintRange{
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

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(60),
				},
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(61),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermission2 := []*types.TimedUpdateWithBadgeIdsPermission{
		{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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
	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	oldActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{{
		TimelineTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		BadgeIds: []*types.UintRange{
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

		TimelineTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		BadgeIds: []*types.UintRange{
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

	newActionPermissions = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermission2 = []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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
	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.Error(t, err)

}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate2(t *testing.T) {
	oldActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermission2 := []*types.TimedUpdateWithBadgeIdsPermission{{

		TimelineTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		},
		BadgeIds: []*types.UintRange{
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

		TimelineTimes: types.InvertUintRanges([]*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(100),
			},
		}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
		BadgeIds: []*types.UintRange{
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

	keeper, _ := keeper.BadgesKeeper(t)
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)
}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate3(t *testing.T) {
	oldActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(20),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(20),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermissions := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			},
			BadgeIds: []*types.UintRange{
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

			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(10),
					End:   sdkmath.NewUint(50),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	newActionPermission2 := []*types.TimedUpdateWithBadgeIdsPermission{
		{

			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			BadgeIds: []*types.UintRange{
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
			TimelineTimes: types.InvertUintRanges([]*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			}, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64)),
			BadgeIds: []*types.UintRange{
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

	keeper, _ := keeper.BadgesKeeper(t)
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(sdk.Context{}, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)
}

func TestValidateCollectionApprovalPermissionsUpdate(t *testing.T) {
	oldActionPermissions := []*types.CollectionApprovalPermission{{

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

	keeper, ctx := keeper.BadgesKeeper(t)
	newActionPermissions = append(newActionPermissions, newActionPermission2...)
	err := keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.NoError(t, err)

	newActionPermissions = []*types.CollectionApprovalPermission{
		{

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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
			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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
	require.Error(t, err)
}

func TestValidateCollectionApprovalPermissionsUpdate2(t *testing.T) {
	oldActionPermissions := []*types.CollectionApprovalPermission{{

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

			BadgeIds: types.InvertUintRanges([]*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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
		BadgeIds: types.InvertUintRanges([]*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.ValidateCollectionApprovalPermissionsUpdate(ctx, oldActionPermissions, newActionPermissions)
	require.Error(t, err)
}

func TestValidateCollectionApprovalPermissionsUpdate3(t *testing.T) {
	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "ABC",
		Addresses:        []string{bob, alice, charlie},
		Whitelist: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "Alice",
		Addresses:        []string{alice},
		Whitelist: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "BobCharlie",
		Addresses:        []string{bob, charlie},
		Whitelist: true,
	})
	require.NoError(t, err)

	oldActionPermissions := []*types.CollectionApprovalPermission{{

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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

			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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
	require.NoError(t, err)
}

func TestValidateCollectionApprovalPermissionsUpdate4Invalid(t *testing.T) {
	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "ABC",
		Addresses:        []string{bob, alice, charlie},
		Whitelist: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "Alice",
		Addresses:        []string{alice},
		Whitelist: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressList(ctx, &types.AddressList{
		ListId:        "BobCharlie",
		Addresses:        []string{bob, charlie},
		Whitelist: true,
	})
	require.NoError(t, err)

	oldActionPermissions := []*types.CollectionApprovalPermission{{

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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
			BadgeIds: []*types.UintRange{
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
			ApprovalId: "All",
			AmountTrackerId: "All",
			ChallengeTrackerId: "All",
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

		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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
		BadgeIds: []*types.UintRange{
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
		ApprovalId: "All",
		AmountTrackerId: "All",
		ChallengeTrackerId: "All",
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
	require.Error(t, err)
}
