package adaptor

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

var (
	getTxsInBlockResponseDemo = `
{
    "data": [
        {
            "blockHeight": 334,
            "userAccountSetHash": {
                "value": "6EUfvvoRf61b26ARtnFpDVGLiQJnfrGLUrgKgPEq4gYLn"
            },
            "contractAccountSetHash": {
                "value": "6EnfsNg2uqDFtpU5Kk5BzzCKZ1VcmCYQZKYxPxuWyj6EN"
            },
            "executionState": "SUCCESS",
            "transactionContent": {
                "ledgerHash": {
                    "value": "6Gw3cK4uazegy4HjoaM81ck9NgYLNoKyBMb7a1TK1jt3d"
                },
                "operations": [
                    {
                        "userID": {
                            "address": {
								"value":"5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6"
							}
                            "pubKey": {
                                "value": "mb5PUfd2HHBhNCmo6vuAzjKZdJT1yYZchwwbc3Bbs6tcgP"
                            }
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
									"bytes":{
										"value":"cGFyYW0xVmFs"
									}
								}
                                "key": "param1"
                            }
                        ],
                        "accountAddress": {
							"value":"5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6"
						}
                    }
                ],
                "hash": {
                    "value": "6LKEtG7JWUa2aN2s4pTymEkGzukR6K7coBaybEPQXGyvf"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "4495h3Kgs3yLMm2xKKrAJq7JdM9vjYRAJiGPsoEHzRREgkmSN7GjDUcYJMwXB5NDtVKQiczngYvbAesF9y9fods8i"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ],
            "adminAccountHash": {
                "value": "6899dc16DboTEnHhRnyCVYvxdtzqh45Ad1CL8rgW512ES"
            },
            "dataAccountSetHash": {
                "value": "6Kq7MUps67w8K3XrzdWUkkxLnqpxzmQGCbykJLrjhWLUV"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "45o5m7EiX6sw2vzykdwHUDeHNoWJBaPfPbXsxx9QhV3qVUSUvhbpu6fwFvN463vtYURkXD3stgXM7dPp6HS1aH8Ax"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ]
        },
        {
            "blockHeight": 334,
            "userAccountSetHash": {
                "value": "6EUfvvoRf61b26ARtnFpDVGLiQJnfrGLUrgKgPEq4gYLn"
            },
            "contractAccountSetHash": {
                "value": "6EnfsNg2uqDFtpU5Kk5BzzCKZ1VcmCYQZKYxPxuWyj6EN"
            },
            "executionState": "SUCCESS",
            "transactionContent": {
                "ledgerHash": {
                    "value": "6Gw3cK4uazegy4HjoaM81ck9NgYLNoKyBMb7a1TK1jt3d"
                },
                "operations": [
                    {
                        "userID": {
                            "address": "5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6",
                            "pubKey": {
                                "value": "mb5PUfd2HHBhNCmo6vuAzjKZdJT1yYZchwwbc3Bbs6tcgP"
                            }
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
									"bytes":{
										"value":"cGFyYW0xVmFs"
									}
								}
                                "key": "param1"
                            }
                        ],
                        "accountAddress": "5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6"
                    }
                ],
                "hash": {
                    "value": "6LKEtG7JWUa2aN2s4pTymEkGzukR6K7coBaybEPQXGyvf"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "4495h3Kgs3yLMm2xKKrAJq7JdM9vjYRAJiGPsoEHzRREgkmSN7GjDUcYJMwXB5NDtVKQiczngYvbAesF9y9fods8i"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ],
            "adminAccountHash": {
                "value": "6899dc16DboTEnHhRnyCVYvxdtzqh45Ad1CL8rgW512ES"
            },
            "dataAccountSetHash": {
                "value": "6Kq7MUps67w8K3XrzdWUkkxLnqpxzmQGCbykJLrjhWLUV"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "45o5m7EiX6sw2vzykdwHUDeHNoWJBaPfPbXsxx9QhV3qVUSUvhbpu6fwFvN463vtYURkXD3stgXM7dPp6HS1aH8Ax"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ]
        },
        {
          "blockHeight": 5,
          "userAccountSetHash": {
            "value": "6FWfm53R7CRqDWiQKN3X7ahq8FvtZoJDZhWVCD5wSxd1L"
          },
          "executionState": "SUCCESS",
          "transactionContent": {
            "ledgerHash": {
              "value": "6Gw3cK4uazegy4HjoaM81ck9NgYLNoKyBMb7a1TK1jt3d"
            },
            "operations": [
              {
                "accountID": {
                  "address": {
					"value": "5Sm2ZLRi8sK6FbpLrJzKnxH1mgkd9XR5eGhy"	
                  }
                  "pubKey": {
                    "value": "mbAyKDhT6KQPU6iheoKgWi5t7L6KGfFpLmKPwpnEdQFx4w"
                  }
                }
              }
            ],
            "hash": {
              "value": "67hMbfGzXCzDQrRaM7Eg5BfeiqnAWGxhfTTLdwq1T3wU9"
            }
          },
          "endpointSignatures": [
            {
              "digest": {
                "value": "43pU3LC8zVCka82jY2DVrBbnun1UrSgjmS2e7RF3iMGnBJ4WQsFjSi8q1ZVggNhB2Lhtukh2XP4aocMzCkFMB1FJS"
              },
              "pubKey": {
                "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
              }
            }
          ],
          "adminAccountHash": {
            "value": "6899dc16DboTEnHhRnyCVYvxdtzqh45Ad1CL8rgW512ES"
          },
          "dataAccountSetHash": {
            "value": "6M1tzVMqGwFfHVta3893TSNfDDEP9CMMEBLJDkXyPh3ob"
          },
          "nodeSignatures": [
            {
              "digest": {
                "value": "43xj1YDb2g7G42Kch9G6fp113VibzxmmLeKnibrhaL8QitSKmzUonuxAERnyk8Kh22L5WtUbryvaAJWvLVeX214xM"
              },
              "pubKey": {
                "value": "mb4Lri3hpQ8f9boJuQJAxWq1GdWFtAjxFT7qKaPyepxDWE"
              }
            }
          ]
        },
    {
        "blockHeight": 171,
        "userAccountSetHash": {
            "value": "6EZ1pzrMZQtwN83ZUHSpSP7epmM2SUkFBN5BiiJVzQeeG"
        },
        "contractAccountSetHash": {
            "value": "6AnYuTomcgLJigoms7uQFH3czQJztntXXY6y6ba6XKjuB"
        },
        "executionState": "SUCCESS",
        "transactionContent": {
            "ledgerHash": {
                    "value": "6Gw3cK4uazegy4HjoaM81ck9NgYLNoKyBMb7a1TK1jt3d"
                },
                "operations": [
                    {
                        "chainCode": "UEsDBBQACAgIACCck00AAAAAAAAAAAAAAAAJAAQATUVUQS1JTkYv/soAAAMAUEsHCAAAAAACAAAAAAAAAFBLAwQUAAgICAAgnJNNAAAAAAAAAAAAAAAAFAAAAE1FVEEtSU5GL01BTklGRVNULk1G803My0xLLS7RDUstKs7Mz7NSMNQz4OVyLkpNLElN0XWqBAlY6BnEG5qYKGj4FyUm56QqOOcXFeQXJZYA1WvycvFyAQBQSwcInnx2U0QAAABFAAAAUEsDBBQACAgIACCck00AAAAAAAAAAAAAAAATAAAAY29udHJhY3QucHJvcGVydGllc1OOKTU1MbCIKTV3dTQDkmZuQLapobkxkHQ2dY0pNTMwMuflUg5PTVFwSU1WMLQEIitjYysDQwXn4BAFIwNDC16u5Py8kqLE5BLb5PxcvawUvaSc/OTs5IzEzDw9mJSeY3FxaokzlGfMywUAUEsHCKuleilrAAAAegAAAFBLAQIUABQACAgIACCck00AAAAAAgAAAAAAAAAJAAQAAAAAAAAAAAAAAAAAAABNRVRBLUlORi/+ygAAUEsBAhQAFAAICAgAIJyTTZ58dlNEAAAARQAAABQAAAAAAAAAAAAAAAAAPQAAAE1FVEEtSU5GL01BTklGRVNULk1GUEsBAhQAFAAICAgAIJyTTauleilrAAAAegAAABMAAAAAAAAAAAAAAAAAwwAAAGNvbnRyYWN0LnByb3BlcnRpZXNQSwUGAAAAAAMAAwC+AAAAbwEAAAAA",
                        "contractID": {
                            "address": {
                                "value": "5SmDz45h8hh1agrkmFf4hn8VCHdTucSdK69y"
                            },
                            "pubKey": {
                                "value": "mawsK9U98hrTL9TgBFh5xUJNhjedYjeTHZuZFYXa5xuAg1"
                            }
                        }
                    }
                ],
                "hash": {
                    "value": "6M4Ykjge3CbHmMQmYFtfGcwZNam8UunFmYrEHnaEwAvTP"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "42ZCrKRLk6bL8V7bDhgpgydcGXkvjKAby87H3rvskdWexUN9pzJibZN19bkuQEn7uCQhHnz4RpaqFFwkPsSFrppJd"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ],
            "adminAccountHash": {
                "value": "6899dc16DboTEnHhRnyCVYvxdtzqh45Ad1CL8rgW512ES"
            },
            "dataAccountSetHash": {
                "value": "6ETzN7tQPcYcxdH2Wt3pVv5Ufn2p5pfMVRk9CjqV7HTCd"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "44cnHp4w5HMKCSmHuVB1Ez82Cat1PDSrdRznNq5wVDZnNwrkoDQX3q51pwbEWT5kGrVRDYTt8M7N7L71PuZGV2RLR"
                    },
                    "pubKey": {
                        "value": "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D"
                    }
                }
            ]
        }
    ],
    "success": true
}
`
)

