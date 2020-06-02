#!/usr/bin/env bash

id=$(curl -XPOST ${NGINX_IP}/counter?to=1000 -s)
echo $id
for i in `seq 1 10`; do curl ${NGINX_IP}/counter/$(echo $id | jq .id -r); echo; sleep 1; done
