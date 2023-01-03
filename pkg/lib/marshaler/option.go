package marshaler

import (
	"github.com/PeerXu/meepo/pkg/lib/option"
	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
)

const OPTION_MARSHALER = "marshaler"

func WithMarshaler(marshaler marshaler_interface.Marshaler) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_MARSHALER] = marshaler
	}
}

func GetMarshaler(o option.Option) (marshaler_interface.Marshaler, error) {
	i := o.Get(OPTION_MARSHALER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_MARSHALER)
	}

	v, ok := i.(marshaler_interface.Marshaler)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}

	return v, nil
}

const OPTION_UNMARSHALER = "unmarshaler"

func WithUnmarshaler(unmarshaler marshaler_interface.Unmarshaler) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_UNMARSHALER] = unmarshaler
	}
}

func GetUnmarshaler(o option.Option) (marshaler_interface.Unmarshaler, error) {
	i := o.Get(OPTION_UNMARSHALER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_UNMARSHALER)
	}

	v, ok := i.(marshaler_interface.Unmarshaler)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}

	return v, nil
}
