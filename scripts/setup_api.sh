#!/usr/bin/env bash

usage() {
    echo test
}

n=$1
for i in `seq 1 $n`; do
  backends=$backends,${i[@]}
done
backends=${backends#,}
echo $backends

BASEDIR=$(dirname $0)
echo $BASEDIR
echo $0

cd $BASEDIR/
echo ########
pwd
echo ##########
docker volume create scripts_vol
docker run --rm --mount source=scripts_vol,target=/etc/nginx/conf.d ansible ansible-playbook localhost.yml -e backends=$backends
docker-compose up -d --scale app=$n
docker exec scripts_rp_1 nginx -s reload