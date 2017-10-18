#!/bin/bash

#!/bin/bash

function exit_on_failure() {
    if [ $? -ne 0 ]; then
        echo "FAILED TO BUILD!"
        exit $?
    fi
}

WORKDIR=$(mktemp -d)
pushd ${WORKDIR}

git clone https://github.com/dandeliondeathray/nona
exit_on_failure

pushd nona

COMMIT_HASH=$(git rev-parse HEAD)
./build_containers.sh $COMMIT_HASH
exit_on_failure

for image in erikedin/nonaslackmessaging erikedin/nonapuzzlestore
do
    sudo docker tag ${image}:${COMMIT_HASH} ${image}:staging
done

popd # nona
popd # ${WORKDIR}
