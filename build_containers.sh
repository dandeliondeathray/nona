#!/bin/bash

function exit_on_failure() {
    if [ $? -ne 0 ]; then
        echo "FAILED TO BUILD!"
        exit 1
    fi
}

TAG=$1
if [ "${TAG}" == "" ]; then
    TAG="dev"
fi

sudo docker build -t erikedin/nonainterface:${TAG} service/nonainterface
exit_on_failure

sudo docker build -t erikedin/nonaplumber:${TAG} service/plumber
exit_on_failure
sudo docker build -t erikedin/nonaslackmessaging:${TAG} --build-arg PLUMBER_TAG=${TAG} service/slackmessaging
exit_on_failure
sudo docker build -t erikedin/nonapuzzlestore:${TAG} --build-arg PLUMBER_TAG=${TAG} service/puzzlestore
exit_on_failure
sudo docker build -t erikedin/nonaslack:${TAG} --build-arg INTERFACE_TAG=${TAG} service/nonaslack
exit_on_failure

sudo docker build -t erikedin/nonastaging:${TAG} --build-arg INTERFACE_TAG=${TAG} test/nonastaging
exit_on_failure
