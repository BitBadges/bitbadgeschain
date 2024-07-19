package types_test

import (
	"testing"

	"bitbadgeschain/testutil/sample"
	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferBadges_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgTransferBadges
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgTransferBadges{
				Creator: "invalid_address",

				CollectionId: sdkmath.NewUint(1),
				Transfers: []*types.Transfer{
					{
						From: sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.UintRange{
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
			msg: types.MsgTransferBadges{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From: sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.UintRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
								OwnershipTimes: []*types.UintRange{
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
			msg: types.MsgTransferBadges{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From: sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								BadgeIds: []*types.UintRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(1),
									},
								},
								OwnershipTimes: []*types.UintRange{
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
			msg: types.MsgTransferBadges{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From:        sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.UintRange{
									{
										Start: sdkmath.NewUint(10),
										End:   sdkmath.NewUint(1),
									},
								},
								OwnershipTimes: []*types.UintRange{
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
			msg: types.MsgTransferBadges{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From:        sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								BadgeIds: []*types.UintRange{
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
			msg: types.MsgTransferBadges{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From:        sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								OwnershipTimes: []*types.UintRange{
									{
										Start: sdkmath.NewUint(10),
										End:   sdkmath.NewUint(1),
									},
								},
								BadgeIds: []*types.UintRange{
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
