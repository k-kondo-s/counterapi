#!/usr/bin/env bash

BASEDIR=$(dirname $0)

usage() {
cat << EOF

    Setup Counter API

    Usage:
        $0 APP_NUM # start apps
        $0 stop    # stop apps

    Example:
        # Start 5 apps
        $0 5

        # Stop
        $0 stop

EOF
}

start() {

    app_num=$1
    proxy_num=${app_num}

    # Even if app count 0, Nginx proxy settings needs at least one, so I do the following.
    if [[ app_num == 0 ]]; then
        proxy_num=1
    fi

    # Convert the sequence like "1 2 3" to "1,2,3"
    backends=""
    for i in `seq 1 ${proxy_num}`; do
        backends=${backends},${i[@]}
    done
    backends=${backends#,}

    cd $BASEDIR/

    # Run docker containers
    # First, create volume used by Nginx
    docker volume create scripts_vol

    # Run ansible container and create the valid proxy settings of Nginx which is according to the num of apps
    docker run --rm --mount source=scripts_vol,target=/etc/nginx/conf.d ansible ansible-playbook localhost.yml -e backends=${backends}

    # Run all applications with the specified number of app
    docker-compose up -d --scale app=${app_num}

    # The number of app can be dynamically changed, so always reload Nginx when the settings changes
    docker exec scripts_rp_1 nginx -s reload

    # Show the result
    docker-compose ps

}

stop() {
    cd $BASEDIR

    # Stop all containers
    docker-compose down

    # Delete volume
    docker volume rm scripts_vol
}

# Validate args num briefly
if [[ $# != 1 ]]; then
    usage
    exit 1
fi

ARG=$1

case $ARG in
    "stop") stop ;;
    *) start ${ARG} ;;
esac