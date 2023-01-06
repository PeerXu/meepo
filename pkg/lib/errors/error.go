package errors

import (
	"fmt"
)

func NewErrorAndErrorFunc[T any](errStr string, tpls ...string) (error, func(T) error) {
	err := fmt.Errorf(errStr)
	return err, NewErrorFunc[T](err, tpls...)
}

func NewErrorFunc[T any](err error, tpls ...string) func(T) error {
	var tpl string
	if len(tpls) > 0 {
		tpl = tpls[0]
	} else {
		tpl = "%w: %v"
	}

	return func(v T) error {
		return fmt.Errorf(tpl, err, v)
	}
}
