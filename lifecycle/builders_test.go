package lifecycle

import (
	sdkmath "cosmossdk.io/math"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const merchant = "bb1zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zql3w7"
const payer = "bb1xvenxvenxvenxvenxvenxvenxvenxvenlrd2nm"
const start = "1893456000000"
const max = "18446744073709551615"

func raw(v any) json.RawMessage {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
func artifact(t *testing.T, id string) *Message {
	t.Helper()
	b, e := os.ReadFile("testdata/artifacts/" + id + ".json")
	require.NoError(t, e)
	var m Message
	require.NoError(t, json.Unmarshal(b, &m))
	var v map[string]any
	require.NoError(t, json.Unmarshal(m.Value, &v))
	require.Equal(t, "", v["creator"])
	require.Equal(t, "", v["manager"])
	v["creator"] = alice
	v["manager"] = alice
	m.Value = raw(v)
	return &m
}
func balance(amount, from, to string) any {
	return []any{map[string]any{"amount": amount, "tokenIds": []any{map[string]string{"start": "1", "end": "1"}}, "ownershipTimes": []any{map[string]string{"start": from, "end": to}}}}
}
func aCoins(address, amount string) Assertion {
	return Assertion{Kind: "coins", Address: address, Denom: "ubadge", Expected: raw(amount)}
}
func aBalances(address string, v any) Assertion {
	return Assertion{Kind: "balances", CollectionID: "1", Address: address, Expected: raw(v)}
}
func aManager() Assertion { return Assertion{Kind: "manager", CollectionID: "1", Expected: raw(alice)} }
func transfer(actor, from, to string, balances any, approval ...string) *Message {
	id := "subscription-tier-1"
	if len(approval) > 0 {
		id = approval[0]
	}
	return &Message{TypeURL: "/tokenization.MsgTransferTokens", Value: raw(map[string]any{"creator": actor, "collectionId": "1", "transfers": []any{map[string]any{"from": from, "toAddresses": []string{to}, "balances": balances, "prioritizedApprovals": []any{map[string]string{"approvalId": id, "approvalLevel": "collection", "approverAddress": "", "version": "0"}}}}})}
}
func msgStep(id, actor string, m *Message, success bool, assertions ...Assertion) Step {
	return Step{ID: id, Actor: actor, Message: m, Expect: &Expect{Success: success}, Assertions: assertions}
}
func scenario(id string, steps ...Step) Scenario {
	return Scenario{Version: 1, ID: id, TimeMs: start, Actors: []Actor{{Name: "manager", Address: alice, Coins: sdk.Coins{}}, {Name: "payer", Address: payer, Coins: sdk.Coins{}}, {Name: "merchant", Address: merchant, Coins: sdk.Coins{}}, {Name: "bob", Address: bob, Coins: sdk.Coins{}}}, Steps: steps}
}
func runScenario(t *testing.T, s Scenario) *Result {
	t.Helper()
	r, e := RunJSON(raw(s))
	require.NoError(t, e)
	if !r.Passed {
		for _, step := range r.Steps {
			if !step.Passed {
				step.Events = nil
				b, _ := json.Marshal(step)
				t.Errorf("failed step: %s", b)
			}
		}
		t.FailNow()
	}
	if os.Getenv("UPDATE_LIFECYCLE_FIXTURES") == "1" {
		require.NoError(t, os.MkdirAll("testdata/scenarios", 0755))
		b, _ := json.MarshalIndent(s, "", "  ")
		require.NoError(t, os.WriteFile("testdata/scenarios/"+s.ID+".json", append(b, 10), 0644))
	}
	return r
}
func TestBuilderArtifactsCreateUnchanged(t *testing.T) {
	for _, id := range []string{"subscription-fixed-price", "invoice-direct-payment", "backed-token", "purchasable-credits", "spendable-service-credits"} {
		t.Run(id, func(t *testing.T) {
			runScenario(t, scenario(id, msgStep("create", "manager", artifact(t, id), true, aManager())))
		})
	}
}
func TestSubscriptionLifecycle(t *testing.T) {
	s := scenario("subscription", msgStep("create", "manager", artifact(t, "subscription-fixed-price"), true, aManager()),
		msgStep("unpaid", "payer", transfer(payer, "Mint", payer, balance("1", start, "1896047999999")), false, aBalances(payer, []any{}), aCoins(merchant, "0")))
	runScenario(t, s)
}

func fund(s *Scenario, address, amount string) {
	for i := range s.Actors {
		if s.Actors[i].Address == address {
			s.Actors[i].Coins = sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewIntFromUint64(1)))
			n, ok := sdkmath.NewIntFromString(amount)
			if !ok {
				panic("amount")
			}
			s.Actors[i].Coins = sdk.NewCoins(sdk.NewCoin("ubadge", n))
		}
	}
}
func TestSubscriptionPaidLifecycle(t *testing.T) {
	s := scenario("subscription-paid", msgStep("create", "manager", artifact(t, "subscription-fixed-price"), true, aManager()),
		msgStep("purchase", "payer", transfer(payer, "Mint", payer, balance("1", start, "1896047999999")), true, aBalances(payer, balance("1", start, "1896047999999")), aCoins(merchant, "5000000000")),
		msgStep("transfer-blocked", "payer", transfer(payer, payer, bob, balance("1", start, "1896047999999")), false, aBalances(bob, []any{})),
		Step{ID: "expiry-inclusive", AdvanceTimeMs: "2591999999", Assertions: []Assertion{{Kind: "balanceAt", CollectionID: "1", Address: payer, TokenID: "1", OwnershipTime: "now", Expected: raw("1")}}},
		Step{ID: "expired", AdvanceTimeMs: "1", Assertions: []Assertion{{Kind: "balanceAt", CollectionID: "1", Address: payer, TokenID: "1", OwnershipTime: "now", Expected: raw("0")}}},
		msgStep("renew", "payer", transfer(payer, "Mint", payer, balance("1", "1896048000000", "1898639999999")), true, aBalances(payer, balance("1", start, "1898639999999")), aCoins(merchant, "10000000000")))
	fund(&s, payer, "20000000000")
	runScenario(t, s)
}

