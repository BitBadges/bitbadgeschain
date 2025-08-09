package eip712

/*
	These are used as fully populated examples to generate EIP712 types.
	This is because the EIP712 type generation code expects all values to be populated an  non-optional.

	We want to make sure the type generation includes all default values and empty values, even for optional fields.
	This is because that is what the SDK does.
*/

// getMerkleChallengeSchema returns the schema for a Merkle challenge
func getMerkleChallengeSchema() string {
	return `{
		"root": "",
		"expectedProofLength": "",
		"useCreatorAddressAsLeaf": false,
		"leafSigner": "",
		"maxUsesPerLeaf": "",
		"challengeTrackerId": "",
		"uri": "",
		"customData": ""
	}`
}

// getETHSignatureChallengeSchema returns the schema for an ETH signature challenge
func getETHSignatureChallengeSchema() string {
	return `{
		"signer": "",
		"challengeTrackerId": "",
		"uri": "",
		"customData": ""
	}`
}

// getMustOwnBadgesSchema returns the schema for must own tokens criteria
func getMustOwnBadgesSchema() string {
	return `{
		"collectionId": "",
		"amountRange": ` + getUintRangeSchema() + `,
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"badgeIds": [` + getUintRangeSchema() + `],
		"overrideWithCurrentTime": false,
		"mustSatisfyForAllAssets": false,
		"ownershipCheckParty": ""
	}`
}

// getCoinTransferSchema returns the schema for coin transfers
func getCoinTransferSchema() string {
	return `{
		"to": "",
		"overrideFromWithApproverAddress": false,
		"overrideToWithInitiator": false,
		"coins": [
			{
				"amount": "",
				"denom": ""
			}
		]
	}`
}

// getBalanceSchema returns the schema for a balance
func getBalanceSchema() string {
	return `{
		"amount": "",
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"badgeIds": [` + getUintRangeSchema() + `]
	}`
}

// getPredeterminedBalancesSchema returns the schema for predetermined balances
func getPredeterminedBalancesSchema() string {
	return `{
		"manualBalances": [
			{
				"balances": [
					` + getBalanceSchema() + `
				]
			}
		],
		"incrementedBalances": {
			"startBalances": [
				` + getBalanceSchema() + `
			],
			"incrementBadgeIdsBy": "",
			"allowOverrideWithAnyValidBadge": false,
			"durationFromTimestamp": "",
			"incrementOwnershipTimesBy": "",
			"allowOverrideTimestamp": true,
			"recurringOwnershipTimes": { 
				"startTime": "",
				"intervalLength": "",
				"chargePeriodLength": ""
			}
		},
		"orderCalculationMethod": {
			"useOverallNumTransfers": false,
			"usePerToAddressNumTransfers": false,
			"usePerFromAddressNumTransfers": false,
			"usePerInitiatedByAddressNumTransfers": false,
			"challengeTrackerId": "",
			"useMerkleChallengeLeafIndex": false
		}
	}`
}

// getApprovalAmountsSchema returns the schema for approval amounts
func getApprovalAmountsSchema() string {
	return `{
		"overallApprovalAmount": "",
		"perToAddressApprovalAmount": "",
		"perFromAddressApprovalAmount": "",
		"amountTrackerId": "",
		"perInitiatedByAddressApprovalAmount": "",
		"resetTimeIntervals": {
			"startTime": "",
			"intervalLength": ""
		}
	}`
}

// getMaxNumTransfersSchema returns the schema for max number of transfers
func getMaxNumTransfersSchema() string {
	return `{
		"overallMaxNumTransfers": "",
		"perToAddressMaxNumTransfers": "",
		"perFromAddressMaxNumTransfers": "",
		"amountTrackerId": "",
		"perInitiatedByAddressMaxNumTransfers": "",
		"resetTimeIntervals": {
			"startTime": "",
			"intervalLength": ""
		}
	}`
}

// getAutoDeletionOptionsSchema returns the schema for auto deletion options
func getAutoDeletionOptionsSchema() string {
	return `{ 
		"afterOneUse": true, 
		"afterOverallMaxNumTransfers": false, 
		"allowCounterpartyPurge": false, 
		"allowPurgeIfExpired": false 
	}`
}

// getDynamicStoreChallengeSchema returns the schema for dynamic store challenges
func getDynamicStoreChallengeSchema() string {
	return `{ "storeId": "" }`
}

