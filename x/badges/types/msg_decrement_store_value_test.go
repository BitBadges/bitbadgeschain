package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgDecrementStoreValue_ValidateBasic(t *testing.T) {
	msg := NewMsgDecrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5), true)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgDecrementStoreValue("", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5), true)
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgDecrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5), true)
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgDecrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "", sdkmath.NewUint(5), true)
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgDecrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0), true)
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgDecrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "invalid-address", sdkmath.NewUint(5), true)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
