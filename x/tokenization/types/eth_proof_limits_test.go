package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestETHProofLimits(t *testing.T) {
	for _, proofs := range [][]*ETHSignatureProof{
		make([]*ETHSignatureProof, 101),
		{nil},
		{{Nonce: strings.Repeat("x", 257), Signature: "00"}},
		{{Nonce: "ticket", Signature: strings.Repeat("0", 133)}},
	} {
		require.Error(t, ValidateETHSignatureProofs(proofs))
	}
	require.NoError(t, ValidateETHSignatureProofs([]*ETHSignatureProof{{Nonce: "ticket-123", Signature: strings.Repeat("00", 65)}}))
}
