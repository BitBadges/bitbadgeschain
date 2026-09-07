package app

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"
	"time"

	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/stretchr/testify/require"
)

func TestV35MigrationsAgainstCapturedMainnetStores(t *testing.T) {
	raw, err := os.ReadFile("testdata/v35-mainnet-11979578.json")
	require.NoError(t, err)
	var fixture struct {
		Height int64 `json:"height"`
		Stores map[string][]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"stores"`
	}
	require.NoError(t, json.Unmarshal(raw, &fixture))
	require.Equal(t, int64(11979578), fixture.Height)
	require.Greater(t, len(fixture.Stores["tokenization"]), 900, "a zero-row rehearsal cannot exercise migrations")
	app := Setup(false)
	ctx := app.NewContext(false).WithBlockHeight(fixture.Height).WithBlockTime(time.Unix(1788739200, 0))
	for module, rows := range fixture.Stores {
		store := ctx.KVStore(app.GetKey(module))
		for _, row := range rows {
			key, err := base64.StdEncoding.DecodeString(row.Key)
			require.NoError(t, err)
			value, err := base64.StdEncoding.DecodeString(row.Value)
			require.NoError(t, err)
			store.Set(key, value)
		}
	}
	keepers := v35.Keepers{Account: app.AccountKeeper, ConsensusParams: app.ConsensusParamsKeeper, FeeMarket: app.FeeMarketKeeper, IBCRateLimit: app.IBCRateLimitKeeper, Tokenization: app.TokenizationKeeper, ManagerSplitter: app.ManagerSplitterKeeper}
	started := time.Now()
	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, keepers))
	t.Logf("captured mainnet module migration completed in %s", time.Since(started))
	snapshot := func() map[string]map[string]string {
		out := map[string]map[string]string{}
		for module := range fixture.Stores {
			out[module] = map[string]string{}
			it := storetypes.KVStorePrefixIterator(ctx.KVStore(app.GetKey(module)), nil)
			for ; it.Valid(); it.Next() {
				out[module][string(it.Key())] = string(it.Value())
			}
			require.NoError(t, it.Close())
		}
		return out
	}
	once := snapshot()
	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, keepers))
	require.Equal(t, once, snapshot(), "populated migrations must be idempotent")
}
