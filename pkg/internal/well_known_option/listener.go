package well_known_option

import (
	"net"

	"github.com/PeerXu/meepo/pkg/internal/option"
)

const (
	OPTION_LISTENER = "listener"
)

func WithListener(l net.Listener) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_LISTENER] = l
	}
}

func GetListener(o option.Option) (net.Listener, error) {
	var x net.Listener

	i := o.Get(OPTION_LISTENER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_LISTENER)
	}

	v, ok := i.(net.Listener)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(x, v)
	}

	return v, nil
}