// getBaseApprovalCriteriaSchema returns the base approval criteria schema without collection-specific fields
func getBaseApprovalCriteriaSchema() string {
	return `{
		"mustOwnBadges": [
			` + getMustOwnBadgesSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getCollectionApprovalCriteriaSchema returns the schema for collection approval criteria
func getCollectionApprovalCriteriaSchema() string {
	return `{
		"mustOwnBadges": [
			` + getMustOwnBadgesSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireToEqualsInitiatedBy": false,
		"requireFromEqualsInitiatedBy": false,
		"requireToDoesNotEqualInitiatedBy": false,
		"requireFromDoesNotEqualInitiatedBy": false,
		"overridesFromOutgoingApprovals": false,
		"userRoyalties": {
			"percentage": "",
			"payoutAddress": ""
		},
		"overridesToIncomingApprovals": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getOutgoingApprovalCriteriaSchema returns the schema for outgoing approval criteria
func getOutgoingApprovalCriteriaSchema() string {
	return `{
		"mustOwnBadges": [
			` + getMustOwnBadgesSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireToEqualsInitiatedBy": false,
		"requireToDoesNotEqualInitiatedBy": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getIncomingApprovalCriteriaSchema returns the schema for incoming approval criteria
func getIncomingApprovalCriteriaSchema() string {
	return `{
		"mustOwnBadges": [
			` + getMustOwnBadgesSchema() + `
		],
		"merkleChallenges": [
			` + getMerkleChallengeSchema() + `
		],
		"coinTransfers": [
			` + getCoinTransferSchema() + `
		],
		"predeterminedBalances": ` + getPredeterminedBalancesSchema() + `,
		"approvalAmounts": ` + getApprovalAmountsSchema() + `,
		"autoDeletionOptions": ` + getAutoDeletionOptionsSchema() + `,
		"maxNumTransfers": ` + getMaxNumTransfersSchema() + `,
		"requireFromEqualsInitiatedBy": false,
		"requireFromDoesNotEqualInitiatedBy": false,
		"dynamicStoreChallenges": [
			` + getDynamicStoreChallengeSchema() + `
		],
		"ethSignatureChallenges": [
			` + getETHSignatureChallengeSchema() + `
		]
	}`
}

// getUintRangeSchema returns the schema for a uint range
func getUintRangeSchema() string {
	return `{"start": "", "end": ""}`
}

// getCollectionApprovalSchema returns the schema for collection approval
func getCollectionApprovalSchema() string {
	return `{
		"fromListId": "",
		"toListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"badgeIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getCollectionApprovalCriteriaSchema() + `
	}`
}

// getIncomingApprovalSchema returns the schema for an incoming approval object
func getIncomingApprovalSchema() string {
	return `{
		"fromListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"badgeIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getIncomingApprovalCriteriaSchema() + `
	}`
}

// getOutgoingApprovalSchema returns the schema for an outgoing approval object
func getOutgoingApprovalSchema() string {
	return `{
		"toListId": "",
		"initiatedByListId": "",
		"transferTimes": [` + getUintRangeSchema() + `],
		"badgeIds": [` + getUintRangeSchema() + `],
		"ownershipTimes": [` + getUintRangeSchema() + `],
		"uri": "",
		"customData": "",
		"approvalId": "",
		"version": "0",
		"approvalCriteria": ` + getOutgoingApprovalCriteriaSchema() + `
	}`
}

// getCollectionPermissionsSchema returns the schema for collection permissions
func getCollectionPermissionsSchema() string {
	return `{
		"canDeleteCollection": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canArchiveCollection": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateOffChainBalancesMetadata": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateStandards": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCustomData": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateManager": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCollectionMetadata": [{"timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateValidBadgeIds": [{"badgeIds": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateBadgeMetadata": [{"badgeIds": [` + getUintRangeSchema() + `], "timelineTimes": [` + getUintRangeSchema() + `], "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateCollectionApprovals": [{"fromListId": "", "toListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "badgeIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}]
	}`
}

// getUserPermissionsSchema returns the schema for user permissions
func getUserPermissionsSchema() string {
	return `{
		"canUpdateOutgoingApprovals": [{"toListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "badgeIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateIncomingApprovals": [{"fromListId": "", "initiatedByListId": "", "transferTimes": [` + getUintRangeSchema() + `], "badgeIds": [` + getUintRangeSchema() + `], "ownershipTimes": [` + getUintRangeSchema() + `], "approvalId": "", "permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}],
		"canUpdateAutoApproveAllIncomingTransfers": [{"permanentlyPermittedTimes": [` + getUintRangeSchema() + `], "permanentlyForbiddenTimes": [` + getUintRangeSchema() + `]}]
	}`
}

// GetSchemas returns all the schemas for the EIP712 types
func GetSchemas() []string {
	schemas := make([]string, 0)

	schemas = append(schemas, `{
		"type": "maps/CreateMap",
		"value": {
			"creator": "",
			"mapId": "",
					"inheritManagerTimelineFrom": "",
		"managerTimeline": [
			{
				"manager": "",
				"timelineTimes": [`+getUintRangeSchema()+`]
			}
		],
		"updateCriteria": {
			"managerOnly": false,
			"collectionId": "",
			"creatorOnly": false,
			"firstComeFirstServe": false
		},
		"valueOptions": {
			"noDuplicates": false,
			"permanentOnceSet": false,
			"expectUint": false,
			"expectBoolean": false,
			"expectAddress": false,
			"expectUri": false
		},
		"defaultValue": "",
		"metadataTimeline": [
			{
				"metadata": {
					"uri": "",
					"customData": ""
				},
				"timelineTimes": [`+getUintRangeSchema()+`]
			}
		],
		"permissions": {
			"canUpdateMetadata": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateManager": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canDeleteMap": [
				{
					"timelineTimes": [`+getUintRangeSchema()+`],
					"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
					"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/UpdateMap",
		"value": {
			"creator": "",
			"mapId": "",
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateMetadataTimeline": false,
			"metadataTimeline": [
				{
					"metadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updatePermissions": false,
			"permissions": {
				"canUpdateMetadata": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				],
				"canDeleteMap": [
					{
						"timelineTimes": [`+getUintRangeSchema()+`],
						"permanentlyPermittedTimes": [`+getUintRangeSchema()+`],
						"permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]
					}
				]
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/DeleteMap",
		"value": {
			"creator": "",
			"mapId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "maps/SetValue",
		"value": {
			"creator": "",
			"mapId": "",
			"key": "",
			"value": "",
			"options": {
				"useMostRecentCollectionId": false
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/CreateProtocol",
		"value": {
			"creator": "",
			"name": "",
			"uri": "",
			"customData": "",
			"isFrozen": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/DeleteProtocol",
		"value": {
			"creator": "",
			"name": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/SetCollectionForProtocol",
		"value": {
			"creator": "",
			"name": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "protocols/UpdateProtocol",
		"value": {
			"creator": "",
			"name": "",
			"uri": "",
			"customData": "",
			"isFrozen": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateAddressLists",
		"value": {
			"creator": "",
			"addressLists": [
				{
					"listId": "",
					"addresses": [],
					"whitelist": false,
					"uri": "",
					"customData": "",
					"createdBy": ""
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteCollection",
		"value": {
			"creator": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/TransferBadges",
		"value": {
			"creator": "",
			"collectionId": "",
			"transfers": [
				{
					"from": "",
					"toAddresses": [],
					"balances": [
						`+getBalanceSchema()+`
					],
					"precalculateBalancesFromApproval": {
						"approvalId": "",
						"approvalLevel": "",
						"approverAddress": "",
						"version": "0"
					},
					"merkleProofs": [
						{
							"leaf": "",
							"aunts": [
								{
									"aunt": "",
									"onRight": false
								}
							],
							"leafSignature": ""
						}
					],
					"ethSignatureProofs": [
						{
							"nonce": "",
							"signature": ""
						}
					],
					"memo": "",
					"prioritizedApprovals": [
						{
							"approvalId": "",
							"approvalLevel": "",
							"approverAddress": "",
							"version": "0"
						}
					],
					"onlyCheckPrioritizedCollectionApprovals": false,
					"onlyCheckPrioritizedIncomingApprovals": false,
					"onlyCheckPrioritizedOutgoingApprovals": false,
					"precalculationOptions": {
						"overrideTimestamp": "0",
						"badgeIdsOverride": [`+getUintRangeSchema()+`]
					},
					"affiliateAddress": "",
					"numAttempts": "0"
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UniversalUpdateCollection",
		"value": {
			"creator": "",
			"collectionId": "",
			"balancesType": "",
			"defaultBalances": {
				"balances": [
					`+getBalanceSchema()+`
				],
				"incomingApprovals": [
					`+getIncomingApprovalSchema()+`
				],
				"outgoingApprovals": [
					`+getOutgoingApprovalSchema()+`
				],
				"userPermissions": `+getUserPermissionsSchema()+`,
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true,
				"autoApproveAllIncomingTransfers": true
			},
			"updateValidBadgeIds": false,
			"validBadgeIds": [`+getUintRangeSchema()+`],
			"updateCollectionPermissions": false,
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateBadgeMetadataTimeline": false,
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateOffChainBalancesMetadataTimeline": false,
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					]
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateUserApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"updateOutgoingApprovals": false,
			"outgoingApprovals": [
				`+getOutgoingApprovalSchema()+`
			],
			"updateIncomingApprovals": false,
			"incomingApprovals": [
				`+getIncomingApprovalSchema()+`
			],
			"updateAutoApproveSelfInitiatedOutgoingTransfers": false,
			"autoApproveSelfInitiatedOutgoingTransfers": false,
			"updateAutoApproveSelfInitiatedIncomingTransfers": false,
			"autoApproveSelfInitiatedIncomingTransfers": false,
			"updateAutoApproveAllIncomingTransfers": false,
			"autoApproveAllIncomingTransfers": false,
			"updateUserPermissions": false,
			"userPermissions": `+getUserPermissionsSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateCollection",
		"value": {
			"creator": "",
			"balancesType": "",
			"defaultBalances": {
										"balances": [
					`+getBalanceSchema()+`
				],
				"incomingApprovals": [
					`+getIncomingApprovalSchema()+`
				],
				"outgoingApprovals": [
					`+getOutgoingApprovalSchema()+`
				],
				"userPermissions": `+getUserPermissionsSchema()+`,
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true,
				"autoApproveAllIncomingTransfers": true
			},
			"validBadgeIds": [`+getUintRangeSchema()+`],
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					]
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateCollection",
		"value": {
			"creator": "",
			"collectionId": "",
			"updateValidBadgeIds": false,
			"validBadgeIds": [`+getUintRangeSchema()+`],
			"updateCollectionPermissions": false,
			"collectionPermissions": `+getCollectionPermissionsSchema()+`,
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateBadgeMetadataTimeline": false,
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [`+getUintRangeSchema()+`]
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateOffChainBalancesMetadataTimeline": false,
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"mintEscrowCoinsToTransfer": [
				{
					"amount": "",
					"denom": ""
				}
			],
			"cosmosCoinWrapperPathsToAdd": [
				{
					"denom": "",
					"balances": [
						`+getBalanceSchema()+`
					],
					"symbol": "",
					"denomUnits": [
						{
							"decimals": "0",
							"symbol": "",
							"isDefaultDisplay": false
						}
					]
				}
			],
			"invariants": {
				"noCustomOwnershipTimes": false
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateDynamicStore",
		"value": {
			"creator": "",
			"defaultValue": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateDynamicStore",
		"value": {
			"creator": "",
			"storeId": "",
			"defaultValue": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteDynamicStore",
		"value": {
			"creator": "",
			"storeId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetDynamicStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"value": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/IncrementStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"amount": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DecrementStoreValue",
		"value": {
			"creator": "",
			"storeId": "",
			"address": "",
			"amount": "",
			"setToZeroOnUnderflow": false
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetIncomingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approval": `+getIncomingApprovalSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteIncomingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approvalId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetOutgoingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approval": `+getOutgoingApprovalSchema()+`
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteOutgoingApproval",
		"value": {
			"creator": "",
			"collectionId": "",
			"approvalId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/PurgeApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"purgeExpired": false,
			"approverAddress": "",
			"purgeCounterpartyApprovals": false,
			"approvalsToPurge": [
				{
					"approvalId": "",
					"approvalLevel": "",
					"approverAddress": "",
					"version": "0"
				}
			]
		}
	}`)

	// UniversalUpdateCollection helper message types
	schemas = append(schemas, `{
		"type": "badges/SetValidBadgeIds",
		"value": {
			"creator": "",
			"collectionId": "",
			"validBadgeIds": [`+getUintRangeSchema()+`],
			"canUpdateValidBadgeIds": [{"badgeIds": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetManager",
		"value": {
			"creator": "",
			"collectionId": "",
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateManager": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCollectionMetadata",
		"value": {
			"creator": "",
			"collectionId": "",
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateCollectionMetadata": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetBadgeMetadata",
		"value": {
			"creator": "",
			"collectionId": "",
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": ""
						}
					],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateBadgeMetadata": [{"badgeIds": [`+getUintRangeSchema()+`], "timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCustomData",
		"value": {
			"creator": "",
			"collectionId": "",
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateCustomData": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetStandards",
		"value": {
			"creator": "",
			"collectionId": "",
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canUpdateStandards": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetCollectionApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"collectionApprovals": [
				`+getCollectionApprovalSchema()+`
			],
			"canUpdateCollectionApprovals": [{"fromListId": "", "toListId": "", "initiatedByListId": "", "transferTimes": [`+getUintRangeSchema()+`], "badgeIds": [`+getUintRangeSchema()+`], "ownershipTimes": [`+getUintRangeSchema()+`], "approvalId": "", "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/SetIsArchived",
		"value": {
			"creator": "",
			"collectionId": "",
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [`+getUintRangeSchema()+`]
				}
			],
			"canArchiveCollection": [{"timelineTimes": [`+getUintRangeSchema()+`], "permanentlyPermittedTimes": [`+getUintRangeSchema()+`], "permanentlyForbiddenTimes": [`+getUintRangeSchema()+`]}]
		}
	}`)

	return schemas
}
