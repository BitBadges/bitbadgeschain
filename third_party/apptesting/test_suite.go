package apptesting

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmod "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"

	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"

	storemetrics "cosmossdk.io/store/metrics"

	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/cosmos/cosmos-sdk/types/module"

	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

type KeeperTestHelper struct {
	suite.Suite

	// defaults to false,
	// set to true if any method that potentially alters baseapp/abci is used.
	// this controls whether or not we can reuse the app instance, or have to set a new one.
	hasUsedAbci bool
	// defaults to false, set to true if we want to use a new app instance with caching enabled.
	// then on new setup test call, we just drop the current cache.
	// this is not always enabled, because some tests may take a painful performance hit due to CacheKv.
	withCaching bool

	App         *app.App
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

// Defines IDs for all supported
// Osmosis pools. Additionally, encapsulates
// an internal gauge ID for each pool.
// This struct is initialized and returned by
// PrepareAllSupportedPools().
type SupportedPoolAndGaugeInfo struct {
	ConcentratedPoolID uint64
	BalancerPoolID     uint64
	StableSwapPoolID   uint64
	CosmWasmPoolID     uint64
	AlloyedPoolID      uint64

	ConcentratedGaugeID uint64
	BalancerGaugeID     uint64
	StableSwapGaugeID   uint64
}

var (
	SecondaryDenom       = "uion"
	SecondaryAmount      = osmomath.NewInt(100000000)
	baseTestAccts        = []sdk.AccAddress{}
	defaultTestStartTime = time.Now().UTC()
	testDescription      = stakingtypes.NewDescription("test_moniker", "test_identity", "test_website", "test_security_contact", "test_details")
)

var (
	WBTC                       = "ibc/D1542AA8762DB13087D8364F3EA6509FD6F009A34F00426AF9E4F9FA85CBBF1F"
	DefaultTickSpacing         = uint64(100)
	DefaultLowerPrice          = osmomath.NewDec(4545)
	DefaultLowerTick           = int64(30545000)
	DefaultUpperPrice          = osmomath.NewDec(5500)
	DefaultUpperTick           = int64(31500000)
	DefaultCurrPrice           = osmomath.NewDec(5000)
	DefaultCurrTick      int64 = 31000000
	DefaultCurrSqrtPrice       = func() osmomath.BigDec {
		curSqrtPrice, _ := osmomath.MonotonicSqrt(DefaultCurrPrice) // 70.710678118654752440
		return osmomath.BigDecFromDec(curSqrtPrice)
	}()
	PerUnitLiqScalingFactor = osmomath.NewDec(1e15).MulMut(osmomath.NewDec(1e12))

	DefaultSpreadRewardAccumCoins = sdk.NewDecCoins(sdk.NewDecCoinFromDec("foo", osmomath.NewDec(50).MulTruncate(PerUnitLiqScalingFactor)))

	DefaultCoinAmount = osmomath.NewInt(1000000000000000000)

	// Default tokens and amounts
	ETH                 = "eth"
	DefaultAmt0         = osmomath.NewInt(1000000)
	DefaultAmt0Expected = osmomath.NewInt(998976)
	DefaultCoin0        = sdk.NewCoin(ETH, DefaultAmt0)
	USDC                = "usdc"
	DefaultAmt1         = osmomath.NewInt(5000000000)
	DefaultAmt1Expected = osmomath.NewInt(5000000000)
	DefaultCoin1        = sdk.NewCoin(USDC, DefaultAmt1)
	DefaultCoins        = sdk.NewCoins(DefaultCoin0, DefaultCoin1)
)

func init() {
	baseTestAccts = CreateRandomAccounts(3)
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts)
// preserves the caching enabled/disabled state.
func (s *KeeperTestHelper) Setup() {
	s.T().Log("Setting up KeeperTestHelper")

	s.App = app.Setup(false)
	s.setupGeneral()
}

// PrepareAllSupportedPoolsCustomProject creates all supported pools and returns their IDs.
// Additionally, attaches an internal gauge ID for each pool.
// Allows the flexibility of being used from outside the Osmosis repository by providing custom project name and transmuter bytecode path.
func (s *KeeperTestHelper) PrepareAllSupportedPoolsCustomProject(projectName, transmuterPath string) SupportedPoolAndGaugeInfo {

	var (
		// Prepare pools and their IDs
		balancerPoolID   = s.PrepareBalancerPool()
		stableswapPoolID = s.PrepareBasicStableswapPool()
	)

	return SupportedPoolAndGaugeInfo{
		BalancerPoolID:   balancerPoolID,
		StableSwapPoolID: stableswapPoolID,
	}
}

// resets the test environment
// requires that all commits go through helpers in s.
// On first reset, will instantiate a new app, with caching enabled.
// NOTE: If you are using ABCI methods, usage of Reset vs Setup has not been well tested.
// It is believed to work, but if you get an odd error, try changing the call to this for setup to sanity check.
// what's supposed to happen is a new setup call, and reset just does that in such a case.
func (s *KeeperTestHelper) Reset() {
	if s.hasUsedAbci || !s.withCaching {
		s.withCaching = true
		s.Setup()
	} else {
		s.App.PoolManagerKeeper.ResetCaches()
		s.setupGeneral()
	}
}

func (s *KeeperTestHelper) setupGeneral() {
	s.setupGeneralCustomChainId("bitbadges-1")
}

func (s *KeeperTestHelper) setupGeneralCustomChainId(chainId string) {
	s.Ctx = s.App.BaseApp.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: chainId, Time: defaultTestStartTime})
	if s.withCaching {
		s.Ctx, _ = s.Ctx.CacheContext()
	}
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = []sdk.AccAddress{}
	s.TestAccs = append(s.TestAccs, baseTestAccts...)
	s.hasUsedAbci = false
}

