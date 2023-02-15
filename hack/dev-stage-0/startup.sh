#! /bin/bash

session="meepo-dev-stage-0-0"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo0"
tmux send-keys -t $session:$window "dlv debug --output meepo-dev-stage-0.0 main.go -- serve -c etc/dev-stage-0.0.yaml" C-m

session="meepo-dev-stage-0-1"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo1"
tmux send-keys -t $session:$window "dlv debug --output meepo-dev-stage-0.1 main.go -- serve -c etc/dev-stage-0.1.yaml" C-m

session="meepo-dev-stage-0-2"
tmux new-session -d -s $session

window=0
tmux rename-window -t $session:$window "meepo2"
tmux send-keys -t $session:$window "dlv debug --output meepo-dev-stage-0.2 main.go -- serve -c etc/dev-stage-0.2.yaml" C-m
