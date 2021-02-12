package signaling

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type Signal struct {
	ICECandidates    []webrtc.ICECandidate   `json:"iceCandidates"`
	ICEParameters    webrtc.ICEParameters    `json:"iceParameters"`
	DTLSParameters   webrtc.DTLSParameters   `json:"dtlsParameters"`
	SCTPCapabilities webrtc.SCTPCapabilities `json:"sctpCapabilities"`
}

type UserData = map[string]interface{}

type Descriptor struct {
	ID       string   `json:"id"`
	Signal   *Signal  `json:"signal"`
	UserData UserData `json:"userData,omitempty"`
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
