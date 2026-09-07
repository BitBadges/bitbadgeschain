package tokenization_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
)

// TestUnmarshalErrorDoesNotEchoInput: a rejected message is reported by its
// cause, not by repeating the caller's input in the error (which is also
// what every node logs).
func TestUnmarshalErrorDoesNotEchoInput(t *testing.T) {
	ts := helpers.NewTestSuite()
	const marker = "ZZMARKERZZ"
	method := ts.Precompile.ABI.Methods["transferTokens"]
	input, err := helpers.PackMethodWithJSON(&method, `{"collectionId":"1","unknownField":"`+marker+`"}`)
	require.NoError(t, err)

	_, err = ts.CallPrecompile(ts.AliceEVM, input)
	require.Error(t, err)
	require.NotContains(t, err.Error(), marker)
}
