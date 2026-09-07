package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestETHSignatureMessageSeparatesFieldsAndChain(t *testing.T) {
	first := ETHSignatureChallengeMessage("bitbadges-1", "nonce", "initiator", "1", "", "collection", "a-b", "c")
	second := ETHSignatureChallengeMessage("bitbadges-1", "nonce", "initiator", "1", "", "collection", "a", "b-c")
	require.NotEqual(t, first, second)
	require.NotEqual(t, first, ETHSignatureChallengeMessage("bitbadges-2", "nonce", "initiator", "1", "", "collection", "a-b", "c"))
	require.Equal(t, "BitBadges ETH Signature Challenge v2\n11:bitbadges-15:nonce9:initiator1:10:10:collection3:a-b1:c", first)
	require.Error(t, ValidateETHSignatureProofs([]*ETHSignatureProof{{Nonce: "ambiguous:nonce", Signature: "00"}}))
}
