#!/usr/bin/env bash

for i in $(curl -s $NGINX_IP/counter | jq .ids[] -r); do
    echo ${i[@]}
    curl $NGINX_IP/counter/${i[@]} -s
    echo
done