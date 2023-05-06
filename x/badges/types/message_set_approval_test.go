package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgSetApproval_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgSetApproval
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgSetApproval{
				Creator: "invalid_address",
				CollectionId: 1,
				Balances: []*types.Balance{
					{
						Amount: 1,
						BadgeIds: []*types.IdRange{
							{
								Start: 0,
								End:   0,
							},
						},
					},
				},
				Address: sample.AccAddress(),
				
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgSetApproval{
				CollectionId: 1,
				Creator: sample.AccAddress(),
				Balances: []*types.Balance{
					{
						Amount: 1,
						BadgeIds: []*types.IdRange{
							{
								Start: 0,
								End:   0,
							},
						},
					},
				},
				Address: sample.AccAddress(),
			},
		}, {
			name: "invalid badgeId range",
			msg: types.MsgSetApproval{
				Creator: sample.AccAddress(),
				CollectionId: 1,
				Balances: []*types.Balance{
					{
						Amount: 1,
						BadgeIds: []*types.IdRange{
							{
								Start: 10,
								End:   1,
							},
						},
					},
				},
				Address: sample.AccAddress(),
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