func (s *KeeperTestHelper) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	s.App = app.Setup(true)
	s.Ctx = s.App.BaseApp.NewContextLegacy(true, cmtproto.Header{})
	// TODO: not sure
	s.hasUsedAbci = true
}

// RunTestCaseWithoutStateUpdates runs the testcase as a callback with the given name.
// Does not persist any state changes. This is useful when test suite uses common state setup
// but desures each test case to be run in isolation.
func (s *KeeperTestHelper) RunTestCaseWithoutStateUpdates(name string, cb func(t *testing.T)) {
	originalCtx := s.Ctx
	s.Ctx, _ = s.Ctx.CacheContext()

	s.T().Run(name, cb)

	s.Ctx = originalCtx
}

// CreateTestContext creates a test context.
func (s *KeeperTestHelper) CreateTestContext() sdk.Context {
	ctx, _ := s.CreateTestContextWithMultiStore()
	return ctx
}

// CreateTestContextWithMultiStore creates a test context and returns it together with multi store.
func (s *KeeperTestHelper) CreateTestContextWithMultiStore() (sdk.Context, storetypes.CommitMultiStore) {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()

	ms := rootmulti.NewStore(db, logger, storemetrics.NewNoOpMetrics())

	return sdk.NewContext(ms, cmtproto.Header{}, false, logger), ms
}

// CreateTestContext creates a test context.
func (s *KeeperTestHelper) Commit() {
	_, err := s.App.FinalizeBlock(&abci.RequestFinalizeBlock{Height: s.Ctx.BlockHeight(), Time: s.Ctx.BlockTime()})
	if err != nil {
		panic(err)
	}
	_, err = s.App.Commit()
	if err != nil {
		panic(err)
	}

	newBlockTime := s.Ctx.BlockTime().Add(time.Second)

	header := s.Ctx.BlockHeader()
	header.Time = newBlockTime
	header.Height++

	s.Ctx = s.App.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
		Height: header.Height,
		Time:   header.Time,
	})

	s.hasUsedAbci = true
}

// FundAcc funds target address with specified amount.
func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := testutil.FundAccount(s.Ctx, s.App.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}

