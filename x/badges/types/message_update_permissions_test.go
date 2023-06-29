package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgUpdateCollectionPermissions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCollectionPermissions
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateCollectionPermissions{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
				Permissions:  &types.CollectionPermissions{},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateCollectionPermissions{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				Permissions:  &types.CollectionPermissions{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