const burn = "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqs7gvmv"
const backing = "bb1f8y7m98kw6yexhrsxse4aqcagfzsvjghgupn8yltwhj5f5tce73s0askw9"

func TestInvoiceLifecycle(t *testing.T) {
	s := scenario("invoice", msgStep("create", "manager", artifact(t, "invoice-direct-payment"), true, aManager()),
		msgStep("wrong-payer", "bob", transfer(bob, "Mint", burn, balance("1", "1", max), "payment-request-pay"), false, aCoins(merchant, "0")),
		msgStep("wrong-recipient", "payer", transfer(payer, "Mint", payer, balance("1", "1", max), "payment-request-pay"), false, aCoins(merchant, "0")),
		msgStep("settle", "payer", transfer(payer, "Mint", burn, balance("1", "1", max), "payment-request-pay"), true, aCoins(merchant, "7000000000"), aCoins(payer, "12993000000")),
		msgStep("duplicate", "payer", transfer(payer, "Mint", burn, balance("1", "1", max), "payment-request-pay"), false, aCoins(merchant, "7000000000"), aCoins(payer, "12993000000")))
	fund(&s, payer, "20000000000")
	fund(&s, bob, "20000000000")
	runScenario(t, s)
}
func TestBackedLifecycle(t *testing.T) {
	s := scenario("backed", msgStep("create", "manager", artifact(t, "backed-token"), true, aManager()),
		msgStep("deposit", "payer", transfer(payer, backing, payer, balance("10", "1", max), "smart-token-deposit"), true, aCoins(payer, "90"), aCoins(backing, "10"), aBalances(payer, balance("10", "1", max))),
		msgStep("withdraw", "payer", transfer(payer, payer, backing, balance("4", "1", max), "smart-token-withdraw"), true, aCoins(payer, "94"), aCoins(backing, "6"), aBalances(payer, balance("6", "1", max))),
		msgStep("over-withdraw", "payer", transfer(payer, payer, backing, balance("7", "1", max), "smart-token-withdraw"), false, aCoins(payer, "94"), aCoins(backing, "6"), aBalances(payer, balance("6", "1", max))),
		msgStep("unbacked", "bob", transfer(bob, backing, bob, balance("1", "1", max), "smart-token-deposit"), false, aBalances(bob, []any{})))
	fund(&s, payer, "100")
	runScenario(t, s)
}
func TestNFTAndFungibleLifecycles(t *testing.T) {
	for _, tc := range []struct{ id, amount string }{{"nft", "1"}, {"fungible", "10"}} {
		t.Run(tc.id, func(t *testing.T) {
			s := scenario(tc.id,
				msgStep("create", "manager", artifact(t, tc.id), true, aManager()),
				msgStep("unauthorized-mint", "bob", transfer(bob, "Mint", payer, balance(tc.amount, "1", max), "manager-mint"), false, aBalances(payer, []any{})),
				msgStep("mint-cap", "manager", transfer(alice, "Mint", payer, balance(tc.amount, "1", max), "manager-mint"), true, aBalances(payer, balance(tc.amount, "1", max))),
				msgStep("cap-exhausted", "manager", transfer(alice, "Mint", payer, balance("1", "1", max), "manager-mint"), false, aBalances(payer, balance(tc.amount, "1", max))),
				msgStep("transfer", "payer", transfer(payer, payer, bob, balance(tc.amount, "1", max), "transferable-approval"), true, aBalances(payer, []any{}), aBalances(bob, balance(tc.amount, "1", max))),
				msgStep("burn", "bob", transfer(bob, bob, burn, balance(tc.amount, "1", max), "transferable-approval"), true, aBalances(bob, []any{})),
				msgStep("no-remint-after-burn", "manager", transfer(alice, "Mint", payer, balance("1", "1", max), "manager-mint"), false, aBalances(payer, []any{})),
				msgStep("unauthorized-update", "bob", &Message{TypeURL: "/tokenization.MsgUpdateCollection", Value: raw(map[string]any{"creator": bob, "collectionId": "1", "updateManager": true, "manager": bob})}, false, aManager()),
				msgStep("manager-update", "manager", &Message{TypeURL: "/tokenization.MsgUpdateCollection", Value: raw(map[string]any{"creator": alice, "collectionId": "1", "updateManager": true, "manager": bob})}, true, Assertion{Kind: "manager", CollectionID: "1", Expected: raw(bob)}))
			runScenario(t, s)
		})
	}
}
func TestSubscriptionConsentAndLockedTerms(t *testing.T) {
	m := artifact(t, "subscription-fixed-price")
	var value map[string]any
	require.NoError(t, json.Unmarshal(m.Value, &value))
	approvals := value["collectionApprovals"].([]any)
	approvals[0].(map[string]any)["approvalCriteria"].(map[string]any)["coinTransfers"] = []any{}
	consent := func(allowed bool) *Message {
		return &Message{TypeURL: "/tokenization.MsgUpdateUserApprovals", Value: raw(map[string]any{"creator": payer, "collectionId": "1", "updateAutoApproveAllIncomingTransfers": true, "autoApproveAllIncomingTransfers": allowed, "updateAutoApproveSelfInitiatedIncomingTransfers": true, "autoApproveSelfInitiatedIncomingTransfers": allowed})}
	}
	s := scenario("subscription-consent-and-locks", msgStep("create", "manager", m, true, aManager()),
		msgStep("deny-incoming", "payer", consent(false), true, aBalances(payer, []any{})),
		msgStep("consent-required", "payer", transfer(payer, "Mint", payer, balance("1", start, "1896047999999")), false, aBalances(payer, []any{}), aCoins(merchant, "0")),
		msgStep("price-bypass", "manager", &Message{TypeURL: "/tokenization.MsgUpdateCollection", Value: raw(map[string]any{"creator": alice, "collectionId": "1", "updateCollectionApprovals": true, "collectionApprovals": approvals})}, false, aCoins(merchant, "0")),
		msgStep("unlock-permissions", "manager", &Message{TypeURL: "/tokenization.MsgUpdateCollection", Value: raw(map[string]any{"creator": alice, "collectionId": "1", "updateCollectionPermissions": true, "collectionPermissions": map[string]any{}})}, false, aManager()),
		msgStep("allow-incoming", "payer", consent(true), true, aBalances(payer, []any{})),
		msgStep("paid-purchase", "payer", transfer(payer, "Mint", payer, balance("1", start, "1896047999999")), true, aBalances(payer, balance("1", start, "1896047999999")), aCoins(merchant, "5000000000")))
	fund(&s, payer, "20000000000")
	runScenario(t, s)
}