func TestParseTransaction(t *testing.T) {
	result := gjson.Parse(getTxsInBlockResponseDemo)
	data := result.Get("data")
	txs := parseTransactions(data, 0, 0)
	assert.Len(t, txs, 4, spew.Sdump(data))
	tx := txs[0]
	assert.Equal(t, "6LKEtG7JWUa2aN2s4pTymEkGzukR6K7coBaybEPQXGyvf", tx.Hash)
	assert.Equal(t, int64(0), tx.IndexInBlock)
	//assert.Equal(t, "5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6", tx.Dataset)
	assert.Equal(t, "SUCCESS", tx.ExecutionState)
	assert.Len(t, tx.Contents, 2, spew.Sdump(data))
	assert.Len(t, tx.EndpointPublicKey, 1)
	assert.Equal(t, "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D", tx.EndpointPublicKey[0])
	assert.Len(t, tx.NodePublicKey, 1)
	assert.Equal(t, "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D", tx.NodePublicKey[0])

	opUser, ok := tx.Contents[0].(*UserOperation)
	assert.True(t, ok, spew.Sdump(opUser))
	assert.Equal(t, "mb5PUfd2HHBhNCmo6vuAzjKZdJT1yYZchwwbc3Bbs6tcgP", opUser.PublicKey, spew.Sdump(opUser, data))
	assert.Equal(t, "5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6", opUser.Address)

	opWrite, ok := tx.Contents[1].(*KVSetOperation)
	assert.True(t, ok, spew.Sdump(opWrite))
	assert.Equal(t, "5SmKyKcUhizt46BABySB8CxXxmixX81HfMP6", opWrite.DataSetAddress)
	assert.Len(t, opWrite.Histories, 1)
	historyHead := opWrite.Histories[0]
	assert.Equal(t, "param1", historyHead.Key)
	assert.Equal(t, "cGFyYW0xVmFs", historyHead.Value)
	assert.Equal(t, int64(-1), historyHead.Version)

	tx = txs[1]
	assert.Equal(t, int64(1), tx.IndexInBlock)

	tx = txs[2]
	assert.Equal(t, int64(2), tx.IndexInBlock)
	assert.Len(t, tx.Contents, 1)
	opAccount, ok := tx.Contents[0].(*DataAccountOperation)
	assert.True(t, ok, spew.Sdump(opAccount))
	assert.Equal(t, "5Sm2ZLRi8sK6FbpLrJzKnxH1mgkd9XR5eGhy", opAccount.Address)
	assert.Equal(t, "mbAyKDhT6KQPU6iheoKgWi5t7L6KGfFpLmKPwpnEdQFx4w", opAccount.PublicKey)

	tx = txs[3]
	assert.Equal(t, int64(3), tx.IndexInBlock)
	assert.Equal(t, "mb3iAYgF6bW8ohkTWYJWkLG1Tovb1oFErTCq2uUMfLAF2D", tx.NodePublicKey[0])
	assert.Len(t, tx.Contents, 1)
	opContract, ok := tx.Contents[0].(*ContractDeployOperation)
	assert.True(t, ok, spew.Sdump(opContract))
	assert.Equal(t, "5SmDz45h8hh1agrkmFf4hn8VCHdTucSdK69y", opContract.Address)
	assert.Equal(t, "mawsK9U98hrTL9TgBFh5xUJNhjedYjeTHZuZFYXa5xuAg1", opContract.PublicKey)
}
