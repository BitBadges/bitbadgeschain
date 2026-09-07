package tokenization

import (
	"strings"
	"testing"

	types "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestMessageGasRejectsUnpricedAndOversizedMessages(t *testing.T) {
	for _, msg := range []sdk.Msg{
		&banktypes.MsgSend{},
		&types.MsgSetValidTokenIds{ValidTokenIds: make([]*types.UintRange, MaxTokenIdRanges+1)},
		&types.MsgSetStandards{Standards: make([]string, MaxApprovalRanges+1)},
		&types.MsgSetStandards{Standards: []string{strings.Repeat("x", MaxMetadataLength+1)}},
		&types.MsgPurgeApprovals{ApprovalsToPurge: make([]*types.ApprovalIdentifierDetails, MaxApprovalRanges+1)},
		&types.MsgCreateDynamicStore{Uri: strings.Repeat("x", MaxMetadataLength+1)},
		&types.MsgUpdateDynamicStore{CustomData: strings.Repeat("x", MaxMetadataLength+1)},
	} {
		t.Run(sdk.MsgTypeURL(msg), func(t *testing.T) { _, err := messageGas(msg); require.Error(t, err) })
	}
}

func TestAdditionalMessageElementsArePriced(t *testing.T) {
	for _, msg := range []sdk.Msg{
		&types.MsgSetValidTokenIds{ValidTokenIds: []*types.UintRange{{}}},
		&types.MsgSetStandards{Standards: []string{"ERC20"}},
		&types.MsgPurgeApprovals{ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{}}},
	} {
		gas, err := messageGas(msg)
		require.NoError(t, err)
		require.Positive(t, gas)
	}
}
