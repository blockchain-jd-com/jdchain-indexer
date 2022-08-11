# 使用说明

穿透式检索`jdchain-indexer`（`Argus`）提供`JD Chain`区块链基础数据索引、自定义键值索引服务。

版本对应关系：
|  jdchain-indexer（Argus）   | JD Chain  |
|  ----  | ----  |
| 0.9.0  | 1.6.0 |

**编译此项目获取`Argus`可执行文件**

## 安装并启动 Dgraph

参照[`Dgraph`官方文档](https://dgraph.io/downloads)下载安装并启动`Dgraph`（要求版本>1.1.0）.

`docker`启动示例：
```bash
docker run -d --rm -it -p 8181:8080 -p 9080:9080 -p 8000:8000 dgraph/standalone:v20.03.0
```

## 数据检索系统的使用（Argus）

### 更新 Schema

首次启动Argus时需要更新Schema，之后重启时不需要

```bash
# 指明 Dgraph 网络地址
argus schema-update  --dgraph 127.0.0.1:9080
```

参数：
- `dgraph` `Dgraph`服务地址

`Argus`针对`JD Chain`基础数据已建立了一些默认索引，参照[基础数据索引](docs/default_schema.md)


### 启动Argus所有服务

```bash
argus --ledger-host 127.0.0.1:8080 --dgraph 127.0.0.1:9080 --production true 
```

参数:

- `dgraph` `dgraph`服务地址，默认值：`127.0.0.1:9080`
- `production` 是否生产模式，默认`false`
- `ledger-host` 区块链网关服务地址，例如：`http://127.0.0.1:8080`
- `api-host` 区块链基础数据检索服务绑定`host`，默认`0.0.0.0`
- `api-port` 区块链基础数据检索服务绑定端口，默认`10001`，对应网关`data.retrieval.url`配置项
- `schema-port` `Schema`服务端口，默认`8082`，对应网关`schema.retrieval.url`配置项

> 其中`api-host`/`argus api-port`与`argus api-server`命令中`host`/`port`同义
> 其中`schema-port`与`argus data`命令中`port`同义
> 其中`task-port`与`argus task`命令中`port`同义

执行上面命令会一键[启动区块链基础数据索引](#启动区块链基础数据索引)，[启动区块链基础数据索引检索服务](#启动区块链基础数据索引检索服务)，[启动Value索引服务](#启动Value索引服务)。

### 启动区块链基础数据索引

```bash
# 指明 区块链网关服务和 Dgraph 网络地址
argus ledger-rdf --ledger-host 127.0.0.1:8080 --dgraph 127.0.0.1:9080 --production true
```

参数：

- `ledger-host` `JD Chain`网关服务地址
- `dgraph` `Dgraph`服务地址
- `production`生产模式

`Argus`将会持续运行，当有新账本和新区块产生时，会自动创建索引

### 启动区块链基础数据索引检索服务

```bash
# 指明 API服务所在服务器地址和所要监听的端口，以及 Dgraph 网络地址
argus api-server --host 127.0.0.1 --port 10001 --dgraph 127.0.0.1:9080 --production true
```

参数：

- `host` 服务绑定`IP`
- `port` 服务绑定端口
- `dgraph` `Dgraph`服务地址
- `production`生产模式

> 对应网关`data.retrieval.url`配置项

提供的接口及参数请参照[账本基础数据检索API](docs/ledger_api.md)

### 启动Value索引服务

`Argus`将会持续运行，针对自定义`Schema`，会自动根据数据账户中键值数据创建对应索引

```bash
argus data --port 8082 --ledger-host http://127.0.0.1:8080 --dgraph 127.0.0.1:9080 --production true
```

参数：

- `port` 服务绑定端口
- `ledger-host` `JD Chain`网关服务地址
- `dgraph` `Dgraph`服务地址
- `production`生产模式

> 对应网关`schema.retrieval.url`配置项

提供的接口及参数请参照[Schema API](docs/schema_api.md)

### 移除索引数据

会将数据库中所有索引移除，慎用！

```bash
# 指明 Dgraph 网络地址
argus drop  --dgraph 127.0.0.1:9080
```

参数：

- `dgraph` `Dgraph`服务地址

### 导出数据到MySQL

使用工具将`JD Chain`数据[导入到关系型数据库](meta_indexer/app/rds_import/readme.md)中， 目前支持的数据库包括: `MySQL`

