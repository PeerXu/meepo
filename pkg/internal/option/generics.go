package option

func Setter[T any](k string) func(T) ApplyOption {
	return func(v T) ApplyOption {
		return func(o Option) {
			o[k] = v
		}
	}
}

func Getter[T any](k string) func(Option) (T, error) {
	return func(o Option) (T, error) {
		var x T
		i := o.Get(k).Inter()
		if i == nil {
			return x, ErrOptionRequiredFn(k)
		}

		v, ok := i.(T)
		if !ok {
			return x, ErrUnexpectedTypeFn(x, i)
		}

		return v, nil
	}
}

func New[T any](k string) (func(T) ApplyOption, func(Option) (T, error)) {
	return Setter[T](k), Getter[T](k)
}
