package keeper_test

import (
	"math"
	"testing"
	"time"

	bankv1beta1 "cosmossdk.io/api/cosmos/bank/v1beta1"
	basev1beta1 "cosmossdk.io/api/cosmos/base/v1beta1"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	"bitbadgeschain/x/offers/keeper"
	"bitbadgeschain/x/offers/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bitbadgesapp "bitbadgeschain/app"

	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"

	// Add these imports if not already present

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

const (
	//Note these are alphanumerically sorted (needed for approvals test)
	alice   = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	bob     = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	charlie = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
)

type TestSuite struct {
	suite.Suite

	app         *bitbadgesapp.App
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

// Bunch of weird config stuff to setup the app. Inherited most from Cosmos SDK tutorials and existing Cosmos SDK modules.
func (suite *TestSuite) SetupTest() {
	app := bitbadgesapp.Setup(
		false,
	)

	ctx := app.BaseApp.NewContext(false)

	// app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.AppCodec().InterfaceRegistry())
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.msgServer = keeper.NewMsgServerImpl(app.OffersKeeper)
	suite.queryClient = queryClient

	bob_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(bob))
	alice_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(alice))
	charlie_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(charlie))

	suite.app.AccountKeeper.SetAccount(suite.ctx, bob_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, alice_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, charlie_acc)

	//initialize bob with 1000 coins

	suite.ctx = suite.ctx.WithBlockTime(time.Now())

	// for i := uint64(0); i < 1000; i++ {
	// 	suite.app.AccountKeeper.SetAccount(suite.ctx, suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.AccAddress([]byte{byte(i)})))
	// }

	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(bob), sdk.NewCoins(sdk.NewInt64Coin("ubadge", 1000)))
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(alice), sdk.NewCoins(sdk.NewInt64Coin("ubadge", 1000)))
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(charlie), sdk.NewCoins(sdk.NewInt64Coin("ubadge", 1000)))
}

func GetFullUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}
}

func GetTestMsgs(creator string) []*types.ExecutableMsgWithOptions {
	msg := &bankv1beta1.MsgSend{
		FromAddress: creator,
		ToAddress:   charlie,
		Amount: []*basev1beta1.Coin{
			{
				Denom:  "ubadge",
				Amount: "1",
			},
		},
	}

	protoMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	finalMsgs := []*types.ExecutableMsgWithOptions{
		{Msg: protoMsg, UseContractAddress: false},
	}

	return finalMsgs
}

func (suite *TestSuite) TestCreateProposal() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	msg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{
				Creator:       bob,
				MsgsToExecute: bobMsgs,
				Accepted:      false,
			},
			{
				Creator:       alice,
				MsgsToExecute: aliceMsgs,
				Accepted:      false,
			},
		},
		ValidTimes: GetFullUintRanges(),
	}

	resp, err := msgServer.CreateProposal(sdk.WrapSDKContext(ctx), msg)

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().NotEmpty(resp.Id)

	proposal, found := suite.app.OffersKeeper.GetProposalFromStore(ctx, resp.Id)
	suite.Require().True(found)
	for idx, party := range proposal.Parties {
		suite.Require().Equal(party.Creator, msg.Parties[idx].Creator)
		// suite.Require().Equal(party.MsgsToExecute, msg.Parties[0].MsgsToExecute)
		suite.Require().Equal(party.Accepted, msg.Parties[idx].Accepted)
	}
}

func (suite *TestSuite) TestAcceptProposal() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, err := msgServer.CreateProposal(sdk.WrapSDKContext(ctx), createMsg)
	suite.Require().NoError(err)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	acceptResp, err := msgServer.AcceptProposal(sdk.WrapSDKContext(ctx), acceptMsg)

	suite.Require().NoError(err)
	suite.Require().NotNil(acceptResp)

	proposal, found := suite.app.OffersKeeper.GetProposalFromStore(ctx, createResp.Id)
	suite.Require().True(found)
	suite.Require().True(proposal.Parties[1].Accepted)
	suite.Require().False(proposal.Parties[0].Accepted)
}

func (suite *TestSuite) TestRejectAndDeleteProposal() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, err := msgServer.CreateProposal(sdk.WrapSDKContext(ctx), createMsg)
	suite.Require().NoError(err)

	rejectMsg := &types.MsgRejectAndDeleteProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	rejectResp, err := msgServer.RejectAndDeleteProposal(sdk.WrapSDKContext(ctx), rejectMsg)

	suite.Require().NoError(err)
	suite.Require().NotNil(rejectResp)

	_, found := suite.app.OffersKeeper.GetProposalFromStore(ctx, createResp.Id)
	suite.Require().False(found)
}

