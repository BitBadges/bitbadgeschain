package sendmanager

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Query_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // keep autocli `params` next to the custom `balance` command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					// Implemented by the custom CLI command so both the
					// documented positional args and the v36 --address/--denom
					// flags are accepted. Autocli PositionalArgs would drop the flags.
					RpcMethod: "Balance",
					Skip:      true,
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					// Implemented by the custom CLI command so both the
					// documented positional args and the v36 --from/--to-address/--amount
					// flags are accepted. Autocli PositionalArgs would drop the flags.
					RpcMethod: "SendWithAliasRouting",
					Skip:      true,
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
