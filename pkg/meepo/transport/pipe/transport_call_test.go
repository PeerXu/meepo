package transport_pipe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PeerXu/meepo/pkg/internal/logging"
	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/lock"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func TestTransportCall(t *testing.T) {
	logger, err := logging.NewLogger(logging.WithLevel("trace"))
	assert.Nil(t, err)

	tsp := &PipeTransport{
		logger:      logger,
		fnsMtx:      lock.NewLock(well_known_option.WithName("fnsMtx")),
		fns:         make(map[string]meepo_interface.HandleFunc),
		marshaler:   marshaler_json.Marshaler,
		unmarshaler: marshaler_json.Unmarshaler,
	}

	tsp.Handle("echo", func(ctx context.Context, req meepo_interface.HandleRequest) (meepo_interface.HandleResponse, error) {
		return req, nil
	})

	var res string
	err = tsp.Call(context.Background(), "echo", "hello, world", &res)
	assert.Nil(t, err)
	assert.Equal(t, "hello, world", res)
}
