package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestETHTrackerKeySeparatesAmbiguousIds(t *testing.T) {
	require.NotEqual(t, ConstructETHSignatureTrackerKey(sdkmath.OneUint(), "", "collection", "a-b", "c", "nonce"), ConstructETHSignatureTrackerKey(sdkmath.OneUint(), "", "collection", "a", "b-c", "nonce"))
}
