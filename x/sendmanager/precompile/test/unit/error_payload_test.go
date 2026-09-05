package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

// TestUnmarshalErrorDoesNotEchoInput: a rejected message is reported by its
// cause, not by repeating the caller's input in the error (which is also
// what every node logs).
func TestUnmarshalErrorDoesNotEchoInput(t *testing.T) {
	ts := helpers.NewTestSuite(t)
	const marker = "ZZMARKERZZ"
	method := ts.Precompile.ABI.Methods["send"]
	input, err := helpers.PackMethodCall(&method, `{"to_address":"`+ts.Bob.String()+`","amount":"`+marker+`"}`)
	require.NoError(t, err)

	contract := ts.CreateMockContract(ts.AliceEVM, input)
	_, err = ts.Precompile.Execute(ts.Ctx, contract, false)
	require.Error(t, err)
	require.NotContains(t, err.Error(), marker)
}
