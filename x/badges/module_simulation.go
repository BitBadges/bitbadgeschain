package badges

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	badgessimulation "github.com/bitbadges/bitbadgeschain/x/badges/simulation"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	sdkmath "cosmossdk.io/math"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = badgessimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgTransferBadges          = "op_weight_msg_transfer_badge"
	defaultWeightMsgTransferBadges int = 1000

	opWeightMsgDeleteCollection          = "op_weight_msg_delete_collection"
	defaultWeightMsgDeleteCollection int = 100

	opWeightMsgUpdateUserApprovals          = "op_weight_msg_update_user_approved_transfers"
	defaultWeightMsgUpdateUserApprovals int = 100

	opWeightMsgUniversalUpdateCollection          = "op_weight_msg_update_collection"
	defaultWeightMsgUniversalUpdateCollection int = 1000

	opWeightMsgCreateAddressLists          = "op_weight_msg_create_address_lists"
	defaultWeightMsgCreateAddressLists int = 100

	opWeightMsgCreateCollection = "op_weight_msg_msg_create_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCollection int = 100

	opWeightMsgUpdateCollection = "op_weight_msg_msg_update_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCollection int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	badgesGenesis := types.GenesisState{
		Params:           types.DefaultParams(),
		PortId:           types.PortID,
		NextCollectionId: sdkmath.NewUint(1),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&badgesGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgTransferBadges int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTransferBadges, &weightMsgTransferBadges, nil,
		func(_ *rand.Rand) {
			weightMsgTransferBadges = defaultWeightMsgTransferBadges
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferBadges,
		badgessimulation.SimulateMsgTransferBadges(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteCollection, &weightMsgDeleteCollection, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCollection = defaultWeightMsgDeleteCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCollection,
		badgessimulation.SimulateMsgDeleteCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateUserApprovals int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUserApprovals, &weightMsgUpdateUserApprovals, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUserApprovals = defaultWeightMsgUpdateUserApprovals
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUserApprovals,
		badgessimulation.SimulateMsgUpdateUserApprovals(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUniversalUpdateCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUniversalUpdateCollection, &weightMsgUniversalUpdateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgUniversalUpdateCollection = defaultWeightMsgUniversalUpdateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUniversalUpdateCollection,
		badgessimulation.SimulateMsgUniversalUpdateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateAddressLists int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateAddressLists, &weightMsgCreateAddressLists, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAddressLists = defaultWeightMsgCreateAddressLists
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAddressLists,
		badgessimulation.SimulateMsgCreateAddressLists(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateCollection, &weightMsgCreateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCollection = defaultWeightMsgCreateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCollection,
		badgessimulation.SimulateMsgCreateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCollection, &weightMsgUpdateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCollection = defaultWeightMsgUpdateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCollection,
		badgessimulation.SimulateMsgUpdateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
