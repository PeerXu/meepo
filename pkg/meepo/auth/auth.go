package auth

import (
	"sync"
)

const (
	CONTEXT_NAME      = "name"
	CONTEXT_SIGNATURE = "signature"
)

type Context = map[string]interface{}

type AuthorizeOption = OFN

type Engine interface {
	Sign(payload Context) (Context, error)
	Verify(payload Context, signature Context) error
}

type NewEngineOption = OFN
type NewEngineFunc func(...NewEngineOption) (Engine, error)

var newEngineFuncs sync.Map

func RegisterNewEngineFunc(name string, fn NewEngineFunc) {
	newEngineFuncs.Store(name, fn)
}

func NewEngine(name string, opts ...NewEngineOption) (Engine, error) {
	fn, ok := newEngineFuncs.Load(name)
	if !ok {
		return nil, UnsupportedAuthEngineError
	}

	return fn.(NewEngineFunc)(opts...)
}
