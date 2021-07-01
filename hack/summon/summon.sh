#! /bin/sh

export MEEPO=${MEEPO:-"/bin/meepo"}

if [ ! -f "/root/config.template.yaml" ]; then
    echo "require config.template.yaml"
    exit 1
fi

mkdir -p /etc/meepo
cp /root/config.template.yaml /etc/meepo/meepo.yaml

if [ "x${MEEPO_SIGNALING_URL}" != "x" ]; then
    ${MEEPO} config set signaling.url="${MEEPO_SIGNALING_URL}"
fi
if [ "x${MEEPO_AS_SIGNALING}" != "x" ]; then
    ${MEEPO} config set asSignaling="${MEEPO_AS_SIGNALING}"
fi

${MEEPO} serve --daemon=false --log-level=trace
