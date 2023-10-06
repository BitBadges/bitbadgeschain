package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateUserPermissions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateUserApprovals
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateUserApprovals{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
				UserPermissions:  &types.UserPermissions{},
				UpdateUserPermissions: true,
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUserApprovals{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				UserPermissions:  &types.UserPermissions{},
				UpdateUserPermissions: true,
			},
		},
		// {
		// 	name: "no permissions",
		// 	msg: types.MsgUpdateUserApprovals{
		// 		Creator:      sample.AccAddress(),
		// 		CollectionId: sdkmath.NewUint(1),
		// 		UpdateUserPermissions: true,
		// 	},
		// 	err: types.ErrPermissionsIsNil,
		// },
		{
			name: "overlap times",
			msg: types.MsgUpdateUserApprovals{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				UpdateUserPermissions: true,
				UserPermissions: &types.UserPermissions{
					CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{
						{
							PermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
							ForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
						},
					},
				},
			},
			err: types.ErrPermissionsIsNil,
		},
		{
			name: "valid",
			msg: types.MsgUpdateUserApprovals{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				UserPermissions:  GetValidUserPermissions(),
				UpdateUserPermissions: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
