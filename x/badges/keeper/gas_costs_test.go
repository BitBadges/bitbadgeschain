package keeper_test

import (
	"github.com/rodaine/table"
)

var tbl = table.New("Function Name", "Gas Consumed")

type GasFunction struct {
	Name      string
	IgnoreGas bool
	F         func()
}

const PRINT_MODE = false

func RunFunctionsAndPrintGasCosts(suite *TestSuite, tbl table.Table, functions []GasFunction) {
	for _, f := range functions {
		if f.IgnoreGas {
			f.F()
		} else {
			startGas := suite.ctx.GasMeter().GasConsumed()
			f.F()
			endGas := suite.ctx.GasMeter().GasConsumed()
			tbl.AddRow(f.Name, endGas-startGas)
		}
	}

	if PRINT_MODE {
		firstColumnFormatter := func(format string, vals ...interface{}) string {
			return ""
		}

		tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()
	}

}

//TODO:

// func (suite *TestSuite) TestGasCosts() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	//Initial variable definitions
// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				BadgeMetadata: "https://example.com/{id}",
// 				Permissions: 46,
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	collectionsToCreate2 := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				BadgeMetadata: "https://example.com/{id}",
// 				Permissions: sdkmath.NewUint(62),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	addressesThroughTenThousand := make([]uint64, 10000)
// 	for j := 0; j < 10000; j++ {
// 		addressesThroughTenThousand[j] = uint64(j)
// 	}

// 	addresses := []uint64{}
// 	for i := 0; i < 1000; i++ {
// 		addresses = append(addresses, aliceAccountNum+uint64(i))
// 	}

// 	collectionsToCreateAllInOne := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// BadgeMetadata: "https://example.com/{id}",
// 				Permissions: sdkmath.NewUint(62),
// 				BadgesToCreate: []*types.BadgeSupplyAndAmount{
// 					{
// 						Supply: 1000000,
// 						Amount: sdkmath.NewUint(1),
// 					},
// 					{
// 						Supply: sdkmath.NewUint(1),
// 						Amount: sdkmath.NewUint(10000),
// 					},
// 					{
// 						Supply: sdkmath.NewUint(10000),
// 						Amount: sdkmath.NewUint(10000),
// 					},
// 				},
// 				FreezeAddressRanges: []*types.UintRange{
// 					{
// 						Start: 1000,
// 						End: sdkmath.NewUint(1000),
// 					},
// 				},
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	RunFunctionsAndPrintGasCosts(suite, tbl, []GasFunction{
// 		{F: func() { CreateCollections(suite, wctx, collectionsToCreate) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(1)) }},
// 		{F: func() { CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: 1000000,
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	}) }},
// 		{F: func() { CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(1),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	}) }},
// 		{F: func() { CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	}) }},
// 		{F: func() { FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(100)0}}) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(1)) }},
// 		{F: func() { FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 9999}}) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(1)) }},
// 		{F: func() { FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), false, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(100)0}}) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(1)) }},
// 		{F: func() { FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), false, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 9999}}) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(1)) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, []uint64{1})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{1})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(3), End: 3}}, []uint64{0})
// 		}},
// 		{IgnoreGas: true, F: func() {
// 			HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(3), End: 3}}, []uint64{0})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}, IgnoreGas: true},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() { HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start:   sdkmath.NewUint(4), End: sdkmath.NewUint(4),}}, []uint64{2}) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{2})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}, IgnoreGas: true},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}},
// 		{F: func() {
// 			HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		}, IgnoreGas: true},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}})
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum) }},
// 		{F: func() { SetApproval(suite, wctx, bob, 10000, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum) }},
// 		{F: func() { SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum) }},
// 		{F: func() { SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}) }},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum) }},
// 		{F: func() { CreateCollections(suite, wctx, collectionsToCreate2) }},
// 		{F: func() { GetCollection(suite, wctx, sdkmath.NewUint(2)) }},
// 		{F: func() { CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	}) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum) }},
// 		{F: func() {
// 			for i := uint64(0); i < 1000; i++ {
// 				TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: i * 2, End: i * 2}}, 0, 0)
// 			}
// 		}},
// 		{F: func() { GetUserBalance(suite, wctx, 2, bobAccountNum) }},
// 		{F: func() { CreateCollections(suite, wctx, collectionsToCreateAllInOne) }},
// 	})
// }

// func (suite *TestSuite) TestGasCostsOldVersionWithRequireChecks() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	for i := 0; i < 1; i++ {
// 		collectionsToCreate := []CollectionsToCreate{
// 			{
// 				Badge: types.MsgNewCollection{
// 					Uri: &types.UriObject{
// 						Uri:                    "example.com/",
// 						Scheme:                 1,
// 						IdxRangeToRemove:       &types.UintRange{},
// 						InsertSubassetBytesIdx: 0,

// 						InsertIdIdx: 10,
// 					},
// 					Permissions: 46,
// 				},
// 				Amount:  sdkmath.NewUint(1),
// 				Creator: bob,
// 			},
// 		}
// 		tbl := table.New("Function Name", "Gas Consumed")

