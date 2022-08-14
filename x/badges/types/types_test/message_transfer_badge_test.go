package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
				Creator:     "invalid_address",
				ToAddresses: []uint64{0},
				Amounts:     []uint64{10},
				From:        1,
				SubbadgeRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgTransferBadge{
				Creator:     sample.AccAddress(),
				ToAddresses: []uint64{0},
				Amounts:     []uint64{10},
				From:        1,
				SubbadgeRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
		}, {
			name: "invalid addresses",
			msg: types.MsgTransferBadge{
				Creator:     sample.AccAddress(),
				ToAddresses: []uint64{0},
				Amounts:     []uint64{10},
				From:        0,
				SubbadgeRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
			err: types.ErrElementCantEqualThis,
		}, {
			name: "invalid amounts",
			msg: types.MsgTransferBadge{
				Creator:     sample.AccAddress(),
				ToAddresses: []uint64{0},
				Amounts:     []uint64{0},
				From:        7,
				SubbadgeRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
			err: types.ErrElementCantEqualThis,
		},
		{
			name: "invalid subbadge range",
			msg: types.MsgTransferBadge{
				Creator:     sample.AccAddress(),
				ToAddresses: []uint64{0},
				Amounts:     []uint64{0},
				From:        7,
				SubbadgeRanges: []*types.IdRange{
					{
						Start: 10,
						End:   1,
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
