package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
)

// TestParseTransferTokensMemo tests parsing of transfer_tokens memo
func TestParseTransferTokensMemo(t *testing.T) {
	t.Run("valid transfer_tokens memo", func(t *testing.T) {
		memo := `{
			"transfer_tokens": {
				"collection_id": "123",
				"transfers": [{
					"from": "Mint",
					"to_addresses": ["bb1recipient"],
					"balances": [{
						"amount": "1",
						"ownership_times": [{"start": "1", "end": "18446744073709551615"}],
						"token_ids": [{"start": "1", "end": "1"}]
					}],
					"merkle_proofs": [],
					"eth_signature_proofs": [],
					"prioritized_approvals": [],
					"only_check_prioritized_collection_approvals": false,
					"only_check_prioritized_incoming_approvals": false,
					"only_check_prioritized_outgoing_approvals": false
				}],
				"fail_on_error": true,
				"recover_address": ""
			}
		}`

		hookData, err := customhookstypes.ParseHookDataFromMemo(memo)
		require.NoError(t, err)
		require.NotNil(t, hookData)
		require.Nil(t, hookData.SwapAndAction)
		require.NotNil(t, hookData.TransferTokens)
		require.Equal(t, "123", hookData.TransferTokens.CollectionId)
		require.True(t, hookData.TransferTokens.FailOnError)
		require.NotNil(t, hookData.TransferTokens.Transfers)
	})

	t.Run("mutual exclusivity — both keys present", func(t *testing.T) {
		memo := `{
			"swap_and_action": {"post_swap_action": {}},
			"transfer_tokens": {"collection_id": "1", "transfers": []}
		}`

		hookData, err := customhookstypes.ParseHookDataFromMemo(memo)
		require.Error(t, err)
		require.Nil(t, hookData)
		require.Contains(t, err.Error(), "swap_and_action and transfer_tokens")
	})

	t.Run("no hook data — returns nil", func(t *testing.T) {
		memo := `{"some_other_key": "value"}`
		hookData, err := customhookstypes.ParseHookDataFromMemo(memo)
		require.NoError(t, err)
		require.Nil(t, hookData)
	})

	t.Run("empty memo — returns nil", func(t *testing.T) {
		hookData, err := customhookstypes.ParseHookDataFromMemo("")
		require.NoError(t, err)
		require.Nil(t, hookData)
	})

	t.Run("swap_and_action only — still works", func(t *testing.T) {
		memo := `{
			"swap_and_action": {
				"user_swap": {"swap_exact_asset_in": {"operations": [{"pool": "1", "denom_in": "uatom", "denom_out": "ubadge"}]}},
				"min_asset": {"native": {"denom": "ubadge", "amount": "1"}},
				"post_swap_action": {"transfer": {"to_address": "bb1abc"}}
			}
		}`
		hookData, err := customhookstypes.ParseHookDataFromMemo(memo)
		require.NoError(t, err)
		require.NotNil(t, hookData)
		require.NotNil(t, hookData.SwapAndAction)
		require.Nil(t, hookData.TransferTokens)
	})
}

