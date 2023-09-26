package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestActionPermissionUpdate(t *testing.T) {
	oldActionPermission := &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			ForbiddenTimes: []*types.UintRange{},
		},
	}

	newActionPermission := &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{

			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{

			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(120),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(t, err)
	newActionPermission = &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{

			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(t, err)
}

func TestActionPermissionUpdateWithBadgeIds(t *testing.T) {
	oldActionPermission := &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			ForbiddenTimes: []*types.UintRange{},
		},
	}

	newActionPermission := &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.BalancesActionDefaultValues{

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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(120),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.BalancesActionPermission{
		Combinations: []*types.BalancesActionCombination{{}, {BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.BalancesActionDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateBalancesActionPermissionUpdate([]*types.BalancesActionPermission{oldActionPermission}, []*types.BalancesActionPermission{newActionPermission})
	require.NoError(t, err)
}

func TestTimedUpdatePermission(t *testing.T) {
	oldActionPermission := &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			ForbiddenTimes: []*types.UintRange{},
		},
	}

	newActionPermission := &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(122),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(120),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)
}

func TestValidateTimedUpdatePermissionWithBadgeIds(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
			},
			ForbiddenTimes: []*types.UintRange{},
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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

			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(120),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(80),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	oldActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)

	oldActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	println("START")

	err = keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.Error(t, err)

}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate2(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate3(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
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
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, _ := keeper.BadgesKeeper(t)
	err := keeper.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestValidateCollectionApprovedTransferPermissionsUpdate(t *testing.T) {
	oldActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, []*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission, newActionPermission2}, "")
	require.NoError(t, err)

	newActionPermission = &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "x",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission2 = &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, []*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission, newActionPermission2}, "")
	require.Error(t, err)
}

func TestValidateCollectionApprovedTransferPermissionsUpdate2(t *testing.T) {
	oldActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "Manager",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, []*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission}, "")
	require.Error(t, err)
}

func TestValidateCollectionApprovedTransferPermissionsUpdate3(t *testing.T) {
	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "ABC",
		Addresses:        []string{bob, alice, charlie},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "Alice",
		Addresses:        []string{alice},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "BobCharlie",
		Addresses:        []string{bob, charlie},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	oldActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "ABC",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "AllWithoutMint",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "AllWithoutMint",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "AllWithoutMint",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, []*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission}, "")
	require.NoError(t, err)
}

func TestValidateCollectionApprovedTransferPermissionsUpdate4Invalid(t *testing.T) {
	keeper, ctx := keeper.BadgesKeeper(t)
	err := keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "ABC",
		Addresses:        []string{bob, alice, charlie},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "Alice",
		Addresses:        []string{alice},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	err = keeper.CreateAddressMapping(ctx, &types.AddressMapping{
		MappingId:        "BobCharlie",
		Addresses:        []string{bob, charlie},
		IncludeAddresses: true,
	})
	require.NoError(t, err)

	oldActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "ABC",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "AllWithoutMint",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	newActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TransferTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
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
					End:  sdkmath.NewUint(50),
				},
			},
			ToMappingId:          "BobCharlie",
			FromMappingId:        "AllWithoutMint",
			InitiatedByMappingId: "AllWithoutMint",
			ApprovalTrackerId: "All",
			ChallengeTrackerId: "All",
			PermittedTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(50),
				},
				{
					Start: sdkmath.NewUint(200),
					End:   sdkmath.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(199),
				},
			},
		},
	}

	err = keeper.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, []*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission}, "")
	require.Error(t, err)
}
