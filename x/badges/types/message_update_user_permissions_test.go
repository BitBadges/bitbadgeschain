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
		msg types.MsgUpdateUserPermissions
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateUserPermissions{
				Creator: "invalid_address",
				CollectionId: sdkmath.NewUint(1),
				Permissions: &types.UserPermissions{},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUserPermissions{
				Creator: sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				Permissions: &types.UserPermissions{},
			},
		},
		{
			name: "no permissions",
			msg: types.MsgUpdateUserPermissions{
				Creator: sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
			},
			err: types.ErrPermissionsIsNil,
		},
		{
			name: "overlap times",
			msg: types.MsgUpdateUserPermissions{
				Creator: sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				Permissions: &types.UserPermissions{
					CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransferPermission{
						{
							DefaultValues: &types.UserApprovedOutgoingTransferDefaultValues{
								PermittedTimes: []*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2) } },
								ForbiddenTimes: []*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2) } },
							},
							Combinations: []*types.UserApprovedOutgoingTransferCombination{{}},
						},
					},
				},
			},
			err: types.ErrPermissionsIsNil,
		},
		{
			name: "valid",
			msg: types.MsgUpdateUserPermissions{
				Creator: sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				Permissions: GetValidUserPermissions(),
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
