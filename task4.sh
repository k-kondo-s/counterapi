#!/usr/bin/env bash

scripts/setup_api.sh 0

curl ${NGINX_IP} -m 10

scripts/setup_api.sh 5

for i in `seq 1 100`; do curl -s ${NGINX_IP}; echo; done | sort | uniq