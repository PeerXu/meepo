package meepo_core

import (
	"context"
)

func (mp *Meepo) context() context.Context {
	return context.Background()
}
