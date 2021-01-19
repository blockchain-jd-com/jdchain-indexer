package worker

import (
	"container/list"
	"fmt"
	"git.jd.com/jd-blockchain/explorer/value_indexer/schema"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestNewValueIndexWorker(t *testing.T) {
	fakeData := NewFakeDataForWorker(testDataSrc)
	schemaInfo, err := NewSchemaInfo(
		"schema-id",
		"ledger-id",
		"5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x",
		`    type Company{
        id(isPrimaryKey: Boolean = true):                   Int
        name:               String
    }`)
	assert.Nil(t, err)
	ns, err := schema.NewSchemaParser().FirstNodeSchema(schemaInfo.Content)
	assert.Nil(t, err)

	kvSchemaBuilder := NewKVSchemaBuilder(schemaInfo, ns)
	worker := NewValueIndexWorker("1", NewSchemaIndexStatusDefault(schemaInfo, "fakeuid"), ns, kvSchemaBuilder, fakeData)
	assert.NotNil(t, worker)

	worker.Start(fakeData)
	time.Sleep(time.Second)

	worker.Clear()

	time.Sleep(15 * time.Second)
}

type FakeDataForWorker struct {
	data *list.List
}

func NewFakeDataForWorker(src string) *FakeDataForWorker {
	source := &FakeDataForWorker{
		data: list.New(),
	}
	source.data.PushBack(src)
	return source
}

func (data *FakeDataForWorker) PushDelete(raw string) (e error) {
	fmt.Println("FakeDataForWorker Delete data:")
	fmt.Println(raw)
	return nil
}
func (data *FakeDataForWorker) PushUpdate(raw string) (e error) {
	fmt.Println("FakeDataForWorker Push data:")
	fmt.Println(raw)
	return nil
}

func (data *FakeDataForWorker) UIDs(predict string) ([]string, error) {
	fmt.Println("uids: ", predict)
	return []string{"1", "2", "3"}, nil
}

func (data *FakeDataForWorker) Stop() {
	fmt.Println("stopped")
}

func (data *FakeDataForWorker) Read() (string, int, error) {
	if data.data.Len() <= 0 {
		return "", 1, io.EOF
	}
	ele := data.data.Front()
	data.data.Remove(ele)
	return ele.Value.(string), 0, nil
}

var (
	testDataSrc = `
    {
    "data": [
        {
            "blockHeight": 1306,
            "executionState": "SUCCESS",
            "transactionContent": {
                "ledgerHash": {
                    "value": "6EEnjgF7vjLbZ49aY9Fg8FezqgZr9r5yT79egJ6eKYUy3"
                },
                "operations": [
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "KDUT3vaaSh8W1WTZ5bDEXEy4pJxRxtbDBP2RrFC2wAJz6J23DJHS"
									}
                                },
                                "key": "2"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "A56PGsTddBhk4Nb3ugY1BPRyhwdkSQPrkPpuiLjUxkQwGRtmaiKqioHRnG4"
									}
                                },
                                "key": "130"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "43i7caezNhTupQyg892B3VQfkzrEES72PttENHFKrbKxsCLyGCM4L7BaCFrVNc"
									}
                                },
                                "key": "19936"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    }
                ],
                "hash": {
                    "value": "64iNEdeFLWi8ZDzHGRX8YsXe3Mp2FD4Cj8zqmMyhJbYW3"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "44kHXr44gvK36pbNnDE4Pk26HLzgcnbi3oeZn62vwMZSjLHTNZ3ffUxjSfeXkoukPxysj3rcpWSU3WDUBvJSYMTko"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ],
            "hash": {
                "value": "69rkDKG1T4Z29CwDbovvqPyckUT4hNZAkkLiZ8bwtEnxs"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "41DZat4Vx3gFXLpCQPq8SU7ZX2Fah2TQ4fv6YxFyiU4whjkJ84vfkN8XdxD2j7ujNMcdAosUixvDWkrE4puXe5sX9"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ]
        },
        {
            "blockHeight": 1306,
            "executionState": "SUCCESS",
            "transactionContent": {
                "ledgerHash": {
                    "value": "6EEnjgF7vjLbZ49aY9Fg8FezqgZr9r5yT79egJ6eKYUy3"
                },
                "operations": [
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "A56PGsTddDghrMJGa1udLZu1GypiRYPtFK5RmRjcKkGJn1nriKeXL45URgp"
									}
                                },
                                "key": "289"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "27xvX7gbYZjUThJiNJ9T6ztMr4MD9ucCkTinkCAM8yMfXqLbifu2wnecFVXbpaTrVzLnAXK9LL4hG3nC"
									}
                                },
                                "key": "306"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "KDUT3vaaSq4TDnocxzTYWswUd6u7Z2Rqn5NydjoDAXciGZfhqTet"
									}
                                },
                                "key": "444"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "h34mK3iYfJhA5Qn9MTHx9NvvFostmhPQ4GjU48Go28BAMaqPzHUrWcDsvLtU"
									}
                                },
                                "key": "574"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    }
                ],
                "hash": {
                    "value": "68mxHqbKpkVAMQpnkq6zDfX37J4Ts7YxPCJZbcSQJrV9G"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "44gWVGkr9FZUZoCRx7szExMCSKm4qnpFT7QJaVSXhZTtg1zPcK5YnQEEH5dGhexFoggfDxZ5bHxZJQroa52QaVSeE"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ],
            "hash": {
                "value": "6M5KSaGYGavh95Nsv67TLkTBAFHfT8j7zFsEGLzNfDfA3"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "45j9M3b2HNY7wzTDTbKt9QXuYKd6jAfimDK9Vmw1FBSSQWMtU2bZabcu1tkCaH7K1awJAwC36AyLmGE8XxejYwrkE"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ]
        },
        {
            "blockHeight": 1306,
            "executionState": "SUCCESS",
            "transactionContent": {
                "ledgerHash": {
                    "value": "6EEnjgF7vjLbZ49aY9Fg8FezqgZr9r5yT79egJ6eKYUy3"
                },
                "operations": [
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "DHcFjJbV4EfhxWe1jGphtXcbuuXGm1T2GVfcrsUZVMd1oJLU"
									}
                                },
                                "key": "5"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "2GTQRVBMZqKzJVKNJSP8swuMYGfnZu3gt8bcwtY"
									}
                                },
                                "key": "10761"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    },
                    {
                        "writeSet": [
                            {
                                "expectedVersion": -1,
                                "value": {
                                    "type": "JSON",
									"bytes": {
										"value": "rYi9LarpdV2gG4KJM5iZNpYtSEVDcGXas6"
									}
                                },
                                "key": "69434"
                            }
                        ],
                        "accountAddress": {
                            "value": "5SmM4MrmpayoatEtNdYdY8J9u6sBeuGzJu1x"
                        }
                    }
                ],
                "hash": {
                    "value": "66NnkGNgJASFaEw9mVv8gJsqfzwXLHMxjoVQvRbTRhHT5"
                }
            },
            "endpointSignatures": [
                {
                    "digest": {
                        "value": "45RktPZUP99MrpM6kZQ7iZtHo5zAHbEM7juRG7GZCkCmDos1RAK3bj4iQV2TxPCCR2M8YzB5UdJpK3SfLiwk2VMNn"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ],
            "hash": {
                "value": "6ArGpXo7oc2xr5TP5smhdbWY4f9xvvyLyk4c4GcyjkFuM"
            },
            "nodeSignatures": [
                {
                    "digest": {
                        "value": "42nur4pHbJWMgh9DQkxcWCzTLoULGwHoomaMzPTGLcoadwc48KA8aWkqgEP2rVQUFQNww9Sj9B5wKDLxk3ueEdArV"
                    },
                    "pubKey": {
                        "value": "mb3m2B5bBULSwJKmX4hxVK7Sjikyk7e8v4WaYJkjpT9jRk"
                    }
                }
            ]
        }
    ],
    "success": true
}
    `
)
