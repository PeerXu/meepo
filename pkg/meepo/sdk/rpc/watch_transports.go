package sdk_rpc

import (
	"context"
	"encoding/json"

	"github.com/PeerXu/meepo/pkg/lib/rand"
	sdk_core "github.com/PeerXu/meepo/pkg/meepo/sdk/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (s *RPCSDK) WatchTransports() (<-chan sdk_interface.TransportView, <-chan error, func(), error) {
	stm, err := s.caller.CallStream(s.context(), sdk_core.METHOD_WATCH_EVENTS)
	if err != nil {
		return nil, nil, nil, err
	}

	cmd := sdk_interface.WatchEventsStream_WatchCommand{
		Command:  "watch",
		Session:  rand.DefaultStringGenerator.Generate(8),
		Policies: []string{"mpo.transport.state.*"},
	}

	cmdMsg, err := stm.Marshaler().Marshal(cmd)
	if err != nil {
		return nil, nil, nil, err
	}

	err = stm.SendMessage(cmdMsg)
	if err != nil {
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithCancel(s.context())
	tvs := make(chan sdk_interface.TransportView)
	errs := make(chan error)

	// TODO: log
	go func() {
		for {
			evtMsg, err := stm.RecvMessage()
			if err != nil {
				errs <- err
				return
			}
			var evt sdk_interface.WatchEventsStream_Event
			if err = evtMsg.Unmarshal(&evt); err != nil {
				errs <- err
				return
			}

			buf, err := json.Marshal(evt.Data)
			if err != nil {
				errs <- err
				return
			}

			var tv sdk_interface.TransportView
			if err = json.Unmarshal(buf, &tv); err != nil {
				errs <- err
				return
			}

			tvs <- tv
		}
	}()

	go func() {
		<-ctx.Done()
		close(tvs)
		close(errs)
		stm.Close()
	}()

	return tvs, errs, cancel, nil
}
