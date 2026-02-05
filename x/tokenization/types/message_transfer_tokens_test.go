package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgTransferTokens
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgTransferTokens{
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
								TokenIds: []*types.UintRange{
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
			msg: types.MsgTransferTokens{
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
								TokenIds: []*types.UintRange{
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
			msg: types.MsgTransferTokens{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From: sample.AccAddress(),
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								TokenIds: []*types.UintRange{
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
			name: "invalid ID range",
			msg: types.MsgTransferTokens{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From:        sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								TokenIds: []*types.UintRange{
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
			msg: types.MsgTransferTokens{
				Creator: sample.AccAddress(),

				Transfers: []*types.Transfer{
					{
						From:        sample.AccAddress(),
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: sdkmath.NewUint(10),
								TokenIds: []*types.UintRange{
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
			msg: types.MsgTransferTokens{
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
								TokenIds: []*types.UintRange{
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