func TestApprovalPermissionAndTrackerAssertions(t *testing.T) {
	m := artifact(t, "invoice-direct-payment")
	var create types.MsgCreateCollection
	require.NoError(t, json.Unmarshal(m.Value, &create))
	for _, approval := range create.CollectionApprovals {
		approval.ApprovalCriteria.ApprovalAmounts = &types.ApprovalAmounts{OverallApprovalAmount: sdkmath.NewUint(0), PerToAddressApprovalAmount: sdkmath.NewUint(0), PerFromAddressApprovalAmount: sdkmath.NewUint(0), PerInitiatedByAddressApprovalAmount: sdkmath.NewUint(0), ResetTimeIntervals: &types.ResetTimeIntervals{StartTime: sdkmath.NewUint(0), IntervalLength: sdkmath.NewUint(0)}}
		approval.ApprovalCriteria.AutoDeletionOptions = &types.AutoDeletionOptions{}
		approval.ApprovalCriteria.MustPrioritize = true
	}
	s := scenario("invoice-evidence",
		msgStep("create", "manager", m, true, Assertion{Kind: "approvals", CollectionID: "1", Expected: raw(create.CollectionApprovals)}, Assertion{Kind: "permissions", CollectionID: "1", Expected: raw(create.CollectionPermissions)}),
		msgStep("settle", "payer", transfer(payer, "Mint", burn, balance("1", "1", max), "payment-request-pay"), true, Assertion{Kind: "tracker", CollectionID: "1", ApprovalID: "payment-request-pay", AmountTrackerID: "payment-request-pay-tracker", Level: "collection", TrackerType: "overall", Expected: raw(map[string]any{"numTransfers": "1", "lastUpdatedAt": start})}))
	fund(&s, payer, "20000000000")
	runScenario(t, s)
}

