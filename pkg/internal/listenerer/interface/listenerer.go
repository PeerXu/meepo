package listenerer_interface

import "context"

type Listenerer interface {
	Listen(ctx context.Context, network, address string, opts ...ListenOption) (Listener, error)
}
