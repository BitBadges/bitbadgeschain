package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgCreateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0))
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgCreateDynamicStore("", sdkmath.NewUint(0))
	err := msg.ValidateBasic()
	require.Error(t, err)
}
