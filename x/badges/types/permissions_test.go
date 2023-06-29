package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestActionPermissionUpdate(t *testing.T) {
	oldActionPermission := &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{
			PermittedTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			ForbiddenTimes: []*types.IdRange{
				
			},
		},
	}

	newActionPermission := &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{
			
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{
			
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(120),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(t, err)
	newActionPermission = &types.ActionPermission{
		Combinations: []*types.ActionCombination{{}},
		DefaultValues: &types.ActionDefaultValues{
			
			PermittedTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(101),
					End:   sdk.NewUint(199),
				},
			},
		
		},
	}

	err = types.ValidateActionPermissionUpdate([]*types.ActionPermission{oldActionPermission}, []*types.ActionPermission{newActionPermission})
	require.Error(t, err)
}

func TestActionPermissionUpdateWithBadgeIds(t *testing.T) {
	oldActionPermission := &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					
				},
			
		},
	}

	newActionPermission := &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)


	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(122),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
			PermittedTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(101),
					End:   sdk.NewUint(199),
				},
			},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			
				BadgeIds: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
				},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(120),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(80),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.ActionWithBadgeIdsPermission{
		Combinations: []*types.ActionWithBadgeIdsCombination{{}, {BadgeIdsOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateActionWithBadgeIdsPermissionUpdate([]*types.ActionWithBadgeIdsPermission{oldActionPermission}, []*types.ActionWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)
}


func TestTimedUpdatePermission(t *testing.T) {
	oldActionPermission := &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					
				},
			
		},
	}

	newActionPermission := &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)


	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(122),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	//TODO:
	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
			PermittedTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
			ForbiddenTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(101),
					End:   sdk.NewUint(199),
				},
			},
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(120),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(80),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdatePermission{
		Combinations: []*types.TimedUpdateCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
			
		},
	}

	err = types.ValidateTimedUpdatePermissionUpdate([]*types.TimedUpdatePermission{oldActionPermission}, []*types.TimedUpdatePermission{newActionPermission})
	require.NoError(t, err)
}

func TestValidateTimedUpdatePermissionWithBadgeIds(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					
				},
			
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)


	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(122),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(122),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)
	

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(80),
				},
			},
			
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
				
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				{
					Start: sdk.NewUint(200),
					End:   sdk.NewUint(300),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(120),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(80),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	oldActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(50),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			},
		
	}


	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(60),
				},
				{
					Start: sdk.NewUint(61),
					End:   sdk.NewUint(100),
				},
				
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(60),
				},
				{
					Start:  sdk.NewUint(61),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			},
		
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.NoError(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(60),
				},
				{
					Start: sdk.NewUint(61),
					End:   sdk.NewUint(100),
				},
				
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(60),
				},
				{
					Start:  sdk.NewUint(62),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
			
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission})
	require.Error(t, err)

	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(61),
					End:   sdk.NewUint(100),
				},
				
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(61),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
				
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
		
			
		},
	}

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
	
	oldActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(50),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
		},
	}


	newActionPermission = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(10),
					End:   sdk.NewUint(50),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
		
			
		},
	}

	newActionPermission2 = &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	println("START")

	err = types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.Error(t, err)

}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate2(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(50),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(10),
					End:   sdk.NewUint(50),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(9),
				},
				{
					Start:  sdk.NewUint(50),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestValidateTimedUpdateWithBadgeIdsPermissionUpdate3(t *testing.T) {
	oldActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(20),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(50),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(10),
					End:   sdk.NewUint(50),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission2 := &types.TimedUpdateWithBadgeIdsPermission{
		Combinations: []*types.TimedUpdateWithBadgeIdsCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.TimedUpdateWithBadgeIdsDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(9),
				},
				{
					Start:  sdk.NewUint(50),
					End:    sdk.NewUint(100),
				},
			},
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	err := types.ValidateTimedUpdateWithBadgeIdsPermissionUpdate([]*types.TimedUpdateWithBadgeIdsPermission{oldActionPermission}, []*types.TimedUpdateWithBadgeIdsPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)
}

func TestValidateCollectionApprovedTransferPermissionsUpdate(t *testing.T) {
	oldActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(20),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			TransferTimes: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			ToMappingId: "test",
			FromMappingId: "dfsaf",
			InitiatedByMappingId: "fdsjhksad",
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(50),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(10),
					End:   sdk.NewUint(50),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
			TransferTimes: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			ToMappingId: "test",
			FromMappingId: "dfsaf",
			InitiatedByMappingId: "fdsjhksad",
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission2 := &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(9),
				},
				{
					Start:  sdk.NewUint(50),
					End:    sdk.NewUint(100),
				},
			},
			TransferTimes: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			ToMappingId: "test",
			FromMappingId: "dfsaf",
			InitiatedByMappingId: "fdsjhksad",
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	
	err := types.ValidateCollectionApprovedTransferPermissionsUpdate([]*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission, newActionPermission2})
	require.NoError(t, err)



	newActionPermission = &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(10),
					End:   sdk.NewUint(50),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(10),
					End:    sdk.NewUint(50),
				},
			},
			TransferTimes: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			ToMappingId: "x",
			FromMappingId: "dfsaf",
			InitiatedByMappingId: "fdsjhksad",
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	newActionPermission2 = &types.CollectionApprovedTransferPermission{
		Combinations: []*types.CollectionApprovedTransferCombination{{}, {TimelineTimesOptions: &types.ValueOptions{InvertDefault: true}}},
		DefaultValues: &types.CollectionApprovedTransferDefaultValues{
			TimelineTimes: []*types.IdRange{
				{
					Start: sdk.NewUint(0),
					End:   sdk.NewUint(100),
				},
			},
			BadgeIds: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(9),
				},
				{
					Start:  sdk.NewUint(50),
					End:    sdk.NewUint(100),
				},
			},
			TransferTimes: []*types.IdRange{
				{
					Start:  sdk.NewUint(0),
					End:    sdk.NewUint(100),
				},
			},
			ToMappingId: "test",
			FromMappingId: "dfsaf",
			InitiatedByMappingId: "fdsjhksad",
				PermittedTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(0),
						End:   sdk.NewUint(100),
					},
					{
						Start: sdk.NewUint(200),
						End:   sdk.NewUint(300),
					},
				},
				ForbiddenTimes: []*types.IdRange{
					{
						Start: sdk.NewUint(101),
						End:   sdk.NewUint(199),
					},
				},
			
		},
	}

	
	err = types.ValidateCollectionApprovedTransferPermissionsUpdate([]*types.CollectionApprovedTransferPermission{oldActionPermission}, []*types.CollectionApprovedTransferPermission{newActionPermission, newActionPermission2})
	require.Error(t, err)
}