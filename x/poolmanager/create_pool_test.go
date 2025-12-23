package poolmanager_test

import (
	"fmt"
	"math"
	"reflect"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	stableswap "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/stableswap"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// TestCreatePool tests that all possible pools are created correctly.
func (s *KeeperTestSuite) TestCreatePool() {
	var (
		validBalancerPoolMsg = balancer.NewMsgCreateBalancerPool(s.TestAccs[0], balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()), []balancer.PoolAsset{
			{
				Token:  sdk.NewCoin(FOO, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
			{
				Token:  sdk.NewCoin(BAR, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
		})

		invalidBalancerPoolMsg = balancer.NewMsgCreateBalancerPool(s.TestAccs[0], balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.NewDecWithPrec(1, 2)), []balancer.PoolAsset{
			{
				Token:  sdk.NewCoin(FOO, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
			{
				Token:  sdk.NewCoin(BAR, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
		})

		DefaultStableswapLiquidity = sdk.NewCoins(
			sdk.NewCoin(FOO, defaultInitPoolAmount),
			sdk.NewCoin(BAR, defaultInitPoolAmount),
		)

		validStableswapPoolMsg = stableswap.NewMsgCreateStableswapPool(s.TestAccs[0], stableswap.PoolParams{SwapFee: osmomath.NewDec(0), ExitFee: osmomath.NewDec(0)}, DefaultStableswapLiquidity, []uint64{}, "")

		invalidStableswapPoolMsg = stableswap.NewMsgCreateStableswapPool(s.TestAccs[0], stableswap.PoolParams{SwapFee: osmomath.NewDec(0), ExitFee: osmomath.NewDecWithPrec(1, 2)}, DefaultStableswapLiquidity, []uint64{}, "")

		// validTransmuterCodeId = uint64(1)

		defaultFundAmount = sdk.NewCoins(sdk.NewCoin(FOO, defaultInitPoolAmount.Mul(osmomath.NewInt(2))), sdk.NewCoin(BAR, defaultInitPoolAmount.Mul(osmomath.NewInt(2))))
	)

	tests := []struct {
		name                                 string
		creatorFundAmount                    sdk.Coins
		isPermissionlessPoolCreationDisabled bool
		msg                                  types.CreatePoolMsg
		expectedModuleType                   reflect.Type
		expectError                          bool
	}{
		{
			name:               "first balancer pool - success",
			creatorFundAmount:  defaultFundAmount,
			msg:                validBalancerPoolMsg,
			expectedModuleType: gammKeeperType,
		},
		{
			name:               "second balancer pool - success",
			creatorFundAmount:  defaultFundAmount,
			msg:                validBalancerPoolMsg,
			expectedModuleType: gammKeeperType,
		},
		{
			name:               "stableswap pool - success",
			creatorFundAmount:  defaultFundAmount,
			msg:                validStableswapPoolMsg,
			expectedModuleType: gammKeeperType,
		},
		// {
		// 	name:               "concentrated pool - success",
		// 	creatorFundAmount:  defaultFundAmount,
		// 	msg:                validConcentratedPoolMsg,
		// 	expectedModuleType: concentratedKeeperType,
		// },
		// {
		// 	name:               "cosmwasm pool - success",
		// 	creatorFundAmount:  defaultFundAmount,
		// 	msg:                validCWPoolMsg,
		// 	expectedModuleType: cosmwasmKeeperType,
		// },
		{
			name:               "error: balancer pool with non zero exit fee",
			creatorFundAmount:  defaultFundAmount,
			msg:                invalidBalancerPoolMsg,
			expectedModuleType: gammKeeperType,
			expectError:        true,
		},
		{
			name:               "error: stableswap pool with non zero exit fee",
			creatorFundAmount:  defaultFundAmount,
			msg:                invalidStableswapPoolMsg,
			expectedModuleType: gammKeeperType,
			expectError:        true,
		},
		// {
		// 	name:                                 "error: pool creation is disabled for concentrated pool via param",
		// 	creatorFundAmount:                    defaultFundAmount,
		// 	isPermissionlessPoolCreationDisabled: true,
		// 	msg:                                  validConcentratedPoolMsg,
		// 	expectedModuleType:                   concentratedKeeperType,
		// 	expectError:                          true,
		// },
	}

	// setup cosmwasm pool
	// codeId := s.StoreCosmWasmPoolContractCode(apptesting.TransmuterContractName)
	// s.Require().Equal(validTransmuterCodeId, codeId)
	// s.App.CosmwasmPoolKeeper.WhitelistCodeId(s.Ctx, codeId)

	execModes := map[string]sdk.ExecMode{
		"check":    sdk.ExecModeCheck,
		"finalize": sdk.ExecModeFinalize,
	}
	totalTestCount := 0
	for _, tc := range tests {
		for execModeName, execMode := range execModes {
			totalTestCount++
			s.Run(fmt.Sprintf("%s-%s", tc.name, execModeName), func() {
				s.Ctx = s.Ctx.WithExecMode(execMode)

				poolmanagerKeeper := s.App.PoolManagerKeeper
				ctx := s.Ctx

				// poolCreationFee := poolmanagerKeeper.GetParams(s.Ctx).PoolCreationFee
				poolCreationFee := sdk.Coins{}
				s.FundAcc(s.TestAccs[0], append(tc.creatorFundAmount, poolCreationFee...))

				poolId, err := poolmanagerKeeper.CreatePool(ctx, tc.msg)

				if tc.expectError {
					s.Require().Error(err)
					return
				}

				// Validate pool.
				s.Require().NoError(err)
				s.Require().Equal(uint64(totalTestCount), poolId)

				// // Validate that mapping pool id -> module type has been persisted.
				// swapModule, err := poolmanagerKeeper.GetPoolModule(ctx, poolId)
				// s.Require().NoError(err)
				// s.Require().Equal(tc.expectedModuleType, reflect.TypeOf(&swapModule))
			})
		}
	}
}

// // Tests that only poolmanager as a pool creator can create a pool via CreatePoolZeroLiquidityNoCreationFee
// func (s *KeeperTestSuite) TestCreatePoolZeroLiquidityNoCreationFee() {
// 	poolManagerModuleAcc := s.App.AccountKeeper.GetModuleAccount(s.Ctx, types.ModuleName)

// 	// withCreator := func(msg clmodel.MsgCreateConcentratedPool, address sdk.AccAddress) clmodel.MsgCreateConcentratedPool {
// 	// 	msg.Sender = address.String()
// 	// 	return msg
// 	// }

// 	balancerPoolMsg := balancer.NewMsgCreateBalancerPool(poolManagerModuleAcc.GetAddress(), balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()), []balancer.PoolAsset{
// 		{
// 			Token:  sdk.NewCoin(FOO, defaultInitPoolAmount),
// 			Weight: osmomath.NewInt(1),
// 		},
// 		{
// 			Token:  sdk.NewCoin(BAR, defaultInitPoolAmount),
// 			Weight: osmomath.NewInt(1),
// 		},
// 	})

// 	// concentratedPoolMsg := clmodel.NewMsgCreateConcentratedPool(poolManagerModuleAcc.GetAddress(), FOO, BAR, 1, defaultPoolSpreadFactor)

// 	tests := []struct {
// 		name               string
// 		msg                types.CreatePoolMsg
// 		expectedModuleType reflect.Type
// 		expectError        error
// 	}{
// 		// {
// 		// 	name:               "pool manager creator for concentrated pool - success",
// 		// 	msg:                concentratedPoolMsg,
// 		// 	expectedModuleType: concentratedKeeperType,
// 		// },
// 		// {
// 		// 	name:        "creator is not pool manager - failure",
// 		// 	msg:         withCreator(concentratedPoolMsg, s.TestAccs[0]),
// 		// 	expectError: types.InvalidPoolCreatorError{CreatorAddresss: s.TestAccs[0].String()},
// 		// },
// 		{
// 			name:        "balancer pool with pool manager creator - error, wrong pool",
// 			msg:         balancerPoolMsg,
// 			expectError: types.InvalidPoolTypeError{PoolType: types.Balancer},
// 		},
// 	}

// 	for i, tc := range tests {
// 		s.Run(tc.name, func() {
// 			poolmanagerKeeper := s.App.PoolManagerKeeper
// 			ctx := s.Ctx

// 			// Note: this is necessary for gauge creation in the after pool created hook.
// 			// There is a check requiring positive supply existing on-chain.
// 			s.MintCoins(sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, osmomath.OneInt())))

// 			pool, err := poolmanagerKeeper.CreateConcentratedPoolAsPoolManager(ctx, tc.msg)

// 			if tc.expectError != nil {
// 				s.Require().Error(err)
// 				s.Require().ErrorIs(err, tc.expectError)
// 				return
// 			}

// 			// Validate pool.
// 			s.Require().NoError(err)
// 			s.Require().Equal(uint64(i+1), pool.GetId())

// 			// Validate that mapping pool id -> module type has been persisted.
// 			swapModule, err := poolmanagerKeeper.GetPoolModule(ctx, pool.GetId())
// 			s.Require().NoError(err)
// 			s.Require().Equal(tc.expectedModuleType, reflect.TypeOf(swapModule))
// 		})
// 	}
// }

func (s *KeeperTestSuite) TestSetAndGetAllPoolRoutes() {
	tests := []struct {
		name         string
		preSetRoutes []types.ModuleRoute
	}{
		{
			name:         "no routes",
			preSetRoutes: []types.ModuleRoute{},
		},
		{
			name: "only balancer",
			preSetRoutes: []types.ModuleRoute{
				{
					PoolType: types.Balancer,
					PoolId:   1,
				},
			},
		},
		{
			name: "two balancer",
			preSetRoutes: []types.ModuleRoute{
				{
					PoolType: types.Balancer,
					PoolId:   1,
				},
				{
					PoolType: types.Balancer,
					PoolId:   2,
				},
			},
		},
		{
			name: "all supported pools",
			preSetRoutes: []types.ModuleRoute{
				{
					PoolType: types.Balancer,
					PoolId:   1,
				},
				{
					PoolType: types.Stableswap,
					PoolId:   2,
				},
				// {
				// 	PoolType: types.Concentrated,
				// 	PoolId:   3,
				// },
				// {
				// 	PoolType: types.CosmWasm,
				// 	PoolId:   4,
				// },
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Setup()
			poolManagerKeeper := s.App.PoolManagerKeeper

			for _, preSetRoute := range tc.preSetRoutes {
				poolManagerKeeper.SetPoolRoute(s.Ctx, preSetRoute.PoolId, preSetRoute.PoolType)
			}

			moduleRoutes := poolManagerKeeper.GetAllPoolRoutes(s.Ctx)

			// Validate.
			s.Require().Len(moduleRoutes, len(tc.preSetRoutes))
			s.Require().EqualValues(tc.preSetRoutes, moduleRoutes)
		})
	}
}

func (s *KeeperTestSuite) TestGetNextPoolIdAndIncrement() {
	tests := []struct {
		name               string
		expectedNextPoolId uint64
	}{
		{
			name:               "small next pool ID",
			expectedNextPoolId: 2,
		},
		{
			name:               "large next pool ID",
			expectedNextPoolId: 2999999,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Setup()

			s.App.PoolManagerKeeper.SetNextPoolId(s.Ctx, tc.expectedNextPoolId)
			nextPoolId := s.App.PoolManagerKeeper.GetNextPoolId(s.Ctx)
			s.Require().Equal(tc.expectedNextPoolId, nextPoolId)

			// System under test.
			nextPoolId = s.App.PoolManagerKeeper.GetNextPoolIdAndIncrement(s.Ctx)
			s.Require().Equal(tc.expectedNextPoolId, nextPoolId)
			s.Require().Equal(tc.expectedNextPoolId+1, s.App.PoolManagerKeeper.GetNextPoolId(s.Ctx))
		})
	}
}

func (s *KeeperTestSuite) TestValidateCreatedPool() {
	tests := []struct {
		name          string
		poolId        uint64
		pool          types.PoolI
		expectedError error
	}{
		{
			name:   "pool ID 1",
			poolId: 1,
			pool: &balancer.Pool{
				Address: types.NewPoolAddress(1).String(),
				Id:      1,
			},
		},
		{
			name:   "pool ID 309",
			poolId: 309,
			pool: &balancer.Pool{
				Address: types.NewPoolAddress(309).String(),
				Id:      309,
			},
		},
		{
			name:   "error: unexpected ID",
			poolId: 1,
			pool: &balancer.Pool{
				Address: types.NewPoolAddress(1).String(),
				Id:      2,
			},
			expectedError: types.IncorrectPoolIdError{ExpectedPoolId: 1, ActualPoolId: 2},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.Setup()

			// System under test.
			err := s.App.PoolManagerKeeper.ValidateCreatedPool(s.Ctx, tc.poolId, tc.pool)
			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expectedError.Error())
				return
			}
			s.Require().NoError(err)
		})
	}
}

// TestCreatePoolWithBadgesDisablePoolCreationInvariant tests that pool creation fails when
// a badges collection has disablePoolCreation set to true.
func (s *KeeperTestSuite) TestCreatePoolWithBadgesDisablePoolCreationInvariant() {
	s.Setup()

	ctx := s.Ctx
	badgesKeeper := s.App.BadgesKeeper
	poolmanagerKeeper := s.App.PoolManagerKeeper
	gammKeeper := s.App.GammKeeper

	// Create a badges collection with allowPoolCreation set to false
	creator := s.TestAccs[0]
	creatorStr := creator.String()

	// Create collection with wrapper path and allowPoolCreation = false
	createCollectionMsg := &badgestypes.MsgCreateCollection{
		Creator: creatorStr,
		ValidTokenIds: []*badgestypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		CollectionPermissions: &badgestypes.CollectionPermissions{},
		Manager: creatorStr,
		CosmosCoinWrapperPathsToAdd: []*badgestypes.CosmosCoinWrapperPathAddObject{
			{
				Denom: "testbadge",
				Balances: []*badgestypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
						TokenIds:       []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					},
				},
				Symbol:              "TEST",
				DenomUnits:          []*badgestypes.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "testbadge", IsDefaultDisplay: true}},
				AllowCosmosWrapping: true,
			},
		},
		Invariants: &badgestypes.InvariantsAddObject{
			DisablePoolCreation: true, // Set to true to block pool creation
		},
	}

	// Create the collection
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(badgesKeeper)
	createCollectionResp, err := badgesMsgServer.CreateCollection(ctx, createCollectionMsg)
	s.Require().NoError(err)
	collectionId := createCollectionResp.CollectionId

	// Get the collection to find the wrapper denom
	collection, found := badgesKeeper.GetCollectionFromStore(ctx, collectionId)
	s.Require().True(found)
	s.Require().NotNil(collection.Invariants)
	s.Require().True(collection.Invariants.DisablePoolCreation, "disablePoolCreation should be true")

	// Get the wrapper denom
	wrapperPath := collection.CosmosCoinWrapperPaths[0]
	badgesDenom := "badges:" + collectionId.String() + ":" + wrapperPath.Denom

	// Try to create a pool with the badges asset that has allowPoolCreation = false
	// This should fail
	poolMsg := balancer.NewMsgCreateBalancerPool(
		creator,
		balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()),
		[]balancer.PoolAsset{
			{
				Token:  sdk.NewCoin(badgesDenom, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
			{
				Token:  sdk.NewCoin(FOO, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
		},
	)

	// Fund the creator
	s.FundAcc(creator, sdk.NewCoins(
		sdk.NewCoin(FOO, defaultInitPoolAmount.Mul(osmomath.NewInt(2))),
	))

	// Try to create the pool - this should fail
	_, err = poolmanagerKeeper.CreatePool(ctx, poolMsg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "disablePoolCreation invariant is set to true")

	// Now test that pool creation validation passes when disablePoolCreation is false
	// Create another collection with disablePoolCreation = false
	createCollectionMsg2 := &badgestypes.MsgCreateCollection{
		Creator: creatorStr,
		ValidTokenIds: []*badgestypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		CollectionPermissions: &badgestypes.CollectionPermissions{},
		Manager: creatorStr,
		CosmosCoinWrapperPathsToAdd: []*badgestypes.CosmosCoinWrapperPathAddObject{
			{
				Denom: "testbadgeallowed",
				Balances: []*badgestypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
						TokenIds:       []*badgestypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					},
				},
				Symbol:              "TESTALLOWED",
				DenomUnits:          []*badgestypes.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "testbadgeallowed", IsDefaultDisplay: true}},
				AllowCosmosWrapping: true,
			},
		},
		Invariants: &badgestypes.InvariantsAddObject{
			DisablePoolCreation: false, // Set to false to allow pool creation
		},
	}

	createCollectionResp2, err := badgesMsgServer.CreateCollection(ctx, createCollectionMsg2)
	s.Require().NoError(err)
	collectionId2 := createCollectionResp2.CollectionId

	collection2, found := badgesKeeper.GetCollectionFromStore(ctx, collectionId2)
	s.Require().True(found)
	s.Require().NotNil(collection2.Invariants)
	s.Require().False(collection2.Invariants.DisablePoolCreation, "disablePoolCreation should be false")

	wrapperPath2 := collection2.CosmosCoinWrapperPaths[0]
	badgesDenom2 := "badges:" + collectionId2.String() + ":" + wrapperPath2.Denom

	// Try to create a pool with the badges asset that has disablePoolCreation = false
	// The validation should pass (disablePoolCreation is false)
	poolMsg2 := balancer.NewMsgCreateBalancerPool(
		creator,
		balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()),
		[]balancer.PoolAsset{
			{
				Token:  sdk.NewCoin(badgesDenom2, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
			{
				Token:  sdk.NewCoin(FOO, defaultInitPoolAmount),
				Weight: osmomath.NewInt(1),
			},
		},
	)

	// The validation should pass (disablePoolCreation is false)
	// Note: The actual pool creation might still fail due to insufficient badges tokens,
	// but the disablePoolCreation check should pass
	err = gammKeeper.ValidatePoolCreationAllowed(ctx, poolMsg2.InitialLiquidity())
	s.Require().NoError(err, "Pool creation should be allowed when disablePoolCreation is false")
}
