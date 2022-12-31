package meepo_core

import "context"

func (mp *Meepo) hdrAPIWhoami(ctx context.Context, _ any) (any, error) {
	return mp.Addr().String(), nil
}
