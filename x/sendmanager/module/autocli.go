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
					Example:   "bitbadgeschaind query sendmanager balance bb1... badgeslp:1:utoken",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "denom"},
					},
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
					Use:       "send-with-alias-routing [from_key_or_address] [to_address] [amount]",
					Short:     "Send coins with alias denom routing",
					Long:      "Send coins from one address to another. Supports both standard coins and alias denoms (e.g., badgeslp:). This mirrors cosmos bank MsgSend but routes through sendmanager.",
					Example:   "bitbadgeschaind tx sendmanager send-with-alias-routing alice bb1... 1000ubadge 5badgeslp:1:utoken",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "from_address"},
						{ProtoField: "to_address"},
						{ProtoField: "amount", Varargs: true},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