func (suite *TestSuite) TestExecuteProposal() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, err := msgServer.CreateProposal(sdk.WrapSDKContext(ctx), createMsg)
	suite.Require().NoError(err)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	AcceptProposal(msgServer, ctx, acceptMsg)

	acceptMsg = &types.MsgAcceptProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	AcceptProposal(msgServer, ctx, acceptMsg)

	executeMsg := &types.MsgExecuteProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	executeResp, err := msgServer.ExecuteProposal(sdk.WrapSDKContext(ctx), executeMsg)

	suite.Require().NoError(err)
	suite.Require().NotNil(executeResp)

	_, found := suite.app.OffersKeeper.GetProposalFromStore(ctx, createResp.Id)
	suite.Require().False(found)

	// Check if the messages were actually executed
	bobBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	charlieBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge")

	suite.Require().Equal(int64(999), bobBalance.Amount.Int64())      // 1000 - 1
	suite.Require().Equal(int64(999), aliceBalance.Amount.Int64())    // 1000 - 1
	suite.Require().Equal(int64(1002), charlieBalance.Amount.Int64()) // 1000 + 1 + 1
}

func CreateProposal(msgServer types.MsgServer, ctx sdk.Context, msg *types.MsgCreateProposal) (*types.MsgCreateProposalResponse, error) {
	//validate basic first
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msgServer.CreateProposal(sdk.WrapSDKContext(ctx), msg)
}

func AcceptProposal(msgServer types.MsgServer, ctx sdk.Context, msg *types.MsgAcceptProposal) (*types.MsgAcceptProposalResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msgServer.AcceptProposal(sdk.WrapSDKContext(ctx), msg)
}

func RejectAndDeleteProposal(msgServer types.MsgServer, ctx sdk.Context, msg *types.MsgRejectAndDeleteProposal) (*types.MsgRejectAndDeleteProposalResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msgServer.RejectAndDeleteProposal(sdk.WrapSDKContext(ctx), msg)
}

func ExecuteProposal(msgServer types.MsgServer, ctx sdk.Context, msg *types.MsgExecuteProposal) (*types.MsgExecuteProposalResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	return msgServer.ExecuteProposal(sdk.WrapSDKContext(ctx), msg)
}

func (suite *TestSuite) TestCreateProposalWithEmptyParties() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	msg := &types.MsgCreateProposal{
		Creator:    bob,
		Parties:    []*types.Parties{},
		ValidTimes: GetFullUintRanges(),
	}

	_, err := CreateProposal(msgServer, ctx, msg)
	suite.Require().Error(err)
}

func (suite *TestSuite) TestCreateProposalWithDuplicateParties() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)

	msg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}

	_, err := CreateProposal(msgServer, ctx, msg)
	suite.Require().Error(err)
}

func (suite *TestSuite) TestAcceptProposalNonExistent() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      sdkmath.NewUint(122),
	}
	_, err := AcceptProposal(msgServer, ctx, acceptMsg)

	suite.Require().Error(err)
}

func (suite *TestSuite) TestAcceptProposalUnauthorized() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, _ := CreateProposal(msgServer, ctx, createMsg)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: charlie,
		Id:      createResp.Id,
	}
	_, err := AcceptProposal(msgServer, ctx, acceptMsg)

	suite.Require().Error(err)
}

func (suite *TestSuite) TestExecuteProposalNotAllAccepted() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, _ := CreateProposal(msgServer, ctx, createMsg)

	executeMsg := &types.MsgExecuteProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	_, err := ExecuteProposal(msgServer, ctx, executeMsg)

	suite.Require().Error(err)
}

func (suite *TestSuite) TestExecuteProposalUnauthorized() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, _ := CreateProposal(msgServer, ctx, createMsg)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	AcceptProposal(msgServer, ctx, acceptMsg)

	acceptMsg = &types.MsgAcceptProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	AcceptProposal(msgServer, ctx, acceptMsg)

	executeMsg := &types.MsgExecuteProposal{
		Creator: charlie,
		Id:      createResp.Id,
	}
	_, err := ExecuteProposal(msgServer, ctx, executeMsg)

	suite.Require().Error(err)
}

func (suite *TestSuite) TestItAutoPopulatesFromAddressWithPartyCreator() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: aliceMsgs, Accepted: false}, // bob is creator but has aliceMsgs
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, err := CreateProposal(msgServer, ctx, createMsg)
	suite.Require().NoError(err)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	_, err = AcceptProposal(msgServer, ctx, acceptMsg)
	suite.Require().NoError(err)

	acceptMsg = &types.MsgAcceptProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	_, err = AcceptProposal(msgServer, ctx, acceptMsg)
	suite.Require().NoError(err)

	executeMsg := &types.MsgExecuteProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	_, err = ExecuteProposal(msgServer, ctx, executeMsg)
	suite.Require().NoError(err)

	// Check if the messages were actually executed
	bobBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	charlieBalance := suite.app.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge")

	suite.Require().Equal(int64(999), bobBalance.Amount.Int64())      // 1000 - 1
	suite.Require().Equal(int64(999), aliceBalance.Amount.Int64())    // 1000 - 1
	suite.Require().Equal(int64(1002), charlieBalance.Amount.Int64()) // 1000 + 1 + 1

	suite.Require().NoError(err)
}