// FundModuleAcc funds target modules with specified amount.
func (s *KeeperTestHelper) FundModuleAcc(moduleName string, amounts sdk.Coins) {
	err := testutil.FundModuleAccount(s.Ctx, s.App.BankKeeper, moduleName, amounts)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) MintCoins(coins sdk.Coins) {
	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins)
	s.Require().NoError(err)
}

// BeginNewBlockWithProposer begins a new block with a proposer.
func (s *KeeperTestHelper) BeginNewBlockWithProposer(executeNextEpoch bool, proposer sdk.ValAddress) {
	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
	s.Assert().NoError(err)

	valConsAddr, err := validator.GetConsAddr()
	s.Require().NoError(err)

	valAddr := valConsAddr

	newBlockTime := s.Ctx.BlockTime().Add(5 * time.Second)
	header := cmtproto.Header{Height: s.Ctx.BlockHeight() + 1, Time: newBlockTime}
	s.Ctx = s.Ctx.WithBlockTime(newBlockTime).WithBlockHeight(s.Ctx.BlockHeight() + 1)
	voteInfos := []abci.VoteInfo{{
		Validator:   abci.Validator{Address: valAddr, Power: 1000},
		BlockIdFlag: cmtproto.BlockIDFlagCommit,
	}}
	s.Ctx = s.Ctx.WithVoteInfos(voteInfos)

	fmt.Println("beginning block ", s.Ctx.BlockHeight())

	_, err = s.App.BeginBlocker(s.Ctx)
	s.Require().NoError(err)

	s.Ctx = s.App.NewContextLegacy(false, header)
	s.hasUsedAbci = true
}

// EndBlock ends the block, and runs commit
func (s *KeeperTestHelper) EndBlock() {
	_, err := s.App.EndBlocker(s.Ctx)
	s.Require().NoError(err)
	s.hasUsedAbci = true
}

func (s *KeeperTestHelper) RunMsg(msg sdk.Msg) (*sdk.Result, error) {
	// cursed that we have to copy this internal logic from SDK
	router := s.App.BaseApp.MsgServiceRouter()
	if handler := router.Handler(msg); handler != nil {
		// ADR 031 request type routing
		return handler(s.Ctx, msg)
	}
	s.FailNow("msg %v could not be ran", msg)
	s.hasUsedAbci = true
	return nil, fmt.Errorf("msg %v could not be ran", msg)
}

// AllocateRewardsToValidator allocates reward tokens to a distribution module then allocates rewards to the validator address.
func (s *KeeperTestHelper) AllocateRewardsToValidator(valAddr sdk.ValAddress, rewardAmt osmomath.Int) {
	validator, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().NoError(err)

	// allocate reward tokens to distribution module
	coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, rewardAmt)}
	err = testutil.FundModuleAccount(s.Ctx, s.App.BankKeeper, distrtypes.ModuleName, coins)
	s.Require().NoError(err)

	// allocate rewards to validator
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
	decTokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: osmomath.NewDec(20000)}}
	err = s.App.DistrKeeper.AllocateTokensToValidator(s.Ctx, validator, decTokens)
	s.Require().NoError(err)
}

// SetupGammPoolsWithBondDenomMultiplier uses given multipliers to set initial pool supply of bond denom.
func (s *KeeperTestHelper) SetupGammPoolsWithBondDenomMultiplier(multipliers []osmomath.Dec) []gammtypes.CFMMPoolI {
	bondDenom, err := s.App.StakingKeeper.BondDenom(s.Ctx)
	s.Require().NoError(err)
	// TODO: use sdk crypto instead of tendermint to generate address
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	pools := []gammtypes.CFMMPoolI{}
	for index, multiplier := range multipliers {
		token := fmt.Sprintf("token%d", index)
		uosmoAmount := gammtypes.InitPoolSharesSupply.ToLegacyDec().Mul(multiplier).RoundInt()

		s.FundAcc(acc1, sdk.NewCoins(
			sdk.NewCoin(bondDenom, uosmoAmount.Mul(osmomath.NewInt(10))),
			sdk.NewInt64Coin(token, 100000),
		))

		var (

			// pool assets
			defaultFooAsset = balancer.PoolAsset{
				Weight: osmomath.NewInt(100),
				Token:  sdk.NewCoin(bondDenom, uosmoAmount),
			}
			defaultBarAsset = balancer.PoolAsset{
				Weight: osmomath.NewInt(100),
				Token:  sdk.NewCoin(token, osmomath.NewInt(10000)),
			}

			poolAssets = []balancer.PoolAsset{defaultFooAsset, defaultBarAsset}
		)

		poolParams := balancer.PoolParams{
			SwapFee: osmomath.NewDecWithPrec(1, 2),
			ExitFee: osmomath.Dec(osmomath.ZeroInt()),
		}
		msg := balancer.NewMsgCreateBalancerPool(acc1, poolParams, poolAssets)

		poolId, err := s.App.PoolManagerKeeper.CreatePool(s.Ctx, msg)
		s.Require().NoError(err)

		pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
		s.Require().NoError(err)

		pools = append(pools, pool)
	}

	return pools
}

