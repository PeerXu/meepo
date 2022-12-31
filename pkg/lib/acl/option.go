package acl

import "github.com/PeerXu/meepo/pkg/internal/option"

const (
	OPTION_ACL = "acl"
)

var WithAcl, GetAcl = option.New[Acl](OPTION_ACL)