// 		startGas := suite.ctx.GasMeter().GasConsumed()
// 		CreateCollections(suite, wctx, collectionsToCreate)
// 		endGas := suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: 1000000,
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge 1 (Supply 10000)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(1),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge 10000 (Supply 1)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge 10000 (Supply 10000)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(100)0}})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Freeze 1 Address", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		addressesThroughTenThousand := make([]uint64, 10000)
// 		for j := 0; j < 10000; j++ {
// 			addressesThroughTenThousand[j] = uint64(j)
// 		}
// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 9999}})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Freeze 10000 Addresses", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), false, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(100)0}})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Unfreeze 1 Address", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), false, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 9999}})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Unfreeze 10000 Address", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Pending 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, []uint64{1})
// 		suite.Require().Nil(err, "Error accepting badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("HandlePendingTransfer - 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Pending 1000 Diff IDs", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{1})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("HandlePendingTransfer - 1000", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Pending 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(3), End: 3}}, []uint64{0})
// 		suite.Require().Nil(err, "Error accepting badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("HandlePendingTransfer (Reject) - 1", endGas-startGas)
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(3), End: 3}}, []uint64{0})
// 		suite.Require().Nil(err, "Error accepting badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Pending 1000", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("HandlePendingTransfer (Reject) - 1000", endGas-startGas)
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error transferring badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("RequestTransferBadge - 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start:   sdkmath.NewUint(4), End: sdkmath.NewUint(4),}}, []uint64{2})
// 		suite.Require().Nil(err, "Error accepting badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Accept RequestTransfer - 1", endGas-startGas)
// 		// err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		// suite.Require().Nil(err, "Error accepting badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Request Transfer- 1000", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{2})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Accept RequestTransfer - 1000", endGas-startGas)
// 		// err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, true, 0, 0)
// 		// suite.Require().Nil(err, "Error transferring badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("RequestTransferBadge - 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error accepting badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Reject RequestTransfer - 1", endGas-startGas)
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error accepting badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Request Transfer- 1000", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Reject RequestTransfer - 1000", endGas-startGas)
// 		err = HandlePendingTransfers(suite, wctx, alice, sdkmath.NewUint(0), []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, []uint64{0})
// 		suite.Require().Nil(err, "Error accepting badge")

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Revoke Badge 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.UintRange{{Start: 1000, End: sdkmath.NewUint(1)999}})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Revoke Badge 1000", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = SetApproval(suite, wctx, bob, 10000, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Set Approval 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Remove Approval 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}})
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("Set Approval 1000 diff Badge IDs", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bobAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		collectionsToCreate = []CollectionsToCreate{
// 			{
// 				Badge: types.MsgNewCollection{
// 					Uri: &types.UriObject{
// 						Uri:                    "example.com/",
// 						Scheme:                 1,
// 						IdxRangeToRemove:       &types.UintRange{},
// 						InsertSubassetBytesIdx: 0,

// 						InsertIdIdx: 10,
// 					},
// 					Permissions: sdkmath.NewUint(62),
// 				},
// 				Amount:  sdkmath.NewUint(1),
// 				Creator: bob,
// 			},
// 		}

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		CreateCollections(suite, wctx, collectionsToCreate)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetCollection(suite, wctx, sdkmath.NewUint(2))
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetBadge", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 2, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 		suite.Require().Nil(err, "Error creating badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("CreateBadge 10000 (Supply 10000)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Forceful 1", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Forceful 1000 (Same Address, Diff Badge ID)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		addresses := []uint64{}
// 		for i := 0; i < 1000; i++ {
// 			addresses = append(addresses, aliceAccountNum+uint64(i))
// 		}
// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Same BadgeId)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		err = TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.UintRange{{Start: sdkmath.NewUint(0), End: 999}}, 0, 0)
// 		suite.Require().Nil(err, "Error transferring badge")
// 		suite.Require().Nil(err, "Error transferring badge")
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Diff BadgeId)", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		_, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), aliceAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		for i := uint64(0); i < 1000; i++ {
// 			err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.UintRange{{Start: i * 2, End: i * 2}}, 0, 0)
// 			suite.Require().Nil(err, "Error transferring badge")
// 			suite.Require().Nil(err, "Error transferring badge")
// 		}
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("TransferBadge - 1000 alternating badge IDs", endGas-startGas)

// 		startGas = suite.ctx.GasMeter().GasConsumed()
// 		UserBalance, _ := GetUserBalance(suite, wctx, 2, bobAccountNum)
// 		endGas = suite.ctx.GasMeter().GasConsumed()
// 		tbl.AddRow("GetUserBalance", endGas-startGas)

// 		_ = UserBalance

// 		firstColumnFormatter := func(format string, vals ...interface{}) string {
// 			return ""
// 		}

// 		tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()

// 	}
// }
