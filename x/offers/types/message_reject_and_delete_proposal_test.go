package types

import (
	"testing"

	"bitbadgeschain/testutil/sample"

	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRejectAndDeleteProposal_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRejectAndDeleteProposal
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRejectAndDeleteProposal{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgRejectAndDeleteProposal{
				Creator: sample.AccAddress(),
				Id:      sdkmath.NewUint(1),
			},
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
