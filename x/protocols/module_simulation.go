package protocols

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	protocolssimulation "github.com/bitbadges/bitbadgeschain/x/protocols/simulation"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = protocolssimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgCreateProtocol = "op_weight_msg_create_protocol"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateProtocol int = 100

	opWeightMsgUpdateProtocol = "op_weight_msg_update_protocol"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateProtocol int = 100

	opWeightMsgDeleteProtocol = "op_weight_msg_delete_protocol"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteProtocol int = 100

	opWeightMsgSetCollectionForProtocol = "op_weight_msg_set_collection_for_protocol"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetCollectionForProtocol int = 100

	opWeightMsgUnsetCollectionForProtocol = "op_weight_msg_unset_collection_for_protocol"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnsetCollectionForProtocol int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	protocolsGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&protocolsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateProtocol int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateProtocol, &weightMsgCreateProtocol, nil,
		func(_ *rand.Rand) {
			weightMsgCreateProtocol = defaultWeightMsgCreateProtocol
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateProtocol,
		protocolssimulation.SimulateMsgCreateProtocol(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateProtocol int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateProtocol, &weightMsgUpdateProtocol, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateProtocol = defaultWeightMsgUpdateProtocol
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateProtocol,
		protocolssimulation.SimulateMsgUpdateProtocol(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteProtocol int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteProtocol, &weightMsgDeleteProtocol, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteProtocol = defaultWeightMsgDeleteProtocol
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteProtocol,
		protocolssimulation.SimulateMsgDeleteProtocol(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetCollectionForProtocol int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetCollectionForProtocol, &weightMsgSetCollectionForProtocol, nil,
		func(_ *rand.Rand) {
			weightMsgSetCollectionForProtocol = defaultWeightMsgSetCollectionForProtocol
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetCollectionForProtocol,
		protocolssimulation.SimulateMsgSetCollectionForProtocol(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnsetCollectionForProtocol int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnsetCollectionForProtocol, &weightMsgUnsetCollectionForProtocol, nil,
		func(_ *rand.Rand) {
			weightMsgUnsetCollectionForProtocol = defaultWeightMsgUnsetCollectionForProtocol
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnsetCollectionForProtocol,
		protocolssimulation.SimulateMsgUnsetCollectionForProtocol(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateProtocol,
			defaultWeightMsgCreateProtocol,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				protocolssimulation.SimulateMsgCreateProtocol(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateProtocol,
			defaultWeightMsgUpdateProtocol,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				protocolssimulation.SimulateMsgUpdateProtocol(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteProtocol,
			defaultWeightMsgDeleteProtocol,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				protocolssimulation.SimulateMsgDeleteProtocol(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetCollectionForProtocol,
			defaultWeightMsgSetCollectionForProtocol,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				protocolssimulation.SimulateMsgSetCollectionForProtocol(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnsetCollectionForProtocol,
			defaultWeightMsgUnsetCollectionForProtocol,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				protocolssimulation.SimulateMsgUnsetCollectionForProtocol(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
