package badges

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/bitbadges/bitbadgeschain/x/badges/testutil/sample"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	badgessimulation "github.com/bitbadges/bitbadgeschain/x/badges/simulation"
)

// avoid unused import issue
var (
	_ = badgessimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateCollection        = "op_weight_msg_create_collection"
	opWeightMsgUniversalUpdateCollection = "op_weight_msg_universal_update_collection"
	opWeightMsgDeleteCollection         = "op_weight_msg_delete_collection"
	opWeightMsgTransferTokens           = "op_weight_msg_transfer_tokens"
	opWeightMsgUpdateUserApprovals      = "op_weight_msg_update_user_approvals"
	opWeightMsgSetIncomingApproval      = "op_weight_msg_set_incoming_approval"
	opWeightMsgSetOutgoingApproval      = "op_weight_msg_set_outgoing_approval"
	opWeightMsgPurgeApprovals           = "op_weight_msg_purge_approvals"
	opWeightMsgCreateAddressLists       = "op_weight_msg_create_address_lists"
	opWeightMsgSetDynamicStoreValue     = "op_weight_msg_set_dynamic_store_value"
	
	// Default weights - higher for more common operations
	defaultWeightMsgCreateCollection         = 100
	defaultWeightMsgUniversalUpdateCollection = 80
	defaultWeightMsgDeleteCollection          = 20
	defaultWeightMsgTransferTokens            = 200 // Most common operation
	defaultWeightMsgUpdateUserApprovals       = 60
	defaultWeightMsgSetIncomingApproval       = 40
	defaultWeightMsgSetOutgoingApproval       = 40
	defaultWeightMsgPurgeApprovals            = 10
	defaultWeightMsgCreateAddressLists        = 30
	defaultWeightMsgSetDynamicStoreValue      = 20
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	badgesGenesis := types.DefaultGenesis()
	// Use default genesis which already initializes NextCollectionId to 1
	// this line is used by starport scaffolding # simapp/module/genesisState
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(badgesGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCollection int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCollection, &weightMsgCreateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCollection = defaultWeightMsgCreateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCollection,
		badgessimulation.SimulateMsgCreateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUniversalUpdateCollection int
	simState.AppParams.GetOrGenerate(opWeightMsgUniversalUpdateCollection, &weightMsgUniversalUpdateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgUniversalUpdateCollection = defaultWeightMsgUniversalUpdateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUniversalUpdateCollection,
		badgessimulation.SimulateMsgUniversalUpdateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCollection int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCollection, &weightMsgDeleteCollection, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCollection = defaultWeightMsgDeleteCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCollection,
		badgessimulation.SimulateMsgDeleteCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferTokens, &weightMsgTransferTokens, nil,
		func(_ *rand.Rand) {
			weightMsgTransferTokens = defaultWeightMsgTransferTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferTokens,
		badgessimulation.SimulateMsgTransferTokens(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateUserApprovals int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateUserApprovals, &weightMsgUpdateUserApprovals, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUserApprovals = defaultWeightMsgUpdateUserApprovals
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUserApprovals,
		badgessimulation.SimulateMsgUpdateUserApprovals(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetIncomingApproval int
	simState.AppParams.GetOrGenerate(opWeightMsgSetIncomingApproval, &weightMsgSetIncomingApproval, nil,
		func(_ *rand.Rand) {
			weightMsgSetIncomingApproval = defaultWeightMsgSetIncomingApproval
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetIncomingApproval,
		badgessimulation.SimulateMsgSetIncomingApproval(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetOutgoingApproval int
	simState.AppParams.GetOrGenerate(opWeightMsgSetOutgoingApproval, &weightMsgSetOutgoingApproval, nil,
		func(_ *rand.Rand) {
			weightMsgSetOutgoingApproval = defaultWeightMsgSetOutgoingApproval
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetOutgoingApproval,
		badgessimulation.SimulateMsgSetOutgoingApproval(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgPurgeApprovals int
	simState.AppParams.GetOrGenerate(opWeightMsgPurgeApprovals, &weightMsgPurgeApprovals, nil,
		func(_ *rand.Rand) {
			weightMsgPurgeApprovals = defaultWeightMsgPurgeApprovals
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPurgeApprovals,
		badgessimulation.SimulateMsgPurgeApprovals(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateAddressLists int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAddressLists, &weightMsgCreateAddressLists, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAddressLists = defaultWeightMsgCreateAddressLists
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAddressLists,
		badgessimulation.SimulateMsgCreateAddressLists(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetDynamicStoreValue int
	simState.AppParams.GetOrGenerate(opWeightMsgSetDynamicStoreValue, &weightMsgSetDynamicStoreValue, nil,
		func(_ *rand.Rand) {
			weightMsgSetDynamicStoreValue = defaultWeightMsgSetDynamicStoreValue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetDynamicStoreValue,
		badgessimulation.SimulateMsgSetDynamicStoreValue(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
