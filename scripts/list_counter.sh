#!/usr/bin/env bash

NGINX_IP=127.0.0.1

for i in `seq 1 100`; do
    curl -s -XPOST $NGINX_IP/counter?to=1000
    echo
done

for i in `curl -s $NGINX_IP/counter | jq .ids[] -r`; do
    curl -s $NGINX_IP/counter/${i[@]}
    echo
done