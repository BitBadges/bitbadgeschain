package types_test

import (
	"strings"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddressListCannotUseUppercaseAccountAsId(t *testing.T) {
	address := sdk.AccAddress([]byte("address-list-account")).String()
	require.Error(t, types.ValidateAddressList(&types.AddressList{ListId: strings.ToUpper(address)}))
}
