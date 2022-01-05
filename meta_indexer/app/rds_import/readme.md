## Import JD Chain to RDS


使用工具将`JD Chain`数据导入到关系型数据库中， 目前支持的数据库包括: `MySQL`

### 使用说明

1. 初始化数据库及表信息

使用提供的`jdchain.sql`脚本初始化数据库和表结构

示例： 初始化数据结构到`MySQL`中:

```sh
mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p $MYSQL_PASSWORD  < ./jdchain.sql
```

2. 执行命令导入`JD Chain`数据到数据库

```sh
./rds_import --ledger-host $API_HOST  --ledger $LEDGER --dsn $DATASOURCE_NAME --from $FROM --to $TO 
```

参数说明:

* `--ledger-host`:  网关地址。 默认值: `http://127.0.0.1:8080`
* `--ledger`:  账本HASH
* `--dsn`:  数据库数据源名称。 示例: *mysql*: `root:root@tcp(127.0.0.1:3306)/jdchain`
* `--from`:  导入区块起始高度。 默认值: `0`
* `--to`:  导入区块截止高度。 默认值: `-1`, 当前账本最新高度
* `--monitor`: 开启实时监控， 当有新区块产生时，实时导入数据。 目前该选项只能实时监听区块和交易数据
* `--refresh`: 开启实时监控， 刷新用户、数据账户、事件账户、合约等数据， 默认`10`, 单位分钟



