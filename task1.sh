#!/usr/bin/env bash

scripts/setup_api.sh 3

for i in `seq 1 100`; do curl -s ${NGINX_IP}; echo; done | sort | uniq
