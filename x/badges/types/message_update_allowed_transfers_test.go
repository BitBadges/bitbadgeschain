package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateAllowedTransfers_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateAllowedTransfers
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateAllowedTransfers{
				Creator: "invalid_address",
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							Addresses: []string{
								"invalid_address",
							},
							ManagerOptions: uint64(types.AddressOptions_None),
							IncludeOnlySpecified: true,
						},
						To: &types.AddressesMapping{
							Addresses: []string{
								"invalid_address",
							},
							ManagerOptions: uint64(types.AddressOptions_None),
							IncludeOnlySpecified: true,
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateAllowedTransfers{
				Creator: sample.AccAddress(),
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							Addresses: []string{
								sample.AccAddress(),
							},
							ManagerOptions: uint64(types.AddressOptions_None),
							IncludeOnlySpecified: true,
						},
						To: &types.AddressesMapping{
							Addresses: []string{
								sample.AccAddress(),
							},
							ManagerOptions: uint64(types.AddressOptions_None),
							IncludeOnlySpecified: true,
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
