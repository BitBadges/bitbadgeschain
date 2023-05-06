package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateBytes_ValidateBasic(t *testing.T) {
	var arr []byte
	for i := 0; i <= 260; i++ {
		arr = append(arr, byte(i))
	}

	tests := []struct {
		name string
		msg  types.MsgUpdateBytes
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateBytes{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateBytes{
				Creator: sample.AccAddress(),
			},
		},
		{
			name: "invalid bytes",
			msg: types.MsgUpdateBytes{
				Creator:  sample.AccAddress(),
				Bytes: string(arr),
			},
			err: types.ErrBytesGreaterThan256,
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
