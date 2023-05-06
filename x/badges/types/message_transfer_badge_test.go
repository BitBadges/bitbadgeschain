package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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
				Creator: "invalid_address",
				From:   sample.AccAddress(),
				CollectionId: 1,
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: 10,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End:   0,
									},
								},
							},
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				From:    sample.AccAddress(),
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: 10,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End:   0,
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
				From:    sample.AccAddress(),
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{
							sample.AccAddress(),
						},
						Balances: []*types.Balance{
							{
								Amount: 0,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End:   0,
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
				From:    sample.AccAddress(),
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{sample.AccAddress()},
						Balances: []*types.Balance{
							{
								Amount: 10,
								BadgeIds: []*types.IdRange{
									{
										Start: 10,
										End:   1,
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
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
