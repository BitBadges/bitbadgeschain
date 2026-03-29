package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", NewUintFromString("1"), sdkmath.NewUint(0))
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgUpdateDynamicStore("", NewUintFromString("1"), sdkmath.NewUint(0))
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", NewUintFromString("0"), sdkmath.NewUint(0))
	err = msg.ValidateBasic()
	require.Error(t, err)

	// Test with numeric value > 1
	msg = NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", NewUintFromString("1"), sdkmath.NewUint(50))
	require.NoError(t, msg.ValidateBasic())
}
