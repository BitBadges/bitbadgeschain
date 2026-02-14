package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgMintAndDistributeTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUniversalUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgUniversalUpdateCollection{
				Creator:       sample.AccAddress(),
				CollectionId:  sdkmath.NewUint(1),
				ValidTokenIds: []*types.UintRange{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
