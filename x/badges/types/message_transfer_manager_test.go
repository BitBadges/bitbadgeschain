package types_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgUpdateManager_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateManager
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateManager{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateManager{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				ManagerTimeline: []*types.ManagerTimeline{
					{
						Times: []*types.IdRange{
							{
								Start: sdkmath.NewUint(0),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
						Manager: sample.AccAddress(),
					},
				},
			},
		},
		{
			name: "invalid address 2",
			msg: types.MsgUpdateManager{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				ManagerTimeline: []*types.ManagerTimeline{
					{
						Times: []*types.IdRange{
							{
								Start: sdkmath.NewUint(0),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
						Manager: "invalid_address",
					},
				},
			},
			err: types.ErrInvalidAddress,
		},
		{
			name: "invalid times",
			msg: types.MsgUpdateManager{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				ManagerTimeline: []*types.ManagerTimeline{
					{
						Manager: "invalid_address",
					},
				},
			},
			err: types.ErrRangeDoesNotOverlap,
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
