#! /bin/bash

session="meepo-dev-stage-0-0"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo0"
tmux send-keys -t $session:$window "PIONS_LOG_TRACE=all dlv debug -r stdout:/tmp/meepo-dev-stage-0-0.log -r stderr:/tmp/meepo-dev-stage-0-0.log --output meepo-dev-stage-0.0 main.go -- serve -c etc/dev-stage-0.0.yaml" C-m
tmux split-window -v
tmux resize-pane -U 10
tmux send-keys -t $session:$window "watch ./meepo-dev-stage-0.0 -H :13345 t l" C-m
tmux split-window -v


session="meepo-dev-stage-0-1"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo1"
tmux send-keys -t $session:$window "PIONS_LOG_TRACE=all dlv debug -r stdout:/tmp/meepo-dev-stage-0-1.log -r stderr:/tmp/meepo-dev-stage-0-1.log --output meepo-dev-stage-0.1 main.go -- serve -c etc/dev-stage-0.1.yaml" C-m
tmux split-window -v
tmux resize-pane -U 10
tmux send-keys -t $session:$window "watch ./meepo-dev-stage-0.1 -H :14345 t l" C-m
tmux split-window -v


session="meepo-dev-stage-0-2"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo2"
tmux send-keys -t $session:$window "PIONS_LOG_TRACE=all dlv debug -r stdout:/tmp/meepo-dev-stage-0-2.log -r stderr:/tmp/meepo-dev-stage-0-2.log --output meepo-dev-stage-0.2 main.go -- serve -c etc/dev-stage-0.2.yaml" C-m
tmux split-window -v
tmux resize-pane -U 10
tmux send-keys -t $session:$window "watch ./meepo-dev-stage-0.2 -H :15345 t l" C-m
tmux split-window -v


session="meepo-dev-stage-0-3"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo3"
tmux send-keys -t $session:$window "PIONS_LOG_TRACE=all dlv debug -r stdout:/tmp/meepo-dev-stage-0-3.log -r stderr:/tmp/meepo-dev-stage-0-3.log --output meepo-dev-stage-0.3 main.go -- serve -c etc/dev-stage-0.3.yaml" C-m
tmux split-window -v
tmux resize-pane -U 10
tmux send-keys -t $session:$window "watch ./meepo-dev-stage-0.3 -H :16345 t l" C-m
tmux split-window -v
