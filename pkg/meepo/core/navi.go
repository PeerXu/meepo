package meepo_core

import "time"

type NaviRequest struct {
	Session   string
	Tracker   Addr
	Candidate Addr
	CreatedAt time.Time
}