// SwapAndSetSpotPrice runs a swap to set Spot price of a pool using arbitrary values
// returns spot price after the arbitrary swap.
func (s *KeeperTestHelper) SwapAndSetSpotPrice(poolId uint64, fromAsset sdk.Coin, toAsset sdk.Coin) osmomath.BigDec {
	// create a dummy account
	acc1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// fund dummy account with tokens to swap
	coins := sdk.Coins{sdk.NewInt64Coin(fromAsset.Denom, 100000000000000)}
	s.FundAcc(acc1, coins)

	route := []poolmanagertypes.SwapAmountOutRoute{
		{
			PoolId:       poolId,
			TokenInDenom: fromAsset.Denom,
		},
	}
	_, err := s.App.PoolManagerKeeper.RouteExactAmountOut(
		s.Ctx,
		acc1,
		route,
		fromAsset.Amount,
		sdk.NewCoin(toAsset.Denom,
			toAsset.Amount.Quo(osmomath.NewInt(4))))
	s.Require().NoError(err)

	spotPrice, err := s.App.GammKeeper.CalculateSpotPrice(s.Ctx, poolId, fromAsset.Denom, toAsset.Denom)
	s.Require().NoError(err)

	return spotPrice
}

// BuildTx builds a transaction.
func (s *KeeperTestHelper) BuildTx(
	txBuilder client.TxBuilder,
	msgs []sdk.Msg,
	sigV2 signing.SignatureV2,
	memo string, txFee sdk.Coins,
	gasLimit uint64,
) authsigning.Tx {
	err := txBuilder.SetMsgs(msgs[0])
	s.Require().NoError(err)

	err = txBuilder.SetSignatures(sigV2)
	s.Require().NoError(err)

	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(txFee)
	txBuilder.SetGasLimit(gasLimit)

	return txBuilder.GetTx()
}

// StateNotAltered validates that app state is not altered. Fails if it is.
func (s *KeeperTestHelper) StateNotAltered() {
	oldState, err := s.App.ExportAppStateAndValidators(false, []string{}, []string{})
	s.Require().NoError(err)
	s.App.CommitMultiStore().Commit()
	newState, err := s.App.ExportAppStateAndValidators(false, []string{}, []string{})
	s.Require().NoError(err)
	s.Require().Equal(oldState, newState)
	s.hasUsedAbci = true
}

func (s *KeeperTestHelper) SkipIfWSL() {
	SkipIfWSL(s.T())
}

