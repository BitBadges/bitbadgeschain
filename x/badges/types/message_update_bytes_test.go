package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
)

func TestMsgUpdateBytes_ValidateBasic(t *testing.T) {
	var arr []byte
	for i := 0; i <= 260; i++ {
		arr = append(arr, byte(i))
	}

	tests := []struct {
		name string
		msg  MsgUpdateBytes
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateBytes{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateBytes{
				Creator: sample.AccAddress(),
			},
		},
		{
			name: "invalid bytes",
			msg: MsgUpdateBytes{
				Creator: sample.AccAddress(),
				NewBytes: arr,
			},
			err: ErrBytesGreaterThan256,
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
