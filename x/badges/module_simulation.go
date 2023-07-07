package badges

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	badgessimulation "github.com/bitbadges/bitbadgeschain/x/badges/simulation"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
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
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgNewCollection          = "op_weight_msg_new_collection"
	defaultWeightMsgNewCollection int = 1000

	opWeightMsgMintAndDistributeBadges          = "op_weight_msg_mint_and_distribute_badges"
	defaultWeightMsgMintAndDistributeBadges int = 250

	opWeightMsgTransferBadge          = "op_weight_msg_transfer_badge"
	defaultWeightMsgTransferBadge int = 1000

	opWeightMsgUpdateCollectionApprovedTransfers          = "op_weight_msg_update_collection_approved_transfers"
	defaultWeightMsgUpdateCollectionApprovedTransfers int = 100

	opWeightMsgUpdateMetadata          = "op_weight_msg_update_metadata"
	defaultWeightMsgUpdateMetadata int = 100

	opWeightMsgUpdateCollectionPermissions          = "op_weight_msg_update_collection_permissions"
	defaultWeightMsgUpdateCollectionPermissions int = 100

	opWeightMsgUpdateManager          = "op_weight_msg_update_manager"
	defaultWeightMsgUpdateManager int = 100

	opWeightMsgUpdateCustomData          = "op_weight_msg_update_custom_data"
	defaultWeightMsgUpdateCustomData int = 100

	opWeightMsgDeleteCollection = "op_weight_msg_delete_collection"
	defaultWeightMsgDeleteCollection int = 100

	opWeightMsgArchiveCollection = "op_weight_msg_archive_collection"
	defaultWeightMsgArchiveCollection int = 100

	opWeightMsgUpdateUserApprovedTransfers = "op_weight_msg_update_user_approved_transfers"
	defaultWeightMsgUpdateUserApprovedTransfers int = 100

	opWeightMsgUpdateUserPermissions = "op_weight_msg_update_user_permissions"
	defaultWeightMsgUpdateUserPermissions int = 100

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

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgNewCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgNewCollection, &weightMsgNewCollection, nil,
		func(_ *rand.Rand) {
			weightMsgNewCollection = defaultWeightMsgNewCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewCollection,
		badgessimulation.SimulateMsgNewCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgMintAndDistributeBadges int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintAndDistributeBadges, &weightMsgMintAndDistributeBadges, nil,
		func(_ *rand.Rand) {
			weightMsgMintAndDistributeBadges = defaultWeightMsgMintAndDistributeBadges
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintAndDistributeBadges,
		badgessimulation.SimulateMsgMintAndDistributeBadges(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTransferBadge, &weightMsgTransferBadge, nil,
		func(_ *rand.Rand) {
			weightMsgTransferBadge = defaultWeightMsgTransferBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferBadge,
		badgessimulation.SimulateMsgTransferBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCollectionApprovedTransfers int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCollectionApprovedTransfers, &weightMsgUpdateCollectionApprovedTransfers, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCollectionApprovedTransfers = defaultWeightMsgUpdateCollectionApprovedTransfers
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCollectionApprovedTransfers,
		badgessimulation.SimulateMsgUpdateCollectionApprovedTransfers(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateMetadata int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateMetadata, &weightMsgUpdateMetadata, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateMetadata = defaultWeightMsgUpdateMetadata
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateMetadata,
		badgessimulation.SimulateMsgUpdateMetadata(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCollectionPermissions int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCollectionPermissions, &weightMsgUpdateCollectionPermissions, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCollectionPermissions = defaultWeightMsgUpdateCollectionPermissions
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCollectionPermissions,
		badgessimulation.SimulateMsgUpdateCollectionPermissions(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateManager int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateManager, &weightMsgUpdateManager, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateManager = defaultWeightMsgUpdateManager
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateManager,
		badgessimulation.SimulateMsgUpdateManager(am.accountKeeper, am.bankKeeper, am.keeper),
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

	var weightMsgArchiveCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgArchiveCollection, &weightMsgArchiveCollection, nil,
		func(_ *rand.Rand) {
			weightMsgArchiveCollection = defaultWeightMsgArchiveCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgArchiveCollection,
		badgessimulation.SimulateMsgArchiveCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateUserApprovedTransfers int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUserApprovedTransfers, &weightMsgUpdateUserApprovedTransfers, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUserApprovedTransfers = defaultWeightMsgUpdateUserApprovedTransfers
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUserApprovedTransfers,
		badgessimulation.SimulateMsgUpdateUserApprovedTransfers(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateUserPermissions int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUserPermissions, &weightMsgUpdateUserPermissions, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUserPermissions = defaultWeightMsgUpdateUserPermissions
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUserPermissions,
		badgessimulation.SimulateMsgUpdateUserPermissions(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
