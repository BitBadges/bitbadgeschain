package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"bitbadgeschain/x/badges/testutil/sample"
	"bitbadgeschain/x/badges/types"
)

func TestMsgUniversalUpdateCollectionPermissions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUniversalUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     "invalid_address",
				CollectionId:                sdkmath.NewUint(1),
				UpdateCollectionPermissions: true,
				CollectionPermissions: &types.CollectionPermissions{
					CanDeleteCollection: []*types.ActionPermission{
						{
							PermanentlyPermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
							PermanentlyForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
						},
					},
				},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:               sample.AccAddress(),
				CollectionId:          sdkmath.NewUint(1),
				CollectionPermissions: &types.CollectionPermissions{},
			},
		},
		{
			name: "invalid permissions",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     sample.AccAddress(),
				CollectionId:                sdkmath.NewUint(1),
				UpdateCollectionPermissions: true,
				CollectionPermissions: &types.CollectionPermissions{
					CanDeleteCollection: []*types.ActionPermission{
						{
							PermanentlyPermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
							PermanentlyForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
						},
					},
				},
			},
			err: types.ErrRangesOverlap,
		},
		{
			name: "valid permissions",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     sample.AccAddress(),
				CollectionId:                sdkmath.NewUint(1),
				UpdateCollectionPermissions: true,
				CollectionPermissions: &types.CollectionPermissions{
					CanDeleteCollection: []*types.ActionPermission{
						{
							PermanentlyPermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
							PermanentlyForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(22)}},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
