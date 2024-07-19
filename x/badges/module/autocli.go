package badges

//TODO: Implement this?

// // AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
// func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
// 	return &autocliv1.ModuleOptions{
// 		Query: &autocliv1.ServiceCommandDescriptor{
// 			Service: modulev1.Query_ServiceDesc.ServiceName,
// 			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
// 				{
// 					RpcMethod: "Params",
// 					Use:       "params",
// 					Short:     "Shows the parameters of the module",
// 				},
// 				{
// 					RpcMethod:      "GetCollection",
// 					Use:            "get-collection [collectionId]",
// 					Short:          "Get a collection by ID",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "collectionId"}},
// 				},
// 				// this line is used by ignite scaffolding # autocli/query
// 			},
// 		},
// 		Tx: &autocliv1.ServiceCommandDescriptor{
// 			Service:              modulev1.Msg_ServiceDesc.ServiceName,
// 			EnhanceCustomCommand: true, // only required if you want to use the custom command
// 			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
// 				{
// 					RpcMethod: "UpdateParams",
// 					Skip:      true, // skipped because authority gated
// 				},
// 			},
// 		},
// 	}
// }
