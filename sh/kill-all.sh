#!/usr/bin/env bash
ps -ef|grep 'argus'|grep -v grep|cut -c 9-15|xargs kill -9
ps -ef|grep 'dgraph'|grep -v grep|cut -c 9-15|xargs kill -9
