#!/usr/bin/env bash

# Build script of Counter API

CURRENTDIR=$(dirname $0)
cd $CURRENTDIR
BASEDIR=$(pwd)

# Build Counter API app container
cd $BASEDIR
cd ../app
docker build -t counterapi .

# Build Ansible container
cd $BASEDIR
cd ../ansible
docker build -t ansible .
