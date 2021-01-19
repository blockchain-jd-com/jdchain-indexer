## 区块链基础数据检索API

### 全局搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/all/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 不能小于20个字符 |

测试用例：
```bash
$ curl http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/all/search?keyword=j5o2vtSv7t7MwNCkhWjDTNWGP5SkrwGhJcFG856ca56jsR | json_pp 
{
   "success" : true,
   "data" : {
      "txs" : [
         {
            "hash" : "j5o2vtSv7t7MwNCkhWjDTNWGP5SkrwGhJcFG856ca56jsR",
            "block_height" : 8,
            "execution_state" : "SUCCESS"
         }
      ]
   },
   "total" : 1
}
```

### 合约搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/contracts/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/contracts/search?keyword=LdeNyTUurinxBWqpvkEhmuYEgVNFxH48dJLP7&from=0&count=10' | json_pp 
{
   "success" : true,
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRGHN8yC5EUNTsQvJQPBUdiHZgDUWUgnNLr2UBQJjSkNLL"
         },
         "address" : {
            "value" : "LdeNyTUurinxBWqpvkEhmuYEgVNFxH48dJLP7"
         }
      }
   ]
}
```

### 合约数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/contracts/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/contracts/count/search?keyword=LdeNyTUurinxBWqpvkEhmuYEgVNFxH48dJLP7' | json_pp 
{
   "success" : true,
   "data" : 1
}
```

### 区块搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/blocks/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/blocks/search?keyword=j5u3rWDjDt4EhEg1VVpKWwXCNgHHJCyRtQRmksxVXGDA4u&from=0&count=10' | json_pp 
{
   "success" : true,
   "data" : [
      {
         "hash" : "j5u3rWDjDt4EhEg1VVpKWwXCNgHHJCyRtQRmksxVXGDA4u",
         "height" : 8
      }
   ]
}
```

### 区块数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/blocks/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/blocks/count/search?keyword=j5u3rWDjDt4EhEg1VVpKWwXCNgHHJCyRtQRmksxVXGDA4u' | json_pp 
{
   "success" : true,
   "data" : 1
}
```

### 交易搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/txs/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/txs/search?keyword=j5o&from=0&count=10' | json_pp
{
   "success" : true,
   "data" : [
      {
         "hash" : "j5o5EF9TVEtkUU9XSLsY7e372fvziHeh7HusnumDNB3scm",
         "block_height" : 1,
         "execution_state" : "SUCCESS"
      },
      {
         "block_height" : 2,
         "execution_state" : "SUCCESS",
         "hash" : "j5ogvqEJaPmaUbMjGcxd1Y6EjLExqrYv3oVQr1mSFG1Sk3"
      },
      {
         "hash" : "j5o2vtSv7t7MwNCkhWjDTNWGP5SkrwGhJcFG856ca56jsR",
         "block_height" : 8,
         "execution_state" : "SUCCESS"
      }
   ]
}
```

### 交易数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/txs/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/txs/count/search?keyword=j5o' | json_pp
{
   "data" : 3,
   "success" : true
}
```

### 终端账户-交易搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/users/txs/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/users/txs/search?keyword=Lde&from=0&count=10' | json_pp
{
   "data" : [
      {
         "block_height" : 1,
         "execution_state" : "SUCCESS",
         "hash" : "j5o5EF9TVEtkUU9XSLsY7e372fvziHeh7HusnumDNB3scm"
      },
      {
         "hash" : "j5ogvqEJaPmaUbMjGcxd1Y6EjLExqrYv3oVQr1mSFG1Sk3",
         "execution_state" : "SUCCESS",
         "block_height" : 2
      },
      {
         "hash" : "j5gZdjc4baXeXUu3nkzo2ymmaBwFsAUSwgwwDMejUUR1Xw",
         "execution_state" : "SUCCESS",
         "block_height" : 3
      },
      {
         "block_height" : 4,
         "execution_state" : "SUCCESS",
         "hash" : "j5vxG64SEVWQ1nq6XK4rzVCDneFrVgUjvuVT5AMFzmYnyg"
      },
      {
         "hash" : "j5moq9Xni7hzCCG1QUtzMhfU3oZBoNV8bqjp9bmoR5CG1S",
         "execution_state" : "SUCCESS",
         "block_height" : 5
      },
      {
         "hash" : "j5ja9myRCNoS2yY1fcUJJXNWKpxGxZputeMEESJ2X1knET",
         "execution_state" : "SUCCESS",
         "block_height" : 6
      },
      {
         "hash" : "j5t4DroKSBnYTTwiQuUM6PKUSVjDeBALGUKKfCWReUFtRZ",
         "block_height" : 7,
         "execution_state" : "SUCCESS"
      },
      {
         "block_height" : 8,
         "execution_state" : "SUCCESS",
         "hash" : "j5o2vtSv7t7MwNCkhWjDTNWGP5SkrwGhJcFG856ca56jsR"
      }
   ],
   "success" : true
}
```

### 终端账户-交易数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/users/txs/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/users/txs/count/search?keyword=Lde' | json_pp
{
   "data" : 8,
   "success" : true
}
```

