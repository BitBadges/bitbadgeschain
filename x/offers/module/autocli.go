package offers

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "bitbadgeschain/api/offers"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},

				{
					RpcMethod: "GetProposal",
					Use:       "proposal",
					Short:     "Shows the proposal",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true,
				},
				//TODO: These are arrays? How to handle that?
				{
					RpcMethod: "CreateProposal",
					Use:       "create-proposal [parties-json]",
					Short:     "Send a createProposal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "parties",
						},
						{
							ProtoField: "validTimes",
						},
					},
				},
				{
					RpcMethod: "AcceptProposal",
					Use:       "accept-proposal [id]",
					Short:     "Send an acceptProposal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "RejectAndDeleteProposal",
					Use:       "reject-and-delete-proposal [id]",
					Short:     "Send a rejectAndDeleteProposal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "ExecuteProposal",
					Use:       "execute-proposal [id]",
					Short:     "Send an executeProposal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
