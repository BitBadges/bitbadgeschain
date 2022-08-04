package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func TestMsgNewSubBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgNewSubBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgNewSubBadge{
				Creator:         "invalid_address",
				Supplys:         []uint64{10},
				AmountsToCreate: []uint64{1},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgNewSubBadge{
				Creator:         sample.AccAddress(),
				Supplys:         []uint64{10},
				AmountsToCreate: []uint64{1},
			},
		},  {
			name: "invalid amount",
			msg: types.MsgNewSubBadge{
				Creator:         sample.AccAddress(),
				Supplys:         []uint64{10},
				AmountsToCreate: []uint64{0},
			},
			err: types.ErrAmountEqualsZero,
		}, {
			name: "mismatching lengths",
			msg: types.MsgNewSubBadge{
				Creator:         sample.AccAddress(),
				Supplys:         []uint64{10, 2},
				AmountsToCreate: []uint64{0},
			},
			err: types.ErrInvalidSupplyAndAmounts,
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
