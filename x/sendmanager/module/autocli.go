package sendmanager

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Balance",
					Use:       "balance [address] [denom]",
					Short:     "Query balance of a specific denom for an address with alias routing",
					Long:      "Query the balance of a specific denomination for an address. Supports both standard coins and alias denoms (e.g., badgeslp:).",
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
					RpcMethod: "SendWithAliasRouting",
					Use:       "send-with-alias-routing [from_address] [to_address] [amount]",
					Short:     "Send coins with alias denom routing",
					Long:      "Send coins from one address to another. Supports both standard coins and alias denoms (e.g., badgeslp:). This mirrors cosmos bank MsgSend but routes through sendmanager.",
					Example:   "bitbadgeschaind tx sendmanager send-with-alias-routing cosmos1abc... cosmos1def... 1000uatom",
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
