package keeper_test

import (
	"context"
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
	alice                  = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob                    = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	charlie                = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	validUri               = "https://example.com/badge.json"
	invalidUri             = "invaliduri"
	firstAccountNumCreated = uint64(7) //Just how it is. I believe the first 6 are validator node accounts
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
}

func TestBadgesKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type BadgesToCreate struct {
	Badge   types.MsgNewBadge
	Amount  uint64
	Creator string
}

func CreateBadges(suite *TestSuite, ctx context.Context, badges []BadgesToCreate) error {
	for _, badge := range badges {
		for i := 0; i < int(badge.Amount); i++ {
			msg := types.NewMsgNewBadge(badge.Creator, badge.Badge.Uri, badge.Badge.Permissions, badge.Badge.SubassetUris)
			_, err := suite.msgServer.NewBadge(ctx, msg)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func CreateSubBadges(suite *TestSuite, ctx context.Context, creator string, badgeId uint64, supplys []uint64, amounts []uint64) error {
	msg := types.NewMsgNewSubBadge(creator, badgeId, supplys, amounts)
	_, err := suite.msgServer.NewSubBadge(ctx, msg)
	return err
}

func RequestTransferBadge(suite *TestSuite, ctx context.Context, creator string, from uint64, amount uint64, badgeId uint64, subbadgeId uint64) error {
	msg := types.NewMsgRequestTransferBadge(creator, from, amount, badgeId, subbadgeId)
	_, err := suite.msgServer.RequestTransferBadge(ctx, msg)
	return err
}

func RevokeBadge(suite *TestSuite, ctx context.Context, creator string, address uint64, amount uint64, badgeId uint64, subbadgeId uint64) error {
	msg := types.NewMsgRevokeBadge(creator, address, amount, badgeId, subbadgeId)
	_, err := suite.msgServer.RevokeBadge(ctx, msg)
	return err
}

func TransferBadge(suite *TestSuite, ctx context.Context, creator string, from uint64, to uint64, amount uint64, badgeId uint64, subbadgeId uint64) error {
	msg := types.NewMsgTransferBadge(creator, from, to, amount, badgeId, subbadgeId)
	_, err := suite.msgServer.TransferBadge(ctx, msg)
	return err
}

func SetApproval(suite *TestSuite, ctx context.Context, creator string, amount uint64, address uint64, badgeId uint64, subbadgeId uint64) error {
	msg := types.NewMsgSetApproval(creator, amount, address, badgeId, subbadgeId)
	_, err := suite.msgServer.SetApproval(ctx, msg)
	return err
}

func HandlePendingTransfer(suite *TestSuite, ctx context.Context, creator string, accept bool, badgeId uint64, subbadgeId uint64, thisNonce uint64) error {
	msg := types.NewMsgHandlePendingTransfer(creator, accept, badgeId, subbadgeId, thisNonce)
	_, err := suite.msgServer.HandlePendingTransfer(ctx, msg)
	return err
}

func FreezeAddress(suite *TestSuite, ctx context.Context, creator string, address uint64, badgeId uint64, subbadgeId uint64, add bool) error {
	msg := types.NewMsgFreezeAddress(creator, address, badgeId, subbadgeId, add)
	_, err := suite.msgServer.FreezeAddress(ctx, msg)
	return err
}

/* Below, we should define all query handlers and use them within the other integration tests. */

func GetBadge(suite *TestSuite, ctx context.Context, id uint64) types.BitBadge {
	res, err := suite.app.BadgesKeeper.GetBadge(ctx, &types.QueryGetBadgeRequest{Id: uint64(id)})
	if err != nil {
		suite.Fail(err.Error())
	}

	return *res.Badge
}

func GetBadgeBalance(suite *TestSuite, ctx context.Context, badgeId uint64, subbadgeId uint64, address uint64) types.BadgeBalanceInfo {
	res, err := suite.app.BadgesKeeper.GetBalance(ctx, &types.QueryGetBalanceRequest{
		BadgeId:    uint64(badgeId),
		SubbadgeId: uint64(subbadgeId),
		Address:    uint64(address),
	})

	if err != nil {
		suite.Fail(err.Error())
	}

	return *res.BalanceInfo
}
