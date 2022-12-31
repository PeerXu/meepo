package sdk_rpc

import (
	"context"
)

func (s *RPCSDK) context() context.Context {
	return context.Background()
}
