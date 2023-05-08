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

	opWeightMsgUpdateAllowedTransfers          = "op_weight_msg_freeze_address"
	defaultWeightMsgUpdateAllowedTransfers int = 100

	opWeightMsgUpdateUris          = "op_weight_msg_update_uris"
	defaultWeightMsgUpdateUris int = 100

	opWeightMsgUpdatePermissions          = "op_weight_msg_update_permissions"
	defaultWeightMsgUpdatePermissions int = 100

	opWeightMsgTransferManager          = "op_weight_msg_transfer_manager"
	defaultWeightMsgTransferManager int = 100

	opWeightMsgRequestTransferManager          = "op_weight_msg_request_transfer_manager"
	defaultWeightMsgRequestTransferManager int = 100

	opWeightMsgUpdateBytes          = "op_weight_msg_update_bytes"
	defaultWeightMsgUpdateBytes int = 100

	opWeightMsgClaimBadge = "op_weight_msg_claim_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimBadge int = 100

	opWeightMsgDeleteCollection = "op_weight_msg_delete_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCollection int = 100

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
		Claims: 				 	[]*types.Claim{},
		ClaimStoreKeys:  	[]string{},
		NextClaimId: 	 		sdk.NewUint(1),
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

	var weightMsgNewBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgNewBadge, &weightMsgNewBadge, nil,
		func(_ *rand.Rand) {
			weightMsgNewBadge = defaultWeightMsgNewBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewBadge,
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

	var weightMsgSetApproval int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSetApproval, &weightMsgSetApproval, nil,
		func(_ *rand.Rand) {
			weightMsgSetApproval = defaultWeightMsgSetApproval
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetApproval,
		badgessimulation.SimulateMsgSetApproval(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAllowedTransfers int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateAllowedTransfers, &weightMsgUpdateAllowedTransfers, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAllowedTransfers = defaultWeightMsgUpdateAllowedTransfers
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAllowedTransfers,
		badgessimulation.SimulateMsgUpdateAllowedTransfers(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateUris int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateUris, &weightMsgUpdateUris, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateUris = defaultWeightMsgUpdateUris
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateUris,
		badgessimulation.SimulateMsgUpdateUris(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePermissions int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePermissions, &weightMsgUpdatePermissions, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePermissions = defaultWeightMsgUpdatePermissions
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePermissions,
		badgessimulation.SimulateMsgUpdatePermissions(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferManager int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTransferManager, &weightMsgTransferManager, nil,
		func(_ *rand.Rand) {
			weightMsgTransferManager = defaultWeightMsgTransferManager
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferManager,
		badgessimulation.SimulateMsgTransferManager(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRequestTransferManager int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRequestTransferManager, &weightMsgRequestTransferManager, nil,
		func(_ *rand.Rand) {
			weightMsgRequestTransferManager = defaultWeightMsgRequestTransferManager
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestTransferManager,
		badgessimulation.SimulateMsgRequestTransferManager(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateBytes int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateBytes, &weightMsgUpdateBytes, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateBytes = defaultWeightMsgUpdateBytes
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateBytes,
		badgessimulation.SimulateMsgUpdateBytes(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimBadge, &weightMsgClaimBadge, nil,
		func(_ *rand.Rand) {
			weightMsgClaimBadge = defaultWeightMsgClaimBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimBadge,
		badgessimulation.SimulateMsgClaimBadge(am.accountKeeper, am.bankKeeper, am.keeper),
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

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
