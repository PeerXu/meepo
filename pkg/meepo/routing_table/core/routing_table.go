package meepo_routing_table_core

import (
	"github.com/PeerXu/meepo/pkg/internal/routing_table"
	meepo_routing_table_interface "github.com/PeerXu/meepo/pkg/meepo/routing_table/interface"
)

type routingTable struct {
	routing_table.RoutingTable
	greenLine int
}

func NewRoutingTable(rt routing_table.RoutingTable, greenLine int) meepo_routing_table_interface.RoutingTable {
	return &routingTable{
		RoutingTable: rt,
		greenLine:    greenLine,
	}
}

func (rt *routingTable) BucketHealthLevel(cpl int) meepo_routing_table_interface.HealthLevel {
	sz := rt.BucketSize(cpl)
	if sz == 0 {
		return meepo_routing_table_interface.HEALTH_LEVEL_RED
	} else if sz < rt.greenLine {
		return meepo_routing_table_interface.HEALTH_LEVEL_YELLOW
	} else {
		return meepo_routing_table_interface.HEALTH_LEVEL_GREEN
	}
}

func (rt *routingTable) HealthReport() *meepo_routing_table_interface.HealthReport {
	report := map[meepo_routing_table_interface.HealthLevel][]int{
		meepo_routing_table_interface.HEALTH_LEVEL_RED:    nil,
		meepo_routing_table_interface.HEALTH_LEVEL_YELLOW: nil,
		meepo_routing_table_interface.HEALTH_LEVEL_GREEN:  nil,
	}

	for i := 0; i < rt.TableSize(); i++ {
		lvl := rt.BucketHealthLevel(i)
		report[lvl] = append(report[lvl], i)
	}

	summary := map[meepo_routing_table_interface.HealthLevel]int{
		meepo_routing_table_interface.HEALTH_LEVEL_RED:    len(report[meepo_routing_table_interface.HEALTH_LEVEL_RED]),
		meepo_routing_table_interface.HEALTH_LEVEL_YELLOW: len(report[meepo_routing_table_interface.HEALTH_LEVEL_YELLOW]),
		meepo_routing_table_interface.HEALTH_LEVEL_GREEN:  len(report[meepo_routing_table_interface.HEALTH_LEVEL_GREEN]),
	}

	return &meepo_routing_table_interface.HealthReport{
		Summary: summary,
		Report:  report,
	}
}
