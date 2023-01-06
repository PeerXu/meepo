package meepo_routing_table_interface

import "github.com/PeerXu/meepo/pkg/lib/routing_table"

type HealthLevel string

func (x HealthLevel) String() string {
	return string(x)
}

const (
	HEALTH_LEVEL_RED    HealthLevel = "red"
	HEALTH_LEVEL_YELLOW HealthLevel = "yellow"
	HEALTH_LEVEL_GREEN  HealthLevel = "green"
)

type HealthReport struct {
	Summary map[HealthLevel]int
	Report  map[HealthLevel][]int
	Detials []map[string]any
}

type RoutingTable interface {
	routing_table.RoutingTable
	BucketHealthLevel(cpl int) HealthLevel
	HealthReport() *HealthReport
}
