package meepo_core

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/jinzhu/copier"

	lib_protocol "github.com/PeerXu/meepo/pkg/lib/protocol"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

func RequireProtocolVersion[IT, OT any](from string, to string) rpc_core.Middleware[IT, OT] {
	var requireFrom, requireTo *semver.Version
	var err error
	if from != "" {
		if requireFrom, err = lib_protocol.ParseProtocolVersion(from); err != nil {
			panic(err)
		}
	}
	if to != "" {
		if requireTo, err = lib_protocol.ParseProtocolVersion(to); err != nil {
			panic(err)
		}
	}

	return func(next func(context.Context, IT) (OT, error)) func(context.Context, IT) (OT, error) {
		return func(ctx context.Context, req IT) (res OT, err error) {
			var s struct{ Protocol string }
			if err = copier.Copy(&s, &req); err != nil {
				return
			}
			pvs := s.Protocol
			if pvs == "" {
				pvs = lib_protocol.UNKNOWN_VERSION_STRING
			}
			pv, err := lib_protocol.ParseProtocolVersion(pvs)
			if err != nil {
				return
			}

			if requireFrom != nil {
				if pv.Compare(requireFrom) < 0 {
					err = ErrRequireProtocolVersionBetweenFn(from, to)
					return
				}
			}

			if requireTo != nil {
				if pv.Compare(requireTo) > 0 {
					err = ErrRequireProtocolVersionBetweenFn(from, to)
					return
				}
			}

			return next(ctx, req)
		}
	}
}
