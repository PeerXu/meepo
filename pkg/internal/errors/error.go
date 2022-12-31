package errors

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/internal/typer"
)

func NewErrorAndErrorFunc[T typer.Typer](errStr string, tpls ...string) (error, func(T) error) {
	err := fmt.Errorf(errStr)
	return err, NewErrorFunc[T](err, tpls...)
}

func NewErrorFunc[T typer.Typer](err error, tpls ...string) func(T) error {
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
