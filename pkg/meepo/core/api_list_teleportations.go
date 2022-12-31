package meepo_core

import "context"

func (mp *Meepo) hdrAPIListTeleportations(ctx context.Context, _req any) (any, error) {
	tps, err := mp.ListTeleportations(ctx)
	if err != nil {
		return nil, err
	}

	return ViewTeleportations(tps), nil
}
