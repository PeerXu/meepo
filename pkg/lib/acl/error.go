package acl

import (
	"fmt"

	"github.com/PeerXu/meepo/pkg/internal/errors"
)

var (
	ErrInvalidEntity, ErrInvalidEntityFn = errors.NewErrorAndErrorFunc[string]("invalid entity")
	ErrNotPermitted                      = fmt.Errorf("not permitted")
)
