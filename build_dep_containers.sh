#!/bin/bash

VERSION=$1
if [ "${VERSION}" == "" ]
then
	VERSION="0.11.1"
fi

sudo docker build --build-arg VERSION=${VERSION} -t erikedin/librdkafka:${VERSION}  kubernetes/containers/librdkafka
