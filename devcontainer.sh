#!/bin/bash

#Get repo absolute location for mounting into the container.
local_workdir=$(cd $(dirname $(dirname "${BASH_SOURCE[0]}")) >/dev/null 2>&1 && pwd)
printf $local_workdir
printf "\n"
printf $PWD
printf "\n"

main() {
# Working directory inside container
local container_workdir=/go/src/bisa-patungan
# Identify container name
local container_name=dev-container-bisa-patungan

docker run --rm -it \
  --name $container_name \
  --volume $local_workdir:$container_workdir \
  --workdir $container_workdir \
  golang
}

main

