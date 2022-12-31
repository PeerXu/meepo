package meepo_core

import (
	"context"
)

func (mp *Meepo) hdrAPIListTransports(ctx context.Context, _ any) (any, error) {
	ts, err := mp.ListTransports(ctx)
	if err != nil {
		return nil, err
	}

	return ViewTransports(ts), nil
}
