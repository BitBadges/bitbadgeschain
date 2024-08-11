package offers

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"bitbadgeschain/testutil/sample"
	offerssimulation "bitbadgeschain/x/offers/simulation"
	"bitbadgeschain/x/offers/types"
)

// avoid unused import issue
var (
	_ = offerssimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateProposal = "op_weight_msg_create_proposal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateProposal int = 100

	opWeightMsgAcceptProposal = "op_weight_msg_accept_proposal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAcceptProposal int = 100

	opWeightMsgRejectAndDeleteProposal = "op_weight_msg_reject_and_delete_proposal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRejectAndDeleteProposal int = 100

	opWeightMsgExecuteProposal = "op_weight_msg_execute_proposal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgExecuteProposal int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	offersGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&offersGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateProposal int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateProposal, &weightMsgCreateProposal, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProposal = defaultWeightMsgCreateProposal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateProposal,
		offerssimulation.SimulateMsgCreateProposal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAcceptProposal int
	simState.AppParams.GetOrGenerate(opWeightMsgAcceptProposal, &weightMsgAcceptProposal, nil,
		func(_ *rand.Rand) {
			weightMsgAcceptProposal = defaultWeightMsgAcceptProposal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAcceptProposal,
		offerssimulation.SimulateMsgAcceptProposal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRejectAndDeleteProposal int
	simState.AppParams.GetOrGenerate(opWeightMsgRejectAndDeleteProposal, &weightMsgRejectAndDeleteProposal, nil,
		func(_ *rand.Rand) {
			weightMsgRejectAndDeleteProposal = defaultWeightMsgRejectAndDeleteProposal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRejectAndDeleteProposal,
		offerssimulation.SimulateMsgRejectAndDeleteProposal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExecuteProposal int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteProposal, &weightMsgExecuteProposal, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteProposal = defaultWeightMsgExecuteProposal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteProposal,
		offerssimulation.SimulateMsgExecuteProposal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateProposal,
			defaultWeightMsgCreateProposal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				offerssimulation.SimulateMsgCreateProposal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAcceptProposal,
			defaultWeightMsgAcceptProposal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				offerssimulation.SimulateMsgAcceptProposal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRejectAndDeleteProposal,
			defaultWeightMsgRejectAndDeleteProposal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				offerssimulation.SimulateMsgRejectAndDeleteProposal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgExecuteProposal,
			defaultWeightMsgExecuteProposal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				offerssimulation.SimulateMsgExecuteProposal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