// SkipIfWSL skips tests if running on WSL
// This is a workaround to enable quickly running full unit test suite locally
// on WSL without failures. The failures are stemming from trying to upload
// wasm code. An OS permissioning issue.
func SkipIfWSL(t *testing.T) {
	t.Helper()
	skip := os.Getenv("SKIP_WASM_WSL_TESTS")
	fmt.Println("SKIP_WASM_WSL_TESTS", skip)
	if skip == "true" {
		t.Skip("Skipping Wasm tests")
	}
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func TestMessageAuthzSerialization(t *testing.T, msg sdk.Msg, module module.AppModuleBasic) {
	someDate := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
	const (
		mockGranter string = "cosmos1abc"
		mockGrantee string = "cosmos1xyz"
	)

	var (
		mockMsgGrant  authz.MsgGrant
		mockMsgRevoke authz.MsgRevoke
		mockMsgExec   authz.MsgExec
	)

	// mutates mockMsg
	testSerDeser := func(msg proto.Message, mockMsg proto.Message) {
		encCfg := moduletestutil.MakeTestEncodingConfig(authzmod.AppModuleBasic{}, module)
		msgGrantBytes := json.RawMessage(sdk.MustSortJSON(encCfg.Codec.MustMarshalJSON(msg)))
		err := encCfg.Codec.UnmarshalJSON(msgGrantBytes, mockMsg)
		require.NoError(t, err)
	}

	// Authz: Grant Msg
	typeURL := sdk.MsgTypeURL(msg)
	expiryTime := someDate.Add(time.Hour)
	grant, err := authz.NewGrant(someDate, authz.NewGenericAuthorization(typeURL), &expiryTime)
	require.NoError(t, err)

	msgGrant := authz.MsgGrant{Granter: mockGranter, Grantee: mockGrantee, Grant: grant}
	testSerDeser(&msgGrant, &mockMsgGrant)

	// Authz: Revoke Msg
	msgRevoke := authz.MsgRevoke{Granter: mockGranter, Grantee: mockGrantee, MsgTypeUrl: typeURL}
	testSerDeser(&msgRevoke, &mockMsgRevoke)

	// Authz: Exec Msg
	msgAny, err := cdctypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	msgExec := authz.MsgExec{Grantee: mockGrantee, Msgs: []*cdctypes.Any{msgAny}}
	testSerDeser(&msgExec, &mockMsgExec)
	require.Equal(t, msgExec.Msgs[0].Value, mockMsgExec.Msgs[0].Value)
}

func GenerateTestAddrs() (string, string) {
	pk1 := ed25519.GenPrivKey().PubKey()
	validAddr := sdk.AccAddress(pk1.Address()).String()
	invalidAddr := sdk.AccAddress("invalid").String()
	return validAddr, invalidAddr
}

// sets up the volume for the pools in the group
// mutates poolIDToVolumeMap
func (s *KeeperTestHelper) SetupVolumeForPools(poolIDs []uint64, volumesForEachPool []osmomath.Int, poolIDToVolumeMap map[uint64]math.Int) {
	bondDenom, err := s.App.StakingKeeper.BondDenom(s.Ctx)
	s.Require().NoError(err)

	s.Require().Equal(len(poolIDs), len(volumesForEachPool))
	for i := 0; i < len(poolIDs); i++ {
		currentPoolID := poolIDs[i]

		currentVolume := volumesForEachPool[i]

		fmt.Printf("currentVolume %d %s\n", i, currentVolume)

		// Retrieve the existing volume to add to it.
		existingVolume := s.App.PoolManagerKeeper.GetOsmoVolumeForPool(s.Ctx, currentPoolID)

		s.App.PoolManagerKeeper.SetVolume(s.Ctx, currentPoolID, sdk.NewCoins(sdk.NewCoin(bondDenom, existingVolume.Add(currentVolume))))

		if existingVolume, ok := poolIDToVolumeMap[currentPoolID]; ok {
			poolIDToVolumeMap[currentPoolID] = existingVolume.Add(currentVolume)
		} else {
			poolIDToVolumeMap[currentPoolID] = currentVolume
		}
	}
}

// initializes or increases the volumes for the given pools
func (s *KeeperTestHelper) IncreaseVolumeForPools(poolIDs []uint64, volumesForEachPool []osmomath.Int) {
	s.SetupVolumeForPools(poolIDs, volumesForEachPool, map[uint64]osmomath.Int{})
}

// RequireDecCoinsSlice compares two slices of DecCoins
func (s *KeeperTestHelper) RequireDecCoinsSlice(expected, actual []sdk.DecCoins) {
	s.Require().Equal(len(expected), len(actual))
	for i := range actual {
		s.Require().Equal(expected[i].String(), actual[i].String())
	}
}
