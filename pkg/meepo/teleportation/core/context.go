package teleportation_core

import "context"

func (tp *teleportation) context() context.Context {
	return context.Background()
}
