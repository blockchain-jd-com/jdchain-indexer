#!/usr/bin/env bash

#jdchain's gateway corresponding version: 1.1.2.RELEASE;
DGRAPH_URL="127.0.0.1:9080"
ARGUS_HOST="11.7.178.75"
ARGUS_SEARCH_PORT="10001"
ARGUS_INDEXER_PORT="8082"
JDCHAIN_GW="http://jdchain2-18081.jd.com"
ARGUS_PATH="/jdchain/argus"

#first download the Dgraph version   : v1.0.16
nohup dgraph zero > dgraph_zero.log 2>&1 &
sleep 2
nohup dgraph alpha --lru_mb 1024 --zero localhost:5080 --log_dir dgraph_log > dgraph_alpha.log 2>&1 &
sleep 8
$ARGUS_PATH schema-update  --dgraph $DGRAPH_URL
sleep 2
nohup $ARGUS_PATH api-server --host $ARGUS_HOST --port $ARGUS_SEARCH_PORT --dgraph $DGRAPH_URL --production true> api-server.out 2>&1 &
sleep 2
nohup $ARGUS_PATH task  --dgraph $DGRAPH_URL > task_monitor.log 2>&1 &
sleep 2
nohup $ARGUS_PATH ledger-rdf  --api $JDCHAIN_GW --dgraph $DGRAPH_URL  --production true > converter2.out 2>&1 &
sleep 2
nohup $ARGUS_PATH data  --ledger-host $JDCHAIN_GW --port $ARGUS_INDEXER_PORT --dgraph $DGRAPH_URL --production true > value_indexer.out 2>&1 &
