package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgIncrementStoreValue_ValidateBasic(t *testing.T) {
	msg := NewMsgIncrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5))
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgIncrementStoreValue("", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5))
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgIncrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(5))
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgIncrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "", sdkmath.NewUint(5))
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgIncrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0))
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgIncrementStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "invalid-address", sdkmath.NewUint(5))
	err = msg.ValidateBasic()
	require.Error(t, err)
}