// TestUnmarshalTransfersFromJSON tests the snake_case JSON → protobuf conversion
func TestUnmarshalTransfersFromJSON(t *testing.T) {
	t.Run("basic transfer", func(t *testing.T) {
		transfersJSON := json.RawMessage(`[{
			"from": "Mint",
			"to_addresses": ["bb1recipient"],
			"balances": [{
				"amount": "100",
				"ownership_times": [{"start": "1", "end": "18446744073709551615"}],
				"token_ids": [{"start": "1", "end": "10"}]
			}],
			"merkle_proofs": [],
			"eth_signature_proofs": [],
			"prioritized_approvals": [],
			"only_check_prioritized_collection_approvals": false,
			"only_check_prioritized_incoming_approvals": false,
			"only_check_prioritized_outgoing_approvals": false
		}]`)

		proto, err := customhookstypes.UnmarshalTransfersFromJSON(transfersJSON)
		require.NoError(t, err)
		require.Len(t, proto, 1)
		require.Equal(t, "Mint", proto[0].From)
		require.Equal(t, []string{"bb1recipient"}, proto[0].ToAddresses)
		require.Len(t, proto[0].Balances, 1)
		require.Equal(t, "100", proto[0].Balances[0].Amount.String())
		require.Len(t, proto[0].Balances[0].TokenIds, 1)
		require.Equal(t, "1", proto[0].Balances[0].TokenIds[0].Start.String())
		require.Equal(t, "10", proto[0].Balances[0].TokenIds[0].End.String())
	})

	t.Run("invalid JSON — error", func(t *testing.T) {
		_, err := customhookstypes.UnmarshalTransfersFromJSON(json.RawMessage(`not json`))
		require.Error(t, err)
	})

	t.Run("with merkle proofs", func(t *testing.T) {
		transfersJSON := json.RawMessage(`[{
			"from": "Mint",
			"to_addresses": ["bb1recipient"],
			"balances": [{
				"amount": "1",
				"ownership_times": [{"start": "1", "end": "1"}],
				"token_ids": [{"start": "1", "end": "1"}]
			}],
			"merkle_proofs": [{
				"leaf": "abc123",
				"aunts": [
					{"aunt": "def456", "on_right": true},
					{"aunt": "ghi789", "on_right": false}
				],
				"leaf_signature": "sig"
			}],
			"eth_signature_proofs": [],
			"prioritized_approvals": []
		}]`)

		proto, err := customhookstypes.UnmarshalTransfersFromJSON(transfersJSON)
		require.NoError(t, err)
		require.Len(t, proto[0].MerkleProofs, 1)
		require.Equal(t, "abc123", proto[0].MerkleProofs[0].Leaf)
		require.Len(t, proto[0].MerkleProofs[0].Aunts, 2)
		require.True(t, proto[0].MerkleProofs[0].Aunts[0].OnRight)
		require.False(t, proto[0].MerkleProofs[0].Aunts[1].OnRight)
		require.Equal(t, "sig", proto[0].MerkleProofs[0].LeafSignature)
	})

	t.Run("with precalculate balances", func(t *testing.T) {
		transfersJSON := json.RawMessage(`[{
			"from": "Mint",
			"to_addresses": ["bb1recipient"],
			"balances": [],
			"precalculate_balances_from_approval": {
				"approval_id": "mint-approval",
				"approval_level": "collection",
				"approver_address": "",
				"version": "1",
				"precalculation_options": {
					"override_timestamp": "1000",
					"token_ids_override": [{"start": "5", "end": "10"}]
				}
			},
			"merkle_proofs": [],
			"eth_signature_proofs": [],
			"prioritized_approvals": []
		}]`)

		proto, err := customhookstypes.UnmarshalTransfersFromJSON(transfersJSON)
		require.NoError(t, err)
		require.NotNil(t, proto[0].PrecalculateBalancesFromApproval)
		require.Equal(t, "mint-approval", proto[0].PrecalculateBalancesFromApproval.ApprovalId)
		require.Equal(t, "collection", proto[0].PrecalculateBalancesFromApproval.ApprovalLevel)
		require.NotNil(t, proto[0].PrecalculateBalancesFromApproval.PrecalculationOptions)
		require.Equal(t, "1000", proto[0].PrecalculateBalancesFromApproval.PrecalculationOptions.OverrideTimestamp.String())
	})

	t.Run("with prioritized approvals", func(t *testing.T) {
		transfersJSON := json.RawMessage(`[{
			"from": "Mint",
			"to_addresses": ["bb1recipient"],
			"balances": [{
				"amount": "1",
				"ownership_times": [{"start": "1", "end": "1"}],
				"token_ids": [{"start": "1", "end": "1"}]
			}],
			"merkle_proofs": [],
			"eth_signature_proofs": [],
			"prioritized_approvals": [{
				"approval_id": "my-approval",
				"approval_level": "collection",
				"approver_address": "",
				"version": "0"
			}],
			"only_check_prioritized_collection_approvals": true,
			"only_check_prioritized_incoming_approvals": false,
			"only_check_prioritized_outgoing_approvals": false
		}]`)

		proto, err := customhookstypes.UnmarshalTransfersFromJSON(transfersJSON)
		require.NoError(t, err)
		require.Len(t, proto[0].PrioritizedApprovals, 1)
		require.Equal(t, "my-approval", proto[0].PrioritizedApprovals[0].ApprovalId)
		require.True(t, proto[0].OnlyCheckPrioritizedCollectionApprovals)
		require.False(t, proto[0].OnlyCheckPrioritizedIncomingApprovals)
		require.False(t, proto[0].OnlyCheckPrioritizedOutgoingApprovals)
	})

	t.Run("with eth signature proofs", func(t *testing.T) {
		transfersJSON := json.RawMessage(`[{
			"from": "Mint",
			"to_addresses": ["bb1recipient"],
			"balances": [{
				"amount": "1",
				"ownership_times": [{"start": "1", "end": "1"}],
				"token_ids": [{"start": "1", "end": "1"}]
			}],
			"merkle_proofs": [],
			"eth_signature_proofs": [{"nonce": "test-nonce", "signature": "0xabc"}],
			"prioritized_approvals": []
		}]`)

		proto, err := customhookstypes.UnmarshalTransfersFromJSON(transfersJSON)
		require.NoError(t, err)
		require.Len(t, proto[0].EthSignatureProofs, 1)
		require.Equal(t, "test-nonce", proto[0].EthSignatureProofs[0].Nonce)
		require.Equal(t, "0xabc", proto[0].EthSignatureProofs[0].Signature)
	})
}

