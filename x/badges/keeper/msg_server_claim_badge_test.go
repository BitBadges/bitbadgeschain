package keeper_test

//TODO: lol

// func (suite *TestSuite) TestSendAllToClaimsAndClaim() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)
// 	err := *new(error)

// 	// rootHash, merkleProofs := merkle.ProofsFromByteSlices([][]byte{[]byte(alice), []byte(bob), []byte(charlie), []byte(charlie)})

// 	aliceLeaf := "-" + alice + "-1-0-0"
// 	bobLeaf := "-" + bob + "-1-0-0"
// 	charlieLeaf := "-" + charlie + "-1-0-0"
// 	// aliceLeaf := alice
// 	// bobLeaf := bob
// 	// charlieLeaf := charlie

// 	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
// 	leafHashes := make([][]byte, len(leafs))
// 	for i, leaf := range leafs {
// 		initialHash := sha256.Sum256(leaf)
// 		leafHashes[i] = initialHash[:]
// 		for j := 0; j < 32; j++ {
// 			print(leafHashes[i][j])
// 			print(" ")
// 		}
// 		println()

// 		// println("leafHashes[i]: ", string(leafHashes[i]))
// 	}
// 	println()

// 	levelTwoHashes := make([][]byte, 2)
// 	for i := 0; i < len(leafHashes); i += 2 {
// 		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
// 		levelTwoHashes[i/2] = iHash[:]
// 		for j := 0; j < 32; j++ {
// 			print(levelTwoHashes[i/2][j])
// 			print(" ")
// 		}
// 		println()
// 	}
// 	println()

// 	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
// 	rootHash := rootHashI[:]

