package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
)

func TestMsgTransferBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTransferBadge{
				Creator: "invalid_address",
				To: 0,
				From: 1,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: MsgTransferBadge{
				Creator: sample.AccAddress(),
				To: 0,
				From: 1,
			},
		}, {
			name: "invalid addresses",
			msg: MsgTransferBadge{
				Creator: sample.AccAddress(),
				To: 0,
				From: 0,
			},
			err: ErrSenderAndReceiverSame,
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