// TestTransferTokensMemoJSONRoundtrip tests that the memo JSON format roundtrips correctly
func TestTransferTokensMemoJSONRoundtrip(t *testing.T) {
	memo := `{
		"transfer_tokens": {
			"collection_id": "456",
			"transfers": [{
				"from": "Mint",
				"to_addresses": ["bb1abc", "bb1def"],
				"balances": [{
					"amount": "50",
					"ownership_times": [{"start": "1", "end": "18446744073709551615"}],
					"token_ids": [{"start": "1", "end": "5"}]
				}],
				"merkle_proofs": [],
				"eth_signature_proofs": [],
				"prioritized_approvals": []
			}],
			"fail_on_error": false,
			"recover_address": "bb1recover"
		}
	}`

	hookData, err := customhookstypes.ParseHookDataFromMemo(memo)
	require.NoError(t, err)
	require.NotNil(t, hookData)
	require.NotNil(t, hookData.TransferTokens)
	require.Equal(t, "456", hookData.TransferTokens.CollectionId)
	require.False(t, hookData.TransferTokens.FailOnError)
	require.Equal(t, "bb1recover", hookData.TransferTokens.RecoverAddress)

	// Verify transfers can be unmarshaled
	proto, err := customhookstypes.UnmarshalTransfersFromJSON(hookData.TransferTokens.Transfers)
	require.NoError(t, err)
	require.Len(t, proto, 1)
	require.Equal(t, "Mint", proto[0].From)
	require.Equal(t, "50", proto[0].Balances[0].Amount.String())
}

// TestMemoSizeLimit tests that the memo size limit still works with transfer_tokens
func TestMemoSizeLimit(t *testing.T) {
	// Create a memo that exceeds 64KB
	bigMemo := `{"transfer_tokens": {"collection_id": "1", "transfers": [{"from": "Mint", "to_addresses": ["bb1abc"], "balances": [{"amount": "1", "ownership_times": [{"start": "1", "end": "1"}], "token_ids": [{"start": "1", "end": "1"}]}], "merkle_proofs": [], "eth_signature_proofs": [], "prioritized_approvals": [], "memo": "`
	// Pad to exceed 64KB
	for len(bigMemo) < 65*1024 {
		bigMemo += "aaaaaaaaaa"
	}
	bigMemo += `"}], "fail_on_error": true}}`

	_, err := customhookstypes.ParseHookDataFromMemo(bigMemo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "memo size exceeds")
}
