package app

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// shippedFreshChainArtifacts are the files a fresh chain is started from. They
// are not covered by the upgrade handler — nothing migrates them — so they drift
// silently the moment the binary's denom or decimals change, and the failure
// shows up as a chain that will not start rather than as a test failure.
//
// genesis-711316.json is deliberately NOT in this list. It is an export of
// mainnet state at height 711316, i.e. a snapshot of the 9-decimal chain, and
// the v35 upgrade handler is what converts that state. Rewriting it by hand
// would produce a genesis labelled 18-decimal carrying 9-decimal amounts.
// Whether the fork-restart path should ship a pre-converted genesis instead of
// upgrading in place is an open decision, not something this test should
// prejudge.
var shippedFreshChainArtifacts = []string{
	"genesis.json",
	"testnet.genesis.json",
	"config.yml",
	"config.testnet.yml",
	"start-chain.sh",
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	return filepath.Dir(wd) // the test runs in app/
}

// The retired denom in a fresh-chain artifact means that chain starts with
// balances in a denom the binary does not recognise as its base unit.
func TestShippedGenesisArtifactsDoNotNameTheRetiredDenom(t *testing.T) {
	root := repoRoot(t)
	// Word boundaries so an ibc/ hash that happens to contain the string does
	// not produce a false positive.
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(appparams.LegacyBaseCoinUnit) + `\b`)

	for _, name := range shippedFreshChainArtifacts {
		body, err := os.ReadFile(filepath.Join(root, name))
		require.NoError(t, err, "shipped artifact %s must exist", name)

		require.Empty(t, re.FindAllString(string(body), -1),
			"%s still names %s; a chain started from it will not agree with the binary's base denom",
			name, appparams.LegacyBaseCoinUnit)
		require.True(t, strings.Contains(string(body), appparams.BaseCoinUnit),
			"%s names neither the retired denom nor the current one", name)
	}
}

// x/vm derives the chain's decimals from the exponent of the display unit in
// bank's denom metadata. A genesis carrying 9 there boots an 18-decimal binary
// into believing it is a 9-decimal chain.
func TestShippedGenesisDenomMetadataMatchesAppParams(t *testing.T) {
	root := repoRoot(t)

	for _, name := range []string{"genesis.json", "testnet.genesis.json"} {
		body, err := os.ReadFile(filepath.Join(root, name))
		require.NoError(t, err)

		var doc struct {
			AppState struct {
				Bank struct {
					DenomMetadata []struct {
						Base       string `json:"base"`
						Display    string `json:"display"`
						DenomUnits []struct {
							Denom    string `json:"denom"`
							Exponent uint32 `json:"exponent"`
						} `json:"denom_units"`
					} `json:"denom_metadata"`
				} `json:"bank"`
			} `json:"app_state"`
		}
		require.NoError(t, json.Unmarshal(body, &doc), "parsing %s", name)

		for _, md := range doc.AppState.Bank.DenomMetadata {
			if md.Base != appparams.BaseCoinUnit && md.Display != appparams.DisplayCoinUnit {
				continue
			}
			require.Equal(t, appparams.BaseCoinUnit, md.Base, "%s: base denom", name)
			for _, unit := range md.DenomUnits {
				if unit.Denom == appparams.DisplayCoinUnit {
					require.Equal(t, uint32(appparams.BaseCoinDecimals), unit.Exponent,
						"%s: the display unit's exponent is what x/vm reads as the chain's decimals", name)
				}
			}
		}
	}
}

// PowerReduction moved from the SDK default 10^6 to 10^15 with the decimals.
// Consensus power is bonded tokens / PowerReduction, so a gentx that bonds less
// than PowerReduction produces a validator with zero power and InitGenesis
// fails with "validator set is empty after InitGenesis" — before any of this
// package's other tests get a chance to run against it.
func TestStartChainScriptBondsEnoughForNonZeroVotingPower(t *testing.T) {
	root := repoRoot(t)
	body, err := os.ReadFile(filepath.Join(root, "start-chain.sh"))
	require.NoError(t, err)

	bonds := regexp.MustCompile(`gentx[^\n]*?"(\d+)ustake"`).FindAllStringSubmatch(string(body), -1)
	require.NotEmpty(t, bonds, "expected start-chain.sh to bond a gentx amount in ustake")

	for _, match := range bonds {
		amount, ok := new(big.Int).SetString(match[1], 10)
		require.True(t, ok, "unparseable bond amount %q", match[1])
		require.True(t, amount.Cmp(appparams.PowerReduction.BigInt()) >= 0,
			"gentx bonds %s ustake but PowerReduction is %s: the validator would have zero voting power",
			match[1], appparams.PowerReduction)
	}
}