// 	for j := 0; j < 32; j++ {
// 		print(rootHash[j])
// 		print(" ")
// 	}

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// 				BadgeUris:            []*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(62),
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	claimToAdd := types.Claim{
// 		Balances: []*types.Balance{
// 			{
// 				Amount: sdk.NewUint(10),
// 				BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 			},
// 		},
// 		BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 		IncrementIdsBy: 0,
// 		AmountPerClaim: 1,
// 		Data:       hex.EncodeToString(rootHash),
// 		Type: 	 	uint64(types.ClaimType_MerkleTree),
// 		Uri: "",
// 		TimeRange: &types.IdRange{
// 			Start: sdk.NewUint(0),
// 			End:   sdk.NewUint(math.MaxUint64),
// 		},
// 	}

// 	err = CreateBadges(suite, wctx, bob, sdk.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10),
// 			Amount: sdk.NewUint(1),
// 		},
// 	},
// 		[]*types.Transfer{},
// 		[]*types.Claim{
// 			&claimToAdd,
// 		},
// 		"https://example.com",
// 		[]*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 	)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount: sdk.NewUint(10),
// 		},
// 	}, badge.MaxSupplys)

// 	claim := badge.Claims[0]

// 	err = ClaimBadge(suite, wctx, bob, sdk.NewUint(0), 0, &types.ClaimProof{
// 			Leaf: aliceLeaf,
// 			Aunts: []*types.ClaimProofItem{
// 				{
// 					Aunt: hex.EncodeToString(leafHashes[1]),
// 					OnRight: true,
// 				},
// 				{
// 					Aunt: hex.EncodeToString(levelTwoHashes[1]),
// 					OnRight: true,
// 				},
// 			},
// 		},
// 		"",
// 		&types.IdRange{
// 			Start: sdk.NewUint(0),
// 			End:   sdk.NewUint(math.MaxUint64),
// 		},
// 	)
// 	suite.Require().Nil(err, "Error claiming badge")

// 	aliceBalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), aliceAccountNum)
// 	suite.Require().Equal(uint64(1), aliceBalance.Balances[0].Balance)
// 	suite.Require().Equal([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, aliceBalance.Balances[0].BadgeIds)

// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
// 	claim = badge.Claims[0]
// 	suite.Require().Equal(uint64(9), claim.Balances[0].Balance)
// 	// suite.Require().Equal([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, aliceBalance.Balances[0].BadgeIds)
// }

// func (suite *TestSuite) TestSendAllToClaimsAccountTypeInvalid() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)
// 	err := *new(error)

// 	aliceLeaf := "-" + alice + "-1-0-0"
// 	bobLeaf := "-" + bob + "-1-0-0"
// 	charlieLeaf := "-" + charlie + "-1-0-0"

// 	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
// 	leafHashes := make([][]byte, len(leafs))
// 	for i, leaf := range leafs {
// 		initialHash := sha256.Sum256(leaf)
// 		leafHashes[i] = initialHash[:]
// 	}

// 	levelTwoHashes := make([][]byte, 2)
// 	for i := 0; i < len(leafHashes); i += 2 {
// 		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
// 		levelTwoHashes[i/2] = iHash[:]
// 	}

// 	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
// 	rootHash := rootHashI[:]

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// 				BadgeUris:            []*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(62),
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	claimToAdd := types.Claim{
// 		Balances: []*types.Balance{
// 			{
// 				Amount: sdk.NewUint(10),
// 				BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 			},
// 		},
// 		BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 		IncrementIdsBy: 0,
// 		AmountPerClaim: 1,
// 		Data:       hex.EncodeToString(rootHash),
// 		Type: 	 	uint64(types.ClaimType_MerkleTree),
// 		Uri: "",
// 		TimeRange: &types.IdRange{
// 			Start: sdk.NewUint(0),
// 			End:   sdk.NewUint(math.MaxUint64),
// 		},
// 	}

// 	err = CreateBadges(suite, wctx, bob, sdk.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10),
// 			Amount: sdk.NewUint(1),
// 		},
// 	},
// 		[]*types.Transfer{},
// 		[]*types.Claim{
// 			&claimToAdd,
// 		}, "https://example.com",
// 		[]*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 	)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount: sdk.NewUint(10),
// 		},
// 	}, badge.MaxSupplys)

// 	claim := badge.Claims[0]
// 	suite.Require().Equal(&claimToAdd, claim)

// 	err = ClaimBadge(suite, wctx, alice, sdk.NewUint(0), 0, &types.ClaimProof{
// 		Leaf: "",
// 		Aunts: []*types.ClaimProofItem{

// 			{
// 				Aunt: hex.EncodeToString(leafHashes[1]),
// 				OnRight: true,
// 			},
// 			{
// 				Aunt: hex.EncodeToString(levelTwoHashes[1]),
// 				OnRight: true,
// 			},
// 		},
// 	}, "", &types.IdRange{
// 		Start: sdk.NewUint(0),
// 		End:   sdk.NewUint(math.MaxUint64),
// 	})
// 	suite.Require().EqualError(err, keeper.ErrRootHashInvalid.Error())
// }

// func (suite *TestSuite) TestSendAllToClaimsAccountTypeCodes() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)
// 	err := *new(error)

// 	aliceLeaf := "-" + alice + "-1-0-0"
// 	bobLeaf := "-" + bob + "-1-0-0"
// 	charlieLeaf := "-" + charlie + "-1-0-0"

// 	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
// 	leafHashes := make([][]byte, len(leafs))
// 	for i, leaf := range leafs {
// 		initialHash := sha256.Sum256(leaf)
// 		leafHashes[i] = initialHash[:]
// 	}

// 	levelTwoHashes := make([][]byte, 2)
// 	for i := 0; i < len(leafHashes); i += 2 {
// 		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
// 		levelTwoHashes[i/2] = iHash[:]
// 	}

// 	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
// 	rootHash := rootHashI[:]

// 	// output := tmhash.Sum(append([]byte{0}, []byte("hello")...))
// 	// // output := tmhash.Sum([]byte("hello"))
// 	// for i := 0; i < len(merkleProofs[0].LeafHash); i++ {
// 	// 	println(merkleProofs[0].LeafHash[i], output[i])
// 	// }

// 	// 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// 				BadgeUris:            []*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(62),
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	claimToAdd := types.Claim{
// 		Balances: []*types.Balance{
// 			{
// 				Amount: sdk.NewUint(10),
// 				BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 			},
// 		},
// 		BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 		IncrementIdsBy: 1,
// 		AmountPerClaim: 1,
// 		Data:       hex.EncodeToString(rootHash),
// 		Type: 	 	uint64(types.ClaimType_FirstCome),
// 		Uri: "",
// 		TimeRange: &types.IdRange{
// 			Start: sdk.NewUint(0),
// 			End:   sdk.NewUint(math.MaxUint64),
// 		},
// 	}

// 	err = CreateBadges(suite, wctx, bob, sdk.NewUint(0), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10),
// 			Amount: sdk.NewUint(1),
// 		},
// 	},
// 		[]*types.Transfer{},
// 		[]*types.Claim{
// 			&claimToAdd,
// 		}, "https://example.com",
// 		[]*types.BadgeUri{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End: sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 	)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount: sdk.NewUint(10),
// 		},
// 	}, badge.MaxSupplys)

// 	claim := badge.Claims[0]
// 	suite.Require().Equal(&claimToAdd, claim)

// 	err = ClaimBadge(suite, wctx, alice, sdk.NewUint(0), 0, &types.ClaimProof{
// 		Leaf: aliceLeaf,
// 		Aunts: []*types.ClaimProofItem{
// 			{
// 				Aunt: hex.EncodeToString(leafHashes[1]),
// 				OnRight: true,
// 			},
// 			{
// 				Aunt: hex.EncodeToString(levelTwoHashes[1]),
// 				OnRight: true,
// 			},
// 		},
// 	}, "", &types.IdRange{
// 		Start: sdk.NewUint(0),
// 		End:   sdk.NewUint(math.MaxUint64),
// 	})
// 	suite.Require().Nil(err, "Error claiming badge")

// 	aliceBalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), aliceAccountNum)
// 	suite.Require().Equal(uint64(1), aliceBalance.Balances[0].Balance)
// 	suite.Require().Equal([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, aliceBalance.Balances[0].BadgeIds)

// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
// 	claim = badge.Claims[0]
// 	suite.Require().Equal(uint64(9), claim.Balances[0].Balance)
// 	suite.Require().Equal([]*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2),}}, claim.BadgeIds)
// }
