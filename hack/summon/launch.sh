#! /bin/bash

if [ "x${MEEPO_PREFIX}" == "x" ]; then
    MEEPO_PREFIX="a"
fi

if [ "x${MEEPO_CLUSTER_SIZE}" == "x" ]; then
    echo "require MEEPO_CLUSTER_SIZE"
    exit 1
fi

tmux new -d -s meepo
for idx in $(seq $((MEEPO_CLUSTER_SIZE-1))); do
    tmux splitw -t 0
    tmux selectl -t 0 tiled
done

for idx in $(seq 0 $((MEEPO_CLUSTER_SIZE-1))); do
    tmux selectp -t ${idx}
    tmux send -t meepo "MEEPO_ID=${MEEPO_PREFIX}${idx} ./start.sh" C-j
done

tmux attach -t meepo