### 用户账户搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/users/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/users/search?keyword=Lde&from=0&count=10' | json_pp
{
   "success" : true,
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRL5vD8fhSLvn2g89GHjRzRRb5CswCvBYoPfb6E2tndgWA"
         },
         "address" : {
            "value" : "LdeP2yzn1dwG7Y81TGiStGp89YftmgaErrz9o"
         }
      },
      {
         "address" : {
            "value" : "LdeNi6zA1fbXEX85TZxpP6DN9tfyW44S8sifn"
         },
         "pubKey" : {
            "value" : "7VeR8396dihD4eLf5hb63eerSSVJwHhHcLTKDaKjboVRyzNT"
         }
      },
      {
         "address" : {
            "value" : "LdeNrnHJm21DjHWyrcXNGenmbvqAn49Psxysg"
         },
         "pubKey" : {
            "value" : "7VeRAi81RAPzPt8Px6ZkjhYGvPdttcmTCA9Nywq3r8kxbMRn"
         }
      },
      {
         "pubKey" : {
            "value" : "7VeRF8fmiGj1S5BHmNMgkB1GzEBjc7HEkX6jvsp8B6P9cs2i"
         },
         "address" : {
            "value" : "LdeNfNXsECpRAJn7wHv4LBToJyNzTLVALnbwp"
         }
      },
      {
         "address" : {
            "value" : "LdeNj8ByBSm2sVHYnwtndPirFbSLNWTVWGq7g"
         },
         "pubKey" : {
            "value" : "7VeRH6qMJ8mQ9ywVBisu2e39JjyvVA5STKbB5EdxpNM6CeH5"
         }
      }
   ]
}
```

### 用户账户数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/users/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/users/count/search?keyword=Lde' | json_pp
{
   "success" : true,
   "data" : 5
}
```

### 数据账户搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/accounts/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/accounts/search?keyword=Lde' | json_pp
{
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRFGB8ysFtshcwv2sqarHnJwNvP3ienxfW1FAiNhfRBHUp"
         },
         "address" : {
            "value" : "LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa"
         }
      },
      {
         "pubKey" : {
            "value" : "7VeRCVvdYt2UYz1HSLtc1Ekkg1zUueXvXxriCmuHzQX7V3Ex"
         },
         "address" : {
            "value" : "LdeNoHDq8CyKBanF3kpNCZaEvHzUVPdcoTuks"
         }
      },
      {
         "address" : {
            "value" : "LdeNgiG1N74XkPXb6VKsN7LkkrJcsaHaQ76iB"
         },
         "pubKey" : {
            "value" : "7VeRJTi3sfomUFfYYJzSUPhmDYnnnLhfG97Z5X632fU1p7jQ"
         }
      }
   ],
   "success" : true
}
```

### 数据账户数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/accounts/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/accounts/count/search?keyword=Lde' | json_pp
{
   "success" : true,
   "data" : 5
}
```