func (suite *TestSuite) TestInvalidTimes() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: aliceMsgs, Accepted: false}, // bob is creator but has aliceMsgs
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(10000), End: sdkmath.NewUint(10000)},
		},
	}
	createResp, err := CreateProposal(msgServer, ctx, createMsg)
	suite.Require().NoError(err)

	acceptMsg := &types.MsgAcceptProposal{
		Creator: alice,
		Id:      createResp.Id,
	}
	_, err = AcceptProposal(msgServer, ctx, acceptMsg)
	suite.Require().NoError(err)

	acceptMsg = &types.MsgAcceptProposal{
		Creator: bob,
		Id:      createResp.Id,
	}
	_, err = AcceptProposal(msgServer, ctx, acceptMsg)
	suite.Require().NoError(err)

	executeMsg := &types.MsgExecuteProposal{
		Creator: bob,
		Id:      createResp.Id,
	}

	_, err = ExecuteProposal(msgServer, ctx, executeMsg)
	suite.Require().Error(err)

}

func (suite *TestSuite) TestCreatorMustFinalize() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes:          GetFullUintRanges(),
		CreatorMustFinalize: true,
	}
	createResp, err := CreateProposal(msgServer, ctx, createMsg)
	suite.Require().NoError(err)

	// Accept proposals
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: bob, Id: createResp.Id})
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: alice, Id: createResp.Id})

	// Try to execute with non-creator
	_, err = ExecuteProposal(msgServer, ctx, &types.MsgExecuteProposal{Creator: alice, Id: createResp.Id})
	suite.Require().Error(err)

	// Execute with creator
	_, err = ExecuteProposal(msgServer, ctx, &types.MsgExecuteProposal{Creator: bob, Id: createResp.Id})
	suite.Require().NoError(err)
}

func (suite *TestSuite) TestAnyoneCanFinalize() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes:        GetFullUintRanges(),
		AnyoneCanFinalize: true,
	}
	createResp, err := CreateProposal(msgServer, ctx, createMsg)
	suite.Require().NoError(err)

	// Accept proposals
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: bob, Id: createResp.Id})
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: alice, Id: createResp.Id})

	// Execute with non-participant
	_, err = ExecuteProposal(msgServer, ctx, &types.MsgExecuteProposal{Creator: charlie, Id: createResp.Id})
	suite.Require().NoError(err)
}

func (suite *TestSuite) TestContractAddressExecution() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	bobMsgs := GetTestMsgs(bob)
	aliceMsgs := GetTestMsgs(alice)

	// Set useContractAddress to true for both parties
	for _, msg := range bobMsgs {
		msg.UseContractAddress = true
	}
	for _, msg := range aliceMsgs {
		msg.UseContractAddress = true
	}

	createMsg := &types.MsgCreateProposal{
		Creator: bob,
		Parties: []*types.Parties{
			{Creator: bob, MsgsToExecute: bobMsgs, Accepted: false},
			{Creator: alice, MsgsToExecute: aliceMsgs, Accepted: false},
		},
		ValidTimes: GetFullUintRanges(),
	}
	createResp, err := CreateProposal(msgServer, ctx, createMsg)
	suite.Require().NoError(err)

	proposal, _ := suite.app.OffersKeeper.GetProposalFromStore(ctx, createResp.Id)

	// fund the contract
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(proposal.ContractAddress), sdk.NewCoins(sdk.NewInt64Coin("ubadge", 1000)))

	// Accept proposals
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: bob, Id: createResp.Id})
	AcceptProposal(msgServer, ctx, &types.MsgAcceptProposal{Creator: alice, Id: createResp.Id})

	// Execute proposal
	_, err = ExecuteProposal(msgServer, ctx, &types.MsgExecuteProposal{Creator: bob, Id: createResp.Id})
	suite.Require().NoError(err)

	// Check if the messages were executed with the contract address
	// Note: You might need to modify this part based on how you're actually implementing
	// the contract address functionality in your keeper
	proposal, _ = suite.app.OffersKeeper.GetProposalFromStore(ctx, createResp.Id)
	for _, party := range proposal.Parties {
		for _, msg := range party.MsgsToExecute {
			// This is a placeholder check. You should replace this with the actual way
			// you're verifying that the contract address was used
			suite.Require().True(msg.UseContractAddress)
		}
	}
}

func TestOffersKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
