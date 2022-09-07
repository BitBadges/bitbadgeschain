package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgFreezeAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgFreezeAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgFreezeAddress{
				Creator: "invalid_address",
				AddressRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgFreezeAddress{
				Creator: sample.AccAddress(),
				AddressRanges: []*types.IdRange{
					{
						Start: 0,
						End:   0,
					},
				},
			},
		}, {
			name: "invalid subbadge range",
			msg: types.MsgFreezeAddress{
				Creator: sample.AccAddress(),
				AddressRanges: []*types.IdRange{
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
