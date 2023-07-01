package config

import "time"

type Poof struct {
	Disable           bool          `yaml:"disable"`
	Interval          time.Duration `yaml:"interval"`
	RequestCandidates int           `yaml:"request_candidates"`
}
