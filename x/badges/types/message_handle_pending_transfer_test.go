package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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

				NonceRanges: []*types.IdRange{
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
				NonceRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
				Actions: []uint64{1},
			},
		}, {
			name: "invalid subbadge range",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.IdRange{
					{
						Start: 10,
						End:   1,
					},
				},
				Actions: []uint64{1},
			},
			err: types.ErrStartGreaterThanEnd,
		},
		{
			name: "actions length == 1",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.IdRange{
					{
						Start: 0,
						End:   1,
					},
					{
						Start: 2,
						End:   3,
					},
				},
				Actions: []uint64{1},
			},
		},
		{
			name: "actions length == nonceRanges length",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.IdRange{
					{
						Start: 0,
						End:   1,
					},
					{
						Start: 2,
						End:   3,
					},
				},
				Actions: []uint64{1, 2},
			},
		},
		{
			name: "actions length != nonceRanges length",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.IdRange{
					{
						Start: 0,
						End:   1,
					},
					{
						Start: 2,
						End:   3,
					},
					{
						Start: 5,
						End:   6,
					},
				},
				Actions: []uint64{1, 2},
			},
			err: types.ErrActionsLengthNotEqualToRangesLength,
		},
		{
			name: "actions > 2",
			msg: types.MsgHandlePendingTransfer{
				Creator: sample.AccAddress(),
				NonceRanges: []*types.IdRange{
					{
						Start: 0,
						End:   1,
					},
					{
						Start: 2,
						End:   3,
					},
					{
						Start: 5,
						End:   6,
					},
				},
				Actions: []uint64{1, 2, 10},
			},
			err: types.ErrActionOutOfRange,
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
