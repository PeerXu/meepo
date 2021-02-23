package signaling

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type UserData = map[string]interface{}

type Descriptor struct {
	ID                 string                     `json:"id"`
	SessionDescription *webrtc.SessionDescription `json:"sessionDescription"`
	UserData           UserData                   `json:"userData,omitempty"`
}

type WireHandler func(*Descriptor) (*Descriptor, error)

type Engine interface {
	Wire(dst, src *Descriptor) (*Descriptor, error)
	OnWire(handler WireHandler)
	Close() error
}

type NewEngineFunc func(...NewEngineOption) (Engine, error)

var newEngineFuncs sync.Map

func RegisterNewEngineFunc(name string, fn NewEngineFunc) {
	newEngineFuncs.Store(name, fn)
}

func NewEngine(name string, opts ...NewEngineOption) (Engine, error) {
	fn, ok := newEngineFuncs.Load(name)
	if !ok {
		return nil, UnsupportedSignalingEngine(name)
	}

	return fn.(NewEngineFunc)(opts...)
}
