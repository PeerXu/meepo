package context

import "context"

func Value[T any](ctx context.Context, key any) (val T, found bool) {
	v := ctx.Value(key)
	if v == nil {
		found = false
		return
	}

	return v.(T), true
}

func MustValue[T any](ctx context.Context, key any) (val T) {
	val, _ = Value[T](ctx, key)
	return
}
