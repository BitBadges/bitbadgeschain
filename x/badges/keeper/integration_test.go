package keeper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"

	bitbadgesapp "github.com/trevormil/bitbadgeschain/app"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	alice             = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob               = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	charlie           = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	bobAccountNum     = uint64(7) //7 is just how it is. I believe the first 6 are validator node accounts
	aliceAccountNum   = uint64(8)
	charlieAccountNum = uint64(9)
)

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type TestSuite struct {
	suite.Suite

	app         *bitbadgesapp.App
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

//Bunch of weird config stuff to setup the app. Inherited most from Cosmos SDK tutorials and existing Cosmos SDK modules.
func (suite *TestSuite) SetupTest() {
	simapp.FlagEnabledValue = true
	simapp.FlagCommitValue = true

	_, db, _, logger, _, err := simapp.SetupSimulation("goleveldb-app-sim", "Simulation")
	if err != nil {
		panic("Error constructing simapp")
	}

	encoding := cosmoscmd.MakeEncodingConfig(bitbadgesapp.ModuleBasics)

	app := bitbadgesapp.NewApp(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		bitbadgesapp.DefaultNodeHome,
		0,
		encoding,
		simapp.EmptyAppOptions{},
	)

	genesisState := bitbadgesapp.NewDefaultGenesisState(app.AppCodec())
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(abci.RequestInitChain{
		Validators:      []abci.ValidatorUpdate{},
		AppStateBytes:   stateBytes,
		ConsensusParams: DefaultConsensusParams,
	})

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.msgServer = keeper.NewMsgServerImpl(app.BadgesKeeper)
	suite.queryClient = queryClient

	bob_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(bob))
	alice_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(alice))
	charlie_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(charlie))

	suite.app.AccountKeeper.SetAccount(suite.ctx, bob_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, alice_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, charlie_acc)

	for i := uint64(0); i < 1000; i++ {
		suite.app.AccountKeeper.SetAccount(suite.ctx, suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.AccAddress([]byte{byte(i)})))
	}
}

func TestBadgesKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
