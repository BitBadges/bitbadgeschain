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
	opWeightMsgNewBadge          = "op_weight_msg_new_badge"
	defaultWeightMsgNewBadge int = 1000

	opWeightMsgMintAndDistributeBadges          = "op_weight_msg_new_sub_badge"
	defaultWeightMsgMintAndDistributeBadges int = 250

	opWeightMsgTransferBadge          = "op_weight_msg_transfer_badge"
	defaultWeightMsgTransferBadge int = 10000

	opWeightMsgSetApproval          = "op_weight_msg_set_approval"
	defaultWeightMsgSetApproval int = 500

	opWeightMsgUpdateCollectionApprovedTransfers          = "op_weight_msg_freeze_address"
	defaultWeightMsgUpdateCollectionApprovedTransfers int = 100

	opWeightMsgUpdateMetadata          = "op_weight_msg_update_uris"
	defaultWeightMsgUpdateMetadata int = 100

	opWeightMsgUpdateCollectionPermissions          = "op_weight_msg_update_permissions"
	defaultWeightMsgUpdateCollectionPermissions int = 100

	opWeightMsgUpdateManager          = "op_weight_msg_transfer_manager"
	defaultWeightMsgUpdateManager int = 100

	opWeightMsgRequestUpdateManager          = "op_weight_msg_request_transfer_manager"
	defaultWeightMsgRequestUpdateManager int = 100

	opWeightMsgUpdateCustomData          = "op_weight_msg_update_bytes"
	defaultWeightMsgUpdateCustomData int = 100

	opWeightMsgClaimBadge = "op_weight_msg_claim_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimBadge int = 100

	opWeightMsgDeleteCollection = "op_weight_msg_delete_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCollection int = 100

	opWeightMsgArchiveCollection = "op_weight_msg_archive_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgArchiveCollection int = 100

	opWeightMsgForkCollection = "op_weight_msg_fork_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgForkCollection int = 100

	opWeightMsgUpdateUserApprovedTransfers = "op_weight_msg_update_user_approved_transfers"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateUserApprovedTransfers int = 100

	opWeightMsgUpdateUserPermissions = "op_weight_msg_update_user_permissions"
	// TODO: Determine the simulation weight value
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
		NextCollectionId: sdk.NewUint(1),
		Collections:      []*types.BadgeCollection{},
		Balances:         []*types.UserBalanceStore{},
		BalanceStoreKeys: []string{},
		// Claims:           []*types.Claim{},
		// ClaimStoreKeys:   []string{},
		NextClaimId:      sdk.NewUint(1),
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

	// var weightMsgNewBadge int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgNewBadge, &weightMsgNewBadge, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgNewBadge = defaultWeightMsgNewBadge
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgNewBadge,
	// 	badgessimulation.SimulateMsgNewCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgMintAndDistributeBadges int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintAndDistributeBadges, &weightMsgMintAndDistributeBadges, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgMintAndDistributeBadges = defaultWeightMsgMintAndDistributeBadges
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgMintAndDistributeBadges,
	// 	badgessimulation.SimulateMsgMintAndDistributeBadges(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgTransferBadge int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTransferBadge, &weightMsgTransferBadge, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgTransferBadge = defaultWeightMsgTransferBadge
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgTransferBadge,
	// 	badgessimulation.SimulateMsgTransferBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgSetApproval int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetApproval, &weightMsgSetApproval, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgSetApproval = defaultWeightMsgSetApproval
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgSetApproval,
	// 	badgessimulation.SimulateMsgSetApproval(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateCollectionApprovedTransfers int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCollectionApprovedTransfers, &weightMsgUpdateCollectionApprovedTransfers, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateCollectionApprovedTransfers = defaultWeightMsgUpdateCollectionApprovedTransfers
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateCollectionApprovedTransfers,
	// 	badgessimulation.SimulateMsgUpdateCollectionApprovedTransfers(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateMetadata int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateMetadata, &weightMsgUpdateMetadata, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateMetadata = defaultWeightMsgUpdateMetadata
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateMetadata,
	// 	badgessimulation.SimulateMsgUpdateMetadata(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateCollectionPermissions int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCollectionPermissions, &weightMsgUpdateCollectionPermissions, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateCollectionPermissions = defaultWeightMsgUpdateCollectionPermissions
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateCollectionPermissions,
	// 	badgessimulation.SimulateMsgUpdateCollectionPermissions(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateManager int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateManager, &weightMsgUpdateManager, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateManager = defaultWeightMsgUpdateManager
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateManager,
	// 	badgessimulation.SimulateMsgUpdateManager(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgRequestUpdateManager int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRequestUpdateManager, &weightMsgRequestUpdateManager, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgRequestUpdateManager = defaultWeightMsgRequestUpdateManager
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgRequestUpdateManager,
	// 	badgessimulation.SimulateMsgRequestUpdateManager(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateCustomData int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateCustomData, &weightMsgUpdateCustomData, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateCustomData = defaultWeightMsgUpdateCustomData
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateCustomData,
	// 	badgessimulation.SimulateMsgUpdateCustomData(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgClaimBadge int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimBadge, &weightMsgClaimBadge, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgClaimBadge = defaultWeightMsgClaimBadge
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgClaimBadge,
	// 	badgessimulation.SimulateMsgClaimBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgDeleteCollection int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteCollection, &weightMsgDeleteCollection, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgDeleteCollection = defaultWeightMsgDeleteCollection
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgDeleteCollection,
	// 	badgessimulation.SimulateMsgDeleteCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgArchiveCollection int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgArchiveCollection, &weightMsgArchiveCollection, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgArchiveCollection = defaultWeightMsgArchiveCollection
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgArchiveCollection,
	// 	badgessimulation.SimulateMsgArchiveCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgForkCollection int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgForkCollection, &weightMsgForkCollection, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgForkCollection = defaultWeightMsgForkCollection
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgForkCollection,
	// 	badgessimulation.SimulateMsgForkCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateUserApprovedTransfers int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUserApprovedTransfers, &weightMsgUpdateUserApprovedTransfers, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateUserApprovedTransfers = defaultWeightMsgUpdateUserApprovedTransfers
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateUserApprovedTransfers,
	// 	badgessimulation.SimulateMsgUpdateUserApprovedTransfers(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// var weightMsgUpdateUserPermissions int
	// simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUserPermissions, &weightMsgUpdateUserPermissions, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgUpdateUserPermissions = defaultWeightMsgUpdateUserPermissions
	// 	},
	// )
	// operations = append(operations, simulation.NewWeightedOperation(
	// 	weightMsgUpdateUserPermissions,
	// 	badgessimulation.SimulateMsgUpdateUserPermissions(am.accountKeeper, am.bankKeeper, am.keeper),
	// ))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
