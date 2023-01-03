package option

import "github.com/stretchr/objx"

type Option = objx.Map

type ApplyOption = func(Option)

func NewOption(x ...map[string]any) Option {
	var m map[string]any
	if len(x) > 0 {
		m = x[0]
	} else {
		m = map[string]any{}
	}
	return objx.New(m)
}

func Apply(opts ...ApplyOption) Option {
	o := NewOption()
	return ApplyWithDefault(o, opts...)
}

func ApplyWithDefault(d Option, opts ...ApplyOption) Option {
	o := d.Copy()
	for _, apply := range opts {
		apply(o)
	}
	return o
}
