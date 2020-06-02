#!/usr/bin/env bash

echo ### Create 100 (and more) Counters ###
for i in `seq 1 100`; do curl -XPOST -s $NGINX_IP/counter?to=1000; done

echo ### Delete all registered counters ###
ids=$(curl -s $NGINX_IP/counter | jq .ids[] -r)
for i in ${ids}; do curl -XPOST -s $NGINX_IP/counter/${i[@]}/stop; done

echo ### check registered counters ###
curl $NGINX_IP/counter