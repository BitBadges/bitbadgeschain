package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteCollection{
				Creator:      "invalid_address",
				CollectionId: sdk.NewUint(1),
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdk.NewUint(1),
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
