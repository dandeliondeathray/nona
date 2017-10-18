#!/bin/bash

function exit_on_failure() {
    if [ $? -ne 0 ]; then
        echo "FAILED TO BUILD!"
        exit $?
    fi
}

TAG=$1
if [ "${TAG}" == "" ]; then
    TAG="dev"
fi

#sudo docker build -t erikedin/nonainterface:${TAG} nonainterface
#exit_on_failure

sudo docker build -t erikedin/nonaplumber:${TAG} plumber
exit_on_failure
sudo docker build -t erikedin/nonaslackmessaging:${TAG} --build-arg PLUMBER_TAG=${TAG} slackmessaging
exit_on_failure
sudo docker build -t erikedin/nonapuzzlestore:${TAG} --build-arg PLUMBER_TAG=${TAG} puzzlestore
exit_on_failure
