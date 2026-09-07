package types_test

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// H-6 regression: only the canonical (lowercase) bech32 spelling of an address is accepted,
// so one key cannot appear under two spellings in balance keys, trackers and address lists.
func TestValidateAddressRequiresCanonicalBech32(t *testing.T) {
	canonical := sample.AccAddress()
	upper := strings.ToUpper(canonical)
	require.NotEqual(t, canonical, upper)

	require.NoError(t, types.ValidateAddress(canonical, false))
	require.ErrorIs(t, types.ValidateAddress(upper, false), types.ErrInvalidAddress)
	require.NoError(t, types.ValidateAddress("Mint", true))
}

func TestMsgTransferTokensRejectsNonCanonicalAddresses(t *testing.T) {
	canonical := sample.AccAddress()
	recipient := sample.AccAddress()
	upper := strings.ToUpper(canonical)

	base := func(creator, to string) *types.MsgTransferTokens {
		return &types.MsgTransferTokens{
			Creator:      creator,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{
				From:        "Mint",
				ToAddresses: []string{to},
				Balances:    []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}}},
			}},
		}
	}

	require.NoError(t, base(canonical, recipient).ValidateBasic())
	require.ErrorIs(t, base(upper, recipient).ValidateBasic(), types.ErrInvalidAddress, "creator")
	require.ErrorIs(t, base(canonical, strings.ToUpper(recipient)).ValidateBasic(), types.ErrInvalidAddress, "recipient")
}

func TestMsgCreateAddressListsRejectsNonCanonicalEntries(t *testing.T) {
	creator := sample.AccAddress()
	member := sample.AccAddress()

	msg := &types.MsgCreateAddressLists{
		Creator:      creator,
		AddressLists: []*types.AddressListInput{{ListId: "list1", Addresses: []string{strings.ToUpper(member)}, Whitelist: true}},
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidAddress)

	msg.AddressLists[0].Addresses = []string{member}
	require.NoError(t, msg.ValidateBasic())
}
