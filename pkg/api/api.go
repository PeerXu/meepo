package api

import (
	"context"
	"fmt"
	"sync"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Wait() error
}

type NewServerFunc func(...NewServerOption) (Server, error)

var newServerFuncs sync.Map

func RegisterNewServerFunc(name string, fn NewServerFunc) {
	newServerFuncs.Store(name, fn)
}

func NewServer(name string, opts ...NewServerOption) (Server, error) {
	fn, ok := newServerFuncs.Load(name)
	if !ok {
		return nil, fmt.Errorf("Unsupported server: %s", name)
	}

	return fn.(NewServerFunc)(opts...)
}
