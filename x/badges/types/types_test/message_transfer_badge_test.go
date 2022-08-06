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
				Creator: "invalid_address",
				ToAddresses:      []uint64{ 0 },
				Amounts:      []uint64{ 10 },
				From:    1,
				NumberRange: &types.NumberRange{
					Start: 0,
					End: 0,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				ToAddresses:      []uint64{ 0 },
				Amounts:      []uint64{ 10 },
				From:    1,
				NumberRange: &types.NumberRange{
					Start: 0,
					End: 0,
				},
			},
		}, {
			name: "invalid addresses",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				ToAddresses:      []uint64{ 0 },
				Amounts:      []uint64{ 10 },
				From:    0,
				NumberRange: &types.NumberRange{
					Start: 0,
					End: 0,
				},
			},
			err: types.ErrSenderAndReceiverSame,
		},  {
			name: "invalid amounts",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				ToAddresses:      []uint64{ 0 },
				Amounts:      []uint64{ 0 },
				From:    7,
				NumberRange: &types.NumberRange{
					Start: 0,
					End: 0,
				},
			},
			err: types.ErrAmountEqualsZero,
		},
		{
			name: "invalid subbadge range",
			msg: types.MsgTransferBadge{
				Creator: sample.AccAddress(),
				ToAddresses:      []uint64{ 0 },
				Amounts:      []uint64{ 0 },
				From:    7,
				NumberRange: &types.NumberRange{
					Start: 10,
					End: 0,
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
