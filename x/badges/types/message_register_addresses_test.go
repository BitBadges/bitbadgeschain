package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRegisterAddresses_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRegisterAddresses
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRegisterAddresses{
				Creator:             "invalid_address",
				AddressesToRegister: []string{},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgRegisterAddresses{
				Creator:             sample.AccAddress(),
				AddressesToRegister: []string{sample.AccAddress()},
			},
		},
		{
			name: "invalid address",
			msg: MsgRegisterAddresses{
				Creator:             sample.AccAddress(),
				AddressesToRegister: []string{},
			},
			err: sdkerrors.ErrInvalidAddress,
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
