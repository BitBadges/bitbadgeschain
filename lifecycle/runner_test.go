package lifecycle

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

const alice = "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
const bob = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

func TestIsolatedExecutionAndRollback(t *testing.T) {
	input := []byte(`{"version":1,"id":"bank-rollback","timeMs":"1700000000000","actors":[{"name":"alice","address":"` + alice + `","coins":[{"denom":"ubadge","amount":"10"}]},{"name":"bob","address":"` + bob + `","coins":[]}],"steps":[{"id":"send","actor":"alice","message":{"typeUrl":"/cosmos.bank.v1beta1.MsgSend","value":{"from_address":"` + alice + `","to_address":"` + bob + `","amount":[{"denom":"ubadge","amount":"4"}]}},"expect":{"success":true},"assertions":[{"kind":"coins","address":"` + bob + `","denom":"ubadge","expected":"4"}]},{"id":"reject","actor":"alice","message":{"typeUrl":"/cosmos.bank.v1beta1.MsgSend","value":{"from_address":"` + alice + `","to_address":"` + bob + `","amount":[{"denom":"ubadge","amount":"7"}]}},"expect":{"success":false,"errorContains":"insufficient"},"assertions":[{"kind":"coins","address":"` + alice + `","denom":"ubadge","expected":"6"},{"kind":"coinDelta","address":"` + bob + `","denom":"ubadge","expected":"0"}]}]}`)
	first, err := RunJSON(input)
	require.NoError(t, err)
	require.True(t, first.Passed, "%+v", first)
	second, err := RunJSON(input)
	require.NoError(t, err)
	require.Equal(t, first, second)
	var raw map[string]any
	require.NoError(t, json.Unmarshal(input, &raw))
	raw["steps"].([]any)[0].(map[string]any)["actor"] = "bob"
	changed, _ := json.Marshal(raw)
	_, err = RunJSON(changed)
	require.ErrorContains(t, err, "signer")
}

func TestRejectInvalidScenarios(t *testing.T) {
	for _, input := range []string{`{}`, `{"version":2}`, `{"version":1,"id":"empty","timeMs":"1","actors":[],"steps":[]}`, `{"version":1,"extra":true}`} {
		_, err := RunJSON([]byte(input))
		require.Error(t, err)
	}
}

func TestRejectUnsupportedAndVacuousAssertions(t *testing.T) {
	for _, a := range []Assertion{{Kind: "unknown", Expected: raw(true)}, {Kind: "coins", Address: alice, Denom: "ubadge"}, {Kind: "balances", Address: alice, CollectionID: "-1", Expected: raw([]any{})}, {Kind: "balanceAt", Address: alice, CollectionID: "1", TokenID: "now", OwnershipTime: "now", Expected: raw("0")}} {
		require.Error(t, validateAssertion(a))
	}
	var e Expect
	require.Error(t, json.Unmarshal([]byte(`{}`), &e))
}
func TestTimeAdvancePreservesExactMilliseconds(t *testing.T) {
	s := scenario("clock", Step{ID: "advance", AdvanceTimeMs: "1", Assertions: []Assertion{aCoins(alice, "0")}})
	r := runScenario(t, s)
	require.Equal(t, "1893456000001", r.Steps[0].TimeMs)
}
