package badges

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"
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
	opWeightMsgCreateCollection          = "op_weight_msg_create_collection"
	opWeightMsgUniversalUpdateCollection = "op_weight_msg_universal_update_collection"
	opWeightMsgDeleteCollection          = "op_weight_msg_delete_collection"
	opWeightMsgTransferTokens            = "op_weight_msg_transfer_tokens"
	opWeightMsgUpdateUserApprovals       = "op_weight_msg_update_user_approvals"
	opWeightMsgSetIncomingApproval       = "op_weight_msg_set_incoming_approval"
	opWeightMsgSetOutgoingApproval       = "op_weight_msg_set_outgoing_approval"
	opWeightMsgPurgeApprovals            = "op_weight_msg_purge_approvals"
	opWeightMsgCreateAddressLists        = "op_weight_msg_create_address_lists"
	opWeightMsgSetDynamicStoreValue      = "op_weight_msg_set_dynamic_store_value"

	// Default weights - higher for more common operations
	defaultWeightMsgCreateCollection          = 100
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
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	badgesGenesis := types.DefaultGenesis()

	// Pre-create collections, dynamic stores, and balances for better simulation starting state
	// This ensures simulation operations have valid state to work with from the start
	if len(simState.Accounts) > 0 {
		// Convert accounts to simtypes.Account format for SetupSimulationState
		simAccounts := make([]simtypes.Account, len(simState.Accounts))
		for i, acc := range simState.Accounts {
			simAccounts[i] = simtypes.Account{
				Address: acc.Address,
				PrivKey: acc.PrivKey,
			}
		}

		// Create a temporary context and keeper to setup state
		// Note: In actual simulation, this will be done during InitGenesis
		// But we can pre-populate the genesis state with collections and dynamic stores
		// The actual state setup will happen when InitGenesis is called
		r := rand.New(rand.NewSource(simState.Rand.Int63()))

		// Pre-create some collections in genesis state
		// Collections will be created during simulation operations, not in genesis
		// This ensures simulation starts with a clean state but operations can create resources

		// Pre-create some dynamic stores in genesis state
		for i := 0; i < badgessimulation.DefaultSimDynamicStoreCount && i < len(simAccounts); i++ {
			creator := simAccounts[i].Address.String()
			defaultValue := r.Intn(2) == 0

			// Create dynamic store that will be initialized in InitGenesis
			dynamicStore := &types.DynamicStore{
				StoreId:       sdkmath.NewUint(uint64(i + 1)),
				CreatedBy:     creator,
				DefaultValue:  defaultValue,
				GlobalEnabled: true,
			}
			badgesGenesis.DynamicStores = append(badgesGenesis.DynamicStores, dynamicStore)
		}

		// Update next IDs to reflect pre-created resources
		if len(badgesGenesis.DynamicStores) > 0 {
			badgesGenesis.NextDynamicStoreId = sdkmath.NewUint(uint64(len(badgesGenesis.DynamicStores) + 1))
		}
	}

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
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgCreateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgUniversalUpdateCollection int
	simState.AppParams.GetOrGenerate(opWeightMsgUniversalUpdateCollection, &weightMsgUniversalUpdateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgUniversalUpdateCollection = defaultWeightMsgUniversalUpdateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUniversalUpdateCollection,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgUniversalUpdateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgDeleteCollection int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCollection, &weightMsgDeleteCollection, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCollection = defaultWeightMsgDeleteCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCollection,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgDeleteCollection(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgTransferTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferTokens, &weightMsgTransferTokens, nil,
		func(_ *rand.Rand) {
			weightMsgTransferTokens = defaultWeightMsgTransferTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferTokens,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgTransferTokens(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgUpdateUserApprovals int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateUserApprovals, &weightMsgUpdateUserApprovals, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUserApprovals = defaultWeightMsgUpdateUserApprovals
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUserApprovals,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgUpdateUserApprovals(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgSetIncomingApproval int
	simState.AppParams.GetOrGenerate(opWeightMsgSetIncomingApproval, &weightMsgSetIncomingApproval, nil,
		func(_ *rand.Rand) {
			weightMsgSetIncomingApproval = defaultWeightMsgSetIncomingApproval
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetIncomingApproval,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgSetIncomingApproval(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgSetOutgoingApproval int
	simState.AppParams.GetOrGenerate(opWeightMsgSetOutgoingApproval, &weightMsgSetOutgoingApproval, nil,
		func(_ *rand.Rand) {
			weightMsgSetOutgoingApproval = defaultWeightMsgSetOutgoingApproval
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetOutgoingApproval,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgSetOutgoingApproval(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgPurgeApprovals int
	simState.AppParams.GetOrGenerate(opWeightMsgPurgeApprovals, &weightMsgPurgeApprovals, nil,
		func(_ *rand.Rand) {
			weightMsgPurgeApprovals = defaultWeightMsgPurgeApprovals
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPurgeApprovals,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgPurgeApprovals(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgCreateAddressLists int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAddressLists, &weightMsgCreateAddressLists, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAddressLists = defaultWeightMsgCreateAddressLists
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAddressLists,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgCreateAddressLists(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	var weightMsgSetDynamicStoreValue int
	simState.AppParams.GetOrGenerate(opWeightMsgSetDynamicStoreValue, &weightMsgSetDynamicStoreValue, nil,
		func(_ *rand.Rand) {
			weightMsgSetDynamicStoreValue = defaultWeightMsgSetDynamicStoreValue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetDynamicStoreValue,
		badgessimulation.MultiRunOperation(
			badgessimulation.SimulateMsgSetDynamicStoreValue(am.accountKeeper, am.bankKeeper, am.keeper),
			badgessimulation.DefaultMultiRunAttempts,
		),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
