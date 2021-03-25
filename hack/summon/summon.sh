#! /bin/sh

export MEEPO=${MEEPO:-"/bin/meepo"}

if [ "x${MEEPO_ID}" == "x" ]; then
    echo "require MEEPO_ID"
    exit 1
fi

if [ ! -f "/root/config.template.yaml" ]; then
    echo "require config.template.yaml"
    exit 1
fi

mkdir -p /root/.meepo/
cp /root/config.template.yaml /root/.meepo/config.yaml

${MEEPO} config set id="${MEEPO_ID}"
if [ "x${MEEPO_SIGNALING_URL}" != "x" ]; then
    ${MEEPO} config set signaling.url="${MEEPO_SIGNALING_URL}"
fi
if [ "x${MEEPO_AS_SIGNALING}" != "x" ]; then
    ${MEEPO} config set asSignaling="${MEEPO_AS_SIGNALING}"
fi

${MEEPO} serve --daemon=false --log-level=trace
