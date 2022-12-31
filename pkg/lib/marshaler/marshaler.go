package marshaler

import (
	"context"

	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
)

type key string

const (
	marshalerKey   key = "marshaler"
	unmarshalerKey key = "unmarshaler"
)

func ContextWithMarshaler(ctx context.Context, marshaler marshaler_interface.Marshaler) context.Context {
	return context.WithValue(ctx, marshalerKey, marshaler)
}

func ContextWithUnmarshaler(ctx context.Context, unmarshaler marshaler_interface.Unmarshaler) context.Context {
	return context.WithValue(ctx, unmarshalerKey, unmarshaler)
}

func ContextWithMarshalerAndUnmarshaler(ctx context.Context, marshaler marshaler_interface.Marshaler, unmarshaler marshaler_interface.Unmarshaler) context.Context {
	return ContextWithUnmarshaler(ContextWithMarshaler(ctx, marshaler), unmarshaler)
}

func Marshal(ctx context.Context, x any) ([]byte, error) {
	marshaler, ok := ctx.Value(marshalerKey).(marshaler_interface.Marshaler)
	if !ok {
		panic("require marshaler")
	}
	return marshaler.Marshal(x)
}

func Unmarshal(ctx context.Context, x []byte, y any) error {
	unmarshaler, ok := ctx.Value(unmarshalerKey).(marshaler_interface.Unmarshaler)
	if !ok {
		panic("require unmarshaler")
	}
	return unmarshaler.Unmarshal(x, y)
}
