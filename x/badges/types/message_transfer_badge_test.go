package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgTransferBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgTransferBadge{
				Creator:      "invalid_address",
				
				CollectionId: sdkmath.NewUint(1),
				Transfers: []*types.Transfer{
					{
						From:         sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},

							},
						},
					},
				},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				
				Transfers: []*types.Transfer{
					{
						From:    sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
								Times:  []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
							},
						},
					},
				},
			},
		}, {
			name: "invalid amounts",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				
				Transfers: []*types.Transfer{
					{
						From:    sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
								Times:  []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
							},
						},
					},
				},
			},
			err: types.ErrAmountEqualsZero,
		},
		{
			name: "invalid badge range",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				
				Transfers: []*types.Transfer{
					{
						From:    sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(10),
										End:   sdkmath.NewUint(1),
									},
								},
								Times:  []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},

							},
						},
					},
				},
			},
			err: types.ErrStartGreaterThanEnd,
		},
		{
			name: "invalid times",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				
				Transfers: []*types.Transfer{
					{
						From:    sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(10),
										End:   sdkmath.NewUint(1),
									},
								},
								

							},
						},
					},
				},
			},
			err: types.ErrStartGreaterThanEnd,
		},
		{
			name: "invalid times 2",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				
				Transfers: []*types.Transfer{
					{
						From:    sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								Times: []*types.IdRange{
									{
										Start: sdkmath.NewUint(10),
										End:   sdkmath.NewUint(1),
									},
								},
								BadgeIds: []*types.IdRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
								

							},
						},
					},
				},
			},
			err: types.ErrStartGreaterThanEnd,
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
