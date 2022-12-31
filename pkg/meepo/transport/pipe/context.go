package transport_pipe

import "context"

func (p *PipeChannel) context() context.Context {
	return context.Background()
}
