#! /bin/bash

if [ "x${MEEPO_ID}" == "x" ]; then
    echo "require MEEPO_ID"
    exit 1
fi

sudo docker exec -it -w /root meepo_${MEEPO_ID} /bin/sh
