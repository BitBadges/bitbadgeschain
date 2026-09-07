package gamm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

// TestUnmarshalErrorDoesNotEchoInput: a rejected message is reported by its
// cause, not by repeating the caller's input in the error (which is also
// what every node logs).
func TestUnmarshalErrorDoesNotEchoInput(t *testing.T) {
	ts := helpers.NewTestSuite(t)
	const marker = "ZZMARKERZZ"
	for _, tc := range []struct{ method, json string }{
		{"joinPool", `{"pool_id":"` + marker + `"}`},
		{"getPool", `{"poolId":{"nested":"` + marker + `"}}`},
	} {
		method := ts.Precompile.ABI.Methods[tc.method]
		input, err := helpers.PackMethodWithJSON(&method, tc.json)
		require.NoError(t, err)
		contract := ts.CreateMockContract(ts.AliceEVM, input)
		contract.Input = input

		_, err = ts.Precompile.Execute(ts.Ctx, contract, tc.method == "getPool")
		require.Error(t, err, tc.method)
		require.NotContains(t, err.Error(), marker, tc.method)
	}
}
