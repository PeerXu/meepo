#! /bin/bash

DOCKER=${DOCKER:-"$(which docker)"}

MEEPO=${MEEPO:-"$(which meepo)"}
MEEPO_CONFIG=${MEEPO_CONFIG:-"$(readlink -f ~/.meepo/config.yaml)"}
CONTAINER_NAME=$1

DOCKER_RUN_OPTS=""
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS --rm"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -it"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS --name ${CONTAINER_NAME}"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS --entrypoint /root/summon.sh"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -v ${MEEPO}:/bin/meepo"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -v ${MEEPO_CONFIG}:/root/config.template.yaml"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -v `pwd`/summon.sh:/root/summon.sh"
DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -e MEEPO_AS_SIGNALING=true"
if [ "x${MEEPO_SIGNALING_URL}" != "x" ]; then
    DOCKER_RUN_OPTS="$DOCKER_RUN_OPTS -e MEEPO_SIGNALING_URL=${MEEPO_SIGNALING_URL}"
fi

sudo ${DOCKER} run \
	${DOCKER_RUN_OPTS} \
	alpine