func TestFailedBatchRollsBackTokensAndTrackers(t *testing.T) {
	m := transfer(alice, "Mint", payer, balance("1", "1", max), "manager-mint")
	var value map[string]any
	require.NoError(t, json.Unmarshal(m.Value, &value))
	first := value["transfers"].([]any)[0]
	value["transfers"] = []any{first, first}
	m.Value = raw(value)
	s := scenario("batch-rollback", msgStep("create", "manager", artifact(t, "nft"), true, aManager()),
		msgStep("partial-mint-rejected", "manager", m, false, aBalances(payer, []any{})),
		msgStep("cap-not-consumed", "manager", transfer(alice, "Mint", payer, balance("1", "1", max), "manager-mint"), true, aBalances(payer, balance("1", "1", max))))
	runScenario(t, s)
}

func TestCheckedInScenarios(t *testing.T) {
	entries, err := os.ReadDir("testdata/scenarios")
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile("testdata/scenarios/" + entry.Name())
			require.NoError(t, err)
			result, err := RunJSON(data)
			require.NoError(t, err)
			require.True(t, result.Passed, "scenario %s", entry.Name())
		})
	}
}

func TestInvoiceDeadlineBoundary(t *testing.T) {
	for _, tc := range []struct {
		time           string
		success        bool
		merchantAmount string
	}{{"1896048000000", true, "7000000000"}, {"1896048000001", false, "0"}} {
		t.Run(tc.time, func(t *testing.T) {
			s := scenario("invoice-deadline-"+tc.time, msgStep("create", "manager", artifact(t, "invoice-direct-payment"), true, aManager()), msgStep("pay", "payer", transfer(payer, "Mint", burn, balance("1", "1", max), "payment-request-pay"), tc.success, aCoins(merchant, tc.merchantAmount)))
			s.TimeMs = tc.time
			fund(&s, payer, "20000000000")
			runScenario(t, s)
		})
	}
}