### 事件账户搜索

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/eventAccounts/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string |  |
| from     | form |int    | >= 0 |
| count    | form |int    | > 0, < 1000|

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/eventAccounts/search?keyword=Lde' | json_pp
{
   "success" : true,
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRFjbocQzGubf7uQgBDXZCXW8W4QsAgS51ShsXFjVXV9CK"
         },
         "address" : {
            "value" : "LdeNujhAnx8hxca4venxoonh6Bchjiz8tfCAX"
         }
      }
   ]
}
```

### 事件账户数量

| Method | Url |
| ------ | ------ |
|GET |  /ledgers/:ledger/eventAccounts/count/search |

参数列表：
| 名称 | 位置 | 类型 | 限制 |
| ------ | ------ | ------ |
| ledger   | path |string | 可以是多个，半角字符,分隔 |
| keyword  | form |string | 为空时查找所有 |

测试用例：
```bash
$ curl 'http://localhost:10001/ledgers/j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC/eventAccounts/count/search?keyword=Lde' | json_pp
{
   "data" : 1,
   "success" : true
}
```

### 账本列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/ledgers |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/ledgers' | json_pp
{
   "success" : true,
   "data" : [
      {
         "height" : 9, // 最新高度
         "hash" : "j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC"
      }
   ]
}
```

### 区块分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/blocks/range |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/blocks/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&from=0&to=2' | json_pp
{
   "success" : true,
   "data" : [
      {
         "hash" : "j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC",
         "height" : 0
      },
      {
         "hash" : "j5r1A2cwNHEgCqLL87jFbjaNUHv2hVtfDaynArUvJBGGyL",
         "height" : 1
      },
      {
         "hash" : "j5i9BJY3WsaETaqGZsFDayNkLzCmCwVX5QHsPkpfSMV57G",
         "height" : 2
      }
   ]
}
```

### 交易分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/txs/range   |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| height   | form |int |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/txs/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&height=0&from=0&count=2' | json_pp
{
   "success" : true,
   "data" : [
      {
         "hash" : "j5sX9HuA2wsKJV3qCWCpmEnAvPMjJo5c5q3J8EgK1sASZn",
         "block_height" : 0,
         "execution_state" : "SUCCESS"
      }
   ]
}
```

### 用户账户分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/users/range |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/users/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&from=0&count=2' | json_pp
{
   "success" : true,
   "data" : [
      {
         "address" : {
            "value" : "LdeP2yzn1dwG7Y81TGiStGp89YftmgaErrz9o"
         },
         "pubKey" : {
            "value" : "7VeRL5vD8fhSLvn2g89GHjRzRRb5CswCvBYoPfb6E2tndgWA"
         }
      },
      {
         "pubKey" : {
            "value" : "7VeR8396dihD4eLf5hb63eerSSVJwHhHcLTKDaKjboVRyzNT"
         },
         "address" : {
            "value" : "LdeNi6zA1fbXEX85TZxpP6DN9tfyW44S8sifn"
         }
      }
   ]
}
```

### 合约分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/contracts/range |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/contracts/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&from=0&count=2' | json_pp
{
   "data" : [
      {
         "address" : {
            "value" : "LdeNyTUurinxBWqpvkEhmuYEgVNFxH48dJLP7"
         },
         "pubKey" : {
            "value" : "7VeRGHN8yC5EUNTsQvJQPBUdiHZgDUWUgnNLr2UBQJjSkNLL"
         }
      }
   ],
   "success" : true
}
```

### 数据账户分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/accounts/range |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/accounts/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&from=0&count=2' | json_pp
{
   "success" : true,
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRFGB8ysFtshcwv2sqarHnJwNvP3ienxfW1FAiNhfRBHUp"
         },
         "address" : {
            "value" : "LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa"
         }
      },
      {
         "address" : {
            "value" : "LdeNgiG1N74XkPXb6VKsN7LkkrJcsaHaQ76iB"
         },
         "pubKey" : {
            "value" : "7VeRJTi3sfomUFfYYJzSUPhmDYnnnLhfG97Z5X632fU1p7jQ"
         }
      }
   ]
}
```

### 事件账户分页列表

| Method | Url |
| ------ | ------ |
|GET |  /api/v1/query/eventAccounts/range |

参数列表：
| 名称 | 位置 | 类型 |
| ------ | ------ |
| ledgers   | form |string |
| from  | form |int |
| to     | form |int    |

测试用例：
```bash
$ curl 'http://localhost:10001/api/v1/query/eventAccounts/range?ledgers=j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC&from=0&count=2' | json_pp
{
   "success" : true,
   "data" : [
      {
         "pubKey" : {
            "value" : "7VeRFjbocQzGubf7uQgBDXZCXW8W4QsAgS51ShsXFjVXV9CK"
         },
         "address" : {
            "value" : "LdeNujhAnx8hxca4venxoonh6Bchjiz8tfCAX"
         }
      }
   ]
}
```