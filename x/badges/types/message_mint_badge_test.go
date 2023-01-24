package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgMintBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgMintBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgMintBadge{
				Creator: "invalid_address",
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 10,
						Amount: 1,
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgMintBadge{
				Creator: sample.AccAddress(),
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 10,
						Amount: 1,
					},
				},
			},
		}, {
			name: "invalid amount",
			msg: types.MsgMintBadge{
				Creator: sample.AccAddress(),
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 10,
						Amount: 0,
					},
				},
			},
			err: types.ErrElementCantEqualThis,
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
