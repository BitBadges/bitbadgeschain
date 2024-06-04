package ante_test

//TODO:
// func (suite AnteTestSuite) TestAnteHandler() {
// 	suite.enableFeemarket = false
// 	suite.SetupTest() // reset

// 	addr, privKey := tests.NewAddrKey()
// 	// to := tests.GenerateAddress()

// 	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
// 	suite.Require().NoError(acc.SetSequence(1))
// 	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

// 	testCases := []struct {
// 		name      string
// 		txFn      func() sdk.Tx
// 		checkTx   bool
// 		reCheckTx bool
// 		expPass   bool
// 	}{
// 		{
// 			"success - DeliverTx EIP712 signed Cosmos Tx with MsgSend",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9000-1", gas, amount)
// 				return txBuilder.GetTx()
// 			}, false, false, true,
// 		},
// 		{
// 			"success - DeliverTx EIP712 signed Cosmos Tx with DelegateMsg",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				coinAmount := sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0))
// 				amount := sdk.NewCoins(coinAmount)
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgDelegate(from, privKey, "ethermint_9000-1", gas, amount)
// 				return txBuilder.GetTx()
// 			}, false, false, true,
// 		},
// 		{
// 			"fails - DeliverTx EIP712 signed Cosmos Tx with wrong Chain ID",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9002-1", gas, amount)
// 				return txBuilder.GetTx()
// 			}, false, false, false,
// 		},
// 		{
// 			"fails - DeliverTx EIP712 signed Cosmos Tx with different gas fees",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9001-1", gas, amount)
// 				txBuilder.SetGasLimit(uint64(300000))
// 				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(30))))
// 				return txBuilder.GetTx()
// 			}, false, false, false,
// 		},
// 		{
// 			"fails - DeliverTx EIP712 signed Cosmos Tx with empty signature",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9001-1", gas, amount)
// 				sigsV2 := signing.SignatureV2{}
// 				txBuilder.SetSignatures(sigsV2)
// 				return txBuilder.GetTx()
// 			}, false, false, false,
// 		},
// 		{
// 			"fails - DeliverTx EIP712 signed Cosmos Tx with invalid sequence",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9001-1", gas, amount)
// 				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
// 				suite.Require().NoError(err)
// 				sigsV2 := signing.SignatureV2{
// 					PubKey: privKey.PubKey(),
// 					Data: &signing.SingleSignatureData{
// 						SignMode: signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
// 					},
// 					Sequence: nonce - 1,
// 				}
// 				txBuilder.SetSignatures(sigsV2)
// 				return txBuilder.GetTx()
// 			}, false, false, false,
// 		},
// 		{
// 			"fails - DeliverTx EIP712 signed Cosmos Tx with invalid signMode",
// 			func() sdk.Tx {
// 				from := acc.GetAddress()
// 				amount := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdkmath.NewInt(0)))
// 				gas := uint64(200000)
// 				txBuilder := suite.CreateTestEIP712TxBuilderMsgSend(from, privKey, "ethermint_9001-1", gas, amount)
// 				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
// 				suite.Require().NoError(err)
// 				sigsV2 := signing.SignatureV2{
// 					PubKey: privKey.PubKey(),
// 					Data: &signing.SingleSignatureData{
// 						SignMode: signing.SignMode_SIGN_MODE_UNSPECIFIED,
// 					},
// 					Sequence: nonce,
// 				}
// 				txBuilder.SetSignatures(sigsV2)
// 				return txBuilder.GetTx()
// 			}, false, false, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.Run(tc.name, func() {
// 			suite.ctx = suite.ctx.WithIsCheckTx(tc.checkTx).WithIsReCheckTx(tc.reCheckTx)

// 			// expConsumed := params.TxGasContractCreation + params.TxGas
// 			_, err := suite.anteHandler(suite.ctx, tc.txFn(), false)

// 			// suite.Require().Equal(consumed, ctx.GasMeter().GasConsumed())

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 				// suite.Require().Equal(int(expConsumed), int(suite.ctx.GasMeter().GasConsumed()))
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }
