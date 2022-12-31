package dialer

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_DIALER = "dialer"
)

func WithDialer(d Dialer) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_DIALER] = d
	}
}

func GetDialer(o option.Option) (Dialer, error) {
	i := o.Get(OPTION_DIALER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_DIALER)
	}

	v, ok := i.(Dialer)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}

	return v, nil
}
