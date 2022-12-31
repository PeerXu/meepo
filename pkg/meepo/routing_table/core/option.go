package meepo_routing_table_core

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_GREEN_LINE = "greenLine"
)

var (
	WithGreenLine, GetGreenLine = option.New[int](OPTION_GREEN_LINE)
)
