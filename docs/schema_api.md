## 区块链基础数据检索API

[添加索引](#添加索引)，[索引列表](#索引列表)，[启动索引](#启动索引)，[停止索引](#停止索引)，[删除索引](#删除索引)，[SQL方式查询](#SQL方式查询)，[Dgraph语句查询](#Dgraph语句查询)

### 添加索引

| Method | Url |
| ------ | ------ |
| POST |   /schema            |

测试用例：

```bash
$ curl localhost:8082/schema -XPOST -d $'
{
  "id": "teacher",
  "ledger": "j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC",
  "associate_account": "LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa",
  "content": "type teacher{ id(isPrimaryKey: Boolean = true):int name:string classes:[uid]}"
}
'
{"success":true,"data":null}
```
参数：
- `id`唯一标识
- `ledger` 账本`hash`
- `associate_account`关联的数据账户地址
- `content` 索引内容，内容格式参照[GraphQL](https://graphql.cn/learn/schema/)，目前字段类型支持：`string`/`int`(`id`)/`float`/`boolean`/`datetime`/`uid`/`[uids]`。必须要有一个`PrimaryKey`用于`Argus`逻辑处理。

`Argus`针对键值对索引会自动添加如下`Predicates`：

- 索引元数据关联关系 `<schema-name>-schema: uid @reverse .`

- 健`<schema-name>-key>: string @index(exact) .`

- 版本`<schema-name>-version>: int @index(int) .`

- 时间`<schema-name>-time>: int @index(int) .`

  

使用`JD Chain SDK`向账本`j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC`下数据账户`LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa`插入数据：
```java
txTemp.dataAccount(Bytes.fromBase58("LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa")).setJSON("teacher one", "{\"id\":\"1\", \"name\":\"teacher one\"}", -1);
txTemp.dataAccount(Bytes.fromBase58("LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa")).setJSON("teacher two", "{\"id\":\"2\", \"name\":\"teacher two\"}", -1);
```


### 索引列表

| Method | Url |
| ------ | ------ |
| GET |    /schema/list       |

测试用例：
```bash
$ curl http://localhost:8082/schema/list | json_pp 
{
   "data" : [
      {
         "Progress" : -1,
         "status" : 0,
         "schema" : {
            "ledger" : "j5ufkRQxKeN7VAwJzh1pBoZbUEsozLuSWnQNoBGuYBpgDC",
            "associate_account" : "LdeNwqJPPKjUiaKQJWcXxqLEtE7wkXFbN7oXa",
            "content" : "type teacher{ id(isPrimaryKey: Boolean = true):int name:string classes:[uid]}",
            "id" : "teacher-j5ufkR-LdeNwq"
         }
      }
   ],
   "success" : true
}
```

### 启动索引

| Method | Url |
| ------ | ------ |
| GET |    /schema/start/:id  |

测试用例：
```bash
$ curl http://localhost:8082/schema/start/teacher-j5ufkR-LdeNwq
{"success":true,"data":null}
```

### 停止索引

| Method | Url |
| ------ | ------ |
| GET |    /schema/stop/:id   |

测试用例：
```bash
$ curl http://localhost:8082/schema/stop/teacher-j5ufkR-LdeNwq
{"success":true,"data":null}
```

### 删除索引

| Method | Url |
| ------ | ------ |
| DELETE | /schema/:id        |

测试用例：
```bash
$ curl -X DELETE http://localhost:8082/schema/teacher-j5ufkR-LdeNwq
{"success":false,"error":{"errorCode":1,"errorMessage":"schema is running, stop first"}}
```
> 删除前必须先停掉索引

### SQL方式查询

| Method | Url |
| ------ | ------ |
| POST |   /schema/querysql   |

测试用例：
```bash
$ curl localhost:8082/schema/querysql -XPOST -d $'
select * from teacher
' | json_pp
{
   "nodes" : [
      {
         "teacher-name" : "teacher one",
         "uid" : "0x3d",
         "teacher-id" : 1
      },
      {
         "uid" : "0x3e",
         "teacher-id" : 2,
         "teacher-name" : "teacher two"
      }
   ]
}
```

> 目前`SQL`的支持还很基础，许多语法都不支持......

### Dgraph语句查询

| Method | Url |
| ------ | ------ |
| POST |   /schema/query      |

测试用例：
```bash
$ curl localhost:8082/schema/query -XPOST -d $'
{
  teachers(func: has(teacher-name)) {
  	uid
	teacher-key
	teacher-version
	teacher-id
	teacher-name
	teacher-classes
  }
}' | json_pp
{
   "teachers" : [
      {
         "teacher-id" : 1,
         "uid" : "0x3d",
         "teacher-name" : "teacher one"
      },
      {
         "teacher-name" : "teacher two",
         "uid" : "0x3e",
         "teacher-id" : 2
      }
   ]
}
```

更多查询语法使用请参照[Dgraph](https://dgraph.io/docs/query-language/)