package badges

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
	badgessimulation "github.com/trevormil/bitbadgeschain/x/badges/simulation"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
	opWeightMsgNewBadge = "op_weight_msg_new_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewBadge int = 100

	opWeightMsgNewSubBadge = "op_weight_msg_new_sub_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgNewSubBadge int = 100

	opWeightMsgTransferBadge = "op_weight_msg_transfer_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferBadge int = 100

	opWeightMsgRequestTransferBadge = "op_weight_msg_request_transfer_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestTransferBadge int = 100

	opWeightMsgHandlePendingTransfer = "op_weight_msg_handle_pending_transfer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgHandlePendingTransfer int = 100

	opWeightMsgSetApproval = "op_weight_msg_set_approval"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetApproval int = 100

	opWeightMsgRevokeBadge = "op_weight_msg_revoke_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRevokeBadge int = 100

	opWeightMsgFreezeAddress = "op_weight_msg_freeze_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgFreezeAddress int = 100

	opWeightMsgUpdateUris = "op_weight_msg_update_uris"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateUris int = 100

	opWeightMsgUpdatePermissions = "op_weight_msg_update_permissions"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePermissions int = 100

	opWeightMsgTransferManager = "op_weight_msg_transfer_manager"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferManager int = 100

	opWeightMsgRequestTransferManager = "op_weight_msg_request_transfer_manager"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestTransferManager int = 100

	opWeightMsgSelfDestructBadge = "op_weight_msg_self_destruct_badge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSelfDestructBadge int = 100

	opWeightMsgPruneBalances = "op_weight_msg_prune_balances"
	// TODO: Determine the simulation weight value
	defaultWeightMsgPruneBalances int = 100

	opWeightMsgUpdateBytes = "op_weight_msg_update_bytes"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateBytes int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	badgesGenesis := types.GenesisState{
		Params:      types.DefaultParams(),
		PortId:      types.PortID,
		NextBadgeId: 0,
		Badges:      []*types.BitBadge{},
		Balances:    []*types.UserBalanceInfo{},
		BalanceIds:  []string{},
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
		badgessimulation.SimulateMsgNewBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgNewSubBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgNewSubBadge, &weightMsgNewSubBadge, nil,
		func(_ *rand.Rand) {
			weightMsgNewSubBadge = defaultWeightMsgNewSubBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgNewSubBadge,
		badgessimulation.SimulateMsgNewSubBadge(am.accountKeeper, am.bankKeeper, am.keeper),
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

	var weightMsgRequestTransferBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRequestTransferBadge, &weightMsgRequestTransferBadge, nil,
		func(_ *rand.Rand) {
			weightMsgRequestTransferBadge = defaultWeightMsgRequestTransferBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestTransferBadge,
		badgessimulation.SimulateMsgRequestTransferBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgHandlePendingTransfer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgHandlePendingTransfer, &weightMsgHandlePendingTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgHandlePendingTransfer = defaultWeightMsgHandlePendingTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgHandlePendingTransfer,
		badgessimulation.SimulateMsgHandlePendingTransfer(am.accountKeeper, am.bankKeeper, am.keeper),
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

	var weightMsgRevokeBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRevokeBadge, &weightMsgRevokeBadge, nil,
		func(_ *rand.Rand) {
			weightMsgRevokeBadge = defaultWeightMsgRevokeBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRevokeBadge,
		badgessimulation.SimulateMsgRevokeBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgFreezeAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgFreezeAddress, &weightMsgFreezeAddress, nil,
		func(_ *rand.Rand) {
			weightMsgFreezeAddress = defaultWeightMsgFreezeAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgFreezeAddress,
		badgessimulation.SimulateMsgFreezeAddress(am.accountKeeper, am.bankKeeper, am.keeper),
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

	var weightMsgSelfDestructBadge int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSelfDestructBadge, &weightMsgSelfDestructBadge, nil,
		func(_ *rand.Rand) {
			weightMsgSelfDestructBadge = defaultWeightMsgSelfDestructBadge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSelfDestructBadge,
		badgessimulation.SimulateMsgSelfDestructBadge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgPruneBalances int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgPruneBalances, &weightMsgPruneBalances, nil,
		func(_ *rand.Rand) {
			weightMsgPruneBalances = defaultWeightMsgPruneBalances
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPruneBalances,
		badgessimulation.SimulateMsgPruneBalances(am.accountKeeper, am.bankKeeper, am.keeper),
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

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
