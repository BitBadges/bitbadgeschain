package app

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

// TestModuleAccountsBlockedFromBankSends checks that every module account
// except gov (which receives proposal deposits) refuses bank sends, both in
// the live bank keeper and in the BlockedAddresses helper used by the
// simulation tests.
func TestModuleAccountsBlockedFromBankSends(t *testing.T) {
	app := Setup(false)

	blocked := BlockedAddresses()
	// All active module accounts except gov, plus the retired precisebank reserve.
	require.Len(t, blocked, len(moduleAccPerms), "all module accounts but gov, plus retired precisebank, must be blocked")
	require.True(t, blocked["precisebank"], "retired precisebank reserve must remain blocked")
	require.True(t, app.BankKeeper.BlockedAddr(authtypes.NewModuleAddress("precisebank")),
		"retired precisebank reserve must be blocked in BankKeeper")

	for _, perm := range moduleAccPerms {
		addr := authtypes.NewModuleAddress(perm.Account)
		if perm.Account == govtypes.ModuleName {
			require.False(t, app.BankKeeper.BlockedAddr(addr), "gov must accept deposits")
			require.False(t, blocked[perm.Account])
			continue
		}
		require.True(t, app.BankKeeper.BlockedAddr(addr), "module account %q must be blocked from bank sends", perm.Account)
		require.True(t, blocked[perm.Account], "BlockedAddresses must list %q", perm.Account)
	}
}
