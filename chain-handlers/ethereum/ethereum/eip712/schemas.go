package eip712

/*
	These are used as fully populated examples to generate EIP712 types.
	This is because the EIP712 type generation code expects all values to be populated an  non-optional.

	We want to make sure the type generation includes all default values and empty values, even for optional fields.
	This is because that is what the SDK does.
*/

//TODO: Store JSONs in a file directory not directly here

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
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
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
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"permissions": {
				"canUpdateMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canDeleteMap": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
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
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateMetadataTimeline": false,
			"metadataTimeline": [
				{
					"metadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updatePermissions": false,
			"permissions": {
				"canUpdateMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canDeleteMap": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
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
				"useMostRecentCollectionId": false,
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
			"creatorOverride": "",
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
			"creatorOverride": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/TransferBadges",
		"value": {
			"creator": "",
			"creatorOverride": "",
			"collectionId": "",
			"transfers": [
				{
					"from": "",
					"toAddresses": [],
					"balances": [
						{
							"amount": "",
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"precalculateBalancesFromApproval": {
						"approvalId": "",
						"approvalLevel": "",
						"approverAddress": ""
					},
					"merkleProofs": [
						{
							"leaf": "",
							"aunts": [
								{
									"aunt": "",
									"onRight": false
								}
							]
						}
					],
					"memo": "",
					"prioritizedApprovals": [
						{
							"approvalId": "",
							"approvalLevel": "",
							"approverAddress": ""
						}
					],
					"onlyCheckPrioritizedCollectionApprovals": false,
					"onlyCheckPrioritizedIncomingApprovals": false,
					"onlyCheckPrioritizedOutgoingApprovals": false
					]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UniversalUpdateCollection",
		"value": {
			"creator": "",
			"creatorOverride": "",
			"collectionId": "",
			"balancesType": "",
			"defaultBalances": {
				"balances": [	{
					"amount": "",
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					]
				}],
				"incomingApprovals":  [
					{
						"fromListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"uri": "",
						"customData": "",
						"approvalId": "",
						"approvalCriteria": {
							"merkleChallenges": [{
								"root": "",
								"expectedProofLength": "",
								"useCreatorAddressAsLeaf": false,
								"maxUsesPerLeaf": "",
								"challengeTrackerId": "",
								"uri": "",
								"customData": ""
							}],
							"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
							"predeterminedBalances": {
								"manualBalances": [
									{
										"balances": [
											{
												"amount": "",
												"ownershipTimes": [
													{
														"start": "",
														"end": ""
													}
												],
												"badgeIds": [
													{
														"start": "",
														"end": ""
													}
												]
											}
										]
									}
								],
								"incrementedBalances": {
									"startBalances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									],
									"incrementBadgeIdsBy": "",
									"incrementOwnershipTimesBy": ""
								},
								"orderCalculationMethod": {
									"useOverallNumTransfers": false,
									"usePerToAddressNumTransfers": false,
									"usePerFromAddressNumTransfers": false,
									"usePerInitiatedByAddressNumTransfers": false,
									"challengeTrackerId": "",
									"useMerkleChallengeLeafIndex": false
								}
							},
							"approvalAmounts": {
								"overallApprovalAmount": "",
								"perToAddressApprovalAmount": "",
								"perFromAddressApprovalAmount": "",
								"amountTrackerId": "",
								"perInitiatedByAddressApprovalAmount": ""
							},
							"maxNumTransfers": {
								"overallMaxNumTransfers": "",
								"perToAddressMaxNumTransfers": "",
								"perFromAddressMaxNumTransfers": "",
								"amountTrackerId": "",
								"perInitiatedByAddressMaxNumTransfers": ""
							},
							"requireFromEqualsInitiatedBy": false,
							"requireFromDoesNotEqualInitiatedBy": false
						}
					}
				],
				"outgoingApprovals": [
					{
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"uri": "",
						"customData": "",
						"approvalId": "",
						"approvalCriteria": {
							"merkleChallenges": [{
								"root": "",
								"expectedProofLength": "",
								"useCreatorAddressAsLeaf": false,
								"maxUsesPerLeaf": "",
								"challengeTrackerId": "",
								"uri": "",
								"customData": ""
							}],
							"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
							"predeterminedBalances": {
								"manualBalances": [
									{
										"balances": [
											{
												"amount": "",
												"ownershipTimes": [
													{
														"start": "",
														"end": ""
													}
												],
												"badgeIds": [
													{
														"start": "",
														"end": ""
													}
												]
											}
										]
									}
								],
								"incrementedBalances": {
									"startBalances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									],
									"incrementBadgeIdsBy": "",
									"incrementOwnershipTimesBy": ""
								},
								"orderCalculationMethod": {
									"useOverallNumTransfers": false,
									"usePerToAddressNumTransfers": false,
									"usePerFromAddressNumTransfers": false,
									"usePerInitiatedByAddressNumTransfers": false,
									"challengeTrackerId": "",
									"useMerkleChallengeLeafIndex": false
								}
							},
							"approvalAmounts": {
								"overallApprovalAmount": "",
								"perToAddressApprovalAmount": "",
								"perFromAddressApprovalAmount": "",
								"amountTrackerId": "",
								"perInitiatedByAddressApprovalAmount": ""
							},
							"maxNumTransfers": {
								"overallMaxNumTransfers": "",
								"perToAddressMaxNumTransfers": "",
								"perFromAddressMaxNumTransfers": "",
								"amountTrackerId": "",
								"perInitiatedByAddressMaxNumTransfers": ""
							},
							"requireToEqualsInitiatedBy": false,
							"requireToDoesNotEqualInitiatedBy": false
						}
					}
				],
				"userPermissions": {
					"canUpdateOutgoingApprovals": [
						{
							"toListId": "",
							"initiatedByListId": "",
							"transferTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							],
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"approvalId": "",
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateIncomingApprovals": [
						{
							"fromListId": "",
							"initiatedByListId": "",
							"transferTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							],
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"approvalId": "",
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [
						{
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [
						{
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					]
				},
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true
			},
			
			"badgeIdsToAdd": [
				{
					"start": "",
					"end": ""
				}
			],
			"updateCollectionPermissions": false,
			"collectionPermissions": {
				"canDeleteCollection": [
					{
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canArchiveCollection": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateOffChainBalancesMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateStandards": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCustomData": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateValidBadgeIds": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateBadgeMetadata": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionApprovals": [
					{
						"fromListId": "",
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"approvalId": "",
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateBadgeMetadataTimeline": false,
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateOffChainBalancesMetadataTimeline": false,
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				{
					"fromListId": "",
					"toListId": "",
					"initiatedByListId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"merkleChallenges": [{
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"challengeTrackerId": "",
							"uri": "",
							"customData": ""
						}],
						"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"challengeTrackerId": "",
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"amountTrackerId": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"amountTrackerId": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireFromEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false,
						"overridesFromOutgoingApprovals": false,
						"overridesToIncomingApprovals": false
					}
				}
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateUserApprovals",
		"value": {
			"creator": "",
			"creatorOverride": "",
			"collectionId": "",
			"updateOutgoingApprovals": false,
			"outgoingApprovals": [
				{
					"toListId": "",
					"initiatedByListId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"merkleChallenges": [{
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"challengeTrackerId": "",
							"uri": "",
							"customData": ""
						}],
						"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"challengeTrackerId": "",
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"amountTrackerId": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"amountTrackerId": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false
					}
				}
			],
			"updateIncomingApprovals": false,
			"incomingApprovals": [
				{
					"fromListId": "",
					"initiatedByListId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"merkleChallenges": [{
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"challengeTrackerId": "",
							"uri": "",
							"customData": ""
						}],
						"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"challengeTrackerId": "",
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"amountTrackerId": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"amountTrackerId": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireFromEqualsInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false
					}
				}
			],
			"updateAutoApproveSelfInitiatedOutgoingTransfers": false,
			"autoApproveSelfInitiatedOutgoingTransfers": false,
			"updateAutoApproveSelfInitiatedIncomingTransfers": false,
			"autoApproveSelfInitiatedIncomingTransfers": false,
			"updateUserPermissions": false,
			"userPermissions": {
				"canUpdateOutgoingApprovals": [
					{
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"approvalId": "",
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateIncomingApprovals": [
					{
						"fromListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"approvalId": "",
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [
					{
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [
					{
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			}
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/CreateCollection",
		"value": {
			"creator": "",
			"creatorOverride": "",
			"balancesType": "",
			"defaultBalances": {
				"balances": [	{
					"amount": "",
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					]
				}],
				"incomingApprovals":  [
					{
						"fromListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"uri": "",
						"customData": "",
						"approvalId": "",
						"approvalCriteria": {
							"merkleChallenges": [{
								"root": "",
								"expectedProofLength": "",
								"useCreatorAddressAsLeaf": false,
								"maxUsesPerLeaf": "",
								"challengeTrackerId": "",
								"uri": "",
								"customData": ""
							}],
							"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
							"predeterminedBalances": {
								"manualBalances": [
									{
										"balances": [
											{
												"amount": "",
												"ownershipTimes": [
													{
														"start": "",
														"end": ""
													}
												],
												"badgeIds": [
													{
														"start": "",
														"end": ""
													}
												]
											}
										]
									}
								],
								"incrementedBalances": {
									"startBalances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									],
									"incrementBadgeIdsBy": "",
									"incrementOwnershipTimesBy": ""
								},
								"orderCalculationMethod": {
									"useOverallNumTransfers": false,
									"usePerToAddressNumTransfers": false,
									"usePerFromAddressNumTransfers": false,
									"usePerInitiatedByAddressNumTransfers": false,
									"challengeTrackerId": "",
									"useMerkleChallengeLeafIndex": false
								}
							},
							"approvalAmounts": {
								"overallApprovalAmount": "",
								"perToAddressApprovalAmount": "",
								"perFromAddressApprovalAmount": "",
								"amountTrackerId": "",
								"perInitiatedByAddressApprovalAmount": ""
							},
							"maxNumTransfers": {
								"overallMaxNumTransfers": "",
								"perToAddressMaxNumTransfers": "",
								"perFromAddressMaxNumTransfers": "",
								"amountTrackerId": "",
								"perInitiatedByAddressMaxNumTransfers": ""
							},
							"requireFromEqualsInitiatedBy": false,
							"requireFromDoesNotEqualInitiatedBy": false
						}
					}
				],
				"outgoingApprovals": [
					{
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"uri": "",
						"customData": "",
						"approvalId": "",
						"approvalCriteria": {
							"merkleChallenges": [{
								"root": "",
								"expectedProofLength": "",
								"useCreatorAddressAsLeaf": false,
								"maxUsesPerLeaf": "",
								"challengeTrackerId": "",
								"uri": "",
								"customData": ""
							}],
							"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
							"predeterminedBalances": {
								"manualBalances": [
									{
										"balances": [
											{
												"amount": "",
												"ownershipTimes": [
													{
														"start": "",
														"end": ""
													}
												],
												"badgeIds": [
													{
														"start": "",
														"end": ""
													}
												]
											}
										]
									}
								],
								"incrementedBalances": {
									"startBalances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									],
									"incrementBadgeIdsBy": "",
									"incrementOwnershipTimesBy": ""
								},
								"orderCalculationMethod": {
									"useOverallNumTransfers": false,
									"usePerToAddressNumTransfers": false,
									"usePerFromAddressNumTransfers": false,
									"usePerInitiatedByAddressNumTransfers": false,
									"challengeTrackerId": "",
									"useMerkleChallengeLeafIndex": false
								}
							},
							"approvalAmounts": {
								"overallApprovalAmount": "",
								"perToAddressApprovalAmount": "",
								"perFromAddressApprovalAmount": "",
								"amountTrackerId": "",
								"perInitiatedByAddressApprovalAmount": ""
							},
							"maxNumTransfers": {
								"overallMaxNumTransfers": "",
								"perToAddressMaxNumTransfers": "",
								"perFromAddressMaxNumTransfers": "",
								"amountTrackerId": "",
								"perInitiatedByAddressMaxNumTransfers": ""
							},
							"requireToEqualsInitiatedBy": false,
							"requireToDoesNotEqualInitiatedBy": false
						}
					}
				],
				"userPermissions": {
					"canUpdateOutgoingApprovals": [
						{
							"toListId": "",
							"initiatedByListId": "",
							"transferTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							],
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
						  "approvalId": "",
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateIncomingApprovals": [
						{
							"fromListId": "",
							"initiatedByListId": "",
							"transferTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							],
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"approvalId": "",
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [
						{
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [
						{
							"permanentlyPermittedTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"permanentlyForbiddenTimes": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					]
				},
				"autoApproveSelfInitiatedIncomingTransfers": true,
				"autoApproveSelfInitiatedOutgoingTransfers": true
			},
			"badgeIdsToAdd": [
				{
					"start": "",
					"end": ""
				}
			],
			"collectionPermissions": {
				"canDeleteCollection": [
					{
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canArchiveCollection": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateOffChainBalancesMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateStandards": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCustomData": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateValidBadgeIds": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateBadgeMetadata": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionApprovals": [
					{
						"fromListId": "",
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"approvalId": "",
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"collectionApprovals": [
				{
					"fromListId": "",
					"toListId": "",
					"initiatedByListId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"merkleChallenges": [{
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"challengeTrackerId": "",
							"uri": "",
							"customData": ""
						}],
						"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"challengeTrackerId": "",
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"amountTrackerId": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"amountTrackerId": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireFromEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false,
						"overridesFromOutgoingApprovals": false,
						"overridesToIncomingApprovals": false
					}
				}
			],
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateCollection",
		"value": {
			"creator": "",
			"creatorOverride": "",
			"collectionId": "",
			"badgeIdsToAdd": [
				{
					"start": "",
					"end": ""
				}
			],
			"updateCollectionPermissions": false,
			"collectionPermissions": {
				"canDeleteCollection": [
					{
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canArchiveCollection": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateOffChainBalancesMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateStandards": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCustomData": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateValidBadgeIds": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateBadgeMetadata": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionApprovals": [
					{
						"fromListId": "",
						"toListId": "",
						"initiatedByListId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"approvalId": "",
						"permanentlyPermittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permanentlyForbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateBadgeMetadataTimeline": false,
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateOffChainBalancesMetadataTimeline": false,
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				{
					"fromListId": "",
					"toListId": "",
					"initiatedByListId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"merkleChallenges": [{
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"challengeTrackerId": "",
							"uri": "",
							"customData": ""
						}],
						"coinTransfers": [
								{
									"to": "",
									"coins": [
										{
											"amount": "",
											"denom": "",
										}
									]
								}
							],
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"challengeTrackerId": "",
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"amountTrackerId": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"amountTrackerId": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireFromEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false,
						"overridesFromOutgoingApprovals": false,
						"overridesToIncomingApprovals": false
					}
				}
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			]
		}
	}`)

	return schemas
}
