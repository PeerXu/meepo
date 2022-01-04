package ofn

import "github.com/stretchr/objx"

type Option = objx.Map

type OFN = func(objx.Map)

func NewOption(ms ...map[string]interface{}) Option {
	var m map[string]interface{}
	if len(ms) > 0 {
		m = ms[0]
	} else {
		m = map[string]interface{}{}
	}
	return objx.New(m)
}
