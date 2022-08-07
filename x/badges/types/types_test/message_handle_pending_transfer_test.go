package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func TestMsgHandlePendingTransfer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgHandlePendingTransfer
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgHandlePendingTransfer{
				Creator: "invalid_address",

				NonceRanges: []*types.NumberRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.NumberRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
		}, {
			name: "invalid subbadge range",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.NumberRange{
					{
						Start: 10,
						End:   0,
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
