package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateManagerApprovedTransfers_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateManagerApprovedTransfers
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateManagerApprovedTransfers{
				Creator: "invalid_address",
				ManagerApprovedTransfers: []*types.TransferMapping{
					{
						From: &types.Addresses{
							AccountNums: []*types.IdRange{
								{
									Start: 0,
									End: 0,
								},
							},
							ManagerOptions: types.ManagerOptions_Neutral,
						
					},
						To: &types.Addresses{
							AccountNums: []*types.IdRange{
								{
									Start: 0,
									End: 0,
								},
							},
							ManagerOptions: types.ManagerOptions_Neutral,
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateManagerApprovedTransfers{
				Creator: sample.AccAddress(),
				ManagerApprovedTransfers: []*types.TransferMapping{
					{
						From: &types.Addresses{
								AccountNums: []*types.IdRange{
									{
										Start: 0,
										End: 0,
									},
								},
								ManagerOptions: types.ManagerOptions_Neutral,
						},
						To: &types.Addresses{
							AccountNums: []*types.IdRange{
								{
									Start: 0,
									End: 0,
								},
							},
							ManagerOptions: types.ManagerOptions_Neutral,
						
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
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
