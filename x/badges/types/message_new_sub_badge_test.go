package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
)

func TestMsgNewSubBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewSubBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewSubBadge{
				Creator: "invalid_address",
				Supplys:  []uint64{ 10 },
				AmountsToCreate: []uint64{ 1 },
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supplys:  []uint64{ 10 },
				AmountsToCreate: []uint64{ 1 },
			},
		}, {
			name: "invalid supply",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supplys:  []uint64{ 0 },
				AmountsToCreate: []uint64{ 1 },
			},
			err: ErrSupplyEqualsZero,
		}, {
			name: "invalid amount",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supplys:  []uint64{ 10 },
				AmountsToCreate: []uint64{ 0 },
			},
			err: ErrAmountEqualsZero,
		}, {
			name: "mismatching lengths",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supplys:  []uint64{ 10, 2 },
				AmountsToCreate: []uint64{ 0 },
			},
			err: ErrAmountEqualsZero,
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