func TestCreditPurchaseLifecycle(t *testing.T) {
	s := scenario("credits-purchase", msgStep("create", "manager", artifact(t, "purchasable-credits"), true, aManager()),
		msgStep("buy-two-units", "payer", transfer(payer, "Mint", payer, balance("74", "1", max), "credit-scaled"), true, aBalances(payer, balance("74", "1", max)), aCoins(merchant, "2")),
		msgStep("non-multiple", "payer", transfer(payer, "Mint", payer, balance("38", "1", max), "credit-scaled"), false, aBalances(payer, balance("74", "1", max)), aCoins(merchant, "2")),
		msgStep("cannot-transfer", "payer", transfer(payer, payer, bob, balance("1", "1", max), "credit-scaled"), false, aBalances(bob, []any{})),
		msgStep("cannot-burn", "payer", transfer(payer, payer, burn, balance("1", "1", max), "credit-scaled"), false, aBalances(payer, balance("74", "1", max))))
	fund(&s, payer, "100")
	runScenario(t, s)
}
func TestSpendableCreditsLifecycle(t *testing.T) {
	s := scenario("spendable-credits", msgStep("create", "manager", artifact(t, "spendable-service-credits"), true, aManager()),
		msgStep("buy-three-packs", "bob", transfer(bob, "Mint", bob, balance("30", "1", max), "spendable-purchase"), true, aBalances(bob, balance("30", "1", max)), aCoins(payer, "3000000")),
		msgStep("over-purchase-cap", "bob", transfer(bob, "Mint", bob, balance("40", "1", max), "spendable-purchase"), false, aBalances(bob, balance("30", "1", max)), aCoins(payer, "3000000")),
		msgStep("provider-cannot-consume", "payer", transfer(payer, bob, burn, balance("1", "1", max), "spendable-consume"), false, aBalances(bob, balance("30", "1", max))),
		msgStep("holder-consumes", "bob", transfer(bob, bob, burn, balance("10", "1", max), "spendable-consume"), true, aBalances(bob, balance("20", "1", max))),
		Step{ID: "after-expiry", AdvanceTimeMs: "2592000001", Assertions: []Assertion{aBalances(bob, balance("20", "1", max))}},
		msgStep("expired-consumption", "bob", transfer(bob, bob, burn, balance("1", "1", max), "spendable-consume"), false, aBalances(bob, balance("20", "1", max))),
		msgStep("expired-purchase", "bob", transfer(bob, "Mint", bob, balance("10", "1", max), "spendable-purchase"), false, aCoins(payer, "3000000")))
	fund(&s, bob, "10000000")
	runScenario(t, s)
}

func TestArtifactProvenance(t *testing.T) {
	data, err := os.ReadFile("testdata/provenance.json")
	require.NoError(t, err)
	var manifest struct {
		Artifacts map[string]string `json:"artifacts"`
	}
	require.NoError(t, json.Unmarshal(data, &manifest))
	require.Len(t, manifest.Artifacts, 7)
	for name, expected := range manifest.Artifacts {
		data, err := os.ReadFile("testdata/artifacts/" + name)
		require.NoError(t, err)
		sum := sha256.Sum256(data)
		require.Equal(t, expected, hex.EncodeToString(sum[:]), name)
	}
}
