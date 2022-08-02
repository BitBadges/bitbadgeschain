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
				Supply:  10,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supply:  10,
			},
		}, {
			name: "invalid supply",
			msg: MsgNewSubBadge{
				Creator: sample.AccAddress(),
				Supply:  0,
			},
			err: ErrSupplyEqualsZero,
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
